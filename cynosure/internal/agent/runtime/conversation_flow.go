package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
)

const (
	assistantDeltaEvent       = "assistant_delta"
	reasoningDeltaEvent       = "reasoning_delta"
	toolCallStartEvent        = "tool_call_start"
	toolCallDoneEvent         = "tool_call_done"
	maxRound                  = 1000
	toolArgsPreviewMax        = 160
	toolResultPreviewMaxLines = 3

	// mainAgentTurnTimeout 限定单个主 Agent 回合的耗时上限，它是在各轮之间
	// 进行检查的软边界。
	mainAgentTurnTimeout = 24 * time.Hour

	// defaultMaxTokens 是单次请求的输出预算；在首次发生截断时会被升级为
	// truncationMaxTokens。
	defaultMaxTokens = 8000
	// truncationMaxTokens 是首次截断后升级的输出预算（为默认值的 8 倍）。
	truncationMaxTokens = 64 * 1024
	// maxResumeAttempts 表示在升级预算后仍发生截断时，最多发起的续写请求次数。
	maxResumeAttempts = 3

	// truncationResumePrompt 作为一条 user 消息注入，用于让模型继续输出被
	// token 上限截断的内容。
	truncationResumePrompt = `Output token limit hit. Resume directly — no apology, no recap of what you were doing. Pick up mid-thought if that is where the cut happened. Break remaining work into smaller pieces.`
)

type bufferedModelDelta struct {
	Event   string
	Content string
}

func (s *Service) RespondToConversation(ctx context.Context, conversation storage.Conversation, user storage.User, userMessage string, writer EventWriter) (storage.Message, error) {
	// 记忆业务开关：系统级能力开启 默认开启
	// 注意：它只控制记忆注入与提取，不控制会话锁与模型历史持久化。
	memoryOn := s.EnableMemory && user.MemoryEnabled
	// 入口获取会话锁：上一轮收尾未完成时阻塞等待，直到拿到锁或等待超时。
	// 获取失败（本地锁异常 / 等待超时）则降级放行，跳过本轮收尾。
	var lockToken string
	var stopRenew func()
	handedOff := false
	if s.EnableMemory {
		token := idgen.New("lock")
		ok, err := s.Store.AcquireConversationLock(ctx, conversation.ID, token, s.Cfg.ConversationLockTTL, s.Cfg.ConversationLockWaitTimeout)
		if err != nil {
			logger.Warn(fmt.Sprintf("conversation lock: acquire failed conversation=%s: %v", conversation.ID, err))
		} else if !ok {
			logger.Warn(fmt.Sprintf("conversation lock: acquire timed out conversation=%s", conversation.ID))
		} else {
			lockToken = token
			stopRenew = s.startLockRenewer(conversation.ID, token)
		}
	}
	// defer 兜底：任何提前 return（出错路径）都停止续期并释放锁；
	// 正常走到异步收尾时通过 handedOff 转交所有权，避免被误释放。
	defer func() {
		if handedOff || lockToken == "" {
			return
		}
		if stopRenew != nil {
			stopRenew()
		}
		_ = s.Store.ReleaseConversationLock(context.Background(), conversation.ID, lockToken)
	}()

	history, err := s.loadConversationMessages(ctx, conversation.ID)
	if err != nil {
		return storage.Message{}, err
	}
	modelHistory := s.loadModelHistory(ctx, conversation.ID, history)
	state := s.newLoopState(conversation, user, userMessage, history, modelHistory, writer)
	if err := s.hookManager().RunUserPromptSubmit(ctx, &UserPromptSubmitContext{State: state}); err != nil {
		return storage.Message{}, err
	}

	snapshot, err := s.buildSkillSnapshot(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	state.SkillSnapshot = snapshot
	state.SystemPrompt = s.buildSystemPrompt(ctx, conversation, user, snapshot, state.History, memoryOn)
	state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
	toolDefs := s.toolDefinitionsForUser(ctx, user.ID)
	round := 0
	roundsSinceTodoWrite := 0
	turnStart := time.Now()
	var cumulativeReasoning strings.Builder

	for {
		round++
		if round > maxRound {
			return storage.Message{}, fmt.Errorf("conversation round limit %d exceeded", maxRound)
		}
		if time.Since(turnStart) > mainAgentTurnTimeout {
			return storage.Message{}, fmt.Errorf("main agent turn timed out after %v", mainAgentTurnTimeout)
		}
		// 压缩结果写回唯一的真实消息历史 ModelHistory：内存态 = 发送态 = 落库态。
		modelHistory, err := s.compressContextBeforeLLM(ctx, state)
		if err != nil {
			return storage.Message{}, err
		}
		state.ModelHistory = modelHistory
		state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, s.Tools, roundsSinceTodoWrite)
		estimator := compression.DefaultTokenEstimator{}
		state.LastContextTokens = estimator.EstimateRequestTokens(state.SystemPrompt, state.ModelHistory, toolDefs)
		state.LastContextBudget = estimator.ContextTokenBudget()
		emitMeta(state)
		req := openai.ChatCompletionRequest{
			Model:     s.Cfg.LLM.ModelID,
			Messages:  state.Messages,
			Tools:     toolDefs,
			MaxTokens: defaultMaxTokens,
		}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundWithRecovery(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("main-agent conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
			return storage.Message{}, err
		}
		if err := ctx.Err(); err != nil {
			return storage.Message{}, err
		}
		cumulativeReasoning.WriteString(msg.ReasoningContent)
		requestMsg := msg
		state.Messages = append(state.Messages, requestMsg)
		if toolCallsInclude(msg.ToolCalls, agenttools.TodoWriteToolName) {
			roundsSinceTodoWrite = 0
		} else {
			roundsSinceTodoWrite++
		}

		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			// 会话结束：把本轮模型最终回复纳入真实消息历史后重算，得到会话结束时的最终
			// token 用量，覆盖请求前的估算值，确保存储与下发的都是最终用量。
			finalAssistant := storage.Message{Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
			finalHistoryForEstimate := append(cloneMessages(state.ModelHistory), finalAssistant)
			finalEstimator := compression.DefaultTokenEstimator{}
			state.LastContextTokens = finalEstimator.EstimateRequestTokens(state.SystemPrompt, finalHistoryForEstimate, toolDefs)
			state.LastContextBudget = finalEstimator.ContextTokenBudget()
			stopCtx := &StopContext{State: state, ModelMessage: msg, Content: fallbackAssistantContent(msg.Content), ReasoningContent: cumulativeReasoning.String()}
			if err := s.hookManager().RunStop(ctx, stopCtx); err != nil {
				return storage.Message{}, err
			}
			// 最终 assistant 追加进唯一的真实消息历史 ModelHistory（与 Stop 钩子把它追加进
			// 展示历史 History 对称）；用带 ID/Meta 的 stopCtx.AssistantMessage，使断点 ID
			// 与展示历史一致。落库与记忆提取统一以 ModelHistory（压缩后真实消息线）为源。
			state.ModelHistory = append(state.ModelHistory, stopCtx.AssistantMessage)
			if s.EnableMemory {
				// 需求2 条件(2)/初次提取：轮次自然结束时评估是否更新会话记忆。
				updateSession := false
				if memoryOn {
					updateSession = s.shouldUpdateSessionMemoryAtTurnEnd(conversation, state.LastContextTokens)
				}
				handedOff = s.scheduleMemoryWork(conversation, user, state.ModelHistory, lockToken, stopRenew, memoryOn, updateSession, state.LastContextTokens)
			}
			return stopCtx.AssistantMessage, nil
		}
		assistantMessage := storage.Message{ID: state.NextMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
		state.History = append(state.History, assistantMessage)
		state.ModelHistory = append(state.ModelHistory, assistantMessage)

		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			if approved, _ := s.approveToolCall(ctx, tc); !approved {
				// 用户拒绝：展示被拒状态后立即结束本轮，且不进行收尾（记忆/历史落库）。
				toolCtx.Outcome = toolExecutionOutcome{Status: "rejected", Result: "Error: user rejected this operation", Audit: toolCtx.Outcome.Audit}
				emitToolCallStart(state, toolCtx)
				emitToolCallDone(state, toolCtx)
				return storage.Message{Role: "assistant", Content: "操作已被拒绝，已结束本轮。"}, nil
			}
			emitToolCallStart(state, toolCtx)
			toolCtx.Outcome = s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Skills: snapshot, PersistedOutputReader: s.newPersistedOutputReader(conversation.ID, user.ID), Todos: state.Todos}, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			emitToolCallDone(state, toolCtx)
			state.ToolCallCount++
			emitMeta(state)
		}
		// 需求2 条件(1)/初次提取：每个 tool_calls round 结束后评估并按需异步刷新会话记忆。
		// 提取数据源用【真实模型线】state.ModelHistory（lockstep 追加，含工具调用与结果）。
		if memoryOn {
			s.maybeUpdateSessionMemoryMidLoop(conversation, user, state.ModelHistory, state.LastContextTokens, len(msg.ToolCalls))
		}
	}
}

// emitMeta 实时下发当前累计的回复元信息。
func emitMeta(state *LoopState) {
	if state == nil || state.Writer == nil {
		return
	}
	_ = state.Writer.Event("meta", map[string]any{
		"tool_call_count": state.ToolCallCount,
		"context_tokens":  state.LastContextTokens,
		"context_budget":  state.LastContextBudget,
	})
}

func emitToolCallStart(state *LoopState, toolCtx *ToolUseContext) {
	if state == nil || state.Writer == nil || toolCtx == nil {
		return
	}
	_ = state.Writer.Event(toolCallStartEvent, map[string]any{
		"tool_call_id": toolCtx.ToolCall.ID,
		"tool_name":    toolCtx.Name,
		"raw_args":     toolCtx.RawArgs,
		"args_preview": previewToolArgs(toolCtx.Name, toolCtx.RawArgs),
		"status":       "running",
	})
}

func emitToolCallDone(state *LoopState, toolCtx *ToolUseContext) {
	if state == nil || state.Writer == nil || toolCtx == nil {
		return
	}
	_ = state.Writer.Event(toolCallDoneEvent, map[string]any{
		"tool_call_id":   toolCtx.ToolCall.ID,
		"tool_name":      toolCtx.Name,
		"raw_args":       toolCtx.RawArgs,
		"args_preview":   previewToolArgs(toolCtx.Name, toolCtx.RawArgs),
		"status":         toolCtx.Outcome.Status,
		"result_preview": previewToolResult(toolCtx.Outcome),
		"audit_summary":  toolCtx.Outcome.AuditSummary(),
	})
}

func previewToolArgs(toolName, rawArgs string) string {
	trimmed := strings.TrimSpace(rawArgs)
	if trimmed == "" {
		return ""
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(trimmed), &args); err != nil {
		return truncatePreview(singleLine(trimmed), toolArgsPreviewMax)
	}
	for _, key := range preferredToolArgKeys(toolName) {
		if value, ok := args[key]; ok {
			return truncatePreview(fmt.Sprintf("%s: %s", key, singleLine(fmt.Sprint(value))), toolArgsPreviewMax)
		}
	}
	return truncatePreview(singleLine(stableJSON(args)), toolArgsPreviewMax)
}

func preferredToolArgKeys(toolName string) []string {
	if toolName == "bash" {
		return []string{"command", "description"}
	}
	return []string{"file_path", "path", "name", "query", "command"}
}

func previewToolResult(outcome toolExecutionOutcome) string {
	preview := outcome.Result
	if strings.TrimSpace(preview) == "" {
		preview = outcome.Audit.OutcomeSummary
	}
	return truncatePreviewLines(preview, toolResultPreviewMaxLines)
}

func stableJSON(args map[string]any) string {
	keys := make([]string, 0, len(args))
	for key := range args {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%q:%q", key, fmt.Sprint(args[key])))
	}
	return "{" + strings.Join(parts, ",") + "}"
}

func singleLine(text string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(text)), " ")
}

func limitPreviewLines(text string, maxLines int) string {
	lines := strings.Split(strings.TrimSpace(text), "\n")
	if len(lines) <= maxLines {
		return strings.Join(lines, "\n")
	}
	return strings.Join(lines[:maxLines], "\n") + " …"
}

func truncatePreviewLines(text string, maxLines int) string {
	text = strings.TrimSpace(text)
	if text == "" || maxLines <= 0 {
		return ""
	}
	lines := strings.Split(text, "\n")
	if len(lines) <= maxLines {
		return strings.Join(lines, "\n")
	}
	omitted := len(lines) - maxLines
	kept := append([]string(nil), lines[:maxLines]...)
	kept = append(kept, fmt.Sprintf("... + %d lines", omitted))
	return strings.Join(kept, "\n")
}

func truncatePreview(text string, maxLen int) string {
	runes := []rune(strings.TrimSpace(text))
	if maxLen <= 0 || len(runes) <= maxLen {
		return string(runes)
	}
	if maxLen == 1 {
		return "…"
	}
	return string(runes[:maxLen-1]) + "…"
}

// runModelRoundWithRecovery 在 runModelRoundStream 之上封装了截断与上下文溢出
// 的恢复逻辑。瞬时的 429/529 重试在 LLM 客户端内部处理，对这里是透明的。
//
//   - 截断（finish_reason == length）：首次发生时，把输出预算升级为
//     truncationMaxTokens，并在不改动 messages 的情况下重试同一请求（丢弃
//     已产生的部分输出）。若仍然截断，则最多发起 maxResumeAttempts 次续写
//     请求，每次都追加被截断的文本以及一条续写提示；各分段以流式方式输出
//     并拼接为一条 assistant 消息。
//   - 上下文溢出（HTTP 413）：执行一次 reactiveCompact，从压缩后的历史重建
//     请求并重试。若再次溢出，则返回给调用方（由既有的兜底边界处理）。
func (s *Service) runModelRoundWithRecovery(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, error) {
	cur := req
	upgraded := false
	compacted := false
	resumeAttempts := 0
	var accumulated strings.Builder
	var accumulatedReasoning strings.Builder

	for {
		msg, finishReason, deltas, err := s.runModelRoundStream(ctx, state, cur)
		if err != nil {
			if llm.IsContextOverflow(err) && !compacted {
				if compactErr := s.reactiveCompact(ctx, state); compactErr != nil {
					logger.Warn(fmt.Sprintf("reactive compact failed conversation=%s: %v", state.Conversation.ID, compactErr))
					return openai.ChatCompletionMessage{}, "", err
				}
				compacted = true
				cur.Messages = state.Messages
				cur.MaxTokens = defaultMaxTokens
				upgraded = false
				resumeAttempts = 0
				accumulated.Reset()
				accumulatedReasoning.Reset()
				continue
			}
			return openai.ChatCompletionMessage{}, "", err
		}

		resuming := accumulated.Len() > 0
		if finishReason != openai.FinishReasonLength {
			if resuming {
				flushContentDeltas(state, deltas)
				return openai.ChatCompletionMessage{
					Role:             "assistant",
					Content:          accumulated.String() + msg.Content,
					ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
					ToolCalls:        msg.ToolCalls,
				}, finishReason, nil
			}
			if state.Writer != nil && shouldEmitAssistantContentDeltas(finishReason, msg.ToolCalls) {
				flushContentDeltas(state, deltas)
			}
			return msg, finishReason, nil
		}

		// 输出被截断。
		if !upgraded {
			// 丢弃部分输出，升级预算，原样重试。
			upgraded = true
			cur.MaxTokens = truncationMaxTokens
			continue
		}
		if resumeAttempts >= maxResumeAttempts {
			// 尽力而为：输出并返回目前已获得的内容。
			flushContentDeltas(state, deltas)
			logger.Warn(fmt.Sprintf("truncation resume exhausted conversation=%s after %d attempts", state.Conversation.ID, resumeAttempts))
			return openai.ChatCompletionMessage{
				Role:             "assistant",
				Content:          accumulated.String() + msg.Content,
				ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
			}, finishReason, nil
		}
		// 保留这一段部分输出，并要求模型继续。
		flushContentDeltas(state, deltas)
		accumulated.WriteString(msg.Content)
		accumulatedReasoning.WriteString(msg.ReasoningContent)
		resumeAttempts++
		cur.Messages = append(cur.Messages,
			openai.ChatCompletionMessage{Role: "assistant", Content: msg.Content},
			openai.ChatCompletionMessage{Role: "user", Content: truncationResumePrompt},
		)
	}
}

func flushContentDeltas(state *LoopState, deltas []bufferedModelDelta) {
	if state.Writer == nil {
		return
	}
	for _, delta := range deltas {
		_ = state.Writer.Event(delta.Event, map[string]any{"content": delta.Content})
	}
}

func (s *Service) runModelRoundStream(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, []bufferedModelDelta, error) {
	stream, err := s.LLM.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, "", nil, err
	}
	defer stream.Close()

	var content strings.Builder
	var reasoningContent strings.Builder
	var finishReason openai.FinishReason
	toolCalls := &streamedToolCallAccumulator{}
	var bufferedContentDeltas []bufferedModelDelta
	seenChoice := false
	seenOutput := false

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return openai.ChatCompletionMessage{}, "", nil, err
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		seenChoice = true
		choice := chunk.Choices[0]
		if choice.Delta.Content != "" {
			seenOutput = true
			content.WriteString(choice.Delta.Content)
			bufferedContentDeltas = append(bufferedContentDeltas, bufferedModelDelta{Event: assistantDeltaEvent, Content: choice.Delta.Content})
		}
		if choice.Delta.ReasoningContent != "" {
			seenOutput = true
			reasoningContent.WriteString(choice.Delta.ReasoningContent)
			if state.Writer != nil {
				_ = state.Writer.Event(reasoningDeltaEvent, map[string]any{"content": choice.Delta.ReasoningContent})
			}
		}
		if len(choice.Delta.ToolCalls) > 0 {
			seenOutput = true
			toolCalls.Add(choice.Delta.ToolCalls)
		}
		if choice.FinishReason != "" {
			finishReason = choice.FinishReason
		}
	}
	if !seenChoice || (!seenOutput && finishReason == "") {
		return openai.ChatCompletionMessage{}, "", nil, fmt.Errorf("model stream returned no choices")
	}
	calls := toolCalls.Calls()

	return openai.ChatCompletionMessage{Role: "assistant", Content: content.String(), ReasoningContent: reasoningContent.String(), ToolCalls: calls}, finishReason, bufferedContentDeltas, nil
}

func shouldEmitAssistantContentDeltas(finishReason openai.FinishReason, toolCalls []openai.ToolCall) bool {
	return finishReason != openai.FinishReasonToolCalls && len(toolCalls) == 0
}

type streamedToolCallAccumulator struct {
	calls []openai.ToolCall
}

func (a *streamedToolCallAccumulator) Add(deltas []openai.ToolCall) {
	for deltaPosition, delta := range deltas {
		index := len(a.calls) - 1
		if delta.Index != nil {
			index = *delta.Index
		} else if len(deltas) > 1 {
			index = deltaPosition
		}
		if index < 0 {
			index = 0
		}
		for len(a.calls) <= index {
			a.calls = append(a.calls, openai.ToolCall{})
		}

		call := a.calls[index]
		if delta.ID != "" {
			call.ID = delta.ID
		}
		if delta.Type != "" {
			call.Type = delta.Type
		}
		if delta.Function.Name != "" {
			call.Function.Name = delta.Function.Name
		}
		if delta.Function.Arguments != "" {
			call.Function.Arguments += delta.Function.Arguments
		}
		a.calls[index] = call
	}
}

func (a *streamedToolCallAccumulator) Calls() []openai.ToolCall {
	if len(a.calls) == 0 {
		return nil
	}
	return append([]openai.ToolCall(nil), a.calls...)
}

func (s *Service) loadConversationMessages(ctx context.Context, conversationID string) ([]storage.Message, error) {
	if cached, ok, err := s.Store.GetConversationCache(ctx, conversationID); err == nil && ok {
		return cached, nil
	}
	messages, err := s.Store.ListMessagesByConversation(ctx, conversationID, 100)
	if err != nil {
		return nil, err
	}
	if err := s.Store.SetConversationCache(ctx, conversationID, messages); err != nil {
		_ = err
	}
	return messages, nil
}

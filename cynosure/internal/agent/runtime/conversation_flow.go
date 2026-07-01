package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/runtime/compression"
	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/idgen"
	"github.com/Cynosure700/cynosure/cynosure/internal/llm"
	"github.com/Cynosure700/cynosure/cynosure/internal/logger"
	agenttools "github.com/Cynosure700/cynosure/cynosure/internal/tools"
)

const (
	assistantDeltaEvent       = "assistant_delta"
	assistantDeltaDoneEvent   = "assistant_delta_done"
	assistantStreamResetEvent = "assistant_stream_reset"
	toolCallStartEvent        = "tool_call_start"
	toolCallDoneEvent         = "tool_call_done"
	toolCallGroupClearEvent   = "tool_call_group_clear"
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
	truncationMaxTokens = 64 * 1000
	// maxResumeAttempts 表示在升级预算后仍发生截断时，最多发起的续写请求次数。
	maxResumeAttempts = 3

	// truncationResumePrompt 作为一条 user 消息注入，用于让模型继续输出被
	// token 上限截断的内容。
	truncationResumePrompt = `Output token limit hit. Resume directly — no apology, no recap of what you were doing. Pick up mid-thought if that is where the cut happened. Break remaining work into smaller pieces.`

	emptyFinalAnswerRetryPrompt = `请继续执行用户的指令。若执行完成，请输出总结；若执行受阻，请输出原因。`
)

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
	memoryIndex := ""
	memorySection := ""
	if memoryOn {
		memoryIndex, memorySection = s.buildMemorySection(ctx, conversation.ID, user, state.History)
	}
	state.SystemPrompt = s.buildSystemPromptWithMemory(user, snapshot, memoryIndex, memorySection)
	state.SystemReminder = s.buildSystemReminderWithMemory(user, snapshot, memoryIndex, memorySection)
	state.Messages = buildOpenAIMessagesWithReminder(state.SystemPrompt, state.SystemReminder, state.ModelHistory)
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
		state.Messages = buildOpenAIMessagesWithReminder(state.SystemPrompt, state.SystemReminder, state.ModelHistory)
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, s.Tools, roundsSinceTodoWrite)
		estimator := compression.DefaultTokenEstimator{ModelID: s.Cfg.LLM.ModelID}
		state.LastContextTokens = estimator.EstimateRequestTokens(estimatePromptWithReminder(state.SystemPrompt, state.SystemReminder), state.ModelHistory, toolDefs)
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
		logger.LogLLMRound(round, fmt.Sprintf("main-agent conversation=%s", conversation.ID), reqBody, respBody, string(finishReason), err)
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
			if finishReason != "tool_calls" && len(msg.ToolCalls) == 0 && strings.TrimSpace(msg.Content) == "" {
				appendInternalUserPrompt(state, wrapSystemReminder(emptyFinalAnswerRetryPrompt))
				continue
			}
			// 会话结束：把本轮模型最终回复纳入真实消息历史后重算，得到会话结束时的最终
			// token 用量，覆盖请求前的估算值，确保存储与下发的都是最终用量。
			finalAssistant := storage.Message{Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)}
			finalHistoryForEstimate := append(cloneMessages(state.ModelHistory), finalAssistant)
			finalEstimator := compression.DefaultTokenEstimator{ModelID: s.Cfg.LLM.ModelID}
			state.LastContextTokens = finalEstimator.EstimateRequestTokens(estimatePromptWithReminder(state.SystemPrompt, state.SystemReminder), finalHistoryForEstimate, toolDefs)
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

		batch, err := s.executeToolCallBatch(ctx, state, msg.ToolCalls, toolBatchOptions{
			ToolContext: ToolContext{User: user, Conversation: conversation, Skills: snapshot, PersistedOutputReader: s.newPersistedOutputReader(conversation.ID, user.ID), Todos: state.Todos, Writer: writer},
		})
		if err != nil {
			return storage.Message{}, err
		}
		if batch.Rejected {
			return storage.Message{Role: "assistant", Content: "操作已被拒绝，已结束本轮。"}, nil
		}
		// 需求2 条件(1)/初次提取：每个 tool_calls round 结束后评估并按需异步刷新会话记忆。
		// 提取数据源用【真实模型线】state.ModelHistory（lockstep 追加，含工具调用与结果）。
		if memoryOn {
			s.maybeUpdateSessionMemoryMidLoop(conversation, user, state.ModelHistory, state.LastContextTokens, len(msg.ToolCalls))
		}
	}
}

type toolBatchOptions struct {
	ToolContext      ToolContext
	Registry         *ToolRegistry
	UseChildRegistry bool
	SuppressResultUI bool
	EphemeralGroupID string
}

type toolBatchResult struct {
	Contexts []*ToolUseContext
	Rejected bool
}

func (s *Service) executeToolCallBatch(ctx context.Context, state *LoopState, calls []openai.ToolCall, opts toolBatchOptions) (toolBatchResult, error) {
	contexts := make([]*ToolUseContext, 0, len(calls))
	for _, tc := range calls {
		toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
		if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
			return toolBatchResult{}, err
		}
		if approved, _ := s.approveToolCall(ctx, tc); !approved {
			toolCtx.Outcome = toolExecutionOutcome{Status: "rejected", Result: "Error: user rejected this operation", Audit: toolCtx.Outcome.Audit}
			emitToolCallStartWithOptions(state, toolCtx, opts)
			emitToolCallDoneWithOptions(state, toolCtx, opts)
			return toolBatchResult{Contexts: append(contexts, toolCtx), Rejected: true}, nil
		}
		contexts = append(contexts, toolCtx)
	}
	for _, toolCtx := range contexts {
		emitToolCallStartWithOptions(state, toolCtx, opts)
	}

	todosSnapshot := append([]agenttools.TodoItem(nil), opts.ToolContext.Todos...)
	var wg sync.WaitGroup
	for _, toolCtx := range contexts {
		toolCtx := toolCtx
		execCtx := opts.ToolContext
		execCtx.Todos = append([]agenttools.TodoItem(nil), todosSnapshot...)
		execCtx.ParentToolCallID = toolCtx.ToolCall.ID
		wg.Add(1)
		go func() {
			defer wg.Done()
			if opts.UseChildRegistry {
				toolCtx.Outcome = s.executeChildToolCall(ctx, opts.Registry, execCtx, toolCtx.Name, toolCtx.RawArgs, toolCtx.Outcome.Audit)
			} else {
				toolCtx.Outcome = s.executeToolCall(ctx, execCtx, toolCtx.Name, toolCtx.RawArgs, toolCtx.Outcome.Audit)
			}
			// 工具刚执行完，文件内容最新，此刻计算 edit_file/multi_edit 的 diff 真实
			// 行号，供事件下发与展示历史持久化（不进入模型上下文）。
			toolCtx.Outcome.EditLineStarts = editLineStartsForToolCtx(state, toolCtx)
			// 并行执行时谁先完成谁先更新展示信息，无需等待整批完成。
			// 落库/历史追加仍在 wg.Wait() 之后串行进行，避免数据竞争。
			emitToolCallDoneWithOptions(state, toolCtx, opts)
		}()
	}
	wg.Wait()

	for _, toolCtx := range contexts {
		if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
			state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
		}
		if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
			return toolBatchResult{}, err
		}
		state.ToolCallCount++
		emitMeta(state)
	}
	return toolBatchResult{Contexts: contexts}, nil
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
	emitToolCallStartWithOptions(state, toolCtx, toolBatchOptions{})
}

func emitToolCallStartWithOptions(state *LoopState, toolCtx *ToolUseContext, opts toolBatchOptions) {
	if state == nil || state.Writer == nil || toolCtx == nil {
		return
	}
	payload := map[string]any{
		"tool_call_id": toolCtx.ToolCall.ID,
		"tool_name":    toolCtx.Name,
		"raw_args":     toolCtx.RawArgs,
		"args_preview": previewToolArgs(toolCtx.Name, toolCtx.RawArgs),
		"status":       "running",
	}
	applyToolEventOptions(payload, opts)
	_ = state.Writer.Event(toolCallStartEvent, payload)
}

func emitToolCallDone(state *LoopState, toolCtx *ToolUseContext) {
	emitToolCallDoneWithOptions(state, toolCtx, toolBatchOptions{})
}

func emitToolCallDoneWithOptions(state *LoopState, toolCtx *ToolUseContext, opts toolBatchOptions) {
	if state == nil || state.Writer == nil || toolCtx == nil {
		return
	}
	payload := map[string]any{
		"tool_call_id":   toolCtx.ToolCall.ID,
		"tool_name":      toolCtx.Name,
		"raw_args":       toolCtx.RawArgs,
		"args_preview":   previewToolArgs(toolCtx.Name, toolCtx.RawArgs),
		"status":         toolCtx.Outcome.Status,
		"result_preview": previewToolResult(toolCtx.Outcome),
		"audit_summary":  toolCtx.Outcome.AuditSummary(),
	}
	if len(toolCtx.Outcome.EditLineStarts) > 0 {
		payload["edit_line_starts"] = toolCtx.Outcome.EditLineStarts
	}
	applyToolEventOptions(payload, opts)
	if opts.SuppressResultUI {
		delete(payload, "result_preview")
	}
	_ = state.Writer.Event(toolCallDoneEvent, payload)
}

// editLineStartsForToolCtx 在工具执行完成后计算 edit_file/multi_edit 的 diff 真实
// 行号。仅对成功的编辑类工具计算，工作区根目录取自运行时环境；其它情况返回 nil。
func editLineStartsForToolCtx(state *LoopState, toolCtx *ToolUseContext) [][]int {
	if state == nil || toolCtx == nil || toolCtx.Outcome.Status != "success" {
		return nil
	}
	return agenttools.EditFileLineStarts(state.RuntimeEnv().WorkspaceRoot, toolCtx.Name, toolCtx.RawArgs)
}

func applyToolEventOptions(payload map[string]any, opts toolBatchOptions) {
	if opts.EphemeralGroupID == "" && !opts.SuppressResultUI {
		return
	}
	payload["scope"] = "subagent"
	if opts.EphemeralGroupID != "" {
		payload["ephemeral_group_id"] = opts.EphemeralGroupID
	}
	if opts.SuppressResultUI {
		payload["suppress_result"] = true
	}
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
//   - 上下文溢出（HTTP 413）：执行一次 reactiveCompact——它是「单次压缩」，内部按对话
//     轮渐进剥离（最多 3 次，满足 token 阈值即提前成功停止；3 次后仍超限则返回错误），
//     成功后从压缩后的历史重建请求并重试一次。若 reactiveCompact 返回错误、或压缩后重试
//     仍溢出，则返回给调用方（由既有的兜底边界处理），不再重复触发。
func (s *Service) runModelRoundWithRecovery(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, error) {
	cur := req
	upgraded := false
	compacted := false
	resumeAttempts := 0
	var accumulated strings.Builder
	var accumulatedReasoning strings.Builder

	for {
		msg, finishReason, streamedContent, streamID, err := s.runModelRoundStream(ctx, state, cur)
		if err != nil {
			if llm.IsContextOverflow(err) && !compacted {
				// 上下文溢出会丢弃这一轮（可能已流式输出的）部分内容并从压缩后的
				// 历史重试，因此先让 UI 清空已显示的半截输出。
				if streamedContent {
					emitStreamReset(state, streamID)
				}
				// 单次压缩：ReactiveCompact 内部完成最多 3 次剥离。压缩后重试一次；
				// 仍溢出则由上层兜底，不再重复触发。
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
			// 内容已在 runModelRoundStream 中边收边发；这里只负责把分段续写
			// 的结果拼接成一条完整 assistant 消息返回。
			if resuming {
				return openai.ChatCompletionMessage{
					Role:             "assistant",
					Content:          accumulated.String() + msg.Content,
					ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
					ToolCalls:        msg.ToolCalls,
				}, finishReason, nil
			}
			return msg, finishReason, nil
		}

		// 输出被截断。
		if !upgraded {
			// 丢弃部分输出，升级预算，原样重试。已流式输出的半截内容需要让
			// UI 先清空，避免与升级后重试得到的完整内容叠加。
			if streamedContent {
				emitStreamReset(state, streamID)
			}
			upgraded = true
			cur.MaxTokens = truncationMaxTokens
			continue
		}
		if resumeAttempts >= maxResumeAttempts {
			// 尽力而为：返回目前已获得的内容（已流式输出，无需重发）。
			logger.Warn(fmt.Sprintf("truncation resume exhausted conversation=%s after %d attempts", state.Conversation.ID, resumeAttempts))
			return openai.ChatCompletionMessage{
				Role:             "assistant",
				Content:          accumulated.String() + msg.Content,
				ReasoningContent: accumulatedReasoning.String() + msg.ReasoningContent,
			}, finishReason, nil
		}
		// 保留这一段部分输出（已流式输出），并要求模型继续。
		accumulated.WriteString(msg.Content)
		accumulatedReasoning.WriteString(msg.ReasoningContent)
		resumeAttempts++
		cur.Messages = append(cur.Messages,
			openai.ChatCompletionMessage{Role: "assistant", Content: msg.Content},
			openai.ChatCompletionMessage{Role: "user", Content: truncationResumePrompt},
		)
	}
}

// emitStreamReset 通知 UI 丢弃本轮已流式显示的半截 assistant 内容。用于输出被
// 截断升级重试、或上下文溢出压缩重试这类「已发出的增量会被废弃并重新生成」的
// 场景，避免被废弃的半截输出残留在界面上与重试结果叠加。
func emitStreamReset(state *LoopState, streamID string) {
	if state == nil || state.Writer == nil {
		return
	}
	payload := map[string]any{}
	if streamID != "" {
		payload["stream_id"] = streamID
	}
	_ = state.Writer.Event(assistantStreamResetEvent, payload)
}

func (s *Service) runModelRoundStream(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, bool, string, error) {
	stream, err := s.LLM.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, "", false, "", err
	}
	defer stream.Close()

	var content strings.Builder
	var reasoningContent strings.Builder
	var finishReason openai.FinishReason
	toolCalls := &streamedToolCallAccumulator{}
	streamedContent := false
	seenChoice := false
	seenOutput := false
	streamID := idgen.New("stream")
	deltaSeq := 0

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return openai.ChatCompletionMessage{}, "", streamedContent, streamID, err
		}
		if len(chunk.Choices) == 0 {
			continue
		}
		seenChoice = true
		choice := chunk.Choices[0]
		if choice.Delta.Content != "" {
			seenOutput = true
			content.WriteString(choice.Delta.Content)
			// 边收边发：每个内容增量都立刻下发给 UI，让回复随生成逐字呈现，
			// 而不是在整段流结束后一次性补发，避免长时间空白后瞬间刷屏。
			if state.Writer != nil {
				deltaSeq++
				_ = state.Writer.Event(assistantDeltaEvent, map[string]any{
					"content":   choice.Delta.Content,
					"stream_id": streamID,
					"delta_seq": deltaSeq,
				})
				streamedContent = true
			}
		}
		if choice.Delta.ReasoningContent != "" {
			seenOutput = true
			// reasoning_content 仍累加并随消息落库/参与记忆，但不再向 UI 流式输出。
			reasoningContent.WriteString(choice.Delta.ReasoningContent)
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
		return openai.ChatCompletionMessage{}, "", streamedContent, streamID, fmt.Errorf("model stream returned no choices")
	}
	calls := toolCalls.Calls()
	if streamedContent && state.Writer != nil {
		_ = state.Writer.Event(assistantDeltaDoneEvent, map[string]any{
			"content":     content.String(),
			"stream_id":   streamID,
			"delta_count": deltaSeq,
		})
	}

	return openai.ChatCompletionMessage{Role: "assistant", Content: content.String(), ReasoningContent: reasoningContent.String(), ToolCalls: calls}, finishReason, streamedContent, streamID, nil
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

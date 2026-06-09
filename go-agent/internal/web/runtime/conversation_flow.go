package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/runtime/compression"
	"nano_cc/internal/web/storage"
)

const (
	assistantDeltaEvent = "assistant_delta"
	reasoningDeltaEvent = "reasoning_delta"
)

type bufferedModelDelta struct {
	Event   string
	Content string
}

func (s *Service) RespondToConversation(ctx context.Context, conversation storage.Conversation, user storage.User, userMessage string, writer EventWriter) (storage.Message, error) {
	// 入口获取会话锁：上一轮收尾未完成时阻塞等待，直到拿到锁或等待超时。
	// 获取失败（Redis 异常 / 等待超时）则降级放行，跳过本轮收尾。
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
	state.SystemPrompt = s.buildSystemPrompt(ctx, conversation, user, snapshot, state.History)
	state.Messages = buildOpenAIMessages(state.SystemPrompt, state.ModelHistory)
	round := 0
	roundsSinceTodoWrite := 0
	var cumulativeReasoning strings.Builder
	var lastRequestHistory []storage.Message

	for {
		round++
		requestHistory, err := s.compressContextBeforeLLM(ctx, state)
		if err != nil {
			return storage.Message{}, err
		}
		lastRequestHistory = requestHistory
		state.Messages = buildOpenAIMessages(state.SystemPrompt, requestHistory)
		roundsSinceTodoWrite = maybeAppendTodoWriteReminder(state, s.Tools, roundsSinceTodoWrite)
		estimator := compression.DefaultTokenEstimator{}
		state.LastContextTokens = estimator.EstimateRequestTokens(state.SystemPrompt, requestHistory, s.Tools.Definitions())
		state.LastContextBudget = estimator.ContextTokenBudget()
		emitMeta(state)
		req := openai.ChatCompletionRequest{
			Model:    s.Cfg.LLM.ModelID,
			Messages: state.Messages,
			Tools:    s.Tools.Definitions(),
		}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundStream(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("main-agent conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
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
			stopCtx := &StopContext{State: state, ModelMessage: msg, Content: fallbackAssistantContent(msg.Content), ReasoningContent: cumulativeReasoning.String()}
			if err := s.hookManager().RunStop(ctx, stopCtx); err != nil {
				return storage.Message{}, err
			}
			if s.EnableMemory {
				finalHistory := append(state.History, stopCtx.AssistantMessage)
				finalModelHistory := append(cloneMessages(lastRequestHistory), stopCtx.AssistantMessage)
				handedOff = s.scheduleMemoryWork(conversation, user, finalHistory, finalModelHistory, lockToken, stopRenew)
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
			toolCtx.Outcome = s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Skills: snapshot, ParentToolCallID: tc.ID, PersistedOutputReader: s.newPersistedOutputReader(conversation.ID, user.ID)}, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if toolCtx.Name == agenttools.TodoWriteToolName && toolCtx.Outcome.Status == "success" {
				state.Todos = append([]agenttools.TodoItem(nil), toolCtx.Outcome.Todos...)
			}
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			state.ToolCallCount++
			emitMeta(state)
		}
	}
}

// emitMeta 通过 SSE 实时下发当前累计的回复元信息。
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

func (s *Service) runModelRoundStream(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, error) {
	stream, err := s.LLM.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, "", err
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
			return openai.ChatCompletionMessage{}, "", err
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
		return openai.ChatCompletionMessage{}, "", fmt.Errorf("model stream returned no choices")
	}
	calls := toolCalls.Calls()
	if state.Writer != nil && shouldEmitAssistantContentDeltas(finishReason, calls) {
		for _, delta := range bufferedContentDeltas {
			_ = state.Writer.Event(delta.Event, map[string]any{"content": delta.Content})
		}
	}

	return openai.ChatCompletionMessage{Role: "assistant", Content: content.String(), ReasoningContent: reasoningContent.String(), ToolCalls: calls}, finishReason, nil
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

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/web/storage"
)

func (s *Service) RespondToConversation(ctx context.Context, conversation storage.Conversation, user storage.User, userMessage string, writer EventWriter) (storage.Message, error) {
	history, err := s.loadConversationMessages(ctx, conversation.ID)
	if err != nil {
		return storage.Message{}, err
	}
	state := s.newLoopState(conversation, user, userMessage, history, writer)
	if err := s.hookManager().RunUserPromptSubmit(ctx, &UserPromptSubmitContext{State: state}); err != nil {
		return storage.Message{}, err
	}

	snapshot, err := s.buildSkillSnapshot(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	state.SkillSnapshot = snapshot
	state.SystemPrompt = s.buildSystemPrompt(user, snapshot)
	state.Messages = buildOpenAIMessages(state.SystemPrompt, state.History)
	round := 0

	for {
		round++
		req := openai.ChatCompletionRequest{
			Model:    s.Cfg.LLM.ModelID,
			Messages: state.Messages,
			Tools:    s.Tools.Definitions(),
		}
		reqBody, _ := json.Marshal(req)
		msg, finishReason, err := s.runModelRoundStream(ctx, state, req)
		respBody, _ := json.Marshal(msg)
		logger.LogLLMRound(round, fmt.Sprintf("web-runtime conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
			return storage.Message{}, err
		}
		requestMsg := msg
		state.Messages = append(state.Messages, requestMsg)

		if finishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			stopCtx := &StopContext{State: state, ModelMessage: msg, Content: fallbackAssistantContent(msg.Content), ReasoningContent: msg.ReasoningContent}
			if err := s.hookManager().RunStop(ctx, stopCtx); err != nil {
				return storage.Message{}, err
			}
			return stopCtx.AssistantMessage, nil
		}
		state.History = append(state.History, storage.Message{ID: state.NextMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "assistant", Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCalls: openAIToolCallsToStorage(msg.ToolCalls)})

		for _, tc := range msg.ToolCalls {
			toolCtx := &ToolUseContext{State: state, ToolCall: tc, Name: tc.Function.Name, RawArgs: tc.Function.Arguments}
			if err := s.hookManager().RunPreToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
			toolCtx.Outcome = s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Skills: snapshot}, tc.Function.Name, tc.Function.Arguments, toolCtx.Outcome.Audit)
			if err := s.hookManager().RunPostToolUse(ctx, toolCtx); err != nil {
				return storage.Message{}, err
			}
		}
	}
}

func (s *Service) runModelRoundStream(ctx context.Context, state *LoopState, req openai.ChatCompletionRequest) (openai.ChatCompletionMessage, openai.FinishReason, error) {
	stream, err := config.Client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return openai.ChatCompletionMessage{}, "", err
	}
	defer stream.Close()

	var content strings.Builder
	var reasoningContent strings.Builder
	var finishReason openai.FinishReason
	toolCalls := &streamedToolCallAccumulator{}
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
			if state.Writer != nil {
				_ = state.Writer.Event("assistant_delta", map[string]any{"content": choice.Delta.Content})
			}
		}
		if choice.Delta.ReasoningContent != "" {
			seenOutput = true
			reasoningContent.WriteString(choice.Delta.ReasoningContent)
			if state.Writer != nil {
				_ = state.Writer.Event("reasoning_delta", map[string]any{"content": choice.Delta.ReasoningContent})
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

	return openai.ChatCompletionMessage{Role: "assistant", Content: content.String(), ReasoningContent: reasoningContent.String(), ToolCalls: toolCalls.Calls()}, finishReason, nil
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

func (s *Service) persistAssistantReply(ctx context.Context, conversation storage.Conversation, userID string, history []storage.Message, content string, reasoningContent string, writer EventWriter) (storage.Message, error) {
	assistant := storage.Message{ID: newMessageID(), ConversationID: conversation.ID, UserID: userID, Role: "assistant", Content: content, ReasoningContent: reasoningContent}
	updatedHistory := append(history, assistant)
	if err := s.Store.SetConversationHistory(ctx, conversation.ID, updatedHistory); err != nil {
		return storage.Message{}, err
	}
	if writer != nil {
		_ = writer.Event("assistant", map[string]any{"content": assistant.Content, "reasoning_content": assistant.ReasoningContent})
	}
	return assistant, nil
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

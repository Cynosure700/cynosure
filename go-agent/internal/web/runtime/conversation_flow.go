package runtime

import (
	"context"
	"encoding/json"
	"fmt"

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
		resp, err := config.Client.CreateChatCompletion(ctx, req)
		respBody, _ := json.Marshal(resp)
		logger.LogLLMRound(round, fmt.Sprintf("web-runtime conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
			return storage.Message{}, err
		}
		if len(resp.Choices) == 0 {
			return storage.Message{}, fmt.Errorf("model returned no choices")
		}
		choice := resp.Choices[0]
		msg := choice.Message
		requestMsg := msg
		state.Messages = append(state.Messages, requestMsg)

		if choice.FinishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			stopCtx := &StopContext{State: state, ModelMessage: msg, Content: fallbackAssistantContent(msg.Content), ReasoningContent: msg.ReasoningContent}
			if err := s.hookManager().RunStop(ctx, stopCtx); err != nil {
				return storage.Message{}, err
			}
			return stopCtx.AssistantMessage, nil
		}

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

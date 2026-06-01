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
	userEntry := storage.Message{ID: newMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "user", Content: userMessage}
	history = append(history, userEntry)
	if shouldInferConversationTitle(conversation.Title) {
		if err := s.Store.UpdateConversationTitle(ctx, conversation.ID, inferConversationTitle(userMessage)); err != nil {
			return storage.Message{}, err
		}
	} else {
		if err := s.Store.TouchConversationActivity(ctx, conversation.ID); err != nil {
			return storage.Message{}, err
		}
	}

	snapshot, err := s.buildSkillSnapshot(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	systemPrompt := s.buildSystemPrompt(user, snapshot)
	messages := buildOpenAIMessages(systemPrompt, history)
	round := 0

	for {
		round++
		req := openai.ChatCompletionRequest{
			Model:    s.Cfg.LLM.ModelID,
			Messages: messages,
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
		messages = append(messages, msg)

		if choice.FinishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			return s.persistAssistantReply(ctx, conversation, user.ID, history, fallbackAssistantContent(msg.Content), writer)
		}

		for _, tc := range msg.ToolCalls {
			outcome := s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Skills: snapshot}, tc.Function.Name, tc.Function.Arguments)
			_ = s.Store.CreateToolCall(ctx, storage.ToolCall{ID: newToolCallID(), ConversationID: conversation.ID, UserID: user.ID, ToolName: tc.Function.Name, Status: outcome.Status, Summary: outcome.AuditSummary()})
			if writer != nil {
				_ = writer.Event("tool", map[string]any{"name": tc.Function.Name, "status": outcome.Status, "result": outcome.Result})
			}
			messages = append(messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: tc.ID, Content: outcome.MessageContent()})
		}
	}
}

func (s *Service) persistAssistantReply(ctx context.Context, conversation storage.Conversation, userID string, history []storage.Message, content string, writer EventWriter) (storage.Message, error) {
	assistant := storage.Message{ID: newMessageID(), ConversationID: conversation.ID, UserID: userID, Role: "assistant", Content: content}
	updatedHistory := append(history, assistant)
	if err := s.Store.SetConversationHistory(ctx, conversation.ID, updatedHistory); err != nil {
		return storage.Message{}, err
	}
	if writer != nil {
		_ = writer.Event("assistant", map[string]any{"content": assistant.Content})
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

package hooks

import (
	"context"

	"cynosure/internal/agent/storage"
)

func persistAssistantStopHook(ctx context.Context, h *StopContext) error {
	state := h.State
	h.AssistantMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: "assistant", Content: h.Content, ReasoningContent: h.ReasoningContent, Meta: stopMessageMeta(state)}
	updatedHistory := append(state.History, h.AssistantMessage)
	return state.Store.SetConversationHistory(ctx, state.Conversation.ID, updatedHistory)
}

func emitAssistantStopHook(ctx context.Context, h *StopContext) error {
	if h.State.Writer == nil {
		return nil
	}
	meta := stopMessageMeta(h.State)
	_ = h.State.Writer.Event("assistant", map[string]any{
		"message_id":        h.AssistantMessage.ID,
		"content":           h.AssistantMessage.Content,
		"reasoning_content": h.AssistantMessage.ReasoningContent,
		"final":             true,
		"tool_call_count":   meta.ToolCallCount,
		"context_tokens":    meta.ContextTokens,
		"context_budget":    meta.ContextBudget,
	})
	return nil
}

func stopMessageMeta(state *LoopState) *storage.MessageMeta {
	return &storage.MessageMeta{
		ToolCallCount: state.ToolCallCount,
		ContextTokens: state.LastContextTokens,
		ContextBudget: state.LastContextBudget,
	}
}

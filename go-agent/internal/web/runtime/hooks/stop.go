package hooks

import (
	"context"

	"nano_cc/internal/web/storage"
)

func persistAssistantStopHook(ctx context.Context, h *StopContext) error {
	state := h.State
	h.AssistantMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: "assistant", Content: h.Content, ReasoningContent: h.ReasoningContent}
	updatedHistory := append(state.History, h.AssistantMessage)
	return state.Store.SetConversationHistory(ctx, state.Conversation.ID, updatedHistory)
}

func emitAssistantStopHook(ctx context.Context, h *StopContext) error {
	if h.State.Writer == nil {
		return nil
	}
	_ = h.State.Writer.Event("assistant", map[string]any{"content": h.AssistantMessage.Content, "reasoning_content": h.AssistantMessage.ReasoningContent})
	return nil
}

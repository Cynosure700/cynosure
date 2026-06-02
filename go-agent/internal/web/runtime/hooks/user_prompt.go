package hooks

import (
	"context"

	"nano_cc/internal/web/storage"
)

func appendUserMessageHook(ctx context.Context, h *UserPromptSubmitContext) error {
	state := h.State
	state.UserMessage = storage.Message{ID: state.NextMessageID(), ConversationID: state.Conversation.ID, UserID: state.User.ID, Role: "user", Content: state.UserInput}
	state.History = append(state.History, state.UserMessage)
	return nil
}

func conversationActivityHook(ctx context.Context, h *UserPromptSubmitContext) error {
	state := h.State
	if state.ShouldInferTitle(state.Conversation.Title) {
		return state.Store.UpdateConversationTitle(ctx, state.Conversation.ID, state.InferTitle(state.UserInput))
	}
	return state.Store.TouchConversationActivity(ctx, state.Conversation.ID)
}

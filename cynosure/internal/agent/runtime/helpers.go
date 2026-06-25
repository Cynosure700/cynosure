package runtime

import (
	"strings"

	"cynosure/internal/agent/storage"
	"cynosure/internal/idgen"
)

func newMessageID() string { return idgen.New("msg") }

func fallbackAssistantContent(content string) string {
	if strings.TrimSpace(content) == "" {
		return "(no response)"
	}
	return content
}

func appendInternalUserPrompt(state *LoopState, content string) {
	state.ModelHistory = append(state.ModelHistory, storage.Message{
		ID:             state.NextMessageID(),
		ConversationID: state.Conversation.ID,
		UserID:         state.User.ID,
		Role:           "system",
		Content:        content,
	})
}

func shouldInferConversationTitle(currentTitle string) bool {
	trimmed := strings.TrimSpace(currentTitle)
	return trimmed == "" || trimmed == "新对话"
}

func inferConversationTitle(userMessage string) string {
	trimmed := strings.TrimSpace(userMessage)
	if len([]rune(trimmed)) > 30 {
		return string([]rune(trimmed)[:30])
	}
	if trimmed == "" {
		return "新对话"
	}
	return trimmed
}

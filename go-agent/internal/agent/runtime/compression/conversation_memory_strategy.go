package compression

import (
	"context"
	"strings"

	"nano_cc/internal/agent/storage"
)

const conversationMemoryStrategyName = "conversation_memory"

const (
	conversationMemoryMinTokens = 10 * 1024
	conversationMemoryMaxTokens = 40 * 1024
	conversationMemoryMinCount  = 5

	conversationMemorySystemPreamble = "以下是本会话的记忆，由系统随对话持续维护，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。"
)

// ConversationMemoryStrategy runs just before FullHistorySummarizationStrategy.
// When the request still exceeds the token budget but a sufficiently rich,
// already-maintained conversation memory exists, it rebuilds RequestHistory from
// that memory (plus a recent tail and the always-kept last message), so the
// fallback full-history summarization is skipped entirely. Otherwise it is a
// no-op and the fallback summarizer handles compression.
type ConversationMemoryStrategy struct{}

func (s *ConversationMemoryStrategy) Name() string {
	return conversationMemoryStrategyName
}

func (s *ConversationMemoryStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Store == nil {
		return nil
	}
	budget := req.Estimator.ContextTokenBudget()
	before := req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools)
	if before <= budget {
		return nil
	}
	if len(req.RequestHistory) < 2 {
		return nil
	}

	items, err := req.Store.ListConversationMemories(ctx, req.Conversation.ID)
	if err != nil || len(items) < conversationMemoryMinCount {
		return nil
	}
	memoryText := renderConversationMemory(items)
	memTokens := estimateTokensFromBytes(len([]byte(memoryText)))
	if memTokens < conversationMemoryMinTokens || memTokens > conversationMemoryMaxTokens {
		return nil
	}

	lastMessage := req.RequestHistory[len(req.RequestHistory)-1]
	summarizable := req.RequestHistory[:len(req.RequestHistory)-1]
	tail := selectRecentTail(req, summarizable, memoryText, budget)
	kept := repairToolCallBoundaries(append(append([]storage.Message{}, tail...), lastMessage))
	req.RequestHistory = buildConversationMemoryHistory(memoryText, kept)
	return nil
}

func renderConversationMemory(items []storage.ConversationMemory) string {
	var b strings.Builder
	for _, m := range items {
		b.WriteString("- ")
		b.WriteString(m.Name)
		if strings.TrimSpace(m.Description) != "" {
			b.WriteString("：")
			b.WriteString(m.Description)
		}
		b.WriteByte('\n')
		if strings.TrimSpace(m.Body) != "" {
			b.WriteString("  ")
			b.WriteString(m.Body)
			b.WriteByte('\n')
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

func buildConversationMemoryHistory(memoryText string, tail []storage.Message) []storage.Message {
	messages := []storage.Message{
		{Role: "system", Content: conversationMemorySystemPreamble},
		{Role: "user", Content: "<conversation-memory>\n" + memoryText + "\n</conversation-memory>"},
	}
	return append(messages, tail...)
}

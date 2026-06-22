package compression

import (
	"context"
	"strings"

	"cynosure/internal/agent/storage"
)

const conversationMemoryStrategyName = "conversation_memory"

const (
	conversationMemoryMinTokens = 10 * 1024
	conversationMemoryMaxTokens = 40 * 1024
	conversationMemoryMinCount  = 5

	conversationMemorySystemPreamble = "以下是本会话的记忆，由系统随对话持续维护，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。"
)

// ConversationMemoryStrategy 在 FullHistorySummarizationStrategy 之前运行。
// 当请求仍超出 token 预算，但存在一份足够丰富、已持续维护的会话记忆时，
// 它会用该记忆（加上一个近期尾部以及始终保留的最后一条消息）重建 RequestHistory，
// 从而完全跳过兜底的全量历史摘要。否则它什么都不做，由兜底摘要器处理压缩。
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

	// 需求3：尾部恰为"未被 session memory 记忆到"的未压缩消息——从断点（含端点）到末尾，
	// 且取自【逐字 display history】（而非可能已带压缩产物的模型线 RequestHistory）。
	// 断点未知/未命中（未持久化、被消息窗口裁剪等）→ no-op，交全量摘要兜底，无数据丢失。
	tail, ok := selectUnfoldedTail(req.DisplayHistory, req.ConversationMemoryBreakpoint)
	if !ok {
		return nil
	}
	kept := repairToolCallBoundaries(tail)
	candidate := buildConversationMemoryHistory(memoryText, kept)
	// 保持"本策略一旦动手，结果须落在预算内"的不变式：未折叠尾部过大则不动手，交全量摘要兜底，
	// 避免记忆块被随后的全量摘要二次压缩。
	if req.Estimator.EstimateRequestTokens(req.SystemPrompt, candidate, req.Tools) > budget {
		return nil
	}
	req.RequestHistory = candidate
	return nil
}

// selectUnfoldedTail 返回断点消息（含端点）之后的尾部消息。命中返回 (tail, true)；
// 断点为空或在历史中未找到返回 (nil, false)。
func selectUnfoldedTail(history []storage.Message, breakpointID string) ([]storage.Message, bool) {
	if breakpointID == "" {
		return nil, false
	}
	for i := range history {
		if history[i].ID == breakpointID {
			return history[i:], true
		}
	}
	return nil, false
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

package compression

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

const fullHistorySummarizationStrategyName = "full_history_summarization"

const (
	summaryTargetTokens   = 8 * 1000
	recentTailMessages    = 5
	summarySystemPreamble = "以下是为满足上下文窗口限制生成的会话摘要，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。"
)

// HistorySummarizer 将请求状态的历史摘要为一段 markdown 摘要。
type HistorySummarizer func(ctx context.Context, req SummaryRequest) (SummaryResult, error)

type SummaryRequest struct {
	Conversation storage.Conversation
	User         storage.User
	History      []storage.Message
	TargetTokens int
	// Aggressive 为 true 时，摘要器应使用更激进的 ReactiveCompact 摘要 prompt
	//（更紧凑、只保留续写必需的最小事实）。由 ReactiveCompactStrategy 传入。
	Aggressive bool
}

type SummaryResult struct {
	Summary string
}

// FullHistorySummarizationStrategy 是兜底层：当确定性的三层压缩后仍超出 token 预算时，
// 它用一段摘要加上一个近期尾部窗口替换 RequestHistory。仅作用于请求，绝不写入展示历史。
type FullHistorySummarizationStrategy struct{}

func (s *FullHistorySummarizationStrategy) Name() string {
	return fullHistorySummarizationStrategyName
}

func (s *FullHistorySummarizationStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Summarizer == nil {
		return nil
	}
	budget := req.Estimator.ContextTokenBudget()
	before := req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools)
	if before <= budget {
		return nil
	}
	// 至少要留下一条可摘要的消息；仅有一条消息时无内容可摘要。
	if len(req.RequestHistory) < 2 {
		return nil
	}

	// 摘要与保留的尾部都基于【逐字 display history】，使保留的最近消息是真实、未经压缩
	// 改写的原始消息；display history 为空/过短时退回 RequestHistory，避免产出空历史。
	source := req.DisplayHistory
	if len(source) < 2 {
		source = req.RequestHistory
	}

	// 保留最近 N 条消息逐字不变；对其之前的全部内容做摘要。
	// 对尾部数量设上限，以确保至少有一条消息可被摘要。
	n := len(source)
	tailCount := recentTailMessages
	if tailCount > n-1 {
		tailCount = n - 1
	}
	summarizable := source[:n-tailCount]
	tail := source[n-tailCount:]

	result, err := req.Summarizer(ctx, SummaryRequest{
		Conversation: req.Conversation,
		User:         req.User,
		History:      summarizable,
		TargetTokens: summaryTargetTokens,
	})
	if err != nil {
		return fmt.Errorf("summarize history: %w", err)
	}
	summary := result.Summary

	kept := repairToolCallBoundaries(tail)
	req.RequestHistory = buildSummaryHistory(summary, kept)
	return nil
}

func buildSummaryHistory(summary string, tail []storage.Message) []storage.Message {
	messages := []storage.Message{
		{Role: "system", Content: summarySystemPreamble},
		{Role: "user", Content: "<conversation-summary>\n" + summary + "\n</conversation-summary>"},
	}
	return append(messages, tail...)
}

func messageBytes(msg storage.Message) int {
	if data, err := json.Marshal(msg); err == nil {
		return len(data)
	}
	return len(msg.Content) + len(msg.ReasoningContent)
}

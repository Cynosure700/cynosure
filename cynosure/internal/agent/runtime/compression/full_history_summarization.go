package compression

import (
	"context"
	"encoding/json"
	"fmt"

	"nano_cc/internal/agent/storage"
)

const fullHistorySummarizationStrategyName = "full_history_summarization"

const (
	summaryTargetTokens   = 8 * 1024
	recentTailMessages    = 5
	summarySystemPreamble = "以下是为满足上下文窗口限制生成的会话摘要，仅用于本次模型推理，不是用户发送的真实消息。请把它当作已发生对话的可靠浓缩。"
)

// HistorySummarizer summarizes a request-state history into a markdown summary.
type HistorySummarizer func(ctx context.Context, req SummaryRequest) (SummaryResult, error)

type SummaryRequest struct {
	Conversation storage.Conversation
	User         storage.User
	History      []storage.Message
	TargetTokens int
}

type SummaryResult struct {
	Summary string
}

// FullHistorySummarizationStrategy is the fallback layer: when the deterministic
// three layers still exceed the token budget, it replaces RequestHistory with a
// summary plus a recent tail window. Request-only; never written to display history.
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
	// At least one message must remain summarizable; with a single message
	// there is nothing to summarize.
	if len(req.RequestHistory) < 2 {
		return nil
	}

	// 摘要与保留的尾部都基于【逐字 display history】，使保留的最近消息是真实、未经压缩
	// 改写的原始消息；display history 为空/过短时退回 RequestHistory，避免产出空历史。
	source := req.DisplayHistory
	if len(source) < 2 {
		source = req.RequestHistory
	}

	// Keep the most recent N messages verbatim; summarize everything before.
	// Cap the tail so at least one message stays summarizable.
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

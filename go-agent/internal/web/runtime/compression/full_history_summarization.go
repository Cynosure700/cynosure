package compression

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"nano_cc/internal/idgen"
	"nano_cc/internal/web/storage"
)

const fullHistorySummarizationStrategyName = "full_history_summarization"

const (
	summaryTargetTokens   = 8 * 1024
	recentTailMinTokens   = 16 * 1024
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

	sourceHash := historyHash(req.RequestHistory)

	// Try cache first.
	summary := ""
	if cached, err := req.Store.GetContextSummaryByHistoryHash(ctx, req.Conversation.ID, req.User.ID, sourceHash); err == nil && cached.Summary != "" {
		summary = cached.Summary
	} else {
		result, err := req.Summarizer(ctx, SummaryRequest{
			Conversation: req.Conversation,
			User:         req.User,
			History:      req.RequestHistory,
			TargetTokens: summaryTargetTokens,
		})
		if err != nil {
			return fmt.Errorf("summarize history: %w", err)
		}
		summary = result.Summary
		after := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
		_ = req.Store.CreateContextSummary(ctx, storage.ContextSummary{
			ID:                    "cs_" + idgen.Hex(),
			ConversationID:        req.Conversation.ID,
			UserID:                req.User.ID,
			SourceHistorySHA256:   sourceHash,
			Strategy:              fullHistorySummarizationStrategyName,
			EstimatedTokensBefore: before,
			EstimatedTokensAfter:  after,
			Summary:               summary,
		})
	}

	tail := selectRecentTail(req, summary, budget)
	req.RequestHistory = buildSummaryHistory(summary, tail)
	return nil
}

// selectRecentTail picks messages from the end of the history within the budget
// remaining after reserving space for the summary, then repairs tool boundaries.
func selectRecentTail(req *Request, summary string, budget int) []storage.Message {
	history := req.RequestHistory
	summaryTokens := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
	tailBudget := budget - summaryTokens
	if tailBudget > recentTailMinTokens {
		// Cap reserved tail space to keep summary authoritative.
		tailBudget = recentTailMinTokens
	}
	if tailBudget <= 0 {
		return nil
	}

	var selected []storage.Message
	used := 0
	for i := len(history) - 1; i >= 0; i-- {
		msgTokens := estimateTokensFromBytes(messageBytes(history[i]))
		if used+msgTokens > tailBudget {
			break
		}
		used += msgTokens
		selected = append([]storage.Message{history[i]}, selected...)
	}
	return repairToolCallBoundaries(selected)
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

func historyHash(history []storage.Message) string {
	data, err := json.Marshal(history)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

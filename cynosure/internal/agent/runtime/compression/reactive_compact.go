package compression

import (
	"context"
	"fmt"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
)

const reactiveCompactStrategyName = "reactive_compact"

const (
	// reactiveSummaryTargetTokens is the summary budget for reactive compaction;
	// more aggressive than FullHistorySummarization's 8K.
	reactiveSummaryTargetTokens = 2 * 1024
	// reactiveTailMinTokens caps how much recent history is kept verbatim;
	// more aggressive than FullHistorySummarization's 16K.
	reactiveTailMinTokens = 4 * 1024
)

// ReactiveCompactStrategy is invoked out-of-band when the LLM rejects a request
// with HTTP 413 (context overflow). It is structurally similar to
// FullHistorySummarizationStrategy but more aggressive: it compacts
// unconditionally (no budget threshold check), targets a smaller summary, keeps
// a shorter recent tail, and drops reasoning_content from kept messages.
type ReactiveCompactStrategy struct{}

func (s *ReactiveCompactStrategy) Name() string {
	return reactiveCompactStrategyName
}

func (s *ReactiveCompactStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Summarizer == nil {
		return nil
	}
	// The last message is always preserved verbatim; only the earlier history
	// is summarized. With a single message there is nothing to summarize.
	if len(req.RequestHistory) < 2 {
		return nil
	}

	before := req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools)
	lastMessage := req.RequestHistory[len(req.RequestHistory)-1]
	summarizable := req.RequestHistory[:len(req.RequestHistory)-1]
	sourceHash := historyHash(summarizable)

	summary := ""
	if cached, err := req.Store.GetContextSummaryByHistoryHash(ctx, req.Conversation.ID, req.User.ID, sourceHash); err == nil && cached.Summary != "" {
		summary = cached.Summary
	} else {
		result, err := req.Summarizer(ctx, SummaryRequest{
			Conversation: req.Conversation,
			User:         req.User,
			History:      summarizable,
			TargetTokens: reactiveSummaryTargetTokens,
		})
		if err != nil {
			return fmt.Errorf("reactive compact summarize history: %w", err)
		}
		summary = result.Summary
		after := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
		_ = req.Store.CreateContextSummary(ctx, storage.ContextSummary{
			ID:                    "cs_" + idgen.Hex(),
			ConversationID:        req.Conversation.ID,
			UserID:                req.User.ID,
			SourceHistorySHA256:   sourceHash,
			Strategy:              reactiveCompactStrategyName,
			EstimatedTokensBefore: before,
			EstimatedTokensAfter:  after,
			Summary:               summary,
		})
	}

	tail := selectReactiveTail(req, summarizable, summary)
	kept := repairToolCallBoundaries(append(append([]storage.Message{}, tail...), lastMessage))
	dropReasoningContent(kept)
	req.RequestHistory = buildSummaryHistory(summary, kept)
	return nil
}

// selectReactiveTail picks messages from the end of the summarizable history
// within an aggressive tail budget reserved after the summary.
func selectReactiveTail(req *Request, summarizable []storage.Message, summary string) []storage.Message {
	summaryTokens := req.Estimator.EstimateRequestTokens(req.SystemPrompt, buildSummaryHistory(summary, nil), req.Tools)
	tailBudget := req.Estimator.ContextTokenBudget() - summaryTokens
	if tailBudget > reactiveTailMinTokens {
		tailBudget = reactiveTailMinTokens
	}
	if tailBudget <= 0 {
		return nil
	}

	var selected []storage.Message
	used := 0
	for i := len(summarizable) - 1; i >= 0; i-- {
		msgTokens := estimateTokensFromBytes(messageBytes(summarizable[i]))
		if used+msgTokens > tailBudget {
			break
		}
		used += msgTokens
		selected = append([]storage.Message{summarizable[i]}, selected...)
	}
	return selected
}

// dropReasoningContent strips reasoning_content from kept messages in place to
// reclaim context space; it is never needed for continuation.
func dropReasoningContent(messages []storage.Message) {
	for i := range messages {
		messages[i].ReasoningContent = ""
	}
}

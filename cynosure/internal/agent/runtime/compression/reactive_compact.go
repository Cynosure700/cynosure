package compression

import (
	"context"
	"fmt"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

const reactiveCompactStrategyName = "reactive_compact"

const (
	// reactiveSummaryTargetTokens 是被动压缩的摘要预算；
	// 比 FullHistorySummarization 的 8K 更激进。
	reactiveSummaryTargetTokens = 2 * 1000
	// reactiveTailMinTokens 限制逐字保留的近期历史量；
	// 比 FullHistorySummarization 的 16K 更激进。
	reactiveTailMinTokens = 4 * 1000
)

// ReactiveCompactStrategy 在 LLM 以 HTTP 413（上下文溢出）拒绝请求时被带外调用。
// 它在结构上与 FullHistorySummarizationStrategy 类似，但更激进：无条件压缩
// （不做预算阈值检查）、目标摘要更小、保留更短的近期尾部，并从保留的消息中丢弃 reasoning_content。
type ReactiveCompactStrategy struct{}

func (s *ReactiveCompactStrategy) Name() string {
	return reactiveCompactStrategyName
}

func (s *ReactiveCompactStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Summarizer == nil {
		return nil
	}
	// 最后一条消息始终逐字保留；只对更早的历史做摘要。
	// 仅有一条消息时无内容可摘要。
	if len(req.RequestHistory) < 2 {
		return nil
	}

	lastMessage := req.RequestHistory[len(req.RequestHistory)-1]
	summarizable := req.RequestHistory[:len(req.RequestHistory)-1]

	result, err := req.Summarizer(ctx, SummaryRequest{
		Conversation: req.Conversation,
		User:         req.User,
		History:      summarizable,
		TargetTokens: reactiveSummaryTargetTokens,
	})
	if err != nil {
		return fmt.Errorf("reactive compact summarize history: %w", err)
	}
	summary := result.Summary

	tail := selectReactiveTail(req, summarizable, summary)
	kept := repairToolCallBoundaries(append(append([]storage.Message{}, tail...), lastMessage))
	dropReasoningContent(kept)
	req.RequestHistory = buildSummaryHistory(summary, kept)
	return nil
}

// selectReactiveTail 在摘要之后预留的激进尾部预算范围内，
// 从可摘要历史的末尾向前挑选消息。
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

// dropReasoningContent 就地从保留的消息中剥除 reasoning_content 以回收上下文空间；
// 它对于续写从不需要。
func dropReasoningContent(messages []storage.Message) {
	for i := range messages {
		messages[i].ReasoningContent = ""
	}
}

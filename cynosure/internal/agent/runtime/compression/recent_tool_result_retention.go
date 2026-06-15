package compression

import (
	"context"

	"nano_cc/internal/textutil"
)

const recentToolResultRetentionStrategyName = "recent_tool_result_retention"

// RecentToolResultRetentionStrategy keeps only the most recent N tool results
// fully inline; earlier tool results are replaced with a one-line placeholder.
type RecentToolResultRetentionStrategy struct{}

func (s *RecentToolResultRetentionStrategy) Name() string {
	return recentToolResultRetentionStrategyName
}

func (s *RecentToolResultRetentionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory

	// Collect tool message indexes in order.
	var toolIndexes []int
	for i := range history {
		if history[i].Role == "tool" {
			toolIndexes = append(toolIndexes, i)
		}
	}
	if len(toolIndexes) <= recentToolResultRetention {
		return nil
	}

	// The last N tool messages (by position) keep their full result.
	cutoff := len(toolIndexes) - recentToolResultRetention
	for _, idx := range toolIndexes[:cutoff] {
		status, result, isJSON := textutil.ParseToolResult(history[idx].Content)
		if isCompactedResult(result) {
			continue
		}
		history[idx].Content = rebuildToolResult(status, earlierToolResultPlaceholder, isJSON)
	}
	return nil
}

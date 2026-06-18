package compression

import (
	"context"

	"nano_cc/internal/textutil"
)

const recentToolResultRetentionStrategyName = "recent_tool_result_retention"

// RecentToolResultRetentionStrategy keeps only the most recent N full inline
// tool results once their count exceeds the micro compaction threshold; earlier
// full inline tool results are replaced with a one-line placeholder.
type RecentToolResultRetentionStrategy struct{}

func (s *RecentToolResultRetentionStrategy) Name() string {
	return recentToolResultRetentionStrategyName
}

func (s *RecentToolResultRetentionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory

	type inlineToolResult struct {
		index  int
		status string
		isJSON bool
	}

	var candidates []inlineToolResult
	for i := range history {
		if history[i].Role != "tool" {
			continue
		}
		status, result, isJSON := textutil.ParseToolResult(history[i].Content)
		if isCompactedResult(result) {
			continue
		}
		candidates = append(candidates, inlineToolResult{index: i, status: status, isJSON: isJSON})
	}
	if len(candidates) <= recentToolResultRetentionThreshold {
		return nil
	}

	cutoff := len(candidates) - recentToolResultRetention
	for _, candidate := range candidates[:cutoff] {
		history[candidate.index].Content = rebuildToolResult(candidate.status, earlierToolResultPlaceholder, candidate.isJSON)
	}
	return nil
}

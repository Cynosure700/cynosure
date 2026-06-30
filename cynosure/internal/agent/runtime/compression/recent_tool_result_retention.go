package compression

import (
	"context"

	"github.com/Cynosure700/cynosure/cynosure/internal/textutil"
)

const recentToolResultRetentionStrategyName = "recent_tool_result_retention"

// RecentToolResultRetentionStrategy 在完整内联的 tool 结果数量超过微压缩阈值后，
// 仅保留最近的 N 个；更早的完整内联 tool 结果会被替换为一行占位符。
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

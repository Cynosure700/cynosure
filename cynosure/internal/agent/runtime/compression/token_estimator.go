package compression

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
)

const (
	// 128KB 上下文限制
	defaultModelContextLimit = 200 * 1024
	// 8KB 最大响应 token 数
	defaultMaxResponseTokens = 8 * 1024
	// 8KB 安全余量
	defaultSafetyMargin = 8 * 1024
)

// TokenEstimator estimates the token footprint of an outgoing request.
type TokenEstimator interface {
	EstimateRequestTokens(systemPrompt string, history []storage.Message, tools []openai.Tool) int
	ContextTokenBudget() int
}

// DefaultTokenEstimator uses a conservative ceil(utf8Bytes/3) approximation.
type DefaultTokenEstimator struct{}

func (DefaultTokenEstimator) ContextTokenBudget() int {
	return defaultModelContextLimit - defaultMaxResponseTokens - defaultSafetyMargin
}

func (e DefaultTokenEstimator) EstimateRequestTokens(systemPrompt string, history []storage.Message, tools []openai.Tool) int {
	bytes := len([]byte(systemPrompt))
	for _, msg := range history {
		if data, err := json.Marshal(msg); err == nil {
			bytes += len(data)
		} else {
			bytes += len(msg.Content) + len(msg.ReasoningContent)
		}
	}
	if len(tools) > 0 {
		if data, err := json.Marshal(tools); err == nil {
			bytes += len(data)
		}
	}
	return estimateTokensFromBytes(bytes)
}

func estimateTokensFromBytes(bytes int) int {
	if bytes <= 0 {
		return 0
	}
	return (bytes + 2) / 3
}

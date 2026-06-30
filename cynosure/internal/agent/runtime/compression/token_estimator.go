package compression

import (
	"encoding/json"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/storage"
)

const (
	// 默认上下文限制（200K）：未命中大上下文模型时采用。
	defaultModelContextLimit = 200 * 1000
	// 大上下文模型（如 deepseek-v4-flash / deepseek-v4-pro / glm-5.2）的上下文限制（1M）。
	largeModelContextLimit = 1000 * 1000
	// 8K 最大响应 token 数
	defaultMaxResponseTokens = 8 * 1000
	// 8K 安全余量
	defaultSafetyMargin = 8 * 1000
)

// largeContextModels 列出上下文限制为 1M 的模型 ID（小写归一化后匹配）。
var largeContextModels = map[string]struct{}{
	"deepseek-v4-flash": {},
	"deepseek-v4-pro":   {},
	"glm-5.2":           {},
}

// ModelContextLimit 返回指定模型的最大上下文 token 限制。
// deepseek-v4-flash、deepseek-v4-pro、glm-5.2 为 1M，其余为 200K。
func ModelContextLimit(modelID string) int {
	if _, ok := largeContextModels[strings.ToLower(strings.TrimSpace(modelID))]; ok {
		return largeModelContextLimit
	}
	return defaultModelContextLimit
}

// TokenEstimator 估算一个外发请求的 token 占用量。
type TokenEstimator interface {
	EstimateRequestTokens(systemPrompt string, history []storage.Message, tools []openai.Tool) int
	ContextTokenBudget() int
}

// DefaultTokenEstimator 使用保守的 ceil(utf8Bytes/3) 近似估算。
// ModelID 决定上下文限制；为空时回退到默认（200K）限制。
type DefaultTokenEstimator struct {
	ModelID string
}

func (e DefaultTokenEstimator) ContextTokenBudget() int {
	return ModelContextLimit(e.ModelID) - defaultMaxResponseTokens - defaultSafetyMargin
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

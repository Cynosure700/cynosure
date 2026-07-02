package compression

import (
	"encoding/json"
	"strings"
	"unicode/utf8"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

const (
	tokenUnitK = 1024
	tokenUnitM = tokenUnitK * tokenUnitK

	// 默认上下文限制（200K）：未命中大上下文模型时采用。
	defaultModelContextLimit = 200 * tokenUnitK
	// 大上下文模型（如 deepseek-v4-flash / deepseek-v4-pro / glm-5.2）的上下文限制（1M）。
	largeModelContextLimit = tokenUnitM
	// 8K 最大响应 token 数
	defaultMaxResponseTokens = 8 * tokenUnitK
	// 8K 安全余量
	defaultSafetyMargin = 8 * tokenUnitK
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

// DefaultTokenEstimator 使用 ceil(字符数/4) 近似估算。
// ModelID 决定上下文限制；为空时回退到默认（200K）限制。
type DefaultTokenEstimator struct {
	ModelID string
}

func (e DefaultTokenEstimator) ContextTokenBudget() int {
	return ModelContextLimit(e.ModelID) - defaultMaxResponseTokens - defaultSafetyMargin
}

func (e DefaultTokenEstimator) EstimateRequestTokens(systemPrompt string, history []storage.Message, tools []openai.Tool) int {
	chars := utf8.RuneCountInString(systemPrompt)
	for _, msg := range history {
		if data, err := json.Marshal(msg); err == nil {
			chars += utf8.RuneCount(data)
		} else {
			chars += utf8.RuneCountInString(msg.Content) + utf8.RuneCountInString(msg.ReasoningContent)
		}
	}
	if len(tools) > 0 {
		if data, err := json.Marshal(tools); err == nil {
			chars += utf8.RuneCount(data)
		}
	}
	return estimateTokensFromChars(chars)
}

func estimateTokensFromChars(chars int) int {
	if chars <= 0 {
		return 0
	}
	return (chars + 3) / 4
}

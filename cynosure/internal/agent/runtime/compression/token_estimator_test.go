package compression

import "testing"

func TestModelContextLimit(t *testing.T) {
	cases := []struct {
		modelID string
		want    int
	}{
		{"deepseek-v4-flash", largeModelContextLimit},
		{"deepseek-v4-pro", largeModelContextLimit},
		{"glm-5.2", largeModelContextLimit},
		// 归一化：大小写与首尾空白不影响匹配。
		{"  GLM-5.2 ", largeModelContextLimit},
		{"DeepSeek-V4-Flash", largeModelContextLimit},
		// 其它模型回退到默认 200K。
		{"gpt-4", defaultModelContextLimit},
		{"deepseek-v3", defaultModelContextLimit},
		{"", defaultModelContextLimit},
	}
	for _, tc := range cases {
		if got := ModelContextLimit(tc.modelID); got != tc.want {
			t.Errorf("ModelContextLimit(%q) = %d, want %d", tc.modelID, got, tc.want)
		}
	}
}

func TestDefaultTokenEstimatorContextTokenBudget(t *testing.T) {
	largeBudget := largeModelContextLimit - defaultMaxResponseTokens - defaultSafetyMargin
	if got := (DefaultTokenEstimator{ModelID: "glm-5.2"}).ContextTokenBudget(); got != largeBudget {
		t.Errorf("large-model budget = %d, want %d", got, largeBudget)
	}
	defaultBudget := defaultModelContextLimit - defaultMaxResponseTokens - defaultSafetyMargin
	if got := (DefaultTokenEstimator{}).ContextTokenBudget(); got != defaultBudget {
		t.Errorf("default budget = %d, want %d", got, defaultBudget)
	}
}

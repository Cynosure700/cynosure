package compression

import "testing"

func TestModelContextLimit(t *testing.T) {
	cases := []struct {
		modelID string
		want    int
	}{
		{"deepseek-v4-flash", 1000 * 1000},
		{"deepseek-v4-pro", 1000 * 1000},
		{"glm-5.2", 1000 * 1000},
		// 归一化：大小写与首尾空白不影响匹配。
		{"  GLM-5.2 ", 1000 * 1000},
		{"DeepSeek-V4-Flash", 1000 * 1000},
		// 其它模型回退到默认 256K。
		{"gpt-4", 256 * 1000},
		{"deepseek-v3", 256 * 1000},
		{"", 256 * 1000},
	}
	for _, tc := range cases {
		if got := ModelContextLimit(tc.modelID); got != tc.want {
			t.Errorf("ModelContextLimit(%q) = %d, want %d", tc.modelID, got, tc.want)
		}
	}
}

func TestDefaultTokenEstimatorContextTokenBudget(t *testing.T) {
	largeBudget := 1000*1000 - 8*1000 - 8*1000
	if got := (DefaultTokenEstimator{ModelID: "glm-5.2"}).ContextTokenBudget(); got != largeBudget {
		t.Errorf("large-model budget = %d, want %d", got, largeBudget)
	}
	defaultBudget := 256*1000 - 8*1000 - 8*1000
	if got := (DefaultTokenEstimator{}).ContextTokenBudget(); got != defaultBudget {
		t.Errorf("default budget = %d, want %d", got, defaultBudget)
	}
}

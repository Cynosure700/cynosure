package runtime

import "testing"

func TestOutputTokenBudgetsUseBinaryK(t *testing.T) {
	if defaultMaxTokens != 8*1024 {
		t.Fatalf("defaultMaxTokens = %d, want 8K tokens in binary units", defaultMaxTokens)
	}
	if truncationMaxTokens != 64*1024 {
		t.Fatalf("truncationMaxTokens = %d, want 64K tokens in binary units", truncationMaxTokens)
	}
}

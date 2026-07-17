package runtime

import "testing"

func TestOutputTokenBudgetsUseDecimalK(t *testing.T) {
	if defaultMaxTokens != 8*1000 {
		t.Fatalf("defaultMaxTokens = %d, want 8K tokens in decimal units", defaultMaxTokens)
	}
	if truncationMaxTokens != 64*1000 {
		t.Fatalf("truncationMaxTokens = %d, want 64K tokens in decimal units", truncationMaxTokens)
	}
}

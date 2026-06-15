package tools

import (
	"strings"
	"testing"
)

func TestLoadSkillToolDescriptionRequiresExactNameBeforeUse(t *testing.T) {
	for _, tool := range AllToolDefs {
		if tool.Function == nil || tool.Function.Name != "load_skill" {
			continue
		}
		description := tool.Function.Description
		for _, want := range []string{"full instructions", "exact name", "before using or following"} {
			if !strings.Contains(description, want) {
				t.Fatalf("expected load_skill description to contain %q, got %q", want, description)
			}
		}
		if strings.Contains(strings.ToLower(description), "database") {
			t.Fatalf("expected load_skill description to avoid database wording, got %q", description)
		}
		return
	}
	t.Fatalf("expected AllToolDefs to include load_skill")
}

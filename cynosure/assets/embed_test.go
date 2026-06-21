package assets

import (
	"strings"
	"testing"
)

func TestFunctionalPromptLoadsEmbeddedMarkdown(t *testing.T) {
	prompt, err := FunctionalPrompt("memory_extraction")
	if err != nil {
		t.Fatalf("load functional prompt: %v", err)
	}
	if strings.TrimSpace(prompt) != prompt {
		t.Fatalf("expected prompt to be trimmed, got %q", prompt)
	}
	for _, want := range []string{
		"project-scoped long-term memory extraction engine",
		"ONLY for the current project",
		`Output ONLY a JSON array: [{"name","type","description","body"}].`,
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected prompt to contain %q, got %q", want, prompt)
		}
	}
}

func TestFunctionalPromptLoadsSubagentTemplates(t *testing.T) {
	for name, want := range map[string]string{
		"general_subagent": "general 子智能体",
		"explore_subagent": "{{current_working_directory}}",
	} {
		prompt, err := FunctionalPrompt(name)
		if err != nil {
			t.Fatalf("load %s prompt: %v", name, err)
		}
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected %s prompt to contain %q, got %q", name, want, prompt)
		}
	}
}

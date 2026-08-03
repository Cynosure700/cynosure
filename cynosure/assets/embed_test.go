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

func TestSystemPromptDescribesUpdateMemoryAsExistingMemoryOnly(t *testing.T) {
	prompt := SystemPrompt()
	want := "update_memory`（更新或修正已有记忆）"
	if !strings.Contains(prompt, want) {
		t.Fatalf("expected system prompt to contain %q", want)
	}
	if strings.Contains(prompt, "update_memory`（新增或修正记忆）") {
		t.Fatal("expected system prompt not to describe update_memory as creating memories")
	}
}

func TestFunctionalPromptLoadsSubagentTemplates(t *testing.T) {
	for name, want := range map[string]string{
		"general_subagent": "general 子智能体",
		"explore_subagent": "{{workspace_root}}",
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

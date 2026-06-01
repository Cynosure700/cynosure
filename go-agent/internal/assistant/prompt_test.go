package assistant

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadBaseSystemPromptFromFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "system_prompt.md")
	if err := os.WriteFile(path, []byte("  Custom base prompt.\n\n"), 0o644); err != nil {
		t.Fatalf("write prompt file: %v", err)
	}

	prompt, err := LoadBaseSystemPrompt(path)
	if err != nil {
		t.Fatalf("load base prompt: %v", err)
	}
	if prompt != "Custom base prompt." {
		t.Fatalf("expected trimmed prompt file content, got %q", prompt)
	}
}

func TestBuildSystemPromptUsesLoadedBasePromptAndAppendsDynamicSections(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{
		BasePrompt:        "Custom base prompt.",
		Surface:           "browser chat",
		WorkingDirectory:  "/workspace",
		ToolNames:         []string{"load_skill", "bash"},
		SkillDescriptions: "- demo: Demo skill",
		MemorySection:     "Remember user preference.",
	})

	for _, want := range []string{
		"Custom base prompt.",
		"You are responding inside browser chat.",
		"Current workspace root: /workspace.",
		"Runtime tools available in this conversation: load_skill, bash.",
		"Available skills:\n- demo: Demo skill",
		"Remember user preference.",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected prompt to contain %q, got %q", want, prompt)
		}
	}
	if strings.Contains(prompt, "You are nano_cc, a general-purpose agent") {
		t.Fatalf("expected custom base prompt to replace compiled default, got %q", prompt)
	}
}

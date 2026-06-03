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
		SkillDescriptions: "<skills>\n<skill>\n<name>demo</name>\n<description>Demo skill</description>\n</skill>\n</skills>",
		MemorySection:     "Remember user preference.",
	})

	for _, want := range []string{
		"# System Instructions",
		"## Identity",
		"Custom base prompt.",
		"---",
		"## Runtime Context",
		"Surface: browser chat",
		"Workspace root: /workspace",
		"Use the workspace root as the default working directory for runtime file and shell operations unless the runtime tells you otherwise.",
		"## Runtime Tools",
		"The following tools are available in this conversation:",
		"- load_skill",
		"- bash",
		"## Skills",
		"The following skills are available as summaries only.",
		"Before using or following a skill, call `load_skill` with the exact skill name to load its full instructions.",
		"Do not infer the full workflow from the summary alone.",
		"Available skills:\n\n<skills>\n<skill>\n<name>demo</name>\n<description>Demo skill</description>\n</skill>\n</skills>",
		"## Memory",
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

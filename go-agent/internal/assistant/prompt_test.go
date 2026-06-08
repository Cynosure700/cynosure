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
		"<identity>",
		"Custom base prompt.",
		"</identity>",
		"<workspace>",
		"Surface: browser chat",
		"Working directory: /workspace",
		"除非运行时另有说明，默认以工作目录作为运行时文件与 Shell 操作的根目录。",
		"</workspace>",
		"<tools>",
		"本次会话可用的工具如下：",
		"- load_skill",
		"- bash",
		"</tools>",
		"<skills>",
		"以下技能只提供摘要。",
		"使用或遵循某个技能前，先用 `load_skill` 以精确的技能名加载其完整说明。",
		"不要仅凭摘要臆测完整的工作流。",
		"可用技能：\n\n<skills>\n<skill>\n<name>demo</name>\n<description>Demo skill</description>\n</skill>\n</skills>",
		"<memory>",
		"Remember user preference.",
		"</memory>",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected prompt to contain %q, got %q", want, prompt)
		}
	}
	if strings.Contains(prompt, "你是 nano_cc") {
		t.Fatalf("expected custom base prompt to replace compiled default, got %q", prompt)
	}
}

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
		BasePrompt:       "Custom base prompt.",
		Surface:          "local TUI",
		WorkingDirectory: "/workspace",
		CynosureMarkdown: CynosureMarkdownContext{
			UserPath:         "/home/alice/.cynosure/CYNOSURE.MD",
			UserContent:      "# User Rule\n全局说明",
			WorkspacePath:    "/workspace/.cynosure/CYNOSURE.MD",
			WorkspaceContent: "# Project Rule\n项目说明",
		},
		ToolNames:         []string{"load_skill", "bash"},
		SkillDescriptions: "<skills>\n<skill>\n<name>demo</name>\n<description>Demo skill</description>\n</skill>\n</skills>",
		MemorySection:     "Remember user preference.",
	})

	for _, want := range []string{
		"<identity>",
		"Custom base prompt.",
		"</identity>",
		"<workspace>",
		"Surface: local TUI",
		"Working directory: /workspace",
		"除非运行时另有说明，默认以工作目录作为运行时文件与 Shell 操作的根目录。",
		"</workspace>",
		"<system-reminder>",
		"# linkMd",
		"/home/alice/.cynosure/CYNOSURE.MD 的内容（用户为所有项目配置的私人全局说明）：",
		"# User Rule\n全局说明",
		"/workspace/.cynosure/CYNOSURE.MD 的内容（项目说明，已提交到代码库或工作区）：",
		"# Project Rule\n项目说明",
		"重要：这些上下文可能与当前任务相关，也可能无关。除非与任务高度相关，否则不要对其作出回应。",
		"</system-reminder>",
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
	workspaceEnd := strings.Index(prompt, "</workspace>")
	reminderStart := strings.Index(prompt, "<system-reminder>")
	toolsStart := strings.Index(prompt, "<tools>")
	if !(workspaceEnd < reminderStart && reminderStart < toolsStart) {
		t.Fatalf("expected link reminder between workspace and tools, got %q", prompt)
	}
}

func TestBuildSystemPromptOmitsEmptyCynosureMarkdownContext(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(prompt, "# linkMd") || strings.Contains(prompt, "<system-reminder>") {
		t.Fatalf("expected empty link markdown context to be omitted, got %q", prompt)
	}
}

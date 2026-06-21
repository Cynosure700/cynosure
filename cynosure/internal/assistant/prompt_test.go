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
		"# cynosureMd",
		"用户全局说明：",
		"# User Rule\n全局说明",
		"项目说明：",
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
		"记忆（Memory）段可能同时包含过往记忆索引和真实有效记忆，两者边界必须严格区分。",
		"memory.md 只提供过往记忆文件索引，仅用于 update_memory/delete_memory 定位、更新或删除记忆文件；索引条目本身不作为任何有用信息",
		"只有标注为真实有效记忆、且来自具体记忆文件的内容，才可作为历史上下文参考。",
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
	if strings.Contains(prompt, "# cynosureMd") || strings.Contains(prompt, "<system-reminder>") {
		t.Fatalf("expected empty link markdown context to be omitted, got %q", prompt)
	}
}

func TestDefaultBaseSystemPromptUsesDomainSections(t *testing.T) {
	for _, want := range []string{
		"## 定义角色",
		"## 安全性声明",
		"## 帮助文档",
		"## 输出风格",
		"## 任务管理",
		"复杂任务、多步骤任务、用户提供多个目标或需要持续验证的任务，必须在开始执行前调用 todo_write",
		"每完成一个待办事项，必须立即调用 todo_write 将该项标记为 completed",
		"调用 todo_list 查询当前待办列表",
		"## 工具调用",
		"## 环境信息",
	} {
		if !strings.Contains(DefaultBaseSystemPrompt, want) {
			t.Fatalf("expected default base prompt to contain domain section %q", want)
		}
	}
	if strings.Contains(DefaultBaseSystemPrompt, "记忆（Memory）段可能同时包含过往记忆索引") {
		t.Fatalf("expected memory guidance to be injected dynamically, not embedded in default base prompt")
	}
}

func TestBuildSystemPromptOmitsMemoryGuidanceWhenMemorySectionIsEmpty(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(prompt, "<memory>") || strings.Contains(prompt, "记忆（Memory）段可能同时包含过往记忆索引") {
		t.Fatalf("expected empty memory section to omit dynamic memory guidance, got %q", prompt)
	}
}

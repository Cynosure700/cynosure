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
		GitStatus:         "This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.\n\nCurrent branch: main",
		CurrentDate:       "2026-06-24",
	})

	for _, want := range []string{
		"<identity>",
		"Custom base prompt.",
		"</identity>",
		"<workspace>",
		"Surface: local TUI",
		"Workspace root: /workspace",
		"以工作区根目录作为运行时文件与 Shell 操作的根目录：相对路径基于工作区根目录解析，绝对路径原样使用。",
		"</workspace>",
		"<system-reminder>",
		"current_day: 2026-06-24",
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
		"<git-status>",
		"This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.",
		"Current branch: main",
		"</git-status>",
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
	if strings.Contains(prompt, "你是 cynosure") {
		t.Fatalf("expected custom base prompt to replace compiled default, got %q", prompt)
	}
	workspaceEnd := strings.Index(prompt, "</workspace>")
	reminderStart := strings.Index(prompt, "<system-reminder>")
	toolsStart := strings.Index(prompt, "<tools>")
	if !(workspaceEnd < reminderStart && reminderStart < toolsStart) {
		t.Fatalf("expected link reminder between workspace and tools, got %q", prompt)
	}
	gitStatusStart := strings.Index(prompt, "<git-status>")
	gitStatusEnd := strings.Index(prompt, "</git-status>")
	memoryStart := strings.Index(prompt, "<memory>")
	if !(toolsStart < gitStatusStart && gitStatusStart < gitStatusEnd && gitStatusEnd < memoryStart) {
		t.Fatalf("expected git-status section after tools and before memory, got %q", prompt)
	}
}

func TestBuildSystemPromptOmitsEmptyGitStatus(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(prompt, "<git-status>") {
		t.Fatalf("expected empty git status to be omitted, got %q", prompt)
	}
}

func TestBuildSystemPromptRendersGitStatusWithoutMemory(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{
		BasePrompt: "Base prompt.",
		Surface:    "local TUI",
		GitStatus:  "Current branch: feature/x",
	})
	if !strings.Contains(prompt, "<git-status>") || !strings.Contains(prompt, "Current branch: feature/x") {
		t.Fatalf("expected git status section even without memory, got %q", prompt)
	}
	if strings.Contains(prompt, "<memory>") {
		t.Fatalf("expected no memory section when memory empty, got %q", prompt)
	}
}

func TestBuildSystemPromptOmitsEmptyCynosureMarkdownContext(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(prompt, "# cynosureMd") || strings.Contains(prompt, "<system-reminder>") {
		t.Fatalf("expected empty link markdown context to be omitted, got %q", prompt)
	}
}

func TestBuildSystemPromptRendersCurrentDayWithoutCynosureMarkdown(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{
		BasePrompt:  "Base prompt.",
		Surface:     "local TUI",
		CurrentDate: "2026-06-24",
	})
	if !strings.Contains(prompt, "<system-reminder>") {
		t.Fatalf("expected system-reminder section for current_day, got %q", prompt)
	}
	if !strings.Contains(prompt, "current_day: 2026-06-24") {
		t.Fatalf("expected current_day line, got %q", prompt)
	}
	if strings.Contains(prompt, "# cynosureMd") {
		t.Fatalf("expected no cynosureMd content when context empty, got %q", prompt)
	}
	workspaceEnd := strings.Index(prompt, "</workspace>")
	reminderStart := strings.Index(prompt, "<system-reminder>")
	if !(workspaceEnd < reminderStart) {
		t.Fatalf("expected system-reminder after workspace, got %q", prompt)
	}
}

func TestBuildSystemPromptOmitsSystemReminderWhenNoDateOrContext(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(prompt, "<system-reminder>") || strings.Contains(prompt, "current_day:") {
		t.Fatalf("expected no system-reminder when neither date nor context present, got %q", prompt)
	}
}

func TestBuildSystemPromptCurrentDayPrecedesCynosureMarkdown(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{
		BasePrompt:  "Base prompt.",
		Surface:     "local TUI",
		CurrentDate: "2026-06-24",
		CynosureMarkdown: CynosureMarkdownContext{
			WorkspacePath:    "/workspace/.cynosure/CYNOSURE.MD",
			WorkspaceContent: "# Project Rule",
		},
	})
	dayIdx := strings.Index(prompt, "current_day: 2026-06-24")
	mdIdx := strings.Index(prompt, "# cynosureMd")
	if dayIdx < 0 || mdIdx < 0 || dayIdx > mdIdx {
		t.Fatalf("expected current_day before cynosureMd content, got %q", prompt)
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
		"当用户的问题涉及项目代码、代码库、工程实现、构建测试或缺陷排查时，除非能够基于当前上下文直接、准确回答，否则必须先调用 todo_write 将工作拆成多步骤任务",
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

func TestDefaultBaseSystemPromptRequiresExploreForSearchSubagents(t *testing.T) {
	for _, want := range []string{
		"spawn_subagent 必须提供 sub_type 与 task",
		"搜索、文件定位、代码探索、实现梳理、证据收集等搜索相关任务必须使用 sub_type=explore",
		"sub_type=general 仅用于需要隔离上下文的综合分析或执行型子任务，不得用于搜索相关任务",
		"调用子智能体时，必须将任务拆成多个轻量级、边界清晰的子任务",
		"禁止把整个项目的探索任务交给单个子智能体",
	} {
		if !strings.Contains(DefaultBaseSystemPrompt, want) {
			t.Fatalf("expected default base prompt to contain %q", want)
		}
	}
}

func TestDefaultBaseSystemPromptGuidesFileAndSearchTools(t *testing.T) {
	for _, want := range []string{
		"read_file 可直接读取本地文件系统中的文件",
		"用户提供文件路径时，默认该路径有效",
		"读取不存在的用户提供路径是允许的",
		"除非路径由用户直接提供，否则必须先确认文件存在再读取",
		"write_file 会覆盖目标路径的既有文件",
		"写入既有文件前必须先使用 read_file 读取当前内容",
		"修改既有文件优先使用 edit_file 或 multi_edit",
		"除非用户明确要求，不要创建文档文件（*.md）或 README 文件",
		"edit_file 与 multi_edit 执行精确字符串替换",
		"替换从 read_file 输出复制的内容时，只匹配行号前缀之后的真实文件内容",
		"文件内容搜索必须优先使用 grep，不要用 bash 调用 grep 或 rg",
		"文件名模式匹配使用 glob，结果按修改时间排序",
	} {
		if !strings.Contains(DefaultBaseSystemPrompt, want) {
			t.Fatalf("expected default base prompt to contain %q", want)
		}
	}
	for _, forbidden := range []string{
		"read_file 只能读取已确认存在的普通文件",
		"用户提供文件路径时也必须先确认存在",
	} {
		if strings.Contains(DefaultBaseSystemPrompt, forbidden) {
			t.Fatalf("expected default base prompt to drop stale read_file restriction %q", forbidden)
		}
	}
}

func TestBuildSystemPromptOmitsMemoryGuidanceWhenMemorySectionIsEmpty(t *testing.T) {
	prompt := BuildSystemPrompt(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(prompt, "<memory>") || strings.Contains(prompt, "记忆（Memory）段可能同时包含过往记忆索引") {
		t.Fatalf("expected empty memory section to omit dynamic memory guidance, got %q", prompt)
	}
}

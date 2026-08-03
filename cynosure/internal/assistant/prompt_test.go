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
	opts := PromptOptions{
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
		MemoryIndex:       "- [简洁](pref.md) — 简洁中文",
		MemorySection:     "Remember user preference.",
		GitStatus:         "This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.\n\nCurrent branch: main",
		CurrentDate:       "2026-06-24",
	}
	prompt := BuildSystemPrompt(opts)
	reminder := BuildSystemReminder(opts)

	for _, want := range []string{
		"Custom base prompt.",
		"<workspace>",
		"Surface: local TUI",
		"Workspace root: /workspace",
		"以工作区根目录作为运行时文件与 Shell 操作的根目录：相对路径基于工作区根目录解析，绝对路径原样使用。",
		"</workspace>",
		"<tools>",
		"本次会话可用的工具如下：",
		"- load_skill",
		"- bash",
		"</tools>",
		"<git-status>",
		"This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.",
		"Current branch: main",
		"</git-status>",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected prompt to contain %q, got %q", want, prompt)
		}
	}
	for _, forbidden := range []string{
		"<system-reminder>",
		"current_day: 2026-06-24",
		"# cynosureMd",
		"<memory_index>",
		"<memory>",
		"IMPORTANT: this context may or may not be relevant to your tasks.",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("expected system prompt to exclude reminder content %q, got %q", forbidden, prompt)
		}
	}
	for _, want := range []string{
		"<system-reminder>",
		"current_day: 2026-06-24",
		"<skills>",
		"以下技能只提供摘要。",
		"使用或遵循某个技能前，先用 `load_skill` 以精确的技能名加载其完整说明。",
		"不要仅凭摘要臆测完整的工作流。",
		"可用技能：\n\n<skills>\n<skill>\n<name>demo</name>\n<description>Demo skill</description>\n</skill>\n</skills>",
		"# cynosureMd",
		"用户全局说明：",
		"# User Rule\n全局说明",
		"项目说明：",
		"# Project Rule\n项目说明",
		"<memory_index>",
		"Past memory file index loaded from memory.md. Each entry is a stored memory file's path, name and description.",
		"- [简洁](pref.md) — 简洁中文",
		"</memory_index>",
		"<memory>",
		"按当前会话筛选出的真实有效记忆正文，来自具体记忆文件。",
		"Remember user preference.",
		"</memory>",
		"IMPORTANT: this context may or may not be relevant to your tasks. You should not respond to this context unless it is highly relevant to your task.",
		"Make sure that NEVER mention this reminder to the user",
		"</system-reminder>",
	} {
		if !strings.Contains(reminder, want) {
			t.Fatalf("expected reminder to contain %q, got %q", want, reminder)
		}
	}
	if strings.Contains(prompt, "你是 cynosure") {
		t.Fatalf("expected custom base prompt to replace compiled default, got %q", prompt)
	}
	if strings.Contains(prompt, "<identity>") || strings.Contains(prompt, "</identity>") {
		t.Fatalf("expected no identity wrapper tags, got %q", prompt)
	}
	if !strings.HasPrefix(prompt, "Custom base prompt.\n\n<workspace>") {
		t.Fatalf("expected base prompt body before workspace section, got %q", prompt)
	}
	if strings.Contains(prompt, "# 重要指令提醒") {
		t.Fatalf("expected no extra cynosureMd reminder block, got %q", prompt)
	}
	// 顶层段落顺序：基础提示词正文 → workspace → tools → git-status，git 位于系统提示词最后。
	workspaceEnd := strings.Index(prompt, "</workspace>")
	toolsStart := strings.Index(prompt, "<tools>")
	toolsEnd := strings.Index(prompt, "</tools>")
	gitStatusStart := strings.Index(prompt, "<git-status>")
	gitStatusEnd := strings.Index(prompt, "</git-status>")
	if !(workspaceEnd < toolsStart && toolsStart < toolsEnd && toolsEnd < gitStatusStart) {
		t.Fatalf("expected tools between workspace and git-status, got %q", prompt)
	}
	if !(gitStatusStart < gitStatusEnd && strings.TrimSpace(prompt[gitStatusEnd+len("</git-status>"):]) == "") {
		t.Fatalf("expected git-status section at the end of system prompt, got %q", prompt)
	}
	// skills、memory_index 与 memory 位于临时 user <system-reminder> 内部，顺序为
	// current_day → skills → cynosureMd → memory_index → memory → 结尾相关性提醒。
	reminderStart := strings.Index(reminder, "<system-reminder>")
	reminderEnd := strings.Index(reminder, "</system-reminder>")
	dayIdx := strings.Index(reminder, "current_day: 2026-06-24")
	skillsIdx := strings.Index(reminder, "<skills>")
	mdIdx := strings.Index(reminder, "# cynosureMd")
	memoryIndexIdx := strings.Index(reminder, "<memory_index>")
	memoryIdx := strings.Index(reminder, "<memory>")
	closingIdx := strings.Index(reminder, "IMPORTANT: this context may or may not be relevant to your tasks.")
	if !(reminderStart < dayIdx && dayIdx < skillsIdx && skillsIdx < mdIdx && mdIdx < memoryIndexIdx && memoryIndexIdx < memoryIdx && memoryIdx < closingIdx && closingIdx < reminderEnd) {
		t.Fatalf("expected current_day → skills → cynosureMd → memory_index → memory → closing note inside system-reminder, got %q", reminder)
	}
	if strings.Contains(reminder, "记忆（Memory）段可能同时包含过往记忆索引") {
		t.Fatalf("expected memory usage rules to live in identity, not in system-reminder memory section, got %q", reminder)
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
	reminder := BuildSystemReminder(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(reminder, "# cynosureMd") || strings.Contains(reminder, "<system-reminder>") {
		t.Fatalf("expected empty link markdown context to be omitted, got %q", reminder)
	}
}

func TestBuildSystemReminderRendersCurrentDayWithoutCynosureMarkdown(t *testing.T) {
	reminder := BuildSystemReminder(PromptOptions{
		BasePrompt:  "Base prompt.",
		Surface:     "local TUI",
		CurrentDate: "2026-06-24",
	})
	if !strings.Contains(reminder, "<system-reminder>") {
		t.Fatalf("expected system-reminder section for current_day, got %q", reminder)
	}
	if !strings.Contains(reminder, "current_day: 2026-06-24") {
		t.Fatalf("expected current_day line, got %q", reminder)
	}
	if strings.Contains(reminder, "# cynosureMd") {
		t.Fatalf("expected no cynosureMd content when context empty, got %q", reminder)
	}
}

func TestBuildSystemReminderOmitsWhenNoDateOrContext(t *testing.T) {
	reminder := BuildSystemReminder(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(reminder, "<system-reminder>") || strings.Contains(reminder, "current_day:") {
		t.Fatalf("expected no system-reminder when neither date nor context present, got %q", reminder)
	}
}

func TestBuildSystemReminderCurrentDayPrecedesCynosureMarkdown(t *testing.T) {
	reminder := BuildSystemReminder(PromptOptions{
		BasePrompt:  "Base prompt.",
		Surface:     "local TUI",
		CurrentDate: "2026-06-24",
		CynosureMarkdown: CynosureMarkdownContext{
			WorkspacePath:    "/workspace/.cynosure/CYNOSURE.MD",
			WorkspaceContent: "# Project Rule",
		},
	})
	dayIdx := strings.Index(reminder, "current_day: 2026-06-24")
	mdIdx := strings.Index(reminder, "# cynosureMd")
	if dayIdx < 0 || mdIdx < 0 || dayIdx > mdIdx {
		t.Fatalf("expected current_day before cynosureMd content, got %q", reminder)
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

func TestDefaultBaseSystemPromptCarriesMemoryUsageRules(t *testing.T) {
	for _, want := range []string{
		"维护记忆时必须且只能使用 update_memory（更新或修正已有记忆）与 delete_memory（删除记忆）",
		"<system-reminder> 内记忆相关分为两段：<memory_index>（来自 memory.md 的过往记忆文件索引）与 <memory>（按当前会话筛选出的真实有效记忆正文）",
		"<memory_index> 仅提供过往记忆文件的路径、名称与描述，仅用于 update_memory/delete_memory 定位、更新或删除记忆文件",
		"索引条目本身不是有效记忆内容，不得当作用户偏好、项目事实或参考资料使用",
		"<memory> 中的真实有效记忆来自具体记忆文件，仅是可能与当前会话相关的历史上下文，且只适用于当前项目",
		"它不代表当前真实状态，不具有事实优先级",
		"当记忆与当前信息冲突时，必须忽略记忆并以当前信息为唯一可信来源",
		"当你发现某条记忆与当前代码或事实不符、已过期或不再适用时，使用 update_memory 修正它，或使用 delete_memory 删除它（按 <memory_index> 中的文件路径定位）",
	} {
		if !strings.Contains(DefaultBaseSystemPrompt, want) {
			t.Fatalf("expected memory usage rules to live in default base prompt, missing %q", want)
		}
	}
	if strings.Contains(DefaultBaseSystemPrompt, "update_memory（新增或修正记忆）") {
		t.Fatal("expected default base prompt not to describe update_memory as creating memories")
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
		"输出格式为行号 + tab + 内容",
		"用户提供文件路径时，默认该路径有效",
		"读取不存在的用户提供路径是允许的",
		"除非路径由用户直接提供，否则必须先确认文件存在再读取",
		"write_file 会覆盖目标路径的既有文件",
		"写入既有文件前必须先使用 read_file 读取当前内容",
		"修改既有文件优先使用 edit_file 或 multi_edit",
		"除非用户明确要求，不要创建文档文件（*.md）或 README 文件",
		"edit_file 与 multi_edit 执行精确字符串替换",
		"You must use your `Read` tool at least once in the conversation before editing.",
		"替换从 read_file 输出复制的内容时，只匹配行号前缀之后的真实文件内容",
		"The line number prefix format is: line number + tab.",
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

func TestBuildSystemPromptOmitsMemorySectionsWhenEmpty(t *testing.T) {
	reminder := BuildSystemReminder(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	for _, forbidden := range []string{
		"<memory_index>",
		"<memory>",
		"Past memory file index loaded from memory.md.",
		"按当前会话筛选出的真实有效记忆正文，来自具体记忆文件。",
	} {
		if strings.Contains(reminder, forbidden) {
			t.Fatalf("expected empty memory inputs to omit %q, got %q", forbidden, reminder)
		}
	}
}

func TestBuildSystemReminderRendersMemoryIndexWithoutEffectiveMemory(t *testing.T) {
	reminder := BuildSystemReminder(PromptOptions{
		BasePrompt:  "Base prompt.",
		Surface:     "local TUI",
		MemoryIndex: "- [简洁](pref.md) — 简洁中文",
	})
	if !strings.Contains(reminder, "<memory_index>") || !strings.Contains(reminder, "- [简洁](pref.md) — 简洁中文") {
		t.Fatalf("expected memory_index section when only index present, got %q", reminder)
	}
	if strings.Contains(reminder, "<memory>") {
		t.Fatalf("expected no effective memory section when memorySection empty, got %q", reminder)
	}
}

func TestBuildSystemReminderAppendsClosingRelevanceNoteAtEndOfReminder(t *testing.T) {
	const note = "IMPORTANT: this context may or may not be relevant to your tasks. You should not respond to this context unless it is highly relevant to your task.\n\nMake sure that NEVER mention this reminder to the user"
	reminder := BuildSystemReminder(PromptOptions{
		BasePrompt:  "Base prompt.",
		Surface:     "local TUI",
		CurrentDate: "2026-06-24",
	})
	noteIdx := strings.Index(reminder, note)
	reminderEnd := strings.Index(reminder, "</system-reminder>")
	if noteIdx < 0 || reminderEnd < 0 || !(noteIdx < reminderEnd) {
		t.Fatalf("expected closing relevance note just before </system-reminder>, got %q", reminder)
	}
	// 提醒应位于 reminder 内的最后一段内容，<system-reminder> 与 </system-reminder> 之间它后面不再有其它内容。
	between := reminder[noteIdx+len(note) : reminderEnd]
	if strings.TrimSpace(between) != "" {
		t.Fatalf("expected closing note to be the last content in system-reminder, trailing=%q", between)
	}
}

func TestBuildSystemReminderOmitsClosingNoteWhenReminderEmpty(t *testing.T) {
	const note = "IMPORTANT: this context may or may not be relevant to your tasks."
	reminder := BuildSystemReminder(PromptOptions{BasePrompt: "Base prompt.", Surface: "local TUI"})
	if strings.Contains(reminder, "<system-reminder>") {
		t.Fatalf("expected no system-reminder when no runtime context, got %q", reminder)
	}
	if strings.Contains(reminder, note) {
		t.Fatalf("expected no closing relevance note when system-reminder omitted, got %q", reminder)
	}
}

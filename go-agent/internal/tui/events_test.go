package tui

import (
	"context"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/sessions"
)

var ansiEscapePattern = regexp.MustCompile("\x1b\\[[0-9;]*m")
var ansiBackgroundPattern = regexp.MustCompile("\x1b\\[[0-9;]*48;[0-9;]*m")
var ansiBlueForegroundPattern = regexp.MustCompile("\x1b\\[[0-9;]*38;5;39m")
var ansiSelectedGreyBackgroundPattern = regexp.MustCompile("\x1b\\[[0-9;]*48;5;238[0-9;]*m")
var ansiReverseVideoPattern = regexp.MustCompile("\x1b\\[[0-9;]*7m")

func plainTerminalText(text string) string {
	return ansiEscapePattern.ReplaceAllString(text, "")
}

func TestAssistantFileTreeDoesNotRenderErrorBackgroundHighlighting(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.renderer = newMarkdownRenderer(app.messageWidth())
	content := "```go\ninternal/\n├── agent/      智能体核心逻辑（含 runtime/compression、runtime/hooks、mcp、storage）\n├── assistant/  对话助手层\n└── tui/        终端用户界面\n```"

	rendered := app.renderMessage(Message{Role: "assistant", Content: content})

	if ansiBackgroundPattern.MatchString(rendered) {
		t.Fatalf("assistant file tree render = %q, should not use background/error highlighting for plain file tree paths", rendered)
	}
	plain := plainTerminalText(rendered)
	for _, want := range []string{"internal/", "agent/", "runtime/compression", "assistant/", "tui/"} {
		if !strings.Contains(plain, want) {
			t.Fatalf("assistant file tree render = %q, want it to contain %q", plain, want)
		}
	}
}

func TestAssistantFileAndDirectoryReferencesRenderBlue(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.renderer = newMarkdownRenderer(app.messageWidth())
	content := "需要检查 go-agent/internal/tui/app.go 和 go-agent/internal/tui/ 目录。"

	rendered := app.renderMessage(Message{Role: "assistant", Content: content})

	for _, want := range []string{"go-agent/internal/tui/app.go", "go-agent/internal/tui/"} {
		if !strings.Contains(plainTerminalText(rendered), want) {
			t.Fatalf("assistant render = %q, want visible path %q", rendered, want)
		}
	}
	if got := len(ansiBlueForegroundPattern.FindAllString(rendered, -1)); got < 2 {
		t.Fatalf("assistant render = %q, want both file and directory references rendered blue", rendered)
	}
}

func TestToolMessageFileAndDirectoryReferencesKeepToolStyle(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:          "Read",
		ArgsPreview:   "file_path=go-agent/internal/tui/app.go",
		Status:        "success",
		ResultPreview: "listed go-agent/internal/tui/",
	}}

	rendered := app.renderMessage(msg)

	if !strings.Contains(plainTerminalText(rendered), "go-agent/internal/tui/app.go") || !strings.Contains(plainTerminalText(rendered), "go-agent/internal/tui/") {
		t.Fatalf("tool render = %q, want visible file and directory references", rendered)
	}
	if got := len(ansiBlueForegroundPattern.FindAllString(rendered, -1)); got != 1 {
		t.Fatalf("tool render = %q, want only the leading tool bullet rendered blue, got %d blue foreground sequences", rendered, got)
	}
}

func TestToolMessageRendersLeadingBlueBullet(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:          "Glob",
		ArgsPreview:   `path="/Users/bytedance/golang_pro/nano_cc", pattern="**/README*"`,
		Status:        "success",
		ResultPreview: "Found 1 file(s).",
	}}

	rendered := app.renderMessage(msg)
	plain := plainTerminalText(rendered)

	if !strings.Contains(plain, "● ✓ Glob") {
		t.Fatalf("tool render = %q, want tool call line prefixed with a small bullet", rendered)
	}
	if !strings.Contains(rendered, ansiForeground(tuiPalette.blue)+"●") {
		t.Fatalf("tool render = %q, want leading bullet rendered blue", rendered)
	}
}

func TestWorkspaceDirectoryRendersBlueInHeader(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", ModelID: "deepseek-v4-flash"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)

	header := model.renderHeader()

	if !strings.Contains(plainTerminalText(header), "/tmp/project") {
		t.Fatalf("header = %q, want visible workspace directory", header)
	}
	if !ansiBlueForegroundPattern.MatchString(header) {
		t.Fatalf("header = %q, want workspace directory rendered blue", header)
	}
}

func TestTypedInputKeepsNormalDisplay(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.input.SetValue("检查 go-agent/internal/tui/app.go 后继续")

	rendered := app.renderInput()

	if !strings.Contains(plainTerminalText(rendered), "检查 go-agent/internal/tui/app.go 后继续") {
		t.Fatalf("input render = %q, want typed user input visible", rendered)
	}
	if ansiReverseVideoPattern.MatchString(rendered) {
		t.Fatalf("input render = %q, want active input to keep normal display without selected styling", rendered)
	}
	if ansiSelectedGreyBackgroundPattern.MatchString(rendered) {
		t.Fatalf("input render = %q, want active input without selected grey background", rendered)
	}
	if !strings.Contains(rendered, "\x1b[0;38;5;255m 后继续") {
		t.Fatalf("input render = %q, want text after a blue file reference to restore normal white input style", rendered)
	}
}

func TestSubmittedUserMessageRendersAsGreySelectedLine(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 56

	rendered := app.renderMessage(Message{Role: "user", Content: "检查 go-agent/internal/tui/app.go 后继续"})

	if !strings.Contains(plainTerminalText(rendered), "检查 go-agent/internal/tui/app.go 后继续") {
		t.Fatalf("user message render = %q, want submitted user input visible", rendered)
	}
	if !ansiSelectedGreyBackgroundPattern.MatchString(rendered) {
		t.Fatalf("user message render = %q, want submitted user input rendered with a grey selected background", rendered)
	}
	if ansiReverseVideoPattern.MatchString(rendered) {
		t.Fatalf("user message render = %q, want grey selected background instead of reverse video", rendered)
	}
	if !strings.Contains(rendered, "\x1b[48;5;238;38;5;255m 后继续") {
		t.Fatalf("user message render = %q, want text after a blue file reference to stay on the grey selected line", rendered)
	}
	line := strings.Split(plainTerminalText(rendered), "\n")[0]
	if got := lipgloss.Width(line); got != app.messageWidth() {
		t.Fatalf("selected user line width = %d, want full message width %d: %q", got, app.messageWidth(), line)
	}
}

func TestPathHighlightRestoresOuterMessageStyleAfterPath(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120

	rendered := app.renderMessage(Message{Role: "error", Content: "打开 go-agent/internal/tui/app.go 失败"})

	if !strings.Contains(plainTerminalText(rendered), "打开 go-agent/internal/tui/app.go 失败") {
		t.Fatalf("error render = %q, want visible error text", rendered)
	}
	if !strings.Contains(rendered, "\x1b[0;38;5;209m 失败") {
		t.Fatalf("error render = %q, want text after a blue file reference to restore error style", rendered)
	}
}

type fakeSessionResumer struct {
	sessions []storage.ResumableSession
	conv     storage.Conversation
	history  []storage.Message
}

func (f *fakeSessionResumer) ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error) {
	return f.sessions, nil
}

func (f *fakeSessionResumer) ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error) {
	return f.conv, f.history, nil
}

func TestEventWriterConvertsRuntimeEvents(t *testing.T) {
	ch := make(chan Event, 4)
	writer := NewEventWriter(ch)

	if err := writer.Event("assistant_delta", map[string]any{"content": "hello"}); err != nil {
		t.Fatalf("Event returned error: %v", err)
	}
	got := <-ch
	if got.Name != "assistant_delta" || got.Content != "hello" {
		t.Fatalf("event = %#v, want assistant delta with content", got)
	}
}

func TestHandleSlashCommand(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", SkillCount: 2, MCPToolCount: 3})
	if handled := app.handleSlashCommand("/cwd"); !handled {
		t.Fatal("/cwd was not handled")
	}
	if len(app.messages) == 0 || app.messages[len(app.messages)-1].Content != "当前工作区：/tmp/project" {
		t.Fatalf("messages = %#v, want cwd message", app.messages)
	}
}

func TestHandleSkillsCommandShowsSkillDetails(t *testing.T) {
	app := NewModel(nil, SessionInfo{Skills: []sessions.SkillSummary{{Name: "project-helper", Source: "workspace", Description: "Project helper", Path: "/project/.link/skills/project-helper/skill.md"}}})

	if handled := app.handleSlashCommand("/skills"); !handled {
		t.Fatal("/skills was not handled")
	}
	content := app.messages[len(app.messages)-1].Content
	for _, want := range []string{"已加载 Skills：1 个", "project-helper", "workspace", "Project helper", "/project/.link/skills/project-helper/skill.md"} {
		if !strings.Contains(content, want) {
			t.Fatalf("/skills output = %q, want it to contain %q", content, want)
		}
	}
}

func TestHandleMCPCommandShowsServerDetails(t *testing.T) {
	app := NewModel(nil, SessionInfo{MCPToolCount: 3, MCPServers: []mcp.ServerStatus{{Name: "filesystem", Scope: "workspace", Transport: "stdio", Command: "npx", Args: []string{"-y", "server"}, Enabled: true, Connected: true, ToolCount: 3}, {Name: "docs", Scope: "workspace", Transport: "sse", URL: "https://example.com/sse", Enabled: true, LastError: "connect timeout"}}})

	if handled := app.handleSlashCommand("/mcp"); !handled {
		t.Fatal("/mcp was not handled")
	}
	content := app.messages[len(app.messages)-1].Content
	for _, want := range []string{"MCP Servers：2 个，工具：3 个", "filesystem", "connected", "npx -y server", "docs", "failed", "connect timeout"} {
		if !strings.Contains(content, want) {
			t.Fatalf("/mcp output = %q, want it to contain %q", content, want)
		}
	}
}

func TestResumeCommandShowsCurrentWorkspaceSessions(t *testing.T) {
	resumer := &fakeSessionResumer{sessions: []storage.ResumableSession{{SessionID: "session-1", Title: "第一段会话", UpdatedAt: time.Date(2026, 6, 12, 10, 0, 0, 0, time.UTC), MessageCount: 2}}}
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", Resumer: resumer})

	if handled := app.handleSlashCommand("/resume"); !handled {
		t.Fatal("/resume was not handled")
	}
	content := app.messages[len(app.messages)-1].Content
	for _, want := range []string{"可恢复的历史会话", "1.", "第一段会话", "session-1", "消息:2"} {
		if !strings.Contains(content, want) {
			t.Fatalf("/resume output = %q, want it to contain %q", content, want)
		}
	}
	if !app.resumeSelecting || len(app.resumeCandidates) != 1 {
		t.Fatalf("resumeSelecting=%v candidates=%#v, want selecting with one candidate", app.resumeSelecting, app.resumeCandidates)
	}
}

func TestResumeSelectionRestoresConversationAndDisplayHistory(t *testing.T) {
	resumer := &fakeSessionResumer{
		sessions: []storage.ResumableSession{{SessionID: "session-1", Title: "第一段会话", UpdatedAt: time.Now(), MessageCount: 3}},
		conv:     storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: "local-user", Title: "第一段会话"},
		history:  []storage.Message{{Role: "user", Content: "hello"}, {Role: "tool", Content: "large result"}, {Role: "assistant", Content: "hi"}},
	}
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", User: storage.User{ID: "local-user"}, Resumer: resumer})
	app.handleSlashCommand("/resume")

	if handled := app.handleResumeSelection("1"); !handled {
		t.Fatal("resume selection was not handled")
	}
	if app.session.Conversation.SessionID != "session-1" {
		t.Fatalf("SessionID = %q, want session-1", app.session.Conversation.SessionID)
	}
	if app.resumeSelecting {
		t.Fatal("resumeSelecting should be false after successful restore")
	}
	rendered := app.renderMessages()
	for _, want := range []string{"hello", "hi"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered messages = %q, want %q", rendered, want)
		}
	}
	if strings.Contains(rendered, "large result") {
		t.Fatalf("rendered messages = %q, should not include tool result", rendered)
	}
}

func TestModelIgnoresStaleGenerationEvents(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 2
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "stale"})
	model := updated.(Model)
	if len(model.messages) != 0 {
		t.Fatalf("messages = %#v, want stale event ignored", model.messages)
	}
}

func TestModelKeepsWaitingAfterStaleGenerationEvent(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 2
	app.running = true
	app.events <- Event{Generation: 2, Name: "assistant_delta", Content: "fresh"}

	_, cmd := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "stale"})

	if cmd == nil {
		t.Fatal("expected stale event to keep waiting while current generation is still running")
	}
	got := cmd()
	if event, ok := got.(Event); !ok || event.Generation != 2 || event.Content != "fresh" {
		t.Fatalf("next message = %#v, want fresh event from current generation", got)
	}
}

func TestRespondSendsTerminalEventThroughEventQueue(t *testing.T) {
	app := NewModel(nil, SessionInfo{})

	msg := app.respond(context.Background(), "hello", 7)()

	if msg != nil {
		t.Fatalf("respond command returned %#v, want nil so waitEvent preserves queued event order", msg)
	}
	got := <-app.events
	if got.Generation != 7 || got.Name != "error" || !strings.Contains(got.Content, "runtime 未初始化") {
		t.Fatalf("queued event = %#v, want runtime error for generation 7", got)
	}
}

func TestModelDisplaysReasoningDeltasAsMutedAssistantThinking(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "reasoning_delta", Content: "先判断是否需要工具"})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())

	for _, want := range []string{"思考", "先判断是否需要工具"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered messages = %q, want it to contain %q", rendered, want)
		}
	}
}

func TestModelDisplaysOnlyCurrentAssistantReasoningWhileRunning(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.messages = []Message{
		{Role: "assistant", Content: "历史答案", ReasoningContent: "历史思考"},
		{Role: "user", Content: "继续"},
	}
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "reasoning_delta", Content: "当前思考"})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())

	if !strings.Contains(rendered, "当前思考") {
		t.Fatalf("rendered messages = %q, want current reasoning", rendered)
	}
	if strings.Contains(rendered, "历史思考") {
		t.Fatalf("rendered messages = %q, should keep historical reasoning collapsed while current reasoning streams", rendered)
	}
}

func TestModelDoesNotExpandHistoricalReasoningBeforeCurrentReasoningStarts(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.messages = []Message{
		{Role: "user", Content: "上一问"},
		{Role: "assistant", Content: "历史答案", ReasoningContent: "历史思考"},
		{Role: "user", Content: "新问题"},
	}
	app.running = true

	rendered := app.renderMessages()

	if strings.Contains(rendered, "历史思考") {
		t.Fatalf("rendered messages = %q, should keep historical reasoning collapsed before current reasoning starts", rendered)
	}
}

func TestModelAssistantFinalEventHidesReasoningAfterDoneAndKeepsMeta(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "reasoning_delta", Content: "分析路径"})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "assistant", Content: "完成", Data: map[string]any{"content": "完成", "reasoning_content": "分析路径", "tool_call_count": 2, "context_tokens": 45000, "context_budget": 100000}})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "done"})
	model := updated.(Model)

	rendered := model.renderMessages()
	if !strings.Contains(rendered, "完成") {
		t.Fatalf("rendered messages = %q, want final answer", rendered)
	}
	if strings.Contains(rendered, "分析路径") {
		t.Fatalf("rendered messages = %q, should hide reasoning after done", rendered)
	}
	view := plainTerminalText(model.View())
	for _, want := range []string{"工具 2", "上下文 45%"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
}

func TestModelMetaEventUpdatesLiveStatus(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "meta", Data: map[string]any{"tool_call_count": 3, "context_tokens": 72000, "context_budget": 100000}})
	model := updated.(Model)

	view := model.View()
	for _, want := range []string{"工具 3", "上下文 72%"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
}

func TestModelDisplaysToolCallLifecycle(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "tool_1",
		"tool_name":    "bash",
		"args_preview": "command: go test ./...",
		"status":       "running",
	}})
	model := updated.(Model)
	if len(model.messages) != 1 || model.messages[0].Role != "tool" || model.messages[0].ToolCall == nil {
		t.Fatalf("messages = %#v, want one tool message", model.messages)
	}
	rendered := model.renderMessages()
	for _, want := range []string{"⏺ Bash", "command: go test ./...", "⎿ running"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}

	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":   "tool_1",
		"tool_name":      "bash",
		"args_preview":   "command: go test ./...",
		"status":         "success",
		"result_preview": "ok nano_cc/internal/tui 0.42s",
	}})
	model = updated.(Model)
	if len(model.messages) != 1 {
		t.Fatalf("messages = %#v, want tool done to update existing message", model.messages)
	}
	rendered = plainTerminalText(model.renderMessages())
	for _, want := range []string{"✓ Bash", "⎿ success · ok nano_cc/internal/tui 0.42s"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
}

func TestModelAppendsToolDoneWhenStartWasMissing(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":   "tool_missing_start",
		"tool_name":      "read_file",
		"args_preview":   "file_path: /tmp/a.go",
		"status":         "rejected",
		"result_preview": "Error: outside workspace",
	}})
	model := updated.(Model)
	if len(model.messages) != 1 || model.messages[0].Role != "tool" || model.messages[0].ToolCall == nil {
		t.Fatalf("messages = %#v, want appended completed tool message", model.messages)
	}
	rendered := plainTerminalText(model.renderMessages())
	for _, want := range []string{"✗ Read", "⎿ rejected · Error: outside workspace"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
}

func TestModelShowsThinkingAndToolResultsUntilAssistantReplyStarts(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.appendMessage("user", "新问题")
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "reasoning_delta", Content: "需要先跑测试"})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "tool_1",
		"tool_name":    "bash",
		"args_preview": "command: go test ./...",
		"status":       "running",
	}})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":   "tool_1",
		"tool_name":      "bash",
		"args_preview":   "command: go test ./...",
		"status":         "success",
		"result_preview": "ok nano_cc/internal/tui 0.42s",
	}})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "reasoning_delta", Content: "再检查结果"})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())
	for _, want := range []string{"需要先跑测试", "再检查结果", "✓ Bash", "ok nano_cc/internal/tui 0.42s"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q visible before assistant reply starts", rendered, want)
		}
	}

	updated, _ = model.Update(Event{Generation: 1, Name: "assistant_delta", Content: "完成"})
	model = updated.(Model)
	rendered = plainTerminalText(model.renderMessages())
	for _, want := range []string{"新问题", "完成"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
	for _, forbidden := range []string{"需要先跑测试", "再检查结果", "✓ Bash", "ok nano_cc/internal/tui 0.42s"} {
		if strings.Contains(rendered, forbidden) {
			t.Fatalf("rendered = %q, should hide %q once assistant reply starts", rendered, forbidden)
		}
	}
}

func TestRunningModelShowsThinkingIndicatorAtTranscriptBottom(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.messages = []Message{
		{Role: "user", Content: "介绍一下当前项目"},
		{Role: "assistant", Content: "正在整理"},
	}
	app.running = true
	app.thinkingStartedAt = time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	app.thinkingNow = app.thinkingStartedAt

	rendered := plainTerminalText(app.renderMessages())

	if !strings.Contains(rendered, "* Thinking... (1s)") {
		t.Fatalf("rendered messages = %q, want Thinking indicator with initial elapsed second", rendered)
	}
	if strings.LastIndex(rendered, "* Thinking... (1s)") < strings.LastIndex(rendered, "正在整理") {
		t.Fatalf("rendered messages = %q, want Thinking indicator below the current assistant reply", rendered)
	}
}

func TestThinkingIndicatorUpdatesElapsedSecondsAndHidesWhenDone(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.messages = []Message{{Role: "user", Content: "hello"}}
	app.generation = 1
	app.running = true
	app.thinkingStartedAt = time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	app.thinkingNow = app.thinkingStartedAt

	updated, _ := app.Update(thinkingTickMsg{generation: 1, at: app.thinkingStartedAt.Add(2 * time.Second)})
	model := updated.(Model)
	if rendered := plainTerminalText(model.renderMessages()); !strings.Contains(rendered, "* Thinking... (2s)") {
		t.Fatalf("rendered messages = %q, want Thinking indicator to advance elapsed seconds", rendered)
	}

	updated, _ = model.Update(Event{Generation: 1, Name: "done"})
	model = updated.(Model)
	if rendered := plainTerminalText(model.renderMessages()); strings.Contains(rendered, "Thinking...") {
		t.Fatalf("rendered messages = %q, should hide Thinking indicator after done", rendered)
	}
}

func TestViewUsesClaudeLikeSeparatedTerminalRegions(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	app.width = 80
	app.height = 24
	app.viewport.Width = 80
	app.viewport.Height = app.viewportHeight()
	app.appendMessage("user", "hello")
	app.appendMessage("assistant", "你好")
	app.refreshViewport()

	view := plainTerminalText(app.View())
	for _, want := range []string{"go-agent", "/tmp/project", "› hello", "你好", "工具 0", "上下文 --", "╭", "╰"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
}

func TestErrorEventReleasesRunningStateForNextPrompt(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "error", Content: "network failed"})
	model := updated.(Model)

	if model.running {
		t.Fatal("model should stop running after an error so the next user message can be sent")
	}
}

func TestViewFitsWithinTerminalHeight(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)
	model.appendMessage("user", "hello")
	model.appendMessage("assistant", strings.Repeat("这一行回答会比较长，用来模拟实际模型输出。", 6))
	model.refreshViewport()

	if got := lipgloss.Height(model.View()); got > model.height {
		t.Fatalf("view height = %d, want <= terminal height %d so history is not clipped", got, model.height)
	}
}

func TestViewKeepsHeaderInTranscriptAndInputFixedLikeClaudeCode(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", SkillCount: 2, MCPToolCount: 3})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)
	model.appendMessage("user", "hello")
	model.appendMessage("assistant", "你好")
	model.refreshViewport()

	view := plainTerminalText(model.View())
	for _, want := range []string{"go-agent", "/tmp/project", "╭", "╰", "Enter 发送", "上下文 --"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want Claude-like separated region marker %q", view, want)
		}
	}
	if model.viewport.Height != model.height-model.inputAreaHeight() {
		t.Fatalf("viewport height = %d, want only fixed input area reserved", model.viewport.Height)
	}
	if lipgloss.Height(view) > model.height {
		t.Fatalf("view height = %d, want <= terminal height %d", lipgloss.Height(view), model.height)
	}
}

func TestHeaderShowsOnlyModelAndWorkspaceMetadata(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", ModelID: "deepseek-v4-flash", SkillCount: 2, MCPToolCount: 3})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)

	header := model.renderHeader()
	for _, want := range []string{"deepseek-v4-flash", "/tmp/project"} {
		if !strings.Contains(header, want) {
			t.Fatalf("header = %q, want metadata %q", header, want)
		}
	}
	for _, forbidden := range []string{"go-agent", "workspace ", "Welcome back!", "API Usage Billing", "Tips for getting started", "Project guide", "Skills 2", "MCP tools 3"} {
		if strings.Contains(header, forbidden) {
			t.Fatalf("header = %q, should not contain noisy metadata %q", header, forbidden)
		}
	}
}

func TestHeaderUsesGreenAccentAndCompactCenteredLinkVersionMascot(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", ModelID: "deepseek-v4-flash"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)

	header := model.renderHeader()
	for _, want := range []string{`Link version: 0.0.1`, `Welcome back`, `/^ ^\`, `/ 0 0 \`, `V\ Y /V`, `/ - \`, `|| (__V`} {
		if !strings.Contains(header, want) {
			t.Fatalf("header = %q, want compact Link version mascot part %q", header, want)
		}
	}
	for _, forbidden := range []string{"/\\_/\\", "( o.o )", "> ^ <", "Doggy Server", "DOGGY API", `^-----^`, `Q /`, `(___\\====`, "Ready ✓"} {
		if strings.Contains(header, forbidden) {
			t.Fatalf("header = %q, should not contain old/oversized mascot part %q", header, forbidden)
		}
	}
	plain := plainTerminalText(header)
	if got := lipgloss.Height(plain); got > 13 {
		t.Fatalf("header height = %d, want compact Link version header height <= 13", got)
	}
	if strings.Index(plain, "Welcome back") > strings.Index(plain, "model deepseek-v4-flash") {
		t.Fatalf("header = %q, want Welcome back above model line", header)
	}
	for _, line := range strings.Split(plain, "\n") {
		if strings.Contains(line, `Link version: 0.0.1`) {
			leftPadding := strings.Index(line, `Link version: 0.0.1`)
			if leftPadding < 30 {
				t.Fatalf("header mascot line = %q, want Link version content centered with substantial left padding", line)
			}
		}
	}
	if headerAccentColor() != tuiPalette.mint {
		t.Fatalf("header accent = %q, want green accent matching mint palette %q", headerAccentColor(), tuiPalette.mint)
	}
}

func TestHeaderBoxClosesWithinTerminalWidth(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/Users/bytedance/golang_pro/nano_cc/go-agent", ModelID: "deepseek-v4-flash"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)

	header := plainTerminalText(model.renderHeader())
	for _, line := range strings.Split(header, "\n") {
		if got := lipgloss.Width(line); got > model.width {
			t.Fatalf("header line width = %d, want <= terminal width %d: %q", got, model.width, line)
		}
	}
	if !strings.Contains(header, "╮") || !strings.Contains(header, "╯") {
		t.Fatalf("header = %q, want closed right-side rounded border", header)
	}
}

func TestAssistantMessageRendersContentWithoutGoAgentLead(t *testing.T) {
	app := NewModel(nil, SessionInfo{})

	rendered := app.renderMessage(Message{Role: "assistant", Content: "你好"})
	if strings.Contains(rendered, "go-agent") {
		t.Fatalf("assistant render = %q, should render answer directly without go-agent label", rendered)
	}
	if !strings.Contains(rendered, "你好") {
		t.Fatalf("assistant render = %q, want assistant content", rendered)
	}
}

func TestHeaderScrollsAwayWithTranscriptHistory(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", ModelID: "deepseek-v4-flash"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 16})
	model := updated.(Model)
	for i := 0; i < 20; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()

	view := plainTerminalText(model.View())
	for _, forbidden := range []string{"deepseek-v4-flash", "/tmp/project"} {
		if strings.Contains(view, forbidden) {
			t.Fatalf("view = %q, should let header metadata %q scroll away with transcript history", view, forbidden)
		}
	}
	if model.viewport.Height != model.height-model.inputAreaHeight() {
		t.Fatalf("viewport height = %d, want only input area reserved from terminal height", model.viewport.Height)
	}
}

func TestViewFitsVerySmallTerminalHeight(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", ModelID: "deepseek-v4-flash"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 5})
	model := updated.(Model)
	model.appendMessage("assistant", "reply")
	model.refreshViewport()

	if got := lipgloss.Height(model.View()); got > model.height {
		t.Fatalf("view height = %d, want <= tiny terminal height %d", got, model.height)
	}
}

func TestHeaderKeepsLongWorkspaceOnOneLine(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/a/very/long/workspace/path/that/should/not/wrap/the/header", SkillCount: 2, MCPToolCount: 3})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 16})
	model := updated.(Model)

	header := model.renderHeader()
	if !strings.Contains(header, "…") {
		t.Fatalf("header = %q, want long workspace truncated with ellipsis", header)
	}
}

func TestViewKeepsWelcomeAndMessagesInScrollableTranscriptWithFixedPrompt(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)
	model.appendMessage("user", "hello")
	model.appendMessage("assistant", "你好")
	model.refreshViewport()

	view := plainTerminalText(model.View())
	for _, want := range []string{"› hello", "go-agent", "你好", "问 go-agent 一件事"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
	if strings.Index(view, "› hello") > strings.Index(view, "你好") {
		t.Fatalf("view = %q, want assistant answer below the submitted user prompt", view)
	}
	if !strings.Contains(model.renderInputArea(), "问 go-agent 一件事") {
		t.Fatalf("input area = %q, want fixed prompt to contain placeholder", model.renderInputArea())
	}
	if strings.Contains(model.renderTranscript(), "问 go-agent 一件事") {
		t.Fatalf("transcript = %q, want active prompt outside scrollable history", model.renderTranscript())
	}
}

func TestViewportScrollsTranscriptToActivePrompt(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 8})
	model := updated.(Model)
	for i := 0; i < 10; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()

	view := model.View()
	if !strings.Contains(view, "问 go-agent 一件事") {
		t.Fatalf("view = %q, want active prompt visible after scrolling to bottom", view)
	}
	if lipgloss.Height(view) > model.height {
		t.Fatalf("view height = %d, want <= terminal height %d", lipgloss.Height(view), model.height)
	}
}

func TestPageUpScrollsTranscriptHistoryWhileKeepingPromptFixed(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 8})
	model := updated.(Model)
	for i := 0; i < 10; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()
	bottomOffset := model.viewport.YOffset

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyPgUp})
	model = updated.(Model)

	if model.viewport.YOffset >= bottomOffset {
		t.Fatalf("viewport offset = %d, want less than bottom offset %d after page up", model.viewport.YOffset, bottomOffset)
	}
	if !strings.Contains(model.View(), "问 go-agent 一件事") {
		t.Fatalf("view = %q, want fixed prompt to remain visible after page up", model.View())
	}
}

func TestMouseWheelUpScrollsTranscriptHistory(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 8})
	model := updated.(Model)
	for i := 0; i < 10; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()
	bottomOffset := model.viewport.YOffset

	updated, _ = model.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp})
	model = updated.(Model)

	if model.viewport.YOffset >= bottomOffset {
		t.Fatalf("viewport offset = %d, want less than bottom offset %d after wheel up", model.viewport.YOffset, bottomOffset)
	}
}

func TestManualScrollUpPreventsAutoScrollOnNewContent(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 8})
	model := updated.(Model)
	for i := 0; i < 10; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()

	updated, _ = model.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp})
	model = updated.(Model)
	scrolledOffset := model.viewport.YOffset

	updated, _ = model.Update(Event{Name: "assistant_delta", Content: "new streaming content"})
	model = updated.(Model)

	if model.viewport.YOffset != scrolledOffset {
		t.Fatalf("viewport offset = %d, want to stay at manually scrolled offset %d when new content arrives", model.viewport.YOffset, scrolledOffset)
	}
	if !strings.Contains(model.View(), "问 go-agent 一件事") {
		t.Fatalf("view = %q, want fixed prompt to remain visible while history stays scrolled", model.View())
	}
}

func TestBottomPositionAutoFollowsNewContent(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 8})
	model := updated.(Model)
	for i := 0; i < 10; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()

	updated, _ = model.Update(Event{Name: "assistant_delta", Content: "new streaming content"})
	model = updated.(Model)

	if !model.viewport.AtBottom() {
		t.Fatalf("viewport offset = %d, want to keep following new content at bottom", model.viewport.YOffset)
	}
	if !strings.Contains(model.View(), "问 go-agent 一件事") {
		t.Fatalf("view = %q, want active prompt visible while following bottom", model.View())
	}
}

func TestScrollingBackToBottomRestoresAutoFollow(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 60, Height: 8})
	model := updated.(Model)
	for i := 0; i < 10; i++ {
		model.appendMessage("user", "hello")
		model.appendMessage("assistant", "reply")
	}
	model.refreshViewport()

	updated, _ = model.Update(tea.MouseMsg{Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyPgDown})
	model = updated.(Model)
	if !model.viewport.AtBottom() {
		t.Fatalf("viewport offset = %d, want page down to return to bottom", model.viewport.YOffset)
	}

	updated, _ = model.Update(Event{Name: "assistant_delta", Content: "new streaming content"})
	model = updated.(Model)

	if !model.viewport.AtBottom() {
		t.Fatalf("viewport offset = %d, want new content to follow after returning to bottom", model.viewport.YOffset)
	}
}

func TestTypingRefreshesInlinePrompt(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	model := updated.(Model)

	if !strings.Contains(plainTerminalText(model.View()), "› h") {
		t.Fatalf("view = %q, want inline prompt to show typed text", model.View())
	}
}

func TestTypingSpaceRefreshesInlinePrompt(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune(" ")})
	model := updated.(Model)

	if !strings.Contains(plainTerminalText(model.renderInput()), "h "+inputCursor) {
		t.Fatalf("input = %q, want typed space to remain visible before cursor", model.renderInput())
	}
}

func TestSubmittingInputPreservesTypedSpaces(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("hello ")})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyEnter})
	model := updated.(Model)

	if len(model.messages) == 0 || model.messages[len(model.messages)-1].Content != "hello " {
		t.Fatalf("messages = %#v, want submitted text to preserve typed trailing space", model.messages)
	}
}

func TestRunConfigDoesNotCaptureMouseClicks(t *testing.T) {
	config := newRunConfig(context.Background(), NewModel(nil, SessionInfo{}))

	if !config.altScreen || config.mouseCellMotion {
		t.Fatalf("run config = %#v, want alt screen without mouse capture so terminal clicks keep native selection behavior", config)
	}
}

func TestConversationFrameDoesNotRenderSelectableBlankPadding(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	model := updated.(Model)

	frame := model.renderConversationFrame()
	for _, line := range strings.Split(frame, "\n") {
		if line != "" && strings.Trim(line, " ") == "" {
			t.Fatalf("conversation frame contains selectable blank padding line %q", line)
		}
		if strings.HasSuffix(line, strings.Repeat(" ", 4)) {
			t.Fatalf("conversation frame line has selectable right padding %q", line)
		}
	}
}

func TestUserMessagesRenderMutedFromAssistantOutput(t *testing.T) {
	if got := userStyle().GetForeground(); got != tuiPalette.muted {
		t.Fatalf("user foreground = %#v, want muted grey foreground to separate it from assistant output", got)
	}
	if userStyle().GetForeground() == tuiPalette.ink {
		t.Fatalf("user foreground = %#v, want a different color from the default assistant output", userStyle().GetForeground())
	}
}

func TestEmptyInputShowsCursor(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)

	input := model.renderInput()
	if !strings.Contains(input, "█") {
		t.Fatalf("input = %q, want visible cursor for empty input", input)
	}
	if !strings.Contains(input, "问 go-agent 一件事") {
		t.Fatalf("input = %q, want placeholder to remain visible", input)
	}
}

func TestTerminalColorProbeResponseDoesNotPolluteInput(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("11;rgb:1818/1818/1818")})
	model := updated.(Model)

	if strings.Contains(model.renderInput(), "rgb:") {
		t.Fatalf("input = %q, want terminal color probe response ignored", model.renderInput())
	}
}

func TestMessagesWrapToTerminalWidth(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 40, Height: 24})
	model := updated.(Model)

	longChinese := strings.Repeat("这是一段很长的中文回答", 8)
	for _, msg := range []Message{
		{Role: "user", Content: longChinese},
		{Role: "assistant", Content: longChinese},
		{Role: "system", Content: longChinese},
	} {
		if got := maxRenderedLineWidth(model.renderMessage(msg)); got > model.width {
			t.Fatalf("%s message rendered line width = %d, want <= terminal width %d; render = %q", msg.Role, got, model.width, model.renderMessage(msg))
		}
	}
}

func maxRenderedLineWidth(s string) int {
	maxWidth := 0
	for _, line := range strings.Split(s, "\n") {
		maxWidth = max(maxWidth, lipgloss.Width(line))
	}
	return maxWidth
}

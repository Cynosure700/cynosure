package tui

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"cynosure/internal/agent/mcp"
	"cynosure/internal/agent/storage"
	"cynosure/internal/sessions"
)

var ansiEscapePattern = regexp.MustCompile("\x1b\\[[0-9;]*m")
var ansiBackgroundPattern = regexp.MustCompile("\x1b\\[[0-9;]*48;[0-9;]*m")
var ansiBlueForegroundPattern = regexp.MustCompile("\x1b\\[[0-9;]*38;5;39m")
var ansiYellowForegroundPattern = regexp.MustCompile("\x1b\\[[0-9;]*38;5;229m")
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

func TestAssistantFileAndDirectoryReferencesRenderYellow(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.renderer = newMarkdownRenderer(app.messageWidth())
	content := "需要检查 cynosure/internal/tui/app.go 和 cynosure/internal/tui/ 目录。"

	rendered := app.renderMessage(Message{Role: "assistant", Content: content})

	for _, want := range []string{"cynosure/internal/tui/app.go", "cynosure/internal/tui/"} {
		if !strings.Contains(plainTerminalText(rendered), want) {
			t.Fatalf("assistant render = %q, want visible path %q", rendered, want)
		}
	}
	if got := len(ansiYellowForegroundPattern.FindAllString(rendered, -1)); got < 2 {
		t.Fatalf("assistant render = %q, want both file and directory references rendered yellow", rendered)
	}
}

func TestAssistantFunctionReferencesRenderYellow(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.renderer = newMarkdownRenderer(app.messageWidth())
	content := "需要检查 renderMessage()、 Model.Update() 和 colorizeFileReferencesWithRestore() 函数。"

	rendered := app.renderMessage(Message{Role: "assistant", Content: content})

	for _, want := range []string{"renderMessage()", "Model.Update()", "colorizeFileReferencesWithRestore()"} {
		if !strings.Contains(plainTerminalText(rendered), want) {
			t.Fatalf("assistant render = %q, want visible function reference %q", rendered, want)
		}
	}
	if got := len(ansiYellowForegroundPattern.FindAllString(rendered, -1)); got < 3 {
		t.Fatalf("assistant render = %q, want function references rendered yellow", rendered)
	}
}

func TestAssistantInlineCodeReferencesRenderYellow(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 160
	app.renderer = newMarkdownRenderer(app.messageWidth())
	content := "`internal/agent/mcp/config.go` 的 `LoadWorkspaceConfig` 函数在 `mcp.Manager` 中使用。"

	rendered := app.renderMessage(Message{Role: "assistant", Content: content})

	for _, want := range []string{"internal/agent/mcp/config.go", "LoadWorkspaceConfig", "mcp.Manager"} {
		if !strings.Contains(plainTerminalText(rendered), want) {
			t.Fatalf("assistant render = %q, want visible inline code reference %q", rendered, want)
		}
	}
	if got := len(ansiYellowForegroundPattern.FindAllString(rendered, -1)); got < 3 {
		t.Fatalf("assistant render = %q, want inline code references rendered yellow", rendered)
	}
}

func TestToolMessageFileAndDirectoryReferencesKeepToolStyle(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:          "Read",
		ArgsPreview:   "file_path=cynosure/internal/tui/app.go",
		Status:        "success",
		ResultPreview: "listed cynosure/internal/tui/",
	}}

	rendered := app.renderMessage(msg)

	if !strings.Contains(plainTerminalText(rendered), "cynosure/internal/tui/app.go") || !strings.Contains(plainTerminalText(rendered), "cynosure/internal/tui/") {
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
		ArgsPreview:   `path="/Users/bytedance/golang_pro/cynosure", pattern="**/README*"`,
		Status:        "success",
		ResultPreview: "Found 1 file(s).",
	}}

	rendered := app.renderMessage(msg)
	plain := plainTerminalText(rendered)

	if !strings.Contains(plain, "● ✓ glob") {
		t.Fatalf("tool render = %q, want tool call line prefixed with a small bullet", rendered)
	}
	if !strings.Contains(rendered, ansiForeground(tuiPalette.blue)+"●") {
		t.Fatalf("tool render = %q, want leading bullet rendered blue", rendered)
	}
}

func TestToolMessageUsesLowercaseFileSearchDisplayNames(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	for _, tt := range []struct {
		name string
		want string
	}{
		{name: "write_file", want: "● ✓ write("},
		{name: "read_file", want: "● ✓ read("},
		{name: "Glob", want: "● ✓ glob("},
		{name: "grep", want: "● ✓ grep("},
		{name: "edit_file", want: "● ✓ edit("},
	} {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{Role: "tool", ToolCall: &ToolCallView{
				Name:        tt.name,
				ArgsPreview: "file_path: README.md",
				Status:      "success",
			}}

			rendered := plainTerminalText(app.renderMessage(msg))

			if !strings.Contains(rendered, tt.want) {
				t.Fatalf("tool render = %q, want %q", rendered, tt.want)
			}
		})
	}
}

func TestSpawnSubagentToolMessageDisplaysSubTypeAsToolName(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	for _, tt := range []struct {
		name    string
		rawArgs string
		want    string
	}{
		{name: "explore", rawArgs: `{"sub_type":"explore","task":"inspect workspace"}`, want: "● ⏺ Explore("},
		{name: "general", rawArgs: `{"sub_type":"general","task":"analyze design"}`, want: "● ⏺ General("},
	} {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{Role: "tool", ToolCall: &ToolCallView{
				Name:        "spawn_subagent",
				RawArgs:     tt.rawArgs,
				ArgsPreview: tt.rawArgs,
				Status:      "running",
			}}

			rendered := plainTerminalText(app.renderMessage(msg))

			if !strings.Contains(rendered, tt.want) {
				t.Fatalf("tool render = %q, want subtype display %q", rendered, tt.want)
			}
			if strings.Contains(rendered, "SpawnSubagent") {
				t.Fatalf("tool render = %q, should hide raw spawn_subagent tool name", rendered)
			}
		})
	}
}

func TestWorkspaceDirectoryRendersYellowInHeader(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", ModelID: "deepseek-v4-flash"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)

	header := model.renderHeader()

	if !strings.Contains(plainTerminalText(header), "/tmp/project") {
		t.Fatalf("header = %q, want visible workspace directory", header)
	}
	if !ansiYellowForegroundPattern.MatchString(header) {
		t.Fatalf("header = %q, want workspace directory rendered yellow", header)
	}
}

func TestTypedInputKeepsNormalDisplay(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.input.SetValue("检查 cynosure/internal/tui/app.go 后继续")

	rendered := app.renderInput()

	if !strings.Contains(plainTerminalText(rendered), "检查 cynosure/internal/tui/app.go 后继续") {
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

func TestEscInterruptsRunningGenerationSilently(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.running = true
	app.generation = 7
	canceled := false
	app.cancel = func() { canceled = true }
	app.messages = []Message{
		{Role: "user", Content: "写一个长回答"},
		{Role: "assistant", Content: "已经输出一半"},
	}

	updated, cmd := app.Update(tea.KeyMsg{Type: tea.KeyEsc})
	model := updated.(Model)

	if cmd != nil {
		t.Fatalf("expected no command after esc interrupt, got %#v", cmd)
	}
	if !canceled {
		t.Fatal("expected esc to cancel the running response context")
	}
	if model.running {
		t.Fatal("expected running to be cleared after esc")
	}
	if model.generation != 8 {
		t.Fatalf("expected generation to advance and ignore stale events, got %d", model.generation)
	}
	if len(model.messages) != 2 {
		t.Fatalf("expected esc interrupt to avoid adding a system message, got %#v", model.messages)
	}
	if model.messages[1].Content != "已经输出一半" {
		t.Fatalf("expected current partial assistant output to remain visible, got %#v", model.messages[1])
	}
}

func TestSubmittedUserMessageRendersAsGreySelectedLine(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 56

	rendered := app.renderMessage(Message{Role: "user", Content: "检查 cynosure/internal/tui/app.go 后继续"})

	if !strings.Contains(plainTerminalText(rendered), "检查 cynosure/internal/tui/app.go 后继续") {
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

	rendered := app.renderMessage(Message{Role: "error", Content: "打开 cynosure/internal/tui/app.go 失败"})

	if !strings.Contains(plainTerminalText(rendered), "打开 cynosure/internal/tui/app.go 失败") {
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

func (f *fakeSessionResumer) StartNewConversation(ctx context.Context, user storage.User) (storage.Conversation, error) {
	return storage.Conversation{ID: "conv-new", SessionID: "session-new", UserID: user.ID}, nil
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

func TestHandleHelpCommandListsAvailableSlashCommands(t *testing.T) {
	app := NewModel(nil, SessionInfo{})

	if handled := app.handleSlashCommand("/help"); !handled {
		t.Fatal("/help was not handled")
	}
	content := app.messages[len(app.messages)-1].Content
	for _, want := range []string{
		"/clear：开启全新对话并清空当前上下文。",
		"/cwd：显示当前工作区。",
		"/skills：显示已加载的 Skill 列表。",
		"/mcp：显示 MCP server 状态和工具数量。",
		"/resume：查看并恢复当前工作区的历史会话。",
		"/cancel：在恢复历史会话选择中取消当前选择。",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("/help output = %q, want it to contain %q", content, want)
		}
	}
	for _, line := range strings.Split(content, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(strings.TrimSpace(line), "/") {
			t.Fatalf("/help line = %q, want each non-empty line to describe a slash command", line)
		}
	}
}

func TestSlashLikeNaturalQuestionIsNotHandledAsCommand(t *testing.T) {
	app := NewModel(nil, SessionInfo{})

	if handled := app.handleSlashCommand("/help会输出哪些内容"); handled {
		t.Fatalf("natural question starting with /help should not be handled as a slash command: %#v", app.messages)
	}
	if len(app.messages) != 0 {
		t.Fatalf("messages = %#v, want no system command output", app.messages)
	}
}

func TestHandleSkillsCommandShowsSkillDetails(t *testing.T) {
	app := NewModel(nil, SessionInfo{Skills: []sessions.SkillSummary{
		{Name: "project-helper", Source: "workspace", Description: "Project helper", Path: "/project/.cynosure/skills/project-helper/skill.md"},
		{Name: "personal-tool", Source: "user", Description: "Personal tool", Path: "/home/.cynosure/skills/personal-tool/SKILL.md"},
		{Name: "skill-creator", Source: "builtin", Description: "Builtin skill", Path: "skill-creator/SKILL.md"},
	}})

	if handled := app.handleSlashCommand("/skills"); !handled {
		t.Fatal("/skills was not handled")
	}
	content := app.messages[len(app.messages)-1].Content
	for _, want := range []string{
		"已加载 Skills：3 个",
		"project-helper", "当前项目", "Project helper",
		"personal-tool", "家目录", "Personal tool",
		"skill-creator", "系统内置", "Builtin skill",
	} {
		if !strings.Contains(content, want) {
			t.Fatalf("/skills output = %q, want it to contain %q", content, want)
		}
	}
	for _, forbidden := range []string{
		"/project/.cynosure/skills/project-helper/skill.md",
		"/home/.cynosure/skills/personal-tool/SKILL.md",
		"path:",
		"[workspace]", "[user]", "[builtin]",
	} {
		if strings.Contains(content, forbidden) {
			t.Fatalf("/skills output = %q, should not contain %q", content, forbidden)
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
		sessions: []storage.ResumableSession{{SessionID: "session-1", Title: "第一段会话", UpdatedAt: time.Now(), MessageCount: 4}},
		conv:     storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: "local-user", Title: "第一段会话"},
		history: []storage.Message{
			{Role: "user", Content: "hello"},
			{Role: "assistant", Content: "let me check", ToolCalls: []storage.MessageToolCall{
				{ID: "call_1", Type: "function", Function: storage.MessageFunctionCall{Name: "read_file", Arguments: `{"file_path":"README.md"}`}},
			}},
			{Role: "tool", ToolCallID: "call_1", Content: `{"status":"success","result":"file contents"}`},
			{Role: "assistant", Content: "hi"},
		},
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
	rendered := plainTerminalText(app.renderMessages())
	// 原样还原历史会话：user/assistant 文本以及工具调用行都应被展示。
	for _, want := range []string{"hello", "let me check", "README.md", "file contents", "hi"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered messages = %q, want %q", rendered, want)
		}
	}
}

func TestResumeReproducesToolCallsInOrder(t *testing.T) {
	history := []storage.Message{
		{Role: "user", Content: "请运行测试"},
		{Role: "assistant", Content: "我先跑一下测试。", ToolCalls: []storage.MessageToolCall{
			{ID: "call_1", Type: "function", Function: storage.MessageFunctionCall{Name: "bash", Arguments: `{"command":"go test ./..."}`}},
		}},
		{Role: "tool", ToolCallID: "call_1", Content: `{"status":"success","result":"ok cynosure/internal/tui 0.42s"}`},
		{Role: "assistant", Content: "测试通过了。"},
	}
	resumer := &fakeSessionResumer{
		sessions: []storage.ResumableSession{{SessionID: "session-1", Title: "会话", UpdatedAt: time.Now(), MessageCount: len(history)}},
		conv:     storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: "local-user", Title: "会话"},
		history:  history,
	}
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project", User: storage.User{ID: "local-user"}, Resumer: resumer})
	app.handleSlashCommand("/resume")
	if handled := app.handleResumeSelection("1"); !handled {
		t.Fatal("resume selection was not handled")
	}

	rendered := plainTerminalText(app.renderMessages())
	for _, want := range []string{"请运行测试", "我先跑一下测试。", "go test ./...", "ok cynosure/internal/tui 0.42s", "测试通过了。"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want it to contain %q", rendered, want)
		}
	}

	// 工具调用行应被还原，并携带成功状态图标与结果预览。
	var toolMsg *Message
	for i := range app.messages {
		if app.messages[i].Role == "tool" && app.messages[i].ToolCall != nil {
			toolMsg = &app.messages[i]
			break
		}
	}
	if toolMsg == nil {
		t.Fatalf("messages = %#v, want a reconstructed tool message", app.messages)
	}
	if toolMsg.ToolCall.Name != "bash" {
		t.Fatalf("tool name = %q, want bash", toolMsg.ToolCall.Name)
	}
	if toolMsg.ToolCall.Status != "success" {
		t.Fatalf("tool status = %q, want success", toolMsg.ToolCall.Status)
	}
	if !strings.Contains(rendered, "✓") {
		t.Fatalf("rendered = %q, want success icon", rendered)
	}

	// 顺序应与历史一致：assistant 文本在工具行之前，最终回答在工具行之后。
	assistantIdx := strings.Index(rendered, "我先跑一下测试。")
	toolIdx := strings.Index(rendered, "go test ./...")
	finalIdx := strings.Index(rendered, "测试通过了。")
	if !(assistantIdx < toolIdx && toolIdx < finalIdx) {
		t.Fatalf("ordering wrong: assistant=%d tool=%d final=%d", assistantIdx, toolIdx, finalIdx)
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

	updated, cmd := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "stale"})
	model := updated.(Model)

	// 合并事件后，本轮 Update 会非阻塞 drain 掉 channel 中已就绪的同代事件，
	// 因此 fresh 事件在同一次 Update 内被应用，而 stale 事件被忽略。
	if len(model.messages) != 1 || model.messages[0].Content != "fresh" {
		t.Fatalf("messages = %#v, want stale ignored and fresh applied", model.messages)
	}
	// 仍在生成中，应继续等待后续事件。
	if cmd == nil {
		t.Fatal("expected model to keep waiting while current generation is still running")
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

func TestUpdateCoalescesPendingStreamingEventsIntoOneFrame(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true
	// 预填若干已就绪的流式增量，模拟底部正在快速输出时 channel 中的堆积。
	app.events <- Event{Generation: 1, Name: "assistant_delta", Content: "二"}
	app.events <- Event{Generation: 1, Name: "assistant_delta", Content: "三"}
	app.events <- Event{Generation: 1, Name: "assistant_delta", Content: "四"}

	updated, cmd := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "一"})
	model := updated.(Model)

	// 四个增量应在一次 Update 内被合并应用，channel 被清空。
	if len(model.events) != 0 {
		t.Fatalf("pending events = %d, want all coalesced and drained within one Update", len(model.events))
	}
	if len(model.messages) != 1 || model.messages[0].Content != "一二三四" {
		t.Fatalf("messages = %#v, want all streaming deltas merged into one assistant message", model.messages)
	}
	if cmd == nil {
		t.Fatal("expected model to keep waiting for further events while running")
	}
}

func TestUpdateStopsDrainingAtTerminalDoneEvent(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true
	app.events <- Event{Generation: 1, Name: "done"}

	updated, cmd := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "答案"})
	model := updated.(Model)

	if model.running {
		t.Fatal("drained done event should stop the running state")
	}
	if cmd != nil {
		t.Fatalf("expected no further waiting command after done, got %#v", cmd)
	}
}

func TestRenderCacheReusesUnchangedMessageRenders(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 80
	app.renderer = newMarkdownRenderer(app.messageWidth())
	app.messages = []Message{{Role: "assistant", Content: "你好"}}

	app.renderCachedMessage(app.messages[0])
	key := app.messageRenderKey(app.messages[0], app.messageWidth())
	if _, ok := app.renderCache.entries[key]; !ok {
		t.Fatal("expected first render to populate the cache")
	}
	// 篡改缓存值，若第二次渲染命中缓存则应返回被篡改的内容，证明未重算。
	app.renderCache.entries[key] = "CACHED-SENTINEL"
	if got := app.renderCachedMessage(app.messages[0]); got != "CACHED-SENTINEL" {
		t.Fatalf("second render = %q, want cached value reused instead of recomputed", got)
	}
}

func TestRenderCacheInvalidatesWhenMessageContentChanges(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 80
	app.renderer = newMarkdownRenderer(app.messageWidth())

	app.renderCachedMessage(Message{Role: "assistant", Content: "你好"})
	changed := app.renderCachedMessage(Message{Role: "assistant", Content: "你好世界"})

	if !strings.Contains(plainTerminalText(changed), "你好世界") {
		t.Fatalf("render = %q, want changed content recomputed rather than serving stale cache", changed)
	}
}

func TestRefreshViewportPrunesStaleRenderCacheEntries(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)
	model.messages = []Message{{Role: "assistant", Content: "第一版"}}
	model.refreshViewport()

	// 模拟流式增量：内容增长产生新的缓存 key，旧 key 应在下一次刷新时被清理。
	model.messages[0].Content = "第一版内容增长"
	model.refreshViewport()

	if len(model.renderCache.entries) != len(model.messages) {
		t.Fatalf("cache entries = %d, want pruned down to the %d live messages", len(model.renderCache.entries), len(model.messages))
	}
	key := model.messageRenderKey(model.messages[0], model.messageWidth())
	if _, ok := model.renderCache.entries[key]; !ok {
		t.Fatal("expected the current message render to remain cached after pruning")
	}
}

func TestModelDoesNotStreamReasoningContent(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "reasoning_delta", Content: "先判断是否需要工具"})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())

	if strings.Contains(rendered, "先判断是否需要工具") {
		t.Fatalf("rendered messages = %q, reasoning_content must not be displayed", rendered)
	}
	if strings.Contains(rendered, "思考中") {
		t.Fatalf("rendered messages = %q, should not render reasoning thinking block", rendered)
	}
}

func TestModelAssistantFinalEventKeepsContentAndMeta(t *testing.T) {
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
		t.Fatalf("rendered messages = %q, reasoning_content must never be displayed", rendered)
	}
	view := plainTerminalText(model.View())
	if model.toolCallCount != 2 {
		t.Fatalf("toolCallCount = %d, want metadata to keep updating internally", model.toolCallCount)
	}
	if strings.Contains(view, "工具 2") {
		t.Fatalf("view = %q, should not display tool call count", view)
	}
	if !strings.Contains(view, "上下文 45%") {
		t.Fatalf("view = %q, want context status to remain visible", view)
	}
}

func TestModelMetaEventUpdatesContextStatusWithoutDisplayingToolCount(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "meta", Data: map[string]any{"tool_call_count": 3, "context_tokens": 72000, "context_budget": 100000}})
	model := updated.(Model)

	view := model.View()
	if model.toolCallCount != 3 {
		t.Fatalf("toolCallCount = %d, want metadata to keep updating internally", model.toolCallCount)
	}
	if strings.Contains(view, "工具 3") {
		t.Fatalf("view = %q, should not display tool call count", view)
	}
	if !strings.Contains(view, "上下文 72%") {
		t.Fatalf("view = %q, want context status to remain visible", view)
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
		"result_preview": "ok cynosure/internal/tui 0.42s",
	}})
	model = updated.(Model)
	if len(model.messages) != 1 {
		t.Fatalf("messages = %#v, want tool done to update existing message", model.messages)
	}
	rendered = plainTerminalText(model.renderMessages())
	for _, want := range []string{"✓ Bash", "⎿ success · ok cynosure/internal/tui 0.42s"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
}

func TestModelDisplaysSubagentToolStatusWithoutResultPreview(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "grep",
		"args_preview":       "pattern: TODO",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "grep",
		"args_preview":       "pattern: TODO",
		"status":             "success",
		"result_preview":     "secret internal result",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)

	rendered := plainTerminalText(model.renderMessages())
	for _, want := range []string{"grep", "pattern: TODO"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
	for _, unwanted := range []string{"✓ grep", "⎿ success", "⎿ running"} {
		if strings.Contains(rendered, unwanted) {
			t.Fatalf("rendered = %q, subagent tool status should not include %q", rendered, unwanted)
		}
	}
	for i, line := range strings.Split(strings.TrimRight(rendered, "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if !strings.HasPrefix(line, "  ") {
			t.Fatalf("line %d = %q, subagent tool block should align with status column", i, line)
		}
	}
	if strings.Contains(rendered, "●") {
		t.Fatalf("rendered = %q, subagent tool status should not show the leading blue bullet", rendered)
	}
	if strings.Contains(rendered, "secret internal result") {
		t.Fatalf("rendered = %q, should hide subagent internal result", rendered)
	}
	rawRendered := model.renderMessages()
	if strings.Contains(rawRendered, ansiForeground(tuiPalette.mint)) {
		t.Fatalf("rendered = %q, subagent tool status should not use success green", rawRendered)
	}
	if !strings.Contains(rawRendered, ansiForeground(tuiPalette.muted)) {
		t.Fatalf("rendered = %q, subagent tool status should use muted gray", rawRendered)
	}
}

func TestSubagentToolStatusAttachesUnderRunningColumn(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "parent_spawn",
		"tool_name":    "spawn_subagent",
		"args_preview": "task: inspect",
		"status":       "running",
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "read_file",
		"args_preview":       "file_path: /Users/bytedance/golang_pro/cynosure/cynosure/main.go",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)

	lines := strings.Split(strings.TrimRight(plainTerminalText(model.renderMessages()), "\n"), "\n")
	if len(lines) < 3 {
		t.Fatalf("lines = %#v, want parent tool plus attached subagent tool", lines)
	}
	last := lines[len(lines)-1]
	if strings.TrimSpace(lines[len(lines)-2]) == "" {
		t.Fatalf("lines = %#v, subagent tool should attach directly below running block", lines)
	}
	statusLine := lines[len(lines)-2]
	runningIndex := strings.Index(statusLine, "running")
	readIndex := strings.Index(last, "read")
	runningColumn := -1
	readColumn := -1
	if runningIndex >= 0 {
		runningColumn = lipgloss.Width(statusLine[:runningIndex])
	}
	if readIndex >= 0 {
		readColumn = lipgloss.Width(last[:readIndex])
	}
	if runningColumn < 0 || readColumn < 0 || readColumn != runningColumn {
		t.Fatalf("lines = %#v, want subagent tool aligned with running text column", lines)
	}
	for _, unwanted := range []string{"✓ read", "⏺ read", "⎿ success", "⎿ running"} {
		if strings.Contains(last, unwanted) {
			t.Fatalf("last line = %q, should not contain %q", last, unwanted)
		}
	}
}

func TestConcurrentSubagentToolStatusAttachesToMatchingParent(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "spawn_1",
		"tool_name":    "spawn_subagent",
		"args_preview": "task: inspect runtime",
		"status":       "running",
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "spawn_2",
		"tool_name":    "spawn_subagent",
		"args_preview": "task: inspect tui",
		"status":       "running",
	}})
	model = updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":        "subagent_1:tool_1",
		"tool_name":           "read_file",
		"args_preview":        "file_path: internal/agent/runtime/context_compression.go",
		"status":              "running",
		"scope":               "subagent",
		"ephemeral_group_id":  "subagent_1",
		"parent_tool_call_id": "spawn_1",
		"suppress_result":     true,
	}})
	model = updated.(Model)

	lines := strings.Split(strings.TrimRight(plainTerminalText(model.renderMessages()), "\n"), "\n")
	firstSpawn := indexLineContaining(lines, "task: inspect runtime")
	secondSpawn := indexLineContaining(lines, "task: inspect tui")
	child := indexLineContaining(lines, "read(file_path: internal/agent/runtime/context_compression.go)")
	if firstSpawn < 0 || secondSpawn < 0 || child < 0 {
		t.Fatalf("lines = %#v, want both parent spawns and child tool row", lines)
	}
	if !(firstSpawn < child && child < secondSpawn) {
		t.Fatalf("lines = %#v, child tool row should render under matching first subagent parent", lines)
	}
}

func TestSubagentToolStatusKeepsBlankLineBeforeNextTool(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "spawn_1",
		"tool_name":    "spawn_subagent",
		"args_preview": "task: inspect local",
		"status":       "running",
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":        "subagent_1:tool_1",
		"tool_name":           "read_file",
		"args_preview":        "file_path: internal/local/store_test.go",
		"status":              "running",
		"scope":               "subagent",
		"ephemeral_group_id":  "subagent_1",
		"parent_tool_call_id": "spawn_1",
		"suppress_result":     true,
	}})
	model = updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "spawn_2",
		"tool_name":    "spawn_subagent",
		"args_preview": "task: inspect logger",
		"status":       "running",
	}})
	model = updated.(Model)

	lines := strings.Split(strings.TrimRight(plainTerminalText(model.renderMessages()), "\n"), "\n")
	child := indexLineContaining(lines, "read(file_path: internal/local/store_test.go)")
	next := indexLineContaining(lines, "task: inspect logger")
	if child < 0 || next < 0 {
		t.Fatalf("lines = %#v, want child tool row and following parent tool", lines)
	}
	if next-child < 2 || strings.TrimSpace(lines[child+1]) != "" {
		t.Fatalf("lines = %#v, want blank line between subagent tool row and following tool", lines)
	}
}

func indexLineContaining(lines []string, needle string) int {
	for i, line := range lines {
		if strings.Contains(line, needle) {
			return i
		}
	}
	return -1
}

func TestModelShowsOnlyLatestSubagentToolStatus(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "grep",
		"args_preview":       "pattern: TODO",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_2",
		"tool_name":          "read_file",
		"args_preview":       "file_path: README.md",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)

	if len(model.messages) != 1 {
		t.Fatalf("messages = %#v, want only the latest child tool message", model.messages)
	}
	rendered := plainTerminalText(model.renderMessages())
	if strings.Contains(rendered, "grep") || strings.Contains(rendered, "pattern: TODO") {
		t.Fatalf("rendered = %q, should replace the previous subagent tool status", rendered)
	}
	for _, want := range []string{"read", "file_path: README.md"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want latest subagent tool detail %q", rendered, want)
		}
	}
}

func TestStaleSubagentToolDoneDoesNotReappendReplacedToolStatus(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "glob",
		"args_preview":       "pattern: go.mod",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_2",
		"tool_name":          "ls",
		"args_preview":       "path: /Users/bytedance/golang_pro/cynosure/cynosure",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "glob",
		"args_preview":       "pattern: go.mod",
		"status":             "success",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)

	rendered := plainTerminalText(model.renderMessages())
	if strings.Contains(rendered, "glob") || strings.Contains(rendered, "pattern: go.mod") {
		t.Fatalf("rendered = %q, stale child tool done should not reappend replaced status", rendered)
	}
	if !strings.Contains(rendered, "Ls") || !strings.Contains(rendered, "path: /Users/bytedance/golang_pro/cynosure/cynosure") {
		t.Fatalf("rendered = %q, want latest child tool status to remain visible", rendered)
	}
}

func TestToolMessageRendersAllRawArgs(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 160
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:        "grep",
		RawArgs:     `{"pattern":"TODO","path":"internal/tui","output_mode":"content"}`,
		ArgsPreview: "pattern: TODO",
		Status:      "running",
	}}

	rendered := plainTerminalText(app.renderMessage(msg))

	for _, want := range []string{"output_mode: content", "path: internal/tui", "pattern: TODO"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("tool render = %q, want complete raw arg %q", rendered, want)
		}
	}
}

func TestToolMessageTruncatesLongArgsToTwoLines(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 60
	longTask := "探索这个 Go 项目（工作区根目录 /Users/bytedance/golang_pro/cynosure/cynosure）的 internal/tools 包。这是 agent 的工具系统。请阅读该目录下所有文件并总结：1. 工具是如何定义和注册的 2. 有哪些内置工具 3. 工具执行的流程 4. 是否支持 MCP 工具 5. 工具结果如何反馈给 LLM。请用中文返回结构化总结。"
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:        "spawn_subagent",
		RawArgs:     `{"sub_type":"explore","task":"` + longTask + `"}`,
		ArgsPreview: "task: " + longTask,
		Status:      "running",
	}}

	rendered := plainTerminalText(app.renderMessage(msg))
	lines := strings.Split(strings.TrimRight(rendered, "\n"), "\n")
	// ⎿ 之前的行都属于工具参数行，最多保留两行。
	toolLines := 0
	for _, line := range lines {
		if strings.Contains(line, "⎿") {
			break
		}
		toolLines++
	}
	if toolLines > 2 {
		t.Fatalf("tool arg lines = %d (%q), want at most 2 wrapped lines", toolLines, lines)
	}
	if !strings.Contains(rendered, "…") {
		t.Fatalf("rendered = %q, want ellipsis marking truncated args", rendered)
	}
}

func TestSubagentToolStatusWrapsContinuationAtStatusColumn(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 36
	app.generation = 1
	app.running = true
	longArgs := "task: 探索 /Users/bytedance/golang_pro/cynosure/cynosure 当前项目结构和技术栈"

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "grep",
		"args_preview":       longArgs,
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model := updated.(Model)

	lines := strings.Split(strings.TrimRight(plainTerminalText(model.renderMessages()), "\n"), "\n")
	if len(lines) < 2 {
		t.Fatalf("rendered lines = %#v, want wrapped compact tool line", lines)
	}
	for i, line := range lines {
		if !strings.HasPrefix(line, "  ") {
			t.Fatalf("wrapped line %d = %q, want continuation aligned with status column", i, line)
		}
	}
	rendered := strings.Join(lines, "\n")
	for _, unwanted := range []string{"⎿ running", "⏺ grep", "✓ grep"} {
		if strings.Contains(rendered, unwanted) {
			t.Fatalf("rendered = %q, compact subagent tool line should not contain %q", rendered, unwanted)
		}
	}
}

func TestModelClearsSubagentToolGroupWithoutBlankLines(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "parent_spawn",
		"tool_name":    "spawn_subagent",
		"args_preview": "task: inspect",
		"status":       "running",
	}})
	model := updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_1",
		"tool_name":          "grep",
		"args_preview":       "pattern: TODO",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)
	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id":       "subagent_1:tool_2",
		"tool_name":          "read_file",
		"args_preview":       "file_path: README.md",
		"status":             "running",
		"scope":              "subagent",
		"ephemeral_group_id": "subagent_1",
		"suppress_result":    true,
	}})
	model = updated.(Model)
	if len(model.messages) != 2 {
		t.Fatalf("messages = %#v, want parent plus latest child tool", model.messages)
	}
	renderedBeforeClear := plainTerminalText(model.renderMessages())
	if strings.Contains(renderedBeforeClear, "grep") || strings.Contains(renderedBeforeClear, "pattern: TODO") {
		t.Fatalf("rendered = %q, should replace older subagent tool status before clear", renderedBeforeClear)
	}

	updated, _ = model.Update(Event{Generation: 1, Name: "tool_call_group_clear", Data: map[string]any{
		"ephemeral_group_id": "subagent_1",
	}})
	model = updated.(Model)
	if len(model.messages) != 1 {
		t.Fatalf("messages = %#v, want only parent spawn_subagent after clear", model.messages)
	}
	if model.messages[0].ToolCall == nil || model.messages[0].ToolCall.Name != "spawn_subagent" {
		t.Fatalf("messages = %#v, want parent spawn_subagent preserved", model.messages)
	}
	rendered := plainTerminalText(model.renderMessages())
	if strings.Contains(rendered, "grep") || strings.Contains(rendered, "file") {
		t.Fatalf("rendered = %q, should remove cleared subagent tools", rendered)
	}
	if strings.Contains(rendered, "\n\n\n\n") {
		t.Fatalf("rendered = %q, should not contain blank cleared tool gaps", rendered)
	}
}

func TestModelDisplaysMultilineToolResultPreview(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":   "tool_multiline",
		"tool_name":      "bash",
		"args_preview":   "command: seq 1 7",
		"status":         "success",
		"result_preview": "line1\nline2\nline3\nline4\nline5\n... + 2 lines",
	}})
	model := updated.(Model)

	rendered := plainTerminalText(model.renderMessages())
	for _, want := range []string{"line1", "line2", "line3", "line4", "line5", "... + 2 lines"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want multiline preview segment %q", rendered, want)
		}
	}
	if strings.Contains(rendered, "line6") || strings.Contains(rendered, "line7") {
		t.Fatalf("rendered = %q, should not include omitted lines", rendered)
	}
}

func TestReadFileToolMessageAlignsMultilineContent(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.generation = 1
	app.running = true

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":   "tool_read",
		"tool_name":      "read_file",
		"args_preview":   "file_path: cynosure/internal/tui/app.go",
		"status":         "success",
		"result_preview": "package tui\n\nimport \"strings\"\n\nfunc render() {}",
	}})
	model := updated.(Model)

	lines := strings.Split(plainTerminalText(model.renderMessages()), "\n")
	wantColumn := -1
	for _, line := range lines {
		if idx := strings.Index(line, "package tui"); idx >= 0 {
			wantColumn = lipgloss.Width(line[:idx])
			break
		}
	}
	if wantColumn < 0 {
		t.Fatalf("rendered lines = %#v, want first file content line", lines)
	}
	for _, want := range []string{"import \"strings\"", "func render() {}"} {
		found := false
		for _, line := range lines {
			if strings.Contains(line, want) {
				found = true
				idx := strings.Index(line, want)
				if got := lipgloss.Width(line[:idx]); got != wantColumn {
					t.Fatalf("line %q starts content at column %d, want %d", line, got, wantColumn)
				}
			}
		}
		if !found {
			t.Fatalf("rendered lines = %#v, want line containing %q", lines, want)
		}
	}
}

func TestWriteFileToolMessageRendersFileContentPreview(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	rawArgs := `{"file_path":"src/foo.ts","content":"import React from 'react'\nimport { useState } from 'react'\n\nconst a = 1\nconst b = 2\nconst c = 3\nconst d = 4\nconst e = 5\nconst f = 6\nconst g = 7\nconst h = 8\nconst i = 9\nconst j = 10\nconst k = 11\nconst l = 12"}`
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:    "write_file",
		RawArgs: rawArgs,
		Status:  "success",
	}}

	rendered := plainTerminalText(app.renderMessage(msg))

	for _, want := range []string{
		"● ✓ Wrote 15 lines to src/foo.ts",
		"  1│ import React from 'react'",
		"  2│ import { useState } from 'react'",
		"  3│",
		" 10│ const g = 7",
		"... +5 lines  [Ctrl+O to expand]",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
	if strings.Contains(rendered, "const h = 8") {
		t.Fatalf("rendered = %q, should truncate write content after 10 lines", rendered)
	}
	if strings.Contains(rendered, "✓ write(") || strings.Contains(rendered, "success ·") {
		t.Fatalf("rendered = %q, should use specialized write display instead of generic tool display", rendered)
	}
}

func TestCtrlOExpandsAllWriteAndEditFileToolMessages(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.height = 40
	app.generation = 1
	writeArgs1 := `{"file_path":"src/foo.ts","content":"import React from 'react'\nimport { useState } from 'react'\n\nconst a = 1\nconst b = 2\nconst c = 3\nconst d = 4\nconst e = 5\nconst f = 6\nconst g = 7\nconst h = 8\nconst i = 9\nconst j = 10\nconst k = 11\nconst l = 12"}`
	writeArgs2 := `{"file_path":"src/bar.ts","content":"line 1\nline 2\nline 3\nline 4\nline 5\nline 6\nline 7\nline 8\nline 9\nline 10\nline 11"}`
	editArgs := `{"file_path":"src/edit.ts","old_text":"old 1\nold 2\nold 3\nold 4\nold 5\nold 6\nold 7\nold 8\nold 9\nold 10\nold 11","new_text":"new 1\nnew 2\nnew 3\nnew 4\nnew 5\nnew 6\nnew 7\nnew 8\nnew 9\nnew 10\nnew 11"}`
	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_write_1",
		"tool_name":    "write_file",
		"raw_args":     writeArgs1,
		"status":       "success",
	}})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_write_2",
		"tool_name":    "write_file",
		"raw_args":     writeArgs2,
		"status":       "success",
	}})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_edit",
		"tool_name":    "edit_file",
		"raw_args":     editArgs,
		"status":       "success",
	}})
	model := updated.(Model)

	collapsed := plainTerminalText(model.renderMessages())
	for _, forbidden := range []string{"const h = 8", "line 11", "old 11", "new 11"} {
		if strings.Contains(collapsed, forbidden) {
			t.Fatalf("collapsed render = %q, should hide %q before expansion", collapsed, forbidden)
		}
	}
	if got := strings.Count(collapsed, "[Ctrl+O to expand]"); got != 3 {
		t.Fatalf("collapsed render = %q, want 3 expand hints, got %d", collapsed, got)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
	model = updated.(Model)
	expanded := plainTerminalText(model.renderMessages())
	for _, want := range []string{
		" 11│ const h = 8",
		" 15│ const l = 12",
		" 11│ line 11",
		`-11│ old 11`,
		`+11│ new 11`,
		"[Ctrl+O to collapse]",
	} {
		if !strings.Contains(expanded, want) {
			t.Fatalf("expanded render = %q, want %q", expanded, want)
		}
	}
	if strings.Contains(expanded, "[Ctrl+O to expand]") {
		t.Fatalf("expanded render = %q, should hide expand truncation hint", expanded)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
	model = updated.(Model)
	collapsedAgain := plainTerminalText(model.renderMessages())
	for _, forbidden := range []string{"const h = 8", "line 11", "old 11", "new 11"} {
		if strings.Contains(collapsedAgain, forbidden) {
			t.Fatalf("collapsed-again render = %q, should hide %q after collapse", collapsedAgain, forbidden)
		}
	}
	if got := strings.Count(collapsedAgain, "[Ctrl+O to expand]"); got != 3 {
		t.Fatalf("collapsed-again render = %q, want 3 expand hints, got %d", collapsedAgain, got)
	}
}

func TestCtrlOExpandsMultiEditToolMessages(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	app.height = 40
	app.generation = 1
	multiEditArgs := `{"files":[{"file_path":"src/multi.ts","edits":[{"old_string":"old 1\nold 2\nold 3\nold 4\nold 5\nold 6\nold 7\nold 8\nold 9\nold 10\nold 11","new_string":"new 1\nnew 2\nnew 3\nnew 4\nnew 5\nnew 6\nnew 7\nnew 8\nnew 9\nnew 10\nnew 11"}]}]}`
	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_multi_edit",
		"tool_name":    "multi_edit",
		"raw_args":     multiEditArgs,
		"status":       "success",
	}})
	model := updated.(Model)

	collapsed := plainTerminalText(model.renderMessages())
	if strings.Contains(collapsed, "old 11") || strings.Contains(collapsed, "new 11") {
		t.Fatalf("collapsed render = %q, should hide long multi_edit diff before expansion", collapsed)
	}
	if !strings.Contains(collapsed, "[Ctrl+O to expand]") {
		t.Fatalf("collapsed render = %q, want expand hint", collapsed)
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyCtrlO})
	model = updated.(Model)
	expanded := plainTerminalText(model.renderMessages())
	for _, want := range []string{
		`-11│ old 11`,
		`+11│ new 11`,
		"[Ctrl+O to collapse]",
	} {
		if !strings.Contains(expanded, want) {
			t.Fatalf("expanded render = %q, want %q", expanded, want)
		}
	}
	if strings.Contains(expanded, "[Ctrl+O to expand]") {
		t.Fatalf("expanded render = %q, should hide expand hint after Ctrl+O", expanded)
	}
}

func TestEditFileToolMessageRendersDiffPreview(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	rawArgs := `{"file_path":"src/foo.ts","old_text":"const old = \"hello\"\nreturn false","new_text":"const new = \"world\"\nreturn true\nconsole.log(\"done\")"}`
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:    "edit_file",
		RawArgs: rawArgs,
		Status:  "success",
	}}

	renderedWithANSI := app.renderMessage(msg)
	rendered := plainTerminalText(renderedWithANSI)

	for _, want := range []string{
		"● ✓ Added 3 lines, removed 2 lines",
		`- 1│ const old = "hello"`,
		`+ 1│ const new = "world"`,
		"...",
		"- 2│ return false",
		"+ 2│ return true",
		`+ 3│ console.log("done")`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
	if strings.Contains(rendered, "✓ edit(") || strings.Contains(rendered, "success ·") {
		t.Fatalf("rendered = %q, should use specialized edit display instead of generic tool display", rendered)
	}
	if !strings.Contains(renderedWithANSI, string(tuiPalette.coral)) || !strings.Contains(renderedWithANSI, string(tuiPalette.mint)) {
		t.Fatalf("rendered = %q, want removed and added lines colored", renderedWithANSI)
	}
}

func TestMultiEditToolMessageRendersDiffPreviewByFile(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	rawArgs := `{"files":[{"file_path":"src/foo.ts","edits":[{"old_string":"const old = \"hello\"","new_string":"const new = \"world\""},{"old_string":"return false","new_string":"return true"}]},{"file_path":"src/bar.ts","edits":[{"old_string":"name = \"old\"","new_string":"name = \"new\""}]}]}`
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:    "multi_edit",
		RawArgs: rawArgs,
		Status:  "success",
	}}

	renderedWithANSI := app.renderMessage(msg)
	rendered := plainTerminalText(renderedWithANSI)

	for _, want := range []string{
		"● ✓ src/foo.ts: Added 2 lines, removed 2 lines",
		`- 1│ const old = "hello"`,
		`+ 1│ const new = "world"`,
		"- 1│ return false",
		"+ 1│ return true",
		"● ✓ src/bar.ts: Added 1 line, removed 1 line",
		`- 1│ name = "old"`,
		`+ 1│ name = "new"`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
	if !strings.Contains(rendered, "\n\n● ✓ src/bar.ts") {
		t.Fatalf("rendered = %q, want a blank line between multi_edit file blocks", rendered)
	}
	if got := len(ansiBlueForegroundPattern.FindAllString(renderedWithANSI, -1)); got != 2 {
		t.Fatalf("rendered = %q, want one blue bullet per edited file, got %d", renderedWithANSI, got)
	}
	if strings.Contains(rendered, "✓ multi_edit(") || strings.Contains(rendered, "success ·") {
		t.Fatalf("rendered = %q, should use specialized multi_edit display instead of generic tool display", rendered)
	}
	if !strings.Contains(renderedWithANSI, string(tuiPalette.coral)) || !strings.Contains(renderedWithANSI, string(tuiPalette.mint)) {
		t.Fatalf("rendered = %q, want removed and added lines colored", renderedWithANSI)
	}
}

func TestWriteAndEditFileToolMessagesKeepGenericDisplayBeforeSuccess(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	for _, tt := range []struct {
		name   string
		status string
		want   string
	}{
		{name: "write_file", status: "running", want: "● ⏺ write("},
		{name: "edit_file", status: "error", want: "● ✗ edit("},
	} {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{Role: "tool", ToolCall: &ToolCallView{
				Name:          tt.name,
				RawArgs:       `{"file_path":"src/foo.ts","content":"x","old_text":"x","new_text":"y"}`,
				Status:        tt.status,
				ResultPreview: "Error: failed",
			}}

			rendered := plainTerminalText(app.renderMessage(msg))

			if !strings.Contains(rendered, tt.want) {
				t.Fatalf("rendered = %q, want generic display %q", rendered, tt.want)
			}
			if strings.Contains(rendered, "Wrote 1 line") || strings.Contains(rendered, "Added 1 line") {
				t.Fatalf("rendered = %q, should not use success preview for status %q", rendered, tt.status)
			}
		})
	}
}

func TestTodoWriteToolMessageRendersCheckboxList(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true
	rawArgs := `{"todos":[{"id":"1","content":"梳理需求","status":"completed"},{"id":"2","content":"实现功能","status":"in_progress"},{"id":"3","content":"运行测试","status":"pending"}]}`

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":   "tool_todo",
		"tool_name":      "todoWrite",
		"raw_args":       rawArgs,
		"args_preview":   "todos: 3 items",
		"status":         "success",
		"result_preview": "Todo list updated: 3 items",
	}})
	model := updated.(Model)

	rawRendered := model.renderMessages()
	rendered := plainTerminalText(rawRendered)
	for _, want := range []string{"✓ Update Todos", "⎿ [✓] 梳理需求", "  [•] 实现功能", "  [ ] 运行测试"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
	if !strings.Contains(rawRendered, ansiForeground(tuiPalette.blue)+"•") {
		t.Fatalf("rendered = %q, want in-progress todo dot rendered blue", rawRendered)
	}
	for _, forbidden := range []string{"todos: 3 items", "Todo list updated"} {
		if strings.Contains(rendered, forbidden) {
			t.Fatalf("rendered = %q, should hide generic todoWrite preview %q", rendered, forbidden)
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
	for _, want := range []string{"✗ read", "⎿ rejected · Error: outside workspace"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q", rendered, want)
		}
	}
}

func TestModelKeepsToolResultsAndContentVisibleAfterAssistantReply(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.appendMessage("user", "新问题")
	app.generation = 1
	app.running = true

	// 第一轮的 content 实时输出（工具轮也输出 content）。
	updated, _ := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "我先跑个测试"})
	updated, _ = updated.(Model).Update(Event{Generation: 1, Name: "reasoning_delta", Content: "需要先跑测试"})
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
		"result_preview": "ok cynosure/internal/tui 0.42s",
	}})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())
	for _, want := range []string{"我先跑个测试", "✓ Bash", "ok cynosure/internal/tui 0.42s"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q visible", rendered, want)
		}
	}
	if strings.Contains(rendered, "需要先跑测试") {
		t.Fatalf("rendered = %q, reasoning_content must not be displayed", rendered)
	}

	// 第二轮最终答案开始后，前面的内容（content 与工具调用）依然保留可见。
	updated, _ = model.Update(Event{Generation: 1, Name: "assistant_delta", Content: "完成"})
	model = updated.(Model)
	rendered = plainTerminalText(model.renderMessages())
	for _, want := range []string{"新问题", "我先跑个测试", "✓ Bash", "ok cynosure/internal/tui 0.42s", "完成"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want %q to stay visible after assistant reply starts", rendered, want)
		}
	}
	if strings.Contains(rendered, "需要先跑测试") {
		t.Fatalf("rendered = %q, reasoning_content must not be displayed", rendered)
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

func TestThinkingIndicatorSwitchesToWorkingWhenAssistantContentStarts(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.messages = []Message{{Role: "user", Content: "hello"}}
	app.generation = 1
	app.running = true
	app.thinkingStartedAt = time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	app.thinkingNow = app.thinkingStartedAt.Add(3 * time.Second)

	updated, _ := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "我先查一下"})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())

	if !strings.Contains(rendered, "我先查一下") {
		t.Fatalf("rendered messages = %q, want assistant reply content", rendered)
	}
	if !strings.Contains(rendered, "* Working... (3s)") {
		t.Fatalf("rendered messages = %q, should switch to Working once assistant content starts", rendered)
	}
	if strings.Contains(rendered, "Thinking...") {
		t.Fatalf("rendered messages = %q, should stop showing Thinking once assistant content starts", rendered)
	}
}

func TestThinkingIndicatorSwitchesToWorkingWhenToolCallStarts(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.messages = []Message{{Role: "user", Content: "hello"}}
	app.generation = 1
	app.running = true
	app.thinkingStartedAt = time.Date(2026, 6, 15, 10, 0, 0, 0, time.UTC)
	app.thinkingNow = app.thinkingStartedAt.Add(4 * time.Second)

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_start", Data: map[string]any{
		"tool_call_id": "tool_1",
		"tool_name":    "bash",
		"args_preview": "command: pwd",
		"status":       "running",
	}})
	model := updated.(Model)
	rendered := plainTerminalText(model.renderMessages())

	for _, want := range []string{"Bash", "command: pwd", "* Working... (4s)"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered messages = %q, want %q visible", rendered, want)
		}
	}
	if strings.Contains(rendered, "Thinking...") {
		t.Fatalf("rendered messages = %q, should stop showing Thinking once tool call starts", rendered)
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
	if rendered := plainTerminalText(model.renderMessages()); strings.Contains(rendered, "Thinking...") || strings.Contains(rendered, "Working...") {
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
	for _, want := range []string{"cynosure", "/tmp/project", "› hello", "你好", "上下文 --", "╭", "╰"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
	if strings.Contains(view, "工具 0") {
		t.Fatalf("view = %q, should not display tool call count", view)
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
	for _, want := range []string{"cynosure", "/tmp/project", "╭", "╰", "Enter 发送", "上下文 --"} {
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
	for _, forbidden := range []string{"workspace ", "Welcome back!", "API Usage Billing", "Tips for getting started", "Project guide", "Skills 2", "MCP tools 3"} {
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
	for _, want := range []string{`cynosure version: 0.0.1`, `Welcome back`, `/^ ^\`, `/ 0 0 \`, `V\ Y /V`, `/ - \`, `|| (__V`} {
		if !strings.Contains(header, want) {
			t.Fatalf("header = %q, want compact Cynosure version mascot part %q", header, want)
		}
	}
	for _, forbidden := range []string{"/\\_/\\", "( o.o )", "> ^ <", "Doggy Server", "DOGGY API", `^-----^`, `Q /`, `(___\\====`, "Ready ✓"} {
		if strings.Contains(header, forbidden) {
			t.Fatalf("header = %q, should not contain old/oversized mascot part %q", header, forbidden)
		}
	}
	plain := plainTerminalText(header)
	if got := lipgloss.Height(plain); got > 13 {
		t.Fatalf("header height = %d, want compact Cynosure version header height <= 13", got)
	}
	if strings.Index(plain, "Welcome back") > strings.Index(plain, "model deepseek-v4-flash") {
		t.Fatalf("header = %q, want Welcome back above model line", header)
	}
	for _, line := range strings.Split(plain, "\n") {
		if strings.Contains(line, `Cynosure version: 0.0.1`) {
			leftPadding := strings.Index(line, `Cynosure version: 0.0.1`)
			if leftPadding < 30 {
				t.Fatalf("header mascot line = %q, want Cynosure version content centered with substantial left padding", line)
			}
		}
	}
	if headerAccentColor() != tuiPalette.mint {
		t.Fatalf("header accent = %q, want green accent matching mint palette %q", headerAccentColor(), tuiPalette.mint)
	}
}

func TestHeaderBoxClosesWithinTerminalWidth(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/Users/bytedance/golang_pro/cynosure/cynosure", ModelID: "deepseek-v4-flash"})
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
	if strings.Contains(rendered, "cynosure") {
		t.Fatalf("assistant render = %q, should render answer directly without cynosure label", rendered)
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
	for _, want := range []string{"› hello", "cynosure", "你好", "问 cynosure 一件事"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
	if strings.Index(view, "› hello") > strings.Index(view, "你好") {
		t.Fatalf("view = %q, want assistant answer below the submitted user prompt", view)
	}
	if !strings.Contains(model.renderInputArea(), "问 cynosure 一件事") {
		t.Fatalf("input area = %q, want fixed prompt to contain placeholder", model.renderInputArea())
	}
	if strings.Contains(model.renderTranscript(), "问 cynosure 一件事") {
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
	if !strings.Contains(view, "问 cynosure 一件事") {
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
	if !strings.Contains(model.View(), "问 cynosure 一件事") {
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
	if !strings.Contains(model.View(), "问 cynosure 一件事") {
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
	if !strings.Contains(model.View(), "问 cynosure 一件事") {
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

func TestLeftArrowMovesVisibleInputCursorForMiddleEditing(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("abc")})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyLeft})
	model := updated.(Model)

	if !strings.Contains(plainTerminalText(model.renderInput()), "ab"+inputCursor+"c") {
		t.Fatalf("input = %q, want cursor rendered before final character after left arrow", model.renderInput())
	}

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("X")})
	model = updated.(Model)

	if model.input.Value() != "abXc" {
		t.Fatalf("input value = %q, want inserted rune at moved cursor", model.input.Value())
	}
	if !strings.Contains(plainTerminalText(model.renderInput()), "abX"+inputCursor+"c") {
		t.Fatalf("input = %q, want cursor rendered after inserted rune", model.renderInput())
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
	if !strings.Contains(input, "问 cynosure 一件事") {
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

func TestStreamResetDiscardsPartialLiveAssistantContent(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true

	// 流式输出半截内容。
	updated, _ := app.Update(Event{Generation: 1, Name: "assistant_delta", Content: "半截被截断的内容"})
	model := updated.(Model)
	if len(model.messages) != 1 || model.messages[len(model.messages)-1].Content != "半截被截断的内容" {
		t.Fatalf("messages = %#v, want partial assistant content streamed", model.messages)
	}
	if !model.workingStarted {
		t.Fatal("workingStarted should be true once a delta has streamed")
	}

	// 运行时丢弃这段输出并重试，发出 reset，UI 应移除半截内容并恢复等待态。
	updated, _ = model.Update(Event{Generation: 1, Name: "assistant_stream_reset"})
	model = updated.(Model)
	if len(model.messages) != 0 {
		t.Fatalf("messages = %#v, want partial assistant content discarded on reset", model.messages)
	}
	if model.workingStarted {
		t.Fatal("workingStarted should reset so the Thinking indicator returns until retry streams")
	}

	// 重试产生完整内容，应作为一条全新的 assistant 消息呈现。
	updated, _ = model.Update(Event{Generation: 1, Name: "assistant_delta", Content: "完整答案"})
	model = updated.(Model)
	if len(model.messages) != 1 || model.messages[0].Content != "完整答案" {
		t.Fatalf("messages = %#v, want only the retried full answer after reset", model.messages)
	}
}

func TestStreamResetWithoutLiveAssistantIsNoop(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.generation = 1
	app.running = true
	app.messages = []Message{{Role: "user", Content: "问题"}}

	updated, _ := app.Update(Event{Generation: 1, Name: "assistant_stream_reset"})
	model := updated.(Model)
	if len(model.messages) != 1 || model.messages[0].Role != "user" {
		t.Fatalf("messages = %#v, want non-assistant tail preserved on reset", model.messages)
	}
}

func TestEditFileToolMessageRendersRealFileLineNumbers(t *testing.T) {
	dir := t.TempDir()
	// 编辑完成后磁盘上的文件内容：new_text 落在第 5 行起。
	edited := "prefix1\nprefix2\nprefix3\nprefix4\nconst new = \"world\"\nreturn true\nconsole.log(\"done\")\nsuffix1\n"
	if err := os.WriteFile(filepath.Join(dir, "foo.ts"), []byte(edited), 0o644); err != nil {
		t.Fatalf("write edited file: %v", err)
	}
	app := NewModel(nil, SessionInfo{CWD: dir})
	app.width = 120
	app.generation = 1
	app.running = true
	rawArgs := `{"file_path":"foo.ts","old_text":"const old = \"hello\"\nreturn false","new_text":"const new = \"world\"\nreturn true\nconsole.log(\"done\")"}`

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_edit",
		"tool_name":    "edit_file",
		"raw_args":     rawArgs,
		"status":       "success",
	}})
	rendered := plainTerminalText(updated.(Model).renderMessages())

	for _, want := range []string{
		// 删除行使用原文件行号（5、6），新增行使用新文件行号（5、6、7）。
		`- 5│ const old = "hello"`,
		"- 6│ return false",
		`+ 5│ const new = "world"`,
		"+ 6│ return true",
		`+ 7│ console.log("done")`,
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want real-file line number %q", rendered, want)
		}
	}
	// 不应再出现"从 1 开始"的旧行号。
	for _, forbidden := range []string{`- 1│ const old = "hello"`, `+ 1│ const new = "world"`} {
		if strings.Contains(rendered, forbidden) {
			t.Fatalf("rendered = %q, should not number diff from snippet-relative line 1 (%q)", rendered, forbidden)
		}
	}
}

func TestMultiEditToolMessageRendersRealFileLineNumbers(t *testing.T) {
	dir := t.TempDir()
	// 两处编辑完成后的最终文件：alpha new 在第 2 行，beta new 在第 6 行。
	edited := "header\nalpha new\nmid1\nmid2\nmid3\nbeta new\nfooter\n"
	if err := os.WriteFile(filepath.Join(dir, "multi.ts"), []byte(edited), 0o644); err != nil {
		t.Fatalf("write edited file: %v", err)
	}
	app := NewModel(nil, SessionInfo{CWD: dir})
	app.width = 120
	app.generation = 1
	app.running = true
	rawArgs := `{"files":[{"file_path":"multi.ts","edits":[{"old_string":"alpha old","new_string":"alpha new"},{"old_string":"beta old","new_string":"beta new"}]}]}`

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_multi",
		"tool_name":    "multi_edit",
		"raw_args":     rawArgs,
		"status":       "success",
	}})
	rendered := plainTerminalText(updated.(Model).renderMessages())

	for _, want := range []string{
		"- 2│ alpha old",
		"+ 2│ alpha new",
		"- 6│ beta old",
		"+ 6│ beta new",
	} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want per-edit real-file line number %q", rendered, want)
		}
	}
	if strings.Contains(rendered, "- 1│ beta old") || strings.Contains(rendered, "+ 1│ beta new") {
		t.Fatalf("rendered = %q, second edit should not restart numbering at 1", rendered)
	}
}

func TestEditFileToolMessageFallsBackToLineOneWhenFileMissing(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: t.TempDir()})
	app.width = 120
	app.generation = 1
	app.running = true
	// 文件不存在 / new_text 无法定位时，行号回退为 1，且渲染不应 panic。
	rawArgs := `{"file_path":"does-not-exist.ts","old_text":"const old = \"hello\"","new_text":"const new = \"world\""}`

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id": "tool_edit",
		"tool_name":    "edit_file",
		"raw_args":     rawArgs,
		"status":       "success",
	}})
	rendered := plainTerminalText(updated.(Model).renderMessages())

	for _, want := range []string{`- 1│ const old = "hello"`, `+ 1│ const new = "world"`} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want fallback line number %q when file is missing", rendered, want)
		}
	}
}

func TestResumeEditFileToolMessageRendersRealFileLineNumbers(t *testing.T) {
	dir := t.TempDir()
	edited := "prefix1\nprefix2\nprefix3\nprefix4\nconst new = \"world\"\nsuffix1\n"
	if err := os.WriteFile(filepath.Join(dir, "foo.ts"), []byte(edited), 0o644); err != nil {
		t.Fatalf("write edited file: %v", err)
	}
	history := []storage.Message{
		{Role: "user", Content: "改一行"},
		{Role: "assistant", Content: "好的", ToolCalls: []storage.MessageToolCall{
			{ID: "call_1", Type: "function", Function: storage.MessageFunctionCall{Name: "edit_file", Arguments: `{"file_path":"foo.ts","old_text":"const old = \"hello\"","new_text":"const new = \"world\""}`}},
		}},
		{Role: "tool", ToolCallID: "call_1", Content: `{"status":"success","result":"The file foo.ts has been updated successfully."}`},
	}
	resumer := &fakeSessionResumer{
		sessions: []storage.ResumableSession{{SessionID: "session-1", Title: "会话", UpdatedAt: time.Now(), MessageCount: len(history)}},
		conv:     storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: "local-user", Title: "会话"},
		history:  history,
	}
	app := NewModel(nil, SessionInfo{CWD: dir, User: storage.User{ID: "local-user"}, Resumer: resumer})
	app.width = 120
	app.handleSlashCommand("/resume")
	if handled := app.handleResumeSelection("1"); !handled {
		t.Fatal("resume selection was not handled")
	}

	rendered := plainTerminalText(app.renderMessages())
	for _, want := range []string{`- 5│ const old = "hello"`, `+ 5│ const new = "world"`} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("resumed render = %q, want real-file line number %q consistent with live display", rendered, want)
		}
	}
}

func TestWriteFileToolMessageUsesWholeFileLineNumbers(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 120
	rawArgs := `{"file_path":"src/foo.ts","content":"line one\nline two\nline three"}`
	msg := Message{Role: "tool", ToolCall: &ToolCallView{
		Name:    "write_file",
		RawArgs: rawArgs,
		Status:  "success",
	}}

	rendered := plainTerminalText(app.renderMessage(msg))
	for _, want := range []string{"  1│ line one", "  2│ line two", "  3│ line three"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want whole-file write line number %q", rendered, want)
		}
	}
}

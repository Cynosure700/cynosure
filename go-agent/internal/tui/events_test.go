package tui

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/sessions"
)

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
	rendered := model.renderMessages()

	for _, want := range []string{"思考", "先判断是否需要工具"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered messages = %q, want it to contain %q", rendered, want)
		}
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
	view := model.View()
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

func TestViewUsesTerminalTranscriptWithoutConversationBox(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	app.width = 80
	app.height = 24
	app.viewport.Width = 80
	app.viewport.Height = 18
	app.appendMessage("user", "hello")
	app.appendMessage("assistant", "你好")
	app.refreshViewport()

	view := app.View()
	for _, forbidden := range []string{"✦ go-agent", "cwd /tmp/project"} {
		if strings.Contains(view, forbidden) {
			t.Fatalf("view = %q, should not render fixed header %q", view, forbidden)
		}
	}
	for _, want := range []string{"nano, but cozy", "› hello", "go-agent", "你好", "工具 0", "上下文 --"} {
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

func TestViewKeepsWelcomeAndPromptInScrollableTranscript(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	model := updated.(Model)
	model.appendMessage("user", "hello")
	model.appendMessage("assistant", "你好")
	model.refreshViewport()

	view := model.View()
	for _, want := range []string{"nano, but cozy", "› hello", "go-agent", "你好", "问 go-agent 一件事"} {
		if !strings.Contains(view, want) {
			t.Fatalf("view = %q, want it to contain %q", view, want)
		}
	}
	if strings.Index(view, "› hello") > strings.Index(view, "go-agent") {
		t.Fatalf("view = %q, want assistant answer below the submitted user prompt", view)
	}
	if strings.Index(view, "你好") > strings.Index(view, "问 go-agent 一件事") {
		t.Fatalf("view = %q, want next prompt below assistant answer", view)
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

func TestPageUpScrollsTranscriptHistory(t *testing.T) {
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
	if strings.Contains(model.View(), "问 go-agent 一件事") {
		t.Fatalf("view = %q, want page up to reveal history instead of staying at active prompt", model.View())
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
	if strings.Contains(model.View(), "问 go-agent 一件事") {
		t.Fatalf("view = %q, want new content refresh not to force active prompt into view", model.View())
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

	if !strings.Contains(model.View(), "› h") {
		t.Fatalf("view = %q, want inline prompt to show typed text", model.View())
	}
}

func TestTypingSpaceRefreshesInlinePrompt(t *testing.T) {
	app := NewModel(nil, SessionInfo{CWD: "/tmp/project"})
	updated, _ := app.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("h")})
	updated, _ = updated.(Model).Update(tea.KeyMsg{Type: tea.KeySpace, Runes: []rune(" ")})
	model := updated.(Model)

	if !strings.Contains(model.renderInput(), "h "+inputCursor) {
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

func TestRunConfigKeepsTerminalCopyFriendly(t *testing.T) {
	config := newRunConfig(context.Background(), NewModel(nil, SessionInfo{}))

	if config.altScreen || config.mouseCellMotion {
		t.Fatalf("run config = %#v, want no alt screen and no mouse capture so terminal text can be selected and copied", config)
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

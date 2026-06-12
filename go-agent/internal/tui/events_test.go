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
	for _, forbidden := range []string{"╭", "╮", "╰", "╯"} {
		if strings.Contains(view, forbidden) {
			t.Fatalf("view = %q, should not render rounded conversation frame %q", view, forbidden)
		}
	}
	for _, want := range []string{"› hello", "go-agent", "你好", "本轮工具 0", "上下文 --"} {
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

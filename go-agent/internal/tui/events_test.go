package tui

import (
	"context"
	"strings"
	"testing"
	"time"

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

package mcp

import (
	"context"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/storage"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	m := NewManager()
	t.Cleanup(m.Close)
	return m
}

func TestCleanupIdleSessionsKeepsRecentlyUsedWorkspaceSession(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.workspaceSessions["server-1"] = &serverSession{server: storage.MCPServer{ID: "server-1", Name: "recent"}, lastUsedAt: now.Add(-9 * time.Minute)}

	m.cleanupIdleSessions(now)

	if _, ok := m.workspaceSessions["server-1"]; !ok {
		t.Fatalf("recently used workspace session was unexpectedly removed")
	}
}

func TestCleanupIdleSessionsRemovesExpiredWorkspaceSession(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.workspaceSessions["server-1"] = &serverSession{server: storage.MCPServer{ID: "server-1", Name: "expired"}, lastUsedAt: now.Add(-idleTimeout - time.Second)}

	m.cleanupIdleSessions(now)

	if _, ok := m.workspaceSessions["server-1"]; ok {
		t.Fatalf("expired workspace session still exists")
	}
}

func TestCleanupIdleSessionsKeepsActiveWorkspaceCall(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.workspaceSessions["server-1"] = &serverSession{server: storage.MCPServer{ID: "server-1", Name: "active"}, lastUsedAt: now.Add(-idleTimeout - time.Second), activeCalls: 1}

	m.cleanupIdleSessions(now)

	if _, ok := m.workspaceSessions["server-1"]; !ok {
		t.Fatalf("active workspace session was unexpectedly removed")
	}
}

func TestToolsForUserRefreshesWorkspaceLastUsedAt(t *testing.T) {
	m := newTestManager(t)
	old := time.Now().Add(-idleTimeout - time.Minute)
	m.workspaceSessions["server-1"] = &serverSession{
		server:     storage.MCPServer{ID: "server-1", Name: "tools"},
		lastUsedAt: old,
		tools: []openai.Tool{{Type: openai.ToolTypeFunction, Function: &openai.FunctionDefinition{
			Name: "mcp__tools__example",
		}}},
	}

	tools := m.ToolsForUser("user-1")

	if len(tools) != 1 {
		t.Fatalf("expected one MCP tool, got %d", len(tools))
	}
	refreshed := m.workspaceSessions["server-1"].lastUsedAt
	if !refreshed.After(old) {
		t.Fatalf("expected lastUsedAt to be refreshed after %s, got %s", old, refreshed)
	}
}

func TestToolsForUserIncludesBuiltinAndWorkspaceTools(t *testing.T) {
	m := newTestManager(t)
	m.builtinSessions["builtin:docs"] = &serverSession{
		server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"},
		tools: []openai.Tool{{Type: openai.ToolTypeFunction, Function: &openai.FunctionDefinition{
			Name: "mcp__builtin_docs__search",
		}}},
	}
	m.workspaceSessions["workspace:docs"] = &serverSession{
		server: storage.MCPServer{ID: "workspace:docs", Name: "workspace_docs"},
		tools: []openai.Tool{{Type: openai.ToolTypeFunction, Function: &openai.FunctionDefinition{
			Name: "mcp__workspace_docs__search",
		}}},
	}

	tools := m.ToolsForUser("user-1")

	if len(tools) != 2 {
		t.Fatalf("expected builtin and workspace tools, got %d", len(tools))
	}
	if tools[0].Function.Name != "mcp__builtin_docs__search" || tools[1].Function.Name != "mcp__workspace_docs__search" {
		t.Fatalf("unexpected sorted tools: %#v", []string{tools[0].Function.Name, tools[1].Function.Name})
	}
}

func TestCleanupIdleSessionsDoesNotRemoveBuiltinSessions(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.builtinSessions["builtin:docs"] = &serverSession{server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"}, lastUsedAt: now.Add(-idleTimeout - time.Hour)}

	m.cleanupIdleSessions(now)

	if _, ok := m.builtinSessions["builtin:docs"]; !ok {
		t.Fatalf("builtin session should not be removed by idle cleanup")
	}
}

func TestInvalidateIsNoopForLocalSessions(t *testing.T) {
	m := newTestManager(t)
	m.builtinSessions["builtin:docs"] = &serverSession{server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"}}
	m.workspaceSessions["workspace:docs"] = &serverSession{server: storage.MCPServer{ID: "workspace:docs", Name: "workspace_docs"}}

	m.Invalidate("user-1")

	if _, ok := m.builtinSessions["builtin:docs"]; !ok {
		t.Fatalf("builtin session should remain after user invalidate")
	}
	if _, ok := m.workspaceSessions["workspace:docs"]; !ok {
		t.Fatalf("workspace session should remain after user invalidate")
	}
}

func TestCloseDefersActiveWorkspaceSessionClose(t *testing.T) {
	m := NewManager()
	sess := &serverSession{server: storage.MCPServer{ID: "server-1", Name: "active-close"}, activeCalls: 1}
	m.workspaceSessions["server-1"] = sess

	m.Close()

	if len(m.workspaceSessions) != 0 {
		t.Fatalf("expected workspace sessions to be cleared, got %#v", m.workspaceSessions)
	}
	if !sess.closing {
		t.Fatalf("expected active session to be marked closing")
	}
	if sess.activeCalls != 1 {
		t.Fatalf("expected active call to remain in-flight, got %d", sess.activeCalls)
	}

	m.finishCall(sess)

	if sess.activeCalls != 0 {
		t.Fatalf("expected active call count to be decremented after finishCall, got %d", sess.activeCalls)
	}
}

func TestCloseClearsSessionsAndIsIdempotent(t *testing.T) {
	m := NewManager()
	m.workspaceSessions["server-1"] = &serverSession{server: storage.MCPServer{ID: "server-1", Name: "close"}}
	m.builtinSessions["builtin:docs"] = &serverSession{server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"}}

	m.Close()
	m.Close()

	if len(m.workspaceSessions) != 0 {
		t.Fatalf("expected workspace sessions to be cleared, got %#v", m.workspaceSessions)
	}
	if len(m.builtinSessions) != 0 {
		t.Fatalf("expected builtin sessions to be cleared, got %#v", m.builtinSessions)
	}
}

func TestManagerSnapshotIncludesWorkspaceServers(t *testing.T) {
	m := newTestManager(t)
	m.SetWorkspaceServers(context.Background(), []storage.MCPServer{{ID: "workspace:disabled", Name: "disabled", Transport: "stdio", Command: "cmd", Enabled: false}})

	snapshot := m.Snapshot("local-user")
	if len(snapshot.Servers) != 1 {
		t.Fatalf("expected one workspace server in snapshot, got %#v", snapshot.Servers)
	}
	server := snapshot.Servers[0]
	if server.Name != "disabled" || server.Scope != "workspace" || server.Transport != "stdio" {
		t.Fatalf("unexpected server snapshot: %#v", server)
	}
	if server.Connected || server.Enabled {
		t.Fatalf("disabled server should be disabled and disconnected: %#v", server)
	}
}

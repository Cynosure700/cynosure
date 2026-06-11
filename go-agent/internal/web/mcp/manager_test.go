package mcp

import (
	"context"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/web/storage"
)

type fakeMCPStore struct{}

func (fakeMCPStore) ListEnabledMCPServersByUser(ctx context.Context, userID string) ([]storage.MCPServer, error) {
	return nil, nil
}

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	m := NewManager(fakeMCPStore{})
	t.Cleanup(m.Close)
	return m
}

func TestCleanupIdleSessionsKeepsRecentlyUsedSession(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {
			server:     storage.MCPServer{ID: "server-1", Name: "recent"},
			lastUsedAt: now.Add(-9 * time.Minute),
		},
	}

	m.cleanupIdleSessions(now)

	if _, ok := m.sessions["user-1"]["server-1"]; !ok {
		t.Fatalf("recently used session was unexpectedly removed")
	}
}

func TestCleanupIdleSessionsRemovesExpiredSession(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {
			server:     storage.MCPServer{ID: "server-1", Name: "expired"},
			lastUsedAt: now.Add(-idleTimeout - time.Second),
		},
	}

	m.cleanupIdleSessions(now)

	if _, ok := m.sessions["user-1"]; ok {
		t.Fatalf("expired session user bucket still exists: %#v", m.sessions["user-1"])
	}
}

func TestCleanupIdleSessionsKeepsActiveCall(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {
			server:      storage.MCPServer{ID: "server-1", Name: "active"},
			lastUsedAt:  now.Add(-idleTimeout - time.Second),
			activeCalls: 1,
		},
	}

	m.cleanupIdleSessions(now)

	if _, ok := m.sessions["user-1"]["server-1"]; !ok {
		t.Fatalf("active session was unexpectedly removed")
	}
}

func TestToolsForUserRefreshesLastUsedAt(t *testing.T) {
	m := newTestManager(t)
	old := time.Now().Add(-idleTimeout - time.Minute)
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {
			server:     storage.MCPServer{ID: "server-1", Name: "tools"},
			lastUsedAt: old,
			tools: []openai.Tool{{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name: "mcp__tools__example",
				},
			}},
		},
	}

	tools := m.ToolsForUser("user-1")

	if len(tools) != 1 {
		t.Fatalf("expected one MCP tool, got %d", len(tools))
	}
	refreshed := m.sessions["user-1"]["server-1"].lastUsedAt
	if !refreshed.After(old) {
		t.Fatalf("expected lastUsedAt to be refreshed after %s, got %s", old, refreshed)
	}
}

func TestToolsForUserIncludesBuiltinTools(t *testing.T) {
	m := newTestManager(t)
	m.builtinSessions["builtin:docs"] = &serverSession{
		server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"},
		tools: []openai.Tool{{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name: "mcp__builtin_docs__search",
			},
		}},
	}
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {
			server: storage.MCPServer{ID: "server-1", Name: "user_docs"},
			tools: []openai.Tool{{
				Type: openai.ToolTypeFunction,
				Function: &openai.FunctionDefinition{
					Name: "mcp__user_docs__search",
				},
			}},
		},
	}

	tools := m.ToolsForUser("user-1")

	if len(tools) != 2 {
		t.Fatalf("expected builtin and user tools, got %d", len(tools))
	}
	if tools[0].Function.Name != "mcp__builtin_docs__search" || tools[1].Function.Name != "mcp__user_docs__search" {
		t.Fatalf("unexpected sorted tools: %#v", []string{tools[0].Function.Name, tools[1].Function.Name})
	}
}

func TestCleanupIdleSessionsDoesNotRemoveBuiltinSessions(t *testing.T) {
	now := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)
	m := newTestManager(t)
	m.builtinSessions["builtin:docs"] = &serverSession{
		server:     storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"},
		lastUsedAt: now.Add(-idleTimeout - time.Hour),
	}

	m.cleanupIdleSessions(now)

	if _, ok := m.builtinSessions["builtin:docs"]; !ok {
		t.Fatalf("builtin session should not be removed by idle cleanup")
	}
}

func TestInvalidateDoesNotRemoveBuiltinSessions(t *testing.T) {
	m := newTestManager(t)
	m.builtinSessions["builtin:docs"] = &serverSession{server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"}}
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {server: storage.MCPServer{ID: "server-1", Name: "user"}},
	}

	m.Invalidate("user-1")

	if _, ok := m.sessions["user-1"]; ok {
		t.Fatalf("expected user sessions to be removed")
	}
	if _, ok := m.builtinSessions["builtin:docs"]; !ok {
		t.Fatalf("builtin session should remain after user invalidate")
	}
}

func TestInvalidateDefersActiveSessionClose(t *testing.T) {
	m := newTestManager(t)
	sess := &serverSession{
		server:      storage.MCPServer{ID: "server-1", Name: "active"},
		activeCalls: 1,
	}
	m.sessions["user-1"] = map[string]*serverSession{"server-1": sess}

	m.Invalidate("user-1")

	if _, ok := m.sessions["user-1"]; ok {
		t.Fatalf("expected invalidated user sessions to be removed from manager")
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

func TestCloseDefersActiveSessionClose(t *testing.T) {
	m := NewManager(fakeMCPStore{})
	sess := &serverSession{
		server:      storage.MCPServer{ID: "server-1", Name: "active-close"},
		activeCalls: 1,
	}
	m.sessions["user-1"] = map[string]*serverSession{"server-1": sess}

	m.Close()

	if len(m.sessions) != 0 {
		t.Fatalf("expected sessions to be cleared, got %#v", m.sessions)
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
	m := NewManager(fakeMCPStore{})
	m.sessions["user-1"] = map[string]*serverSession{
		"server-1": {server: storage.MCPServer{ID: "server-1", Name: "close"}},
	}
	m.builtinSessions["builtin:docs"] = &serverSession{server: storage.MCPServer{ID: "builtin:docs", Name: "builtin_docs"}}

	m.Close()
	m.Close()

	if len(m.sessions) != 0 {
		t.Fatalf("expected sessions to be cleared, got %#v", m.sessions)
	}
	if len(m.builtinSessions) != 0 {
		t.Fatalf("expected builtin sessions to be cleared, got %#v", m.builtinSessions)
	}
}

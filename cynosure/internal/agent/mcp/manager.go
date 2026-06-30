package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/logger"
)

const (
	toolNamePrefix  = "mcp__"
	connectTimeout  = 30 * time.Second
	callTimeout     = 120 * time.Second
	idleTimeout     = 10 * time.Minute
	cleanupInterval = 1 * time.Minute
)

// serverSession 持有一个已连接 MCP 服务器的会话与其工具元信息。
type serverSession struct {
	server      storage.MCPServer
	signature   string // 配置指纹，用于判断配置是否变更需重连
	session     *mcpsdk.ClientSession
	tools       []openai.Tool     // 已转换并加前缀的工具定义
	toolNames   map[string]string // 前缀名 -> 原始 MCP 工具名
	lastUsedAt  time.Time
	activeCalls int
	closing     bool
}

type ServerStatus struct {
	ID        string
	Name      string
	Scope     string
	Transport string
	Command   string
	Args      []string
	URL       string
	Enabled   bool
	Connected bool
	ToolCount int
	LastError string
}

type Snapshot struct {
	Servers    []ServerStatus
	ToolCount  int
	ErrorCount int
}

// Manager 管理本地 MCP 客户端连接，提供工具发现与调用能力。
type Manager struct {
	mu                sync.Mutex
	done              chan struct{}
	closeOnce         sync.Once
	builtinServers    map[string]storage.MCPServer
	builtinSessions   map[string]*serverSession
	workspaceServers  map[string]storage.MCPServer
	workspaceSessions map[string]*serverSession
	workspaceErrors   map[string]string
}

func NewManager() *Manager {
	manager := &Manager{
		done:              make(chan struct{}),
		builtinServers:    make(map[string]storage.MCPServer),
		builtinSessions:   make(map[string]*serverSession),
		workspaceServers:  make(map[string]storage.MCPServer),
		workspaceSessions: make(map[string]*serverSession),
		workspaceErrors:   make(map[string]string),
	}
	go manager.cleanupLoop()
	return manager
}

func serverSignature(s storage.MCPServer) string {
	data, _ := json.Marshal(struct {
		Transport string
		Command   string
		Args      []string
		Env       map[string]string
		URL       string
		Headers   map[string]string
	}{s.Transport, s.Command, s.Args, s.Env, s.URL, s.Headers})
	return string(data)
}

func sanitizeName(name string) string {
	var sb strings.Builder
	for _, r := range name {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '_':
			sb.WriteRune(r)
		default:
			sb.WriteRune('_')
		}
	}
	if sb.Len() == 0 {
		return "server"
	}
	return sb.String()
}

func prefixedToolName(serverName, toolName string) string {
	return toolNamePrefix + sanitizeName(serverName) + "__" + toolName
}

func (m *Manager) SetBuiltinServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.builtinSessions
	m.builtinServers = make(map[string]storage.MCPServer, len(servers))
	m.builtinSessions = make(map[string]*serverSession)
	for _, server := range servers {
		m.builtinServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureBuiltinSessions(ctx)
}

func (m *Manager) EnsureBuiltinSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.builtinServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.builtinSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectBuiltinAndDiscover(ctx, server)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: builtin stdio connect failed server=%s: %v", server.Name, err))
			continue
		}
		m.mu.Lock()
		if old := m.builtinSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.builtinSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: builtin stdio connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func (m *Manager) SetWorkspaceServers(ctx context.Context, servers []storage.MCPServer) {
	m.mu.Lock()
	old := m.workspaceSessions
	m.workspaceServers = make(map[string]storage.MCPServer, len(servers))
	m.workspaceSessions = make(map[string]*serverSession)
	m.workspaceErrors = make(map[string]string)
	for _, server := range servers {
		m.workspaceServers[server.ID] = server
	}
	m.mu.Unlock()

	for _, sess := range old {
		closeSession(sess)
	}
	m.EnsureWorkspaceSessions(ctx)
}

func (m *Manager) EnsureWorkspaceSessions(ctx context.Context) {
	m.mu.Lock()
	missing := make([]storage.MCPServer, 0)
	for id, server := range m.workspaceServers {
		if !server.Enabled {
			continue
		}
		if _, ok := m.workspaceSessions[id]; !ok {
			missing = append(missing, server)
		}
	}
	m.mu.Unlock()

	for _, server := range missing {
		sess, err := connectAndDiscover(ctx, server)
		m.mu.Lock()
		if err != nil {
			m.workspaceErrors[server.ID] = err.Error()
			m.mu.Unlock()
			logger.Warn(fmt.Sprintf("mcp: workspace connect failed server=%s: %v", server.Name, err))
			continue
		}
		delete(m.workspaceErrors, server.ID)
		if old := m.workspaceSessions[server.ID]; old != nil {
			closeSession(old)
		}
		m.workspaceSessions[server.ID] = sess
		m.mu.Unlock()
		logger.Info(fmt.Sprintf("mcp: workspace connected server=%s tools=%d", server.Name, len(sess.tools)))
	}
}

func connectAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildTransport(server)
	if err != nil {
		return nil, err
	}
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "cynosure", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func connectBuiltinAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildBuiltinStdioTransport(server)
	if err != nil {
		return nil, err
	}
	return connectWithTransport(ctx, server, transport)
}

func connectWithTransport(ctx context.Context, server storage.MCPServer, transport mcpsdk.Transport) (*serverSession, error) {
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "cynosure", Version: "1.0.0"}, nil)
	session, err := client.Connect(connCtx, transport, nil)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	tools, names, err := discoverTools(connCtx, session, server.Name)
	if err != nil {
		_ = session.Close()
		return nil, fmt.Errorf("list tools: %w", err)
	}
	return &serverSession{
		server:     server,
		signature:  serverSignature(server),
		session:    session,
		tools:      tools,
		toolNames:  names,
		lastUsedAt: time.Now(),
	}, nil
}

func discoverTools(ctx context.Context, session *mcpsdk.ClientSession, serverName string) ([]openai.Tool, map[string]string, error) {
	result, err := session.ListTools(ctx, nil)
	if err != nil {
		return nil, nil, err
	}
	tools := make([]openai.Tool, 0, len(result.Tools))
	names := make(map[string]string, len(result.Tools))
	for _, tool := range result.Tools {
		prefixed := prefixedToolName(serverName, tool.Name)
		params, err := json.Marshal(tool.InputSchema)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: marshal schema failed server=%s tool=%s: %v", serverName, tool.Name, err))
			continue
		}
		tools = append(tools, openai.Tool{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        prefixed,
				Description: tool.Description,
				Parameters:  json.RawMessage(params),
			},
		})
		names[prefixed] = tool.Name
	}
	return tools, names, nil
}

// ToolsForUser 返回所有已连接本地 MCP 服务器发现到的工具定义（带前缀）。
func (m *Manager) ToolsForUser(userID string) []openai.Tool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if len(m.builtinSessions) == 0 && len(m.workspaceSessions) == 0 {
		return nil
	}
	tools := make([]openai.Tool, 0)
	now := time.Now()
	for _, sess := range m.builtinSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	for _, sess := range m.workspaceSessions {
		sess.lastUsedAt = now
		tools = append(tools, sess.tools...)
	}
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Function.Name < tools[j].Function.Name
	})
	return tools
}

// CallTool 把带前缀的工具名路由到对应 MCP 服务器执行，返回文本结果。
func (m *Manager) CallTool(ctx context.Context, userID, prefixedName, rawArgs string) (string, error) {
	m.mu.Lock()
	var target *serverSession
	var originalName string
	var builtinID string
	for _, sess := range m.workspaceSessions {
		if name, ok := sess.toolNames[prefixedName]; ok {
			target = sess
			originalName = name
			target.activeCalls++
			break
		}
	}
	if target == nil {
		for id, sess := range m.builtinSessions {
			if name, ok := sess.toolNames[prefixedName]; ok {
				target = sess
				originalName = name
				builtinID = id
				target.activeCalls++
				break
			}
		}
	}
	m.mu.Unlock()

	if target == nil {
		return "", fmt.Errorf("mcp tool %s not found", prefixedName)
	}
	defer m.finishCall(target)

	var args map[string]any
	if strings.TrimSpace(rawArgs) != "" {
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return "", fmt.Errorf("invalid tool arguments: %w", err)
		}
	}
	if builtinID != "" {
		return m.callBuiltinTool(ctx, builtinID, target, originalName, args)
	}

	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()

	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) callBuiltinTool(ctx context.Context, builtinID string, target *serverSession, originalName string, args map[string]any) (string, error) {
	output, err := callSessionTool(ctx, target, originalName, args)
	if err == nil {
		return output, nil
	}
	logger.Warn(fmt.Sprintf("mcp: builtin stdio call failed, reconnecting server=%s: %v", target.server.Name, err))
	m.reconnectBuiltin(ctx, builtinID, target)

	m.mu.Lock()
	retryTarget := m.builtinSessions[builtinID]
	if retryTarget != nil {
		retryTarget.activeCalls++
	}
	m.mu.Unlock()
	if retryTarget == nil {
		return "", err
	}
	defer m.finishCall(retryTarget)
	return callSessionTool(ctx, retryTarget, originalName, args)
}

func callSessionTool(ctx context.Context, target *serverSession, originalName string, args map[string]any) (string, error) {
	callCtx, cancel := context.WithTimeout(ctx, callTimeout)
	defer cancel()
	result, err := target.session.CallTool(callCtx, &mcpsdk.CallToolParams{Name: originalName, Arguments: args})
	if err != nil {
		return "", err
	}
	output := serializeContent(result)
	if result.IsError {
		return "", fmt.Errorf("mcp tool error: %s", output)
	}
	return output, nil
}

func (m *Manager) reconnectBuiltin(ctx context.Context, builtinID string, old *serverSession) {
	m.mu.Lock()
	server, ok := m.builtinServers[builtinID]
	if current := m.builtinSessions[builtinID]; current == old {
		delete(m.builtinSessions, builtinID)
	}
	m.mu.Unlock()
	if !ok || !server.Enabled {
		return
	}
	closeSession(old)
	sess, err := connectBuiltinAndDiscover(ctx, server)
	if err != nil {
		logger.Warn(fmt.Sprintf("mcp: builtin stdio reconnect failed server=%s: %v", server.Name, err))
		return
	}
	m.mu.Lock()
	m.builtinSessions[builtinID] = sess
	m.mu.Unlock()
	logger.Info(fmt.Sprintf("mcp: builtin stdio reconnected server=%s tools=%d", server.Name, len(sess.tools)))
}

// Invalidate 保留兼容入口。本地 TUI 不再维护用户 DB MCP 连接，因此该方法不影响 builtin/workspace 会话。
func (m *Manager) Invalidate(userID string) {
}

// TestServer 临时连接一个配置并发现工具，返回工具名列表，供连接测试使用。
func (m *Manager) TestServer(ctx context.Context, server storage.MCPServer) ([]string, error) {
	sess, err := connectAndDiscover(ctx, server)
	if err != nil {
		return nil, err
	}
	defer sess.session.Close()
	names := make([]string, 0, len(sess.toolNames))
	for _, original := range sess.toolNames {
		names = append(names, original)
	}
	sort.Strings(names)
	return names, nil
}

// Close 关闭全部连接，用于服务退出。
func (m *Manager) Close() {
	m.closeOnce.Do(func() { close(m.done) })
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.builtinSessions {
		m.closeWhenIdleLocked(sess)
	}
	for _, sess := range m.workspaceSessions {
		m.closeWhenIdleLocked(sess)
	}
	m.builtinSessions = make(map[string]*serverSession)
	m.workspaceSessions = make(map[string]*serverSession)
}

func (m *Manager) Snapshot(userID string) Snapshot {
	m.mu.Lock()
	defer m.mu.Unlock()
	statuses := make([]ServerStatus, 0, len(m.builtinServers)+len(m.workspaceServers))
	addConfigured := func(scope string, server storage.MCPServer, sess *serverSession, lastErr string) {
		status := ServerStatus{ID: server.ID, Name: server.Name, Scope: scope, Transport: server.Transport, Command: server.Command, Args: append([]string(nil), server.Args...), URL: server.URL, Enabled: server.Enabled, Connected: sess != nil, LastError: lastErr}
		if sess != nil {
			status.ToolCount = len(sess.tools)
		}
		statuses = append(statuses, status)
	}
	for id, server := range m.builtinServers {
		addConfigured("builtin", server, m.builtinSessions[id], "")
	}
	for id, server := range m.workspaceServers {
		addConfigured("workspace", server, m.workspaceSessions[id], m.workspaceErrors[id])
	}
	sort.Slice(statuses, func(i, j int) bool {
		if statuses[i].Scope != statuses[j].Scope {
			return statuses[i].Scope < statuses[j].Scope
		}
		return statuses[i].Name < statuses[j].Name
	})
	snapshot := Snapshot{Servers: statuses}
	for _, status := range statuses {
		snapshot.ToolCount += status.ToolCount
		if status.LastError != "" {
			snapshot.ErrorCount++
		}
	}
	return snapshot
}

func (m *Manager) finishCall(sess *serverSession) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if sess.activeCalls > 0 {
		sess.activeCalls--
	}
	sess.lastUsedAt = time.Now()
	if sess.closing && sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) closeWhenIdleLocked(sess *serverSession) {
	if sess == nil {
		return
	}
	sess.closing = true
	if sess.activeCalls == 0 {
		closeSession(sess)
	}
}

func (m *Manager) cleanupLoop() {
	ticker := time.NewTicker(cleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			m.cleanupIdleSessions(time.Now())
		case <-m.done:
			return
		}
	}
}

func (m *Manager) cleanupIdleSessions(now time.Time) {
	type idleSession struct {
		scope string
		id    string
		sess  *serverSession
	}
	var idle []idleSession

	m.mu.Lock()
	for serverID, sess := range m.workspaceSessions {
		if sess.activeCalls > 0 || now.Sub(sess.lastUsedAt) < idleTimeout {
			continue
		}
		idle = append(idle, idleSession{scope: "workspace", id: serverID, sess: sess})
		delete(m.workspaceSessions, serverID)
	}
	m.mu.Unlock()

	for _, item := range idle {
		closeSession(item.sess)
		logger.Info(fmt.Sprintf("mcp: idle session closed scope=%s id=%s server=%s", item.scope, item.id, item.sess.server.Name))
	}
}

func closeSession(sess *serverSession) {
	if sess == nil || sess.session == nil {
		return
	}
	_ = sess.session.Close()
}

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

	"nano_cc/internal/logger"
	"nano_cc/internal/web/storage"
)

const (
	toolNamePrefix = "mcp__"
	connectTimeout = 30 * time.Second
	callTimeout    = 120 * time.Second
)

// Store 是 Manager 依赖的存储能力子集。
type Store interface {
	ListEnabledMCPServersByUser(ctx context.Context, userID string) ([]storage.MCPServer, error)
}

// serverSession 持有一个已连接 MCP 服务器的会话与其工具元信息。
type serverSession struct {
	server    storage.MCPServer
	signature string // 配置指纹，用于判断配置是否变更需重连
	session   *mcpsdk.ClientSession
	tools     []openai.Tool     // 已转换并加前缀的工具定义
	toolNames map[string]string // 前缀名 -> 原始 MCP 工具名
}

// Manager 按用户维度管理 MCP 客户端连接，提供工具发现与调用能力。
type Manager struct {
	store Store
	mu    sync.Mutex
	// userID -> (serverID -> session)
	sessions map[string]map[string]*serverSession
}

func NewManager(store Store) *Manager {
	return &Manager{store: store, sessions: make(map[string]map[string]*serverSession)}
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

// EnsureUserSessions 懒加载：对该用户的启用配置建立/复用连接，关闭已删除或禁用的连接。
// 单个服务器连接失败不影响其他服务器，仅记录日志并跳过。
func (m *Manager) EnsureUserSessions(ctx context.Context, userID string) {
	servers, err := m.store.ListEnabledMCPServersByUser(ctx, userID)
	if err != nil {
		logger.Warn(fmt.Sprintf("mcp: list enabled servers failed user=%s: %v", userID, err))
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.sessions[userID]
	if current == nil {
		current = make(map[string]*serverSession)
		m.sessions[userID] = current
	}

	wanted := make(map[string]struct{}, len(servers))
	for _, server := range servers {
		wanted[server.ID] = struct{}{}
		sig := serverSignature(server)
		if existing, ok := current[server.ID]; ok && existing.signature == sig {
			continue // 配置未变，复用连接
		}
		if existing, ok := current[server.ID]; ok {
			_ = existing.session.Close() // 配置已变，先关闭旧连接
			delete(current, server.ID)
		}
		sess, err := connectAndDiscover(ctx, server)
		if err != nil {
			logger.Warn(fmt.Sprintf("mcp: connect failed user=%s server=%s: %v", userID, server.Name, err))
			continue
		}
		current[server.ID] = sess
	}

	// 关闭已删除/禁用的连接
	for id, sess := range current {
		if _, ok := wanted[id]; !ok {
			_ = sess.session.Close()
			delete(current, id)
		}
	}
}

func connectAndDiscover(ctx context.Context, server storage.MCPServer) (*serverSession, error) {
	transport, err := buildTransport(server)
	if err != nil {
		return nil, err
	}
	connCtx, cancel := context.WithTimeout(ctx, connectTimeout)
	defer cancel()

	client := mcpsdk.NewClient(&mcpsdk.Implementation{Name: "nano_cc", Version: "1.0.0"}, nil)
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
		server:    server,
		signature: serverSignature(server),
		session:   session,
		tools:     tools,
		toolNames: names,
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

// ToolsForUser 返回该用户所有已连接 MCP 服务器发现到的工具定义（带前缀）。
func (m *Manager) ToolsForUser(userID string) []openai.Tool {
	m.mu.Lock()
	defer m.mu.Unlock()
	current := m.sessions[userID]
	if len(current) == 0 {
		return nil
	}
	tools := make([]openai.Tool, 0)
	for _, sess := range current {
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
	for _, sess := range m.sessions[userID] {
		if name, ok := sess.toolNames[prefixedName]; ok {
			target = sess
			originalName = name
			break
		}
	}
	m.mu.Unlock()

	if target == nil {
		return "", fmt.Errorf("mcp tool %s not found", prefixedName)
	}

	var args map[string]any
	if strings.TrimSpace(rawArgs) != "" {
		if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
			return "", fmt.Errorf("invalid tool arguments: %w", err)
		}
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

// Invalidate 关闭并清除该用户的所有连接，下次对话时按最新配置重连。
func (m *Manager) Invalidate(userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sess := range m.sessions[userID] {
		_ = sess.session.Close()
	}
	delete(m.sessions, userID)
}

// TestServer 临时连接一个配置并发现工具，返回工具名列表，供前端“测试连接”使用。
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
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, sessions := range m.sessions {
		for _, sess := range sessions {
			_ = sess.session.Close()
		}
	}
	m.sessions = make(map[string]map[string]*serverSession)
}

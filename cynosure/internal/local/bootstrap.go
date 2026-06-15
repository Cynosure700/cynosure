package local

import (
	"context"
	"fmt"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/idgen"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
)

const LocalUserID = "local-user"

type Bundle struct {
	Runtime      *runtime.Service
	Store        *Store
	MCP          *mcp.Manager
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	SkillCount   int
	Skills       []sessions.SkillSummary
	MCPToolCount int
	MCPServers   []mcp.ServerStatus
}

func Bootstrap(ctx context.Context, cwd string) (*Bundle, error) {
	cfg, err := config.LoadLocalConfig(cwd)
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(cfg.LogsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}
	userSkillsDir, err := config.CynosureSkillsDir()
	if err != nil {
		return nil, err
	}
	builtinSkills, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: userSkillsDir, Source: "user"},
		{Path: config.WorkspaceCynosureSkillsDir(cfg.WorkspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, fmt.Errorf("load cynosure skills: %w", err)
	}
	basePrompt, err := assistant.LoadBaseSystemPrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	cynosureMarkdown, err := config.LoadCynosureMarkdownContext(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("load CYNOSURE.MD context: %w", err)
	}
	store, err := NewStoreWithMemory(cfg.WorkspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("init memory store: %w", err)
	}
	client := llm.NewDeepseekClient(cfg.LLM.BaseURL, cfg.LLM.APIKey)
	runtimeService := runtime.NewService(store, cfg, client)
	runtimeService.EnableMemory = true
	runtimeService.SetBuiltinSkills(builtinSkills)
	runtimeService.SetBasePrompt(basePrompt)
	runtimeService.SetCynosureMarkdownContext(cynosureMarkdown)
	mcpManager := mcp.NewManager()
	workspaceMCPServers, err := mcp.LoadWorkspaceConfig(config.WorkspaceMCPConfigPath(cfg.WorkspaceRoot))
	if err != nil {
		mcpManager.Close()
		return nil, fmt.Errorf("load workspace mcp config: %w", err)
	}
	mcpManager.SetWorkspaceServers(ctx, workspaceMCPServers)
	runtimeService.SetMCPManager(mcpManager)
	user := storage.User{ID: LocalUserID, Email: "local@cynosure", Username: "local", MemoryEnabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	conversation := storage.Conversation{ID: idgen.New("conv"), SessionID: idgen.UUID(), UserID: user.ID, RootMessageID: idgen.New("msg"), Title: "TUI 会话", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := store.CreateConversation(ctx, conversation); err != nil {
		mcpManager.Close()
		return nil, err
	}
	mcpManager.EnsureWorkspaceSessions(ctx)
	mcpTools := mcpManager.ToolsForUser(user.ID)
	mcpSnapshot := mcpManager.Snapshot(user.ID)
	skills := builtinSkills.Summaries()
	return &Bundle{Runtime: runtimeService, Store: store, MCP: mcpManager, User: user, Conversation: conversation, CWD: cfg.WorkspaceRoot, SkillCount: len(skills), Skills: skills, MCPToolCount: len(mcpTools), MCPServers: mcpSnapshot.Servers}, nil
}

func (b *Bundle) Close() {
	if b != nil && b.MCP != nil {
		b.MCP.Close()
	}
}

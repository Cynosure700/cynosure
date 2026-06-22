package local

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cynosure/assets"
	"cynosure/internal/agent/mcp"
	"cynosure/internal/agent/runtime"
	"cynosure/internal/agent/storage"
	"cynosure/internal/assistant"
	"cynosure/internal/config"
	"cynosure/internal/idgen"
	"cynosure/internal/llm"
	"cynosure/internal/logger"
	"cynosure/internal/sessions"
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
	sessionID := idgen.UUID()
	logsDir, err := config.CynosureSessionLogsDir(cfg.WorkspaceRoot, sessionID)
	if err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(logsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	}
	userSkillsDir, err := config.CynosureSkillsDir()
	if err != nil {
		return nil, err
	}
	builtinSkills, err := loadSkills(cfg.WorkspaceRoot, userSkillsDir)
	if err != nil {
		return nil, fmt.Errorf("load cynosure skills: %w", err)
	}
	basePrompt, err := loadBasePrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	functionalPrompts, err := runtime.LoadFunctionalPrompts()
	if err != nil {
		return nil, err
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
	runtimeService.Prompts = functionalPrompts
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
	user := storage.User{ID: LocalUserID, Username: "local", MemoryEnabled: true, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	conversation := storage.Conversation{ID: idgen.New("conv"), SessionID: sessionID, UserID: user.ID, Title: "TUI 会话", CreatedAt: time.Now(), UpdatedAt: time.Now()}
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

// loadSkills 合并内置（嵌入二进制）、用户级与工作区级 skills。
// 优先级：workspace > user > builtin（后者被前者同名覆盖）。
func loadSkills(workspaceRoot, userSkillsDir string) (*sessions.SkillLoader, error) {
	builtin, err := sessions.LoadSkillsFromFS(assets.BuiltinSkillsFS(), "builtin")
	if err != nil {
		return nil, err
	}
	userAndWorkspace, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: userSkillsDir, Source: "user"},
		{Path: config.WorkspaceCynosureSkillsDir(workspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, err
	}
	return sessions.MergeSkillLoaders(builtin, userAndWorkspace), nil
}

// loadBasePrompt 优先使用用户覆盖文件（~/.cynosure/system_prompt.md），
// 否则使用嵌入二进制的内置 system prompt。
func loadBasePrompt(overridePath string) (string, error) {
	if strings.TrimSpace(overridePath) != "" {
		return assistant.LoadBaseSystemPrompt(overridePath)
	}
	return assets.SystemPrompt(), nil
}

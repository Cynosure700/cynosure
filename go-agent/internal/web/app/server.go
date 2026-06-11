package app

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/mcp"
	"nano_cc/internal/web/runtime"
	"nano_cc/internal/web/storage"
)

type serverStore interface {
	HealthCheck(ctx context.Context) error
	RunMigrations(ctx context.Context) error
	UpdateUserMemoryEnabled(ctx context.Context, userID string, enabled bool) error
	ListSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error)
	CreateSkill(ctx context.Context, skill storage.Skill) error
	GetSkillByID(ctx context.Context, skillID string) (storage.Skill, error)
	UpdateSkill(ctx context.Context, skill storage.Skill) error
	DeleteSkill(ctx context.Context, skillID string) error
	ListMCPServersByUser(ctx context.Context, userID string) ([]storage.MCPServer, error)
	CreateMCPServer(ctx context.Context, server storage.MCPServer) error
	GetMCPServerByID(ctx context.Context, id string) (storage.MCPServer, error)
	UpdateMCPServer(ctx context.Context, server storage.MCPServer) error
	DeleteMCPServer(ctx context.Context, id string) error
	ListConversationsByUser(ctx context.Context, userID string) ([]storage.Conversation, error)
	CreateConversation(ctx context.Context, conversation storage.Conversation) error
	GetConversationByID(ctx context.Context, conversationID string) (storage.Conversation, error)
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	DeleteConversation(ctx context.Context, conversationID string) error
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
}

type Server struct {
	cfg           config.AppConfig
	store         serverStore
	authService   *auth.Service
	runtime       *runtime.Service
	builtinSkills *sessions.SkillLoader
	mcpManager    *mcp.Manager
	mux           *http.ServeMux
}

func NewServer() (*Server, error) {
	cfg, err := config.LoadWebConfig()
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	llmClient := llm.NewDeepseekClient(cfg.LLM.BaseURL, cfg.LLM.APIKey)
	store, err := storage.NewStore(cfg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := store.HealthCheck(ctx); err != nil {
		return nil, err
	}
	if err := store.RunMigrations(ctx); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerAt(cfg.LogsDir); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	} else {
		logger.Info(fmt.Sprintf("LLM logs -> %s", logger.LogFilePath()))
	}
	builtinSkills, err := sessions.LoadBuiltinSkillsFromDir(cfg.BuiltinSkillsDir)
	if err != nil {
		return nil, fmt.Errorf("load builtin skills: %w", err)
	}
	basePrompt, err := assistant.LoadBaseSystemPrompt(cfg.SystemPromptPath)
	if err != nil {
		return nil, fmt.Errorf("load system prompt: %w", err)
	}
	runtimeService := runtime.NewService(store, cfg, llmClient)
	runtimeService.SetBuiltinSkills(builtinSkills)
	runtimeService.SetBasePrompt(basePrompt)
	mcpManager := mcp.NewManager(store)
	builtinMCPServers, err := mcp.LoadBuiltinConfig(filepath.Join(cfg.AppHome, "mcp_config.json"))
	if err != nil {
		mcpManager.Close()
		return nil, fmt.Errorf("load builtin mcp config: %w", err)
	}
	builtinCtx, builtinCancel := context.WithTimeout(context.Background(), 35*time.Second)
	defer builtinCancel()
	mcpManager.SetBuiltinServers(builtinCtx, builtinMCPServers)
	runtimeService.SetMCPManager(mcpManager)
	server := &Server{
		cfg:           cfg,
		store:         store,
		authService:   auth.NewService(store, cfg),
		runtime:       runtimeService,
		builtinSkills: builtinSkills,
		mcpManager:    mcpManager,
		mux:           http.NewServeMux(),
	}
	server.routes()
	return server, nil
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.cfg.ServerAddr, s.withCORS(s.mux))
}

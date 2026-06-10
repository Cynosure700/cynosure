package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/auth"
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
	server := &Server{
		cfg:           cfg,
		store:         store,
		authService:   auth.NewService(store, cfg),
		runtime:       runtimeService,
		builtinSkills: builtinSkills,
		mux:           http.NewServeMux(),
	}
	server.routes()
	return server, nil
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.cfg.ServerAddr, s.withCORS(s.mux))
}

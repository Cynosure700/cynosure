package runtime

import (
	"context"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/storage"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	CreateMessage(ctx context.Context, message storage.Message) error
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	ListEnabledSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error)
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
}

type Service struct {
	Store         conversationStore
	Cfg           config.AppConfig
	Tools         *ToolRegistry
	BuiltinSkills *sessions.SkillLoader
	BasePrompt    string
}

func NewService(store *storage.Store, cfg config.AppConfig) *Service {
	return &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

package runtime

import (
	"context"

	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/runtime/compression"
	"nano_cc/internal/web/storage"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	ListEnabledSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error)
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error
	GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
	CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error
	GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	ListSemanticMemories(ctx context.Context) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountSemanticMemories(ctx context.Context) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceSemanticMemories(ctx context.Context, items []storage.Memory) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
}

func NewService(store *storage.Store, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Hooks: NewDefaultHookManager(), EnableMemory: true}
}

func (s *Service) hookManager() *HookManager {
	if s.Hooks == nil {
		s.Hooks = NewDefaultHookManager()
	}
	return s.Hooks
}

func (s *Service) SetBuiltinSkills(loader *sessions.SkillLoader) {
	s.BuiltinSkills = loader
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

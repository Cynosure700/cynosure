package runtime

import (
	"context"
	"time"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/llm"
	"nano_cc/internal/sessions"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
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
	ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	CountProjectFactMemories(ctx context.Context, userID string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error)
	UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error
	AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error)
	RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error)
	ReleaseConversationLock(ctx context.Context, conversationID, token string) error
}

type Service struct {
	Store             conversationStore
	Cfg               config.AppConfig
	LLM               llm.Client
	Tools             *ToolRegistry
	BuiltinSkills     *sessions.SkillLoader
	BasePrompt        string
	CynosureMarkdown  config.CynosureMarkdownContext
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
	Approver          ApprovalDecider
}

func (s *Service) SetApprover(approver ApprovalDecider) {
	s.Approver = approver
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
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

func (s *Service) SetCynosureMarkdownContext(ctx config.CynosureMarkdownContext) {
	s.CynosureMarkdown = ctx
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

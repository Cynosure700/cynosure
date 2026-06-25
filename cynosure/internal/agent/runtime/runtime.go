package runtime

import (
	"context"
	"sync"
	"time"

	"cynosure/internal/agent/mcp"
	"cynosure/internal/agent/runtime/compression"
	"cynosure/internal/agent/storage"
	"cynosure/internal/config"
	"cynosure/internal/llm"
	"cynosure/internal/sessions"
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
	CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error
	GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error)
	GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error)
	ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error)
	ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error)
	InsertMemory(ctx context.Context, m storage.Memory) error
	CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error)
	DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error
	ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error
	LoadMemoryIndexForPrompt(ctx context.Context) (string, bool, int)
	ScanRecentMemories(ctx context.Context) ([]storage.ScannedMemory, error)
	ReadMemoryFile(ctx context.Context, path string) (storage.Memory, error)
	ShouldInjectMemory(ctx context.Context, conversationID, path string, modTime time.Time) (bool, error)
	ForgetInjectedMemory(ctx context.Context, path string) error
	UpdateMemoryFile(ctx context.Context, path string, update storage.MemoryUpdate) (string, error)
	DeleteMemoryFile(ctx context.Context, path string) error
	LoadConsolidationState(ctx context.Context) (storage.ConsolidationState, error)
	SaveConsolidationState(ctx context.Context, state storage.ConsolidationState) error
	ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error)
	ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error
	LoadConversationMemoryBreakpoint(ctx context.Context, conversationID string) (string, error)
	SaveConversationMemoryBreakpoint(ctx context.Context, conversationID, breakpointID string) error
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
	UserSkillsDir     string
	BasePrompt        string
	Prompts           FunctionalPrompts
	CynosureMarkdown  config.CynosureMarkdownContext
	GitStatus         string
	Hooks             *HookManager
	ContextCompressor *compression.Compressor
	EnableMemory      bool
	MCP               *mcp.Manager
	Approver          ApprovalDecider

	sessionMemoryMu       sync.Mutex
	sessionMemoryProgress map[string]*sessionMemoryProgress // key=会话 ID
}

// sessionMemoryProgress 跟踪某会话的会话记忆触发进度：记录上次更新时的上下文基线与
// 自上次以来的工具调用数。进程内态，重启后按"已有会话记忆是否存在"重建。
// 断点不在此保存——它持久化在会话记忆文件，压缩时从文件读取。
type sessionMemoryProgress struct {
	extracted          bool // 是否已完成初次提取（决定走 10K 门槛还是增量条件）
	baselineTokens     int  // 上次更新时的上下文 token（增长基线）
	toolCallsSinceBase int  // 自上次更新以来累计工具调用次数
	updating           bool // 是否有一次会话记忆更新正在进行（单航班守卫）
}

func (s *Service) SetApprover(approver ApprovalDecider) {
	s.Approver = approver
}

func NewService(store conversationStore, cfg config.AppConfig, client llm.Client) *Service {
	return &Service{Store: store, Cfg: cfg, LLM: client, Tools: NewToolRegistry(cfg), Prompts: defaultFunctionalPrompts(), Hooks: NewDefaultHookManager(), EnableMemory: true}
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

// SetUserSkillsDir 设置用户级 skill 目录（~/.cynosure/skills）。
// buildSkillSnapshot 每轮从该目录与工作区目录重新加载非内置 skill。
func (s *Service) SetUserSkillsDir(dir string) {
	s.UserSkillsDir = dir
}

func (s *Service) SetBasePrompt(prompt string) {
	s.BasePrompt = prompt
}

func (s *Service) SetCynosureMarkdownContext(ctx config.CynosureMarkdownContext) {
	s.CynosureMarkdown = ctx
}

func (s *Service) SetGitStatusContext(text string) {
	s.GitStatus = text
}

func (s *Service) SetMCPManager(manager *mcp.Manager) {
	s.MCP = manager
}

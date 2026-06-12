package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .link/.mcp.json 读取，在 Go 侧以强类型表示。
type MCPServer struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   bool              `json:"enabled"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	SessionID     string    `json:"session_id,omitempty"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ResumableSession struct {
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	WorkspaceRoot  string    `json:"workspace_root"`
	Title          string    `json:"title"`
	MessageCount   int       `json:"message_count"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// 该表已经不使用，作为存储的结构体，用于存储消息历史
type Message struct {
	ID               string            `json:"id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	Meta             *MessageMeta      `json:"meta,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

// MessageMeta 记录助手回复的元信息，仅对最终 assistant 消息填充。
type MessageMeta struct {
	ToolCallCount int `json:"tool_call_count"`          // 调用工具次数（0 也需序列化，保证历史展示一致）
	ContextTokens int `json:"context_tokens,omitempty"` // 当前上下文估算 token
	ContextBudget int `json:"context_budget,omitempty"` // 上下文预算（用于算占比）
}

type MessageToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"`
	Function MessageFunctionCall `json:"function"`
}

type MessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolCall struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	ToolName       string    `json:"tool_name"`
	Status         string    `json:"status"`
	Summary        string    `json:"summary"`
	CreatedAt      time.Time `json:"created_at"`
}

type PersistedOutput struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Content        string    `json:"content"`
	Preview        string    `json:"preview"`
	CreatedAt      time.Time `json:"created_at"`
}

type ToolResultLogEntry struct {
	ConversationID string    `json:"conversation_id"`
	SessionID      string    `json:"session_id,omitempty"`
	UserID         string    `json:"user_id"`
	ToolCallID     string    `json:"tool_call_id"`
	ToolName       string    `json:"tool_name"`
	RawArgs        string    `json:"raw_args"`
	Status         string    `json:"status"`
	Result         string    `json:"result"`
	AuditSummary   string    `json:"audit_summary,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}

type ContextSummary struct {
	ID                    string    `json:"id"`
	ConversationID        string    `json:"conversation_id"`
	UserID                string    `json:"user_id"`
	SourceHistorySHA256   string    `json:"source_history_sha256"`
	Strategy              string    `json:"strategy"`
	EstimatedTokensBefore int       `json:"estimated_tokens_before"`
	EstimatedTokensAfter  int       `json:"estimated_tokens_after"`
	Summary               string    `json:"summary"`
	CreatedAt             time.Time `json:"created_at"`
}

type SubagentMessage struct {
	ID               string            `json:"id"`
	RunID            string            `json:"run_id"`
	ParentToolCallID string            `json:"parent_tool_call_id"`
	ConversationID   string            `json:"conversation_id"`
	UserID           string            `json:"user_id"`
	SequenceNo       int               `json:"sequence_no"`
	Role             string            `json:"role"`
	Content          string            `json:"content"`
	ReasoningContent string            `json:"reasoning_content,omitempty"`
	ToolCallID       string            `json:"tool_call_id,omitempty"`
	ToolCalls        []MessageToolCall `json:"tool_calls,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ConversationMemory 是单个会话维度、随每轮对话增量维护的“当前会话主干信息”
// 条目。它不注入 system prompt，仅在上下文压缩触发全量摘要时作为替代品使用。
type ConversationMemory struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Body           string    `json:"body"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

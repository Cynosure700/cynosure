package storage

import "time"

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	MemoryEnabled bool      `json:"memory_enabled"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// MCPServer 是用户配置的一个 MCP 服务器连接。Transport 取值 stdio/sse/streamable。
// Args/Env/Headers 从本地 .cynosure/.mcp.json 读取，在 Go 侧以强类型表示。
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
	ID        string    `json:"id"`
	SessionID string    `json:"session_id,omitempty"`
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

// Message 是单条对话消息的存储结构，用于消息历史的编解码与展示。
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
	// EditLineStarts 仅用于 TUI 展示 edit_file/multi_edit 的 diff 真实行号：
	// [文件索引][该文件内改动索引] = new_string 在文件中的 1-based 起始行。
	// 在工具执行后（文件内容最新）计算并随展示历史持久化，使 /resume 在文件
	// 后续被改动甚至进程重启后仍能还原准确行号。它绝不进入发送给模型的消息
	// （buildOpenAIMessages 不拷贝该字段），只做持久化与展示用途。
	EditLineStarts [][]int   `json:"edit_line_starts,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
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

// ScannedMemory 是确定性扫描得到的候选记忆的轻量描述（只读 frontmatter）。
type ScannedMemory struct {
	Path        string
	Name        string
	Description string
	Type        string
	ModTime     time.Time
}

// MemoryUpdate 描述对单条记忆文件的部分更新；空字段表示保持原值。
type MemoryUpdate struct {
	Name        *string
	Description *string
	Body        *string
}

// ConsolidationState 记录上次定时去重的时间与累计会话计数，跨进程持久。
type ConsolidationState struct {
	LastRunAt    time.Time `json:"last_run_at"`
	SessionCount int       `json:"session_count"`
}

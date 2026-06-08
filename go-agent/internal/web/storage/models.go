package storage

import "time"

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AuthSession struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Status    string     `json:"status"`
	ExpiresAt time.Time  `json:"expires_at"`
	RevokedAt *time.Time `json:"revoked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Skill struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Status      string    `json:"status"`
	Source      string    `json:"source,omitempty"`
	ReadOnly    bool      `json:"readonly,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Conversation struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	RootMessageID string    `json:"root_message_id"`
	Title         string    `json:"title"`
	HistoryJSON   string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
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
	CreatedAt        time.Time         `json:"created_at"`
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

type UserProfile struct {
	UserID      string    `json:"user_id"`
	ProfileJSON string    `json:"profile_json"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ConversationTopics struct {
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	TopicsJSON     string    `json:"topics_json"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Memory struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"` // "" 表示系统级语义记忆
	Type        string    `json:"type"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

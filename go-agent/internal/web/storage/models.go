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

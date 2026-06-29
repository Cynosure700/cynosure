package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/storage"
	agenttools "cynosure/internal/tools"
)

type Store interface {
	UpdateConversationTitle(ctx context.Context, conversationID, title string) error
	TouchConversationActivity(ctx context.Context, conversationID string) error
	SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error
	CreateToolCall(ctx context.Context, tc storage.ToolCall) error
}

type EventWriter interface {
	Event(name string, data any) error
}

type UserPromptSubmitHook func(ctx context.Context, h *UserPromptSubmitContext) error
type PreToolUseHook func(ctx context.Context, h *ToolUseContext) error
type PostToolUseHook func(ctx context.Context, h *ToolUseContext) error
type StopHook func(ctx context.Context, h *StopContext) error

type HookManager struct {
	UserPromptSubmit []UserPromptSubmitHook
	PreToolUse       []PreToolUseHook
	PostToolUse      []PostToolUseHook
	Stop             []StopHook
}

type LoopState struct {
	Store                        Store
	NewMessageID                 func() string
	ToolRuntimeEnv               func() agenttools.RuntimeEnv
	ShouldInferConversationTitle func(string) bool
	InferConversationTitle       func(string) string

	Conversation storage.Conversation
	User         storage.User
	UserInput    string

	History     []storage.Message
	UserMessage storage.Message

	// ModelHistory 是实际发送给模型的那条唯一“真实消息历史”。它从上一回合
	// 持久化下来的压缩历史开始，在循环中与 History 同步追加，并在每一轮被压缩
	// 流水线的输出覆盖（因此内存态 == 发送态 == 落库态）。它是压缩与记忆提取
	// 的统一数据源。History 为了展示/持久化而保持完整且逐字；ModelHistory 可能
	// 携带压缩产物（摘要前言、占位符、persisted-output 标记、被裁剪的消息）。
	ModelHistory []storage.Message

	SkillSnapshot  *agenttools.SkillSnapshot
	SystemPrompt   string
	SystemReminder string
	Messages       []openai.ChatCompletionMessage
	Todos          []agenttools.TodoItem

	ToolCallCount     int // 累计工具调用次数
	LastContextTokens int // 最近一轮发送给模型的上下文 token
	LastContextBudget int // 上下文 token 预算

	Writer EventWriter
}

type UserPromptSubmitContext struct {
	State *LoopState
}

type ToolUseContext struct {
	State    *LoopState
	ToolCall openai.ToolCall
	Name     string
	RawArgs  string
	Outcome  ToolExecutionOutcome
}

type StopContext struct {
	State            *LoopState
	ModelMessage     openai.ChatCompletionMessage
	Content          string
	ReasoningContent string
	AssistantMessage storage.Message
}

type ToolExecutionOutcome struct {
	Status string                `json:"status"`
	Result string                `json:"result"`
	Audit  ToolExecutionAudit    `json:"-"`
	Todos  []agenttools.TodoItem `json:"-"`
	// EditLineStarts 仅用于 TUI 展示 edit_file/multi_edit 的 diff 真实行号，
	// 在工具执行后立即计算（此刻文件内容最新）。`json:"-"` 确保它不进入发送
	// 给模型的 tool 消息内容，只用于事件下发与展示历史持久化。
	EditLineStarts [][]int `json:"-"`
}

type ToolExecutionAudit struct {
	ResolvedCWD         string `json:"resolved_cwd,omitempty"`
	ResolvedCommandPath string `json:"resolved_command_path,omitempty"`
	OutcomeSummary      string `json:"outcome_summary,omitempty"`
	DenialReason        string `json:"denial_reason,omitempty"`
}

func (s *LoopState) NextMessageID() string {
	if s == nil || s.NewMessageID == nil {
		return ""
	}
	return s.NewMessageID()
}

func (s *LoopState) RuntimeEnv() agenttools.RuntimeEnv {
	if s == nil || s.ToolRuntimeEnv == nil {
		return agenttools.RuntimeEnv{}
	}
	return s.ToolRuntimeEnv()
}

func (s *LoopState) ShouldInferTitle(title string) bool {
	return s != nil && s.ShouldInferConversationTitle != nil && s.ShouldInferConversationTitle(title)
}

func (s *LoopState) InferTitle(input string) string {
	if s == nil || s.InferConversationTitle == nil {
		return input
	}
	return s.InferConversationTitle(input)
}

func (o ToolExecutionOutcome) MessageContent() string {
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf(`{"status":%q,"result":%q}`, o.Status, o.Result)
	}
	return string(data)
}

func (o ToolExecutionOutcome) AuditSummary() string {
	data, err := json.Marshal(o.Audit)
	if err != nil {
		return fmt.Sprintf(`{"resolved_cwd":%q,"resolved_command_path":%q,"outcome_summary":%q,"denial_reason":%q}`,
			o.Audit.ResolvedCWD,
			o.Audit.ResolvedCommandPath,
			o.Audit.OutcomeSummary,
			o.Audit.DenialReason,
		)
	}
	return string(data)
}

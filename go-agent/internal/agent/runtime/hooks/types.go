package hooks

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	agenttools "nano_cc/internal/tools"
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

	// ModelHistory is the "model line": the history actually fed to the
	// compression pipeline / model. It starts from the previous turn's
	// compressed request history (loaded from storage) and is appended in
	// lockstep with History during the loop. History stays full and verbatim
	// for display/persistence; ModelHistory may already carry compression
	// artifacts (summary preamble, placeholders, trimmed messages).
	ModelHistory []storage.Message

	SkillSnapshot *agenttools.SkillSnapshot
	SystemPrompt  string
	Messages      []openai.ChatCompletionMessage
	Todos         []agenttools.TodoItem

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
}

type ToolExecutionAudit struct {
	ResolvedCWD           string `json:"resolved_cwd,omitempty"`
	ResolvedCommandPath   string `json:"resolved_command_path,omitempty"`
	CommandArtifactPath   string `json:"command_artifact_path,omitempty"`
	CommandArtifactSource string `json:"command_artifact_source,omitempty"`
	OutcomeSummary        string `json:"outcome_summary,omitempty"`
	DenialReason          string `json:"denial_reason,omitempty"`
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
		return fmt.Sprintf(`{"resolved_cwd":%q,"resolved_command_path":%q,"command_artifact_path":%q,"command_artifact_source":%q,"outcome_summary":%q,"denial_reason":%q}`,
			o.Audit.ResolvedCWD,
			o.Audit.ResolvedCommandPath,
			o.Audit.CommandArtifactPath,
			o.Audit.CommandArtifactSource,
			o.Audit.OutcomeSummary,
			o.Audit.DenialReason,
		)
	}
	return string(data)
}

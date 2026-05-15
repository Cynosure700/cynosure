package runtime

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/assistant"
	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/storage"
)

type EventWriter interface {
	Event(name string, data any) error
}

type conversationStore interface {
	CreateMessage(ctx context.Context, message storage.Message) error
	TouchConversation(ctx context.Context, conversationID, title string) error
	ListEnabledSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error)
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
}

type toolExecutionOutcome struct {
	Status string             `json:"status"`
	Result string             `json:"result"`
	Audit  toolExecutionAudit `json:"-"`
}

type toolExecutionAudit struct {
	ResolvedCWD         string `json:"resolved_cwd,omitempty"`
	CommandArtifactPath string `json:"command_artifact_path,omitempty"`
	OutcomeSummary      string `json:"outcome_summary,omitempty"`
	DenialReason        string `json:"denial_reason,omitempty"`
}

func NewService(store *storage.Store, cfg config.AppConfig, builtinSkills *sessions.SkillLoader) *Service {
	return &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(store, cfg), BuiltinSkills: builtinSkills}
}

func (s *Service) RespondToConversation(ctx context.Context, conversation storage.Conversation, user storage.User, userMessage string, writer EventWriter) (storage.Message, error) {
	history, err := s.loadConversationMessages(ctx, conversation.ID)
	if err != nil {
		return storage.Message{}, err
	}
	if err := s.Store.CreateMessage(ctx, storage.Message{ID: newMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "user", Content: userMessage}); err != nil {
		return storage.Message{}, err
	}
	history = append(history, storage.Message{ConversationID: conversation.ID, UserID: user.ID, Role: "user", Content: userMessage})
	if err := s.Store.TouchConversation(ctx, conversation.ID, inferConversationTitle(conversation.Title, userMessage)); err != nil {
		return storage.Message{}, err
	}
	workspaceRoot, err := s.resolveUserWorkspace(user.ID)
	if err != nil {
		return storage.Message{}, err
	}

	skills, err := s.Store.ListEnabledSkillsByUser(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	loader := s.buildConversationSkillLoader(skills)
	systemPrompt := s.buildSystemPrompt(user, loader)
	messages := buildOpenAIMessages(systemPrompt, history)
	round := 0

	for {
		round++
		req := openai.ChatCompletionRequest{
			Model:    s.Cfg.LLM.ModelID,
			Messages: messages,
			Tools:    s.Tools.Definitions(loader),
		}
		reqBody, _ := json.Marshal(req)
		resp, err := config.Client.CreateChatCompletion(ctx, req)
		respBody, _ := json.Marshal(resp)
		logger.LogLLMRound(round, fmt.Sprintf("web-runtime conversation=%s", conversation.ID), reqBody, respBody, err)
		if err != nil {
			return storage.Message{}, err
		}
		if len(resp.Choices) == 0 {
			return storage.Message{}, fmt.Errorf("model returned no choices")
		}
		choice := resp.Choices[0]
		msg := choice.Message
		messages = append(messages, msg)

		if choice.FinishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			return s.persistAssistantReply(ctx, conversation, user.ID, history, fallbackAssistantContent(msg.Content), writer)
		}

		for _, tc := range msg.ToolCalls {
			outcome := s.executeToolCall(ctx, ToolContext{User: user, Conversation: conversation, Loader: loader, WorkspaceRoot: workspaceRoot}, tc.Function.Name, tc.Function.Arguments)
			_ = s.Store.CreateToolCall(ctx, storage.ToolCall{ID: newToolCallID(), ConversationID: conversation.ID, UserID: user.ID, ToolName: tc.Function.Name, Status: outcome.Status, Summary: outcome.AuditSummary()})
			if writer != nil {
				_ = writer.Event("tool", map[string]any{"name": tc.Function.Name, "status": outcome.Status, "result": outcome.Result})
			}
			messages = append(messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: tc.ID, Content: outcome.MessageContent()})
		}
	}
}

func (s *Service) resolveUserWorkspace(userID string) (string, error) {
	_ = userID
	base := strings.TrimSpace(s.Cfg.WorkspaceRoot)
	if base == "" {
		return "", fmt.Errorf("workspace root is not configured")
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", fmt.Errorf("create workspace root: %w", err)
	}
	resolvedBase, err := filepath.Abs(base)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	return filepath.Clean(resolvedBase), nil
}

func (s *Service) executeToolCall(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) toolExecutionOutcome {
	audit := toolExecutionAudit{
		ResolvedCWD:         strings.TrimSpace(toolCtx.WorkspaceRoot),
		CommandArtifactPath: resolveCommandArtifactPath(name, rawArgs, s.Cfg.CommandBinDir, s.Cfg.CommandScriptDir),
	}
	result, err := s.Tools.Execute(ctx, toolCtx, name, rawArgs)
	if err != nil {
		audit.DenialReason = err.Error()
		return toolExecutionOutcome{Status: "rejected", Result: fmt.Sprintf("Error: %v", err), Audit: audit}
	}
	audit.OutcomeSummary = truncate(result, 500)
	return toolExecutionOutcome{Status: "success", Result: result, Audit: audit}
}

func (o toolExecutionOutcome) MessageContent() string {
	data, err := json.Marshal(o)
	if err != nil {
		return fmt.Sprintf(`{"status":%q,"result":%q}`, o.Status, o.Result)
	}
	return string(data)
}

func (o toolExecutionOutcome) AuditSummary() string {
	data, err := json.Marshal(o.Audit)
	if err != nil {
		return fmt.Sprintf(`{"resolved_cwd":%q,"command_artifact_path":%q,"outcome_summary":%q,"denial_reason":%q}`,
			o.Audit.ResolvedCWD,
			o.Audit.CommandArtifactPath,
			o.Audit.OutcomeSummary,
			o.Audit.DenialReason,
		)
	}
	return string(data)
}

func resolveCommandArtifactPath(toolName, rawArgs string, roots ...string) string {
	if toolName != "bash" {
		return ""
	}
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return ""
	}
	command, _ := args["command"].(string)
	for _, token := range strings.Fields(command) {
		candidate := strings.Trim(token, "\"'`;,()[]{}")
		if candidate == "" || !filepath.IsAbs(candidate) {
			continue
		}
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			continue
		}
		cleanResolved := filepath.Clean(resolved)
		for _, root := range roots {
			if root == "" {
				continue
			}
			resolvedRoot, err := filepath.Abs(root)
			if err != nil {
				continue
			}
			cleanRoot := filepath.Clean(resolvedRoot)
			if cleanResolved == cleanRoot || strings.HasPrefix(cleanResolved, cleanRoot+string(os.PathSeparator)) {
				return cleanResolved
			}
		}
	}
	return ""
}

func (s *Service) persistAssistantReply(ctx context.Context, conversation storage.Conversation, userID string, history []storage.Message, content string, writer EventWriter) (storage.Message, error) {
	assistant := storage.Message{ID: newMessageID(), ConversationID: conversation.ID, UserID: userID, Role: "assistant", Content: content}
	if err := s.Store.CreateMessage(ctx, assistant); err != nil {
		return storage.Message{}, err
	}
	updatedHistory := append(history, assistant)
	if err := s.Store.SetConversationCache(ctx, conversation.ID, updatedHistory); err != nil {
		_ = err
	}
	if writer != nil {
		_ = writer.Event("assistant", map[string]any{"content": assistant.Content})
	}
	return assistant, nil
}

func (s *Service) loadConversationMessages(ctx context.Context, conversationID string) ([]storage.Message, error) {
	if cached, ok, err := s.Store.GetConversationCache(ctx, conversationID); err == nil && ok {
		return cached, nil
	}
	messages, err := s.Store.ListMessagesByConversation(ctx, conversationID, 100)
	if err != nil {
		return nil, err
	}
	if err := s.Store.SetConversationCache(ctx, conversationID, messages); err != nil {
		_ = err
	}
	return messages, nil
}

func (s *Service) buildSystemPrompt(user storage.User, loader *sessions.SkillLoader) string {
	return assistant.BuildSystemPrompt(assistant.PromptOptions{
		Surface:           fmt.Sprintf("the browser chat experience for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
	})
}

func (s *Service) buildConversationSkillLoader(skills []storage.Skill) *sessions.SkillLoader {
	return sessions.MergeSkillLoaders(s.BuiltinSkills, buildDBSkillLoader(skills))
}

func buildDBSkillLoader(skills []storage.Skill) *sessions.SkillLoader {
	entries := make(map[string]*sessions.SkillEntry, len(skills))
	for _, skill := range skills {
		entries[skill.Slug] = &sessions.SkillEntry{
			Meta: map[string]string{"description": skill.Description},
			Body: skill.Content,
			Path: "db://skills/" + skill.ID,
		}
	}
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(entries)
	return loader
}

func buildOpenAIMessages(systemPrompt string, history []storage.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content})
	}
	return messages
}

func fallbackAssistantContent(content string) string {
	if strings.TrimSpace(content) == "" {
		return "(no response)"
	}
	return content
}

func inferConversationTitle(currentTitle, userMessage string) string {
	if strings.TrimSpace(currentTitle) != "" && currentTitle != "新对话" {
		return currentTitle
	}
	trimmed := strings.TrimSpace(userMessage)
	if len([]rune(trimmed)) > 30 {
		return string([]rune(trimmed)[:30])
	}
	if trimmed == "" {
		return "新对话"
	}
	return trimmed
}

func truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max]
}

func newMessageID() string  { return newID("msg") }
func newToolCallID() string { return newID("tool") }

func newID(prefix string) string {
	buf := make([]byte, 12)
	_, _ = rand.Read(buf)
	return fmt.Sprintf("%s_%x", prefix, buf)
}

type SSEWriter struct {
	W http.ResponseWriter
}

func (s SSEWriter) Event(name string, data any) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(s.W, "event: %s\ndata: %s\n\n", name, string(bytes)); err != nil {
		return err
	}
	if flusher, ok := s.W.(http.Flusher); ok {
		flusher.Flush()
	}
	return nil
}

func SkillNames(skills []storage.Skill) []string {
	names := make([]string, 0, len(skills))
	for _, skill := range skills {
		names = append(names, skill.Slug)
	}
	sort.Strings(names)
	return names
}

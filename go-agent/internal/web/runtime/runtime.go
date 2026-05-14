package runtime

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/storage"
)

type EventWriter interface {
	Event(name string, data any) error
}

type Service struct {
	Store *storage.Store
	Cfg   config.AppConfig
	Tools *ToolRegistry
}

func NewService(store *storage.Store, cfg config.AppConfig) *Service {
	return &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(store, cfg)}
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

	skills, err := s.Store.ListEnabledSkillsByUser(ctx, user.ID)
	if err != nil {
		return storage.Message{}, err
	}
	loader := buildDBSkillLoader(skills)
	systemPrompt := s.buildSystemPrompt(user, loader)
	messages := buildOpenAIMessages(systemPrompt, history)

	for {
		req := openai.ChatCompletionRequest{
			Model:    s.Cfg.LLM.ModelID,
			Messages: messages,
			Tools:    s.Tools.Definitions(loader),
		}
		resp, err := config.Client.CreateChatCompletion(ctx, req)
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
			assistant := storage.Message{ID: newMessageID(), ConversationID: conversation.ID, UserID: user.ID, Role: "assistant", Content: fallbackAssistantContent(msg.Content)}
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

		for _, tc := range msg.ToolCalls {
			result, err := s.Tools.Execute(ctx, ToolContext{User: user, Conversation: conversation, Loader: loader}, tc.Function.Name, tc.Function.Arguments)
			status := "success"
			if err != nil {
				status = "rejected"
				result = fmt.Sprintf("Error: %v", err)
			}
			_ = s.Store.CreateToolCall(ctx, storage.ToolCall{ID: newToolCallID(), ConversationID: conversation.ID, UserID: user.ID, ToolName: tc.Function.Name, Status: status, Summary: truncate(result, 500)})
			if writer != nil {
				_ = writer.Event("tool", map[string]any{"name": tc.Function.Name, "status": status, "result": result})
			}
			messages = append(messages, openai.ChatCompletionMessage{Role: "tool", ToolCallID: tc.ID, Content: result})
		}
	}
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
	workspace := filepath.Join(s.Cfg.WorkspaceRoot, user.ID)
	base := fmt.Sprintf("You are a web coding agent for user %s. Only use registered tools. All file operations must stay inside %s.", user.Username, workspace)
	if descriptions := loader.GetDescriptions(); descriptions != "" {
		base += "\nAvailable user skills:\n" + descriptions
	}
	return base
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
	loader := &sessions.SkillLoader{Skills: make(map[string]*sessions.SkillEntry)}
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

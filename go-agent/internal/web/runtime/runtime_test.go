package runtime

import (
	"context"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

type fakeStore struct {
	messages      []storage.Message
	toolCalls     []storage.ToolCall
	cached        []storage.Message
	enabledSkills []storage.Skill
	conversation  storage.Conversation
}

func (f *fakeStore) CreateMessage(ctx context.Context, message storage.Message) error {
	f.messages = append(f.messages, message)
	return nil
}

func (f *fakeStore) TouchConversation(ctx context.Context, conversationID, title string) error {
	f.conversation.ID = conversationID
	f.conversation.Title = title
	return nil
}

func (f *fakeStore) ListEnabledSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error) {
	return f.enabledSkills, nil
}

func (f *fakeStore) SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error {
	f.cached = append([]storage.Message(nil), messages...)
	return nil
}

func (f *fakeStore) GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error) {
	if len(f.cached) == 0 {
		return nil, false, nil
	}
	return append([]storage.Message(nil), f.cached...), true, nil
}

func (f *fakeStore) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error) {
	return nil, nil
}

func (f *fakeStore) CreateToolCall(ctx context.Context, tc storage.ToolCall) error {
	f.toolCalls = append(f.toolCalls, tc)
	return nil
}

type fakeLLMClient struct {
	responses []openai.ChatCompletionResponse
	calls     int
	lastReq   openai.ChatCompletionRequest
}

func (f *fakeLLMClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	f.calls++
	f.lastReq = req
	if len(f.responses) == 0 {
		return openai.ChatCompletionResponse{}, nil
	}
	resp := f.responses[0]
	f.responses = f.responses[1:]
	return resp, nil
}

func TestRespondToConversation_DirectAnswerWithoutTools(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "当然可以，我来帮你规划。"},
		}},
	}}}
	config.Client = llm

	store := &fakeStore{}
	service := &Service{Store: store, Cfg: config.AppConfig{LLM: config.Config{ModelID: "test-model"}}, Tools: NewToolRegistry(nil, config.AppConfig{}), BuiltinSkills: nil}
	conversation := storage.Conversation{ID: "conv_1", Title: "新对话"}
	user := storage.User{ID: "usr_1", Username: "alice"}

	message, err := service.RespondToConversation(context.Background(), conversation, user, "帮我做一个今天的学习计划", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "当然可以，我来帮你规划。" {
		t.Fatalf("unexpected assistant content: %q", message.Content)
	}
	if llm.calls != 1 {
		t.Fatalf("expected llm to be called once, got %d", llm.calls)
	}
	if len(llm.lastReq.Tools) != 0 {
		t.Fatalf("expected no tools for direct answer, got %d", len(llm.lastReq.Tools))
	}
	if len(store.toolCalls) != 0 {
		t.Fatalf("expected no tool calls to be stored, got %d", len(store.toolCalls))
	}
	if len(store.messages) != 2 {
		t.Fatalf("expected user and assistant messages to be persisted, got %d", len(store.messages))
	}
	if len(store.cached) != 2 {
		t.Fatalf("expected cached conversation with 2 messages, got %d", len(store.cached))
	}
}

func TestRespondToConversation_BrowserCapabilityBoundarySkipsModel(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{}
	config.Client = llm

	store := &fakeStore{}
	service := &Service{Store: store, Cfg: config.AppConfig{LLM: config.Config{ModelID: "test-model"}}, Tools: NewToolRegistry(nil, config.AppConfig{}), BuiltinSkills: nil}
	conversation := storage.Conversation{ID: "conv_2", Title: "新对话", UpdatedAt: time.Now()}
	user := storage.User{ID: "usr_2", Username: "bob"}

	message, err := service.RespondToConversation(context.Background(), conversation, user, "请帮我执行 shell 命令 ls 查看本地目录", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if llm.calls != 0 {
		t.Fatalf("expected llm not to be called for browser boundary request, got %d", llm.calls)
	}
	if len(store.toolCalls) != 0 {
		t.Fatalf("expected no tool calls for browser boundary request, got %d", len(store.toolCalls))
	}
	if got := message.Content; got == "" || !contains(got, "不能访问你的本地 shell") {
		t.Fatalf("expected browser capability explanation, got %q", got)
	}
	if len(store.messages) != 2 {
		t.Fatalf("expected user and assistant messages to be persisted, got %d", len(store.messages))
	}
}

func TestRespondToConversation_MergesBuiltinAndUserSkillsIntoPrompt(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "我会结合这些能力来回答。"},
		}},
	}}}
	config.Client = llm

	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {
			Meta: map[string]string{"description": "Builtin description"},
			Body: "builtin body",
			Path: "builtin://builtin-skill",
		},
	})

	store := &fakeStore{enabledSkills: []storage.Skill{{
		ID:          "skill_1",
		UserID:      "usr_3",
		Name:        "User Skill",
		Slug:        "user-skill",
		Description: "User description",
		Content:     "user body",
		Status:      "enabled",
	}}}
	service := &Service{Store: store, Cfg: config.AppConfig{LLM: config.Config{ModelID: "test-model"}}, Tools: NewToolRegistry(nil, config.AppConfig{}), BuiltinSkills: builtin}
	conversation := storage.Conversation{ID: "conv_3", Title: "新对话"}
	user := storage.User{ID: "usr_3", Username: "carol"}

	_, err := service.RespondToConversation(context.Background(), conversation, user, "请结合现有技能回答", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if llm.calls != 1 {
		t.Fatalf("expected llm to be called once, got %d", llm.calls)
	}
	if len(llm.lastReq.Tools) != 1 {
		t.Fatalf("expected load_skill tool to be exposed, got %d tools", len(llm.lastReq.Tools))
	}
	systemPrompt := llm.lastReq.Messages[0].Content
	if !contains(systemPrompt, "- builtin-skill: Builtin description") {
		t.Fatalf("expected builtin skill description in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "- user-skill: User description") {
		t.Fatalf("expected user skill description in prompt, got %q", systemPrompt)
	}
}

func TestToolRegistryDefinitions_UsesRegisteredToolDefinition(t *testing.T) {
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	registry := NewToolRegistry(nil, config.AppConfig{})

	defs := registry.Definitions(loader)
	if len(defs) != 1 {
		t.Fatalf("expected one web tool definition, got %d", len(defs))
	}
	expected, ok := lookupRegisteredTool("load_skill")
	if !ok {
		t.Fatalf("expected load_skill to exist in registered tool definitions")
	}
	if defs[0].Function == nil || expected.Function == nil {
		t.Fatalf("expected function definitions to be present")
	}
	if defs[0].Function.Name != expected.Function.Name {
		t.Fatalf("expected tool name %q, got %q", expected.Function.Name, defs[0].Function.Name)
	}
	if defs[0].Function.Description != expected.Function.Description {
		t.Fatalf("expected tool description %q, got %q", expected.Function.Description, defs[0].Function.Description)
	}
}

func TestToolRegistryExecute_LoadSkillReturnsSkillContent(t *testing.T) {
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	registry := NewToolRegistry(nil, config.AppConfig{})

	content, err := registry.Execute(context.Background(), ToolContext{Loader: loader}, "load_skill", `{"name":"builtin-skill"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(content, "<skill name=\"builtin-skill\">") || !contains(content, "builtin body") {
		t.Fatalf("expected loaded skill content, got %q", content)
	}
	if _, ok := agenttools.Handlers["load_skill"]; !ok {
		t.Fatalf("expected load_skill handler to remain registered in shared tool registry")
	}
}

func TestRegisteredTools_UsesCurrentWebAllowList(t *testing.T) {
	tools := RegisteredTools()
	if len(tools) != 1 || tools[0] != "load_skill" {
		t.Fatalf("expected only load_skill to be registered for web runtime, got %v", tools)
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

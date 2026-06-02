package runtime

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

type fakeStore struct {
	messages         []storage.Message
	historyUpdates   [][]storage.Message
	toolCalls        []storage.ToolCall
	cached           []storage.Message
	enabledSkills    []storage.Skill
	enabledSkillSets [][]storage.Skill
	listEnabledCalls int
	updatedTitle     string
	updatedID        string
	touchedID        string
}

func (f *fakeStore) CreateMessage(ctx context.Context, message storage.Message) error {
	f.messages = append(f.messages, message)
	return nil
}

func (f *fakeStore) SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error {
	f.historyUpdates = append(f.historyUpdates, append([]storage.Message(nil), messages...))
	f.cached = append([]storage.Message(nil), messages...)
	return nil
}

func (f *fakeStore) UpdateConversationTitle(ctx context.Context, conversationID, title string) error {
	f.updatedID = conversationID
	f.updatedTitle = title
	return nil
}

func (f *fakeStore) TouchConversationActivity(ctx context.Context, conversationID string) error {
	f.touchedID = conversationID
	return nil
}

func (f *fakeStore) ListEnabledSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error) {
	f.listEnabledCalls++
	if len(f.enabledSkillSets) > 0 {
		index := f.listEnabledCalls - 1
		if index >= len(f.enabledSkillSets) {
			index = len(f.enabledSkillSets) - 1
		}
		return f.enabledSkillSets[index], nil
	}
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
	reqs      []openai.ChatCompletionRequest
}

type captureEventWriter struct {
	events []capturedEvent
}

type capturedEvent struct {
	name string
	data any
}

func (w *captureEventWriter) Event(name string, data any) error {
	w.events = append(w.events, capturedEvent{name: name, data: data})
	return nil
}

func testAppConfig(t *testing.T) config.AppConfig {
	t.Helper()
	return config.AppConfig{
		LLM:           config.Config{ModelID: "test-model"},
		WorkspaceRoot: t.TempDir(),
	}
}

func (f *fakeLLMClient) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	f.calls++
	f.lastReq = req
	f.reqs = append(f.reqs, req)
	if len(f.responses) == 0 {
		return openai.ChatCompletionResponse{}, nil
	}
	resp := f.responses[0]
	f.responses = f.responses[1:]
	return resp, nil
}

func TestHookManagerRunsUserPromptSubmitHooksInOrder(t *testing.T) {
	var order []string
	manager := &HookManager{UserPromptSubmit: []UserPromptSubmitHook{
		func(ctx context.Context, h *UserPromptSubmitContext) error {
			order = append(order, "first")
			return nil
		},
		func(ctx context.Context, h *UserPromptSubmitContext) error {
			order = append(order, "second")
			return nil
		},
	}}

	if err := manager.RunUserPromptSubmit(context.Background(), &UserPromptSubmitContext{State: &LoopState{}}); err != nil {
		t.Fatalf("run user prompt hooks: %v", err)
	}
	if strings.Join(order, ",") != "first,second" {
		t.Fatalf("expected hooks to run in order, got %v", order)
	}
}

func TestHookManagerStopsOnUserPromptSubmitError(t *testing.T) {
	var order []string
	manager := &HookManager{UserPromptSubmit: []UserPromptSubmitHook{
		func(ctx context.Context, h *UserPromptSubmitContext) error {
			order = append(order, "first")
			return context.Canceled
		},
		func(ctx context.Context, h *UserPromptSubmitContext) error {
			order = append(order, "second")
			return nil
		},
	}}

	if err := manager.RunUserPromptSubmit(context.Background(), &UserPromptSubmitContext{State: &LoopState{}}); err == nil {
		t.Fatalf("expected hook error")
	}
	if strings.Join(order, ",") != "first" {
		t.Fatalf("expected hooks to stop on error, got %v", order)
	}
}

func TestDefaultHooksPersistAssistantOnStop(t *testing.T) {
	store := &fakeStore{}
	writer := &captureEventWriter{}
	state := &LoopState{
		Store:        store,
		NewMessageID: newMessageID,
		Conversation: storage.Conversation{ID: "conv_stop"},
		User:         storage.User{ID: "usr_stop"},
		History:      []storage.Message{{ID: "msg_user", ConversationID: "conv_stop", UserID: "usr_stop", Role: "user", Content: "hello"}},
		Writer:       writer,
	}
	stopCtx := &StopContext{State: state, Content: "answer", ReasoningContent: "thinking"}

	if err := NewDefaultHookManager().RunStop(context.Background(), stopCtx); err != nil {
		t.Fatalf("run stop hooks: %v", err)
	}
	if stopCtx.AssistantMessage.Content != "answer" || stopCtx.AssistantMessage.ReasoningContent != "thinking" {
		t.Fatalf("expected assistant message to be filled, got %#v", stopCtx.AssistantMessage)
	}
	if len(store.historyUpdates) != 1 || len(store.historyUpdates[0]) != 2 || store.historyUpdates[0][1].ReasoningContent != "thinking" {
		t.Fatalf("expected stop hook to persist full history, got %#v", store.historyUpdates)
	}
	if len(writer.events) != 1 || writer.events[0].name != "assistant" {
		t.Fatalf("expected assistant event, got %#v", writer.events)
	}
}

func TestDefaultToolHooksAuditPersistEmitAndAppendToolMessage(t *testing.T) {
	store := &fakeStore{}
	writer := &captureEventWriter{}
	workspace := t.TempDir()
	tools := NewToolRegistry(config.AppConfig{WorkspaceRoot: workspace})
	state := &LoopState{
		Store:          store,
		NewMessageID:   newMessageID,
		ToolRuntimeEnv: tools.runtimeEnv,
		Conversation:   storage.Conversation{ID: "conv_tool"},
		User:           storage.User{ID: "usr_tool"},
		Messages:       []openai.ChatCompletionMessage{{Role: "system", Content: "system"}},
		Writer:         writer,
	}
	toolCtx := &ToolUseContext{
		State:    state,
		ToolCall: openai.ToolCall{ID: "tool_1", Function: openai.FunctionCall{Name: "bash", Arguments: `{"command":"pwd"}`}},
		Name:     "bash",
		RawArgs:  `{"command":"pwd"}`,
		Outcome:  toolExecutionOutcome{Status: "success", Result: workspace},
	}
	manager := NewDefaultHookManager()

	if err := manager.RunPreToolUse(context.Background(), toolCtx); err != nil {
		t.Fatalf("run pre tool hooks: %v", err)
	}
	if err := manager.RunPostToolUse(context.Background(), toolCtx); err != nil {
		t.Fatalf("run post tool hooks: %v", err)
	}
	if toolCtx.Outcome.Audit.ResolvedCWD != workspace || toolCtx.Outcome.Audit.OutcomeSummary == "" {
		t.Fatalf("expected audit to be filled, got %#v", toolCtx.Outcome.Audit)
	}
	if len(store.toolCalls) != 1 || store.toolCalls[0].Status != "success" {
		t.Fatalf("expected persisted tool call, got %#v", store.toolCalls)
	}
	if len(writer.events) != 1 || writer.events[0].name != "tool" {
		t.Fatalf("expected tool event, got %#v", writer.events)
	}
	if len(state.Messages) != 2 || state.Messages[1].Role != "tool" || state.Messages[1].ToolCallID != "tool_1" {
		t.Fatalf("expected tool message to be appended, got %#v", state.Messages)
	}
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
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
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
	if len(llm.lastReq.Tools) != 1 || llm.lastReq.Tools[0].Function == nil || llm.lastReq.Tools[0].Function.Name != "load_skill" {
		t.Fatalf("expected startup-loaded default load_skill tool, got %#v", llm.lastReq.Tools)
	}
	if len(store.toolCalls) != 0 {
		t.Fatalf("expected no tool calls to be stored, got %d", len(store.toolCalls))
	}
	if len(store.messages) != 0 {
		t.Fatalf("expected no legacy message row writes, got %d", len(store.messages))
	}
	if len(store.historyUpdates) != 1 {
		t.Fatalf("expected one full history update, got %d", len(store.historyUpdates))
	}
	if got := store.historyUpdates[0]; len(got) != 2 || got[0].Role != "user" || got[1].Role != "assistant" {
		t.Fatalf("expected full user+assistant history update, got %#v", got)
	}
	if len(store.cached) != 2 {
		t.Fatalf("expected cached conversation with 2 messages, got %d", len(store.cached))
	}
	if store.updatedID != "conv_1" || store.updatedTitle != "帮我做一个今天的学习计划" {
		t.Fatalf("expected default-title conversation to infer title, got id=%q title=%q", store.updatedID, store.updatedTitle)
	}
	if store.touchedID != "" {
		t.Fatalf("expected no activity-only touch when inferring title, got %q", store.touchedID)
	}
}

func TestRespondToConversation_PersistsReasoningContent(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "最终答案", ReasoningContent: "内部推理过程"},
		}},
	}}}
	config.Client = llm

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	writer := &captureEventWriter{}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_reasoning", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "解释一下", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "最终答案" || message.ReasoningContent != "内部推理过程" {
		t.Fatalf("expected content and reasoning content to be persisted, got %#v", message)
	}
	if got := store.historyUpdates[0][1].ReasoningContent; got != "内部推理过程" {
		t.Fatalf("expected history reasoning content, got %q", got)
	}
	if len(writer.events) != 1 || writer.events[0].name != "assistant" {
		t.Fatalf("expected one assistant event, got %#v", writer.events)
	}
	payload, ok := writer.events[0].data.(map[string]any)
	if !ok || payload["reasoning_content"] != "内部推理过程" {
		t.Fatalf("expected assistant event to include reasoning_content, got %#v", writer.events[0].data)
	}
}

func TestRespondToConversation_ShellRequestsReachModelWithWorkspaceTools(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "我会在 workspace 中执行命令。"},
		}},
	}}}
	config.Client = llm

	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.WebAllowedTools = []string{"bash"}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	conversation := storage.Conversation{ID: "conv_2", Title: "新对话", UpdatedAt: time.Now()}
	user := storage.User{ID: "usr_2", Username: "bob"}

	message, err := service.RespondToConversation(context.Background(), conversation, user, "请帮我执行 shell 命令 ls 查看本地目录", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if llm.calls != 1 {
		t.Fatalf("expected llm to be called for shell request, got %d", llm.calls)
	}
	if len(store.toolCalls) != 0 {
		t.Fatalf("expected no tool execution for direct model response, got %d", len(store.toolCalls))
	}
	if got := message.Content; got != "我会在 workspace 中执行命令。" {
		t.Fatalf("unexpected assistant content: %q", got)
	}
	if len(llm.lastReq.Tools) != 1 || llm.lastReq.Tools[0].Function == nil || llm.lastReq.Tools[0].Function.Name != "bash" {
		t.Fatalf("expected shell request to expose bash tool, got %#v", llm.lastReq.Tools)
	}
	systemPrompt := llm.lastReq.Messages[0].Content
	if !contains(systemPrompt, "Current workspace root: "+cfg.WorkspaceRoot) {
		t.Fatalf("expected system prompt to include workspace root, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "rather than a chat-only assistant") {
		t.Fatalf("expected system prompt to position the model as a full agent, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "Runtime tools available in this conversation: bash.") {
		t.Fatalf("expected system prompt to include available tools, got %q", systemPrompt)
	}
	if len(store.historyUpdates) != 1 {
		t.Fatalf("expected one full history update, got %d", len(store.historyUpdates))
	}
	if got := store.historyUpdates[0]; len(got) != 2 || got[0].Content != "请帮我执行 shell 命令 ls 查看本地目录" || got[1].Content != "我会在 workspace 中执行命令。" {
		t.Fatalf("expected full shell conversation history, got %#v", got)
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
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
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
	if !contains(systemPrompt, "Current workspace root: "+cfg.WorkspaceRoot) {
		t.Fatalf("expected workspace root in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "Runtime tools available in this conversation: load_skill.") {
		t.Fatalf("expected load_skill to be listed in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "- builtin-skill: Builtin description") {
		t.Fatalf("expected builtin skill description in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "- user-skill: User description") {
		t.Fatalf("expected user skill description in prompt, got %q", systemPrompt)
	}
}

func TestBuildOpenAIMessagesCarriesHistoryReasoningContent(t *testing.T) {
	messages := buildOpenAIMessages("system prompt", []storage.Message{
		{Role: "user", Content: "问题"},
		{Role: "assistant", Content: "回答", ReasoningContent: "历史推理"},
	})

	if len(messages) != 3 {
		t.Fatalf("expected system plus 2 history messages, got %d", len(messages))
	}
	if messages[2].Content != "回答" || messages[2].ReasoningContent != "历史推理" {
		t.Fatalf("expected assistant history message to carry reasoning content, got %#v", messages[2])
	}
}

func TestBuildOpenAIMessagesCarriesHistoricalToolResult(t *testing.T) {
	messages := buildOpenAIMessages("system prompt", []storage.Message{
		{Role: "assistant", ToolCalls: []storage.MessageToolCall{{ID: "tool_1", Type: "function", Function: storage.MessageFunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}}}},
		{Role: "tool", ToolCallID: "tool_1", Content: `{"status":"success","result":"loaded"}`},
	})

	if len(messages) != 3 {
		t.Fatalf("expected system plus 2 history messages, got %d", len(messages))
	}
	if len(messages[1].ToolCalls) != 1 || messages[1].ToolCalls[0].ID != "tool_1" || messages[1].ToolCalls[0].Function.Name != "load_skill" {
		t.Fatalf("expected assistant history message to carry tool call metadata, got %#v", messages[1])
	}
	if messages[2].Role != "tool" || messages[2].ToolCallID != "tool_1" || messages[2].Content == "" {
		t.Fatalf("expected historical tool result message to carry tool_call_id and content, got %#v", messages[2])
	}
}

func TestRespondToConversation_LoadsHistoricalToolResultBeforeAgentLoop(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "继续回答"},
		}},
	}}}
	config.Client = llm

	store := &fakeStore{cached: []storage.Message{
		{ID: "msg_old_user", ConversationID: "conv_history", UserID: "usr_history", Role: "user", Content: "先加载技能"},
		{ID: "msg_old_assistant_tool", ConversationID: "conv_history", UserID: "usr_history", Role: "assistant", ToolCalls: []storage.MessageToolCall{{ID: "tool_1", Type: "function", Function: storage.MessageFunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}}}},
		{ID: "msg_old_tool", ConversationID: "conv_history", UserID: "usr_history", Role: "tool", ToolCallID: "tool_1", Content: `{"status":"success","result":"loaded"}`},
		{ID: "msg_old_assistant", ConversationID: "conv_history", UserID: "usr_history", Role: "assistant", Content: "已经加载"},
	}}
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_history", Title: "历史会话"}, storage.User{ID: "usr_history", Username: "history-user"}, "继续", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if llm.calls != 1 {
		t.Fatalf("expected llm to be called once, got %d", llm.calls)
	}
	reqMessages := llm.reqs[0].Messages
	if len(reqMessages) != 6 {
		t.Fatalf("expected system, 4 historical messages, and current user message, got %#v", reqMessages)
	}
	if len(reqMessages[2].ToolCalls) != 1 || reqMessages[2].ToolCalls[0].ID != "tool_1" || reqMessages[2].ToolCalls[0].Function.Name != "load_skill" {
		t.Fatalf("expected historical assistant tool call before agent loop, got %#v", reqMessages[2])
	}
	if reqMessages[3].Role != "tool" || reqMessages[3].ToolCallID != "tool_1" || reqMessages[3].Content == "" {
		t.Fatalf("expected historical tool result before agent loop, got %#v", reqMessages[3])
	}
	if reqMessages[5].Role != "user" || reqMessages[5].Content != "继续" {
		t.Fatalf("expected current user message after loaded history, got %#v", reqMessages[5])
	}
}

func TestRespondToConversation_PreservesExplicitConversationTitle(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "好的，我继续回答。"},
		}},
	}}}
	config.Client = llm

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	conversation := storage.Conversation{ID: "conv_custom", Title: "周报整理"}
	user := storage.User{ID: "usr_custom", Username: "grace"}

	_, err := service.RespondToConversation(context.Background(), conversation, user, "请继续补充今天的结论", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.updatedID != "" || store.updatedTitle != "" {
		t.Fatalf("expected explicit title to avoid title rewrite, got id=%q title=%q", store.updatedID, store.updatedTitle)
	}
	if store.touchedID != "conv_custom" {
		t.Fatalf("expected explicit-title conversation to touch activity only, got %q", store.touchedID)
	}
}

func TestBuildSystemPrompt_DoesNotRequireToolRegistry(t *testing.T) {
	cfg := testAppConfig(t)
	service := &Service{Cfg: cfg}
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	prompt := service.buildSystemPrompt(storage.User{ID: "usr_6", Username: "frank"}, agenttools.NewSkillSnapshot(nil, loader))
	if !contains(prompt, "Current workspace root: "+cfg.WorkspaceRoot) {
		t.Fatalf("expected prompt to include workspace root, got %q", prompt)
	}
	if !contains(prompt, "rather than a chat-only assistant") {
		t.Fatalf("expected prompt to describe a full agent, got %q", prompt)
	}
	if contains(prompt, "Runtime tools available in this conversation:") {
		t.Fatalf("expected prompt without tool registry to omit tool list, got %q", prompt)
	}
	if !contains(prompt, "- builtin-skill: Builtin description") {
		t.Fatalf("expected prompt to include skill descriptions, got %q", prompt)
	}
}

func TestBuildSkillSnapshotRefreshesUserSkillsFromStoreEachCall(t *testing.T) {
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	store := &fakeStore{enabledSkillSets: [][]storage.Skill{
		{{ID: "skill_1", UserID: "usr_refresh", Name: "User Skill", Slug: "user-skill", Description: "old description", Content: "old body", Status: "enabled"}},
		{{ID: "skill_1", UserID: "usr_refresh", Name: "User Skill", Slug: "user-skill", Description: "new description", Content: "new body", Status: "enabled"}},
	}}
	service := &Service{Store: store, Cfg: testAppConfig(t), BuiltinSkills: builtin}

	first, err := service.buildSkillSnapshot(context.Background(), "usr_refresh")
	if err != nil {
		t.Fatalf("first snapshot: %v", err)
	}
	second, err := service.buildSkillSnapshot(context.Background(), "usr_refresh")
	if err != nil {
		t.Fatalf("second snapshot: %v", err)
	}

	if store.listEnabledCalls != 2 {
		t.Fatalf("expected store to be queried for each snapshot, got %d calls", store.listEnabledCalls)
	}
	if got := first.Merged.Skills["user-skill"].Body; got != "old body" {
		t.Fatalf("expected first snapshot to use old DB content, got %q", got)
	}
	if got := second.Merged.Skills["user-skill"].Body; got != "new body" {
		t.Fatalf("expected second snapshot to refresh DB content, got %q", got)
	}
	if _, ok := second.Merged.Skills["builtin-skill"]; !ok {
		t.Fatalf("expected refreshed snapshot to keep builtin skills")
	}
}

func TestSkillSnapshotLoadSkillPrefersUserSkillThenFallsBackToLocal(t *testing.T) {
	userSkills := sessions.NewSkillLoader()
	userSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"shared-skill": {Meta: map[string]string{"description": "DB description", "tags": "db"}, Body: "db body", Path: "db://skills/skill_1"},
	})
	localSkills := sessions.NewSkillLoader()
	localSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"shared-skill":  {Meta: map[string]string{"description": "Local description"}, Body: "local body", Path: "/skills/shared/SKILL.md"},
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "/skills/builtin/SKILL.md"},
	})
	snapshot := agenttools.NewSkillSnapshot(userSkills, localSkills)

	loaded, err := snapshot.LoadSkill("shared-skill")
	if err != nil {
		t.Fatalf("load shared skill: %v", err)
	}
	if loaded.Source != "db" || loaded.Entry.Body != "db body" {
		t.Fatalf("expected db skill to win, got source=%q body=%q", loaded.Source, loaded.Entry.Body)
	}

	loaded, err = snapshot.LoadSkill("builtin-skill")
	if err != nil {
		t.Fatalf("load builtin skill: %v", err)
	}
	if loaded.Source != "local" || loaded.Entry.Body != "builtin body" {
		t.Fatalf("expected local fallback, got source=%q body=%q", loaded.Source, loaded.Entry.Body)
	}
}

func TestRespondToConversation_ReturnsSuccessfulToolResultIntoLoop(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonToolCalls,
				Message: openai.ChatCompletionMessage{Role: "assistant", ToolCalls: []openai.ToolCall{{
					ID:   "tool_1",
					Type: "function",
					Function: openai.FunctionCall{
						Name:      "load_skill",
						Arguments: `{"name":"builtin-skill"}`,
					},
				}}},
			}},
		},
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonStop,
				Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "final answer"},
			}},
		},
	}}
	config.Client = llm

	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
	conversation := storage.Conversation{ID: "conv_4", Title: "新对话"}
	user := storage.User{ID: "usr_4", Username: "dave"}

	message, err := service.RespondToConversation(context.Background(), conversation, user, "请加载技能后回答", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "final answer" {
		t.Fatalf("expected final assistant answer, got %q", message.Content)
	}
	if llm.calls != 2 {
		t.Fatalf("expected llm to be called twice, got %d", llm.calls)
	}
	if len(store.toolCalls) != 1 {
		t.Fatalf("expected one stored tool call, got %d", len(store.toolCalls))
	}
	if store.toolCalls[0].Status != "success" {
		t.Fatalf("expected tool call status success, got %q", store.toolCalls[0].Status)
	}
	var audit toolExecutionAudit
	if err := json.Unmarshal([]byte(store.toolCalls[0].Summary), &audit); err != nil {
		t.Fatalf("expected tool call summary to be audit JSON, got error: %v, content: %q", err, store.toolCalls[0].Summary)
	}
	if audit.ResolvedCWD == "" || !strings.HasPrefix(audit.ResolvedCWD, cfg.WorkspaceRoot) {
		t.Fatalf("expected audit resolved cwd under workspace root %q, got %q", cfg.WorkspaceRoot, audit.ResolvedCWD)
	}
	if audit.OutcomeSummary == "" || !contains(audit.OutcomeSummary, "<skill source=\"local\" name=\"builtin-skill\">") {
		t.Fatalf("expected audit outcome summary to include tool result, got %#v", audit)
	}
	if audit.DenialReason != "" {
		t.Fatalf("expected no denial reason for successful tool call, got %#v", audit)
	}
	toolMessage := findToolMessage(t, llm.reqs[1].Messages)
	var outcome toolExecutionOutcome
	if err := json.Unmarshal([]byte(toolMessage.Content), &outcome); err != nil {
		t.Fatalf("expected JSON tool message content, got error: %v, content: %q", err, toolMessage.Content)
	}
	if outcome.Status != "success" {
		t.Fatalf("expected tool loop status success, got %q", outcome.Status)
	}
	if !contains(outcome.Result, "<skill source=\"local\" name=\"builtin-skill\">") {
		t.Fatalf("expected tool result to include skill content, got %q", outcome.Result)
	}
	if len(store.historyUpdates) != 1 {
		t.Fatalf("expected one full history update, got %d", len(store.historyUpdates))
	}
	storedHistory := store.historyUpdates[0]
	if len(storedHistory) != 4 {
		t.Fatalf("expected user, assistant tool call, tool result, final assistant in history, got %#v", storedHistory)
	}
	if storedHistory[1].Role != "assistant" || len(storedHistory[1].ToolCalls) != 1 || storedHistory[1].ToolCalls[0].ID != "tool_1" {
		t.Fatalf("expected assistant tool call message in history, got %#v", storedHistory[1])
	}
	if storedHistory[2].Role != "tool" || storedHistory[2].ToolCallID != "tool_1" || storedHistory[2].Content != toolMessage.Content {
		t.Fatalf("expected tool result message in history, got %#v", storedHistory[2])
	}
	if storedHistory[3].Role != "assistant" || storedHistory[3].Content != "final answer" {
		t.Fatalf("expected final assistant message at end of history, got %#v", storedHistory[3])
	}
}

func TestRespondToConversation_ReturnsRejectedToolResultIntoLoop(t *testing.T) {
	originalClient := config.Client
	defer func() { config.Client = originalClient }()

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonToolCalls,
				Message: openai.ChatCompletionMessage{Role: "assistant", ToolCalls: []openai.ToolCall{{
					ID:   "tool_2",
					Type: "function",
					Function: openai.FunctionCall{
						Name:      "load_skill",
						Arguments: `{"name":"missing-skill"}`,
					},
				}}},
			}},
		},
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonStop,
				Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "handled rejection"},
			}},
		},
	}}
	config.Client = llm

	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
	conversation := storage.Conversation{ID: "conv_5", Title: "新对话"}
	user := storage.User{ID: "usr_5", Username: "erin"}

	message, err := service.RespondToConversation(context.Background(), conversation, user, "请尝试加载不存在的技能", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "handled rejection" {
		t.Fatalf("expected final assistant answer, got %q", message.Content)
	}
	if llm.calls != 2 {
		t.Fatalf("expected llm to be called twice, got %d", llm.calls)
	}
	if len(store.toolCalls) != 1 {
		t.Fatalf("expected one stored tool call, got %d", len(store.toolCalls))
	}
	if store.toolCalls[0].Status != "rejected" {
		t.Fatalf("expected tool call status rejected, got %q", store.toolCalls[0].Status)
	}
	var audit toolExecutionAudit
	if err := json.Unmarshal([]byte(store.toolCalls[0].Summary), &audit); err != nil {
		t.Fatalf("expected tool call summary to be audit JSON, got error: %v, content: %q", err, store.toolCalls[0].Summary)
	}
	if audit.ResolvedCWD == "" || !strings.HasPrefix(audit.ResolvedCWD, cfg.WorkspaceRoot) {
		t.Fatalf("expected audit resolved cwd under workspace root %q, got %q", cfg.WorkspaceRoot, audit.ResolvedCWD)
	}
	if audit.DenialReason == "" || !contains(audit.DenialReason, "unknown skill \"missing-skill\"") {
		t.Fatalf("expected audit denial reason to include tool error, got %#v", audit)
	}
	if audit.OutcomeSummary != "" {
		t.Fatalf("expected no outcome summary for rejected tool call, got %#v", audit)
	}
	toolMessage := findToolMessage(t, llm.reqs[1].Messages)
	var outcome toolExecutionOutcome
	if err := json.Unmarshal([]byte(toolMessage.Content), &outcome); err != nil {
		t.Fatalf("expected JSON tool message content, got error: %v, content: %q", err, toolMessage.Content)
	}
	if outcome.Status != "rejected" {
		t.Fatalf("expected tool loop status rejected, got %q", outcome.Status)
	}
	if !contains(outcome.Result, "Error: unknown skill \"missing-skill\"") {
		t.Fatalf("expected rejection result to include tool error, got %q", outcome.Result)
	}
	if len(store.historyUpdates) != 1 {
		t.Fatalf("expected one full history update, got %d", len(store.historyUpdates))
	}
	storedHistory := store.historyUpdates[0]
	if len(storedHistory) != 4 {
		t.Fatalf("expected user, assistant tool call, rejected tool result, final assistant in history, got %#v", storedHistory)
	}
	if storedHistory[1].Role != "assistant" || len(storedHistory[1].ToolCalls) != 1 || storedHistory[1].ToolCalls[0].ID != "tool_2" {
		t.Fatalf("expected assistant tool call message in history, got %#v", storedHistory[1])
	}
	if storedHistory[2].Role != "tool" || storedHistory[2].ToolCallID != "tool_2" || storedHistory[2].Content != toolMessage.Content {
		t.Fatalf("expected rejected tool result message in history, got %#v", storedHistory[2])
	}
}

func TestRespondToConversation_UsesWorkspaceRootAcrossMultiToolTurn(t *testing.T) {
	originalClient := config.Client
	originalBash := agenttools.Handlers["bash"]
	defer func() {
		config.Client = originalClient
		agenttools.Handlers["bash"] = originalBash
	}()

	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.CurrentWorkingDir, nil
	}

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonToolCalls,
				Message: openai.ChatCompletionMessage{Role: "assistant", ToolCalls: []openai.ToolCall{
					{ID: "tool_load", Type: "function", Function: openai.FunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}},
					{ID: "tool_bash", Type: "function", Function: openai.FunctionCall{Name: "bash", Arguments: `{"command":"pwd"}`}},
				}},
			}},
		},
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonStop,
				Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "final answer"},
			}},
		},
	}}
	config.Client = llm

	workspace := t.TempDir()
	skillDir := filepath.Join(workspace, "skills", "builtin-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: filepath.Join(skillDir, "SKILL.md")},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.WorkspaceRoot = workspace
	cfg.WebAllowedTools = []string{"load_skill", "bash"}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_skill", Title: "新对话"}, storage.User{ID: "usr_skill", Username: "skill-user"}, "请加载技能后执行 bash", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "final answer" {
		t.Fatalf("expected final answer, got %q", message.Content)
	}
	if len(store.toolCalls) != 2 {
		t.Fatalf("expected 2 tool calls, got %d", len(store.toolCalls))
	}
	var audit toolExecutionAudit
	if err := json.Unmarshal([]byte(store.toolCalls[1].Summary), &audit); err != nil {
		t.Fatalf("parse second tool audit: %v", err)
	}
	if filepath.Clean(audit.ResolvedCWD) != filepath.Clean(workspace) {
		t.Fatalf("expected bash audit cwd %q, got %#v", workspace, audit)
	}
	toolMessage := findToolMessageByID(t, llm.reqs[1].Messages, "tool_bash")
	var outcome toolExecutionOutcome
	if err := json.Unmarshal([]byte(toolMessage.Content), &outcome); err != nil {
		t.Fatalf("parse bash tool outcome: %v", err)
	}
	if filepath.Clean(outcome.Result) != filepath.Clean(workspace) {
		t.Fatalf("expected bash tool result %q, got %q", workspace, outcome.Result)
	}
	if len(store.historyUpdates) != 1 {
		t.Fatalf("expected one full history update, got %d", len(store.historyUpdates))
	}
	storedHistory := store.historyUpdates[0]
	if len(storedHistory) != 5 {
		t.Fatalf("expected user, assistant tool calls, two tool results, final assistant in history, got %#v", storedHistory)
	}
	if storedHistory[1].Role != "assistant" || len(storedHistory[1].ToolCalls) != 2 {
		t.Fatalf("expected assistant message with two tool calls in history, got %#v", storedHistory[1])
	}
	if storedHistory[2].Role != "tool" || storedHistory[2].ToolCallID != "tool_load" {
		t.Fatalf("expected first tool result in history, got %#v", storedHistory[2])
	}
	if storedHistory[3].Role != "tool" || storedHistory[3].ToolCallID != "tool_bash" {
		t.Fatalf("expected second tool result in history, got %#v", storedHistory[3])
	}
	if storedHistory[4].Role != "assistant" || storedHistory[4].Content != "final answer" {
		t.Fatalf("expected final assistant message at end of history, got %#v", storedHistory[4])
	}
}

func TestRespondToConversation_DBSkillFallsBackToWorkspaceWithinMultiToolTurn(t *testing.T) {
	originalClient := config.Client
	originalBash := agenttools.Handlers["bash"]
	defer func() {
		config.Client = originalClient
		agenttools.Handlers["bash"] = originalBash
	}()

	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.CurrentWorkingDir, nil
	}

	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonToolCalls,
				Message: openai.ChatCompletionMessage{Role: "assistant", ToolCalls: []openai.ToolCall{
					{ID: "tool_load", Type: "function", Function: openai.FunctionCall{Name: "load_skill", Arguments: `{"name":"db-skill"}`}},
					{ID: "tool_bash", Type: "function", Function: openai.FunctionCall{Name: "bash", Arguments: `{"command":"pwd"}`}},
				}},
			}},
		},
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonStop,
				Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "final answer"},
			}},
		},
	}}
	config.Client = llm

	workspace := t.TempDir()
	store := &fakeStore{enabledSkills: []storage.Skill{{
		ID:          "skill_1",
		UserID:      "usr_db",
		Name:        "DB Skill",
		Slug:        "db-skill",
		Description: "db description",
		Content:     "db body",
		Status:      "enabled",
	}}}
	cfg := testAppConfig(t)
	cfg.WorkspaceRoot = workspace
	cfg.WebAllowedTools = []string{"load_skill", "bash"}
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_db", Title: "新对话"}, storage.User{ID: "usr_db", Username: "db-user"}, "请加载 db skill 后执行 bash", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "final answer" {
		t.Fatalf("expected final answer, got %q", message.Content)
	}
	if len(store.toolCalls) != 2 {
		t.Fatalf("expected 2 tool calls, got %d", len(store.toolCalls))
	}
	var audit toolExecutionAudit
	if err := json.Unmarshal([]byte(store.toolCalls[1].Summary), &audit); err != nil {
		t.Fatalf("parse second tool audit: %v", err)
	}
	if filepath.Clean(audit.ResolvedCWD) != filepath.Clean(workspace) {
		t.Fatalf("expected workspace fallback cwd %q, got %#v", workspace, audit)
	}
	toolMessage := findToolMessageByID(t, llm.reqs[1].Messages, "tool_bash")
	var outcome toolExecutionOutcome
	if err := json.Unmarshal([]byte(toolMessage.Content), &outcome); err != nil {
		t.Fatalf("parse bash tool outcome: %v", err)
	}
	if filepath.Clean(outcome.Result) != filepath.Clean(workspace) {
		t.Fatalf("expected bash tool result %q, got %q", workspace, outcome.Result)
	}
}

func TestToolRegistryDefinitions_UsesRegisteredToolDefinition(t *testing.T) {
	registry := NewToolRegistry(config.AppConfig{})

	defs := registry.Definitions()
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

func TestToolRegistryDefinitions_AreLoadedAtRegistryCreation(t *testing.T) {
	cfg := config.AppConfig{WebAllowedTools: []string{"bash"}}
	registry := NewToolRegistry(cfg)
	cfg.WebAllowedTools = []string{"read_file"}

	defs := registry.Definitions()
	if len(defs) != 1 || defs[0].Function == nil || defs[0].Function.Name != "bash" {
		t.Fatalf("expected registry to keep startup-loaded bash definition, got %#v", defs)
	}
}

func TestToolRegistryExecute_LoadSkillReturnsFullSkillInfo(t *testing.T) {
	localSkills := sessions.NewSkillLoader()
	localSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description", "tags": "go,agent"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	snapshot := agenttools.NewSkillSnapshot(nil, localSkills)
	registry := NewToolRegistry(config.AppConfig{
		AppHome:          "/deploy/app",
		WorkspaceRoot:    "/deploy/app/workspace",
		CommandBinDir:    "/deploy/app/workspace/bin",
		CommandScriptDir: "/deploy/app/workspace/cmd",
	})

	result, err := registry.Execute(context.Background(), ToolContext{Skills: snapshot}, "load_skill", `{"name":"builtin-skill"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := result.Output
	if !contains(content, "<skill source=\"local\" name=\"builtin-skill\">") || !contains(content, "builtin body") {
		t.Fatalf("expected loaded skill content, got %q", content)
	}
	if !contains(content, "<metadata>") || !contains(content, "description: Builtin description") || !contains(content, "tags: go,agent") || !contains(content, "path: builtin://builtin-skill") {
		t.Fatalf("expected loaded skill to include full metadata, got %q", content)
	}
	if !contains(content, "<runtime-paths>") || !contains(content, "COMMAND_BIN_DIR=/deploy/app/workspace/bin") || !contains(content, "COMMAND_SCRIPT_DIR=/deploy/app/workspace/cmd") {
		t.Fatalf("expected loaded skill content to include runtime paths, got %q", content)
	}
	if !contains(content, "WORKSPACE_ROOT=/deploy/app/workspace") {
		t.Fatalf("expected loaded skill content to include workspace root, got %q", content)
	}
	if _, ok := agenttools.Handlers["load_skill"]; !ok {
		t.Fatalf("expected load_skill handler to remain registered in shared tool registry")
	}
}

func TestToolRegistryExecute_LoadSkillUsesConfiguredLocalWorkspacePaths(t *testing.T) {
	localSkills := sessions.NewSkillLoader()
	localSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	snapshot := agenttools.NewSkillSnapshot(nil, localSkills)
	registry := NewToolRegistry(config.AppConfig{
		AppHome:          "/repo/app",
		WorkspaceRoot:    "/repo/app/workspace",
		CommandBinDir:    "/repo/app/workspace/bin",
		CommandScriptDir: "/repo/app/workspace/cmd",
	})

	result, err := registry.Execute(context.Background(), ToolContext{Skills: snapshot}, "load_skill", `{"name":"builtin-skill"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := result.Output
	if !contains(content, "COMMAND_BIN_DIR=/repo/app/workspace/bin") || !contains(content, "COMMAND_SCRIPT_DIR=/repo/app/workspace/cmd") {
		t.Fatalf("expected loaded skill content to include configured local workspace paths, got %q", content)
	}
	if !contains(content, "WORKSPACE_ROOT=/repo/app/workspace") {
		t.Fatalf("expected loaded skill content to include local workspace root, got %q", content)
	}
}

func TestToolRegistryExecute_InjectsRuntimeEnvIntoToolHandler(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()

	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.AppHome + "|" + env.CommandBinDir + "|" + env.CommandScriptDir, nil
	}

	registry := NewToolRegistry(config.AppConfig{
		AppHome:          "/deploy/app",
		WorkspaceRoot:    "/deploy/app/workspace",
		CommandBinDir:    "/deploy/app/bin",
		CommandScriptDir: "/deploy/app/cmd",
		WebAllowedTools:  []string{"bash"},
	})

	execResult, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if execResult.Output != "/deploy/app|/deploy/app/bin|/deploy/app/cmd" {
		t.Fatalf("expected runtime env in handler context, got %q", execResult.Output)
	}
}

func TestToolRegistryExecute_UsesConfiguredCommandDirs(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()

	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.WorkspaceRoot + "|" + env.CommandBinDir + "|" + env.CommandScriptDir, nil
	}

	registry := NewToolRegistry(config.AppConfig{
		WorkspaceRoot:    "/deploy/app/workspace",
		CommandBinDir:    "/deploy/app/workspace/bin",
		CommandScriptDir: "/deploy/app/workspace/cmd",
		WebAllowedTools:  []string{"bash"},
	})
	execResult, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if execResult.Output != "/deploy/app/workspace|/deploy/app/workspace/bin|/deploy/app/workspace/cmd" {
		t.Fatalf("expected configured runtime env, got %q", execResult.Output)
	}
}

func TestToolRegistryExecute_InjectsWorkspaceIntoToolHandler(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()

	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.WorkspaceRoot, nil
	}

	registry := NewToolRegistry(config.AppConfig{WorkspaceRoot: "/workspaces/usr_1", WebAllowedTools: []string{"bash"}})
	execResult, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if execResult.Output != "/workspaces/usr_1" {
		t.Fatalf("expected workspace root in handler context, got %q", execResult.Output)
	}
}

func TestToolRegistryExecute_BashUsesWorkspaceRootEvenWhenSkillDirWasPreviouslyLoaded(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()
	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.CurrentWorkingDir, nil
	}

	workspace := t.TempDir()
	registry := NewToolRegistry(config.AppConfig{WorkspaceRoot: workspace, WebAllowedTools: []string{"bash"}})

	result, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != filepath.Clean(workspace) {
		t.Fatalf("expected workspace root in runtime env, got %#v", result)
	}
}

func TestToolRegistryExecute_WorkspaceRootPreservesConfiguredWorkspacePaths(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()
	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.WorkspaceRoot + "|" + env.CommandBinDir + "|" + env.CommandScriptDir + "|" + env.CurrentWorkingDir, nil
	}

	workspace := t.TempDir()
	cfg := config.AppConfig{
		WorkspaceRoot:    workspace,
		CommandBinDir:    filepath.Join(workspace, "bin"),
		CommandScriptDir: filepath.Join(workspace, "cmd"),
		WebAllowedTools:  []string{"bash"},
	}
	registry := NewToolRegistry(cfg)

	result, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Clean(workspace) + "|" + filepath.Join(workspace, "bin") + "|" + filepath.Join(workspace, "cmd") + "|" + filepath.Clean(workspace)
	if result.Output != expected {
		t.Fatalf("expected workspace-derived paths to remain stable, got %q", result.Output)
	}
}

func TestToolRegistryExecute_PropagatesBashSafetyFlags(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()
	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return boolString(env.AllowOutsideWorkspace) + "|" + boolString(env.AllowDangerousCommands), nil
	}

	workspace := t.TempDir()
	cfg := config.AppConfig{
		WorkspaceRoot:              workspace,
		WebAllowedTools:            []string{"bash"},
		BashAllowOutsideWorkspace:  true,
		BashAllowDangerousCommands: true,
	}
	registry := NewToolRegistry(cfg)

	result, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "true|true" {
		t.Fatalf("expected safety flags in runtime env, got %q", result.Output)
	}
}

func TestExecuteToolCall_AuditUsesWorkspaceRootAsResolvedCWD(t *testing.T) {
	workspace := t.TempDir()
	skillDir := filepath.Join(workspace, "skills", "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	command := filepath.Join(skillDir, "pwd.sh")
	if err := os.WriteFile(command, []byte("#!/bin/sh\npwd\n"), 0o755); err != nil {
		t.Fatalf("write helper: %v", err)
	}

	cfg := config.AppConfig{WorkspaceRoot: workspace, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+command+`"}`)

	if outcome.Status != "success" {
		t.Fatalf("expected success, got %#v", outcome)
	}
	if filepath.Clean(outcome.Audit.ResolvedCWD) != filepath.Clean(workspace) {
		t.Fatalf("expected audit cwd %q, got %#v", workspace, outcome.Audit)
	}
}

func TestRegisteredTools_UsesCurrentWebAllowList(t *testing.T) {
	tools := RegisteredTools(config.AppConfig{})
	if len(tools) != 1 || tools[0] != "load_skill" {
		t.Fatalf("expected only load_skill to be registered for web runtime, got %v", tools)
	}
}

func TestRegisteredTools_UsesConfiguredWebAllowList(t *testing.T) {
	tools := RegisteredTools(config.AppConfig{WebAllowedTools: []string{"load_skill", "missing_tool", "load_skill"}})
	if len(tools) != 1 || tools[0] != "load_skill" {
		t.Fatalf("expected configured allowlist to be filtered to registered tools, got %v", tools)
	}
}

func TestExecuteToolCall_AuditCapturesCommandArtifactPath(t *testing.T) {
	root := t.TempDir()
	binDir := filepath.Join(root, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create bin dir: %v", err)
	}
	script := filepath.Join(binDir, "helper")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatalf("write helper script: %v", err)
	}

	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("create workspace dir: %v", err)
	}

	cfg := config.AppConfig{WorkspaceRoot: workspace, CommandBinDir: binDir, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+script+` --flag"}`)

	if outcome.Audit.CommandArtifactPath != script {
		t.Fatalf("expected audit command artifact path %q, got %#v", script, outcome.Audit)
	}
	if outcome.Audit.ResolvedCommandPath != script {
		t.Fatalf("expected audit resolved command path %q, got %#v", script, outcome.Audit)
	}
}

func TestExecuteToolCall_AuditCapturesRejectedWorkspaceEscape(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("create workspace dir: %v", err)
	}
	outside := filepath.Join(root, "outside.sh")
	if err := os.WriteFile(outside, []byte("#!/bin/sh\necho no\n"), 0o755); err != nil {
		t.Fatalf("write outside script: %v", err)
	}

	cfg := config.AppConfig{WorkspaceRoot: workspace, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+outside+`"}`)

	if outcome.Status != "rejected" {
		t.Fatalf("expected rejected outcome, got %#v", outcome)
	}
	if outcome.Audit.ResolvedCWD != workspace {
		t.Fatalf("expected audit cwd %q, got %#v", workspace, outcome.Audit)
	}
	if outcome.Audit.ResolvedCommandPath != outside {
		t.Fatalf("expected audit resolved command path %q, got %#v", outside, outcome.Audit)
	}
	if outcome.Audit.CommandArtifactPath != "" {
		t.Fatalf("expected no command artifact path for outside command, got %#v", outcome.Audit)
	}
	if !contains(outcome.Audit.DenialReason, "command path escapes workspace") {
		t.Fatalf("expected denial reason to mention workspace escape, got %#v", outcome.Audit)
	}
}

func TestExecuteToolCall_AuditClassifiesWorkspaceCommandArtifactSource(t *testing.T) {
	appHome := t.TempDir()
	workspace := filepath.Join(appHome, "workspace")
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create workspace bin dir: %v", err)
	}
	command := filepath.Join(binDir, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatalf("write workspace helper: %v", err)
	}

	cfg := config.AppConfig{AppHome: appHome, WorkspaceRoot: workspace, CommandBinDir: binDir, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+command+`"}`)

	if outcome.Audit.CommandArtifactSource != "workspace" {
		t.Fatalf("expected workspace artifact source, got %#v", outcome.Audit)
	}
}

func TestExecuteToolCall_AuditClassifiesCustomCommandArtifactSource(t *testing.T) {
	appHome := t.TempDir()
	workspace := filepath.Join(appHome, "workspace")
	customRoot := t.TempDir()
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("create workspace dir: %v", err)
	}
	command := filepath.Join(customRoot, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatalf("write custom helper: %v", err)
	}

	cfg := config.AppConfig{AppHome: appHome, WorkspaceRoot: workspace, CommandBinDir: customRoot, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+command+`"}`)

	if outcome.Audit.CommandArtifactSource != "custom" {
		t.Fatalf("expected custom artifact source, got %#v", outcome.Audit)
	}
}

func TestExecuteToolCall_AuditCapturesAppWorkspaceResolution(t *testing.T) {
	appHome := t.TempDir()
	workspace := filepath.Join(appHome, "workspace")
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create workspace bin dir: %v", err)
	}
	command := filepath.Join(binDir, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\npwd\n"), 0o755); err != nil {
		t.Fatalf("write workspace helper: %v", err)
	}

	cfg := config.AppConfig{
		AppHome:         appHome,
		WorkspaceRoot:   workspace,
		CommandBinDir:   binDir,
		WebAllowedTools: []string{"bash"},
	}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+command+`"}`)

	if outcome.Status != "success" {
		t.Fatalf("expected successful outcome, got %#v", outcome)
	}
	if filepath.Clean(outcome.Result) != filepath.Clean(workspace) {
		t.Fatalf("expected command to run in app workspace %q, got %#v", workspace, outcome)
	}
	if outcome.Audit.ResolvedCWD != workspace {
		t.Fatalf("expected audit cwd %q, got %#v", workspace, outcome.Audit)
	}
	if outcome.Audit.CommandArtifactPath != command {
		t.Fatalf("expected audit command artifact path %q, got %#v", command, outcome.Audit)
	}
	if outcome.Audit.CommandArtifactSource != "workspace" {
		t.Fatalf("expected workspace artifact source, got %#v", outcome.Audit)
	}
}

func executeToolCallWithDefaultHooks(t *testing.T, service *Service, name, rawArgs string) toolExecutionOutcome {
	t.Helper()
	if service.Store == nil {
		service.Store = &fakeStore{}
	}
	if service.Hooks == nil {
		service.Hooks = NewDefaultHookManager()
	}
	ctx := &ToolUseContext{
		State:    service.newLoopState(storage.Conversation{ID: "conv_audit"}, storage.User{ID: "usr_audit"}, "", nil, nil),
		ToolCall: openai.ToolCall{ID: "tool_audit", Function: openai.FunctionCall{Name: name, Arguments: rawArgs}},
		Name:     name,
		RawArgs:  rawArgs,
	}
	if err := service.Hooks.RunPreToolUse(context.Background(), ctx); err != nil {
		t.Fatalf("run pre tool hooks: %v", err)
	}
	ctx.Outcome = service.executeToolCall(context.Background(), ToolContext{}, name, rawArgs, ctx.Outcome.Audit)
	if err := service.Hooks.RunPostToolUse(context.Background(), ctx); err != nil {
		t.Fatalf("run post tool hooks: %v", err)
	}
	return ctx.Outcome
}

func findToolMessage(t *testing.T, messages []openai.ChatCompletionMessage) openai.ChatCompletionMessage {
	t.Helper()
	for _, message := range messages {
		if message.Role == "tool" {
			return message
		}
	}
	t.Fatalf("expected tool message in runtime loop")
	return openai.ChatCompletionMessage{}
}

func findToolMessageByID(t *testing.T, messages []openai.ChatCompletionMessage, toolCallID string) openai.ChatCompletionMessage {
	t.Helper()
	for _, message := range messages {
		if message.Role == "tool" && message.ToolCallID == toolCallID {
			return message
		}
	}
	t.Fatalf("expected tool message for toolCallID %q", toolCallID)
	return openai.ChatCompletionMessage{}
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

func boolString(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

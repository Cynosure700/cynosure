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
	reqs      []openai.ChatCompletionRequest
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
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
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
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
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
	cfg := testAppConfig(t)
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(nil, cfg), BuiltinSkills: builtin}
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

func TestBuildSystemPrompt_DoesNotRequireToolRegistry(t *testing.T) {
	cfg := testAppConfig(t)
	service := &Service{Cfg: cfg}
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	prompt := service.buildSystemPrompt(storage.User{ID: "usr_6", Username: "frank"}, loader)
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
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(nil, cfg), BuiltinSkills: builtin}
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
	if audit.OutcomeSummary == "" || !contains(audit.OutcomeSummary, "<skill name=\"builtin-skill\">") {
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
	if !contains(outcome.Result, "<skill name=\"builtin-skill\">") {
		t.Fatalf("expected tool result to include skill content, got %q", outcome.Result)
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
	service := &Service{Store: store, Cfg: cfg, Tools: NewToolRegistry(nil, cfg), BuiltinSkills: builtin}
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
}

func TestResolveUserWorkspace_UsesSharedWorkspaceRoot(t *testing.T) {
	root := t.TempDir()
	service := &Service{Cfg: config.AppConfig{WorkspaceRoot: root}}

	workspace, err := service.resolveUserWorkspace("usr_123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := root
	if workspace != expected {
		t.Fatalf("expected workspace %q, got %q", expected, workspace)
	}
	info, err := os.Stat(workspace)
	if err != nil {
		t.Fatalf("expected workspace to exist: %v", err)
	}
	if !info.IsDir() {
		t.Fatalf("expected workspace path to be directory")
	}
}

func TestResolveUserWorkspace_SharesSameDirectoryAcrossUsers(t *testing.T) {
	root := t.TempDir()
	service := &Service{Cfg: config.AppConfig{WorkspaceRoot: root}}

	first, err := service.resolveUserWorkspace("usr_a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	second, err := service.resolveUserWorkspace("usr_b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if first != second {
		t.Fatalf("expected shared workspace, got %q and %q", first, second)
	}
}

func TestResolveUserWorkspace_AllowsWorkspaceWithEmbeddedDeploymentResources(t *testing.T) {
	root := t.TempDir()
	builtinSkillsDir := filepath.Join(root, "skills")
	commandBinDir := filepath.Join(root, "bin")
	commandScriptDir := filepath.Join(root, "cmd")
	if err := os.MkdirAll(builtinSkillsDir, 0o755); err != nil {
		t.Fatalf("mkdir builtin skills dir: %v", err)
	}
	if err := os.MkdirAll(commandBinDir, 0o755); err != nil {
		t.Fatalf("mkdir command bin dir: %v", err)
	}
	if err := os.MkdirAll(commandScriptDir, 0o755); err != nil {
		t.Fatalf("mkdir command script dir: %v", err)
	}
	service := &Service{Cfg: config.AppConfig{WorkspaceRoot: root, BuiltinSkillsDir: builtinSkillsDir, CommandBinDir: commandBinDir, CommandScriptDir: commandScriptDir}}

	workspace, err := service.resolveUserWorkspace("usr_overlap")
	if err != nil {
		t.Fatalf("expected workspace root with embedded deployment resources to be allowed, got %v", err)
	}
	if workspace != root {
		t.Fatalf("expected workspace %q, got %q", root, workspace)
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
	registry := NewToolRegistry(nil, config.AppConfig{
		AppHome:          "/deploy/app",
		WorkspaceRoot:    "/deploy/app/output/workspace",
		CommandBinDir:    "/deploy/app/output/workspace/bin",
		CommandScriptDir: "/deploy/app/output/workspace/cmd",
	})

	content, err := registry.Execute(context.Background(), ToolContext{Loader: loader}, "load_skill", `{"name":"builtin-skill"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(content, "<skill name=\"builtin-skill\">") || !contains(content, "builtin body") {
		t.Fatalf("expected loaded skill content, got %q", content)
	}
	if !contains(content, "<deployment-paths>") || !contains(content, "COMMAND_BIN_DIR=/deploy/app/output/workspace/bin") || !contains(content, "COMMAND_SCRIPT_DIR=/deploy/app/output/workspace/cmd") {
		t.Fatalf("expected loaded skill content to include deployment paths, got %q", content)
	}
	if !contains(content, "WORKSPACE_ROOT=/deploy/app/output/workspace") {
		t.Fatalf("expected loaded skill content to include workspace root, got %q", content)
	}
	if _, ok := agenttools.Handlers["load_skill"]; !ok {
		t.Fatalf("expected load_skill handler to remain registered in shared tool registry")
	}
}

func TestToolRegistryExecute_LoadSkillFallsBackToLocalWorkspacePaths(t *testing.T) {
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	registry := NewToolRegistry(nil, config.AppConfig{AppHome: "/repo/app", WorkspaceRoot: "/repo/app/workspace"})

	content, err := registry.Execute(context.Background(), ToolContext{Loader: loader}, "load_skill", `{"name":"builtin-skill"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(content, "COMMAND_BIN_DIR=/repo/app/workspace/bin") || !contains(content, "COMMAND_SCRIPT_DIR=/repo/app/workspace/cmd") {
		t.Fatalf("expected loaded skill content to derive command paths from local workspace, got %q", content)
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

	registry := NewToolRegistry(nil, config.AppConfig{
		AppHome:          "/deploy/app",
		WorkspaceRoot:    "/deploy/app/workspace",
		CommandBinDir:    "/deploy/app/bin",
		CommandScriptDir: "/deploy/app/cmd",
		WebAllowedTools:  []string{"bash"},
	})

	result, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "/deploy/app|/deploy/app/bin|/deploy/app/cmd" {
		t.Fatalf("expected runtime env in handler context, got %q", result)
	}
}

func TestToolRegistryExecute_UsesWorkspaceDerivedCommandDirsWhenUnset(t *testing.T) {
	original := agenttools.Handlers["bash"]
	defer func() { agenttools.Handlers["bash"] = original }()

	agenttools.Handlers["bash"] = func(ctx context.Context, args map[string]any) (string, error) {
		env, ok := agenttools.RuntimeEnvFromContext(ctx)
		if !ok {
			return "", nil
		}
		return env.WorkspaceRoot + "|" + env.CommandBinDir + "|" + env.CommandScriptDir, nil
	}

	registry := NewToolRegistry(nil, config.AppConfig{WorkspaceRoot: "/deploy/app/workspace", WebAllowedTools: []string{"bash"}})
	result, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "/deploy/app/workspace|/deploy/app/workspace/bin|/deploy/app/workspace/cmd" {
		t.Fatalf("expected workspace-derived runtime env, got %q", result)
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

	registry := NewToolRegistry(nil, config.AppConfig{WebAllowedTools: []string{"bash"}})
	result, err := registry.Execute(context.Background(), ToolContext{WorkspaceRoot: "/workspaces/usr_1"}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "/workspaces/usr_1" {
		t.Fatalf("expected workspace root in handler context, got %q", result)
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

	cfg := config.AppConfig{CommandBinDir: binDir, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
	outcome := service.executeToolCall(context.Background(), ToolContext{WorkspaceRoot: filepath.Join(root, "workspace")}, "bash", `{"command":"`+script+` --flag"}`)

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
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
	outcome := service.executeToolCall(context.Background(), ToolContext{WorkspaceRoot: workspace}, "bash", `{"command":"`+outside+`"}`)

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

func TestExecuteToolCall_AuditClassifiesDeploymentCommandArtifactSource(t *testing.T) {
	appHome := t.TempDir()
	workspace := filepath.Join(appHome, "output", "workspace")
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create deployment bin dir: %v", err)
	}
	command := filepath.Join(binDir, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatalf("write deployment helper: %v", err)
	}

	cfg := config.AppConfig{AppHome: appHome, WorkspaceRoot: workspace, CommandBinDir: binDir, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
	outcome := service.executeToolCall(context.Background(), ToolContext{WorkspaceRoot: workspace}, "bash", `{"command":"`+command+`"}`)

	if outcome.Audit.CommandArtifactSource != "deployment" {
		t.Fatalf("expected deployment artifact source, got %#v", outcome.Audit)
	}
}

func TestExecuteToolCall_AuditClassifiesLocalCommandArtifactSource(t *testing.T) {
	appHome := t.TempDir()
	workspace := filepath.Join(appHome, "workspace")
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create local bin dir: %v", err)
	}
	command := filepath.Join(binDir, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatalf("write local helper: %v", err)
	}

	cfg := config.AppConfig{AppHome: appHome, WorkspaceRoot: workspace, CommandBinDir: binDir, WebAllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
	outcome := service.executeToolCall(context.Background(), ToolContext{WorkspaceRoot: workspace}, "bash", `{"command":"`+command+`"}`)

	if outcome.Audit.CommandArtifactSource != "local" {
		t.Fatalf("expected local artifact source, got %#v", outcome.Audit)
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
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
	outcome := service.executeToolCall(context.Background(), ToolContext{WorkspaceRoot: workspace}, "bash", `{"command":"`+command+`"}`)

	if outcome.Audit.CommandArtifactSource != "custom" {
		t.Fatalf("expected custom artifact source, got %#v", outcome.Audit)
	}
}

func TestExecuteToolCall_AuditCapturesDeploymentWorkspaceResolution(t *testing.T) {
	appHome := t.TempDir()
	workspace := filepath.Join(appHome, "output", "workspace")
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create deployment bin dir: %v", err)
	}
	command := filepath.Join(binDir, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\npwd\n"), 0o755); err != nil {
		t.Fatalf("write deployment helper: %v", err)
	}

	cfg := config.AppConfig{
		AppHome:       appHome,
		WorkspaceRoot: workspace,
		CommandBinDir: binDir,
		WebAllowedTools: []string{"bash"},
	}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(nil, cfg)}
	outcome := service.executeToolCall(context.Background(), ToolContext{WorkspaceRoot: workspace}, "bash", `{"command":"`+command+`"}`)

	if outcome.Status != "success" {
		t.Fatalf("expected successful outcome, got %#v", outcome)
	}
	if filepath.Clean(outcome.Result) != filepath.Clean(workspace) {
		t.Fatalf("expected command to run in deployment workspace %q, got %#v", workspace, outcome)
	}
	if outcome.Audit.ResolvedCWD != workspace {
		t.Fatalf("expected audit cwd %q, got %#v", workspace, outcome.Audit)
	}
	if outcome.Audit.CommandArtifactPath != command {
		t.Fatalf("expected audit command artifact path %q, got %#v", command, outcome.Audit)
	}
	if outcome.Audit.CommandArtifactSource != "deployment" {
		t.Fatalf("expected deployment artifact source, got %#v", outcome.Audit)
	}
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

package runtime

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	llmpkg "nano_cc/internal/llm"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
)

type fakeStore struct {
	historyUpdates       [][]storage.Message
	toolCalls            []storage.ToolCall
	subagentMessages     []storage.SubagentMessage
	persistedOutputs     []storage.PersistedOutput
	contextSummaries     []storage.ContextSummary
	memories             []storage.Memory
	conversationMemories []storage.ConversationMemory
	replacedConvMemories []storage.ConversationMemory
	modelHistory         []storage.Message
	modelHistoryExists   bool
	modelHistoryErr      error
	upsertedModelHistory [][]storage.Message
	deletedOldest        [][3]string
	replacedMemories     []storage.Memory
	cached               []storage.Message
	updatedTitle         string
	updatedID            string
	touchedID            string
	toolResultLogs       []storage.ToolResultLogEntry
	lockReleased         int
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

func (f *fakeStore) CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error {
	f.subagentMessages = append(f.subagentMessages, message)
	return nil
}

func (f *fakeStore) AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error {
	f.toolResultLogs = append(f.toolResultLogs, entry)
	return nil
}

func (f *fakeStore) CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error {
	f.persistedOutputs = append(f.persistedOutputs, output)
	return nil
}

func (f *fakeStore) GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error) {
	for _, o := range f.persistedOutputs {
		if o.ID == id && o.UserID == userID && o.ConversationID == conversationID {
			return o, nil
		}
	}
	return storage.PersistedOutput{}, errors.New("persisted output not found")
}

func (f *fakeStore) GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error) {
	for _, o := range f.persistedOutputs {
		if o.ConversationID == conversationID && o.UserID == userID && o.MessageID == messageID && o.ToolCallID == toolCallID && o.Strategy == strategy && o.ContentSHA256 == contentSHA256 {
			return o, nil
		}
	}
	return storage.PersistedOutput{}, errors.New("persisted output not found")
}

func (f *fakeStore) CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error {
	f.contextSummaries = append(f.contextSummaries, summary)
	return nil
}

func (f *fakeStore) GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error) {
	for _, c := range f.contextSummaries {
		if c.ConversationID == conversationID && c.UserID == userID && c.SourceHistorySHA256 == sourceHistorySHA256 {
			return c, nil
		}
	}
	return storage.ContextSummary{}, errors.New("context summary not found")
}

func (f *fakeStore) ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	var result []storage.Memory
	for _, m := range f.memories {
		if m.UserID == userID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (f *fakeStore) ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error) {
	var result []storage.Memory
	for _, m := range f.memories {
		if m.UserID == userID && m.Type == memType {
			result = append(result, m)
		}
	}
	return result, nil
}

func (f *fakeStore) InsertMemory(ctx context.Context, m storage.Memory) error {
	f.memories = append(f.memories, m)
	return nil
}

func (f *fakeStore) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	count := 0
	for _, m := range f.memories {
		if m.UserID == userID && m.Type == memType {
			count++
		}
	}
	return count, nil
}

func (f *fakeStore) DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error {
	f.deletedOldest = append(f.deletedOldest, [3]string{userID, memType, strconv.Itoa(n)})
	return nil
}

func (f *fakeStore) ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error {
	f.replacedMemories = append(f.replacedMemories, items...)
	return nil
}

func (f *fakeStore) ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error) {
	var result []storage.ConversationMemory
	for _, m := range f.conversationMemories {
		if m.ConversationID == conversationID {
			result = append(result, m)
		}
	}
	return result, nil
}

func (f *fakeStore) ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error {
	f.replacedConvMemories = append(f.replacedConvMemories, items...)
	return nil
}

func (f *fakeStore) GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error) {
	if f.modelHistoryErr != nil {
		return nil, false, f.modelHistoryErr
	}
	if !f.modelHistoryExists {
		return nil, false, nil
	}
	return f.modelHistory, true, nil
}

func (f *fakeStore) UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error {
	f.upsertedModelHistory = append(f.upsertedModelHistory, append([]storage.Message(nil), messages...))
	f.modelHistory = append([]storage.Message(nil), messages...)
	f.modelHistoryExists = true
	return nil
}

func (f *fakeStore) AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error) {
	return true, nil
}

func (f *fakeStore) RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error) {
	return true, nil
}

func (f *fakeStore) ReleaseConversationLock(ctx context.Context, conversationID, token string) error {
	f.lockReleased++
	return nil
}

type fakeLLMClient struct {
	responses       []openai.ChatCompletionResponse
	streamChunks    []openai.ChatCompletionStreamResponse
	streamChunkSets [][]openai.ChatCompletionStreamResponse
	calls           int
	lastReq         openai.ChatCompletionRequest
	reqs            []openai.ChatCompletionRequest
}

type fakeChatStream struct {
	chunks []openai.ChatCompletionStreamResponse
	errAt  int
	index  int
}

func (s *fakeChatStream) Recv() (openai.ChatCompletionStreamResponse, error) {
	if s.errAt > 0 && s.index == s.errAt {
		return openai.ChatCompletionStreamResponse{}, errors.New("stream failed")
	}
	if s.index >= len(s.chunks) {
		return openai.ChatCompletionStreamResponse{}, io.EOF
	}
	chunk := s.chunks[s.index]
	s.index++
	return chunk, nil
}

func (s *fakeChatStream) Close() error { return nil }

func intPointer(value int) *int { return &value }

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

// nonMetaEvents 返回除去 meta 事件的事件序列，便于断言主流事件顺序。
func (w *captureEventWriter) nonMetaEvents() []capturedEvent {
	out := make([]capturedEvent, 0, len(w.events))
	for _, e := range w.events {
		if e.name == "meta" {
			continue
		}
		out = append(out, e)
	}
	return out
}

func nonToolLifecycleEvents(events []capturedEvent) []capturedEvent {
	out := make([]capturedEvent, 0, len(events))
	for _, e := range events {
		if e.name == toolCallStartEvent || e.name == toolCallDoneEvent {
			continue
		}
		out = append(out, e)
	}
	return out
}

func eventMap(t *testing.T, data any) map[string]any {
	t.Helper()
	m, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("event data = %#v, want map[string]any", data)
	}
	return m
}

func stringValue(value any) string {
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}

func testAppConfig(t *testing.T) config.AppConfig {
	t.Helper()
	return config.AppConfig{
		LLM:           config.Config{ModelID: "test-model"},
		WorkspaceRoot: t.TempDir(),
	}
}

func TestBuildSubagentSystemPromptAppendsStructuredChildAgentContext(t *testing.T) {
	service := &Service{
		Cfg:        testAppConfig(t),
		BasePrompt: "Base prompt.",
	}

	prompt := service.buildSubagentSystemPrompt(storage.User{Username: "agent-user"}, nil)

	for _, want := range []string{
		"<subagent>",
		"你是由 `spawn_subagent` 派生出来的子智能体。",
		"- 你看不到父对话的历史记录。",
		"- 只能依据当前任务和工作区文件来工作。",
		"- 不要调用 `spawn_subagent`。",
		"- 完成后，只输出一段简洁的摘要，说明你做了什么、关键发现以及尚未解决的问题。",
		"</subagent>",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected subagent prompt to contain %q, got %q", want, prompt)
		}
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

func (f *fakeLLMClient) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (llmpkg.ChatCompletionStream, error) {
	f.calls++
	f.lastReq = req
	f.reqs = append(f.reqs, req)
	if len(f.streamChunkSets) > 0 {
		chunks := f.streamChunkSets[0]
		f.streamChunkSets = f.streamChunkSets[1:]
		return &fakeChatStream{chunks: chunks}, nil
	}
	if len(f.streamChunks) > 0 {
		stream := &fakeChatStream{chunks: f.streamChunks}
		f.streamChunks = nil
		return stream, nil
	}
	if len(f.responses) > 0 {
		resp := f.responses[0]
		f.responses = f.responses[1:]
		chunks := make([]openai.ChatCompletionStreamResponse, 0, len(resp.Choices))
		for _, choice := range resp.Choices {
			chunks = append(chunks, openai.ChatCompletionStreamResponse{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: choice.Message.Content, ReasoningContent: choice.Message.ReasoningContent, ToolCalls: choice.Message.ToolCalls}, FinishReason: choice.FinishReason}}})
		}
		return &fakeChatStream{chunks: chunks}, nil
	}
	stream := &fakeChatStream{}
	return stream, nil
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
	longResult := strings.Join([]string{"line1", "line2", "line3", "line4", "line5", "line6", "line7"}, "\n")
	toolCtx := &ToolUseContext{
		State:    state,
		ToolCall: openai.ToolCall{ID: "tool_1", Function: openai.FunctionCall{Name: "bash", Arguments: `{"command":"pwd"}`}},
		Name:     "bash",
		RawArgs:  `{"command":"pwd"}`,
		Outcome:  toolExecutionOutcome{Status: "success", Result: longResult},
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
	if len(writer.events) != 0 {
		t.Fatalf("expected no tool event to be emitted, got %#v", writer.events)
	}
	if len(state.Messages) != 2 || state.Messages[1].Role != "tool" || state.Messages[1].ToolCallID != "tool_1" {
		t.Fatalf("expected tool message to be appended, got %#v", state.Messages)
	}
	if len(store.toolResultLogs) != 1 || store.toolResultLogs[0].ToolName != "bash" || store.toolResultLogs[0].Result != longResult || store.toolResultLogs[0].RawArgs != `{"command":"pwd"}` {
		t.Fatalf("expected tool result log to be appended, got %#v", store.toolResultLogs)
	}
}

func TestRespondToConversation_DirectAnswerWithoutTools(t *testing.T) {
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "当然可以，我来帮你规划。"},
		}},
	}}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
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
	if len(llm.lastReq.Tools) != 2 || !toolNamesInclude(llm.lastReq.Tools, "load_skill") || !toolNamesInclude(llm.lastReq.Tools, "read_persisted_output") {
		t.Fatalf("expected startup-loaded default load_skill tool plus read_persisted_output, got %#v", llm.lastReq.Tools)
	}
	if len(store.toolCalls) != 0 {
		t.Fatalf("expected no tool calls to be stored, got %d", len(store.toolCalls))
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

func TestRespondToConversation_StreamsAssistantContentDeltas(t *testing.T) {
	llm := &fakeLLMClient{streamChunks: []openai.ChatCompletionStreamResponse{
		{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "你"}}}},
		{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "好"}, FinishReason: openai.FinishReasonStop}}},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	writer := &captureEventWriter{}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_stream", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "打招呼", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "你好" {
		t.Fatalf("expected streamed content to be persisted, got %q", message.Content)
	}
	events := writer.nonMetaEvents()
	if len(events) != 3 {
		t.Fatalf("expected 2 deltas and final assistant, got %#v", events)
	}
	if events[0].name != "assistant_delta" || events[1].name != "assistant_delta" || events[2].name != "assistant" {
		t.Fatalf("unexpected events: %#v", events)
	}
}

func TestRespondToConversation_ContextCancelSkipsStopHooksAndReleasesLock(t *testing.T) {
	llm := &fakeLLMClient{streamChunks: []openai.ChatCompletionStreamResponse{
		{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "部"}}}},
		{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "分"}, FinishReason: openai.FinishReasonStop}}},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.ConversationLockTTL = time.Minute
	cfg.ConversationLockWaitTimeout = time.Millisecond
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), EnableMemory: true}
	writer := &captureEventWriter{}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := service.RespondToConversation(ctx, storage.Conversation{ID: "conv_cancel", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice", MemoryEnabled: true}, "写长回答", writer)
	if err == nil {
		t.Fatal("expected canceled context error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if len(store.historyUpdates) != 0 {
		t.Fatalf("expected no stop hook history persistence after cancel, got %#v", store.historyUpdates)
	}
	if len(store.upsertedModelHistory) != 0 || len(store.replacedConvMemories) != 0 {
		t.Fatalf("expected no memory/model-history handoff after cancel, got model=%#v memories=%#v", store.upsertedModelHistory, store.replacedConvMemories)
	}
	if store.lockReleased != 1 {
		t.Fatalf("expected conversation lock to be released once, got %d", store.lockReleased)
	}
	for _, event := range writer.nonMetaEvents() {
		if event.name == "assistant" {
			t.Fatalf("expected no final assistant event after cancel, got %#v", writer.events)
		}
	}
}

func TestRespondToConversation_PersistsReasoningContent(t *testing.T) {
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "最终答案", ReasoningContent: "内部推理过程"},
		}},
	}}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
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
	events := writer.nonMetaEvents()
	if len(events) != 3 || events[2].name != "assistant" {
		t.Fatalf("expected streamed deltas and final assistant event, got %#v", events)
	}
	payload, ok := events[2].data.(map[string]any)
	if !ok || payload["reasoning_content"] != "内部推理过程" {
		t.Fatalf("expected assistant event to include reasoning_content, got %#v", events[2].data)
	}
}

func TestRespondToConversation_StreamsReasoningContentDeltas(t *testing.T) {
	llm := &fakeLLMClient{streamChunks: []openai.ChatCompletionStreamResponse{
		{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ReasoningContent: "先思考"}}}},
		{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "答案"}, FinishReason: openai.FinishReasonStop}}},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	writer := &captureEventWriter{}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_reasoning_stream", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "解释", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "答案" || message.ReasoningContent != "先思考" {
		t.Fatalf("expected streamed reasoning and content, got %#v", message)
	}
	events := writer.nonMetaEvents()
	if len(events) != 3 {
		t.Fatalf("expected reasoning delta, content delta and final assistant, got %#v", events)
	}
	if events[0].name != "reasoning_delta" || events[1].name != "assistant_delta" || events[2].name != "assistant" {
		t.Fatalf("unexpected events: %#v", events)
	}
	payload, ok := events[2].data.(map[string]any)
	if !ok || payload["message_id"] == "" || payload["final"] != true || payload["reasoning_content"] != "先思考" {
		t.Fatalf("expected final assistant metadata and reasoning content, got %#v", events[2].data)
	}
}

func TestRespondToConversation_StreamsToolCallsAcrossChunks(t *testing.T) {
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: "tool_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "bash"}}}}}}},
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), Function: openai.FunctionCall{Arguments: `{"command"`}}}}}}},
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), Function: openai.FunctionCall{Arguments: `:"pwd"}`}}}}, FinishReason: openai.FinishReasonToolCalls}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "工具执行完成"}, FinishReason: openai.FinishReasonStop}}},
		},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"bash"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	writer := &captureEventWriter{}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_tool_stream", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "执行 pwd", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "工具执行完成" {
		t.Fatalf("expected final streamed answer after tool call, got %q", message.Content)
	}
	if len(store.toolCalls) != 1 || store.toolCalls[0].ToolName != "bash" || store.toolCalls[0].Status != "success" {
		t.Fatalf("expected one successful bash tool call, got %#v", store.toolCalls)
	}
	events := nonToolLifecycleEvents(writer.nonMetaEvents())
	if len(events) < 2 || events[0].name != "assistant_delta" || events[len(events)-1].name != "assistant" {
		t.Fatalf("expected streamed final answer after tool call, got %#v", events)
	}
	finalPayload, ok := events[len(events)-1].data.(map[string]any)
	if !ok || finalPayload["tool_call_count"] != 1 {
		t.Fatalf("expected final assistant event to report 1 tool call, got %#v", events[len(events)-1].data)
	}
	if tokens, _ := finalPayload["context_tokens"].(int); tokens <= 0 {
		t.Fatalf("expected positive context_tokens in final assistant event, got %#v", finalPayload["context_tokens"])
	}
	if len(store.historyUpdates) == 0 {
		t.Fatalf("expected history to be persisted")
	}
	lastUpdate := store.historyUpdates[len(store.historyUpdates)-1]
	finalMsg := lastUpdate[len(lastUpdate)-1]
	if finalMsg.Meta == nil || finalMsg.Meta.ToolCallCount != 1 || finalMsg.Meta.ContextTokens <= 0 {
		t.Fatalf("expected persisted assistant message meta with tool count and tokens, got %#v", finalMsg.Meta)
	}
	var metaEvents int
	for _, e := range writer.events {
		if e.name == "meta" {
			metaEvents++
		}
	}
	if metaEvents == 0 {
		t.Fatalf("expected at least one meta event to be emitted during streaming")
	}
}

func TestRespondToConversation_TodoWriteUpdatesLoopStateTodos(t *testing.T) {
	todoArgs := `{"todos":[{"id":"1","content":"梳理需求","status":"completed"},{"id":"2","content":"实现功能","status":"in_progress"}]}`
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: "todo_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "todo_write", Arguments: todoArgs}}}}, FinishReason: openai.FinishReasonToolCalls}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "已更新计划"}, FinishReason: openai.FinishReasonStop}}},
		},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"todo_write"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	manager := NewDefaultHookManager()
	var observedTodos []agenttools.TodoItem
	manager.PostToolUse = append(manager.PostToolUse, func(ctx context.Context, h *ToolUseContext) error {
		observedTodos = append([]agenttools.TodoItem(nil), h.State.Todos...)
		return nil
	})
	service.Hooks = manager

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_todo", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "实现 todo 工具", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []agenttools.TodoItem{{ID: "1", Content: "梳理需求", Status: "completed"}, {ID: "2", Content: "实现功能", Status: "in_progress"}}
	if len(observedTodos) != len(expected) {
		t.Fatalf("expected observed todos %#v, got %#v", expected, observedTodos)
	}
	for i := range expected {
		if observedTodos[i] != expected[i] {
			t.Fatalf("todo[%d] expected %#v, got %#v", i, expected[i], observedTodos[i])
		}
	}
}

func TestRespondToConversation_TodoListReadsLatestTodoWriteState(t *testing.T) {
	todoArgs := `{"todos":[{"id":"1","content":"梳理需求","status":"completed"},{"id":"2","content":"实现功能","status":"in_progress"}]}`
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: "todo_write_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "todo_write", Arguments: todoArgs}}}}, FinishReason: openai.FinishReasonToolCalls}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: "todo_list_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "todo_list", Arguments: `{}`}}}}, FinishReason: openai.FinishReasonToolCalls}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "已查询计划"}, FinishReason: openai.FinishReasonStop}}},
		},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"todo_write", "todo_list"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_todo_list", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "查询当前 todo", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(store.toolCalls) != 2 {
		t.Fatalf("expected 2 tool calls, got %#v", store.toolCalls)
	}
	toolMessage := findToolMessageByID(t, llm.reqs[2].Messages, "todo_list_1")
	var outcome toolExecutionOutcome
	if err := json.Unmarshal([]byte(toolMessage.Content), &outcome); err != nil {
		t.Fatalf("expected JSON tool message content, got error: %v, content: %q", err, toolMessage.Content)
	}
	result := outcome.Result
	for _, want := range []string{
		"Todo list: 2 items (pending: 0, in_progress: 1, completed: 1).",
		"[completed] 1: 梳理需求",
		"[in_progress] 2: 实现功能",
	} {
		if !strings.Contains(result, want) {
			t.Fatalf("expected todo_list result to contain %q, got %q", want, result)
		}
	}
}

func TestRespondToConversation_InjectsTodoWriteReminderAfterThreeRoundsWithoutTodoWrite(t *testing.T) {
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		bashToolRound("tool_1"),
		bashToolRound("tool_2"),
		bashToolRound("tool_3"),
		{{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "完成"}, FinishReason: openai.FinishReasonStop}}}},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"bash", "todo_write"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_reminder", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "连续执行工具", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(llm.reqs) != 4 {
		t.Fatalf("expected 4 model requests, got %d", len(llm.reqs))
	}
	if !requestContainsTodoWriteReminder(llm.reqs[3]) {
		t.Fatalf("expected fourth request to contain todo_write reminder, got %#v", llm.reqs[3].Messages)
	}
}

func TestRespondToConversation_TodoWriteResetsReminderCounter(t *testing.T) {
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		bashToolRound("tool_1"),
		todoWriteToolRound("todo_1"),
		bashToolRound("tool_2"),
		bashToolRound("tool_3"),
		{{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "完成"}, FinishReason: openai.FinishReasonStop}}}},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"bash", "todo_write"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_reminder_reset", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "执行并更新计划", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, req := range llm.reqs {
		if requestContainsTodoWriteReminder(req) {
			t.Fatalf("todo_write should reset reminder counter; request %d contained reminder: %#v", i, req.Messages)
		}
	}
}

func TestRespondToConversation_DoesNotInjectTodoWriteReminderWhenToolDisabled(t *testing.T) {
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		bashToolRound("tool_1"),
		bashToolRound("tool_2"),
		bashToolRound("tool_3"),
		{{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "完成"}, FinishReason: openai.FinishReasonStop}}}},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"bash"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_reminder_disabled", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "连续执行工具", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for i, req := range llm.reqs {
		if requestContainsTodoWriteReminder(req) {
			t.Fatalf("todo_write disabled; request %d should not contain reminder: %#v", i, req.Messages)
		}
	}
}

func TestRespondToConversation_StreamsReasoningForToolRoundsButNotContent(t *testing.T) {
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ReasoningContent: "先思考要不要调用工具"}}}},
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "我先查一下"}}}},
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: "tool_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}}}}, FinishReason: openai.FinishReasonToolCalls}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ReasoningContent: "工具已返回，组织最终答案"}}}},
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "最终答案"}, FinishReason: openai.FinishReasonStop}}},
		},
	}}
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
	writer := &captureEventWriter{}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_tool_round_text", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "先查资料再回答", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedReasoning := "先思考要不要调用工具工具已返回，组织最终答案"
	if message.Content != "最终答案" || message.ReasoningContent != expectedReasoning {
		t.Fatalf("expected final content and reasoning, got %#v", message)
	}

	var sawToolRoundReasoningDelta bool
	var sawFinalReasoningDelta bool
	var sawFinalAssistantDelta bool
	var sawFinalAssistantWithFullReasoning bool
	for _, event := range writer.events {
		payload, _ := event.data.(map[string]any)
		content, _ := payload["content"].(string)
		if event.name == "reasoning_delta" && content == "先思考要不要调用工具" {
			sawToolRoundReasoningDelta = true
		}
		if event.name == "assistant_delta" && content == "我先查一下" {
			t.Fatalf("tool round assistant delta should not be streamed: %#v", writer.events)
		}
		if event.name == "reasoning_delta" && content == "工具已返回，组织最终答案" {
			sawFinalReasoningDelta = true
		}
		if event.name == "assistant_delta" && content == "最终答案" {
			sawFinalAssistantDelta = true
		}
		if event.name == "assistant" && payload["reasoning_content"] == expectedReasoning {
			sawFinalAssistantWithFullReasoning = true
		}
	}
	if !sawToolRoundReasoningDelta || !sawFinalReasoningDelta || !sawFinalAssistantDelta || !sawFinalAssistantWithFullReasoning {
		t.Fatalf("expected tool round reasoning, final round reasoning, final assistant delta and full final reasoning, got %#v", writer.events)
	}
}

func TestShouldEmitAssistantContentDeltas(t *testing.T) {
	tests := []struct {
		name         string
		finishReason openai.FinishReason
		toolCalls    []openai.ToolCall
		expected     bool
	}{
		{name: "stop without tool calls", finishReason: openai.FinishReasonStop, expected: true},
		{name: "tool calls finish reason", finishReason: openai.FinishReasonToolCalls, toolCalls: []openai.ToolCall{{ID: "tool_1"}}, expected: false},
		{name: "tool calls with stop finish reason", finishReason: openai.FinishReasonStop, toolCalls: []openai.ToolCall{{ID: "tool_1"}}, expected: false},
		{name: "empty finish reason without tool calls", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldEmitAssistantContentDeltas(tt.finishReason, tt.toolCalls); got != tt.expected {
				t.Fatalf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestRespondToConversation_SpawnSubagentUsesFreshMessagesStoresTraceAndDoesNotEmitToolEvents(t *testing.T) {
	spawnArgs := `{"task":"inspect workspace only","cwd":"."}`
	llm := &fakeLLMClient{streamChunkSets: [][]openai.ChatCompletionStreamResponse{
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: "spawn_1", Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "spawn_subagent", Arguments: spawnArgs}}}}, FinishReason: openai.FinishReasonToolCalls}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "sub summary"}, FinishReason: openai.FinishReasonStop}}},
		},
		{
			{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{Content: "final answer"}, FinishReason: openai.FinishReasonStop}}},
		},
	}}
	store := &fakeStore{cached: []storage.Message{{ID: "old_msg", ConversationID: "conv_spawn", UserID: "usr_1", Role: "user", Content: "previous context must not leak"}}}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"spawn_subagent"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	writer := &captureEventWriter{}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_spawn", Title: "已有标题"}, storage.User{ID: "usr_1", Username: "alice"}, "parent request", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if message.Content != "final answer" {
		t.Fatalf("expected final answer, got %q", message.Content)
	}
	if llm.calls != 3 {
		t.Fatalf("expected main, subagent, main llm calls, got %d", llm.calls)
	}
	subReq := llm.reqs[1]
	if len(subReq.Messages) != 2 || subReq.Messages[1].Role != "user" || subReq.Messages[1].Content != "inspect workspace only" {
		t.Fatalf("expected subagent to receive fresh system+task messages only, got %#v", subReq.Messages)
	}
	for _, msg := range subReq.Messages {
		if contains(msg.Content, "previous context must not leak") || contains(msg.Content, "parent request") {
			t.Fatalf("subagent messages leaked parent context: %#v", subReq.Messages)
		}
	}
	for _, event := range writer.events {
		if event.name != "tool" {
			continue
		}
		payload, _ := event.data.(map[string]any)
		if payload["name"] == "spawn_subagent" {
			t.Fatalf("spawn_subagent tool event should not be emitted to frontend: %#v", writer.events)
		}
	}
	if len(store.subagentMessages) < 2 {
		t.Fatalf("expected subagent trace messages to be stored, got %#v", store.subagentMessages)
	}
	if store.subagentMessages[0].Role != "user" || store.subagentMessages[0].Content != "inspect workspace only" {
		t.Fatalf("expected stored subagent task, got %#v", store.subagentMessages[0])
	}
	lastTrace := store.subagentMessages[len(store.subagentMessages)-1]
	if lastTrace.Role != "assistant" || lastTrace.Content != "sub summary" || lastTrace.ParentToolCallID != "spawn_1" {
		t.Fatalf("expected stored subagent summary tied to parent tool call, got %#v", lastTrace)
	}
}

func TestRespondToConversation_ReturnsErrorForEmptyModelStream(t *testing.T) {
	llm := &fakeLLMClient{streamChunks: []openai.ChatCompletionStreamResponse{{}}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_empty_stream", Title: "新对话"}, storage.User{ID: "usr_1", Username: "alice"}, "空响应", nil)
	if err == nil || !strings.Contains(err.Error(), "model stream returned no choices") {
		t.Fatalf("expected empty stream error, got %v", err)
	}
	if len(store.historyUpdates) != 0 {
		t.Fatalf("expected no assistant history update for empty stream, got %#v", store.historyUpdates)
	}
}

func TestRespondToConversation_ShellRequestsReachModelWithWorkspaceTools(t *testing.T) {
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "我会在 workspace 中执行命令。"},
		}},
	}}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.AllowedTools = []string{"bash"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
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
	if len(llm.lastReq.Tools) != 2 || !toolNamesInclude(llm.lastReq.Tools, "bash") || !toolNamesInclude(llm.lastReq.Tools, "read_persisted_output") {
		t.Fatalf("expected shell request to expose bash tool plus read_persisted_output, got %#v", llm.lastReq.Tools)
	}
	systemPrompt := llm.lastReq.Messages[0].Content
	if !contains(systemPrompt, "Working directory: "+cfg.WorkspaceRoot) {
		t.Fatalf("expected system prompt to include workspace root, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "而不是只能聊天的助手") {
		t.Fatalf("expected system prompt to position the model as a full agent, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "本次会话可用的工具如下：\n\n- bash") {
		t.Fatalf("expected system prompt to include available tools, got %q", systemPrompt)
	}
	if len(store.historyUpdates) != 1 {
		t.Fatalf("expected one full history update, got %d", len(store.historyUpdates))
	}
	if got := store.historyUpdates[0]; len(got) != 2 || got[0].Content != "请帮我执行 shell 命令 ls 查看本地目录" || got[1].Content != "我会在 workspace 中执行命令。" {
		t.Fatalf("expected full shell conversation history, got %#v", got)
	}
}

func TestRespondToConversation_IncludesLocalSkillsInPrompt(t *testing.T) {
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "我会结合这些能力来回答。"},
		}},
	}}}
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {
			Meta: map[string]string{"description": "Builtin description"},
			Body: "builtin body",
			Path: "builtin://builtin-skill",
		},
		"user-skill": {
			Meta:   map[string]string{"description": "User description"},
			Body:   "user body",
			Source: "user",
			Path:   "/home/.cynosure/skills/user-skill/SKILL.md",
		},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
	conversation := storage.Conversation{ID: "conv_3", Title: "新对话"}
	user := storage.User{ID: "usr_3", Username: "carol"}

	_, err := service.RespondToConversation(context.Background(), conversation, user, "请结合现有技能回答", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if llm.calls != 1 {
		t.Fatalf("expected llm to be called once, got %d", llm.calls)
	}
	if len(llm.lastReq.Tools) != 2 || !toolNamesInclude(llm.lastReq.Tools, "load_skill") {
		t.Fatalf("expected load_skill tool to be exposed, got %d tools", len(llm.lastReq.Tools))
	}
	systemPrompt := llm.lastReq.Messages[0].Content
	if !contains(systemPrompt, "Working directory: "+cfg.WorkspaceRoot) {
		t.Fatalf("expected workspace root in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "本次会话可用的工具如下：\n\n- load_skill") {
		t.Fatalf("expected load_skill to be listed in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "<name>builtin-skill</name>\n<description>Builtin description</description>") {
		t.Fatalf("expected builtin skill description in prompt, got %q", systemPrompt)
	}
	if !contains(systemPrompt, "<name>user-skill</name>\n<description>User description</description>") {
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
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "继续回答"},
		}},
	}}}
	store := &fakeStore{cached: []storage.Message{
		{ID: "msg_old_user", ConversationID: "conv_history", UserID: "usr_history", Role: "user", Content: "先加载技能"},
		{ID: "msg_old_assistant_tool", ConversationID: "conv_history", UserID: "usr_history", Role: "assistant", ToolCalls: []storage.MessageToolCall{{ID: "tool_1", Type: "function", Function: storage.MessageFunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}}}},
		{ID: "msg_old_tool", ConversationID: "conv_history", UserID: "usr_history", Role: "tool", ToolCallID: "tool_1", Content: `{"status":"success","result":"loaded"}`},
		{ID: "msg_old_assistant", ConversationID: "conv_history", UserID: "usr_history", Role: "assistant", Content: "已经加载"},
	}}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}

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
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{{
		Choices: []openai.ChatCompletionChoice{{
			FinishReason: openai.FinishReasonStop,
			Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "好的，我继续回答。"},
		}},
	}}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
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

	prompt := service.buildSystemPromptWithMemory(storage.User{ID: "usr_6", Username: "frank"}, agenttools.NewSkillSnapshot(nil, loader), "")
	if !contains(prompt, "Working directory: "+cfg.WorkspaceRoot) {
		t.Fatalf("expected prompt to include workspace root, got %q", prompt)
	}
	if !contains(prompt, "而不是只能聊天的助手") {
		t.Fatalf("expected prompt to describe a full agent, got %q", prompt)
	}
	if contains(prompt, "本次会话可用的工具如下：") {
		t.Fatalf("expected prompt without tool registry to omit tool list, got %q", prompt)
	}
	if !contains(prompt, "<name>builtin-skill</name>\n<description>Builtin description</description>") {
		t.Fatalf("expected prompt to include skill descriptions, got %q", prompt)
	}
}

func TestBuildSystemPromptIncludesCynosureMarkdownContext(t *testing.T) {
	cfg := testAppConfig(t)
	service := &Service{Cfg: cfg}
	service.SetCynosureMarkdownContext(config.CynosureMarkdownContext{
		UserPath:         "/home/alice/.cynosure/CYNOSURE.MD",
		UserContent:      "# User Rule\n全局说明",
		WorkspacePath:    filepath.Join(cfg.WorkspaceRoot, ".cynosure", "CYNOSURE.MD"),
		WorkspaceContent: "# Project Rule\n项目说明",
	})

	prompt := service.buildSystemPromptWithMemory(storage.User{ID: "usr_link", Username: "link-user"}, nil, "")
	for _, want := range []string{
		"<system-reminder>",
		"# cynosureMd",
		"/home/alice/.cynosure/CYNOSURE.MD 的内容（用户为所有项目配置的私人全局说明）：",
		"# User Rule\n全局说明",
		filepath.Join(cfg.WorkspaceRoot, ".cynosure", "CYNOSURE.MD") + " 的内容（项目说明，已提交到代码库或工作区）：",
		"# Project Rule\n项目说明",
	} {
		if !contains(prompt, want) {
			t.Fatalf("expected prompt to contain %q, got %q", want, prompt)
		}
	}
}

func TestBuildSkillSnapshotUsesLocalSkillsWithoutStoreLookup(t *testing.T) {
	localSkills := sessions.NewSkillLoader()
	localSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"local-skill": {Meta: map[string]string{"description": "Local description"}, Body: "local body", Source: "workspace", Path: "/project/.cynosure/skills/local-skill/SKILL.md"},
	})
	service := &Service{Cfg: testAppConfig(t), BuiltinSkills: localSkills}

	snapshot, err := service.buildSkillSnapshot(context.Background(), "usr_refresh")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	loaded, err := snapshot.LoadSkill("local-skill")
	if err != nil {
		t.Fatalf("load local skill: %v", err)
	}
	if loaded.Source != "workspace" || loaded.Entry.Body != "local body" {
		t.Fatalf("loaded = source %q body %q, want workspace local body", loaded.Source, loaded.Entry.Body)
	}
}

func TestSkillSnapshotLoadSkillPrefersWorkspaceSkillThenFallsBackToUser(t *testing.T) {
	userSkills := sessions.NewSkillLoader()
	userSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"shared-skill": {Meta: map[string]string{"description": "User description", "tags": "user"}, Body: "user body", Source: "user", Path: "/home/.cynosure/skills/shared/skill.md"},
		"user-skill":   {Meta: map[string]string{"description": "User only"}, Body: "user only body", Source: "user", Path: "/home/.cynosure/skills/user/skill.md"},
	})
	workspaceSkills := sessions.NewSkillLoader()
	workspaceSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"shared-skill":    {Meta: map[string]string{"description": "Workspace description"}, Body: "workspace body", Source: "workspace", Path: "/project/.cynosure/skills/shared/skill.md"},
		"workspace-skill": {Meta: map[string]string{"description": "Workspace only"}, Body: "workspace only body", Source: "workspace", Path: "/project/.cynosure/skills/workspace/skill.md"},
	})
	snapshot := agenttools.NewSkillSnapshot(userSkills, workspaceSkills)

	loaded, err := snapshot.LoadSkill("shared-skill")
	if err != nil {
		t.Fatalf("load shared skill: %v", err)
	}
	if loaded.Source != "workspace" || loaded.Entry.Body != "workspace body" {
		t.Fatalf("expected workspace skill to win, got source=%q body=%q", loaded.Source, loaded.Entry.Body)
	}

	loaded, err = snapshot.LoadSkill("user-skill")
	if err != nil {
		t.Fatalf("load user skill: %v", err)
	}
	if loaded.Source != "user" || loaded.Entry.Body != "user only body" {
		t.Fatalf("expected user fallback, got source=%q body=%q", loaded.Source, loaded.Entry.Body)
	}
}

func TestSkillSnapshotLoadSkillUsesMergedWorkspaceOverride(t *testing.T) {
	userSkills := sessions.NewSkillLoader()
	userSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"shared-skill": {Meta: map[string]string{"description": "User description"}, Body: "user body", Source: "user", Path: "/home/.cynosure/skills/shared/skill.md"},
	})
	workspaceSkills := sessions.NewSkillLoader()
	workspaceSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"shared-skill": {Meta: map[string]string{"description": "Workspace description"}, Body: "workspace body", Source: "workspace", Path: "/project/.cynosure/skills/shared/skill.md"},
	})

	snapshot := agenttools.NewSkillSnapshot(userSkills, workspaceSkills)
	loaded, err := snapshot.LoadSkill("shared-skill")
	if err != nil {
		t.Fatalf("LoadSkill returned error: %v", err)
	}
	if loaded.Entry.Body != "workspace body" {
		t.Fatalf("LoadSkill body = %q, want workspace body", loaded.Entry.Body)
	}
	if loaded.Source != "workspace" {
		t.Fatalf("LoadSkill source = %q, want workspace", loaded.Source)
	}
}

func TestRespondToConversation_ReturnsSuccessfulToolResultIntoLoop(t *testing.T) {
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
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
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

func TestRespondToConversation_EmitsToolLifecycleEvents(t *testing.T) {
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonToolCalls,
				Message: openai.ChatCompletionMessage{Role: "assistant", ToolCalls: []openai.ToolCall{{
					ID:   "tool_lifecycle_1",
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
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
	writer := &captureEventWriter{}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_tool_events", Title: "新对话"}, storage.User{ID: "usr_tool_events", Username: "tool-user"}, "加载技能", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := writer.nonMetaEvents()
	if len(events) < 3 {
		t.Fatalf("events = %#v, want at least tool start, tool done, assistant", events)
	}
	if events[0].name != toolCallStartEvent || events[1].name != toolCallDoneEvent {
		t.Fatalf("event order = %#v, want tool_call_start then tool_call_done before assistant streaming", events)
	}
	start := eventMap(t, events[0].data)
	if start["tool_call_id"] != "tool_lifecycle_1" || start["tool_name"] != "load_skill" || start["status"] != "running" {
		t.Fatalf("start event = %#v, want running load_skill", start)
	}
	if !strings.Contains(stringValue(start["args_preview"]), "builtin-skill") {
		t.Fatalf("start args_preview = %#v, want skill name", start["args_preview"])
	}
	done := eventMap(t, events[1].data)
	if done["tool_call_id"] != "tool_lifecycle_1" || done["tool_name"] != "load_skill" || done["status"] != "success" {
		t.Fatalf("done event = %#v, want success load_skill", done)
	}
	if !strings.Contains(stringValue(done["result_preview"]), "builtin-skill") {
		t.Fatalf("done result_preview = %#v, want loaded skill preview", done["result_preview"])
	}
}

func TestRespondToConversation_EmitsRejectedToolLifecycleEvent(t *testing.T) {
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{
			Choices: []openai.ChatCompletionChoice{{
				FinishReason: openai.FinishReasonToolCalls,
				Message: openai.ChatCompletionMessage{Role: "assistant", ToolCalls: []openai.ToolCall{{
					ID:   "tool_lifecycle_2",
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
				Message:      openai.ChatCompletionMessage{Role: "assistant", Content: "handled"},
			}},
		},
	}}
	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: sessions.NewSkillLoader()}
	writer := &captureEventWriter{}

	_, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_tool_rejected", Title: "新对话"}, storage.User{ID: "usr_tool_rejected", Username: "tool-user"}, "加载不存在技能", writer)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := writer.nonMetaEvents()
	if len(events) < 2 || events[0].name != toolCallStartEvent || events[1].name != toolCallDoneEvent {
		t.Fatalf("events = %#v, want tool lifecycle events", events)
	}
	done := eventMap(t, events[1].data)
	if done["status"] != "rejected" {
		t.Fatalf("done event status = %#v, want rejected", done["status"])
	}
	if !strings.Contains(stringValue(done["result_preview"]), "missing-skill") {
		t.Fatalf("done result_preview = %#v, want rejection reason", done["result_preview"])
	}
}

func TestRespondToConversation_ReturnsRejectedToolResultIntoLoop(t *testing.T) {
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
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	store := &fakeStore{}
	cfg := testAppConfig(t)
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}
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
	originalBash := agenttools.Handlers["bash"]
	defer func() {
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
	cfg.AllowedTools = []string{"load_skill", "bash"}
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: builtin}

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

func TestRespondToConversation_LocalSkillUsesWorkspaceWithinMultiToolTurn(t *testing.T) {
	originalBash := agenttools.Handlers["bash"]
	defer func() {
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
					{ID: "tool_load", Type: "function", Function: openai.FunctionCall{Name: "load_skill", Arguments: `{"name":"local-skill"}`}},
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
	workspace := t.TempDir()
	store := &fakeStore{}
	cfg := testAppConfig(t)
	cfg.WorkspaceRoot = workspace
	cfg.AllowedTools = []string{"load_skill", "bash"}
	localSkills := sessions.NewSkillLoader()
	localSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"local-skill": {Meta: map[string]string{"description": "local description"}, Body: "local body", Source: "workspace", Path: filepath.Join(workspace, ".cynosure", "skills", "local-skill", "SKILL.md")},
	})
	service := &Service{LLM: llm, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg), BuiltinSkills: localSkills}

	message, err := service.RespondToConversation(context.Background(), storage.Conversation{ID: "conv_local", Title: "新对话"}, storage.User{ID: "usr_local", Username: "local-user"}, "请加载本地 skill 后执行 bash", nil)
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
	if len(defs) != 2 || !toolNamesInclude(defs, "read_persisted_output") {
		t.Fatalf("expected configured tool plus auto read_persisted_output, got %#v", defs)
	}
	expected, ok := lookupRegisteredTool("load_skill")
	if !ok {
		t.Fatalf("expected load_skill to exist in registered tool definitions")
	}
	loadSkill := findToolDef(t, defs, "load_skill")
	if loadSkill.Function == nil || expected.Function == nil {
		t.Fatalf("expected function definitions to be present")
	}
	if loadSkill.Function.Description != expected.Function.Description {
		t.Fatalf("expected tool description %q, got %q", expected.Function.Description, loadSkill.Function.Description)
	}
}

func TestToolRegistryDefinitions_AreLoadedAtRegistryCreation(t *testing.T) {
	cfg := config.AppConfig{AllowedTools: []string{"bash"}}
	registry := NewToolRegistry(cfg)
	cfg.AllowedTools = []string{"read_file"}

	defs := registry.Definitions()
	if len(defs) != 2 || !toolNamesInclude(defs, "bash") {
		t.Fatalf("expected registry to keep startup-loaded bash definition, got %#v", defs)
	}
}

func findToolDef(t *testing.T, tools []openai.Tool, name string) openai.Tool {
	t.Helper()
	for _, tool := range tools {
		if tool.Function != nil && tool.Function.Name == name {
			return tool
		}
	}
	t.Fatalf("expected tool %q in definitions", name)
	return openai.Tool{}
}

func TestToolRegistryExecute_LoadSkillReturnsSkillContentWithoutRuntimePaths(t *testing.T) {
	localSkills := sessions.NewSkillLoader()
	localSkills.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin description", "tags": "go,agent"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	snapshot := agenttools.NewSkillSnapshot(nil, localSkills)
	registry := NewToolRegistry(config.AppConfig{
		AppHome:       "/deploy/app",
		WorkspaceRoot: "/deploy/app/workspace",
	})

	result, err := registry.Execute(context.Background(), ToolContext{Skills: snapshot}, "load_skill", `{"name":"builtin-skill"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	content := result.Output
	if !contains(content, "<skill source=\"local\" name=\"builtin-skill\">") || !contains(content, "builtin body") {
		t.Fatalf("expected loaded skill content, got %q", content)
	}
	if !contains(content, "<metadata>") || !contains(content, "description: Builtin description") || !contains(content, "tags: go,agent") {
		t.Fatalf("expected loaded skill to include non-path metadata, got %q", content)
	}
	if contains(content, "path: builtin://builtin-skill") || contains(content, "<runtime-paths>") {
		t.Fatalf("expected loaded skill output to omit local paths, got %q", content)
	}
	if contains(content, "/deploy/app") || contains(content, "COMMAND_BIN_DIR") || contains(content, "COMMAND_SCRIPT_DIR") || contains(content, "WORKSPACE_ROOT") {
		t.Fatalf("expected loaded skill output to omit runtime path details, got %q", content)
	}
	if _, ok := agenttools.Handlers["load_skill"]; !ok {
		t.Fatalf("expected load_skill handler to remain registered in shared tool registry")
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

	registry := NewToolRegistry(config.AppConfig{WorkspaceRoot: "/workspaces/usr_1", AllowedTools: []string{"bash"}})
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
	registry := NewToolRegistry(config.AppConfig{WorkspaceRoot: workspace, AllowedTools: []string{"bash"}})

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
		return env.WorkspaceRoot + "|" + env.CurrentWorkingDir, nil
	}

	workspace := t.TempDir()
	cfg := config.AppConfig{
		WorkspaceRoot: workspace,
		AllowedTools:  []string{"bash"},
	}
	registry := NewToolRegistry(cfg)

	result, err := registry.Execute(context.Background(), ToolContext{}, "bash", `{"command":"pwd"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := filepath.Clean(workspace) + "|" + filepath.Clean(workspace)
	if result.Output != expected {
		t.Fatalf("expected workspace-derived paths to remain stable, got %q", result.Output)
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

	cfg := config.AppConfig{WorkspaceRoot: workspace, AllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+command+`"}`)

	if outcome.Status != "success" {
		t.Fatalf("expected success, got %#v", outcome)
	}
	if filepath.Clean(outcome.Audit.ResolvedCWD) != filepath.Clean(workspace) {
		t.Fatalf("expected audit cwd %q, got %#v", workspace, outcome.Audit)
	}
}

func TestRegisteredTools_UsesCurrentAllowList(t *testing.T) {
	tools := RegisteredTools(config.AppConfig{})
	if len(tools) != 1 || tools[0] != "load_skill" {
		t.Fatalf("expected only load_skill to be registered for local runtime, got %v", tools)
	}
}

func TestRegisteredTools_UsesConfiguredAllowList(t *testing.T) {
	tools := RegisteredTools(config.AppConfig{AllowedTools: []string{"load_skill", "missing_tool", "load_skill"}})
	if len(tools) != 1 || tools[0] != "load_skill" {
		t.Fatalf("expected configured allowlist to be filtered to registered tools, got %v", tools)
	}
}

func TestExecuteToolCall_AuditCapturesResolvedCommandPath(t *testing.T) {
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create bin dir: %v", err)
	}
	script := filepath.Join(binDir, "helper")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho ok\n"), 0o755); err != nil {
		t.Fatalf("write helper script: %v", err)
	}

	cfg := config.AppConfig{WorkspaceRoot: workspace, AllowedTools: []string{"bash"}}
	service := &Service{Cfg: cfg, Tools: NewToolRegistry(cfg)}
	outcome := executeToolCallWithDefaultHooks(t, service, "bash", `{"command":"`+script+` --flag"}`)

	if outcome.Audit.ResolvedCommandPath != script {
		t.Fatalf("expected audit resolved command path %q, got %#v", script, outcome.Audit)
	}
}

func TestExecuteToolCall_AuditCapturesAppWorkspaceResolution(t *testing.T) {
	workspace := t.TempDir()
	binDir := filepath.Join(workspace, "bin")
	if err := os.MkdirAll(binDir, 0o755); err != nil {
		t.Fatalf("create workspace bin dir: %v", err)
	}
	command := filepath.Join(binDir, "helper")
	if err := os.WriteFile(command, []byte("#!/bin/sh\npwd\n"), 0o755); err != nil {
		t.Fatalf("write workspace helper: %v", err)
	}

	cfg := config.AppConfig{
		WorkspaceRoot: workspace,
		AllowedTools:  []string{"bash"},
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
	if outcome.Audit.ResolvedCommandPath != command {
		t.Fatalf("expected audit resolved command path %q, got %#v", command, outcome.Audit)
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
		State:    service.newLoopState(storage.Conversation{ID: "conv_audit"}, storage.User{ID: "usr_audit"}, "", nil, nil, nil),
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

func bashToolRound(id string) []openai.ChatCompletionStreamResponse {
	return []openai.ChatCompletionStreamResponse{{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: id, Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "bash", Arguments: `{"command":"pwd"}`}}}}, FinishReason: openai.FinishReasonToolCalls}}}}
}

func todoWriteToolRound(id string) []openai.ChatCompletionStreamResponse {
	return []openai.ChatCompletionStreamResponse{{Choices: []openai.ChatCompletionStreamChoice{{Delta: openai.ChatCompletionStreamChoiceDelta{ToolCalls: []openai.ToolCall{{Index: intPointer(0), ID: id, Type: openai.ToolTypeFunction, Function: openai.FunctionCall{Name: "todo_write", Arguments: `{"todos":[{"id":"1","content":"更新计划","status":"in_progress"}]}`}}}}, FinishReason: openai.FinishReasonToolCalls}}}}
}

func requestContainsTodoWriteReminder(req openai.ChatCompletionRequest) bool {
	for _, message := range req.Messages {
		if strings.Contains(message.Content, "not called todo_write for 3 consecutive model rounds") {
			return true
		}
	}
	return false
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

func toolNamesInclude(tools []openai.Tool, name string) bool {
	for _, tool := range tools {
		if tool.Function != nil && tool.Function.Name == name {
			return true
		}
	}
	return false
}

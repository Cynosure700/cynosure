package compression

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"nano_cc/internal/textutil"
	"nano_cc/internal/web/storage"
)

type fakeStore struct {
	persistedOutputs     []storage.PersistedOutput
	contextSummaries     []storage.ContextSummary
	conversationMemories []storage.ConversationMemory
}

func (f *fakeStore) CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error {
	f.persistedOutputs = append(f.persistedOutputs, output)
	return nil
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

func (f *fakeStore) ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error) {
	var result []storage.ConversationMemory
	for _, m := range f.conversationMemories {
		if m.ConversationID == conversationID {
			result = append(result, m)
		}
	}
	return result, nil
}

func toolMsg(id, status, result string) storage.Message {
	content, _ := json.Marshal(toolResultMessageContent{Status: status, Result: result})
	return storage.Message{ID: "m_" + id, Role: "tool", ToolCallID: id, Content: string(content)}
}

func assistantToolCallMsg(callID string) storage.Message {
	return storage.Message{Role: "assistant", ToolCalls: []storage.MessageToolCall{{ID: callID, Type: "function", Function: storage.MessageFunctionCall{Name: "bash", Arguments: "{}"}}}}
}

func resultOf(t *testing.T, content string) string {
	t.Helper()
	_, result, _ := textutil.ParseToolResult(content)
	return result
}

// --- ToolResultCompressionStrategy ---

func TestToolResultCompression_NoopUnderThreshold(t *testing.T) {
	store := &fakeStore{}
	history := []storage.Message{
		{Role: "user", Content: "hi"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", strings.Repeat("a", 1024)),
	}
	req := &Request{User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"}, RequestHistory: history, Store: store}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 0 {
		t.Fatalf("expected no persisted outputs, got %d", len(store.persistedOutputs))
	}
}

func TestToolResultCompression_PersistsLargestUntilUnderThreshold(t *testing.T) {
	store := &fakeStore{}
	big := strings.Repeat("b", 250*1024)
	small := strings.Repeat("s", 1024)
	history := []storage.Message{
		{Role: "user", Content: "go"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", big),
		assistantToolCallMsg("c2"),
		toolMsg("c2", "success", small),
	}
	req := &Request{User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"}, RequestHistory: history, Store: store}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 1 {
		t.Fatalf("expected 1 persisted output, got %d", len(store.persistedOutputs))
	}
	if store.persistedOutputs[0].Content != big {
		t.Fatalf("expected full big content persisted")
	}
	// big result now a marker; small stays full.
	if !strings.Contains(resultOf(t, history[2].Content), PersistedOutputMarkerPrefix) {
		t.Fatalf("expected big tool result replaced with marker, got %q", history[2].Content)
	}
	if resultOf(t, history[4].Content) != small {
		t.Fatalf("expected small tool result to stay inline")
	}
}

func TestToolResultCompression_PreservesJSONAndPreview(t *testing.T) {
	store := &fakeStore{}
	big := strings.Repeat("x", 300*1024)
	history := []storage.Message{
		{Role: "user", Content: "go"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", big),
	}
	req := &Request{User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"}, RequestHistory: history, Store: store}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	status, result, isJSON := textutil.ParseToolResult(history[2].Content)
	if !isJSON || status != "success" {
		t.Fatalf("expected JSON wrapper preserved with status, got %q", history[2].Content)
	}
	preview := strings.Repeat("x", toolResultPreviewRunes)
	if !strings.Contains(result, preview) {
		t.Fatalf("expected 2000-char preview embedded")
	}
}

func TestToolResultCompression_SkipsAlreadyMarked(t *testing.T) {
	store := &fakeStore{}
	marker := buildPersistedOutputMarker("po_existing", 300*1024, "preview")
	history := []storage.Message{
		{Role: "user", Content: "go"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", marker),
	}
	req := &Request{User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"}, RequestHistory: history, Store: store}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 0 {
		t.Fatalf("expected already-marked result to be skipped")
	}
}

func TestToolResultCompression_IgnoresToolsBeforeLastUser(t *testing.T) {
	store := &fakeStore{}
	big := strings.Repeat("b", 300*1024)
	history := []storage.Message{
		assistantToolCallMsg("old"),
		toolMsg("old", "success", big),
		{Role: "user", Content: "new turn"},
	}
	req := &Request{User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"}, RequestHistory: history, Store: store}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 0 {
		t.Fatalf("expected tool results before last user to be ignored")
	}
}

// --- MessageWindowCompressionStrategy ---

func TestMessageWindow_NoopAtLimit(t *testing.T) {
	history := make([]storage.Message, messageWindowLimit)
	for i := range history {
		history[i] = storage.Message{Role: "user", Content: "m"}
	}
	req := &Request{RequestHistory: history}
	if err := (&MessageWindowCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(req.RequestHistory) != messageWindowLimit {
		t.Fatalf("expected no trim at limit, got %d", len(req.RequestHistory))
	}
}

func TestMessageWindow_TrimsMiddleKeepingHeadTail(t *testing.T) {
	const n = 60
	history := make([]storage.Message, n)
	for i := range history {
		history[i] = storage.Message{Role: "user", Content: string(rune('A' + i%26))}
	}
	first := history[0].Content
	lastTailFirst := history[n-messageWindowTail].Content
	req := &Request{RequestHistory: history}
	if err := (&MessageWindowCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	got := req.RequestHistory
	if len(got) != messageWindowHead+messageWindowTail {
		t.Fatalf("expected %d messages, got %d", messageWindowHead+messageWindowTail, len(got))
	}
	if got[0].Content != first {
		t.Fatalf("expected head preserved")
	}
	if got[messageWindowHead].Content != lastTailFirst {
		t.Fatalf("expected tail to start at original index %d", n-messageWindowTail)
	}
}

func TestMessageWindow_DropsOrphanToolAtCut(t *testing.T) {
	// Build 60 messages; ensure the tail boundary starts with an orphan tool.
	history := make([]storage.Message, 0, 60)
	history = append(history, storage.Message{Role: "user", Content: "u0"})
	history = append(history, storage.Message{Role: "user", Content: "u1"})
	history = append(history, storage.Message{Role: "user", Content: "u2"})
	for i := 3; i < 60; i++ {
		history = append(history, storage.Message{Role: "user", Content: "filler"})
	}
	// Make the first tail message an orphan tool (no preceding assistant call in window).
	history[60-messageWindowTail] = storage.Message{Role: "tool", ToolCallID: "orphan", Content: `{"status":"success","result":"x"}`}
	req := &Request{RequestHistory: history}
	if err := (&MessageWindowCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	for _, msg := range req.RequestHistory {
		if msg.Role == "tool" && msg.ToolCallID == "orphan" {
			t.Fatalf("expected orphan tool message to be dropped")
		}
	}
}

// --- RecentToolResultRetentionStrategy ---

func TestRecentToolRetention_NoopAtOrBelowThree(t *testing.T) {
	history := []storage.Message{
		toolMsg("a", "success", "ra"),
		toolMsg("b", "success", "rb"),
		toolMsg("c", "success", "rc"),
	}
	req := &Request{RequestHistory: history}
	if err := (&RecentToolResultRetentionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if resultOf(t, history[0].Content) != "ra" {
		t.Fatalf("expected no compaction at threshold")
	}
}

func TestRecentToolRetention_CompactsOlderKeepsRecentThree(t *testing.T) {
	history := []storage.Message{
		toolMsg("a", "success", "ra"),
		toolMsg("b", "success", "rb"),
		toolMsg("c", "success", "rc"),
		toolMsg("d", "success", "rd"),
		toolMsg("e", "success", "re"),
	}
	req := &Request{RequestHistory: history}
	if err := (&RecentToolResultRetentionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if resultOf(t, history[0].Content) != earlierToolResultPlaceholder {
		t.Fatalf("expected oldest tool compacted, got %q", history[0].Content)
	}
	if resultOf(t, history[1].Content) != earlierToolResultPlaceholder {
		t.Fatalf("expected second oldest tool compacted")
	}
	// last three kept full
	for _, idx := range []int{2, 3, 4} {
		if isCompactedResult(resultOf(t, history[idx].Content)) {
			t.Fatalf("expected recent tool %d to keep full result", idx)
		}
	}
	// JSON status preserved
	status, _, isJSON := textutil.ParseToolResult(history[0].Content)
	if !isJSON || status != "success" {
		t.Fatalf("expected JSON wrapper preserved when placeholdering")
	}
}

func TestRecentToolRetention_NonJSONReplacedWithPlainPlaceholder(t *testing.T) {
	history := []storage.Message{
		{Role: "tool", ToolCallID: "a", Content: "raw not json"},
		toolMsg("b", "success", "rb"),
		toolMsg("c", "success", "rc"),
		toolMsg("d", "success", "rd"),
	}
	req := &Request{RequestHistory: history}
	if err := (&RecentToolResultRetentionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if history[0].Content != earlierToolResultPlaceholder {
		t.Fatalf("expected non-JSON tool content replaced with plain placeholder, got %q", history[0].Content)
	}
}

// --- FullHistorySummarizationStrategy ---

func TestFullHistorySummarization_NoopUnderBudget(t *testing.T) {
	store := &fakeStore{}
	called := false
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: []storage.Message{{Role: "user", Content: "tiny"}},
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer: func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
			called = true
			return SummaryResult{Summary: "s"}, nil
		},
	}
	if err := (&FullHistorySummarizationStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if called {
		t.Fatalf("expected summarizer not called under budget")
	}
}

func TestFullHistorySummarization_SummarizesOverBudgetAndCaches(t *testing.T) {
	store := &fakeStore{}
	calls := 0
	var lastInput []storage.Message
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		calls++
		lastInput = r.History
		return SummaryResult{Summary: "STRUCTURED SUMMARY"}, nil
	}
	huge := storage.Message{Role: "user", Content: strings.Repeat("z", 600*1024)}
	last := storage.Message{Role: "user", Content: "the latest message"}
	mkReq := func() *Request {
		return &Request{
			User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
			RequestHistory: []storage.Message{huge, last},
			Store:          store,
			Estimator:      DefaultTokenEstimator{},
			Summarizer:     summarizer,
		}
	}
	req := mkReq()
	if err := (&FullHistorySummarizationStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected summarizer called once, got %d", calls)
	}
	// Summarizer must only receive the history before the last message.
	if len(lastInput) != 1 || lastInput[0].Content != huge.Content {
		t.Fatalf("expected summarizer input to exclude the last message, got %#v", lastInput)
	}
	if len(req.RequestHistory) < 2 || req.RequestHistory[1].Role != "user" || !strings.Contains(req.RequestHistory[1].Content, "STRUCTURED SUMMARY") {
		t.Fatalf("expected summary injected as request context, got %#v", req.RequestHistory)
	}
	// The last message must be preserved verbatim as the final entry.
	final := req.RequestHistory[len(req.RequestHistory)-1]
	if final.Content != last.Content {
		t.Fatalf("expected last message preserved verbatim, got %#v", final)
	}
	if len(store.contextSummaries) != 1 {
		t.Fatalf("expected summary cached once")
	}
	// Second run with same source hash should hit cache (no extra summarizer call).
	req2 := mkReq()
	if err := (&FullHistorySummarizationStrategy{}).Apply(context.Background(), req2); err != nil {
		t.Fatalf("apply2: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected cache hit to avoid second summarizer call, got %d", calls)
	}
}

func TestFullHistorySummarization_NoopWithSingleMessage(t *testing.T) {
	store := &fakeStore{}
	called := false
	huge := storage.Message{Role: "user", Content: strings.Repeat("z", 600*1024)}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: []storage.Message{huge},
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer: func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
			called = true
			return SummaryResult{Summary: "s"}, nil
		},
	}
	if err := (&FullHistorySummarizationStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if called {
		t.Fatalf("expected no summarization with a single message")
	}
	if len(req.RequestHistory) != 1 || req.RequestHistory[0].Content != huge.Content {
		t.Fatalf("expected single message left untouched, got %#v", req.RequestHistory)
	}
}

func TestFullHistorySummarization_PreservesLastToolMessageWithPairedCall(t *testing.T) {
	store := &fakeStore{}
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		return SummaryResult{Summary: "S"}, nil
	}
	// Earlier huge history forces summarization; the latest turn is an
	// assistant tool_call followed by a large tool result as the last message.
	huge := storage.Message{Role: "user", Content: strings.Repeat("z", 600*1024)}
	callMsg := assistantToolCallMsg("call_latest")
	lastTool := toolMsg("call_latest", "success", strings.Repeat("r", 600*1024))
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: []storage.Message{huge, callMsg, lastTool},
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	if err := (&FullHistorySummarizationStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	final := req.RequestHistory[len(req.RequestHistory)-1]
	if final.Role != "tool" || final.ToolCallID != "call_latest" {
		t.Fatalf("expected last tool message preserved, got %#v", final)
	}
	// Its paired assistant tool_call must also be present to keep the request valid.
	foundCall := false
	for _, msg := range req.RequestHistory {
		if msg.Role == "assistant" {
			for _, c := range msg.ToolCalls {
				if c.ID == "call_latest" {
					foundCall = true
				}
			}
		}
	}
	if !foundCall {
		t.Fatalf("expected paired assistant tool_call kept for the last tool message")
	}
}

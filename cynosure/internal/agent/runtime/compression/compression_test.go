package compression

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/textutil"
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

func TestToolResultCompression_NoopUnderDefaultToolLimit(t *testing.T) {
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

func TestToolResultCompression_PersistsOnlyResultOverDefaultToolLimit(t *testing.T) {
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

func TestToolResultCompression_PersistsSingleResultOverToolLimit(t *testing.T) {
	store := &fakeStore{}
	history := []storage.Message{
		{Role: "user", Content: "go"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", "12345678901"),
	}
	req := &Request{
		User:                   storage.User{ID: "u"},
		Conversation:           storage.Conversation{ID: "c"},
		RequestHistory:         history,
		Store:                  store,
		ToolMaxResultSizeChars: func(toolName string) int { return 10 },
	}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 1 {
		t.Fatalf("expected 1 persisted output, got %d", len(store.persistedOutputs))
	}
	if !strings.Contains(resultOf(t, history[2].Content), PersistedOutputMarkerPrefix) {
		t.Fatalf("expected over-limit result replaced with marker, got %q", history[2].Content)
	}
}

func TestToolResultCompression_DoesNotCompressByCombinedTotal(t *testing.T) {
	store := &fakeStore{}
	first := strings.Repeat("a", 120*1024)
	second := strings.Repeat("b", 120*1024)
	history := []storage.Message{
		{Role: "user", Content: "go"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", first),
		assistantToolCallMsg("c2"),
		toolMsg("c2", "success", second),
	}
	req := &Request{
		User:                   storage.User{ID: "u"},
		Conversation:           storage.Conversation{ID: "c"},
		RequestHistory:         history,
		Store:                  store,
		ToolMaxResultSizeChars: func(toolName string) int { return 200 * 1024 },
	}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 0 {
		t.Fatalf("expected no persisted outputs, got %d", len(store.persistedOutputs))
	}
	if resultOf(t, history[2].Content) != first || resultOf(t, history[4].Content) != second {
		t.Fatalf("expected both under-limit results to stay inline")
	}
}

func TestToolResultCompression_UsesDefaultLimitWhenToolNameMissing(t *testing.T) {
	store := &fakeStore{}
	history := []storage.Message{
		{Role: "user", Content: "go"},
		toolMsg("missing", "success", strings.Repeat("x", 50001)),
	}
	req := &Request{User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"}, RequestHistory: history, Store: store}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 1 {
		t.Fatalf("expected default limit to persist over-limit result, got %d persisted outputs", len(store.persistedOutputs))
	}
}

func TestToolResultCompression_CountsRunesForLimit(t *testing.T) {
	store := &fakeStore{}
	result := strings.Repeat("界", 10)
	history := []storage.Message{
		{Role: "user", Content: "go"},
		assistantToolCallMsg("c1"),
		toolMsg("c1", "success", result),
	}
	req := &Request{
		User:                   storage.User{ID: "u"},
		Conversation:           storage.Conversation{ID: "c"},
		RequestHistory:         history,
		Store:                  store,
		ToolMaxResultSizeChars: func(toolName string) int { return 10 },
	}
	if err := (&ToolResultCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(store.persistedOutputs) != 0 {
		t.Fatalf("expected 10-rune result at limit to stay inline, got %d persisted outputs", len(store.persistedOutputs))
	}
	if resultOf(t, history[2].Content) != result {
		t.Fatalf("expected CJK result to stay inline at rune limit")
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
	// Latest user at index 0 followed by many assistant messages, so the
	// current turn count (60) exceeds the limit and trimming triggers.
	const n = 60
	history := make([]storage.Message, n)
	history[0] = storage.Message{Role: "user", Content: "latest-user"}
	for i := 1; i < n; i++ {
		history[i] = storage.Message{Role: "assistant", Content: string(rune('A' + i%26))}
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
		t.Fatalf("expected head to start at latest user")
	}
	if got[messageWindowHead].Content != lastTailFirst {
		t.Fatalf("expected tail to start at original index %d", n-messageWindowTail)
	}
}

func TestMessageWindow_KeepsLatestUserPlusTwoAsHead(t *testing.T) {
	// 60 messages, latest user at index 5 with only non-user messages after it.
	// Current turn count is 60-5=55 > 50, so trimming triggers and head keeps
	// "latest user + 2" while tail keeps the most recent 46.
	const n = 60
	history := make([]storage.Message, n)
	for i := range history {
		history[i] = storage.Message{Role: "assistant", Content: string(rune('A' + i%26))}
	}
	history[5] = storage.Message{Role: "user", Content: "latest-user"}
	req := &Request{RequestHistory: history}
	if err := (&MessageWindowCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	got := req.RequestHistory
	if len(got) != messageWindowHead+messageWindowTail {
		t.Fatalf("expected %d messages, got %d", messageWindowHead+messageWindowTail, len(got))
	}
	// head = history[5:8], so latest user is the first kept message.
	if got[0].Content != "latest-user" {
		t.Fatalf("expected latest user as first head message, got %q", got[0].Content)
	}
	if got[1].Content != history[6].Content || got[2].Content != history[7].Content {
		t.Fatalf("expected the two messages after latest user preserved")
	}
	// tail = history[n-46:], so head (3) is followed by original index n-46.
	if got[messageWindowHead].Content != history[n-messageWindowTail].Content {
		t.Fatalf("expected tail to start at original index %d", n-messageWindowTail)
	}
}

func TestMessageWindow_DropsOrphanToolAtCut(t *testing.T) {
	// Latest user at index 0 so trimming triggers; make the first tail message
	// an orphan tool (no preceding assistant call in the window).
	const n = 60
	history := make([]storage.Message, n)
	history[0] = storage.Message{Role: "user", Content: "latest-user"}
	for i := 1; i < n; i++ {
		history[i] = storage.Message{Role: "assistant", Content: "filler"}
	}
	history[n-messageWindowTail] = storage.Message{Role: "tool", ToolCallID: "orphan", Content: `{"status":"success","result":"x"}`}
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

func TestMessageWindow_NoopWhenCurrentTurnWithinLimit(t *testing.T) {
	// 80 messages total but the latest user turn is short (10 messages after the
	// latest user), so the current turn count (10) is within the limit and no
	// trimming happens despite the large earlier history.
	const n = 80
	history := make([]storage.Message, n)
	for i := range history {
		history[i] = storage.Message{Role: "assistant", Content: "old"}
	}
	history[n-10] = storage.Message{Role: "user", Content: "latest-user"}
	req := &Request{RequestHistory: history}
	if err := (&MessageWindowCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if len(req.RequestHistory) != n {
		t.Fatalf("expected no trim when current turn within limit, got %d", len(req.RequestHistory))
	}
}

func TestMessageWindow_FallsBackToFirstThreeWhenNoUser(t *testing.T) {
	// No user message at all: turn count = whole history. With 60 messages the
	// fallback keeps the first message of the turn plus the following 2, then the
	// most recent 46.
	const n = 60
	history := make([]storage.Message, n)
	for i := range history {
		history[i] = storage.Message{Role: "assistant", Content: string(rune('A' + i%26))}
	}
	first := history[0].Content
	req := &Request{RequestHistory: history}
	if err := (&MessageWindowCompressionStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	got := req.RequestHistory
	if len(got) != messageWindowHead+messageWindowTail {
		t.Fatalf("expected %d messages, got %d", messageWindowHead+messageWindowTail, len(got))
	}
	if got[0].Content != first {
		t.Fatalf("expected head to start at the first message, got %q", got[0].Content)
	}
	if got[messageWindowHead].Content != history[n-messageWindowTail].Content {
		t.Fatalf("expected tail to start at original index %d", n-messageWindowTail)
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

package compression

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

// realUserMsg 构造一条真实用户消息（新对话轮的起点）。
func realUserMsg(content string) storage.Message {
	return storage.Message{Role: "user", Content: content}
}

// reminderUserMsg 构造一条内部注入的 <system-reminder> user 消息（不新起一轮）。
func reminderUserMsg(content string) storage.Message {
	return storage.Message{Role: "user", Content: "<system-reminder>\n" + content + "\n</system-reminder>"}
}

func assistantTextMsg(content string) storage.Message {
	return storage.Message{Role: "assistant", Content: content}
}

// makeTurnHistory 构造 n 个简单对话轮（每轮 = user + assistant），
// content 用给定长度填充，便于控制是否超出 token 预算。
func makeTurnHistory(n int, fill int) []storage.Message {
	var h []storage.Message
	for i := 0; i < n; i++ {
		h = append(h,
			realUserMsg("u"+strings.Repeat("x", fill)),
			assistantTextMsg("a"+strings.Repeat("y", fill)),
		)
	}
	return h
}

// --- 轮切分 ---

func TestSplitConversationTurns_InjectedReminderNotNewTurn(t *testing.T) {
	history := []storage.Message{
		realUserMsg("first"),
		assistantTextMsg("resp1"),
		reminderUserMsg("continue please"), // 内部注入，归入第一轮
		assistantTextMsg("resp2"),
		realUserMsg("second"), // 真实用户，第二轮起点
		assistantTextMsg("resp3"),
	}
	turns := splitConversationTurns(history)
	if len(turns) != 2 {
		t.Fatalf("expected 2 turns, got %d: %#v", len(turns), turns)
	}
	if len(turns[0]) != 4 {
		t.Fatalf("expected first turn to absorb injected reminder (4 msgs), got %d", len(turns[0]))
	}
	if len(turns[1]) != 2 {
		t.Fatalf("expected second turn 2 msgs, got %d", len(turns[1]))
	}
}

func TestSplitConversationTurns_LeadingNonUserPrefix(t *testing.T) {
	// 以 system + summary 前缀开头（上一次压缩产物），它们不属于任何对话轮。
	history := []storage.Message{
		{Role: "system", Content: summarySystemPreamble},
		{Role: "user", Content: "<conversation-summary>\nprev\n</conversation-summary>"},
		realUserMsg("hello"),
		assistantTextMsg("hi"),
	}
	turns := splitConversationTurns(history)
	// 前缀两条不成轮；只有 1 个真实对话轮。
	if len(turns) != 1 {
		t.Fatalf("expected 1 turn ignoring summary prefix, got %d: %#v", len(turns), turns)
	}
}

// --- 单次压缩：满足预算即停止（只剥一次） ---

func TestReactiveCompact_StopsAfterFirstStripWhenWithinBudget(t *testing.T) {
	store := &fakeStore{}
	summaryCalls := 0
	var summarizedInput []storage.Message
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		if !r.Aggressive {
			t.Fatalf("expected aggressive summary prompt for reactive compact")
		}
		summaryCalls++
		summarizedInput = r.History
		return SummaryResult{Summary: "SUM"}, nil
	}
	// 小历史：剥离并摘要后（摘要很短）远低于预算 → 首次剥离后即停止。
	history := makeTurnHistory(10, 20)
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: history,
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	if err := (&ReactiveCompactStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	if summaryCalls != 1 {
		t.Fatalf("expected exactly 1 strip/summary when already within budget after first strip, got %d", summaryCalls)
	}
	if len(summarizedInput) == 0 {
		t.Fatalf("expected summarizer to receive stripped turns")
	}
	// 结果 = system 前导 + summary + 未剥离轮。
	if req.RequestHistory[0].Role != "system" {
		t.Fatalf("expected system preamble first, got %#v", req.RequestHistory[0])
	}
	if !strings.Contains(req.RequestHistory[1].Content, "SUM") {
		t.Fatalf("expected summary as second message, got %#v", req.RequestHistory[1])
	}
}

// --- 单次压缩：始终超预算时，最多剥离 3 次后仍超限则返回错误 ---

func TestReactiveCompact_LoopsAtMostThreeStripsThenErrorsWhenStillOverBudget(t *testing.T) {
	store := &fakeStore{}
	summaryCalls := 0
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		summaryCalls++
		// 摘要本身也很大，确保重组后仍超预算，从而持续触发剥离直到 3 次上限。
		return SummaryResult{Summary: strings.Repeat("S", 700*1024)}, nil
	}
	// 每轮内容都很大，且轮数多，使每次剥离后仍超预算。
	history := makeTurnHistory(60, 40*1024)
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: history,
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	// 3 次剥离后仍超预算 → 返回错误。
	if err := (&ReactiveCompactStrategy{}).Apply(context.Background(), req); err == nil {
		t.Fatalf("expected error when still over budget after 3 strips")
	}
	if summaryCalls != reactiveMaxStrips {
		t.Fatalf("expected exactly %d strips (max), got %d", reactiveMaxStrips, summaryCalls)
	}
}

// --- 摘要调用自身 413/报错 → 直接返回错误，不做收缩重试 ---

func TestReactiveCompact_SummaryErrorReturnsErrorNoRetry(t *testing.T) {
	store := &fakeStore{}
	calls := 0
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		calls++
		return SummaryResult{}, errors.New("boom 413")
	}
	history := makeTurnHistory(10, 20)
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: history,
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	err := (&ReactiveCompactStrategy{}).Apply(context.Background(), req)
	if err == nil {
		t.Fatalf("expected error propagated when summary fails")
	}
	if calls != 1 {
		t.Fatalf("expected summarizer called exactly once (no shrink retry), got %d", calls)
	}
}

// --- 重摘要：第二次剥离时输入包含上一次的摘要 ---

func TestReactiveCompact_ReSummarizesWithPriorSummary(t *testing.T) {
	store := &fakeStore{}
	var inputs [][]storage.Message
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		inputs = append(inputs, r.History)
		// 保持超预算，逼出第二次剥离。
		return SummaryResult{Summary: strings.Repeat("S", 700*1024)}, nil
	}
	history := makeTurnHistory(60, 40*1024)
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: history,
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	// 保持超预算 → 3 次剥离后返回错误（本用例只关注重摘要输入，忽略该错误）。
	_ = (&ReactiveCompactStrategy{}).Apply(context.Background(), req)
	if len(inputs) < 2 {
		t.Fatalf("expected at least 2 strips to observe re-summarize, got %d", len(inputs))
	}
	// 第二次摘要输入应包含上一次摘要（<conversation-summary> 前置消息）。
	foundPrior := false
	for _, m := range inputs[1] {
		if strings.Contains(m.Content, "<conversation-summary>") {
			foundPrior = true
		}
	}
	if !foundPrior {
		t.Fatalf("expected second summarize input to include prior summary")
	}
}

// --- 保留消息丢弃 reasoning_content ---

func TestReactiveCompact_DropsReasoningAndRepairsBoundaries(t *testing.T) {
	store := &fakeStore{}
	summarizer := func(ctx context.Context, r SummaryRequest) (SummaryResult, error) {
		return SummaryResult{Summary: "S"}, nil
	}
	history := []storage.Message{
		realUserMsg("t1"), assistantTextMsg("a1"),
		realUserMsg("t2"), assistantTextMsg("a2"),
		realUserMsg("t3"),
		{Role: "assistant", ReasoningContent: "SECRET THOUGHT", ToolCalls: []storage.MessageToolCall{{ID: "c1", Type: "function", Function: storage.MessageFunctionCall{Name: "bash", Arguments: "{}"}}}},
		toolMsg("c1", "success", "ok"),
	}
	req := &Request{
		User: storage.User{ID: "u"}, Conversation: storage.Conversation{ID: "c"},
		RequestHistory: history,
		Store:          store,
		Estimator:      DefaultTokenEstimator{},
		Summarizer:     summarizer,
	}
	if err := (&ReactiveCompactStrategy{}).Apply(context.Background(), req); err != nil {
		t.Fatalf("apply: %v", err)
	}
	for _, m := range req.RequestHistory {
		if m.ReasoningContent != "" {
			t.Fatalf("expected reasoning_content dropped, found in %#v", m)
		}
	}
}

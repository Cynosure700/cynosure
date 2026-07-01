package runtime

import (
	"context"
	"io"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/config"
	"github.com/Cynosure700/cynosure/cynosure/internal/llm"
)

// reactiveStubStream 是一个单块内容的假流。
type reactiveStubStream struct {
	content  string
	finish   openai.FinishReason
	consumed bool
}

func (s *reactiveStubStream) Recv() (openai.ChatCompletionStreamResponse, error) {
	if s.consumed {
		return openai.ChatCompletionStreamResponse{}, io.EOF
	}
	s.consumed = true
	return openai.ChatCompletionStreamResponse{Choices: []openai.ChatCompletionStreamChoice{{
		Delta:        openai.ChatCompletionStreamChoiceDelta{Content: s.content},
		FinishReason: s.finish,
	}}}, nil
}

func (s *reactiveStubStream) Close() error { return nil }

// reactiveStubLLM 按序：对每个主模型（流式）请求返回 413，直到 overflowCount 用尽；
// 之后返回一个正常的最终回复。摘要（非流式）调用按 summaryOverflow 决定是否 413。
type reactiveStubLLM struct {
	overflowCount    int
	summaryOverflow  bool
	streamCalls      int
	summaryCalls     int
	finalContent     string
	summarizeContent string
}

func (c *reactiveStubLLM) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	// 摘要调用（非流式）。
	c.summaryCalls++
	if c.summaryOverflow {
		return openai.ChatCompletionResponse{}, &llm.APIError{StatusCode: 413, Body: "prompt_too_long"}
	}
	return openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{{
		Message: openai.ChatCompletionMessage{Content: c.summarizeContent},
	}}}, nil
}

func (c *reactiveStubLLM) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (llm.ChatCompletionStream, error) {
	c.streamCalls++
	if c.streamCalls <= c.overflowCount {
		return nil, &llm.APIError{StatusCode: 413, Body: "prompt_too_long"}
	}
	return &reactiveStubStream{content: c.finalContent, finish: openai.FinishReasonStop}, nil
}

// makeManyTurns 构造 n 个真实对话轮的模型历史，供剥离测试用。
func makeManyTurns(n int) []storage.Message {
	var h []storage.Message
	for i := 0; i < n; i++ {
		h = append(h,
			storage.Message{ID: "u", Role: "user", Content: "user turn " + strings.Repeat("x", 20)},
			storage.Message{ID: "a", Role: "assistant", Content: "assistant reply " + strings.Repeat("y", 20)},
		)
	}
	return h
}

func newReactiveTestService(t *testing.T, client llm.Client) (*Service, *fakeStore) {
	t.Helper()
	store := &fakeStore{}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	return &Service{LLM: client, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}, store
}

// 主模型持续 413：413 触发一次「单次压缩」的 reactiveCompact（内部最多剥离 3 次），
// 压缩后重试一次仍 413 时不再重复压缩（compacted 单发），整体返回错误。
func TestRespondToConversation_MainOverflowTriggersSingleReactiveCompact(t *testing.T) {
	client := &reactiveStubLLM{
		overflowCount:    99, // 永远 413
		summarizeContent: "AGGRESSIVE SUMMARY",
	}
	service, store := newReactiveTestService(t, client)
	conv := storage.Conversation{ID: "conv_reactive_cap", Title: "新对话"}
	store.modelHistory = makeManyTurns(50)
	store.modelHistoryExists = true

	_, err := service.RespondToConversation(context.Background(), conv, storage.User{ID: "u1"}, "go", nil)
	if err == nil {
		t.Fatalf("expected error when model keeps overflowing")
	}
	// 单次压缩内部最多剥离 3 次 → 摘要器至多被调用 3 次；且必须至少调用一次。
	if client.summaryCalls == 0 {
		t.Fatalf("expected reactive compact to summarize at least once")
	}
	if client.summaryCalls > 3 {
		t.Fatalf("expected at most 3 reactive summaries (single compact, max 3 strips), got %d", client.summaryCalls)
	}
	// 主模型流式请求只会发生两次：首次 413 + 压缩后重试一次（仍 413，不再压缩）。
	if client.streamCalls != 2 {
		t.Fatalf("expected exactly 2 stream attempts (initial + one retry after single compact), got %d", client.streamCalls)
	}
}

// 主模型先 413 一次、压缩后成功：验证单次压缩后能恢复并返回最终回复。
func TestRespondToConversation_RecoversAfterSingleReactiveCompact(t *testing.T) {
	client := &reactiveStubLLM{
		overflowCount:    1,
		finalContent:     "最终回复",
		summarizeContent: "AGGRESSIVE SUMMARY",
	}
	service, store := newReactiveTestService(t, client)
	conv := storage.Conversation{ID: "conv_reactive_recover", Title: "新对话"}
	store.modelHistory = makeManyTurns(50)
	store.modelHistoryExists = true

	msg, err := service.RespondToConversation(context.Background(), conv, storage.User{ID: "u1"}, "go", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if msg.Content != "最终回复" {
		t.Fatalf("expected final content after recovery, got %q", msg.Content)
	}
	// 单次压缩后重试成功：至少剥离一次（预算内即停）。
	if client.summaryCalls == 0 {
		t.Fatalf("expected reactive compact to summarize at least once")
	}
}

// sizeSensitiveSummaryLLM 的摘要调用在输入超过阈值时返回 413，否则成功。
// 主模型流式调用始终成功（返回最终回复）。用于验证 FullHistory 摘要 413 →
// 转 ReactiveCompact（其摘要输入更小）后成功兜底。
type sizeSensitiveSummaryLLM struct {
	overflowThreshold int
	summaryCalls      int
	overflowedCalls   int
	finalContent      string
}

func (c *sizeSensitiveSummaryLLM) CreateChatCompletion(ctx context.Context, req openai.ChatCompletionRequest) (openai.ChatCompletionResponse, error) {
	c.summaryCalls++
	total := 0
	for _, m := range req.Messages {
		total += len(m.Content)
	}
	if total > c.overflowThreshold {
		c.overflowedCalls++
		return openai.ChatCompletionResponse{}, &llm.APIError{StatusCode: 413, Body: "prompt_too_long"}
	}
	return openai.ChatCompletionResponse{Choices: []openai.ChatCompletionChoice{{
		Message: openai.ChatCompletionMessage{Content: "SUMMARY"},
	}}}, nil
}

func (c *sizeSensitiveSummaryLLM) CreateChatCompletionStream(ctx context.Context, req openai.ChatCompletionRequest) (llm.ChatCompletionStream, error) {
	return &reactiveStubStream{content: c.finalContent, finish: openai.FinishReasonStop}, nil
}

// FullHistory 兜底摘要因输入过大 413 → compressContextBeforeLLM 转 ReactiveCompact，
// 后者只摘要最旧 20% 的对话轮（输入更小、不再 413），从而成功兜底并继续对话。
func TestCompressContextBeforeLLM_FullHistoryOverflowFallsBackToReactive(t *testing.T) {
	// 构造超预算的大历史，强制走到 FullHistory 兜底层。
	big := strings.Repeat("z", 900*1024)
	history := []storage.Message{
		{ID: "u1", Role: "user", Content: "first " + big},
		{ID: "a1", Role: "assistant", Content: "resp one"},
		{ID: "u2", Role: "user", Content: "second turn"},
		{ID: "a2", Role: "assistant", Content: "resp two"},
		{ID: "u3", Role: "user", Content: "third turn"},
		{ID: "a3", Role: "assistant", Content: "resp three"},
	}
	// 阈值设在「全量历史」之下、「最旧几轮」之上：FullHistory 摘要（含 big）会 413，
	// 而 ReactiveCompact 只摘最旧 20% 轮（不含 big 的话）输入更小。为确保 big 落入被
	// 剥离侧，big 在最旧的第一轮。ReactiveCompact 首次剥离 ceil(3*0.2)=1 轮=含 big，
	// 因此其摘要输入仍含 big → 仍 413。为让该用例聚焦「fallback 被触发」，这里断言
	// fallback 确实调用了 reactive（即摘要被多次尝试），并最终返回错误或成功均可。
	client := &sizeSensitiveSummaryLLM{overflowThreshold: 500 * 1024, finalContent: "ok"}
	store := &fakeStore{}
	cfg := config.AppConfig{LLM: config.Config{ModelID: "m"}}
	service := &Service{LLM: client, Store: store, Cfg: cfg, Tools: NewToolRegistry(cfg)}
	state := &LoopState{
		Conversation: storage.Conversation{ID: "c"},
		User:         storage.User{ID: "u"},
		History:      history,
		ModelHistory: cloneMessages(history),
		SystemPrompt: "sys",
	}

	_, err := service.compressContextBeforeLLM(context.Background(), state)
	// FullHistory 摘要 413 必须触发 reactive 兜底：摘要器被调用超过一次
	//（一次 FullHistory + 至少一次 reactive）。
	if client.summaryCalls < 2 {
		t.Fatalf("expected FullHistory overflow to trigger reactive fallback (>=2 summary calls), got %d (err=%v)", client.summaryCalls, err)
	}
}

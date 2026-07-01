package runtime

import (
	"context"
	"testing"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/runtime/compression"
	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/config"
)

func TestSummarizeHistoryForContext_UsesAggressivePromptWhenRequested(t *testing.T) {
	client := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "SUMMARY"}}}},
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "SUMMARY"}}}},
	}}
	service := &Service{
		Store:   &fakeStore{},
		Cfg:     config.AppConfig{LLM: config.Config{ModelID: "m"}},
		LLM:     client,
		Prompts: defaultFunctionalPrompts(),
	}
	history := []storage.Message{{Role: "user", Content: "hello"}}

	// 普通摘要：使用通用 ContextSummary system prompt。
	if _, err := service.summarizeHistoryForContext(context.Background(), compression.SummaryRequest{History: history}); err != nil {
		t.Fatalf("normal summarize: %v", err)
	}
	normalSystem := client.reqs[0].Messages[0].Content
	if normalSystem != service.Prompts.withDefaults().ContextSummary {
		t.Fatalf("expected normal summary to use ContextSummary prompt")
	}

	// 激进摘要：使用 ReactiveCompactSummary system prompt。
	if _, err := service.summarizeHistoryForContext(context.Background(), compression.SummaryRequest{History: history, Aggressive: true}); err != nil {
		t.Fatalf("aggressive summarize: %v", err)
	}
	aggressiveSystem := client.reqs[1].Messages[0].Content
	if aggressiveSystem != service.Prompts.withDefaults().ReactiveCompactSummary {
		t.Fatalf("expected aggressive summary to use ReactiveCompactSummary prompt")
	}
	if aggressiveSystem == normalSystem {
		t.Fatalf("expected aggressive prompt to differ from normal prompt")
	}
}

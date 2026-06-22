package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/runtime/compression"
	"cynosure/internal/agent/storage"
)

// summarizeHistoryForContext 执行一次非工具的 LLM 调用，对请求态历史进行摘要。
// 它绝不会改动 state.History。
func (s *Service) summarizeHistoryForContext(ctx context.Context, req compression.SummaryRequest) (compression.SummaryResult, error) {
	if s.LLM == nil {
		return compression.SummaryResult{}, fmt.Errorf("llm client is not configured")
	}
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: s.Prompts.withDefaults().ContextSummary},
		{Role: "user", Content: renderHistoryForSummary(req.History)},
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    s.Cfg.LLM.ModelID,
		Messages: messages,
	})
	if err != nil {
		return compression.SummaryResult{}, err
	}
	if len(resp.Choices) == 0 {
		return compression.SummaryResult{}, fmt.Errorf("summary response had no choices")
	}
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if summary == "" {
		return compression.SummaryResult{}, fmt.Errorf("summary response was empty")
	}
	return compression.SummaryResult{Summary: summary}, nil
}

func renderHistoryForSummary(history []storage.Message) string {
	var b strings.Builder
	b.WriteString("Summarize the following conversation history:\n\n")
	for _, msg := range history {
		role := msg.Role
		switch role {
		case "tool":
			b.WriteString("[tool result] ")
			b.WriteString(msg.Content)
		case "assistant":
			b.WriteString("[assistant] ")
			if strings.TrimSpace(msg.Content) != "" {
				b.WriteString(msg.Content)
			}
			for _, call := range msg.ToolCalls {
				fmt.Fprintf(&b, "\n[tool call] %s(%s)", call.Function.Name, call.Function.Arguments)
			}
		default:
			fmt.Fprintf(&b, "[%s] %s", role, msg.Content)
		}
		b.WriteString("\n\n")
	}
	return b.String()
}

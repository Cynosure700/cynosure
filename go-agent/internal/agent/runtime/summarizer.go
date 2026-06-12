package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
)

const summarizerSystemPrompt = `You are a context compaction engine. Summarize the conversation so far into a structured Markdown brief that another AI agent can use to continue the task without losing critical information.

Rules:
- Do NOT call tools or invent information that is not present.
- Preserve concrete details: current user goal, completed actions and key conclusions, important file paths, function names, commands, errors, and decisions.
- For any <persisted-output ...> marker you see, keep its id and note that read_persisted_output can fetch the full content.
- If a tool result shows "[Earlier result compacted. Re-run if needed]", note that it must be re-run if needed.
- List unfinished items, things to verify, and next steps.

Output only the Markdown summary.`

// summarizeHistoryForContext runs a single non-tool LLM call to summarize the
// request-state history. It never mutates state.History.
func (s *Service) summarizeHistoryForContext(ctx context.Context, req compression.SummaryRequest) (compression.SummaryResult, error) {
	if s.LLM == nil {
		return compression.SummaryResult{}, fmt.Errorf("llm client is not configured")
	}
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: summarizerSystemPrompt},
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

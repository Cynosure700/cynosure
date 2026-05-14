package sessions

import (
	"context"
	"encoding/json"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/tools"
)

const maxSubagentRounds = 30

const subagentSystemPrompt = `You are a coding subagent. Complete the given task using the available tools, then summarize your findings concisely.`

func RunSubagent(ctx context.Context, prompt string) (string, error) {
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: subagentSystemPrompt},
		{Role: "user", Content: prompt},
	}

	lastContent := "(no summary)"

	for round := 0; round < maxSubagentRounds; round++ {
		req := openai.ChatCompletionRequest{
			Model:    config.ModelID,
			Messages: messages,
			Tools:    tools.ChildToolDefs,
		}
		reqBody, _ := json.Marshal(req)

		resp, err := config.Client.CreateChatCompletion(ctx, req)

		respBody, _ := json.Marshal(resp)
		logger.LogLLMRound(round+1, "subagent", reqBody, respBody, err)

		if err != nil {
			return "", fmt.Errorf("subagent API error: %w", err)
		}

		if len(resp.Choices) == 0 {
			return lastContent, nil
		}

		choice := resp.Choices[0]
		msg := choice.Message
		messages = append(messages, msg)

		if msg.Content != "" {
			lastContent = msg.Content
		}

		if choice.FinishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			break
		}

		for _, tc := range msg.ToolCalls {
			handler, ok := tools.Handlers[tc.Function.Name]
			if !ok {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       "tool",
					ToolCallID: tc.ID,
					Content:    fmt.Sprintf("Unknown tool: %s", tc.Function.Name),
				})
				continue
			}

			var args map[string]any
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				args = map[string]any{}
			}

			result, err := handler(ctx, args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    result,
			})
		}
	}

	return lastContent, nil
}

func init() {
	tools.SetHandler("task", func(ctx context.Context, args map[string]any) (string, error) {
		prompt, _ := args["prompt"].(string)
		if prompt == "" {
			return "", fmt.Errorf("prompt is required")
		}
		return RunSubagent(ctx, prompt)
	})
}

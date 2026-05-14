package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
	"nano_cc/internal/tools"
)

func AgentLoop(systemPrompt string, messages []openai.ChatCompletionMessage) (string, []openai.ChatCompletionMessage, error) {
	ctx := context.Background()

	roundsSinceTodo := 0
	round := 0
	var prevRoundSummary string

	for {
		round++

		// Inject previous round summary
		if prevRoundSummary != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    "user",
				Content: prevRoundSummary,
			})
			prevRoundSummary = ""
		}

		// Nag reminder
		if roundsSinceTodo >= 3 {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    "user",
				Content: "<reminder>Update your todos using the todo tool.</reminder>",
			})
		}

		// Layer 1: micro compact
		sessions.MicroCompact(messages)

		// Layer 2: auto compact
		if sessions.ShouldCompact(messages) {
			compressed, err := sessions.AutoCompact(messages)
			if err != nil {
				logger.Warn(fmt.Sprintf("auto compact failed: %v", err))
			} else {
				messages = append([]openai.ChatCompletionMessage{
					{Role: "system", Content: systemPrompt},
				}, compressed...)
			}
		}

		req := openai.ChatCompletionRequest{
			Model:    config.ModelID,
			Messages: messages,
			Tools:    tools.ParentToolDefs,
		}
		reqBody, _ := json.Marshal(req)

		resp, err := config.Client.CreateChatCompletion(ctx, req)

		respBody, _ := json.Marshal(resp)
		logger.LogLLMRound(round, "agent", reqBody, respBody, err)

		if err != nil {
			return "", messages, fmt.Errorf("API error: %w", err)
		}

		if len(resp.Choices) == 0 {
			return "", messages, fmt.Errorf("no response from model")
		}

		choice := resp.Choices[0]
		msg := choice.Message
		messages = append(messages, msg)

		// No tool calls → return text
		if choice.FinishReason != "tool_calls" || len(msg.ToolCalls) == 0 {
			if msg.Content != "" {
				return msg.Content, messages, nil
			}
			return "(no response)", messages, nil
		}

		usedTodo := false
		manualCompact := false

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

			logger.Tool(tc.Function.Name, tc.Function.Arguments)

			result, err := handler(ctx, args)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:       "tool",
				ToolCallID: tc.ID,
				Content:    result,
			})

			if tc.Function.Name == "todo" {
				usedTodo = true
			}
			if tc.Function.Name == "compact" {
				manualCompact = true
			}
		}

		// Build previous round summary for next iteration
		prevRoundSummary = buildRoundSummary(msg, round)

		// Layer 3: manual compact
		if manualCompact {
			compressed, err := sessions.AutoCompact(messages)
			if err != nil {
				logger.Warn(fmt.Sprintf("manual compact failed: %v", err))
			} else {
				messages = append([]openai.ChatCompletionMessage{
					{Role: "system", Content: systemPrompt},
				}, compressed...)
			}
		}

		if usedTodo {
			roundsSinceTodo = 0
		} else {
			roundsSinceTodo++
		}
	}
}

func buildRoundSummary(msg openai.ChatCompletionMessage, round int) string {
	var parts []string

	if msg.Content != "" {
		parts = append(parts, fmt.Sprintf("Thought: %s", msg.Content))
	}

	if len(msg.ToolCalls) > 0 {
		toolNames := make([]string, len(msg.ToolCalls))
		for i, tc := range msg.ToolCalls {
			toolNames[i] = tc.Function.Name
		}
		parts = append(parts, fmt.Sprintf("Tools used: %s", strings.Join(toolNames, ", ")))
	}

	if len(parts) == 0 {
		return ""
	}

	return fmt.Sprintf("<system-reminder>Round %d completed. %s</system-reminder>", round, strings.Join(parts, " | "))
}

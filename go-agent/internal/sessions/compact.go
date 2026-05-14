package sessions

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/tools"
)

const (
	tokenThreshold  = 50000
	keepRecent      = 3
	transcriptDir   = ".transcripts"
	summaryMaxChars = 80000
)

func EstimateTokens(messages []openai.ChatCompletionMessage) int {
	b, err := json.Marshal(messages)
	if err != nil {
		return 0
	}
	return len(b) / 4
}

func ShouldCompact(messages []openai.ChatCompletionMessage) bool {
	return EstimateTokens(messages) > tokenThreshold
}

func MicroCompact(messages []openai.ChatCompletionMessage) {
	toolIndices := make([]int, 0)
	for i, msg := range messages {
		if msg.Role == "tool" {
			toolIndices = append(toolIndices, i)
		}
	}

	if len(toolIndices) <= keepRecent {
		return
	}

	// Build tool_call_id → function_name mapping from assistant messages
	toolNames := make(map[string]string)
	for _, msg := range messages {
		if msg.Role == "assistant" {
			for _, tc := range msg.ToolCalls {
				toolNames[tc.ID] = tc.Function.Name
			}
		}
	}

	keepFrom := toolIndices[len(toolIndices)-keepRecent]
	for _, idx := range toolIndices {
		if idx >= keepFrom {
			break
		}
		toolName := toolNames[messages[idx].ToolCallID]
		if toolName == "" {
			toolName = "unknown"
		}
		messages[idx].Content = fmt.Sprintf("[Previous: used %s]", toolName)
	}
}

func AutoCompact(messages []openai.ChatCompletionMessage) ([]openai.ChatCompletionMessage, error) {
	// Save transcript
	if err := os.MkdirAll(transcriptDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create transcript dir: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	transcriptPath := filepath.Join(transcriptDir, fmt.Sprintf("transcript_%s.jsonl", timestamp))

	f, err := os.Create(transcriptPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create transcript file: %w", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	for _, msg := range messages {
		if err := enc.Encode(msg); err != nil {
			return nil, fmt.Errorf("failed to write transcript: %w", err)
		}
	}

	// Build summary prompt
	raw, _ := json.Marshal(messages)
	content := string(raw)
	if len(content) > summaryMaxChars {
		content = content[:summaryMaxChars]
	}

	summaryPrompt := fmt.Sprintf(
		"Summarize the following conversation between an AI agent and its tools. Focus on what was accomplished, key decisions made, and current state:\n\n%s",
		content,
	)

	req := openai.ChatCompletionRequest{
		Model: config.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "user", Content: summaryPrompt},
		},
	}
	reqBody, _ := json.Marshal(req)

	resp, err := config.Client.CreateChatCompletion(context.Background(), req)

	respBody, _ := json.Marshal(resp)
	logger.LogLLMRound(0, "compact", reqBody, respBody, err)

	if err != nil {
		return nil, fmt.Errorf("summary API error: %w", err)
	}

	summary := "No summary available."
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Content != "" {
		summary = resp.Choices[0].Message.Content
	}

	return []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: fmt.Sprintf("[Conversation compressed. Transcript: %s]\n\n%s", transcriptPath, summary),
		},
		{
			Role:    "assistant",
			Content: "Understood. I have the context from the summary. Continuing.",
		},
	}, nil
}

func init() {
	tools.SetHandler("compact", func(ctx context.Context, args map[string]any) (string, error) {
		return "Manual compression requested.", nil
	})
}

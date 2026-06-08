package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/logger"
	"nano_cc/internal/web/storage"
)

const topicSummarizerSystemPrompt = `你是一个对话话题提炼引擎。把给定的对话压缩成 3-6 个话题短语。

规则：
- 每个话题是一句话级别的简短短语，不超过 30 个字。
- 只描述"聊过什么事"，不要包含对话原文、不要展开结论细节。
- 只输出一个 JSON 字符串数组，例如 ["话题一","话题二"]，不要输出任何其它内容。`

const (
	maxTopicCount  = 6
	maxTopicRunes  = 30
	maxTopicChars  = 4000
)

// updateConversationTopics extracts a lightweight topic list from the
// conversation history with a single non-tool LLM call and overwrites the
// stored topic list. It is best-effort: failures are logged and swallowed so
// the user-facing response is never affected.
func (s *Service) updateConversationTopics(ctx context.Context, conversation storage.Conversation, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	transcript := renderHistoryForTopics(history)
	if strings.TrimSpace(transcript) == "" {
		return
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: topicSummarizerSystemPrompt},
			{Role: "user", Content: transcript},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: topic extraction failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	topics := parseTopics(resp.Choices[0].Message.Content)
	if len(topics) == 0 {
		return
	}
	topicsJSON, err := json.Marshal(topics)
	if err != nil {
		return
	}
	if err := s.Store.UpsertConversationTopics(ctx, storage.ConversationTopics{
		ConversationID: conversation.ID,
		UserID:         user.ID,
		TopicsJSON:     string(topicsJSON),
	}); err != nil {
		logger.Warn(fmt.Sprintf("memory: persist topics failed: %v", err))
	}
}

// parseTopics extracts a JSON string array from the model output, tolerating
// surrounding code fences or prose, then trims, dedups, and caps the result.
func parseTopics(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if start := strings.IndexByte(trimmed, '['); start >= 0 {
		if end := strings.LastIndexByte(trimmed, ']'); end > start {
			trimmed = trimmed[start : end+1]
		}
	}
	var parsed []string
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	seen := make(map[string]struct{})
	result := make([]string, 0, len(parsed))
	for _, topic := range parsed {
		topic = strings.TrimSpace(topic)
		if topic == "" {
			continue
		}
		if runes := []rune(topic); len(runes) > maxTopicRunes {
			topic = string(runes[:maxTopicRunes])
		}
		if _, exists := seen[topic]; exists {
			continue
		}
		seen[topic] = struct{}{}
		result = append(result, topic)
		if len(result) >= maxTopicCount {
			break
		}
	}
	return result
}

// renderHistoryForTopics builds a plain-text transcript of user and assistant
// messages only, dropping tool noise to keep the extraction call cheap.
func renderHistoryForTopics(history []storage.Message) string {
	var b strings.Builder
	for _, msg := range history {
		var content string
		switch msg.Role {
		case "user":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[user] ")
		case "assistant":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[assistant] ")
		default:
			continue
		}
		b.WriteString(content)
		b.WriteString("\n\n")
		if b.Len() > maxTopicChars {
			break
		}
	}
	return b.String()
}

type updateMemoryArgs struct {
	Profile string `json:"profile"`
}

// updateUserProfile validates the model-provided profile JSON and overwrites
// the user's profile card. Invalid JSON is rejected so the model can retry.
func (s *Service) updateUserProfile(ctx context.Context, toolCtx ToolContext, rawArgs string) (string, error) {
	var args updateMemoryArgs
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", fmt.Errorf("invalid update_memory arguments: %w", err)
	}
	profile := strings.TrimSpace(args.Profile)
	if profile == "" {
		return "", fmt.Errorf("profile is required")
	}
	if !json.Valid([]byte(profile)) {
		return "", fmt.Errorf("profile must be a valid JSON string")
	}
	if err := s.Store.UpsertUserProfile(ctx, storage.UserProfile{UserID: toolCtx.User.ID, ProfileJSON: profile}); err != nil {
		return "", fmt.Errorf("persist user profile: %w", err)
	}
	return "用户档案卡已更新。", nil
}

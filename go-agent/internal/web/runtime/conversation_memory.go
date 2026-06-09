package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
	"nano_cc/internal/web/storage"
)

const conversationMemorySystemPrompt = `你是“当前会话记忆”维护引擎。给你“已有会话记忆条目”和“最新一轮对话”，
请输出更新后的【完整】会话记忆条目列表，使其准确反映本场会话至今的主干信息：
- 覆盖：当前用户目标、关键决策与结论、已完成/已产出的内容、重要约束与上下文、待办与下一步。
- 合并重复、用新信息更新旧条目、删除已过时或被推翻的内容；只重组已知信息，不编造。
- name：短标题(<=80字，可用 [目标]/[决策]/[产出]/[待办] 等前缀)。description：一句话要点(<=300字)。body：支撑细节(<=2000字)。
- 仅输出 JSON 数组：[{"name","description","body"}]。无可记录时输出 []。`

type extractedConversationMemory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// updateConversationMemory runs once at the end of a conversation turn: it asks
// the LLM to rewrite the full conversation memory list based on existing entries
// plus the latest dialogue, then replaces the stored entries for this
// conversation. It is best-effort: failures are logged and swallowed so the
// user-facing response is never affected. The result is never injected into the
// system prompt; it is only consumed by the context compression pipeline.
func (s *Service) updateConversationMemory(ctx context.Context, conversation storage.Conversation, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderDialogueForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListConversationMemories(ctx, conversation.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: load existing failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: conversationMemorySystemPrompt},
			{Role: "user", Content: buildConversationMemoryUserPrompt(existing, dialogue)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: update failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	items := parseConversationMemories(resp.Choices[0].Message.Content)
	refined := make([]storage.ConversationMemory, 0, len(items))
	for _, it := range items {
		refined = append(refined, storage.ConversationMemory{
			ID:          idgen.New("cm"),
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		})
	}
	if err := s.Store.ReplaceConversationMemories(ctx, conversation.ID, user.ID, refined); err != nil {
		logger.Warn(fmt.Sprintf("conversation memory: replace failed: %v", err))
	}
}

// parseConversationMemories extracts a JSON array of memory objects from model
// output, tolerating surrounding code fences or prose, then truncates fields.
func parseConversationMemories(raw string) []extractedConversationMemory {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []extractedConversationMemory
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	result := make([]extractedConversationMemory, 0, len(parsed))
	for _, it := range parsed {
		it.Name = strings.TrimSpace(it.Name)
		it.Description = strings.TrimSpace(it.Description)
		it.Body = strings.TrimSpace(it.Body)
		if it.Name == "" || it.Body == "" {
			continue
		}
		it.Name = truncateRunes(it.Name, maxMemoryNameRunes)
		it.Description = truncateRunes(it.Description, maxMemoryDescRunes)
		it.Body = truncateRunes(it.Body, maxMemoryBodyRunes)
		result = append(result, it)
	}
	return result
}

func buildConversationMemoryUserPrompt(existing []storage.ConversationMemory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Existing conversation memory entries:\n")
	if len(existing) == 0 {
		b.WriteString("(none)")
	} else {
		for i, m := range existing {
			b.WriteString(fmt.Sprintf("[%d] %s: %s\n", i, m.Name, m.Description))
			if strings.TrimSpace(m.Body) != "" {
				b.WriteString("    body: " + m.Body + "\n")
			}
		}
	}
	b.WriteString("\n\nLatest dialogue:\n")
	b.WriteString(dialogue)
	return b.String()
}

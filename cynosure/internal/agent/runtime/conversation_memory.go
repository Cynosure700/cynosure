package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
)

type extractedConversationMemory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

func (s *Service) conversationMemorySystemPrompt() string {
	return s.Prompts.withDefaults().ConversationMemoryUpdate
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
			{Role: "system", Content: s.conversationMemorySystemPrompt()},
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

// scheduleMemoryWork 在一轮对话结束后异步执行收尾操作（模型历史持久化 + 记忆提取
// + 会话记忆更新）。
// 它接管入口持有的会话锁（token）：在独立的 background context 中执行，期间持续
// 续期，完成后停止续期并释放锁。返回 true 表示已接管锁所有权（调用方应跳过 defer
// 释放）；返回 false 表示未持锁（已降级），调用方按原逻辑处理。
// memoryOn 仅控制记忆提取与会话记忆更新；锁释放与模型历史持久化始终执行。
func (s *Service) scheduleMemoryWork(conv storage.Conversation, user storage.User, history []storage.Message, modelHistory []storage.Message, token string, stopRenew func(), memoryOn bool) bool {
	if token == "" {
		// 入口未持锁（已降级）→ 跳过收尾，不接管锁。
		return false
	}
	// 停止请求期看门狗，收尾 goroutine 内重新启动一个。
	if stopRenew != nil {
		stopRenew()
	}
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn(fmt.Sprintf("memory work: panic recovered conversation=%s: %v", conv.ID, r))
			}
		}()
		defer s.Store.ReleaseConversationLock(context.Background(), conv.ID, token)
		stop := s.startLockRenewer(conv.ID, token)
		defer stop()

		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.MemoryWorkTimeout)
		defer cancel()
		// 先落库模型历史：即使后续记忆相关的 LLM 调用超时，也不丢失本轮压缩成果。
		if len(modelHistory) > 0 {
			if err := s.Store.UpsertConversationModelHistory(ctx, conv.ID, user.ID, modelHistory); err != nil {
				logger.Warn(fmt.Sprintf("model history: persist failed conversation=%s: %v", conv.ID, err))
			}
		}
		if memoryOn {
			s.extractMemories(ctx, user, history)
			s.updateConversationMemory(ctx, conv, user, history)
		}
	}()
	return true
}

// startLockRenewer 启动一个后台看门狗，按 TTL/3 周期为会话锁续期，返回的 stop
// 函数用于停止续期（幂等）。续期失败（锁已不属于当前 token）时记录告警并停止。
func (s *Service) startLockRenewer(conversationID, token string) func() {
	ttl := s.Cfg.ConversationLockTTL
	interval := ttl / 3
	if interval <= 0 {
		interval = ttl
	}
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				ok, err := s.Store.RenewConversationLock(context.Background(), conversationID, token, ttl)
				if err != nil {
					logger.Warn(fmt.Sprintf("conversation lock: renew failed conversation=%s: %v", conversationID, err))
					return
				}
				if !ok {
					logger.Warn(fmt.Sprintf("conversation lock: renew lost ownership conversation=%s", conversationID))
					return
				}
			}
		}
	}()
	var once sync.Once
	return func() {
		once.Do(func() { close(done) })
	}
}

package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/idgen"
	"github.com/Cynosure700/cynosure/cynosure/internal/logger"
)

type extractedConversationMemory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

func (s *Service) conversationMemorySystemPrompt() string {
	return s.Prompts.withDefaults().ConversationMemoryUpdate
}

const (
	sessionMemoryTokenUnitK = 1024
	// sessionMemoryInitialTokens 是会话记忆初次提取的上下文 token 门槛（需求1）。
	sessionMemoryInitialTokens = 10 * sessionMemoryTokenUnitK
	// sessionMemoryTokenGrowth 是触发增量更新所需的 token 增长（需求2）。
	sessionMemoryTokenGrowth = 5 * sessionMemoryTokenUnitK
	// sessionMemoryToolCallsMin 是 round 中途触发更新所需的工具调用次数（需求2 条件1）。
	sessionMemoryToolCallsMin = 3
)

// sessionMemoryProgressFor 获取/创建某会话的会话记忆进度态。首次创建时按"已有会话
// 记忆是否存在"判定 extracted，避免进程重启后对老会话重新走 10K 门槛。调用方须持有
// sessionMemoryMu。
func (s *Service) sessionMemoryProgressForLocked(ctx context.Context, conversationID string) *sessionMemoryProgress {
	if s.sessionMemoryProgress == nil {
		s.sessionMemoryProgress = make(map[string]*sessionMemoryProgress)
	}
	if p, ok := s.sessionMemoryProgress[conversationID]; ok {
		return p
	}
	extracted := false
	if existing, err := s.Store.ListConversationMemories(ctx, conversationID); err == nil && len(existing) > 0 {
		extracted = true
	}
	p := &sessionMemoryProgress{extracted: extracted}
	s.sessionMemoryProgress[conversationID] = p
	return p
}

// shouldUpdate 判定是否应更新会话记忆。turnEnded=true 表示轮次自然结束的评估（条件2），
// false 表示 round 中途评估（条件1）。它只判定，不提交基线/断点/extracted——真正的提交
// 在更新成功后由调用方完成，避免失败时误推进。它会就地将基线下移到压缩后的低点。
// 调用方须持有 sessionMemoryMu。
func (p *sessionMemoryProgress) shouldUpdate(currentTokens int, turnEnded bool) bool {
	// 压缩后 token 可能回落到基线以下；将基线下移到当前值，使"增长"始终从最近低点起算。
	if currentTokens < p.baselineTokens {
		p.baselineTokens = currentTokens
	}
	if !p.extracted {
		// 需求1：初次提取仅在上下文达到 10K 时触发。
		return currentTokens >= sessionMemoryInitialTokens
	}
	growth := currentTokens - p.baselineTokens
	if growth < sessionMemoryTokenGrowth {
		return false
	}
	// 需求2：turnEnded（条件2）token 增长达标即可；round 中途（条件1）还需工具调用 >= 3。
	if !turnEnded && p.toolCallsSinceBase < sessionMemoryToolCallsMin {
		return false
	}
	return true
}

// maybeUpdateSessionMemoryMidLoop 在每个 tool_calls round 结束后评估条件(1)/初次提取，
// 满足时异步触发一次会话记忆更新。它累计本会话的工具调用数，并用单航班守卫避免并发更新。
func (s *Service) maybeUpdateSessionMemoryMidLoop(conversation storage.Conversation, user storage.User, history []storage.Message, currentTokens, roundToolCalls int) {
	s.sessionMemoryMu.Lock()
	p := s.sessionMemoryProgressForLocked(context.Background(), conversation.ID)
	p.toolCallsSinceBase += roundToolCalls
	if p.updating {
		// 上一次更新仍在跑 → 跳过，不重复触发、不动基线/断点。
		s.sessionMemoryMu.Unlock()
		return
	}
	if !p.shouldUpdate(currentTokens, false) {
		s.sessionMemoryMu.Unlock()
		return
	}
	p.updating = true
	s.sessionMemoryMu.Unlock()

	snapshot := cloneMessages(history)
	bpID := lastMessageID(snapshot)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Warn(fmt.Sprintf("session memory mid-loop: panic recovered conversation=%s: %v", conversation.ID, r))
			}
		}()
		ctx, cancel := context.WithTimeout(context.Background(), s.Cfg.MemoryWorkTimeout)
		defer cancel()
		ok := s.updateConversationMemory(ctx, conversation, user, snapshot)
		if ok {
			s.commitSessionMemoryBreakpoint(ctx, conversation.ID, currentTokens, bpID)
		}
		s.clearSessionMemoryUpdating(conversation.ID)
	}()
}

// shouldUpdateSessionMemoryAtTurnEnd 在轮次自然结束时评估条件(2)/初次提取，返回是否
// 应在收尾中更新会话记忆。返回 true 时会置位单航班守卫（由 scheduleMemoryWork 负责
// 在更新结束后清除）。若已有更新在跑则返回 false（best-effort，让中途那次覆盖最近状态）。
func (s *Service) shouldUpdateSessionMemoryAtTurnEnd(conversation storage.Conversation, currentTokens int) bool {
	s.sessionMemoryMu.Lock()
	defer s.sessionMemoryMu.Unlock()
	p := s.sessionMemoryProgressForLocked(context.Background(), conversation.ID)
	if p.updating {
		return false
	}
	if !p.shouldUpdate(currentTokens, true) {
		return false
	}
	p.updating = true
	return true
}

// commitSessionMemoryBreakpoint 在一次会话记忆更新成功后提交基线、标记 extracted、清零
// 工具调用计数，并把断点【持久化到会话记忆文件】（不在内存态保存）。
func (s *Service) commitSessionMemoryBreakpoint(ctx context.Context, conversationID string, baselineTokens int, breakpointID string) {
	s.sessionMemoryMu.Lock()
	p := s.sessionMemoryProgressForLocked(ctx, conversationID)
	p.extracted = true
	p.baselineTokens = baselineTokens
	p.toolCallsSinceBase = 0
	s.sessionMemoryMu.Unlock()
	if breakpointID != "" {
		if err := s.Store.SaveConversationMemoryBreakpoint(ctx, conversationID, breakpointID); err != nil {
			logger.Warn(fmt.Sprintf("session memory: persist breakpoint failed conversation=%s: %v", conversationID, err))
		}
	}
}

// clearSessionMemoryUpdating 清除某会话的单航班守卫。
func (s *Service) clearSessionMemoryUpdating(conversationID string) {
	s.sessionMemoryMu.Lock()
	defer s.sessionMemoryMu.Unlock()
	p := s.sessionMemoryProgressForLocked(context.Background(), conversationID)
	p.updating = false
}

// loadConversationMemoryBreakpoint 从会话记忆文件读取持久化的断点消息 ID（空表示未知/
// 未持久化）。best-effort：失败返回空串，压缩侧据此降级到全量摘要兜底。
func (s *Service) loadConversationMemoryBreakpoint(ctx context.Context, conversationID string) string {
	bp, err := s.Store.LoadConversationMemoryBreakpoint(ctx, conversationID)
	if err != nil {
		logger.Warn(fmt.Sprintf("session memory: load breakpoint failed conversation=%s: %v", conversationID, err))
		return ""
	}
	return bp
}

// lastMessageID 返回历史中最后一条消息的 ID（空历史返回空串）。
func lastMessageID(history []storage.Message) string {
	if len(history) == 0 {
		return ""
	}
	return history[len(history)-1].ID
}

// updateConversationMemory 在一个会话回合结束时、或满足触发条件的回合中途运行：
// 它请求 LLM 基于既有条目加上最新对话来重写完整的会话记忆列表，然后替换该会话
// 已存储的条目。它是尽力而为的：失败会被记录并吞掉，从而绝不影响面向用户的响应。
// 结果绝不会注入系统提示词，它只被上下文压缩流水线消费。
//
// 仅当重写成功且产生了非空替换时返回 true，以便调用方提交会话记忆断点/基线。
func (s *Service) updateConversationMemory(ctx context.Context, conversation storage.Conversation, user storage.User, history []storage.Message) bool {
	if s.LLM == nil {
		return false
	}
	dialogue := renderModelHistoryForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return false
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
		return false
	}
	if len(resp.Choices) == 0 {
		return false
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
		return false
	}
	return true
}

// parseConversationMemories 从模型输出中提取一个记忆对象的 JSON 数组，容忍其
// 周围的代码围栏或散文，然后截断各字段。
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
// history 是唯一的真实消息历史 state.ModelHistory（压缩后真实消息线，含本轮最终
// assistant）：既落库为下轮复用，又作为记忆提取/会话记忆更新的源。
// 它接管入口持有的会话锁（token）：在独立的 background context 中执行，期间持续
// 续期，完成后停止续期并释放锁。返回 true 表示已接管锁所有权（调用方应跳过 defer
// 释放）；返回 false 表示未持锁（已降级），调用方按原逻辑处理。
// memoryOn 仅控制记忆提取与会话记忆更新；锁释放与模型历史持久化始终执行。
// updateSessionMemory 控制本轮收尾是否更新会话记忆（由触发条件评估决定）；为 true 时
// 还会在更新成功后提交会话记忆断点/基线（断点=history 最后一条消息，基线=baselineTokens），
// 并清除该会话的单航班守卫。
func (s *Service) scheduleMemoryWork(conv storage.Conversation, user storage.User, history []storage.Message, token string, stopRenew func(), memoryOn bool, updateSessionMemory bool, baselineTokens int) bool {
	if token == "" {
		// 入口未持锁（已降级）→ 跳过收尾，不接管锁。
		// 但仍需清除可能已被轮末评估置位的单航班守卫，避免该会话永久阻塞后续更新。
		if updateSessionMemory {
			s.clearSessionMemoryUpdating(conv.ID)
		}
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
		if len(history) > 0 {
			if err := s.Store.UpsertConversationModelHistory(ctx, conv.ID, user.ID, history); err != nil {
				logger.Warn(fmt.Sprintf("model history: persist failed conversation=%s: %v", conv.ID, err))
			}
		}
		if memoryOn {
			s.extractMemories(ctx, user, history)
			if updateSessionMemory {
				ok := s.updateConversationMemory(ctx, conv, user, history)
				if ok {
					s.commitSessionMemoryBreakpoint(ctx, conv.ID, baselineTokens, lastMessageID(history))
				}
				s.clearSessionMemoryUpdating(conv.ID)
			}
			s.maybeRunConsolidation(ctx, user)
		} else if updateSessionMemory {
			s.clearSessionMemoryUpdating(conv.ID)
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

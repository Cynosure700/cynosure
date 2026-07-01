package compression

import (
	"context"
	"fmt"
	"strings"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

const reactiveCompactStrategyName = "reactive_compact"

const (
	// reactiveSummaryTargetTokens 是被动压缩的摘要预算；
	// 比 FullHistorySummarization 的 8K 更激进。
	reactiveSummaryTargetTokens = 2 * 1000
	// reactiveMaxStrips 是单次被动压缩内部最多剥离的次数。
	reactiveMaxStrips = 3
)

// ReactiveCompactStrategy 在上下文溢出（HTTP 413）时被带外调用，是一次「单次压缩」：
// 内部按对话轮从最旧端渐进剥离，每次剥离剩余对话轮的 20%（至少 1 轮），对「已有摘要 +
// 本次新剥离轮」重摘要，并把结果重组为「系统前导 + 会话摘要 + 未剥离轮」。每次剥离后
// 检查 token 阈值：一旦重组后的历史回到上下文预算内，就停止继续剥离并返回该组合（成功）；
// 最多剥离 3 次——若 3 次后（或已无可剥离对话轮时）仍超出上下文预算，则**返回错误**。
// 保留的消息会丢弃 reasoning_content。摘要调用自身失败（含 413）时也直接返回错误。
type ReactiveCompactStrategy struct{}

func (s *ReactiveCompactStrategy) Name() string {
	return reactiveCompactStrategyName
}

func (s *ReactiveCompactStrategy) Apply(ctx context.Context, req *Request) error {
	if req.Estimator == nil || req.Summarizer == nil {
		return nil
	}
	budget := req.Estimator.ContextTokenBudget()
	summary := ""

	// 单次压缩内部最多剥离 reactiveMaxStrips 次；满足 token 阈值即提前成功返回。
	for strips := 0; strips < reactiveMaxStrips; strips++ {
		turns := splitConversationTurns(req.RequestHistory)
		if len(turns) <= 1 {
			// 已无可剥离的对话轮（无内容或只剩 1 轮）：无法继续压缩，跳出循环由下方
			// 统一按 token 阈值判定成功/报错。
			break
		}

		// 本次新剥离轮数 = ceil(当前对话轮 × 20%)，至少 1，且至少给未剥离侧留 1 轮。
		strip := stripCount(len(turns))
		strippedTurnMsgs := flattenTurns(turns[:strip])

		// 重摘要：把已有摘要与本次新剥离轮一起喂给摘要器，产出一段连贯的整合摘要。
		result, err := req.Summarizer(ctx, SummaryRequest{
			Conversation: req.Conversation,
			User:         req.User,
			History:      buildReactiveSummaryInput(summary, strippedTurnMsgs),
			TargetTokens: reactiveSummaryTargetTokens,
			Aggressive:   true,
		})
		if err != nil {
			// 摘要调用失败（含 413）→ 直接返回错误，不做收缩重试。
			return fmt.Errorf("reactive compact summarize: %w", err)
		}
		summary = result.Summary

		// 未剥离的对话轮（逐字保留），修复 tool_call 配对并丢弃 reasoning_content。
		kept := repairToolCallBoundaries(flattenTurns(turns[strip:]))
		dropReasoningContent(kept)
		req.RequestHistory = buildSummaryHistory(summary, kept)

		// 满足 token 阈值即成功停止：返回当前组合结果
		//（系统前导 + 会话摘要 + 未剥离轮）。
		if req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools) <= budget {
			return nil
		}
	}

	// 已达最多剥离次数（或无更多对话轮可剥），仍超出上下文预算 → 返回错误。
	if req.Estimator.EstimateRequestTokens(req.SystemPrompt, req.RequestHistory, req.Tools) > budget {
		return fmt.Errorf("reactive compact: still over context budget after %d strips", reactiveMaxStrips)
	}
	return nil
}

// stripCount 返回本次要剥离的对话轮数：ceil(total × 20%)，至少 1，
// 且至多 total-1（给未剥离侧至少留一个对话轮）。
func stripCount(total int) int {
	if total <= 1 {
		return 0
	}
	// 整数向上取整：ceil(total * 0.2) == ceil(total/5) == (total*2 + 9) / 10。
	strip := (total*2 + 9) / 10
	if strip < 1 {
		strip = 1
	}
	if strip > total-1 {
		strip = total - 1
	}
	return strip
}

// buildReactiveSummaryInput 组装重摘要的输入：把已有摘要（若有）作为前置上下文消息，
// 拼接本次新剥离轮的原始消息，交给摘要器整合为一段新摘要。
func buildReactiveSummaryInput(priorSummary string, strippedTurns []storage.Message) []storage.Message {
	if strings.TrimSpace(priorSummary) == "" {
		return strippedTurns
	}
	head := storage.Message{
		Role:    "user",
		Content: "<conversation-summary>\n" + priorSummary + "\n</conversation-summary>",
	}
	return append([]storage.Message{head}, strippedTurns...)
}

// flattenTurns 把对话轮列表展平为一条消息序列。
func flattenTurns(turns [][]storage.Message) []storage.Message {
	var out []storage.Message
	for _, turn := range turns {
		out = append(out, turn...)
	}
	return out
}

// splitConversationTurns 把消息序列按「对话轮」切分。一个对话轮从一条真实用户消息开始，
// 到下一条真实用户消息之前结束（含中间的 assistant/tool/内部注入消息）。
//   - 内部注入的 user 消息（以 <system-reminder> 包裹）不作为新轮起点，归入当前轮。
//   - 开头出现的非用户消息（system 前导 / <conversation-summary> 前缀等，即上一次压缩
//     产物）不属于任何对话轮，被丢弃（它们会被本次重摘要重新整合）。
func splitConversationTurns(history []storage.Message) [][]storage.Message {
	var turns [][]storage.Message
	var current []storage.Message
	for _, msg := range history {
		if isRealUserMessage(msg) {
			if len(current) > 0 {
				turns = append(turns, current)
			}
			current = []storage.Message{msg}
			continue
		}
		if len(current) == 0 {
			// 尚未进入任何对话轮：忽略开头的 system/summary 前缀等非轮消息。
			continue
		}
		current = append(current, msg)
	}
	if len(current) > 0 {
		turns = append(turns, current)
	}
	return turns
}

// isRealUserMessage 判定一条消息是否为真实用户输入（对话轮起点）。
// role=user 且内容不是内部注入的 <system-reminder> 包裹、也不是 <conversation-summary>
// 前缀时，视为真实用户消息。
func isRealUserMessage(msg storage.Message) bool {
	if msg.Role != "user" {
		return false
	}
	trimmed := strings.TrimSpace(msg.Content)
	if strings.HasPrefix(trimmed, "<system-reminder>") {
		return false
	}
	if strings.HasPrefix(trimmed, "<conversation-summary>") {
		return false
	}
	return true
}

// dropReasoningContent 就地从保留的消息中剥除 reasoning_content 以回收上下文空间；
// 它对于续写从不需要。
func dropReasoningContent(messages []storage.Message) {
	for i := range messages {
		messages[i].ReasoningContent = ""
	}
}

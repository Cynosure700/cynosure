package compression

import (
	"context"
	"strings"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

const messageWindowCompressionStrategyName = "message_window_compression"
const enabled = false

// MessageWindowCompressionStrategy 在请求副本超出窗口上限时裁剪中间历史，
// 保留首部与尾部，随后修复因裁剪而产生的悬空 OpenAI tool_call / tool_result 配对。
type MessageWindowCompressionStrategy struct{}

func (s *MessageWindowCompressionStrategy) Name() string {
	return messageWindowCompressionStrategyName
}

func (s *MessageWindowCompressionStrategy) Apply(ctx context.Context, req *Request) error {
	if !enabled {
		return nil
	}
	history := req.RequestHistory

	// 定位最近一条用户消息（从尾部向前扫描）。
	lastUser := -1
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			lastUser = i
			break
		}
	}

	// 以当前回合的消息数量（从最近一条用户消息到尾部）作为触发条件，不含更早的历史。
	// 若没有用户消息，则回退为统计整个历史。
	turnStart := lastUser
	if turnStart < 0 {
		turnStart = 0
	}
	if len(history)-turnStart <= messageWindowLimit {
		return nil
	}

	// head 保留最近一条用户消息及其后的 2 条消息（共 3 条）。
	// 若没有用户消息，则回退为本回合的第一条消息及其后的 2 条
	//（系统提示词由调用方单独注入）。
	headStart := lastUser
	if headStart < 0 {
		headStart = 0
	}
	head := history[headStart : headStart+messageWindowHead]
	tail := history[len(history)-messageWindowTail:]

	windowed := make([]storage.Message, 0, messageWindowHead+messageWindowTail)
	windowed = append(windowed, head...)
	windowed = append(windowed, tail...)
	req.RequestHistory = repairToolCallBoundaries(windowed)
	return nil
}

// RepairToolCallBoundaries 移除会让 OpenAI 请求非法的 tool_call / tool
// 片段。OpenAI 要求 role=tool 消息必须紧跟在带 tool_calls 的 assistant
// 消息之后，并且只能响应该 assistant 声明的调用。
func RepairToolCallBoundaries(history []storage.Message) []storage.Message {
	repaired := make([]storage.Message, 0, len(history))
	for i := 0; i < len(history); i++ {
		msg := history[i]
		if msg.Role == "tool" {
			// 丢弃没有紧邻 assistant tool_calls 的孤立 tool 消息。
			continue
		}
		if msg.Role != "assistant" || len(msg.ToolCalls) == 0 {
			repaired = append(repaired, msg)
			continue
		}

		allowed := make(map[string]storage.MessageToolCall, len(msg.ToolCalls))
		for _, call := range msg.ToolCalls {
			allowed[call.ID] = call
		}

		j := i + 1
		matched := make(map[string]struct{}, len(msg.ToolCalls))
		toolMessages := make([]storage.Message, 0, len(msg.ToolCalls))
		for j < len(history) && history[j].Role == "tool" {
			toolMsg := history[j]
			if _, ok := allowed[toolMsg.ToolCallID]; ok {
				if _, duplicate := matched[toolMsg.ToolCallID]; !duplicate {
					matched[toolMsg.ToolCallID] = struct{}{}
					toolMessages = append(toolMessages, toolMsg)
				}
			}
			j++
		}

		keptCalls := make([]storage.MessageToolCall, 0, len(toolMessages))
		for _, call := range msg.ToolCalls {
			if _, ok := matched[call.ID]; ok {
				keptCalls = append(keptCalls, call)
			}
		}
		msg.ToolCalls = keptCalls
		if len(keptCalls) == 0 {
			if strings.TrimSpace(msg.Content) != "" || strings.TrimSpace(msg.ReasoningContent) != "" {
				repaired = append(repaired, msg)
			}
			i = j - 1
			continue
		}

		repaired = append(repaired, msg)
		repaired = append(repaired, toolMessages...)
		i = j - 1
	}
	return repaired
}

func repairToolCallBoundaries(history []storage.Message) []storage.Message {
	return RepairToolCallBoundaries(history)
}

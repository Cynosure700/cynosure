package compression

import (
	"context"
	"strings"

	"nano_cc/internal/agent/storage"
)

const messageWindowCompressionStrategyName = "message_window_compression"

// MessageWindowCompressionStrategy 在请求副本超出窗口上限时裁剪中间历史，
// 保留首部与尾部，随后修复因裁剪而产生的悬空 OpenAI tool_call / tool_result 配对。
type MessageWindowCompressionStrategy struct{}

func (s *MessageWindowCompressionStrategy) Name() string {
	return messageWindowCompressionStrategyName
}

func (s *MessageWindowCompressionStrategy) Apply(ctx context.Context, req *Request) error {
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

// repairToolCallBoundaries 移除孤立的 tool 消息（没有前置的 assistant tool_call），
// 并清除丢失了对应 tool 结果的 assistant tool_calls。
func repairToolCallBoundaries(history []storage.Message) []storage.Message {
	// 收集 assistant 消息仍然暴露的 tool_call id。
	assistantCallIDs := make(map[string]struct{})
	for _, msg := range history {
		if msg.Role == "assistant" {
			for _, call := range msg.ToolCalls {
				assistantCallIDs[call.ID] = struct{}{}
			}
		}
	}
	// 收集窗口中存在的 tool 结果 id。
	toolResultIDs := make(map[string]struct{})
	for _, msg := range history {
		if msg.Role == "tool" && msg.ToolCallID != "" {
			toolResultIDs[msg.ToolCallID] = struct{}{}
		}
	}

	repaired := make([]storage.Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "tool":
			// 丢弃没有匹配 assistant 调用的孤立 tool 消息。
			if _, ok := assistantCallIDs[msg.ToolCallID]; !ok {
				continue
			}
			repaired = append(repaired, msg)
		case "assistant":
			if len(msg.ToolCalls) > 0 {
				kept := msg.ToolCalls[:0:0]
				for _, call := range msg.ToolCalls {
					if _, ok := toolResultIDs[call.ID]; ok {
						kept = append(kept, call)
					}
				}
				msg.ToolCalls = kept
				if len(kept) == 0 && strings.TrimSpace(msg.Content) == "" && strings.TrimSpace(msg.ReasoningContent) == "" {
					continue
				}
			}
			repaired = append(repaired, msg)
		default:
			repaired = append(repaired, msg)
		}
	}
	return repaired
}

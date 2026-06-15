package compression

import (
	"context"
	"strings"

	"nano_cc/internal/agent/storage"
)

const messageWindowCompressionStrategyName = "message_window_compression"

// MessageWindowCompressionStrategy trims middle history when the request copy
// exceeds the window limit, keeping the head and tail, then repairs any
// dangling OpenAI tool_call / tool_result pairs created by the cut.
type MessageWindowCompressionStrategy struct{}

func (s *MessageWindowCompressionStrategy) Name() string {
	return messageWindowCompressionStrategyName
}

func (s *MessageWindowCompressionStrategy) Apply(ctx context.Context, req *Request) error {
	history := req.RequestHistory
	if len(history) <= messageWindowLimit {
		return nil
	}
	head := history[:messageWindowHead]
	tail := history[len(history)-messageWindowTail:]
	windowed := make([]storage.Message, 0, messageWindowHead+messageWindowTail)
	windowed = append(windowed, head...)
	windowed = append(windowed, tail...)
	req.RequestHistory = repairToolCallBoundaries(windowed)
	return nil
}

// repairToolCallBoundaries removes orphan tool messages (no preceding assistant
// tool_call) and clears assistant tool_calls that lost their tool results.
func repairToolCallBoundaries(history []storage.Message) []storage.Message {
	// Collect tool_call ids that the assistant messages still expose.
	assistantCallIDs := make(map[string]struct{})
	for _, msg := range history {
		if msg.Role == "assistant" {
			for _, call := range msg.ToolCalls {
				assistantCallIDs[call.ID] = struct{}{}
			}
		}
	}
	// Collect tool result ids present in the window.
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
			// Drop orphan tool messages without a matching assistant call.
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

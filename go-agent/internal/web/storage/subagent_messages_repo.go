package storage

import (
	"context"
	"encoding/json"
)

func (s *Store) CreateSubagentMessage(ctx context.Context, message SubagentMessage) error {
	toolCallsJSON, err := json.Marshal(message.ToolCalls)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO subagent_messages (id, run_id, parent_tool_call_id, conversation_id, user_id, sequence_no, role, content, reasoning_content, tool_call_id, tool_calls_json, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW())
	`, message.ID, message.RunID, message.ParentToolCallID, message.ConversationID, message.UserID, message.SequenceNo, message.Role, message.Content, message.ReasoningContent, message.ToolCallID, string(toolCallsJSON))
	return err
}

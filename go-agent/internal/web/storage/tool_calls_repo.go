package storage

import "context"

func (s *Store) CreateToolCall(ctx context.Context, tc ToolCall) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO tool_calls (id, conversation_id, user_id, tool_name, status, summary, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`, tc.ID, tc.ConversationID, tc.UserID, tc.ToolName, tc.Status, tc.Summary)
	return err
}

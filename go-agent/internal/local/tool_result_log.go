package local

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"nano_cc/internal/agent/storage"
)

func (s *Store) AppendToolResultLog(ctx context.Context, entry storage.ToolResultLogEntry) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.workspaceRoot == "" {
		return nil
	}
	conv := s.conversations[entry.ConversationID]
	sessionID := entry.SessionID
	if sessionID == "" {
		sessionID = conv.SessionID
	}
	if !validSessionID(sessionID) {
		return fmt.Errorf("invalid session_id: %q", sessionID)
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}
	path := filepath.Join(s.workspaceRoot, "task_outputs", sessionID, "tools.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	unlock, err := lockFile(path + ".lock")
	if err != nil {
		return err
	}
	defer unlock()
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(renderToolResultLogEntry(sessionID, entry))
	return err
}

func lockFile(path string) (func(), error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	if err := syscall.Flock(int(file.Fd()), syscall.LOCK_EX); err != nil {
		_ = file.Close()
		return nil, err
	}
	return func() {
		_ = syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
		_ = file.Close()
	}, nil
}

func renderToolResultLogEntry(sessionID string, entry storage.ToolResultLogEntry) string {
	var b strings.Builder
	b.WriteString("## ")
	b.WriteString(entry.CreatedAt.UTC().Format(time.RFC3339))
	b.WriteString(" · ")
	b.WriteString(entry.ToolName)
	b.WriteString(" · ")
	b.WriteString(entry.Status)
	b.WriteString("\n\n")
	b.WriteString("- conversation_id: ")
	b.WriteString(entry.ConversationID)
	b.WriteString("\n- session_id: ")
	b.WriteString(sessionID)
	b.WriteString("\n- tool_call_id: ")
	b.WriteString(entry.ToolCallID)
	b.WriteString("\n\n")
	b.WriteString("### Arguments\n\n")
	b.WriteString(markdownFence("json", entry.RawArgs))
	b.WriteString("\n\n")
	b.WriteString("### Result\n\n")
	b.WriteString(markdownFence("text", entry.Result))
	if strings.TrimSpace(entry.AuditSummary) != "" {
		b.WriteString("\n\n### Audit\n\n")
		b.WriteString(markdownFence("json", entry.AuditSummary))
	}
	b.WriteString("\n\n")
	return b.String()
}

func markdownFence(language, body string) string {
	fence := strings.Repeat("`", longestBacktickRun(body)+1)
	if len(fence) < 3 {
		fence = "```"
	}
	return fence + language + "\n" + body + "\n" + fence
}

func longestBacktickRun(body string) int {
	longest := 0
	current := 0
	for _, r := range body {
		if r == '`' {
			current++
			if current > longest {
				longest = current
			}
			continue
		}
		current = 0
	}
	return longest
}

package local

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
)

type SessionHistoryStore struct {
	dir string
}

type SessionHistoryDocument struct {
	Version        int               `json:"version"`
	SessionID      string            `json:"session_id"`
	ConversationID string            `json:"conversation_id"`
	UserID         string            `json:"user_id"`
	WorkspaceRoot  string            `json:"workspace_root"`
	Title          string            `json:"title"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	Messages       []storage.Message `json:"messages"`
}

func NewSessionHistoryStore(workspaceRoot string) (*SessionHistoryStore, error) {
	dir, err := config.CynosureSessionDir(workspaceRoot)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("create session dir: %w", err)
	}
	return &SessionHistoryStore{dir: dir}, nil
}

func (s *SessionHistoryStore) SaveHistory(ctx context.Context, workspaceRoot string, conversation storage.Conversation, messages []storage.Message) error {
	return s.save(ctx, workspaceRoot, conversation, messages, "history")
}

func (s *SessionHistoryStore) LoadHistory(ctx context.Context, sessionID string) (SessionHistoryDocument, error) {
	return s.load(ctx, sessionID, "history")
}

func (s *SessionHistoryStore) SaveModelHistory(ctx context.Context, workspaceRoot string, conversation storage.Conversation, messages []storage.Message) error {
	return s.save(ctx, workspaceRoot, conversation, messages, "model_history")
}

func (s *SessionHistoryStore) LoadModelHistory(ctx context.Context, sessionID string) (SessionHistoryDocument, error) {
	return s.load(ctx, sessionID, "model_history")
}

func (s *SessionHistoryStore) ListSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error) {
	currentWorkspace, err := normalizeWorkspacePath(workspaceRoot)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var result []storage.ResumableSession
	for _, entry := range entries {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		if !entry.IsDir() || !validSessionID(entry.Name()) {
			continue
		}
		doc, err := s.LoadHistory(ctx, entry.Name())
		if err != nil {
			if os.IsNotExist(err) || err == sql.ErrNoRows {
				continue
			}
			return nil, err
		}
		docWorkspace, err := normalizeWorkspacePath(doc.WorkspaceRoot)
		if err != nil || docWorkspace != currentWorkspace {
			continue
		}
		result = append(result, storage.ResumableSession{
			SessionID:      doc.SessionID,
			ConversationID: doc.ConversationID,
			WorkspaceRoot:  doc.WorkspaceRoot,
			Title:          doc.Title,
			MessageCount:   len(doc.Messages),
			CreatedAt:      doc.CreatedAt,
			UpdatedAt:      doc.UpdatedAt,
		})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	return result, nil
}

func (s *SessionHistoryStore) save(ctx context.Context, workspaceRoot string, conversation storage.Conversation, messages []storage.Message, name string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !validSessionID(conversation.SessionID) {
		return fmt.Errorf("invalid session_id: %q", conversation.SessionID)
	}
	now := time.Now()
	createdAt := conversation.CreatedAt
	if createdAt.IsZero() {
		createdAt = now
	}
	updatedAt := conversation.UpdatedAt
	if updatedAt.IsZero() || updatedAt.Before(now) {
		updatedAt = now
	}
	doc := SessionHistoryDocument{
		Version:        1,
		SessionID:      conversation.SessionID,
		ConversationID: conversation.ID,
		UserID:         conversation.UserID,
		WorkspaceRoot:  workspaceRoot,
		Title:          conversation.Title,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		Messages:       cloneMessages(messages),
	}
	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	return atomicWriteFile(s.path(conversation.SessionID, name), data, 0o644)
}

func (s *SessionHistoryStore) load(ctx context.Context, sessionID, name string) (SessionHistoryDocument, error) {
	if err := ctx.Err(); err != nil {
		return SessionHistoryDocument{}, err
	}
	if !validSessionID(sessionID) {
		return SessionHistoryDocument{}, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	data, err := os.ReadFile(s.path(sessionID, name))
	if err != nil {
		return SessionHistoryDocument{}, err
	}
	var doc SessionHistoryDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return SessionHistoryDocument{}, fmt.Errorf("parse session %s %s: %w", sessionID, name, err)
	}
	if doc.SessionID == "" {
		doc.SessionID = sessionID
	}
	if doc.SessionID != sessionID {
		return SessionHistoryDocument{}, fmt.Errorf("session_id mismatch: path=%s document=%s", sessionID, doc.SessionID)
	}
	return doc, nil
}

func (s *SessionHistoryStore) path(sessionID, name string) string {
	return filepath.Join(s.dir, sessionID, name)
}

func validSessionID(sessionID string) bool {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" || sessionID == "." || sessionID == ".." || strings.Contains(sessionID, "..") {
		return false
	}
	for _, r := range sessionID {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}

func normalizeWorkspacePath(path string) (string, error) {
	if strings.TrimSpace(path) == "" {
		path = "."
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	clean := filepath.Clean(abs)
	if evaluated, err := filepath.EvalSymlinks(clean); err == nil {
		clean = evaluated
	}
	return clean, nil
}

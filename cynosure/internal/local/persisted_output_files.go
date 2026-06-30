package local

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
	"github.com/Cynosure700/cynosure/cynosure/internal/config"
)

type persistedOutputMetadata struct {
	Version        int       `json:"version"`
	ID             string    `json:"id"`
	SessionID      string    `json:"session_id"`
	ConversationID string    `json:"conversation_id"`
	UserID         string    `json:"user_id"`
	MessageID      string    `json:"message_id"`
	ToolCallID     string    `json:"tool_call_id"`
	Kind           string    `json:"kind"`
	Strategy       string    `json:"strategy"`
	OriginalBytes  int       `json:"original_bytes"`
	ContentSHA256  string    `json:"content_sha256"`
	Preview        string    `json:"preview"`
	ContentFile    string    `json:"content_file"`
	CreatedAt      time.Time `json:"created_at"`
}

func persistOutputToWorkspace(ctx context.Context, workspaceRoot, sessionID string, output storage.PersistedOutput) (storage.PersistedOutput, error) {
	if err := ctx.Err(); err != nil {
		return output, err
	}
	if !validSessionID(sessionID) {
		return output, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	if !validStoredOutputID(output.ID) {
		return output, fmt.Errorf("invalid persisted output id: %q", output.ID)
	}
	if output.CreatedAt.IsZero() {
		output.CreatedAt = time.Now()
	}
	if output.OriginalBytes == 0 {
		output.OriginalBytes = len([]byte(output.Content))
	}
	sum := sha256.Sum256([]byte(output.Content))
	computedSHA := hex.EncodeToString(sum[:])
	if output.ContentSHA256 == "" {
		output.ContentSHA256 = computedSHA
	}

	contentFile := persistedOutputContentFileName(output.ID)
	metadataFile := persistedOutputMetadataFileName(output.ID)
	dir, err := persistedOutputDir(workspaceRoot, sessionID)
	if err != nil {
		return output, err
	}
	metadata := persistedOutputMetadata{
		Version:        1,
		ID:             output.ID,
		SessionID:      sessionID,
		ConversationID: output.ConversationID,
		UserID:         output.UserID,
		MessageID:      output.MessageID,
		ToolCallID:     output.ToolCallID,
		Kind:           output.Kind,
		Strategy:       output.Strategy,
		OriginalBytes:  output.OriginalBytes,
		ContentSHA256:  output.ContentSHA256,
		Preview:        output.Preview,
		ContentFile:    contentFile,
		CreatedAt:      output.CreatedAt,
	}
	metadataBytes, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return output, err
	}
	if err := atomicWriteFile(filepath.Join(dir, contentFile), []byte(output.Content), 0o644); err != nil {
		return output, err
	}
	if err := atomicWriteFile(filepath.Join(dir, metadataFile), metadataBytes, 0o644); err != nil {
		return output, err
	}
	return output, nil
}

func loadPersistedOutputFromWorkspace(ctx context.Context, workspaceRoot, sessionID, id string) (storage.PersistedOutput, error) {
	if err := ctx.Err(); err != nil {
		return storage.PersistedOutput{}, err
	}
	if !validSessionID(sessionID) {
		return storage.PersistedOutput{}, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	if !validStoredOutputID(id) {
		return storage.PersistedOutput{}, fmt.Errorf("invalid persisted output id: %q", id)
	}
	dir, err := persistedOutputDir(workspaceRoot, sessionID)
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	metadataBytes, err := os.ReadFile(filepath.Join(dir, persistedOutputMetadataFileName(id)))
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	var metadata persistedOutputMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return storage.PersistedOutput{}, err
	}
	if metadata.ID != id || metadata.SessionID != sessionID || metadata.Kind != "tool_result" {
		return storage.PersistedOutput{}, fmt.Errorf("persisted output metadata mismatch for %q", id)
	}
	if filepath.Base(metadata.ContentFile) != metadata.ContentFile || metadata.ContentFile == "" {
		return storage.PersistedOutput{}, fmt.Errorf("invalid persisted output content file: %q", metadata.ContentFile)
	}
	contentBytes, err := os.ReadFile(filepath.Join(dir, metadata.ContentFile))
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	sum := sha256.Sum256(contentBytes)
	actualSHA := hex.EncodeToString(sum[:])
	if metadata.ContentSHA256 != "" && actualSHA != metadata.ContentSHA256 {
		return storage.PersistedOutput{}, fmt.Errorf("persisted output %q sha256 mismatch", id)
	}
	return storage.PersistedOutput{
		ID:             metadata.ID,
		ConversationID: metadata.ConversationID,
		UserID:         metadata.UserID,
		MessageID:      metadata.MessageID,
		ToolCallID:     metadata.ToolCallID,
		Kind:           metadata.Kind,
		Strategy:       metadata.Strategy,
		OriginalBytes:  metadata.OriginalBytes,
		ContentSHA256:  metadata.ContentSHA256,
		Content:        string(contentBytes),
		Preview:        metadata.Preview,
		CreatedAt:      metadata.CreatedAt,
	}, nil
}

func findPersistedOutputByMessageHashInWorkspace(ctx context.Context, workspaceRoot, sessionID, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error) {
	if err := ctx.Err(); err != nil {
		return storage.PersistedOutput{}, err
	}
	if !validSessionID(sessionID) {
		return storage.PersistedOutput{}, fmt.Errorf("invalid session_id: %q", sessionID)
	}
	dir, err := persistedOutputDir(workspaceRoot, sessionID)
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.PersistedOutput{}, sql.ErrNoRows
		}
		return storage.PersistedOutput{}, err
	}
	for _, entry := range entries {
		if ctx.Err() != nil {
			return storage.PersistedOutput{}, ctx.Err()
		}
		name := entry.Name()
		if entry.IsDir() || !strings.HasSuffix(name, ".json") {
			continue
		}
		metadataBytes, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return storage.PersistedOutput{}, err
		}
		var metadata persistedOutputMetadata
		if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
			return storage.PersistedOutput{}, err
		}
		if metadata.SessionID != sessionID || metadata.ConversationID != conversationID || metadata.UserID != userID || metadata.MessageID != messageID || metadata.ToolCallID != toolCallID || metadata.Strategy != strategy || metadata.ContentSHA256 != contentSHA256 {
			continue
		}
		return loadPersistedOutputFromWorkspace(ctx, workspaceRoot, sessionID, metadata.ID)
	}
	return storage.PersistedOutput{}, sql.ErrNoRows
}

func persistedOutputDir(workspaceRoot, sessionID string) (string, error) {
	if !validSessionID(sessionID) {
		return "", fmt.Errorf("invalid session_id: %q", sessionID)
	}
	dir, err := taskOutputsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, workspaceDirName(workspaceRoot), sessionID, "tool-results"), nil
}

func taskOutputsDir() (string, error) {
	return config.CynosureTaskOutputsDir()
}

func persistedOutputContentFileName(id string) string {
	return id + ".txt"
}

func persistedOutputMetadataFileName(id string) string {
	return id + ".json"
}

func validStoredOutputID(id string) bool {
	id = strings.TrimSpace(id)
	if id == "" || id == "." || id == ".." || strings.Contains(id, "..") {
		return false
	}
	for _, r := range id {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}

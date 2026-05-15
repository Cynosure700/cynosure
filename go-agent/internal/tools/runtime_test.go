package tools

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"nano_cc/internal/safety"
)

func TestHandleRead_UsesWorkspaceRootForRelativePath(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "note.txt")
	if err := os.WriteFile(file, []byte("hello workspace"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	result, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"path": "note.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello workspace" {
		t.Fatalf("expected workspace file content, got %q", result)
	}
}

func TestHandleWrite_UsesWorkspaceRootForRelativePath(t *testing.T) {
	root := t.TempDir()
	_, err := handleWrite(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"path": "nested/out.txt", "content": "created"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(root, "nested", "out.txt"))
	if err != nil {
		t.Fatalf("expected file in workspace: %v", err)
	}
	if string(data) != "created" {
		t.Fatalf("expected written content, got %q", string(data))
	}
}

func TestHandleEdit_UsesWorkspaceRootForRelativePath(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "edit.txt")
	if err := os.WriteFile(file, []byte("before"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := handleEdit(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"path": "edit.txt", "old_text": "before", "new_text": "after"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(file)
	if err != nil {
		t.Fatalf("read edited file: %v", err)
	}
	if string(data) != "after" {
		t.Fatalf("expected edited content, got %q", string(data))
	}
}

func TestHandleBash_UsesWorkspaceRootAsDefaultDir(t *testing.T) {
	root := t.TempDir()
	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"command": "pwd"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Clean(result) != filepath.Clean(root) {
		t.Fatalf("expected pwd to run in workspace %q, got %q", root, result)
	}
}

func TestSafePathFromRoot_RejectsPathEscape(t *testing.T) {
	root := t.TempDir()
	_, err := safety.SafePathFromRoot(root, "../escape.txt")
	if err == nil {
		t.Fatalf("expected path escape to be rejected")
	}
}

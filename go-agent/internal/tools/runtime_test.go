package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
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

func TestHandleBash_UsesDeploymentWorkspaceAsDefaultDir(t *testing.T) {
	appRoot := t.TempDir()
	workspace := filepath.Join(appRoot, "output", "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("mkdir deployment workspace: %v", err)
	}
	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace}), map[string]any{"command": "pwd"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Clean(result) != filepath.Clean(workspace) {
		t.Fatalf("expected pwd to run in deployment workspace %q, got %q", workspace, result)
	}
}

func TestHandleBash_RejectsAbsolutePathOutsideWorkspace(t *testing.T) {
	root := t.TempDir()
	_, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"command": "cat /tmp/outside.txt"})
	if err == nil {
		t.Fatalf("expected absolute path outside workspace to be rejected")
	}
	if !strings.Contains(err.Error(), "command path escapes workspace") {
		t.Fatalf("expected workspace escape error, got %v", err)
	}
}

func TestSafePathFromRoot_RejectsPathEscape(t *testing.T) {
	root := t.TempDir()
	_, err := safety.SafePathFromRoot(root, "../escape.txt")
	if err == nil {
		t.Fatalf("expected path escape to be rejected")
	}
}

func TestHandleWrite_DoesNotTouchDeploymentCommandDir(t *testing.T) {
	workspace := t.TempDir()
	commandDir := filepath.Join(workspace, "..", "cmd")
	ctx := WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CommandScriptDir: commandDir})

	_, err := handleWrite(ctx, map[string]any{"path": filepath.Join("..", "cmd", "script.py"), "content": "print('x')"})
	if err == nil {
		t.Fatalf("expected write escaping workspace to be rejected")
	}
}

func TestValidatedWorkspaceRootFromContext_RequiresWorkspaceRoot(t *testing.T) {
	_, err := validatedWorkspaceRootFromContext(context.Background())
	if err == nil {
		t.Fatalf("expected missing workspace root to be rejected")
	}
	if !strings.Contains(err.Error(), "workspace root is required") {
		t.Fatalf("expected missing workspace error, got %v", err)
	}
}

func TestValidatedWorkspaceRootFromContext_RejectsNonDirectory(t *testing.T) {
	root := filepath.Join(t.TempDir(), "workspace.txt")
	if err := os.WriteFile(root, []byte("not a dir"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	_, err := validatedWorkspaceRootFromContext(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}))
	if err == nil {
		t.Fatalf("expected non-directory workspace root to be rejected")
	}
	if !strings.Contains(err.Error(), "workspace root is not a directory") {
		t.Fatalf("expected non-directory workspace error, got %v", err)
	}
}

func TestHandlers_RejectInvalidWorkspaceRoot(t *testing.T) {
	handlers := []struct {
		name string
		fn   ToolHandler
		args map[string]any
	}{
		{name: "bash", fn: handleBash, args: map[string]any{"command": "pwd"}},
		{name: "read", fn: handleRead, args: map[string]any{"path": "note.txt"}},
		{name: "write", fn: handleWrite, args: map[string]any{"path": "note.txt", "content": "hello"}},
		{name: "edit", fn: handleEdit, args: map[string]any{"path": "note.txt", "old_text": "before", "new_text": "after"}},
	}

	t.Run("missing workspace root", func(t *testing.T) {
		for _, tc := range handlers {
			t.Run(tc.name, func(t *testing.T) {
				_, err := tc.fn(context.Background(), tc.args)
				if err == nil {
					t.Fatalf("expected missing workspace root to be rejected")
				}
				if !strings.Contains(err.Error(), "workspace root is required") {
					t.Fatalf("expected missing workspace error, got %v", err)
				}
			})
		}
	})

	t.Run("non-directory workspace root", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "workspace.txt")
		if err := os.WriteFile(root, []byte("not a dir"), 0o644); err != nil {
			t.Fatalf("write fixture: %v", err)
		}

		for _, tc := range handlers {
			t.Run(tc.name, func(t *testing.T) {
				ctx := WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root})
				_, err := tc.fn(ctx, tc.args)
				if err == nil {
					t.Fatalf("expected non-directory workspace root to be rejected")
				}
				if !strings.Contains(err.Error(), "workspace root is not a directory") {
					t.Fatalf("expected non-directory workspace error, got %v", err)
				}
			})
		}
	})

	t.Run("unavailable workspace root", func(t *testing.T) {
		root := filepath.Join(t.TempDir(), "missing")

		for _, tc := range handlers {
			t.Run(tc.name, func(t *testing.T) {
				ctx := WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root})
				_, err := tc.fn(ctx, tc.args)
				if err == nil {
					t.Fatalf("expected unavailable workspace root to be rejected")
				}
				if !strings.Contains(err.Error(), "workspace root is unavailable") {
					t.Fatalf("expected unavailable workspace error, got %v", err)
				}
			})
		}
	})
}

package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleRead_UsesWorkspaceRootForRelativePath(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "note.txt")
	if err := os.WriteFile(file, []byte("hello workspace"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	result, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"file_path": "note.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "1\thello workspace" {
		t.Fatalf("expected workspace file content, got %q", result)
	}
}

func TestHandleRead_ReturnsLineNumbersByDefault(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "note.txt")
	if err := os.WriteFile(file, []byte("alpha\n\tbeta\n\nomega\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	result, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"file_path": "note.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "1\talpha\n2\t\tbeta\n3\t\n4\tomega"
	if result != want {
		t.Fatalf("expected numbered file content\nwant: %q\n got: %q", want, result)
	}
}

func TestHandleRead_AppliesOneBasedOffsetAndLimit(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "note.txt")
	if err := os.WriteFile(file, []byte("one\ntwo\nthree\nfour\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	result, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{
		"file_path": "note.txt",
		"offset":    float64(2),
		"limit":     float64(2),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "2\ttwo\n3\tthree"
	if result != want {
		t.Fatalf("expected offset-limited numbered content\nwant: %q\n got: %q", want, result)
	}
}

func TestHandleWrite_UsesWorkspaceRootForRelativePath(t *testing.T) {
	root := t.TempDir()
	result, err := handleWrite(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"file_path": "nested/out.txt", "content": "created"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "File created successfully at: nested/out.txt" {
		t.Fatalf("expected created result, got %q", result)
	}
	data, err := os.ReadFile(filepath.Join(root, "nested", "out.txt"))
	if err != nil {
		t.Fatalf("expected file in workspace: %v", err)
	}
	if string(data) != "created" {
		t.Fatalf("expected written content, got %q", string(data))
	}
}

func TestHandleWrite_ReportsExistingFileUpdate(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "out.txt")
	if err := os.WriteFile(file, []byte("before"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	result, err := handleWrite(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"file_path": "out.txt", "content": "after"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "The file out.txt has been updated successfully." {
		t.Fatalf("expected update result, got %q", result)
	}
}

func TestHandleEdit_UsesWorkspaceRootForRelativePath(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "edit.txt")
	if err := os.WriteFile(file, []byte("before"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	result, err := handleEdit(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"file_path": "edit.txt", "old_text": "before", "new_text": "after"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "The file edit.txt has been updated successfully." {
		t.Fatalf("expected edit result, got %q", result)
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

func TestHandleBash_DoesNotTruncateLargeOutput(t *testing.T) {
	root := t.TempDir()
	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"command": "printf '%*s' 50001 '' | tr ' ' x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 50001 {
		t.Fatalf("expected 50001 chars, got %d", len(result))
	}
}

func TestHandleBash_UsesConfiguredAppWorkspaceAsDefaultDir(t *testing.T) {
	appRoot := t.TempDir()
	workspace := filepath.Join(appRoot, "workspace")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatalf("mkdir app workspace: %v", err)
	}
	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace}), map[string]any{"command": "pwd"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Clean(result) != filepath.Clean(workspace) {
		t.Fatalf("expected pwd to run in app workspace %q, got %q", workspace, result)
	}
}

func TestHandleBash_AllowsCommonCommandsWithoutPathArguments(t *testing.T) {
	root := t.TempDir()
	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"command": "echo hello"})
	if err != nil {
		t.Fatalf("expected common command without path to be allowed: %v", err)
	}
	if result != "hello" {
		t.Fatalf("expected echo output, got %q", result)
	}
}

func TestHandleBash_AllowsOutsideWorkspacePath(t *testing.T) {
	root := t.TempDir()
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("outside ok"), 0o644); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}

	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"command": "cat " + outsideFile})
	if err != nil {
		t.Fatalf("expected outside path to be allowed: %v", err)
	}
	if result != "outside ok" {
		t.Fatalf("expected outside file content, got %q", result)
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
		{name: "read", fn: handleRead, args: map[string]any{"file_path": "note.txt"}},
		{name: "write", fn: handleWrite, args: map[string]any{"file_path": "note.txt", "content": "hello"}},
		{name: "edit", fn: handleEdit, args: map[string]any{"file_path": "note.txt", "old_text": "before", "new_text": "after"}},
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

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

func TestHandleBash_UsesCurrentWorkingDirWhenProvided(t *testing.T) {
	workspace := t.TempDir()
	skillDir := filepath.Join(workspace, "skills", "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CurrentWorkingDir: skillDir}), map[string]any{"command": "pwd"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if filepath.Clean(result) != filepath.Clean(skillDir) {
		t.Fatalf("expected pwd to run in skill dir %q, got %q", skillDir, result)
	}
}

func TestHandleRead_UsesCurrentWorkingDirForRelativePath(t *testing.T) {
	workspace := t.TempDir()
	skillDir := filepath.Join(workspace, "skills", "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	file := filepath.Join(skillDir, "note.txt")
	if err := os.WriteFile(file, []byte("hello skill dir"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	result, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CurrentWorkingDir: skillDir}), map[string]any{"path": "note.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "hello skill dir" {
		t.Fatalf("expected skill-dir file content, got %q", result)
	}
}

func TestHandleWrite_UsesCurrentWorkingDirForRelativePath(t *testing.T) {
	workspace := t.TempDir()
	skillDir := filepath.Join(workspace, "skills", "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	_, err := handleWrite(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CurrentWorkingDir: skillDir}), map[string]any{"path": "nested/out.txt", "content": "created from skill dir"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	data, err := os.ReadFile(filepath.Join(skillDir, "nested", "out.txt"))
	if err != nil {
		t.Fatalf("expected file in skill dir: %v", err)
	}
	if string(data) != "created from skill dir" {
		t.Fatalf("expected written content, got %q", string(data))
	}
}

func TestHandleEdit_UsesCurrentWorkingDirForRelativePath(t *testing.T) {
	workspace := t.TempDir()
	skillDir := filepath.Join(workspace, "skills", "demo-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	file := filepath.Join(skillDir, "edit.txt")
	if err := os.WriteFile(file, []byte("before"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	_, err := handleEdit(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CurrentWorkingDir: skillDir}), map[string]any{"path": "edit.txt", "old_text": "before", "new_text": "after"})
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

func TestHandleBash_RejectsSkillDirOutsideWorkspace(t *testing.T) {
	workspace := t.TempDir()
	outside := t.TempDir()
	_, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CurrentWorkingDir: outside}), map[string]any{"command": "pwd"})
	if err == nil {
		t.Fatalf("expected outside skill dir to be rejected")
	}
	if !strings.Contains(err.Error(), "current working directory escapes workspace") {
		t.Fatalf("expected skill-dir escape error, got %v", err)
	}
}

func TestHandleRead_RejectsSkillDirOutsideWorkspace(t *testing.T) {
	workspace := t.TempDir()
	outside := t.TempDir()
	_, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace, CurrentWorkingDir: outside}), map[string]any{"path": "note.txt"})
	if err == nil {
		t.Fatalf("expected outside skill dir to be rejected")
	}
	if !strings.Contains(err.Error(), "current working directory escapes workspace") {
		t.Fatalf("expected skill-dir escape error, got %v", err)
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

func TestHandleBash_AllowsOutsideWorkspacePathWhenConfigured(t *testing.T) {
	root := t.TempDir()
	outsideDir := t.TempDir()
	outsideFile := filepath.Join(outsideDir, "outside.txt")
	if err := os.WriteFile(outsideFile, []byte("outside ok"), 0o644); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}

	result, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root, AllowOutsideWorkspace: true}), map[string]any{"command": "cat " + outsideFile})
	if err != nil {
		t.Fatalf("expected configured outside path to be allowed: %v", err)
	}
	if result != "outside ok" {
		t.Fatalf("expected outside file content, got %q", result)
	}
}

func TestHandleBash_RejectsDangerousCommandByDefault(t *testing.T) {
	root := t.TempDir()
	_, err := handleBash(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"command": "rm temp.txt"})
	if err == nil {
		t.Fatalf("expected dangerous command to be rejected")
	}
	if !strings.Contains(err.Error(), "dangerous command blocked") {
		t.Fatalf("expected dangerous command error, got %v", err)
	}
}

func TestDangerousCommandPattern_IgnoresNonCommandArguments(t *testing.T) {
	if pattern, ok := dangerousCommandPattern("echo rm"); ok {
		t.Fatalf("expected non-command argument to be allowed, got dangerous pattern %q", pattern)
	}
	if pattern, ok := dangerousCommandPattern("echo ok; rm temp.txt"); !ok || pattern != "rm" {
		t.Fatalf("expected command after separator to be dangerous, got pattern=%q ok=%v", pattern, ok)
	}
}

func TestHandleRead_RejectsAbsolutePathOutsideWorkspace(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}

	_, err := handleRead(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"path": outside})
	if err == nil {
		t.Fatalf("expected absolute path outside workspace to be rejected")
	}
	if !strings.Contains(err.Error(), "path escapes workspace") {
		t.Fatalf("expected workspace escape error, got %v", err)
	}
}

func TestHandleWrite_RejectsAbsolutePathOutsideWorkspace(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")

	_, err := handleWrite(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"path": outside, "content": "blocked"})
	if err == nil {
		t.Fatalf("expected absolute path outside workspace to be rejected")
	}
	if !strings.Contains(err.Error(), "path escapes workspace") {
		t.Fatalf("expected workspace escape error, got %v", err)
	}
}

func TestHandleEdit_RejectsAbsolutePathOutsideWorkspace(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside.txt")
	if err := os.WriteFile(outside, []byte("before"), 0o644); err != nil {
		t.Fatalf("write outside fixture: %v", err)
	}

	_, err := handleEdit(WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root}), map[string]any{"path": outside, "old_text": "before", "new_text": "after"})
	if err == nil {
		t.Fatalf("expected absolute path outside workspace to be rejected")
	}
	if !strings.Contains(err.Error(), "path escapes workspace") {
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
	ctx := WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: workspace})

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

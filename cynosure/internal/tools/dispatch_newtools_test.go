package tools

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
)

// TestDispatchNewToolsEndToEnd exercises the real Dispatch entry point with a
// RuntimeEnv to confirm the new tools are wired, sandboxed and functional.
func TestDispatchNewToolsEndToEnd(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "main.go"), "package main\n// TODO fix\n")
	writeFile(t, filepath.Join(root, "notes.txt"), "alpha beta\n")

	ctx := WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root})

	// grep
	res, err := Dispatch(ctx, "grep", map[string]any{"pattern": "TODO", "output_mode": "files_with_matches"})
	if err != nil || !strings.Contains(res.Output, "main.go") {
		t.Fatalf("grep dispatch failed: out=%q err=%v", res.Output, err)
	}

	// glob
	res, err = Dispatch(ctx, "glob", map[string]any{"pattern": "**/*.go"})
	if err != nil || !strings.Contains(res.Output, "main.go") {
		t.Fatalf("glob dispatch failed: out=%q err=%v", res.Output, err)
	}

	// ls (absolute path required)
	res, err = Dispatch(ctx, "ls", map[string]any{"path": root})
	if err != nil || !strings.Contains(res.Output, "main.go") {
		t.Fatalf("ls dispatch failed: out=%q err=%v", res.Output, err)
	}

	// ls rejects relative path
	if _, err = Dispatch(ctx, "ls", map[string]any{"path": "."}); err == nil {
		t.Fatalf("ls should reject relative path")
	}

	// read_file
	res, err = Dispatch(ctx, "read_file", map[string]any{"file_path": filepath.Join(root, "notes.txt")})
	if err != nil || !strings.Contains(res.Output, "alpha beta") {
		t.Fatalf("read_file dispatch failed: out=%q err=%v", res.Output, err)
	}

	// write_file
	res, err = Dispatch(ctx, "write_file", map[string]any{
		"file_path": filepath.Join(root, "created.txt"),
		"content":   "created content\n",
	})
	if err != nil || !strings.Contains(res.Output, "File created successfully at:") {
		t.Fatalf("write_file dispatch failed: out=%q err=%v", res.Output, err)
	}

	// edit_file
	res, err = Dispatch(ctx, "edit_file", map[string]any{
		"file_path": filepath.Join(root, "created.txt"),
		"old_text":  "created",
		"new_text":  "updated",
	})
	if err != nil || !strings.Contains(res.Output, "has been updated successfully") {
		t.Fatalf("edit_file dispatch failed: out=%q err=%v", res.Output, err)
	}

	// multi_edit
	res, err = Dispatch(ctx, "multi_edit", map[string]any{
		"file_path": filepath.Join(root, "notes.txt"),
		"edits": []any{
			map[string]any{"old_string": "alpha", "new_string": "ALPHA"},
			map[string]any{"old_string": "beta", "new_string": "BETA"},
		},
	})
	if err != nil || !strings.Contains(res.Output, "Applied 2 edits") {
		t.Fatalf("multi_edit dispatch failed: out=%q err=%v", res.Output, err)
	}

	// web_search placeholder
	res, err = Dispatch(ctx, "web_search", map[string]any{"query": "go"})
	if err != nil || !strings.Contains(res.Output, "not configured") {
		t.Fatalf("web_search dispatch failed: out=%q err=%v", res.Output, err)
	}
}

package tools

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunMultiEditFromRoot(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "f.txt")
	writeFile(t, path, "alpha beta alpha gamma")

	out, err := RunMultiEditFromRoot(root, "f.txt", []Edit{
		{OldString: "beta", NewString: "BETA"},
		{OldString: "alpha", NewString: "X", ReplaceAll: true},
	})
	if err != nil {
		t.Fatalf("multi_edit: %v", err)
	}
	if out != "The file f.txt has been updated successfully." {
		t.Fatalf("multi_edit output = %q, want edit_file-compatible result", out)
	}
	data, _ := os.ReadFile(path)
	if string(data) != "X BETA X gamma" {
		t.Fatalf("unexpected content: %q", string(data))
	}
}

func TestHandleMultiEditReturnsOneEditCompatibleResultPerFile(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.txt"), "alpha beta")
	writeFile(t, filepath.Join(root, "b.txt"), "gamma delta")
	ctx := WithRuntimeEnv(context.Background(), RuntimeEnv{WorkspaceRoot: root})

	out, err := handleMultiEdit(ctx, map[string]any{
		"files": []any{
			map[string]any{
				"file_path": "a.txt",
				"edits": []any{
					map[string]any{"old_string": "alpha", "new_string": "ALPHA"},
				},
			},
			map[string]any{
				"file_path": "b.txt",
				"edits": []any{
					map[string]any{"old_string": "delta", "new_string": "DELTA"},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("multi_edit files: %v", err)
	}
	for _, want := range []string{
		"The file a.txt has been updated successfully.",
		"The file b.txt has been updated successfully.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("multi-file output = %q, want %q", out, want)
		}
	}
	if got := strings.Count(out, "has been updated successfully"); got != 2 {
		t.Fatalf("multi-file output = %q, want two file results", out)
	}
	if data, _ := os.ReadFile(filepath.Join(root, "a.txt")); string(data) != "ALPHA beta" {
		t.Fatalf("a.txt content = %q", string(data))
	}
	if data, _ := os.ReadFile(filepath.Join(root, "b.txt")); string(data) != "gamma DELTA" {
		t.Fatalf("b.txt content = %q", string(data))
	}
}

func TestRunMultiEditFromRootAtomicFailure(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "f.txt")
	original := "hello world"
	writeFile(t, path, original)

	_, err := RunMultiEditFromRoot(root, path, []Edit{
		{OldString: "hello", NewString: "hi"},
		{OldString: "missing", NewString: "x"},
	})
	if err == nil {
		t.Fatalf("expected error for missing text")
	}
	data, _ := os.ReadFile(path)
	if string(data) != original {
		t.Fatalf("file must be unchanged on failure, got %q", string(data))
	}
}

func TestApplyEditRejectsSame(t *testing.T) {
	if _, err := applyEdit("abc", "a", "a", false); err == nil {
		t.Fatalf("expected error when old==new")
	}
}

package tools

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunMultiEditFromRoot(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "f.txt")
	writeFile(t, path, "alpha beta alpha gamma")

	out, err := RunMultiEditFromRoot(root, path, []Edit{
		{OldString: "beta", NewString: "BETA"},
		{OldString: "alpha", NewString: "X", ReplaceAll: true},
	})
	if err != nil {
		t.Fatalf("multi_edit: %v", err)
	}
	if out == "" {
		t.Fatalf("expected summary output")
	}
	data, _ := os.ReadFile(path)
	if string(data) != "X BETA X gamma" {
		t.Fatalf("unexpected content: %q", string(data))
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

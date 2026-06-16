package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunGrepFromRoot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.go"), "package main\nfunc Foo() {}\n")
	writeFile(t, filepath.Join(root, "b.txt"), "hello FOO world\n")
	writeFile(t, filepath.Join(root, ".git", "config"), "Foo in git\n")

	// files_with_matches default mode, only .go via glob.
	out, err := RunGrepFromRoot(root, root, "Foo", "*.go", "", false, false, 0)
	if err != nil {
		t.Fatalf("grep: %v", err)
	}
	if out != "a.go" {
		t.Fatalf("expected a.go, got %q", out)
	}

	// case-insensitive content with line numbers across all files; .git skipped.
	out, err = RunGrepFromRoot(root, root, "foo", "", "content", true, true, 0)
	if err != nil {
		t.Fatalf("grep content: %v", err)
	}
	if !strings.Contains(out, "a.go:2:") || !strings.Contains(out, "b.txt:1:") {
		t.Fatalf("unexpected content output: %q", out)
	}
	if strings.Contains(out, "git") {
		t.Fatalf("expected .git to be skipped, got %q", out)
	}

	// count mode.
	out, err = RunGrepFromRoot(root, root, "Foo", "*.go", "count", false, false, 0)
	if err != nil {
		t.Fatalf("grep count: %v", err)
	}
	if out != "a.go:1" {
		t.Fatalf("expected a.go:1, got %q", out)
	}

	// no matches.
	out, _ = RunGrepFromRoot(root, root, "zzz", "", "", false, false, 0)
	if out != "No matches found" {
		t.Fatalf("expected no matches, got %q", out)
	}
}

func TestRunGlobFromRoot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "src", "main.go"), "x")
	time.Sleep(10 * time.Millisecond)
	writeFile(t, filepath.Join(root, "src", "util", "helper.go"), "y")
	writeFile(t, filepath.Join(root, "README.md"), "z")

	out, err := RunGlobFromRoot(root, "**/*.go", 0)
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	lines := strings.Split(out, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 go files, got %q", out)
	}
	// newest first: helper.go was written last.
	if !strings.HasSuffix(lines[0], "helper.go") {
		t.Fatalf("expected helper.go first, got %q", out)
	}

	out, err = RunGlobFromRoot(root, "*.md", 0)
	if err != nil {
		t.Fatalf("glob md: %v", err)
	}
	if !strings.HasSuffix(out, "README.md") {
		t.Fatalf("expected README.md, got %q", out)
	}
}

func TestRunLsFromRoot(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "a.go"), "x")
	writeFile(t, filepath.Join(root, "b.log"), "y")
	if err := os.MkdirAll(filepath.Join(root, "sub"), 0o755); err != nil {
		t.Fatal(err)
	}

	out, err := RunLsFromRoot(root, []string{"*.log"})
	if err != nil {
		t.Fatalf("ls: %v", err)
	}
	if !strings.Contains(out, "a.go") || !strings.Contains(out, "sub/") {
		t.Fatalf("unexpected ls output: %q", out)
	}
	if strings.Contains(out, "b.log") {
		t.Fatalf("expected b.log to be ignored, got %q", out)
	}
}

func TestMatchGlob(t *testing.T) {
	cases := []struct {
		pattern, name string
		want          bool
	}{
		{"*.go", "main.go", true},
		{"*.go", "a/b/main.go", true},
		{"**/*.go", "a/b/main.go", true},
		{"**/*.go", "main.go", true},
		{"src/**/*.ts", "src/a/b.ts", true},
		{"src/*.ts", "src/a/b.ts", false},
		{"*.go", "main.txt", false},
	}
	for _, c := range cases {
		if got := matchGlob(c.pattern, c.name); got != c.want {
			t.Errorf("matchGlob(%q,%q)=%v want %v", c.pattern, c.name, got, c.want)
		}
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

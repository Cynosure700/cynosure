package gitcontext

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormatIncludesAllFieldsInOrder(t *testing.T) {
	out := Format(Status{
		IsRepo:        true,
		Branch:        "main",
		MainBranch:    "main",
		WorkTree:      "M src/foo.ts\n?? docs/bar.md",
		RecentCommits: "b78dd22 init",
		UserName:      "dunxing.7",
	})

	for _, want := range []string{
		"This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.",
		"Current branch: main",
		"Main branch (you will usually use this for PRs): main",
		"Git user: dunxing.7",
		"Status:\nM src/foo.ts\n?? docs/bar.md",
		"Recent commits:\nb78dd22 init",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected formatted output to contain %q, got %q", want, out)
		}
	}

	branchIdx := strings.Index(out, "Current branch:")
	statusIdx := strings.Index(out, "Status:")
	commitsIdx := strings.Index(out, "Recent commits:")
	if !(branchIdx < statusIdx && statusIdx < commitsIdx) {
		t.Fatalf("expected branch < status < commits ordering, got %q", out)
	}
}

func TestFormatCleanWorkTree(t *testing.T) {
	out := Format(Status{IsRepo: true, Branch: "main", WorkTree: "   "})
	if !strings.Contains(out, "Status:\n(clean)") {
		t.Fatalf("expected clean work tree rendered as (clean), got %q", out)
	}
}

func TestFormatOmitsEmptyOptionalFields(t *testing.T) {
	out := Format(Status{IsRepo: true, WorkTree: "M a.go"})
	for _, omitted := range []string{"Current branch:", "Main branch", "Git user:", "Recent commits:"} {
		if strings.Contains(out, omitted) {
			t.Fatalf("expected %q to be omitted when field empty, got %q", omitted, out)
		}
	}
	if !strings.Contains(out, "Status:\nM a.go") {
		t.Fatalf("expected status section to remain, got %q", out)
	}
}

func TestFormatNonRepoReturnsEmpty(t *testing.T) {
	if out := Format(Status{IsRepo: false, Branch: "main"}); out != "" {
		t.Fatalf("expected empty output for non-repo, got %q", out)
	}
}

func TestFormatTruncatesLongStatus(t *testing.T) {
	long := strings.Repeat("x", MaxStatusChars+50)
	out := Format(Status{IsRepo: true, WorkTree: long})
	if !strings.Contains(out, statusTruncationNotice) {
		t.Fatalf("expected truncation notice for long status, got len=%d", len(out))
	}
	wantStatus := "Status:\n" + strings.Repeat("x", MaxStatusChars) + "\n" + statusTruncationNotice
	if !strings.Contains(out, wantStatus) {
		t.Fatalf("expected status truncated to %d chars followed by notice, got %q", MaxStatusChars, out)
	}
}

func TestFormatExactBoundaryNotTruncated(t *testing.T) {
	exact := strings.Repeat("y", MaxStatusChars)
	out := Format(Status{IsRepo: true, WorkTree: exact})
	if strings.Contains(out, statusTruncationNotice) {
		t.Fatalf("expected no truncation at exact boundary, got truncation")
	}
}

func TestCollectRealRepository(t *testing.T) {
	requireGit(t)
	dir := tempRepoDir(t)
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.name", "Test User")
	runGit(t, dir, "config", "user.email", "test@example.com")
	if err := writeFile(dir, "a.txt", "hello"); err != nil {
		t.Fatalf("write file: %v", err)
	}
	runGit(t, dir, "add", "a.txt")
	runGit(t, dir, "commit", "-m", "init commit")

	status := Collect(dir)
	if !status.IsRepo {
		t.Fatalf("expected IsRepo=true for git repo")
	}
	if status.Branch == "" {
		t.Fatalf("expected non-empty branch")
	}
	if status.UserName != "Test User" {
		t.Fatalf("expected user name 'Test User', got %q", status.UserName)
	}
	if !strings.Contains(status.RecentCommits, "init commit") {
		t.Fatalf("expected recent commits to contain init commit, got %q", status.RecentCommits)
	}

	if err := writeFile(dir, "b.txt", "world"); err != nil {
		t.Fatalf("write file: %v", err)
	}
	status = Collect(dir)
	if strings.TrimSpace(status.WorkTree) == "" {
		t.Fatalf("expected non-empty work tree after adding untracked file")
	}
}

func TestCollectNonRepository(t *testing.T) {
	requireGit(t)
	dir := t.TempDir()
	if status := Collect(dir); status.IsRepo {
		t.Fatalf("expected IsRepo=false for non-git directory")
	}
}

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(cmd.Environ(),
		"GIT_CONFIG_GLOBAL=/dev/null",
		"GIT_CONFIG_SYSTEM=/dev/null",
		"GIT_AUTHOR_NAME=Test User",
		"GIT_AUTHOR_EMAIL=test@example.com",
		"GIT_COMMITTER_NAME=Test User",
		"GIT_COMMITTER_EMAIL=test@example.com",
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
}

func writeFile(dir, name, content string) error {
	return os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
}

// tempRepoDir 返回一个测试结束后尽力清理的临时目录。Git 在 macOS 上可能留下
// fsmonitor 等后台文件，使 t.TempDir 的严格清理报 "directory not empty"，故自管清理。
func tempRepoDir(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "gitcontext")
	if err != nil {
		t.Fatalf("mkdir temp: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return dir
}

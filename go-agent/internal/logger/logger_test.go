package logger

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitFileLoggerUnderWorkspaceRoot_WritesIntoLogsDirectory(t *testing.T) {
	workspaceRoot := t.TempDir()

	if err := InitFileLoggerUnderWorkspaceRoot(workspaceRoot); err != nil {
		t.Fatalf("init file logger under workspace root: %v", err)
	}

	path := LogFilePath()
	if !strings.HasPrefix(path, filepath.Join(workspaceRoot, "logs")+string(filepath.Separator)) {
		t.Fatalf("expected log file under workspace root logs directory, got %q", path)
	}
	if filepath.Base(filepath.Dir(path)) != "logs" {
		t.Fatalf("expected parent directory to be logs, got %q", filepath.Dir(path))
	}
}

func TestInfoWritesDebugLogIntoWorkspaceLogs(t *testing.T) {
	workspaceRoot := t.TempDir()

	if err := InitFileLoggerUnderWorkspaceRoot(workspaceRoot); err != nil {
		t.Fatalf("init file logger under workspace root: %v", err)
	}

	Info("debug message")

	content, err := os.ReadFile(LogFilePath())
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(content), "debug message") {
		t.Fatalf("expected debug log in file, got %q", string(content))
	}
}

func TestConsoleDisabledSuppressesTerminalOutputButKeepsFileLog(t *testing.T) {
	workspaceRoot := t.TempDir()
	if err := InitFileLoggerUnderWorkspaceRoot(workspaceRoot); err != nil {
		t.Fatalf("init file logger under workspace root: %v", err)
	}
	SetConsoleEnabled(false)
	t.Cleanup(func() { SetConsoleEnabled(true) })

	output := captureStdout(t, func() {
		Info("tui-only log")
	})

	if output != "" {
		t.Fatalf("stdout = %q, want no terminal output while console logging is disabled", output)
	}
	content, err := os.ReadFile(LogFilePath())
	if err != nil {
		t.Fatalf("read log file: %v", err)
	}
	if !strings.Contains(string(content), "tui-only log") {
		t.Fatalf("expected debug log in file, got %q", string(content))
	}
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("create stdout pipe: %v", err)
	}
	os.Stdout = w
	defer func() { os.Stdout = old }()

	fn()
	if err := w.Close(); err != nil {
		t.Fatalf("close stdout pipe: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout pipe: %v", err)
	}
	return string(out)
}

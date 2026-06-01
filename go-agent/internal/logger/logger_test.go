package logger

import (
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

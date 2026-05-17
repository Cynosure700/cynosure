package logger

import (
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

package logger

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestInitFileLoggerUnderAppHome_WritesIntoLogsDirectory(t *testing.T) {
	appHome := t.TempDir()

	if err := InitFileLoggerUnderAppHome(appHome); err != nil {
		t.Fatalf("init file logger under app home: %v", err)
	}

	path := LogFilePath()
	if !strings.HasPrefix(path, filepath.Join(appHome, "logs")+string(filepath.Separator)) {
		t.Fatalf("expected log file under app home logs directory, got %q", path)
	}
	if filepath.Base(filepath.Dir(path)) != "logs" {
		t.Fatalf("expected parent directory to be logs, got %q", filepath.Dir(path))
	}
}

package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"nano_cc/internal/safety"
)

func RunWrite(path, content string) (string, error) {
	resolved, err := safety.SafePath(path)
	if err != nil {
		return "", err
	}

	dir := filepath.Dir(resolved)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create parent directory: %w", err)
	}

	if err := os.WriteFile(resolved, []byte(content), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Wrote %d bytes to %s", len(content), path), nil
}

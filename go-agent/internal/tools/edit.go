package tools

import (
	"fmt"
	"os"
	"strings"

	"nano_cc/internal/safety"
)

func RunEdit(path, oldText, newText string) (string, error) {
	resolved, err := safety.SafePath(path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	content := string(data)
	if !strings.Contains(content, oldText) {
		return "", fmt.Errorf("text not found in %s", path)
	}

	updated := strings.Replace(content, oldText, newText, 1)
	if err := os.WriteFile(resolved, []byte(updated), 0o644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return fmt.Sprintf("Edited %s", path), nil
}

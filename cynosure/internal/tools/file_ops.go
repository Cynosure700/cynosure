package tools

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nano_cc/internal/safety"
)

const maxReadLen = 50000

func RunRead(path string, limit int) (string, error) {
	return RunReadFromRoot("", path, limit)
}

func RunReadFromRoot(root, path string, limit int) (string, error) {
	resolved, err := safety.SafePathFromRoot(root, path)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(resolved)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return "(empty file)", nil
	}

	lines := strings.Split(string(data), "\n")
	if limit > 0 && limit < len(lines) {
		lines = append(lines[:limit], fmt.Sprintf("... (%d more lines)", len(lines)-limit))
	}

	result := strings.Join(lines, "\n")
	if len(result) > maxReadLen {
		result = result[:maxReadLen]
	}

	return result, nil
}

func RunWrite(path, content string) (string, error) {
	return RunWriteFromRoot("", path, content)
}

func RunWriteFromRoot(root, path, content string) (string, error) {
	resolved, err := safety.SafePathFromRoot(root, path)
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

func RunEdit(path, oldText, newText string) (string, error) {
	return RunEditFromRoot("", path, oldText, newText)
}

func RunEditFromRoot(root, path, oldText, newText string) (string, error) {
	resolved, err := safety.SafePathFromRoot(root, path)
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

package tools

import (
	"fmt"
	"os"
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

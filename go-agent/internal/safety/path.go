package safety

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SafePath(p string) (string, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}

	resolved, err := filepath.Abs(filepath.Join(workdir, p))
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	if !strings.HasPrefix(resolved, workdir) {
		return "", fmt.Errorf("path escapes workspace: %s", p)
	}

	return resolved, nil
}

package safety

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func SafePath(p string) (string, error) {
	return SafePathFromRoot("", p)
}

func SafePathFromRoot(root, p string) (string, error) {
	workdir, err := os.Getwd()
	if strings.TrimSpace(root) != "" {
		workdir = root
	} else {
		if err != nil {
			return "", fmt.Errorf("failed to get working directory: %w", err)
		}
	}

	workdir, err = filepath.Abs(workdir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve workspace root: %w", err)
	}

	resolvedPath := p
	if !filepath.IsAbs(resolvedPath) {
		resolvedPath = filepath.Join(workdir, resolvedPath)
	}
	resolved, err := filepath.Abs(resolvedPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve path: %w", err)
	}

	workdir = filepath.Clean(workdir)
	resolved = filepath.Clean(resolved)
	if resolved != workdir && !strings.HasPrefix(resolved, workdir+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes workspace: %s", p)
	}

	return resolved, nil
}

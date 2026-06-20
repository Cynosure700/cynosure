package safety

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SafePathFromRoot 将 p 解析为基于 root 的绝对路径并清理。
// root 为空时回退到当前工作目录。
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

	return filepath.Clean(resolved), nil
}

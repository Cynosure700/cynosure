package runtime

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func (s *Service) resolveUserWorkspace(userID string) (string, error) {
	_ = userID
	base := strings.TrimSpace(s.Cfg.WorkspaceRoot)
	if base == "" {
		return "", fmt.Errorf("workspace root is not configured")
	}
	if err := os.MkdirAll(base, 0o755); err != nil {
		return "", fmt.Errorf("create workspace root: %w", err)
	}
	resolvedBase, err := filepath.Abs(base)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	return filepath.Clean(resolvedBase), nil
}

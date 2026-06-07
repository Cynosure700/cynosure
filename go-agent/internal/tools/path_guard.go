package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"nano_cc/internal/safety"
)

func validatedWorkspaceRootFromContext(ctx context.Context) (string, error) {
	root := strings.TrimSpace(workspaceRootFromContext(ctx))
	if root == "" {
		return "", fmt.Errorf("workspace root is required")
	}
	resolved, err := filepath.Abs(root)
	if err != nil {
		return "", fmt.Errorf("resolve workspace root: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("workspace root is unavailable: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("workspace root is not a directory")
	}
	return filepath.Clean(resolved), nil
}

func validatedCurrentWorkingDirFromContext(ctx context.Context) (string, error) {
	workspaceRoot, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", err
	}
	workingDir := strings.TrimSpace(currentWorkingDirFromContext(ctx))
	if workingDir == "" {
		return workspaceRoot, nil
	}
	resolved, err := filepath.Abs(workingDir)
	if err != nil {
		return "", fmt.Errorf("resolve current working directory: %w", err)
	}
	info, err := os.Stat(resolved)
	if err != nil {
		return "", fmt.Errorf("current working directory is unavailable: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("current working directory is not a directory")
	}
	resolved = filepath.Clean(resolved)
	if !safety.Contains(workspaceRoot, resolved) {
		return "", fmt.Errorf("current working directory escapes workspace: %s", workingDir)
	}
	return resolved, nil
}

func resolvePathFromContext(ctx context.Context, path string) (string, string, error) {
	workspaceRoot, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", "", err
	}
	workingDir, err := validatedCurrentWorkingDirFromContext(ctx)
	if err != nil {
		return "", "", err
	}
	resolvedPath := path
	if !filepath.IsAbs(resolvedPath) {
		resolvedPath = filepath.Join(workingDir, resolvedPath)
	}
	resolvedPath, err = safety.SafePathFromRoot(workspaceRoot, resolvedPath)
	if err != nil {
		return "", "", err
	}
	return workspaceRoot, resolvedPath, nil
}

func validateBashCommandPaths(root, command string, allowOutsideWorkspace bool) error {
	cleanRoot := filepath.Clean(root)
	for _, token := range splitShellFields(command) {
		candidate := cleanShellPathToken(token)
		if candidate == "" || !isShellPathArgument(candidate) {
			continue
		}
		if !filepath.IsAbs(candidate) {
			candidate = filepath.Join(cleanRoot, candidate)
		}
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			return fmt.Errorf("resolve command path: %w", err)
		}
		cleanResolved := filepath.Clean(resolved)
		if allowOutsideWorkspace {
			continue
		}
		if !safety.Contains(cleanRoot, cleanResolved) {
			return fmt.Errorf("command path escapes workspace: %s", token)
		}
	}
	return nil
}

func splitShellFields(command string) []string {
	fields := make([]string, 0)
	var current strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	for _, r := range command {
		if escaped {
			current.WriteRune(r)
			escaped = false
			continue
		}
		if r == '\\' && !inSingle {
			escaped = true
			continue
		}
		switch r {
		case '\'':
			if !inDouble {
				inSingle = !inSingle
				continue
			}
		case '"':
			if !inSingle {
				inDouble = !inDouble
				continue
			}
		case ' ', '\t', '\n', '\r':
			if !inSingle && !inDouble {
				appendShellField(&fields, &current)
				continue
			}
		case ';', '&', '|':
			if !inSingle && !inDouble {
				appendShellField(&fields, &current)
				fields = append(fields, string(r))
				continue
			}
		}
		current.WriteRune(r)
	}
	appendShellField(&fields, &current)
	return fields
}

func appendShellField(fields *[]string, current *strings.Builder) {
	if current.Len() == 0 {
		return
	}
	*fields = append(*fields, current.String())
	current.Reset()
}

func cleanShellPathToken(token string) string {
	trimmed := strings.Trim(token, "\"'`;,()[]{}")
	if strings.Contains(trimmed, "://") {
		return ""
	}
	if idx := strings.IndexAny(trimmed, "<>|"); idx >= 0 {
		trimmed = trimmed[:idx]
	}
	if strings.Contains(trimmed, "=") && !strings.HasPrefix(trimmed, "/") {
		return ""
	}
	return strings.TrimSpace(trimmed)
}

func isShellPathArgument(token string) bool {
	if token == "" || strings.HasPrefix(token, "-") || token == "." || token == ".." {
		return false
	}
	if filepath.IsAbs(token) || strings.HasPrefix(token, "./") || strings.HasPrefix(token, "../") || strings.ContainsRune(token, os.PathSeparator) {
		return true
	}
	return strings.ContainsRune(token, '.')
}

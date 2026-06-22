package tools

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cynosure/internal/safety"
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
	return filepath.Clean(resolved), nil
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

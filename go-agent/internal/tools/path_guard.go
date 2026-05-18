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
	if resolved != workspaceRoot && !strings.HasPrefix(resolved, workspaceRoot+string(os.PathSeparator)) {
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

func validateBashCommandPaths(root, command string) error {
	for _, token := range strings.Fields(command) {
		candidate := strings.Trim(token, "\"'`;,()[]{}")
		if candidate == "" || !filepath.IsAbs(candidate) {
			continue
		}
		resolved, err := filepath.Abs(candidate)
		if err != nil {
			return fmt.Errorf("resolve command path: %w", err)
		}
		cleanRoot := filepath.Clean(root)
		cleanResolved := filepath.Clean(resolved)
		if cleanResolved != cleanRoot && !strings.HasPrefix(cleanResolved, cleanRoot+string(os.PathSeparator)) {
			return fmt.Errorf("command path escapes workspace: %s", candidate)
		}
	}
	return nil
}

package tools

import (
	"context"
	"fmt"
)

var Handlers = map[string]ToolHandler{
	"bash":       handleBash,
	"read_file":  handleRead,
	"write_file": handleWrite,
	"edit_file":  handleEdit,
	"load_skill": handleLoadSkill,
}

func handleBash(ctx context.Context, args map[string]any) (string, error) {
	cmd, _ := args["command"].(string)
	if cmd == "" {
		return "", fmt.Errorf("command is required")
	}
	root, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", err
	}
	workingDir, err := validatedCurrentWorkingDirFromContext(ctx)
	if err != nil {
		return "", err
	}
	if err := validateBashCommandPaths(root, cmd, allowOutsideWorkspaceFromContext(ctx)); err != nil {
		return "", err
	}
	return RunBashInDirWithOptions(cmd, workingDir, allowDangerousCommandsFromContext(ctx))
}

func handleRead(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	limit := 0
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunReadFromRoot(root, resolvedPath, limit)
}

func handleWrite(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	content, _ := args["content"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunWriteFromRoot(root, resolvedPath, content)
}

func handleEdit(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	oldText, _ := args["old_text"].(string)
	newText, _ := args["new_text"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunEditFromRoot(root, resolvedPath, oldText, newText)
}

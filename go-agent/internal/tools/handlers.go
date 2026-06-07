package tools

import (
	"context"
	"fmt"
)

// TodoWriteToolName is the tool that updates the task plan and returns
// structured todo items in addition to a textual summary.
const TodoWriteToolName = "todo_write"

// ExecResult is the unified result of executing a stateless tool. Todos is
// populated only by the todo_write tool.
type ExecResult struct {
	Output string
	Todos  []TodoItem
}

// Handlers maps stateless tool names to their textual handlers. todo_write is
// dispatched separately because it returns structured todos (see Dispatch).
var Handlers = map[string]ToolHandler{
	"bash":                  handleBash,
	"read_file":             handleRead,
	"write_file":            handleWrite,
	"edit_file":             handleEdit,
	"load_skill":            handleLoadSkill,
	"read_persisted_output": handleReadPersistedOutput,
}

// Dispatch is the single entry point for executing a stateless tool by name.
// It is the authority for tool execution semantics, including todo_write's
// structured output.
func Dispatch(ctx context.Context, name string, args map[string]any) (ExecResult, error) {
	if name == TodoWriteToolName {
		result, err := ExecuteTodoWrite(ctx, args)
		if err != nil {
			return ExecResult{}, err
		}
		return ExecResult{Output: result.Output, Todos: result.Todos}, nil
	}
	handler, ok := Handlers[name]
	if !ok || handler == nil {
		return ExecResult{}, fmt.Errorf("tool %s has no handler", name)
	}
	output, err := handler(ctx, args)
	if err != nil {
		return ExecResult{}, err
	}
	return ExecResult{Output: output}, nil
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
	if err := validateBashCommandPaths(root, cmd, allowOutsideWorkspaceFromContext(ctx), systemAssetDirsFromContext(ctx)...); err != nil {
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

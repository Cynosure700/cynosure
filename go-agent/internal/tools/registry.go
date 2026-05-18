package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/safety"
)

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

type RuntimeEnv struct {
	AppHome          string
	CommandBinDir    string
	CommandScriptDir string
	WorkspaceRoot    string
	CurrentWorkingDir string
}

type contextKey string

const runtimeEnvContextKey contextKey = "runtime_env"

var Handlers = map[string]ToolHandler{
	"bash":          handleBash,
	"read_file":     handleRead,
	"write_file":    handleWrite,
	"edit_file":     handleEdit,
	"load_skill":    nil, // handled by web runtime registry
}

func WithRuntimeEnv(ctx context.Context, env RuntimeEnv) context.Context {
	return context.WithValue(ctx, runtimeEnvContextKey, env)
}

func RuntimeEnvFromContext(ctx context.Context) (RuntimeEnv, bool) {
	env, ok := ctx.Value(runtimeEnvContextKey).(RuntimeEnv)
	return env, ok
}

func workspaceRootFromContext(ctx context.Context) string {
	env, ok := RuntimeEnvFromContext(ctx)
	if !ok {
		return ""
	}
	return env.WorkspaceRoot
}

func currentWorkingDirFromContext(ctx context.Context) string {
	env, ok := RuntimeEnvFromContext(ctx)
	if !ok {
		return ""
	}
	if strings.TrimSpace(env.CurrentWorkingDir) != "" {
		return env.CurrentWorkingDir
	}
	return env.WorkspaceRoot
}

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
	if err := validateBashCommandPaths(root, cmd); err != nil {
		return "", err
	}
	return RunBashInDir(cmd, workingDir)
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

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func toolDef(name, desc string, params any) openai.Tool {
	return openai.Tool{
		Type: "function",
		Function: &openai.FunctionDefinition{
			Name:        name,
			Description: desc,
			Parameters:  mustMarshal(params),
		},
	}
}

func strParam(desc string, required bool) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": desc,
	}
}

func intParam(desc string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"description": desc,
	}
}

var baseToolDefs = []openai.Tool{
	toolDef("bash", "Execute a shell command via bash -c", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": strParam("The shell command to execute", true),
		},
		"required": []string{"command"},
	}),
	toolDef("read_file", "Read a file from the filesystem", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":  strParam("Path to the file to read", true),
			"limit": intParam("Maximum number of lines to read"),
		},
		"required": []string{"path"},
	}),
	toolDef("write_file", "Write content to a file", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    strParam("Path to the file to write", true),
			"content": strParam("Content to write to the file", true),
		},
		"required": []string{"path", "content"},
	}),
	toolDef("edit_file", "Replace text in a file by exact match", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     strParam("Path to the file to edit", true),
			"old_text": strParam("Exact text to find and replace", true),
			"new_text": strParam("Text to replace it with", true),
		},
		"required": []string{"path", "old_text", "new_text"},
	}),
	toolDef("load_skill", "Load a skill by name to access specialized knowledge", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": strParam("Name of the skill to load", true),
		},
		"required": []string{"name"},
	}),
}

var ChildToolDefs = baseToolDefs

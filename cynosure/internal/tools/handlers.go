package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// TodoWriteToolName is the tool that updates the task plan and returns
// structured todo items in addition to a textual summary.
const TodoWriteToolName = "todo_write"

// TodoListToolName is the read-only tool that returns the current task plan.
const TodoListToolName = "todo_list"

// normalToolTimeout bounds the execution of non-terminal tools. The terminal
// tool (bash) enforces its own timeout in bash.go.
const normalToolTimeout = 60 * time.Second

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
	"multi_edit":            handleMultiEdit,
	"grep":                  handleGrep,
	"glob":                  handleGlob,
	"ls":                    handleLs,
	"web_fetch":             handleWebFetch,
	"web_search":            handleWebSearch,
	"load_skill":            handleLoadSkill,
	"todo_list":             handleTodoList,
	"read_persisted_output": handleReadPersistedOutput,
}

// Dispatch is the single entry point for executing a stateless tool by name.
// It is the authority for tool execution semantics, including todo_write's
// structured output.
//
// The terminal tool (bash) enforces its own 120s timeout internally; all other
// tools are bounded by normalToolTimeout via a watchdog here.
func Dispatch(ctx context.Context, name string, args map[string]any) (ExecResult, error) {
	if name == "bash" {
		return dispatchOnce(ctx, name, args)
	}
	return runWithTimeout(normalToolTimeout, name, func() (ExecResult, error) {
		return dispatchOnce(ctx, name, args)
	})
}

func dispatchOnce(ctx context.Context, name string, args map[string]any) (ExecResult, error) {
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

// runWithTimeout executes fn in a goroutine and abandons the wait after d,
// returning a timeout error. The result channel is buffered so the goroutine
// never blocks on send even after a timeout.
func runWithTimeout(d time.Duration, name string, fn func() (ExecResult, error)) (ExecResult, error) {
	type result struct {
		out ExecResult
		err error
	}
	ch := make(chan result, 1)
	go func() {
		out, err := fn()
		ch <- result{out: out, err: err}
	}()
	select {
	case r := <-ch:
		return r.out, r.err
	case <-time.After(d):
		return ExecResult{}, fmt.Errorf("tool %s timed out after %v", name, d)
	}
}

func handleBash(ctx context.Context, args map[string]any) (string, error) {
	cmd, _ := args["command"].(string)
	if cmd == "" {
		return "", fmt.Errorf("command is required")
	}
	workingDir, err := validatedCurrentWorkingDirFromContext(ctx)
	if err != nil {
		return "", err
	}
	return RunBashInDir(cmd, workingDir)
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

func handleMultiEdit(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["file_path"].(string)
	if path == "" {
		return "", fmt.Errorf("file_path is required")
	}
	rawEdits, ok := args["edits"].([]any)
	if !ok || len(rawEdits) == 0 {
		return "", fmt.Errorf("edits is required")
	}
	edits := make([]Edit, 0, len(rawEdits))
	for i, raw := range rawEdits {
		m, ok := raw.(map[string]any)
		if !ok {
			return "", fmt.Errorf("edit %d is not an object", i+1)
		}
		oldStr, _ := m["old_string"].(string)
		newStr, _ := m["new_string"].(string)
		replaceAll, _ := m["replace_all"].(bool)
		edits = append(edits, Edit{OldString: oldStr, NewString: newStr, ReplaceAll: replaceAll})
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunMultiEditFromRoot(root, resolvedPath, edits)
}

func handleGrep(ctx context.Context, args map[string]any) (string, error) {
	pattern, _ := args["pattern"].(string)
	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}
	root, resolvedPath, err := resolveSearchPathFromContext(ctx, args)
	if err != nil {
		return "", err
	}
	globPattern, _ := args["glob"].(string)
	outputMode, _ := args["output_mode"].(string)
	ignoreCase, _ := args["-i"].(bool)
	showLineNumbers, _ := args["-n"].(bool)
	headLimit := 0
	if l, ok := args["head_limit"].(float64); ok {
		headLimit = int(l)
	}
	return RunGrepFromRoot(root, resolvedPath, pattern, globPattern, outputMode, ignoreCase, showLineNumbers, headLimit)
}

func handleGlob(ctx context.Context, args map[string]any) (string, error) {
	pattern, _ := args["pattern"].(string)
	if pattern == "" {
		return "", fmt.Errorf("pattern is required")
	}
	_, resolvedPath, err := resolveSearchPathFromContext(ctx, args)
	if err != nil {
		return "", err
	}
	headLimit := 0
	if l, ok := args["head_limit"].(float64); ok {
		headLimit = int(l)
	}
	return RunGlobFromRoot(resolvedPath, pattern, headLimit)
}

func handleLs(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["path"].(string)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	if !filepath.IsAbs(path) {
		return "", fmt.Errorf("path must be an absolute path: %s", path)
	}
	_, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	ignore := stringSliceArg(args["ignore"])
	return RunLsFromRoot(resolvedPath, ignore)
}

func handleWebFetch(ctx context.Context, args map[string]any) (string, error) {
	url, _ := args["url"].(string)
	if url == "" {
		return "", fmt.Errorf("url is required")
	}
	prompt, _ := args["prompt"].(string)
	return RunWebFetch(ctx, url, prompt)
}

func handleWebSearch(ctx context.Context, args map[string]any) (string, error) {
	query, _ := args["query"].(string)
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	return RunWebSearch(ctx, query)
}

// resolveSearchPathFromContext resolves the optional path argument for grep and
// glob, defaulting to the current working directory and constraining the result
// to the workspace.
func resolveSearchPathFromContext(ctx context.Context, args map[string]any) (string, string, error) {
	path, _ := args["path"].(string)
	if strings.TrimSpace(path) == "" {
		path = "."
	}
	return resolvePathFromContext(ctx, path)
}

func stringSliceArg(value any) []string {
	raw, ok := value.([]any)
	if !ok {
		return nil
	}
	result := make([]string, 0, len(raw))
	for _, item := range raw {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}
	return result
}

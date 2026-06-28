package tools

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// TodoWriteToolName 是用于更新任务计划的工具，除返回文本摘要外，
// 还会返回结构化的待办项。
const TodoWriteToolName = "todo_write"

// TodoListToolName 是只读工具，返回当前任务计划。
const TodoListToolName = "todo_list"

// normalToolTimeout 限定非终端类工具的执行时长。终端类工具（bash）
// 在 bash.go 中自行实现超时。
const normalToolTimeout = 60 * time.Second

// ExecResult 是执行无状态工具的统一结果。Todos 仅由 todo_write 工具填充。
type ExecResult struct {
	Output string
	Todos  []TodoItem
}

// Handlers 将无状态工具名映射到对应的文本处理函数。todo_write 单独分发，
// 因为它会返回结构化的待办项（见 Dispatch）。
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

// Dispatch 是按名称执行无状态工具的唯一入口。它是工具执行语义的权威所在，
// 包括 todo_write 的结构化输出。
//
// 终端类工具（bash）在内部自行实现 120 秒超时；其余所有工具均由这里的
// 看门狗按 normalToolTimeout 限定执行时长。
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

// runWithTimeout 在一个 goroutine 中执行 fn，并在超过 d 后放弃等待，
// 返回超时错误。结果通道带缓冲，因此即便超时后 goroutine 也不会在发送时阻塞。
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
	workingDir, err := validatedWorkspaceRootFromContext(ctx)
	if err != nil {
		return "", err
	}
	return RunBashInDir(cmd, workingDir)
}

func handleRead(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["file_path"].(string)
	if path == "" {
		return "", fmt.Errorf("file_path is required")
	}
	offset := 1
	if o, ok := args["offset"].(float64); ok {
		offset = int(o)
	}
	limit := 0
	if l, ok := args["limit"].(float64); ok {
		limit = int(l)
	}
	root, resolvedPath, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunReadFromRoot(root, resolvedPath, offset, limit)
}

func handleWrite(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["file_path"].(string)
	content, _ := args["content"].(string)
	if path == "" {
		return "", fmt.Errorf("file_path is required")
	}
	root, _, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunWriteFromRoot(root, path, content)
}

func handleEdit(ctx context.Context, args map[string]any) (string, error) {
	path, _ := args["file_path"].(string)
	oldText, _ := args["old_text"].(string)
	newText, _ := args["new_text"].(string)
	if path == "" {
		return "", fmt.Errorf("file_path is required")
	}
	root, _, err := resolvePathFromContext(ctx, path)
	if err != nil {
		return "", err
	}
	return RunEditFromRoot(root, path, oldText, newText)
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

// resolveSearchPathFromContext 解析 grep 和 glob 的可选 path 参数，
// 默认使用工作区根目录，并将结果限制在工作区内。
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

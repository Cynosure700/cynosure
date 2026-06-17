package tools

import (
	"context"
	"strings"
)

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

type RuntimeEnv struct {
	WorkspaceRoot     string
	CurrentWorkingDir string
}

type contextKey string

const runtimeEnvContextKey contextKey = "runtime_env"
const todoSnapshotContextKey contextKey = "todo_snapshot"

func WithRuntimeEnv(ctx context.Context, env RuntimeEnv) context.Context {
	return context.WithValue(ctx, runtimeEnvContextKey, env)
}

func RuntimeEnvFromContext(ctx context.Context) (RuntimeEnv, bool) {
	env, ok := ctx.Value(runtimeEnvContextKey).(RuntimeEnv)
	return env, ok
}

func WithTodoSnapshot(ctx context.Context, todos []TodoItem) context.Context {
	return context.WithValue(ctx, todoSnapshotContextKey, append([]TodoItem(nil), todos...))
}

func TodoSnapshotFromContext(ctx context.Context) ([]TodoItem, bool) {
	todos, ok := ctx.Value(todoSnapshotContextKey).([]TodoItem)
	if !ok {
		return nil, false
	}
	return append([]TodoItem(nil), todos...), true
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

// WebProcessor processes fetched web content with an LLM given a user prompt.
// It is injected by the runtime layer so the tools package does not depend on
// the llm/runtime packages.
type WebProcessor func(ctx context.Context, prompt, content string) (string, error)

const webProcessorContextKey contextKey = "web_processor"

func WithWebProcessor(ctx context.Context, fn WebProcessor) context.Context {
	return context.WithValue(ctx, webProcessorContextKey, fn)
}

func webProcessorFromContext(ctx context.Context) (WebProcessor, bool) {
	fn, ok := ctx.Value(webProcessorContextKey).(WebProcessor)
	return fn, ok && fn != nil
}

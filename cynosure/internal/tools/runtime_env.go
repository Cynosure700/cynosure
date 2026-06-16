package tools

import (
	"context"
	"strings"
)

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

type RuntimeEnv struct {
	WorkspaceRoot          string
	CurrentWorkingDir      string
	AllowOutsideWorkspace  bool
	AllowDangerousCommands bool
}

type contextKey string

const runtimeEnvContextKey contextKey = "runtime_env"

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

func allowOutsideWorkspaceFromContext(ctx context.Context) bool {
	env, ok := RuntimeEnvFromContext(ctx)
	return ok && env.AllowOutsideWorkspace
}

func allowDangerousCommandsFromContext(ctx context.Context) bool {
	env, ok := RuntimeEnvFromContext(ctx)
	return ok && env.AllowDangerousCommands
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

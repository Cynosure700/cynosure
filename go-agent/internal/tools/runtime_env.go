package tools

import (
	"context"
	"strings"
)

type ToolHandler func(ctx context.Context, args map[string]any) (string, error)

type RuntimeEnv struct {
	AppHome                string
	CommandBinDir          string
	CommandScriptDir       string
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

func systemAssetDirsFromContext(ctx context.Context) []string {
	env, ok := RuntimeEnvFromContext(ctx)
	if !ok {
		return nil
	}
	dirs := make([]string, 0, 3)
	for _, dir := range []string{env.CommandBinDir, env.CommandScriptDir} {
		if strings.TrimSpace(dir) != "" {
			dirs = append(dirs, dir)
		}
	}
	return dirs
}

func allowDangerousCommandsFromContext(ctx context.Context) bool {
	env, ok := RuntimeEnvFromContext(ctx)
	return ok && env.AllowDangerousCommands
}

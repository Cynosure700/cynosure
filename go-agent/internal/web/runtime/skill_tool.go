package runtime

import (
	"strings"

	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
)

func (r *ToolRegistry) loadSkillContent(loader *sessions.SkillLoader, skillName string) (ToolExecutionResult, error) {
	content, err := loader.GetContent(skillName)
	if err != nil {
		return ToolExecutionResult{}, err
	}
	envNote := formatRuntimeEnvNote(r.runtimeEnv())
	if envNote == "" {
		return ToolExecutionResult{Output: content}, nil
	}
	return ToolExecutionResult{Output: content + "\n\n" + envNote}, nil
}

func formatRuntimeEnvNote(env agenttools.RuntimeEnv) string {
	lines := make([]string, 0, 5)
	if env.AppHome != "" {
		lines = append(lines, "APP_HOME="+env.AppHome)
	}
	if env.CommandBinDir != "" {
		lines = append(lines, "COMMAND_BIN_DIR="+env.CommandBinDir)
	}
	if env.CommandScriptDir != "" {
		lines = append(lines, "COMMAND_SCRIPT_DIR="+env.CommandScriptDir)
	}
	if env.WorkspaceRoot != "" {
		lines = append(lines, "WORKSPACE_ROOT="+env.WorkspaceRoot)
	}
	if env.CurrentWorkingDir != "" {
		lines = append(lines, "CURRENT_WORKING_DIR="+env.CurrentWorkingDir)
	}
	if len(lines) == 0 {
		return ""
	}
	return "<runtime-paths>\n" + strings.Join(lines, "\n") + "\n</runtime-paths>"
}

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type runtimePaths struct {
	workspaceRoot    string
	builtinSkillsDir string
	commandBinDir    string
	commandScriptDir string
}

type runtimeAssetSpec struct {
	label     string
	envKey    string
	fileValue string
	subdir    string
}

func resolveAppHome(fileCfg fileConfig) (string, error) {
	appHome := getenv("APP_HOME", firstNonEmpty(fileCfg.AppHome, "."))
	resolved, err := filepath.Abs(appHome)
	if err != nil {
		return "", fmt.Errorf("resolve APP_HOME: %w", err)
	}
	return filepath.Clean(resolved), nil
}

func resolvePath(appHome, pathValue string) (string, error) {
	if strings.TrimSpace(pathValue) == "" {
		pathValue = "."
	}
	if !filepath.IsAbs(pathValue) {
		pathValue = filepath.Join(appHome, pathValue)
	}
	resolved, err := filepath.Abs(pathValue)
	if err != nil {
		return "", err
	}
	return filepath.Clean(resolved), nil
}

func resolveWorkspaceRoot(appHome, configured string) (string, error) {
	configured = strings.TrimSpace(configured)
	if configured != "" {
		return resolvePath(appHome, configured)
	}

	for _, candidate := range []string{filepath.Join("output", "workspace"), "workspace"} {
		resolved, err := resolvePath(appHome, candidate)
		if err != nil {
			return "", err
		}
		if workspaceExists(resolved) {
			return resolved, nil
		}
	}

	return resolvePath(appHome, "workspace")
}

func resolveRuntimePaths(appHome string, fileCfg fileConfig) (runtimePaths, error) {
	workspaceRoot, err := resolveWorkspaceRoot(appHome, getenv("WORKSPACE_ROOT", strings.TrimSpace(fileCfg.WorkspaceRoot)))
	if err != nil {
		return runtimePaths{}, fmt.Errorf("resolve workspace root: %w", err)
	}
	paths := runtimePaths{workspaceRoot: workspaceRoot}
	for _, spec := range []runtimeAssetSpec{
		{label: "builtin skills dir", envKey: "BUILTIN_SKILLS_DIR", fileValue: fileCfg.BuiltinSkillsDir, subdir: "skills"},
		{label: "command bin dir", envKey: "COMMAND_BIN_DIR", fileValue: fileCfg.CommandBinDir, subdir: "bin"},
		{label: "command script dir", envKey: "COMMAND_SCRIPT_DIR", fileValue: fileCfg.CommandScriptDir, subdir: "cmd"},
	} {
		resolved, err := resolveRuntimeAssetFromSpec(appHome, workspaceRoot, spec)
		if err != nil {
			return runtimePaths{}, err
		}
		switch spec.subdir {
		case "skills":
			paths.builtinSkillsDir = resolved
		case "bin":
			paths.commandBinDir = resolved
		case "cmd":
			paths.commandScriptDir = resolved
		}
	}
	return paths, nil
}

func resolveRuntimeAssetFromSpec(appHome, workspaceRoot string, spec runtimeAssetSpec) (string, error) {
	configured := getenv(spec.envKey, strings.TrimSpace(spec.fileValue))
	resolved, err := resolveRuntimeAssetDir(appHome, workspaceRoot, configured, spec.subdir)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", spec.label, err)
	}
	return resolved, nil
}

func resolveRuntimeAssetDir(appHome, workspaceRoot, configured, subdir string) (string, error) {
	expected, err := resolvePath(workspaceRoot, subdir)
	if err != nil {
		return "", err
	}

	configured = strings.TrimSpace(configured)
	if configured == "" {
		return expected, nil
	}

	resolved, err := resolvePath(appHome, configured)
	if err != nil {
		return "", err
	}
	if resolved != expected {
		return "", fmt.Errorf("runtime asset dir must stay under workspace root: expected %q", expected)
	}
	return resolved, nil
}

func workspaceExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadWebConfig_ReadsPathSettingsFromConfigFile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"builtin_skills_dir": "shared-workspace/skills",
		"command_bin_dir": "shared-workspace/bin",
		"command_script_dir": "shared-workspace/cmd",
		"workspace_root": "shared-workspace",
		"web_allowed_tools": "load_skill,bash"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("BUILTIN_SKILLS_DIR", "")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WORKSPACE_ROOT", "")
	t.Setenv("WEB_ALLOWED_TOOLS", "")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.AppHome != filepath.Clean(root) {
		t.Fatalf("expected app home %q, got %q", filepath.Clean(root), cfg.AppHome)
	}
	if cfg.BuiltinSkillsDir != filepath.Join(root, "shared-workspace", "skills") {
		t.Fatalf("unexpected builtin skills dir: %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "shared-workspace", "bin") {
		t.Fatalf("unexpected command bin dir: %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "shared-workspace", "cmd") {
		t.Fatalf("unexpected command script dir: %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "shared-workspace") {
		t.Fatalf("unexpected workspace root: %q", cfg.WorkspaceRoot)
	}
	if cfg.SystemPromptPath != filepath.Join(root, "system_prompt.md") {
		t.Fatalf("unexpected system prompt path: %q", cfg.SystemPromptPath)
	}
	if !reflect.DeepEqual(cfg.WebAllowedTools, []string{"load_skill", "bash"}) {
		t.Fatalf("unexpected web allowed tools: %#v", cfg.WebAllowedTools)
	}
	if cfg.BashAllowOutsideWorkspace || cfg.BashAllowDangerousCommands {
		t.Fatalf("expected bash safety flags to default false")
	}
}

func TestLoadWebConfig_EnvironmentOverridesConfigFilePaths(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"builtin_skills_dir": "config-skills",
		"command_bin_dir": "config-bin",
		"command_script_dir": "config-cmd",
		"workspace_root": "config-workspace",
		"web_allowed_tools": "load_skill,bash"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", root)
	t.Setenv("BUILTIN_SKILLS_DIR", filepath.Join("env-workspace", "skills"))
	t.Setenv("COMMAND_BIN_DIR", filepath.Join("env-workspace", "bin"))
	t.Setenv("COMMAND_SCRIPT_DIR", filepath.Join("env-workspace", "cmd"))
	t.Setenv("WORKSPACE_ROOT", "env-workspace")
	t.Setenv("WEB_ALLOWED_TOOLS", "load_skill,edit_file")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.BuiltinSkillsDir != filepath.Join(root, "env-workspace", "skills") {
		t.Fatalf("expected env builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "env-workspace", "bin") {
		t.Fatalf("expected env command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "env-workspace", "cmd") {
		t.Fatalf("expected env command script dir, got %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "env-workspace") {
		t.Fatalf("expected env workspace root, got %q", cfg.WorkspaceRoot)
	}
	if !reflect.DeepEqual(cfg.WebAllowedTools, []string{"load_skill", "edit_file"}) {
		t.Fatalf("expected env web tools, got %#v", cfg.WebAllowedTools)
	}
}

func TestLoadWebConfig_ReadsSystemPromptPathFromConfigAndEnv(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"system_prompt_path": "config-prompt.md"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("SYSTEM_PROMPT_PATH", "env-prompt.md")
	t.Setenv("BUILTIN_SKILLS_DIR", "")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WORKSPACE_ROOT", "")
	t.Setenv("WEB_ALLOWED_TOOLS", "")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}
	if cfg.SystemPromptPath != filepath.Join(root, "env-prompt.md") {
		t.Fatalf("expected env system prompt path, got %q", cfg.SystemPromptPath)
	}
}

func TestLoadWebConfig_ReadsBashSafetyFlagsFromConfigAndEnv(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"bash_allow_outside_workspace": true,
		"bash_allow_dangerous_commands": false
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("BUILTIN_SKILLS_DIR", "")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WORKSPACE_ROOT", "")
	t.Setenv("WEB_ALLOWED_TOOLS", "")
	t.Setenv("BASH_ALLOW_OUTSIDE_WORKSPACE", "false")
	t.Setenv("BASH_ALLOW_DANGEROUS_COMMANDS", "true")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}
	if cfg.BashAllowOutsideWorkspace {
		t.Fatalf("expected env to override outside-workspace allow flag to false")
	}
	if !cfg.BashAllowDangerousCommands {
		t.Fatalf("expected env to override dangerous-command allow flag to true")
	}
}

func TestLoadWebConfig_UsesAppHomeScopedDefaultPaths(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() {
		_ = os.Chdir(oldWD)
	}()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("BUILTIN_SKILLS_DIR", "")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WORKSPACE_ROOT", "")
	t.Setenv("WEB_ALLOWED_TOOLS", "")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.BuiltinSkillsDir != filepath.Join(root, "workspace", "skills") {
		t.Fatalf("expected default builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "workspace", "bin") {
		t.Fatalf("expected default command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "workspace", "cmd") {
		t.Fatalf("expected default command script dir, got %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "workspace") {
		t.Fatalf("expected default workspace root, got %q", cfg.WorkspaceRoot)
	}
	expectedTools := []string{"load_skill", "bash", "read_file", "write_file", "edit_file"}
	if !reflect.DeepEqual(cfg.WebAllowedTools, expectedTools) {
		t.Fatalf("expected default web tools %v, got %#v", expectedTools, cfg.WebAllowedTools)
	}
}

func TestLoadWebConfig_DefaultsToAppHomeWorkspaceEvenWhenOutputWorkspaceExists(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	for _, path := range []string{
		filepath.Join(root, "workspace"),
		filepath.Join(root, "output", "workspace"),
	} {
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatalf("mkdir workspace path %q: %v", path, err)
		}
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("BUILTIN_SKILLS_DIR", "")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WORKSPACE_ROOT", "")
	t.Setenv("WEB_ALLOWED_TOOLS", "")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	expectedWorkspace := filepath.Join(root, "workspace")
	if cfg.WorkspaceRoot != expectedWorkspace {
		t.Fatalf("expected app home workspace root %q, got %q", expectedWorkspace, cfg.WorkspaceRoot)
	}
	if cfg.BuiltinSkillsDir != filepath.Join(expectedWorkspace, "skills") {
		t.Fatalf("expected default builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(expectedWorkspace, "bin") {
		t.Fatalf("expected default command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(expectedWorkspace, "cmd") {
		t.Fatalf("expected default command script dir, got %q", cfg.CommandScriptDir)
	}
}

func TestLoadWebConfig_ExplicitWorkspaceOverrideDrivesDefaultAssetDirs(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("BUILTIN_SKILLS_DIR", "")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WORKSPACE_ROOT", filepath.Join("custom", "runtime-workspace"))
	t.Setenv("WEB_ALLOWED_TOOLS", "")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	expectedWorkspace := filepath.Join(root, "custom", "runtime-workspace")
	if cfg.WorkspaceRoot != expectedWorkspace {
		t.Fatalf("expected explicit workspace root %q, got %q", expectedWorkspace, cfg.WorkspaceRoot)
	}
	if cfg.BuiltinSkillsDir != filepath.Join(expectedWorkspace, "skills") {
		t.Fatalf("expected derived builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(expectedWorkspace, "bin") {
		t.Fatalf("expected derived command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(expectedWorkspace, "cmd") {
		t.Fatalf("expected derived command script dir, got %q", cfg.CommandScriptDir)
	}
}

func TestResolveRuntimePaths_ReturnsCanonicalWorkspaceDerivedDirs(t *testing.T) {
	root := t.TempDir()

	paths, err := resolveRuntimePaths(root, fileConfig{
		WorkspaceRoot:    "custom/runtime-workspace",
		BuiltinSkillsDir: filepath.Join("custom", "runtime-workspace", "skills"),
		CommandBinDir:    filepath.Join("custom", "runtime-workspace", "bin"),
		CommandScriptDir: filepath.Join("custom", "runtime-workspace", "cmd"),
	})
	if err != nil {
		t.Fatalf("resolve runtime paths: %v", err)
	}

	expectedWorkspace := filepath.Join(root, "custom", "runtime-workspace")
	if paths.workspaceRoot != expectedWorkspace {
		t.Fatalf("expected workspace root %q, got %q", expectedWorkspace, paths.workspaceRoot)
	}
	if paths.builtinSkillsDir != filepath.Join(expectedWorkspace, "skills") {
		t.Fatalf("expected builtin skills dir under workspace root, got %q", paths.builtinSkillsDir)
	}
	if paths.commandBinDir != filepath.Join(expectedWorkspace, "bin") {
		t.Fatalf("expected command bin dir under workspace root, got %q", paths.commandBinDir)
	}
	if paths.commandScriptDir != filepath.Join(expectedWorkspace, "cmd") {
		t.Fatalf("expected command script dir under workspace root, got %q", paths.commandScriptDir)
	}
}

func TestLoadWebConfig_RejectsRuntimeAssetDirsOutsideWorkspaceRoot(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"api_key": "test-key",
		"model_id": "test-model",
		"app_home": "` + root + `"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	t.Setenv("OPENAI_BASE_URL", "")
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("MODEL_ID", "")
	t.Setenv("APP_HOME", "")
	t.Setenv("WORKSPACE_ROOT", filepath.Join("custom", "runtime-workspace"))
	t.Setenv("BUILTIN_SKILLS_DIR", "other-skills")
	t.Setenv("COMMAND_BIN_DIR", "")
	t.Setenv("COMMAND_SCRIPT_DIR", "")
	t.Setenv("WEB_ALLOWED_TOOLS", "")

	_, err = LoadWebConfig()
	if err == nil {
		t.Fatalf("expected runtime asset override outside workspace root to be rejected")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "runtime asset dir must stay under workspace root") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureAppLayout_CreatesExpectedDirectories(t *testing.T) {
	root := t.TempDir()
	cfg := AppConfig{
		AppHome:          root,
		BuiltinSkillsDir: filepath.Join(root, "workspace", "skills"),
		CommandBinDir:    filepath.Join(root, "workspace", "bin"),
		CommandScriptDir: filepath.Join(root, "workspace", "cmd"),
		WorkspaceRoot:    filepath.Join(root, "workspace"),
	}

	if err := EnsureAppLayout(cfg); err != nil {
		t.Fatalf("ensure app layout: %v", err)
	}

	paths := []string{
		root,
		filepath.Join(root, "workspace", "skills"),
		filepath.Join(root, "workspace", "bin"),
		filepath.Join(root, "workspace", "cmd"),
		filepath.Join(root, "workspace"),
		filepath.Join(root, "workspace", "logs"),
	}
	for _, path := range paths {
		info, err := os.Stat(path)
		if err != nil {
			t.Fatalf("expected %q to exist: %v", path, err)
		}
		if !info.IsDir() {
			t.Fatalf("expected %q to be a directory", path)
		}
	}
}

func TestValidateAppLayout_AcceptsPreparedDirectories(t *testing.T) {
	root := t.TempDir()
	cfg := AppConfig{
		AppHome:          root,
		BuiltinSkillsDir: filepath.Join(root, "workspace", "skills"),
		CommandBinDir:    filepath.Join(root, "workspace", "bin"),
		CommandScriptDir: filepath.Join(root, "workspace", "cmd"),
		WorkspaceRoot:    filepath.Join(root, "workspace"),
	}

	if err := EnsureAppLayout(cfg); err != nil {
		t.Fatalf("ensure app layout: %v", err)
	}
	if err := ValidateAppLayout(cfg); err != nil {
		t.Fatalf("validate app layout: %v", err)
	}
}

func TestValidateAppLayout_RejectsFilePath(t *testing.T) {
	root := t.TempDir()
	filePath := filepath.Join(root, "not-a-dir")
	if err := os.WriteFile(filePath, []byte("file"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	cfg := AppConfig{
		AppHome:          root,
		BuiltinSkillsDir: filepath.Join(root, "workspace", "skills"),
		CommandBinDir:    filePath,
		CommandScriptDir: filepath.Join(root, "workspace", "cmd"),
		WorkspaceRoot:    filepath.Join(root, "workspace"),
	}
	if err := os.MkdirAll(cfg.BuiltinSkillsDir, 0o755); err != nil {
		t.Fatalf("mkdir builtin skills dir: %v", err)
	}
	if err := os.MkdirAll(cfg.CommandScriptDir, 0o755); err != nil {
		t.Fatalf("mkdir command script dir: %v", err)
	}
	if err := os.MkdirAll(cfg.WorkspaceRoot, 0o755); err != nil {
		t.Fatalf("mkdir workspace root: %v", err)
	}

	err := ValidateAppLayout(cfg)
	if err == nil {
		t.Fatalf("expected validation error for file path")
	}
	if err.Error() != "command bin dir is not a directory" {
		t.Fatalf("unexpected error: %v", err)
	}
}

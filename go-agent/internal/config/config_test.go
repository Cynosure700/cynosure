package config

import (
	"os"
	"path/filepath"
	"reflect"
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
		"builtin_skills_dir": "custom-skills",
		"command_bin_dir": "custom-bin",
		"command_script_dir": "custom-cmd",
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
	if cfg.BuiltinSkillsDir != filepath.Join(root, "custom-skills") {
		t.Fatalf("unexpected builtin skills dir: %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "custom-bin") {
		t.Fatalf("unexpected command bin dir: %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "custom-cmd") {
		t.Fatalf("unexpected command script dir: %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "shared-workspace") {
		t.Fatalf("unexpected workspace root: %q", cfg.WorkspaceRoot)
	}
	if !reflect.DeepEqual(cfg.WebAllowedTools, []string{"load_skill", "bash"}) {
		t.Fatalf("unexpected web allowed tools: %#v", cfg.WebAllowedTools)
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
	t.Setenv("BUILTIN_SKILLS_DIR", "env-skills")
	t.Setenv("COMMAND_BIN_DIR", "env-bin")
	t.Setenv("COMMAND_SCRIPT_DIR", "env-cmd")
	t.Setenv("WORKSPACE_ROOT", "env-workspace")
	t.Setenv("WEB_ALLOWED_TOOLS", "load_skill,edit_file")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.BuiltinSkillsDir != filepath.Join(root, "env-skills") {
		t.Fatalf("expected env builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "env-bin") {
		t.Fatalf("expected env command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "env-cmd") {
		t.Fatalf("expected env command script dir, got %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "env-workspace") {
		t.Fatalf("expected env workspace root, got %q", cfg.WorkspaceRoot)
	}
	if !reflect.DeepEqual(cfg.WebAllowedTools, []string{"load_skill", "edit_file"}) {
		t.Fatalf("expected env web tools, got %#v", cfg.WebAllowedTools)
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

	if cfg.BuiltinSkillsDir != filepath.Join(root, "workspaces", "skills") {
		t.Fatalf("expected default builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "workspaces", "bin") {
		t.Fatalf("expected default command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "workspaces", "cmd") {
		t.Fatalf("expected default command script dir, got %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "workspace") {
		t.Fatalf("expected default workspace root, got %q", cfg.WorkspaceRoot)
	}
}

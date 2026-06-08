package config

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestMain 为依赖 LoadWebConfig 的用例提供默认的敏感信息环境变量；
// 个别用例可通过 t.Setenv 覆盖。
func TestMain(m *testing.M) {
	_ = os.Setenv("OPENAI_API_KEY", "test-key")
	os.Exit(m.Run())
}

func TestLoadWebConfig_ReadsPathSettingsFromConfigFile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"builtin_skills_dir": "skills",
		"command_bin_dir": "bin",
		"command_script_dir": "cmd",
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

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.AppHome != filepath.Clean(root) {
		t.Fatalf("expected app home %q, got %q", filepath.Clean(root), cfg.AppHome)
	}
	if cfg.BuiltinSkillsDir != filepath.Join(root, "skills") {
		t.Fatalf("unexpected builtin skills dir: %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "bin") {
		t.Fatalf("unexpected command bin dir: %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "cmd") {
		t.Fatalf("unexpected command script dir: %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "shared-workspace") {
		t.Fatalf("unexpected workspace root: %q", cfg.WorkspaceRoot)
	}
	if cfg.SystemPromptPath != filepath.Join(root, "system_prompt.md") {
		t.Fatalf("unexpected system prompt path: %q", cfg.SystemPromptPath)
	}
	if cfg.LogsDir != filepath.Join(root, "logs") {
		t.Fatalf("unexpected logs dir: %q", cfg.LogsDir)
	}
	if !reflect.DeepEqual(cfg.WebAllowedTools, []string{"load_skill", "bash"}) {
		t.Fatalf("unexpected web allowed tools: %#v", cfg.WebAllowedTools)
	}
	if cfg.BashAllowOutsideWorkspace || cfg.BashAllowDangerousCommands {
		t.Fatalf("expected bash safety flags to default false")
	}
}

func TestLoadWebConfig_ReadsServerAndDatastoreSettingsFromConfigFile(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"server_addr": ":9090",
		"allowed_origin": "http://example.test",
		"redis_addr": "redis:6379",
		"redis_db": 3,
		"session_cookie_name": "custom_cookie",
		"session_ttl_minutes": 120
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

	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("DATABASE_URL", "user:pass@tcp(db:3306)/app")
	t.Setenv("REDIS_PASSWORD", "redis-secret")
	t.Setenv("JWT_SECRET", "custom-secret")

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.ServerAddr != ":9090" {
		t.Fatalf("unexpected server addr: %q", cfg.ServerAddr)
	}
	if cfg.AllowedOrigin != "http://example.test" {
		t.Fatalf("unexpected allowed origin: %q", cfg.AllowedOrigin)
	}
	if cfg.DatabaseURL != "user:pass@tcp(db:3306)/app" {
		t.Fatalf("unexpected database url: %q", cfg.DatabaseURL)
	}
	if cfg.RedisAddr != "redis:6379" {
		t.Fatalf("unexpected redis addr: %q", cfg.RedisAddr)
	}
	if cfg.RedisPassword != "redis-secret" {
		t.Fatalf("unexpected redis password: %q", cfg.RedisPassword)
	}
	if cfg.RedisDB != 3 {
		t.Fatalf("unexpected redis db: %d", cfg.RedisDB)
	}
	if cfg.JWTSecret != "custom-secret" {
		t.Fatalf("unexpected jwt secret: %q", cfg.JWTSecret)
	}
	if cfg.CookieName != "custom_cookie" {
		t.Fatalf("unexpected cookie name: %q", cfg.CookieName)
	}
	if cfg.SessionTTLMinutes != 120 {
		t.Fatalf("unexpected session ttl: %d", cfg.SessionTTLMinutes)
	}
}

func TestLoadWebConfig_ReadsSecretsFromEnvAndRejectsMissingAPIKey(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
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

	t.Setenv("OPENAI_API_KEY", "")

	if _, err := LoadWebConfig(); err == nil {
		t.Fatalf("expected missing OPENAI_API_KEY to be rejected")
	}

	t.Setenv("OPENAI_API_KEY", "env-key")
	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}
	if cfg.LLM.APIKey != "env-key" {
		t.Fatalf("expected api key from env, got %q", cfg.LLM.APIKey)
	}
}

func TestLoadWebConfig_ReadsSystemPromptPathFromConfig(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
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

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}
	if cfg.SystemPromptPath != filepath.Join(root, "config-prompt.md") {
		t.Fatalf("expected config system prompt path, got %q", cfg.SystemPromptPath)
	}
}

func TestLoadWebConfig_ReadsBashSafetyFlagsFromConfig(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"bash_allow_outside_workspace": true,
		"bash_allow_dangerous_commands": true
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

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}
	if !cfg.BashAllowOutsideWorkspace {
		t.Fatalf("expected outside-workspace allow flag to be true")
	}
	if !cfg.BashAllowDangerousCommands {
		t.Fatalf("expected dangerous-command allow flag to be true")
	}
}

func TestLoadWebConfig_UsesAppHomeScopedDefaultPaths(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
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

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	if cfg.BuiltinSkillsDir != filepath.Join(root, "skills") {
		t.Fatalf("expected default builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "bin") {
		t.Fatalf("expected default command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "cmd") {
		t.Fatalf("expected default command script dir, got %q", cfg.CommandScriptDir)
	}
	if cfg.WorkspaceRoot != filepath.Join(root, "workspace") {
		t.Fatalf("expected default workspace root, got %q", cfg.WorkspaceRoot)
	}
	expectedTools := []string{"load_skill", "bash", "read_file", "write_file", "edit_file", "todo_write", "update_memory"}
	if !reflect.DeepEqual(cfg.WebAllowedTools, expectedTools) {
		t.Fatalf("expected default web tools %v, got %#v", expectedTools, cfg.WebAllowedTools)
	}
}

func TestLoadWebConfig_DefaultsToAppHomeWorkspaceEvenWhenOutputWorkspaceExists(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
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

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	expectedWorkspace := filepath.Join(root, "workspace")
	if cfg.WorkspaceRoot != expectedWorkspace {
		t.Fatalf("expected app home workspace root %q, got %q", expectedWorkspace, cfg.WorkspaceRoot)
	}
	if cfg.BuiltinSkillsDir != filepath.Join(root, "skills") {
		t.Fatalf("expected default builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "bin") {
		t.Fatalf("expected default command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "cmd") {
		t.Fatalf("expected default command script dir, got %q", cfg.CommandScriptDir)
	}
}

func TestLoadWebConfig_AssetDirsStayUnderAppHomeRegardlessOfWorkspace(t *testing.T) {
	root := t.TempDir()
	configPath := filepath.Join(root, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"model_id": "test-model",
		"app_home": "` + root + `",
		"workspace_root": "` + filepath.Join("custom", "runtime-workspace") + `"
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

	cfg, err := LoadWebConfig()
	if err != nil {
		t.Fatalf("load web config: %v", err)
	}

	expectedWorkspace := filepath.Join(root, "custom", "runtime-workspace")
	if cfg.WorkspaceRoot != expectedWorkspace {
		t.Fatalf("expected explicit workspace root %q, got %q", expectedWorkspace, cfg.WorkspaceRoot)
	}
	if cfg.BuiltinSkillsDir != filepath.Join(root, "skills") {
		t.Fatalf("expected app home builtin skills dir, got %q", cfg.BuiltinSkillsDir)
	}
	if cfg.CommandBinDir != filepath.Join(root, "bin") {
		t.Fatalf("expected app home command bin dir, got %q", cfg.CommandBinDir)
	}
	if cfg.CommandScriptDir != filepath.Join(root, "cmd") {
		t.Fatalf("expected app home command script dir, got %q", cfg.CommandScriptDir)
	}
}

func TestResolveRuntimePaths_ReturnsCanonicalAppHomeDerivedDirs(t *testing.T) {
	root := t.TempDir()

	paths, err := resolveRuntimePaths(root, fileConfig{
		WorkspaceRoot:    "custom/runtime-workspace",
		BuiltinSkillsDir: "skills",
		CommandBinDir:    "bin",
		CommandScriptDir: "cmd",
	})
	if err != nil {
		t.Fatalf("resolve runtime paths: %v", err)
	}

	expectedWorkspace := filepath.Join(root, "custom", "runtime-workspace")
	if paths.workspaceRoot != expectedWorkspace {
		t.Fatalf("expected workspace root %q, got %q", expectedWorkspace, paths.workspaceRoot)
	}
	if paths.builtinSkillsDir != filepath.Join(root, "skills") {
		t.Fatalf("expected builtin skills dir under app home, got %q", paths.builtinSkillsDir)
	}
	if paths.commandBinDir != filepath.Join(root, "bin") {
		t.Fatalf("expected command bin dir under app home, got %q", paths.commandBinDir)
	}
	if paths.commandScriptDir != filepath.Join(root, "cmd") {
		t.Fatalf("expected command script dir under app home, got %q", paths.commandScriptDir)
	}
	if paths.logsDir != filepath.Join(root, "logs") {
		t.Fatalf("expected logs dir under app home, got %q", paths.logsDir)
	}
}

func TestLoadWebConfig_RejectsRuntimeAssetDirsOutsideAppHome(t *testing.T) {
	root := t.TempDir()
	appHome := filepath.Join(root, "app-home")
	if err := os.MkdirAll(appHome, 0o755); err != nil {
		t.Fatalf("mkdir app home: %v", err)
	}
	configPath := filepath.Join(appHome, "config.json")
	configBody := `{
		"base_url": "https://example.com",
		"model_id": "test-model",
		"app_home": "` + appHome + `",
		"builtin_skills_dir": "` + filepath.Join("..", "other-skills") + `"
	}`
	if err := os.WriteFile(configPath, []byte(configBody), 0o644); err != nil {
		t.Fatalf("write config file: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(appHome); err != nil {
		t.Fatalf("chdir to temp root: %v", err)
	}
	defer func() { _ = os.Chdir(oldWD) }()

	_, err = LoadWebConfig()
	if err == nil {
		t.Fatalf("expected runtime asset override outside app home to be rejected")
	}
	if got := err.Error(); got == "" || !strings.Contains(got, "runtime asset dir must stay under app home") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestEnsureAppLayout_CreatesExpectedDirectories(t *testing.T) {
	root := t.TempDir()
	cfg := AppConfig{
		AppHome:          root,
		BuiltinSkillsDir: filepath.Join(root, "skills"),
		CommandBinDir:    filepath.Join(root, "bin"),
		CommandScriptDir: filepath.Join(root, "cmd"),
		WorkspaceRoot:    filepath.Join(root, "workspace"),
		LogsDir:          filepath.Join(root, "logs"),
	}

	if err := EnsureAppLayout(cfg); err != nil {
		t.Fatalf("ensure app layout: %v", err)
	}

	paths := []string{
		root,
		filepath.Join(root, "skills"),
		filepath.Join(root, "bin"),
		filepath.Join(root, "cmd"),
		filepath.Join(root, "workspace"),
		filepath.Join(root, "logs"),
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
		BuiltinSkillsDir: filepath.Join(root, "skills"),
		CommandBinDir:    filepath.Join(root, "bin"),
		CommandScriptDir: filepath.Join(root, "cmd"),
		WorkspaceRoot:    filepath.Join(root, "workspace"),
		LogsDir:          filepath.Join(root, "logs"),
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
		BuiltinSkillsDir: filepath.Join(root, "skills"),
		CommandBinDir:    filePath,
		CommandScriptDir: filepath.Join(root, "cmd"),
		WorkspaceRoot:    filepath.Join(root, "workspace"),
		LogsDir:          filepath.Join(root, "logs"),
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

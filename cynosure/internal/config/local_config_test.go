package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeCynosureSettings(t *testing.T, home string, body string) {
	t.Helper()
	linkDir := filepath.Join(home, ".cynosure")
	if err := os.MkdirAll(linkDir, 0o755); err != nil {
		t.Fatalf("mkdir .cynosure: %v", err)
	}
	if err := os.WriteFile(filepath.Join(linkDir, "settings.json"), []byte(body), 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}
}

func writeCynosureMarkdown(t *testing.T, root string, body string) string {
	t.Helper()
	linkDir := filepath.Join(root, ".cynosure")
	if err := os.MkdirAll(linkDir, 0o755); err != nil {
		t.Fatalf("mkdir .cynosure: %v", err)
	}
	path := filepath.Join(linkDir, "CYNOSURE.MD")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write CYNOSURE.MD: %v", err)
	}
	return path
}

func TestLoadLocalConfigUsesLaunchCWDAsWorkspace(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "config.json"), []byte(`{
		"app_home":".",
		"system_prompt_path":"system_prompt.md",
		"builtin_skills_dir":"skills",
		"command_bin_dir":"bin",
		"command_script_dir":"cmd"
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(tmp, "home")
	t.Setenv("HOME", home)
	t.Setenv("OPENAI_API_KEY", "ignored-env-key")
	writeCynosureSettings(t, home, `{"env":{"open_auth_token":" cynosure-key ","open_model":" cynosure-model ","open_base_url":" https://cynosure.example.com "}}`)

	cwd := filepath.Join(tmp, "project")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg, err := LoadLocalConfig(cwd)
	if err != nil {
		t.Fatalf("LoadLocalConfig returned error: %v", err)
	}
	if cfg.WorkspaceRoot != cwd {
		t.Fatalf("WorkspaceRoot = %q, want %q", cfg.WorkspaceRoot, cwd)
	}
	if cfg.LLM.APIKey != "cynosure-key" {
		t.Fatalf("APIKey = %q, want cynosure-key", cfg.LLM.APIKey)
	}
	if cfg.LLM.ModelID != "cynosure-model" {
		t.Fatalf("ModelID = %q, want cynosure-model", cfg.LLM.ModelID)
	}
	if cfg.LLM.BaseURL != "https://cynosure.example.com" {
		t.Fatalf("BaseURL = %q, want https://cynosure.example.com", cfg.LLM.BaseURL)
	}
	if got := cfg.AllowedTools; len(got) == 0 || got[0] != "load_skill" {
		t.Fatalf("AllowedTools = %#v, want default local tools", got)
	}
}

func TestLoadLocalConfigRequiresCynosureSettingsFields(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "config.json"), []byte(`{"app_home":"."}`), 0o644); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(tmp, "home")
	t.Setenv("HOME", home)
	writeCynosureSettings(t, home, `{"env":{"open_auth_token":"secret-token","open_model":"","open_base_url":""}}`)

	_, err = LoadLocalConfig(tmp)
	if err == nil {
		t.Fatal("LoadLocalConfig returned nil error without required cynosure settings fields")
	}
	message := err.Error()
	if !strings.Contains(message, "open_model") || !strings.Contains(message, "open_base_url") {
		t.Fatalf("error = %q, want missing open_model and open_base_url", message)
	}
	if strings.Contains(message, "secret-token") {
		t.Fatalf("error leaked token: %q", message)
	}
}

func TestLoadLocalConfigRejectsMissingWorkspaceCWD(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "config.json"), []byte(`{"app_home":"."}`), 0o644); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(tmp, "home")
	t.Setenv("HOME", home)
	writeCynosureSettings(t, home, `{"env":{"open_auth_token":"test-key","open_model":"test-model","open_base_url":"https://example.com"}}`)
	missing := filepath.Join(tmp, "missing-project")

	if _, err := LoadLocalConfig(missing); err == nil {
		t.Fatal("LoadLocalConfig returned nil error for missing cwd")
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatalf("missing cwd was created or stat failed unexpectedly: %v", err)
	}
}

func TestLoadCynosureMarkdownContextReadsUserAndWorkspaceFiles(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	userPath := writeCynosureMarkdown(t, home, "\n# User Rule\n全局说明\n")
	workspacePath := writeCynosureMarkdown(t, workspace, "\n# Project Rule\n项目说明\n")

	ctx, err := LoadCynosureMarkdownContext(workspace)
	if err != nil {
		t.Fatalf("LoadCynosureMarkdownContext returned error: %v", err)
	}
	if ctx.UserPath != userPath {
		t.Fatalf("UserPath = %q, want %q", ctx.UserPath, userPath)
	}
	if ctx.UserContent != "# User Rule\n全局说明" {
		t.Fatalf("UserContent = %q", ctx.UserContent)
	}
	if ctx.WorkspacePath != workspacePath {
		t.Fatalf("WorkspacePath = %q, want %q", ctx.WorkspacePath, workspacePath)
	}
	if ctx.WorkspaceContent != "# Project Rule\n项目说明" {
		t.Fatalf("WorkspaceContent = %q", ctx.WorkspaceContent)
	}
}

func TestLoadCynosureMarkdownContextIgnoresMissingAndEmptyFiles(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	writeCynosureMarkdown(t, home, "\n\t  \n")

	ctx, err := LoadCynosureMarkdownContext(workspace)
	if err != nil {
		t.Fatalf("LoadCynosureMarkdownContext returned error: %v", err)
	}
	if ctx.UserContent != "" || ctx.WorkspaceContent != "" {
		t.Fatalf("expected empty context for missing/blank CYNOSURE.MD, got %#v", ctx)
	}
}

func TestLoadCynosureMarkdownContextRejectsDirectoryCynosureMarkdown(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	badPath := filepath.Join(workspace, ".cynosure", "CYNOSURE.MD")
	if err := os.MkdirAll(badPath, 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := LoadCynosureMarkdownContext(workspace)
	if err == nil {
		t.Fatal("LoadCynosureMarkdownContext returned nil error for directory CYNOSURE.MD")
	}
	if !strings.Contains(err.Error(), badPath) {
		t.Fatalf("error = %q, want path %q", err.Error(), badPath)
	}
}

func TestLoadCynosureMarkdownContextRejectsUnreadableCynosureMarkdown(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can read chmod 000 files")
	}
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	badPath := writeCynosureMarkdown(t, workspace, "secret")
	if err := os.Chmod(badPath, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(badPath, 0o644) })

	_, err := LoadCynosureMarkdownContext(workspace)
	if err == nil {
		t.Fatal("LoadCynosureMarkdownContext returned nil error for unreadable CYNOSURE.MD")
	}
	if !strings.Contains(err.Error(), badPath) {
		t.Fatalf("error = %q, want path %q", err.Error(), badPath)
	}
}

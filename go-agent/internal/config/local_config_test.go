package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeLinkSettings(t *testing.T, home string, body string) {
	t.Helper()
	linkDir := filepath.Join(home, ".link")
	if err := os.MkdirAll(linkDir, 0o755); err != nil {
		t.Fatalf("mkdir .link: %v", err)
	}
	if err := os.WriteFile(filepath.Join(linkDir, "settings.json"), []byte(body), 0o644); err != nil {
		t.Fatalf("write settings: %v", err)
	}
}

func writeLinkMarkdown(t *testing.T, root string, body string) string {
	t.Helper()
	linkDir := filepath.Join(root, ".link")
	if err := os.MkdirAll(linkDir, 0o755); err != nil {
		t.Fatalf("mkdir .link: %v", err)
	}
	path := filepath.Join(linkDir, "LINK.MD")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write LINK.MD: %v", err)
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
	writeLinkSettings(t, home, `{"env":{"open_auth_token":" link-key ","open_model":" link-model ","open_base_url":" https://link.example.com "}}`)

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
	if cfg.LLM.APIKey != "link-key" {
		t.Fatalf("APIKey = %q, want link-key", cfg.LLM.APIKey)
	}
	if cfg.LLM.ModelID != "link-model" {
		t.Fatalf("ModelID = %q, want link-model", cfg.LLM.ModelID)
	}
	if cfg.LLM.BaseURL != "https://link.example.com" {
		t.Fatalf("BaseURL = %q, want https://link.example.com", cfg.LLM.BaseURL)
	}
	if got := cfg.AllowedTools; len(got) == 0 || got[0] != "load_skill" {
		t.Fatalf("AllowedTools = %#v, want default local tools", got)
	}
}

func TestLoadLocalConfigRequiresLinkSettingsFields(t *testing.T) {
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
	writeLinkSettings(t, home, `{"env":{"open_auth_token":"secret-token","open_model":"","open_base_url":""}}`)

	_, err = LoadLocalConfig(tmp)
	if err == nil {
		t.Fatal("LoadLocalConfig returned nil error without required link settings fields")
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
	writeLinkSettings(t, home, `{"env":{"open_auth_token":"test-key","open_model":"test-model","open_base_url":"https://example.com"}}`)
	missing := filepath.Join(tmp, "missing-project")

	if _, err := LoadLocalConfig(missing); err == nil {
		t.Fatal("LoadLocalConfig returned nil error for missing cwd")
	}
	if _, err := os.Stat(missing); !os.IsNotExist(err) {
		t.Fatalf("missing cwd was created or stat failed unexpectedly: %v", err)
	}
}

func TestLoadLinkMarkdownContextReadsUserAndWorkspaceFiles(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	userPath := writeLinkMarkdown(t, home, "\n# User Rule\n全局说明\n")
	workspacePath := writeLinkMarkdown(t, workspace, "\n# Project Rule\n项目说明\n")

	ctx, err := LoadLinkMarkdownContext(workspace)
	if err != nil {
		t.Fatalf("LoadLinkMarkdownContext returned error: %v", err)
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

func TestLoadLinkMarkdownContextIgnoresMissingAndEmptyFiles(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	writeLinkMarkdown(t, home, "\n\t  \n")

	ctx, err := LoadLinkMarkdownContext(workspace)
	if err != nil {
		t.Fatalf("LoadLinkMarkdownContext returned error: %v", err)
	}
	if ctx.UserContent != "" || ctx.WorkspaceContent != "" {
		t.Fatalf("expected empty context for missing/blank LINK.MD, got %#v", ctx)
	}
}

func TestLoadLinkMarkdownContextRejectsDirectoryLinkMarkdown(t *testing.T) {
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	badPath := filepath.Join(workspace, ".link", "LINK.MD")
	if err := os.MkdirAll(badPath, 0o755); err != nil {
		t.Fatal(err)
	}

	_, err := LoadLinkMarkdownContext(workspace)
	if err == nil {
		t.Fatal("LoadLinkMarkdownContext returned nil error for directory LINK.MD")
	}
	if !strings.Contains(err.Error(), badPath) {
		t.Fatalf("error = %q, want path %q", err.Error(), badPath)
	}
}

func TestLoadLinkMarkdownContextRejectsUnreadableLinkMarkdown(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can read chmod 000 files")
	}
	tmp := t.TempDir()
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	badPath := writeLinkMarkdown(t, workspace, "secret")
	if err := os.Chmod(badPath, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(badPath, 0o644) })

	_, err := LoadLinkMarkdownContext(workspace)
	if err == nil {
		t.Fatal("LoadLinkMarkdownContext returned nil error for unreadable LINK.MD")
	}
	if !strings.Contains(err.Error(), badPath) {
		t.Fatalf("error = %q, want path %q", err.Error(), badPath)
	}
}

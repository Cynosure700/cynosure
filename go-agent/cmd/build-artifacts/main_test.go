package main

import (
	"path/filepath"
	"testing"
)

func TestResolveBuildConfig_UsesWorkspaceDefaults(t *testing.T) {
	appHome := t.TempDir()

	built, err := resolveBuildConfig(appHome, "", "", "")
	if err != nil {
		t.Fatalf("resolve build config: %v", err)
	}

	if built.AppHome != filepath.Clean(appHome) {
		t.Fatalf("expected app home %q, got %q", filepath.Clean(appHome), built.AppHome)
	}
	if built.CommandSource != filepath.Join(appHome, "cmd") {
		t.Fatalf("expected command source under app home, got %q", built.CommandSource)
	}
	if built.CommandBinDir != filepath.Join(appHome, "workspace", "bin") {
		t.Fatalf("expected command bin dir under workspace/bin, got %q", built.CommandBinDir)
	}
	if built.CommandScriptDir != filepath.Join(appHome, "workspace", "cmd") {
		t.Fatalf("expected command script dir under workspace/cmd, got %q", built.CommandScriptDir)
	}
	if built.WorkspaceRoot != filepath.Join(appHome, "workspace") {
		t.Fatalf("expected workspace root under app home, got %q", built.WorkspaceRoot)
	}
	if built.BuiltinSkillsDir != filepath.Join(appHome, "workspace", "skills") {
		t.Fatalf("expected builtin skills dir under workspace/skills, got %q", built.BuiltinSkillsDir)
	}
}

func TestResolveBuildConfig_PreservesExplicitCommandPaths(t *testing.T) {
	appHome := t.TempDir()
	built, err := resolveBuildConfig(appHome, filepath.Join(appHome, "custom-cmd-src"), filepath.Join(appHome, "custom-bin"), filepath.Join(appHome, "custom-scripts"))
	if err != nil {
		t.Fatalf("resolve build config: %v", err)
	}

	if built.CommandSource != filepath.Join(appHome, "custom-cmd-src") {
		t.Fatalf("expected explicit command source to be preserved, got %q", built.CommandSource)
	}
	if built.CommandBinDir != filepath.Join(appHome, "custom-bin") {
		t.Fatalf("expected explicit command bin dir to be preserved, got %q", built.CommandBinDir)
	}
	if built.CommandScriptDir != filepath.Join(appHome, "custom-scripts") {
		t.Fatalf("expected explicit command script dir to be preserved, got %q", built.CommandScriptDir)
	}
}

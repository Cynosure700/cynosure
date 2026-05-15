package deploy

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestDiscoverGoCommands(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "web"), 0o755); err != nil {
		t.Fatalf("mkdir web: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "helper"), 0o755); err != nil {
		t.Fatalf("mkdir helper: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "web", "main.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "helper", "helper.go"), []byte("package helper"), 0o644); err != nil {
		t.Fatalf("write helper.go: %v", err)
	}

	commands, err := DiscoverGoCommands(root)
	if err != nil {
		t.Fatalf("discover commands: %v", err)
	}
	expected := []string{"web"}
	if !reflect.DeepEqual(commands, expected) {
		t.Fatalf("expected %v, got %v", expected, commands)
	}
}

func TestCopyScriptAssets(t *testing.T) {
	sourceRoot := t.TempDir()
	targetRoot := t.TempDir()

	scriptDir := filepath.Join(sourceRoot, "helpers")
	if err := os.MkdirAll(scriptDir, 0o755); err != nil {
		t.Fatalf("mkdir helpers: %v", err)
	}
	pythonPath := filepath.Join(scriptDir, "demo.py")
	if err := os.WriteFile(pythonPath, []byte("print('ok')\n"), 0o755); err != nil {
		t.Fatalf("write python script: %v", err)
	}
	if err := os.WriteFile(filepath.Join(scriptDir, "ignore.go"), []byte("package main"), 0o644); err != nil {
		t.Fatalf("write go file: %v", err)
	}

	if err := CopyScriptAssets(sourceRoot, targetRoot); err != nil {
		t.Fatalf("copy script assets: %v", err)
	}

	copiedPath := filepath.Join(targetRoot, "helpers", "demo.py")
	data, err := os.ReadFile(copiedPath)
	if err != nil {
		t.Fatalf("read copied script: %v", err)
	}
	if string(data) != "print('ok')\n" {
		t.Fatalf("unexpected script contents: %q", string(data))
	}
	if _, err := os.Stat(filepath.Join(targetRoot, "helpers", "ignore.go")); !os.IsNotExist(err) {
		t.Fatalf("go source should not be copied as a script asset")
	}
}

func TestBuildCommandArtifacts_BuildsCommandsAndCopiesScripts(t *testing.T) {
	appHome := t.TempDir()
	commandSource := filepath.Join(appHome, "cmd")
	if err := os.MkdirAll(filepath.Join(commandSource, "demo"), 0o755); err != nil {
		t.Fatalf("mkdir demo command: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandSource, "demo", "main.go"), []byte("package main\nfunc main() {}\n"), 0o644); err != nil {
		t.Fatalf("write main.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(commandSource, "helper.py"), []byte("print('ok')\n"), 0o755); err != nil {
		t.Fatalf("write script asset: %v", err)
	}

	logPath := filepath.Join(appHome, "build.log")
	fakeGo := filepath.Join(appHome, "fake-go.sh")
	fakeGoScript := strings.Join([]string{
		"#!/bin/sh",
		"set -eu",
		"if [ \"$1\" != \"build\" ] || [ \"$2\" != \"-o\" ]; then",
		"  exit 99",
		"fi",
		"output=\"$3\"",
		"pkg=\"$4\"",
		"printf '%s|%s\\n' \"$PWD\" \"$pkg\" > \"$BUILD_LOG\"",
		"printf 'binary for %s\\n' \"$pkg\" > \"$output\"",
	}, "\n")
	if err := os.WriteFile(fakeGo, []byte(fakeGoScript), 0o755); err != nil {
		t.Fatalf("write fake go binary: %v", err)
	}
	t.Setenv("BUILD_LOG", logPath)

	binDir := filepath.Join(appHome, "bin")
	scriptDir := filepath.Join(appHome, "runtime-scripts")
	if err := BuildCommandArtifacts(BuildOptions{
		AppHome:          appHome,
		CommandSource:    commandSource,
		CommandBinDir:    binDir,
		CommandScriptDir: scriptDir,
		GoBinary:         fakeGo,
	}); err != nil {
		t.Fatalf("build command artifacts: %v", err)
	}

	binaryData, err := os.ReadFile(filepath.Join(binDir, "demo"))
	if err != nil {
		t.Fatalf("read built binary: %v", err)
	}
	if string(binaryData) != "binary for ./cmd/demo\n" {
		t.Fatalf("unexpected built binary contents: %q", string(binaryData))
	}

	logData, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read build log: %v", err)
	}
	if string(logData) != appHome+"|./cmd/demo\n" {
		t.Fatalf("expected go build to run from app home with cmd package path, got %q", string(logData))
	}

	scriptData, err := os.ReadFile(filepath.Join(scriptDir, "helper.py"))
	if err != nil {
		t.Fatalf("read copied helper script: %v", err)
	}
	if string(scriptData) != "print('ok')\n" {
		t.Fatalf("unexpected copied script contents: %q", string(scriptData))
	}
}

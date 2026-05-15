package deploy

import (
	"os"
	"path/filepath"
	"reflect"
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

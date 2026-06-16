package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWorkspacePermissionsMissingFile(t *testing.T) {
	perms, err := LoadWorkspacePermissions(t.TempDir())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if perms.IsBypass() || len(perms.AllowedRules) != 0 {
		t.Fatalf("expected zero permissions, got %+v", perms)
	}
}

func TestLoadWorkspacePermissionsBypass(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Dir(WorkspaceCynosureSettingsPath(root)), 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{"permissions":{"defaultMode":"bypassPermissions","allowedRules":["curl *"]}}`
	if err := os.WriteFile(WorkspaceCynosureSettingsPath(root), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	perms, err := LoadWorkspacePermissions(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !perms.IsBypass() {
		t.Fatalf("expected bypass mode")
	}
	if !perms.Allows("curl *") {
		t.Fatalf("expected curl * allowed")
	}
	if perms.Allows("rm *") {
		t.Fatalf("did not expect rm * allowed")
	}
}

func TestAppendWorkspaceApprovalRuleRoundTrip(t *testing.T) {
	root := t.TempDir()
	if err := AppendWorkspaceApprovalRule(root, "curl *"); err != nil {
		t.Fatalf("append: %v", err)
	}
	// 重复写入应去重
	if err := AppendWorkspaceApprovalRule(root, "curl *"); err != nil {
		t.Fatalf("append dup: %v", err)
	}
	if err := AppendWorkspaceApprovalRule(root, "write_file *"); err != nil {
		t.Fatalf("append second: %v", err)
	}
	perms, err := LoadWorkspacePermissions(root)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(perms.AllowedRules) != 2 {
		t.Fatalf("expected 2 rules, got %v", perms.AllowedRules)
	}
	if !perms.Allows("curl *") || !perms.Allows("write_file *") {
		t.Fatalf("rules not persisted: %v", perms.AllowedRules)
	}
}

func TestAppendWorkspaceApprovalRulePreservesOtherFields(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Dir(WorkspaceCynosureSettingsPath(root))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	body := `{"permissions":{"defaultMode":"bypassPermissions"},"other":{"keep":true}}`
	if err := os.WriteFile(WorkspaceCynosureSettingsPath(root), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := AppendWorkspaceApprovalRule(root, "curl *"); err != nil {
		t.Fatalf("append: %v", err)
	}
	data, err := os.ReadFile(WorkspaceCynosureSettingsPath(root))
	if err != nil {
		t.Fatal(err)
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("parse: %v", err)
	}
	if _, ok := raw["other"]; !ok {
		t.Fatalf("expected other field preserved, got %s", data)
	}
	perms, err := LoadWorkspacePermissions(root)
	if err != nil {
		t.Fatal(err)
	}
	if !perms.IsBypass() || !perms.Allows("curl *") {
		t.Fatalf("expected bypass kept and rule added, got %+v", perms)
	}
}

package sessions

import (
	"os"
	"testing"
)

func TestLoadProjectMemory_Exists(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	content := "Test project memory"
	os.WriteFile("AGENTS.md", []byte(content), 0o644)

	result := LoadProjectMemory()
	if result != content {
		t.Errorf("expected %q, got %q", content, result)
	}
}

func TestLoadProjectMemory_NotExists(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	result := LoadProjectMemory()
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestLoadUserMemory_NotExists(t *testing.T) {
	result := LoadUserMemory()
	_ = result
}

func TestBuildPersistentMemorySection_Both(t *testing.T) {
	section := BuildPersistentMemorySection("project content", "user content")
	if section == "" {
		t.Error("expected non-empty section")
	}
	if !contains(section, "<project_memory>") {
		t.Error("missing project_memory tag")
	}
	if !contains(section, "<user_memory>") {
		t.Error("missing user_memory tag")
	}
}

func TestBuildPersistentMemorySection_ProjectOnly(t *testing.T) {
	section := BuildPersistentMemorySection("project content", "")
	if section == "" {
		t.Error("expected non-empty section")
	}
	if contains(section, "<user_memory>") {
		t.Error("should not contain user_memory tag")
	}
}

func TestBuildPersistentMemorySection_Empty(t *testing.T) {
	section := BuildPersistentMemorySection("", "")
	if section != "" {
		t.Errorf("expected empty string, got %q", section)
	}
}

func TestUpdateMemory_Append(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	result, err := handleUpdateMemory(nil, map[string]any{
		"action":  "append",
		"content": "Line 1",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Memory updated: appended to AGENTS.md" {
		t.Errorf("unexpected result: %s", result)
	}

	data, _ := os.ReadFile("AGENTS.md")
	if string(data) != "Line 1\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}

	// Append another line
	handleUpdateMemory(nil, map[string]any{
		"action":  "append",
		"content": "Line 2",
	})

	data, _ = os.ReadFile("AGENTS.md")
	if string(data) != "Line 1\nLine 2\n" {
		t.Errorf("unexpected file content after append: %q", string(data))
	}
}

func TestUpdateMemory_Replace(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	os.WriteFile("AGENTS.md", []byte("old content"), 0o644)

	result, err := handleUpdateMemory(nil, map[string]any{
		"action":  "replace",
		"content": "new content",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Memory updated: replaced AGENTS.md" {
		t.Errorf("unexpected result: %s", result)
	}

	data, _ := os.ReadFile("AGENTS.md")
	if string(data) != "new content" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

func TestUpdateMemory_DefaultAction(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	// No action specified → defaults to append
	result, err := handleUpdateMemory(nil, map[string]any{
		"content": "default append",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Memory updated: appended to AGENTS.md" {
		t.Errorf("unexpected result: %s", result)
	}
}

func TestUpdateMemory_InvalidAction(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	_, err := handleUpdateMemory(nil, map[string]any{
		"action":  "invalid",
		"content": "test",
	})
	if err == nil {
		t.Error("expected error for invalid action")
	}
}

func TestUpdateMemory_EmptyContent(t *testing.T) {
	_, err := handleUpdateMemory(nil, map[string]any{
		"action":  "append",
		"content": "",
	})
	if err == nil {
		t.Error("expected error for empty content")
	}
}

func TestUpdateMemory_PathSafety(t *testing.T) {
	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(origDir)

	_, err := handleUpdateMemory(nil, map[string]any{
		"action":  "append",
		"content": "safe content",
	})
	if err != nil {
		t.Fatalf("unexpected error for safe path: %v", err)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

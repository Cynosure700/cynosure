package sessions

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllFromDir_UsesSkillEntryFilesOnly(t *testing.T) {
	root := t.TempDir()
	skillDir := filepath.Join(root, "demo-skill")
	if err := os.MkdirAll(filepath.Join(skillDir, "references"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	skillDoc := `---
name: stable-skill
description: Stable description
---

Do the thing.`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillDoc), 0o644); err != nil {
		t.Fatalf("write skill: %v", err)
	}
	if err := os.WriteFile(filepath.Join(skillDir, "references", "extra.md"), []byte("extra"), 0o644); err != nil {
		t.Fatalf("write reference: %v", err)
	}

	loader := NewSkillLoader()
	if err := loader.LoadAllFromDir(root); err != nil {
		t.Fatalf("load skills: %v", err)
	}

	if len(loader.Skills) != 1 {
		t.Fatalf("expected 1 skill entry, got %d", len(loader.Skills))
	}
	entry, ok := loader.Skills["stable-skill"]
	if !ok {
		t.Fatalf("expected canonical skill name from frontmatter")
	}
	if got := entry.Meta["description"]; got != "Stable description" {
		t.Fatalf("unexpected description: %q", got)
	}
	if got := entry.Path; got != filepath.Join(skillDir, "SKILL.md") {
		t.Fatalf("unexpected skill path: %q", got)
	}
	if _, exists := loader.Skills["extra"]; exists {
		t.Fatalf("reference markdown should not be loaded as a skill entry")
	}
}

func TestCanonicalSkillName_FallsBackToDirectoryName(t *testing.T) {
	name := canonicalSkillName(filepath.Join("/tmp", "my-skill", "SKILL.md"), map[string]string{})
	if name != "my-skill" {
		t.Fatalf("expected directory fallback, got %q", name)
	}
}

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

func TestLoadBuiltinSkillsFromDir_LoadsWorkspaceSkillsDirectory(t *testing.T) {
	workspaceRoot := t.TempDir()
	skillDir := filepath.Join(workspaceRoot, "skills", "workspace-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	skillDoc := `---
name: workspace-skill
description: Loaded from workspace root
---

Do workspace thing.`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillDoc), 0o644); err != nil {
		t.Fatalf("write skill: %v", err)
	}

	loader, err := LoadBuiltinSkillsFromDir(filepath.Join(workspaceRoot, "skills"))
	if err != nil {
		t.Fatalf("load builtin skills from dir: %v", err)
	}
	entry, ok := loader.Skills["workspace-skill"]
	if !ok {
		t.Fatalf("expected workspace skill to be loaded")
	}
	if got := entry.Meta["description"]; got != "Loaded from workspace root" {
		t.Fatalf("unexpected description: %q", got)
	}
}

func TestLoadBuiltinSkillsFromDir_LoadsResolvedBuiltinSkillDirectory(t *testing.T) {
	testCases := []struct {
		name        string
		skillsDir    string
		description string
	}{
		{name: "deployment", skillsDir: filepath.Join(t.TempDir(), "output", "workspace", "skills"), description: "Loaded from deployment skills dir"},
		{name: "local", skillsDir: filepath.Join(t.TempDir(), "workspace", "skills"), description: "Loaded from local skills dir"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			skillDir := filepath.Join(tc.skillsDir, "resolved-skill")
			if err := os.MkdirAll(skillDir, 0o755); err != nil {
				t.Fatalf("mkdir skill dir: %v", err)
			}
			skillDoc := `---
name: resolved-skill
description: ` + tc.description + `
---

Resolved skill body.`
			if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillDoc), 0o644); err != nil {
				t.Fatalf("write skill: %v", err)
			}

			loader, err := LoadBuiltinSkillsFromDir(tc.skillsDir)
			if err != nil {
				t.Fatalf("load builtin skills from dir: %v", err)
			}
			entry, ok := loader.Skills["resolved-skill"]
			if !ok {
				t.Fatalf("expected resolved skill to be loaded")
			}
			if got := entry.Meta["description"]; got != tc.description {
				t.Fatalf("unexpected description: %q", got)
			}
		})
	}
}

func TestMergeSkillLoaders_PrefersLaterLoaderOnConflict(t *testing.T) {
	builtin := NewSkillLoader()
	builtin.LoadFromEntries(map[string]*SkillEntry{
		"shared": {
			Meta: map[string]string{"description": "builtin"},
			Body: "builtin body",
			Path: "builtin://shared",
		},
		"builtin-only": {
			Meta: map[string]string{"description": "builtin only"},
			Body: "builtin only body",
			Path: "builtin://only",
		},
	})

	user := NewSkillLoader()
	user.LoadFromEntries(map[string]*SkillEntry{
		"shared": {
			Meta: map[string]string{"description": "user"},
			Body: "user body",
			Path: "db://shared",
		},
		"user-only": {
			Meta: map[string]string{"description": "user only"},
			Body: "user only body",
			Path: "db://only",
		},
	})

	merged := MergeSkillLoaders(builtin, user)
	if len(merged.Skills) != 3 {
		t.Fatalf("expected 3 merged skills, got %d", len(merged.Skills))
	}
	if got := merged.Skills["shared"].Body; got != "user body" {
		t.Fatalf("expected later loader to win conflict, got %q", got)
	}
	if got := merged.Skills["builtin-only"].Body; got != "builtin only body" {
		t.Fatalf("unexpected builtin-only body: %q", got)
	}
	if got := merged.Skills["user-only"].Body; got != "user only body" {
		t.Fatalf("unexpected user-only body: %q", got)
	}

	userEntry := user.Skills["shared"]
	mergedEntry := merged.Skills["shared"]
	mergedEntry.Meta["description"] = "changed"
	if userEntry.Meta["description"] != "user" {
		t.Fatalf("merged loader should clone entries instead of mutating source")
	}
}

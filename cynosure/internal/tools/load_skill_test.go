package tools

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cynosure/internal/sessions"
)

func TestLoadSkillReturnsBaseDirForDiskSkill(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("user home: %v", err)
	}
	skillPath := filepath.Join(home, ".cynosure", "skills", "skill-creator", "SKILL.md")
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"skill-creator": {Meta: map[string]string{"description": "Disk skill"}, Body: "disk body", Source: "user", Path: skillPath},
	})
	snapshot := NewSkillSnapshot(loader, nil)

	loaded, err := snapshot.LoadSkill("skill-creator")
	if err != nil {
		t.Fatalf("LoadSkill: %v", err)
	}
	wantBase := filepath.Join(home, ".cynosure", "skills", "skill-creator")
	if loaded.BaseDir != wantBase {
		t.Fatalf("BaseDir = %q, want %q", loaded.BaseDir, wantBase)
	}
	rendered := renderLoadedSkill(loaded)
	if !strings.Contains(rendered, "Base directory for this skill: "+wantBase) {
		t.Fatalf("rendered output missing base dir line, got %q", rendered)
	}
}

func TestLoadSkillMaterializesBuiltinSkillAndReturnsSystemBaseDir(t *testing.T) {
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"skill-creator": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Source: "builtin", Path: "skill-creator/SKILL.md"},
	})
	snapshot := NewSkillSnapshot(nil, loader)
	var materializedName string
	snapshot.BuiltinMaterializer = func(name string) (string, error) {
		materializedName = name
		return "/home/u/.cynosure/system/skills/" + name, nil
	}

	loaded, err := snapshot.LoadSkill("skill-creator")
	if err != nil {
		t.Fatalf("LoadSkill: %v", err)
	}
	if materializedName != "skill-creator" {
		t.Fatalf("expected materializer called with skill-creator, got %q", materializedName)
	}
	if loaded.BaseDir != "/home/u/.cynosure/system/skills/skill-creator" {
		t.Fatalf("BaseDir = %q, want materialized system path", loaded.BaseDir)
	}
	rendered := renderLoadedSkill(loaded)
	if !strings.Contains(rendered, "Base directory for this skill: /home/u/.cynosure/system/skills/skill-creator") {
		t.Fatalf("rendered output missing materialized base dir, got %q", rendered)
	}
	if !strings.Contains(rendered, "builtin body") {
		t.Fatalf("rendered output missing skill body, got %q", rendered)
	}
}

func TestLoadSkillBuiltinWithoutMaterializerStillReturnsContent(t *testing.T) {
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(map[string]*sessions.SkillEntry{
		"skill-creator": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Source: "builtin", Path: "skill-creator/SKILL.md"},
	})
	snapshot := NewSkillSnapshot(nil, loader)

	loaded, err := snapshot.LoadSkill("skill-creator")
	if err != nil {
		t.Fatalf("LoadSkill: %v", err)
	}
	rendered := renderLoadedSkill(loaded)
	if !strings.Contains(rendered, "builtin body") {
		t.Fatalf("expected skill body even without materializer, got %q", rendered)
	}
}

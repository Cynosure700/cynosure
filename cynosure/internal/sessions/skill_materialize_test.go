package sessions

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestMaterializeBuiltinSkillWritesEntireTree(t *testing.T) {
	fsys := fstest.MapFS{
		"skill-creator/SKILL.md":              {Data: []byte("---\nname: skill-creator\n---\nbody")},
		"skill-creator/scripts/run.py":        {Data: []byte("print('hi')")},
		"skill-creator/references/schemas.md": {Data: []byte("schemas")},
		"other-skill/SKILL.md":                {Data: []byte("other")},
	}
	destRoot := filepath.Join(t.TempDir(), "system", "skills")

	base, err := MaterializeBuiltinSkill(fsys, "skill-creator", destRoot)
	if err != nil {
		t.Fatalf("MaterializeBuiltinSkill: %v", err)
	}
	wantBase := filepath.Join(destRoot, "skill-creator")
	if base != wantBase {
		t.Fatalf("base = %q, want %q", base, wantBase)
	}
	for rel, want := range map[string]string{
		filepath.Join("SKILL.md"):                 "---\nname: skill-creator\n---\nbody",
		filepath.Join("scripts", "run.py"):        "print('hi')",
		filepath.Join("references", "schemas.md"): "schemas",
	} {
		data, err := os.ReadFile(filepath.Join(wantBase, rel))
		if err != nil {
			t.Fatalf("read %s: %v", rel, err)
		}
		if string(data) != want {
			t.Fatalf("file %s = %q, want %q", rel, string(data), want)
		}
	}
	// 不应连带落盘其它 skill。
	if _, err := os.Stat(filepath.Join(destRoot, "other-skill")); !os.IsNotExist(err) {
		t.Fatalf("expected other-skill not materialized, stat err = %v", err)
	}
}

func TestMaterializeBuiltinSkillSkipsWhenDestExists(t *testing.T) {
	fsys := fstest.MapFS{
		"skill-creator/SKILL.md": {Data: []byte("fresh embedded body")},
	}
	destRoot := filepath.Join(t.TempDir(), "system", "skills")
	existing := filepath.Join(destRoot, "skill-creator")
	if err := os.MkdirAll(existing, 0o755); err != nil {
		t.Fatal(err)
	}
	sentinel := filepath.Join(existing, "SKILL.md")
	if err := os.WriteFile(sentinel, []byte("user-modified body"), 0o644); err != nil {
		t.Fatal(err)
	}

	base, err := MaterializeBuiltinSkill(fsys, "skill-creator", destRoot)
	if err != nil {
		t.Fatalf("MaterializeBuiltinSkill: %v", err)
	}
	if base != existing {
		t.Fatalf("base = %q, want %q", base, existing)
	}
	data, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatalf("read sentinel: %v", err)
	}
	if string(data) != "user-modified body" {
		t.Fatalf("expected existing dir untouched, got %q", string(data))
	}
}

func TestMaterializeBuiltinSkillErrorsWhenSkillMissing(t *testing.T) {
	fsys := fstest.MapFS{
		"skill-creator/SKILL.md": {Data: []byte("body")},
	}
	destRoot := filepath.Join(t.TempDir(), "system", "skills")

	if _, err := MaterializeBuiltinSkill(fsys, "does-not-exist", destRoot); err == nil {
		t.Fatal("expected error for missing skill subtree, got nil")
	}
}

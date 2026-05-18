package sessions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type SkillEntry struct {
	Meta map[string]string
	Body string
	Path string
}

type SkillLoader struct {
	mu     sync.RWMutex
	Skills map[string]*SkillEntry
}

const skillEntryFileName = "SKILL.md"

func NewSkillLoader() *SkillLoader {
	return &SkillLoader{Skills: make(map[string]*SkillEntry)}
}

func LoadBuiltinSkillsFromDir(dir string) (*SkillLoader, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil, fmt.Errorf("builtin skills dir is required")
	}
	loader := NewSkillLoader()
	if err := loader.LoadAllFromDir(dir); err != nil {
		return nil, err
	}
	return loader, nil
}

func (sl *SkillLoader) LoadAllFromDir(dir string) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.Skills = make(map[string]*SkillEntry)

	err := filepath.WalkDir(dir, func(fullPath string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || d.Name() != skillEntryFileName {
			return nil
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil
		}

		meta, body := parseFrontmatter(string(data))
		name := canonicalSkillName(fullPath, meta)
		sl.Skills[name] = &SkillEntry{
			Meta: meta,
			Body: body,
			Path: fullPath,
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return nil
}

func (sl *SkillLoader) LoadFromEntries(entries map[string]*SkillEntry) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.Skills = make(map[string]*SkillEntry, len(entries))
	for name, entry := range entries {
		sl.Skills[name] = cloneSkillEntry(entry)
	}
}

func (sl *SkillLoader) Entries() map[string]*SkillEntry {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	entries := make(map[string]*SkillEntry, len(sl.Skills))
	for name, entry := range sl.Skills {
		entries[name] = cloneSkillEntry(entry)
	}
	return entries
}

func canonicalSkillName(fullPath string, meta map[string]string) string {
	if name := strings.TrimSpace(meta["name"]); name != "" {
		return name
	}
	return filepath.Base(filepath.Dir(fullPath))
}

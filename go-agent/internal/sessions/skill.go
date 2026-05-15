package sessions

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"nano_cc/internal/tools"
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

var Skills = &SkillLoader{Skills: make(map[string]*SkillEntry)}

const skillsDir = "skills"
const skillEntryFileName = "SKILL.md"

func (sl *SkillLoader) LoadAll() error {
	return sl.LoadAllFromDir(skillsDir)
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

func NewSkillLoader() *SkillLoader {
	return &SkillLoader{Skills: make(map[string]*SkillEntry)}
}

func LoadBuiltinSkillsFromWorkspaceRoot(workspaceRoot string) (*SkillLoader, error) {
	workspaceRoot = strings.TrimSpace(workspaceRoot)
	if workspaceRoot == "" {
		return nil, fmt.Errorf("workspace root is required")
	}
	loader := NewSkillLoader()
	if err := loader.LoadAllFromDir(filepath.Join(workspaceRoot, skillsDir)); err != nil {
		return nil, err
	}
	return loader, nil
}

func canonicalSkillName(fullPath string, meta map[string]string) string {
	if name := strings.TrimSpace(meta["name"]); name != "" {
		return name
	}
	return filepath.Base(filepath.Dir(fullPath))
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

func MergeSkillLoaders(loaders ...*SkillLoader) *SkillLoader {
	merged := NewSkillLoader()
	entries := make(map[string]*SkillEntry)
	for _, loader := range loaders {
		if loader == nil {
			continue
		}
		for name, entry := range loader.Entries() {
			entries[name] = entry
		}
	}
	merged.LoadFromEntries(entries)
	return merged
}

func cloneSkillEntry(entry *SkillEntry) *SkillEntry {
	if entry == nil {
		return nil
	}
	meta := make(map[string]string, len(entry.Meta))
	for key, value := range entry.Meta {
		meta[key] = value
	}
	return &SkillEntry{
		Meta: meta,
		Body: entry.Body,
		Path: entry.Path,
	}
}

func parseFrontmatter(text string) (map[string]string, string) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "---") {
		return map[string]string{}, text
	}

	rest := text[3:]
	idx := strings.Index(rest, "\n---")
	if idx == -1 {
		return map[string]string{}, text
	}

	fm := rest[:idx]
	body := strings.TrimSpace(rest[idx+4:])

	meta := make(map[string]string)
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			meta[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	return meta, body
}

func (sl *SkillLoader) GetDescriptions() string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	if len(sl.Skills) == 0 {
		return ""
	}

	names := make([]string, 0, len(sl.Skills))
	for name := range sl.Skills {
		names = append(names, name)
	}
	sort.Strings(names)

	var lines []string
	for _, name := range names {
		skill := sl.Skills[name]
		desc := skill.Meta["description"]
		tags := skill.Meta["tags"]
		line := fmt.Sprintf("- %s: %s", name, desc)
		if tags != "" {
			line += fmt.Sprintf(" [%s]", tags)
		}
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func (sl *SkillLoader) GetContent(name string) (string, error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	skill, ok := sl.Skills[name]
	if !ok {
		available := make([]string, 0, len(sl.Skills))
		for n := range sl.Skills {
			available = append(available, n)
		}
		sort.Strings(available)
		return "", fmt.Errorf("unknown skill %q. Available: %s", name, strings.Join(available, ", "))
	}

	return fmt.Sprintf("<skill name=\"%s\">\n%s\n</skill>", name, skill.Body), nil
}

func init() {
	tools.SetHandler("load_skill", func(ctx context.Context, args map[string]any) (string, error) {
		name, _ := args["name"].(string)
		if name == "" {
			return "", fmt.Errorf("skill name is required")
		}
		content, err := Skills.GetContent(name)
		if err != nil {
			return err.Error(), nil
		}
		return content, nil
	})
}

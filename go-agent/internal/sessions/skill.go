package sessions

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
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

const (
	skillEntryFileName      = "SKILL.md"
	defaultSkillDescription = "No description provided."
)

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
	lines = append(lines, "<skills>")
	for _, name := range names {
		skill := sl.Skills[name]
		desc := strings.TrimSpace(skill.Meta["description"])
		if desc == "" {
			desc = defaultSkillDescription
		}
		tags := strings.TrimSpace(skill.Meta["tags"])
		lines = append(lines,
			"<skill>",
			fmt.Sprintf("<name>%s</name>", html.EscapeString(name)),
			fmt.Sprintf("<description>%s</description>", html.EscapeString(desc)),
		)
		if tags != "" {
			lines = append(lines, fmt.Sprintf("<tags>%s</tags>", html.EscapeString(tags)))
		}
		lines = append(lines, "</skill>")
	}
	lines = append(lines, "</skills>")

	return strings.Join(lines, "\n")
}

func (sl *SkillLoader) GetContent(name string) (string, error) {
	entry, err := sl.GetEntry(name)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("<skill name=\"%s\">\n%s\n</skill>", name, entry.Body), nil
}

func (sl *SkillLoader) GetEntry(name string) (*SkillEntry, error) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	skill, ok := sl.Skills[name]
	if !ok {
		available := make([]string, 0, len(sl.Skills))
		for n := range sl.Skills {
			available = append(available, n)
		}
		sort.Strings(available)
		return nil, fmt.Errorf("unknown skill %q. Available: %s", name, strings.Join(available, ", "))
	}

	return cloneSkillEntry(skill), nil
}

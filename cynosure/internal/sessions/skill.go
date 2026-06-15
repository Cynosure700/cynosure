package sessions

import (
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

type SkillEntry struct {
	Meta   map[string]string
	Body   string
	Path   string
	Source string
}

type SkillLoader struct {
	mu     sync.RWMutex
	Skills map[string]*SkillEntry
}

const (
	skillEntryFileName      = "SKILL.md"
	defaultSkillDescription = "No description provided."
)

type SkillDir struct {
	Path   string
	Source string
}

type SkillSummary struct {
	Name        string
	Description string
	Source      string
	Path        string
}

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
	return sl.LoadAllFromDirWithSource(dir, "")
}

func (sl *SkillLoader) LoadAllFromDirWithSource(dir string, source string) error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.Skills = make(map[string]*SkillEntry)
	seenInDir := make(map[string]string)

	err := filepath.WalkDir(dir, func(fullPath string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || !isSkillEntryFile(d.Name()) {
			return nil
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil
		}

		meta, body := parseFrontmatter(string(data))
		name := canonicalSkillName(fullPath, meta)
		if previous, exists := seenInDir[name]; exists {
			return fmt.Errorf("duplicate skill %q in %s and %s", name, previous, fullPath)
		}
		seenInDir[name] = fullPath
		sl.Skills[name] = &SkillEntry{
			Meta:   meta,
			Body:   body,
			Path:   fullPath,
			Source: source,
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

func isSkillEntryFile(name string) bool {
	return strings.EqualFold(name, skillEntryFileName)
}

// LoadSkillsFromFS 从一个 fs.FS（如 go:embed 文件系统）加载内置 skills，
// 返回带来源标记的 loader。用于把嵌入二进制的内置 skills 接入加载链。
func LoadSkillsFromFS(fsys fs.FS, source string) (*SkillLoader, error) {
	loader := NewSkillLoader()
	if fsys == nil {
		return loader, nil
	}
	entries := make(map[string]*SkillEntry)
	seenInDir := make(map[string]string)
	err := fs.WalkDir(fsys, ".", func(fullPath string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || !isSkillEntryFile(d.Name()) {
			return nil
		}
		data, err := fs.ReadFile(fsys, fullPath)
		if err != nil {
			return nil
		}
		meta, body := parseFrontmatter(string(data))
		name := canonicalSkillName(fullPath, meta)
		if previous, exists := seenInDir[name]; exists {
			return fmt.Errorf("duplicate skill %q in %s and %s", name, previous, fullPath)
		}
		seenInDir[name] = fullPath
		entries[name] = &SkillEntry{
			Meta:   meta,
			Body:   body,
			Path:   fullPath,
			Source: strings.TrimSpace(source),
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	loader.LoadFromEntries(entries)
	return loader, nil
}

func LoadSkillsFromDirs(dirs []SkillDir) (*SkillLoader, error) {
	merged := NewSkillLoader()
	entries := make(map[string]*SkillEntry)
	for _, dir := range dirs {
		path := strings.TrimSpace(dir.Path)
		if path == "" {
			continue
		}
		loader := NewSkillLoader()
		if err := loader.LoadAllFromDirWithSource(path, strings.TrimSpace(dir.Source)); err != nil {
			return nil, err
		}
		for name, entry := range loader.Entries() {
			entries[name] = entry
		}
	}
	merged.LoadFromEntries(entries)
	return merged, nil
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
		Meta:   meta,
		Body:   entry.Body,
		Path:   entry.Path,
		Source: entry.Source,
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

func (sl *SkillLoader) Summaries() []SkillSummary {
	if sl == nil {
		return nil
	}
	entries := sl.Entries()
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	summaries := make([]SkillSummary, 0, len(names))
	for _, name := range names {
		entry := entries[name]
		description := strings.TrimSpace(entry.Meta["description"])
		if description == "" {
			description = defaultSkillDescription
		}
		summaries = append(summaries, SkillSummary{Name: name, Description: description, Source: entry.Source, Path: entry.Path})
	}
	return summaries
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

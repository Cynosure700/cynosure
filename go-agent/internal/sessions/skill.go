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

func (sl *SkillLoader) LoadAll() error {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.Skills = make(map[string]*SkillEntry)

	err := filepath.WalkDir(skillsDir, func(fullPath string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		relPath, _ := filepath.Rel(skillsDir, fullPath)
		name := strings.TrimSuffix(relPath, ".md")

		data, err := os.ReadFile(fullPath)
		if err != nil {
			return nil
		}

		meta, body := parseFrontmatter(string(data))
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
		sl.Skills[name] = entry
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

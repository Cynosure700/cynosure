package sessions

import (
	"fmt"
	"sort"
	"strings"
)

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

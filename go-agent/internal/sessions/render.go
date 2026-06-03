package sessions

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

const defaultSkillDescription = "No description provided."

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

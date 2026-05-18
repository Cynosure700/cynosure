package app

import (
	"fmt"
	"strings"

	"nano_cc/internal/sessions"
	"nano_cc/internal/web/storage"
)

func slugify(input string) string {
	value := strings.ToLower(strings.TrimSpace(input))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			return r
		default:
			return '-'
		}
	}, value)
	value = strings.Trim(value, "-")
	if value == "" {
		return "skill"
	}
	return value
}

func normalizeSkillStatus(status string) string {
	switch status {
	case "enabled", "disabled", "draft":
		return status
	default:
		return "draft"
	}
}

func validateSkill(skill storage.Skill) error {
	if strings.TrimSpace(skill.Name) == "" {
		return fmt.Errorf("skill name is required")
	}
	if strings.TrimSpace(skill.Content) == "" {
		return fmt.Errorf("skill content is required")
	}
	if strings.TrimSpace(skill.Slug) == "" {
		return fmt.Errorf("skill slug is required")
	}
	return nil
}

func validateNoBuiltinConflict(skill storage.Skill, builtin *sessions.SkillLoader) error {
	if builtin == nil {
		return nil
	}
	if _, exists := builtin.Entries()[skill.Slug]; exists {
		return fmt.Errorf("skill slug conflicts with builtin skill")
	}
	return nil
}

func defaultConversationTitle(title string) string {
	if strings.TrimSpace(title) == "" {
		return "新对话"
	}
	return title
}

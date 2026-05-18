package runtime

import (
	"sort"
	"strings"

	"nano_cc/internal/web/storage"
)

func fallbackAssistantContent(content string) string {
	if strings.TrimSpace(content) == "" {
		return "(no response)"
	}
	return content
}

func inferConversationTitle(currentTitle, userMessage string) string {
	if strings.TrimSpace(currentTitle) != "" && currentTitle != "新对话" {
		return currentTitle
	}
	trimmed := strings.TrimSpace(userMessage)
	if len([]rune(trimmed)) > 30 {
		return string([]rune(trimmed)[:30])
	}
	if trimmed == "" {
		return "新对话"
	}
	return trimmed
}

func truncate(text string, max int) string {
	if len(text) <= max {
		return text
	}
	return text[:max]
}

func SkillNames(skills []storage.Skill) []string {
	names := make([]string, 0, len(skills))
	for _, skill := range skills {
		names = append(names, skill.Slug)
	}
	sort.Strings(names)
	return names
}

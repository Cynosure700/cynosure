package app

import (
	"sort"

	"nano_cc/internal/sessions"
	"nano_cc/internal/web/storage"
)

func appendBuiltinSkills(skills []storage.Skill, builtin *sessions.SkillLoader) []storage.Skill {
	entries := builtinSkillEntries(builtin)
	if len(entries) == 0 {
		return skills
	}
	merged := make([]storage.Skill, 0, len(entries)+len(skills))
	merged = append(merged, entries...)
	merged = append(merged, skills...)
	return merged
}

func builtinSkillByID(skillID string, builtin *sessions.SkillLoader) (storage.Skill, bool) {
	for _, skill := range builtinSkillEntries(builtin) {
		if skill.ID == skillID {
			return skill, true
		}
	}
	return storage.Skill{}, false
}

func builtinSkillEntries(builtin *sessions.SkillLoader) []storage.Skill {
	if builtin == nil {
		return nil
	}
	entries := builtin.Entries()
	if len(entries) == 0 {
		return nil
	}
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	skills := make([]storage.Skill, 0, len(names))
	for _, name := range names {
		entry := entries[name]
		description := ""
		if entry != nil {
			description = entry.Meta["description"]
		}
		content := ""
		if entry != nil {
			content = entry.Body
		}
		skills = append(skills, storage.Skill{
			ID:          builtinSkillID(name),
			Name:        name,
			Slug:        name,
			Description: description,
			Content:     content,
			Status:      "enabled",
			Source:      "builtin",
			ReadOnly:    true,
		})
	}
	return skills
}

func builtinSkillID(name string) string {
	return "builtin:" + name
}

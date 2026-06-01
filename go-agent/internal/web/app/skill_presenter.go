package app

import "nano_cc/internal/sessions"

func isBuiltinSkillID(skillID string, builtin *sessions.SkillLoader) bool {
	if builtin == nil || len(builtin.Entries()) == 0 {
		return false
	}
	const prefix = "builtin:"
	if len(skillID) <= len(prefix) || skillID[:len(prefix)] != prefix {
		return false
	}
	_, exists := builtin.Entries()[skillID[len(prefix):]]
	return exists
}

func builtinSkillID(name string) string {
	return "builtin:" + name
}

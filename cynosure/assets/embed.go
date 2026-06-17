package assets

import (
	"embed"
	"io/fs"
	"strings"
)

//go:embed system_prompt.md
var systemPrompt string

//go:embed all:skills
var skillsFS embed.FS

//go:embed prompts/*.md
var promptsFS embed.FS

// SystemPrompt 返回嵌入的基础 system prompt 内容。
func SystemPrompt() string {
	return systemPrompt
}

// FunctionalPrompt 返回嵌入的功能性提示词内容。
func FunctionalPrompt(name string) (string, error) {
	content, err := promptsFS.ReadFile("prompts/" + name + ".md")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

// BuiltinSkillsFS 返回嵌入的内置 skills 目录子树（根为 skills/）。
func BuiltinSkillsFS() fs.FS {
	sub, err := fs.Sub(skillsFS, "skills")
	if err != nil {
		// skills 目录由 go:embed 在编译期保证存在，理论上不会发生。
		panic(err)
	}
	return sub
}

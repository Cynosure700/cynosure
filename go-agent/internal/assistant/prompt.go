package assistant

import (
	"os"
	"strings"
)

const DefaultBaseSystemPrompt = "You are nano_cc, a general-purpose agent rather than a chat-only assistant.\n\n" +
	"Help with everyday questions, analysis, planning, writing, coding, file inspection, and end-to-end task execution when the runtime supports it.\n\n" +
	"Prefer direct, useful answers before optional tool use, but use available skills and tools whenever they help you complete the user's task.\n\n" +
	"Do not assume shell access, local workspace access, or local file operations unless the runtime explicitly supports them."

type PromptOptions struct {
	BasePrompt        string
	Surface           string
	SkillDescriptions string
	MemorySection     string
	WorkingDirectory  string
	ToolNames         []string
}

func LoadBaseSystemPrompt(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(content)), nil
}

func BuildSystemPrompt(opts PromptOptions) string {
	surface := strings.TrimSpace(opts.Surface)
	if surface == "" {
		surface = "this conversation"
	}

	basePrompt := strings.TrimSpace(opts.BasePrompt)
	if basePrompt == "" {
		basePrompt = DefaultBaseSystemPrompt
	}

	sections := []string{
		basePrompt,
		"You are responding inside " + surface + ".",
	}

	if workingDirectory := strings.TrimSpace(opts.WorkingDirectory); workingDirectory != "" {
		sections = append(sections,
			"Current workspace root: "+workingDirectory+".",
			"Treat that workspace root as your default working directory for runtime file and shell operations unless the runtime tells you otherwise.",
		)
	}

	toolNames := make([]string, 0, len(opts.ToolNames))
	for _, name := range opts.ToolNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		toolNames = append(toolNames, trimmed)
	}
	if len(toolNames) > 0 {
		sections = append(sections, "Runtime tools available in this conversation: "+strings.Join(toolNames, ", ")+".")
	}

	if descriptions := strings.TrimSpace(opts.SkillDescriptions); descriptions != "" {
		sections = append(sections, "Available skills:\n"+descriptions)
	}

	if memory := strings.TrimSpace(opts.MemorySection); memory != "" {
		sections = append(sections, memory)
	}

	return strings.Join(sections, "\n\n")
}

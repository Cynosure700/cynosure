package assistant

import "strings"

type PromptOptions struct {
	Surface           string
	SkillDescriptions string
	MemorySection     string
}

func BuildSystemPrompt(opts PromptOptions) string {
	surface := strings.TrimSpace(opts.Surface)
	if surface == "" {
		surface = "this conversation"
	}

	sections := []string{
		"You are nano_cc, a general-purpose agent assistant.",
		"Help with everyday questions, analysis, planning, writing, and coding when the user asks for it.",
		"Prefer direct, useful answers before optional tool use.",
		"Do not assume shell access, local workspace access, or local file operations unless the runtime explicitly supports them.",
		"You are responding inside " + surface + ".",
	}

	if descriptions := strings.TrimSpace(opts.SkillDescriptions); descriptions != "" {
		sections = append(sections, "Available skills:\n"+descriptions)
	}

	if memory := strings.TrimSpace(opts.MemorySection); memory != "" {
		sections = append(sections, memory)
	}

	return strings.Join(sections, "\n\n")
}

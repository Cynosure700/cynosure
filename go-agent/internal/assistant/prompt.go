package assistant

import (
	"os"
	"strings"
)

const sectionSeparator = "---"

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
		"# System Instructions",
		renderSection("Identity", basePrompt),
	}

	runtimeContext := []string{"Surface: " + surface}

	if workingDirectory := strings.TrimSpace(opts.WorkingDirectory); workingDirectory != "" {
		runtimeContext = append(runtimeContext,
			"Workspace root: "+workingDirectory,
			"Use the workspace root as the default working directory for runtime file and shell operations unless the runtime tells you otherwise.",
		)
	}
	sections = append(sections, renderSection("Runtime Context", strings.Join(runtimeContext, "\n\n")))

	toolNames := make([]string, 0, len(opts.ToolNames))
	for _, name := range opts.ToolNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		toolNames = append(toolNames, trimmed)
	}
	if len(toolNames) > 0 {
		sections = append(sections, renderSection("Runtime Tools", "The following tools are available in this conversation:\n\n"+renderList(toolNames)))
	}

	if descriptions := strings.TrimSpace(opts.SkillDescriptions); descriptions != "" {
		skillBody := strings.Join([]string{
			"The following skills are available as summaries only.",
			"Important rules:\n" + renderList([]string{
				"Each listed skill has a name and description.",
				"Before using or following a skill, call `load_skill` with the exact skill name to load its full instructions.",
				"Do not infer the full workflow from the summary alone.",
				"If multiple skills seem relevant, load the most specific matching skill first.",
			}),
			"Available skills:\n\n" + descriptions,
		}, "\n\n")
		sections = append(sections, renderSection("Skills", skillBody))
	}

	if memory := strings.TrimSpace(opts.MemorySection); memory != "" {
		sections = append(sections, renderSection("Memory", memory))
	}

	return strings.Join(sections, "\n\n"+sectionSeparator+"\n\n")
}

func renderSection(title, body string) string {
	title = strings.TrimSpace(title)
	body = strings.TrimSpace(body)
	if title == "" {
		return body
	}
	if body == "" {
		return "## " + title
	}
	return "## " + title + "\n\n" + body
}

func renderList(items []string) string {
	lines := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		lines = append(lines, "- "+trimmed)
	}
	return strings.Join(lines, "\n")
}

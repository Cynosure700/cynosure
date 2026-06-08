package assistant

import (
	"os"
	"strings"
)

const persistedOutputGuidance = "当较早的消息中出现 `<persisted-output ...>` 标记时，表示完整的工具输出已存入数据库，内联的只是预览。" +
	"如果预览不足以完成任务，请用标记中的 id 和偏移量调用 `read_persisted_output` 分块读取更多内容，不要猜测被省略的部分。" +
	"当看到 `[Earlier result compacted. Re-run if needed]` 时，请重新执行相关工具以再次获取该结果。"

const DefaultBaseSystemPrompt = "你是 nano_cc，一个运行在浏览器对话场景中的通用型智能体（general-purpose agent），而不是只能聊天的助手。\n\n" +
	"帮助用户处理日常问答、分析、规划、写作、编码、文件检查，以及在运行时支持时执行端到端的任务。优先给出直接、有用的回答；当工具或技能能帮助你完成任务时，主动使用它们。\n\n" +
	"不要假设自己拥有 Shell 访问、工作区访问或本地文件操作能力，除非运行时通过 <tools> 和 <workspace> 明确提供了这些能力。"

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

	sections := []string{renderTag("identity", basePrompt)}

	workspaceLines := []string{"Surface: " + surface}
	if workingDirectory := strings.TrimSpace(opts.WorkingDirectory); workingDirectory != "" {
		workspaceLines = append(workspaceLines,
			"Working directory: "+workingDirectory,
			"除非运行时另有说明，默认以工作目录作为运行时文件与 Shell 操作的根目录。",
		)
	}
	sections = append(sections, renderTag("workspace", strings.Join(workspaceLines, "\n")))

	toolNames := make([]string, 0, len(opts.ToolNames))
	for _, name := range opts.ToolNames {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		toolNames = append(toolNames, trimmed)
	}
	if len(toolNames) > 0 {
		toolBody := "本次会话可用的工具如下：\n\n" + renderList(toolNames)
		if toolNamesContain(toolNames, "read_persisted_output") {
			toolBody += "\n\n" + persistedOutputGuidance
		}
		sections = append(sections, renderTag("tools", toolBody))
	}

	if descriptions := strings.TrimSpace(opts.SkillDescriptions); descriptions != "" {
		skillBody := strings.Join([]string{
			"以下技能只提供摘要。",
			"重要规则：\n" + renderList([]string{
				"每个技能都有名称和描述。",
				"使用或遵循某个技能前，先用 `load_skill` 以精确的技能名加载其完整说明。",
				"不要仅凭摘要臆测完整的工作流。",
				"若多个技能看起来都相关，先加载最匹配、最具体的那个。",
			}),
			"可用技能：\n\n" + descriptions,
		}, "\n\n")
		sections = append(sections, renderTag("skills", skillBody))
	}

	if memory := strings.TrimSpace(opts.MemorySection); memory != "" {
		sections = append(sections, renderTag("memory", memory))
	}

	return strings.Join(sections, "\n\n")
}

func renderTag(tag, body string) string {
	tag = strings.TrimSpace(tag)
	body = strings.TrimSpace(body)
	if body == "" {
		return "<" + tag + ">\n</" + tag + ">"
	}
	return "<" + tag + ">\n" + body + "\n</" + tag + ">"
}

func toolNamesContain(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}
	return false
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

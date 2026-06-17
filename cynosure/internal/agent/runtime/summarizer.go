package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/runtime/compression"
	"nano_cc/internal/agent/storage"
)

const summarizerSystemPrompt = `你是上下文压缩引擎。请把目前为止的对话压缩成结构化的 Markdown 简报，使另一个 AI 助手在不丢失关键信息的前提下继续完成任务。

必须覆盖以下部分（用对应小标题，无内容则写"无"）：
1. 用户最新问题：用户当前最想解决的问题 / 目标。
2. 已完成的工作：做了哪些事、得到了哪些关键结论；改动过的文件务必保留完整路径。
3. 关键决策与原因：做了哪些技术选择，为什么。
4. 当前进展与待办：已完成到哪一步、还没做完的部分、明确的下一步 / TODO。
5. 重要报错、命令输出结论与约定：关键报错信息、命令输出的结论性信息、双方达成的约定。

规则：
- 不要调用工具，不要编造不存在的信息。
- 保留具体细节：文件路径、函数名、命令、报错、关键数值。
- 遇到 <persisted-output ...> 标记，保留其 id 并说明可用 read_persisted_output 取回完整内容。
- 工具结果若显示 "[Earlier result compacted. Re-run if needed]"，需注明该结果如有需要应重新执行。

只输出 Markdown 简报本身。`

// summarizeHistoryForContext runs a single non-tool LLM call to summarize the
// request-state history. It never mutates state.History.
func (s *Service) summarizeHistoryForContext(ctx context.Context, req compression.SummaryRequest) (compression.SummaryResult, error) {
	if s.LLM == nil {
		return compression.SummaryResult{}, fmt.Errorf("llm client is not configured")
	}
	messages := []openai.ChatCompletionMessage{
		{Role: "system", Content: summarizerSystemPrompt},
		{Role: "user", Content: renderHistoryForSummary(req.History)},
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    s.Cfg.LLM.ModelID,
		Messages: messages,
	})
	if err != nil {
		return compression.SummaryResult{}, err
	}
	if len(resp.Choices) == 0 {
		return compression.SummaryResult{}, fmt.Errorf("summary response had no choices")
	}
	summary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if summary == "" {
		return compression.SummaryResult{}, fmt.Errorf("summary response was empty")
	}
	return compression.SummaryResult{Summary: summary}, nil
}

func renderHistoryForSummary(history []storage.Message) string {
	var b strings.Builder
	b.WriteString("Summarize the following conversation history:\n\n")
	for _, msg := range history {
		role := msg.Role
		switch role {
		case "tool":
			b.WriteString("[tool result] ")
			b.WriteString(msg.Content)
		case "assistant":
			b.WriteString("[assistant] ")
			if strings.TrimSpace(msg.Content) != "" {
				b.WriteString(msg.Content)
			}
			for _, call := range msg.ToolCalls {
				fmt.Fprintf(&b, "\n[tool call] %s(%s)", call.Function.Name, call.Function.Arguments)
			}
		default:
			fmt.Fprintf(&b, "[%s] %s", role, msg.Content)
		}
		b.WriteString("\n\n")
	}
	return b.String()
}

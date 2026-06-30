package runtime

import (
	"fmt"
	openai "github.com/sashabaranov/go-openai"
	"strings"

	agenttools "github.com/Cynosure700/cynosure/cynosure/internal/tools"
)

const (
	todoWriteReminderThreshold = 10
	todoWriteReminderText      = "The todo_write tool hasn't been used recently. If you're working on tasks that work on tasks that could benefit from tracking progress, consider using the TodoWrite tool to track progress. Make sure that the clean up the task list when it's no relevant to the current task. Make sure that NEVER mention this reminder to the user"
)

func toolCallsInclude(toolCalls []openai.ToolCall, name string) bool {
	for _, call := range toolCalls {
		if call.Function.Name == name {
			return true
		}
	}
	return false
}

func todoWriteReminderMessage(todos []agenttools.TodoItem) openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{Role: "user", Content: renderTodoWriteReminder(todos)}
}

func renderTodoWriteReminder(todos []agenttools.TodoItem) string {
	return wrapSystemReminder(
		todoWriteReminderText,
		"Here are the existing contents of your todo list:\n\n"+renderTodoReminderList(todos),
	)
}

func renderTodoReminderList(todos []agenttools.TodoItem) string {
	if len(todos) == 0 {
		return "[]"
	}
	var b strings.Builder
	b.WriteString("[")
	for i, todo := range todos {
		if i > 0 {
			b.WriteByte('\n')
		}
		label := strings.TrimSpace(todo.ID)
		if label == "" {
			label = fmt.Sprintf("%d", i+1)
		}
		fmt.Fprintf(&b, "%s. [%s] %s", label, todo.Status, todo.Content)
	}
	b.WriteString("]")
	return b.String()
}

// todosAllCompleted 在存在待办且全部为 completed 时返回 true。收尾阶段（无 pending /
// in_progress 项）注入提醒只会诱导模型重复发出"全部已完成"的空操作 todo_write，因此跳过。
func todosAllCompleted(todos []agenttools.TodoItem) bool {
	if len(todos) == 0 {
		return false
	}
	for _, todo := range todos {
		if todo.Status != agenttools.TodoStatusCompleted {
			return false
		}
	}
	return true
}

func maybeAppendTodoWriteReminder(state *LoopState, tools *ToolRegistry, roundsSinceTodoWrite int) int {
	if roundsSinceTodoWrite < todoWriteReminderThreshold || tools == nil || !tools.isAllowed(agenttools.TodoWriteToolName) {
		return roundsSinceTodoWrite
	}
	if todosAllCompleted(state.Todos) {
		return roundsSinceTodoWrite
	}
	state.Messages = append(state.Messages, todoWriteReminderMessage(state.Todos))
	return 0
}

package runtime

import (
	openai "github.com/sashabaranov/go-openai"

	agenttools "cynosure/internal/tools"
)

const (
	todoWriteReminderThreshold = 5
	todoWriteReminderText      = "You have not called todo_write for 3 consecutive model rounds. This is a reminder — NOT an obligation to call it. Only call todo_write when you have actual task changes to record (new tasks, status updates, or plan changes). If no tasks need creating or updating, do NOT call todo_write — proceed normally."
)

func toolCallsInclude(toolCalls []openai.ToolCall, name string) bool {
	for _, call := range toolCalls {
		if call.Function.Name == name {
			return true
		}
	}
	return false
}

func todoWriteReminderMessage() openai.ChatCompletionMessage {
	return openai.ChatCompletionMessage{Role: "system", Content: todoWriteReminderText}
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
	state.Messages = append(state.Messages, todoWriteReminderMessage())
	return 0
}

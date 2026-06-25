package runtime

import (
	openai "github.com/sashabaranov/go-openai"

	agenttools "cynosure/internal/tools"
)

const (
	todoWriteReminderThreshold = 10
	todoWriteReminderText      = "<system-reminder> "+ "/n" + "The task tools haven't been used recently. If you're the type of person who likes to use task tracking to keep track of your progress, you can use todo_write to add new tasks and todo_write to update task status (set to in_progress when starting, completed when done). Please consider cleaning up the task list if it is not needed." + "/n" + "</system-reminder>"		
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
	return openai.ChatCompletionMessage{Role: "user", Content: todoWriteReminderText}
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

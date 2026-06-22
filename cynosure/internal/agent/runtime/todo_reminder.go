package runtime

import (
	openai "github.com/sashabaranov/go-openai"

	agenttools "cynosure/internal/tools"
)

const (
	todoWriteReminderThreshold = 3
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

func maybeAppendTodoWriteReminder(state *LoopState, tools *ToolRegistry, roundsSinceTodoWrite int) int {
	if roundsSinceTodoWrite < todoWriteReminderThreshold || tools == nil || !tools.isAllowed(agenttools.TodoWriteToolName) {
		return roundsSinceTodoWrite
	}
	state.Messages = append(state.Messages, todoWriteReminderMessage())
	return 0
}

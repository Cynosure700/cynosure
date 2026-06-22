package runtime

import (
	openai "github.com/sashabaranov/go-openai"

	agenttools "cynosure/internal/tools"
)

const (
	todoWriteReminderThreshold = 3
	todoWriteReminderText      = "<system-reminder>\nYou have not called todo_write for 3 consecutive model rounds. If the task is multi-step or your plan has changed, call todo_write to create or update the current task plan before continuing. If todo_write is unnecessary for this simple step, continue normally.\n</system-reminder>"
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

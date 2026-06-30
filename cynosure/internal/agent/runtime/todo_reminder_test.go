package runtime

import (
	"strings"
	"testing"

	"github.com/Cynosure700/cynosure/cynosure/internal/config"
	agenttools "github.com/Cynosure700/cynosure/cynosure/internal/tools"
)

func newTodoWriteRegistry(t *testing.T) *ToolRegistry {
	t.Helper()
	tools := NewToolRegistry(config.AppConfig{AllowedTools: []string{agenttools.TodoWriteToolName}})
	if !tools.isAllowed(agenttools.TodoWriteToolName) {
		t.Fatalf("expected %s to be allowed", agenttools.TodoWriteToolName)
	}
	return tools
}

func reminderInjected(state *LoopState) bool {
	for _, msg := range state.Messages {
		if strings.Contains(msg.Content, todoWriteReminderText) {
			return true
		}
	}
	return false
}

func TestMaybeAppendTodoWriteReminder_BelowThreshold(t *testing.T) {
	tools := newTodoWriteRegistry(t)
	state := &LoopState{}
	got := maybeAppendTodoWriteReminder(state, tools, todoWriteReminderThreshold-1)
	if got != todoWriteReminderThreshold-1 {
		t.Fatalf("expected counter unchanged, got %d", got)
	}
	if reminderInjected(state) {
		t.Fatalf("expected no reminder below threshold")
	}
}

func TestMaybeAppendTodoWriteReminder_InjectsWhenPendingExists(t *testing.T) {
	tools := newTodoWriteRegistry(t)
	state := &LoopState{Todos: []agenttools.TodoItem{
		{ID: "1", Content: "a", Status: agenttools.TodoStatusCompleted},
		{ID: "2", Content: "b", Status: agenttools.TodoStatusInProgress},
	}}
	got := maybeAppendTodoWriteReminder(state, tools, todoWriteReminderThreshold)
	if got != 0 {
		t.Fatalf("expected counter reset to 0, got %d", got)
	}
	if !reminderInjected(state) {
		t.Fatalf("expected reminder injected when work remains")
	}
	reminder := state.Messages[len(state.Messages)-1]
	expectedTodos := "[1. [completed] a\n2. [in_progress] b]"
	for _, want := range []string{
		"<system-reminder>",
		"The todo_write tool hasn't been used recently.",
		"consider using the TodoWrite tool to track progress.",
		"Make sure that the clean up the task list when it's no relevant to the current task.",
		"Make sure that NEVER mention this reminder to the user",
		"Here are the existing contents of your todo list:",
		expectedTodos,
		"</system-reminder>",
	} {
		if !strings.Contains(reminder.Content, want) {
			t.Fatalf("expected reminder to contain %q, got %q", want, reminder.Content)
		}
	}
	if !strings.HasPrefix(reminder.Content, "<system-reminder>\n") || !strings.HasSuffix(reminder.Content, "\n</system-reminder>") {
		t.Fatalf("expected reminder to be wrapped in system-reminder tags, got %q", reminder.Content)
	}
}

func TestMaybeAppendTodoWriteReminder_SkipsWhenAllCompleted(t *testing.T) {
	tools := newTodoWriteRegistry(t)
	state := &LoopState{Todos: []agenttools.TodoItem{
		{ID: "1", Content: "a", Status: agenttools.TodoStatusCompleted},
		{ID: "2", Content: "b", Status: agenttools.TodoStatusCompleted},
	}}
	got := maybeAppendTodoWriteReminder(state, tools, todoWriteReminderThreshold)
	if got != todoWriteReminderThreshold {
		t.Fatalf("expected counter preserved when skipping, got %d", got)
	}
	if reminderInjected(state) {
		t.Fatalf("expected no reminder when all todos completed")
	}
}

func TestMaybeAppendTodoWriteReminder_InjectsWhenNoTodos(t *testing.T) {
	tools := newTodoWriteRegistry(t)
	state := &LoopState{}
	got := maybeAppendTodoWriteReminder(state, tools, todoWriteReminderThreshold)
	if got != 0 {
		t.Fatalf("expected counter reset to 0, got %d", got)
	}
	if !reminderInjected(state) {
		t.Fatalf("expected reminder injected when no todos exist yet")
	}
}

func TestTodosAllCompleted(t *testing.T) {
	cases := []struct {
		name  string
		todos []agenttools.TodoItem
		want  bool
	}{
		{"empty", nil, false},
		{"all completed", []agenttools.TodoItem{{Status: agenttools.TodoStatusCompleted}}, true},
		{"has pending", []agenttools.TodoItem{{Status: agenttools.TodoStatusCompleted}, {Status: agenttools.TodoStatusPending}}, false},
		{"has in_progress", []agenttools.TodoItem{{Status: agenttools.TodoStatusInProgress}}, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := todosAllCompleted(tc.todos); got != tc.want {
				t.Fatalf("todosAllCompleted(%v) = %v, want %v", tc.todos, got, tc.want)
			}
		})
	}
}

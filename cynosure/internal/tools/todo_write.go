package tools

import (
	"context"
	"fmt"
	"strings"
)

const (
	TodoStatusPending    = "pending"
	TodoStatusInProgress = "in_progress"
	TodoStatusCompleted  = "completed"
)

type TodoItem struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	ActiveForm string `json:"activeForm"`
	Status     string `json:"status"`
}

type TodoWriteResult struct {
	Output string
	Todos  []TodoItem
}

func ExecuteTodoWrite(ctx context.Context, args map[string]any) (TodoWriteResult, error) {
	rawTodos, ok := args["todos"].([]any)
	if !ok {
		return TodoWriteResult{}, fmt.Errorf("todos is required")
	}

	todos := make([]TodoItem, 0, len(rawTodos))
	counts := map[string]int{TodoStatusPending: 0, TodoStatusInProgress: 0, TodoStatusCompleted: 0}
	for i, raw := range rawTodos {
		item, ok := raw.(map[string]any)
		if !ok {
			return TodoWriteResult{}, fmt.Errorf("todos[%d] must be an object", i)
		}
		id, _ := item["id"].(string)
		content, _ := item["content"].(string)
		activeForm, _ := item["activeForm"].(string)
		status, _ := item["status"].(string)
		id = strings.TrimSpace(id)
		content = strings.TrimSpace(content)
		activeForm = strings.TrimSpace(activeForm)
		status = strings.TrimSpace(status)
		if id == "" {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].id is required", i)
		}
		if content == "" {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].content is required", i)
		}
		if activeForm == "" {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].activeForm is required", i)
		}
		if !isValidTodoStatus(status) {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].status must be one of pending, in_progress, completed", i)
		}
		todos = append(todos, TodoItem{ID: id, Content: content, ActiveForm: activeForm, Status: status})
		counts[status]++
	}

	return TodoWriteResult{
		Output: fmt.Sprintf("Todo list updated: %d items (pending: %d, in_progress: %d, completed: %d).", len(todos), counts[TodoStatusPending], counts[TodoStatusInProgress], counts[TodoStatusCompleted]),
		Todos:  todos,
	}, nil
}

func handleTodoList(ctx context.Context, args map[string]any) (string, error) {
	if len(args) > 0 {
		return "", fmt.Errorf("todo_list does not accept arguments")
	}
	todos, _ := TodoSnapshotFromContext(ctx)
	if len(todos) == 0 {
		return "Todo list is empty. Use todo_write to create or update the current task plan.", nil
	}
	counts := map[string]int{TodoStatusPending: 0, TodoStatusInProgress: 0, TodoStatusCompleted: 0}
	var b strings.Builder
	for _, todo := range todos {
		counts[todo.Status]++
	}
	fmt.Fprintf(&b, "Todo list: %d items (pending: %d, in_progress: %d, completed: %d).", len(todos), counts[TodoStatusPending], counts[TodoStatusInProgress], counts[TodoStatusCompleted])
	for _, todo := range todos {
		fmt.Fprintf(&b, "\n[%s] %s: %s", todo.Status, todo.ID, todo.Content)
	}
	return b.String(), nil
}

func isValidTodoStatus(status string) bool {
	switch status {
	case TodoStatusPending, TodoStatusInProgress, TodoStatusCompleted:
		return true
	default:
		return false
	}
}

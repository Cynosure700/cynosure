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
	ID      string `json:"id"`
	Content string `json:"content"`
	Status  string `json:"status"`
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
		status, _ := item["status"].(string)
		id = strings.TrimSpace(id)
		content = strings.TrimSpace(content)
		status = strings.TrimSpace(status)
		if id == "" {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].id is required", i)
		}
		if content == "" {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].content is required", i)
		}
		if !isValidTodoStatus(status) {
			return TodoWriteResult{}, fmt.Errorf("todos[%d].status must be one of pending, in_progress, completed", i)
		}
		todos = append(todos, TodoItem{ID: id, Content: content, Status: status})
		counts[status]++
	}

	return TodoWriteResult{
		Output: fmt.Sprintf("Todo list updated: %d items (pending: %d, in_progress: %d, completed: %d).", len(todos), counts[TodoStatusPending], counts[TodoStatusInProgress], counts[TodoStatusCompleted]),
		Todos:  todos,
	}, nil
}

func isValidTodoStatus(status string) bool {
	switch status {
	case TodoStatusPending, TodoStatusInProgress, TodoStatusCompleted:
		return true
	default:
		return false
	}
}

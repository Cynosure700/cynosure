package tools

import (
	"fmt"
	"strings"
	"sync"
)

type TodoItem struct {
	ID     string `json:"id"`
	Text   string `json:"text"`
	Status string `json:"status"`
}

type TodoManager struct {
	mu    sync.Mutex
	Items []TodoItem
}

var Todo = &TodoManager{}

const maxTodos = 20

func (tm *TodoManager) Update(items []map[string]any) (string, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if len(items) > maxTodos {
		return "", fmt.Errorf("max %d todos allowed", maxTodos)
	}

	validated := make([]TodoItem, 0, len(items))
	inProgressCount := 0

	for i, item := range items {
		id, _ := item["id"].(string)
		if id == "" {
			id = fmt.Sprintf("%d", i+1)
		}

		text, _ := item["text"].(string)
		text = strings.TrimSpace(text)
		if text == "" {
			return "", fmt.Errorf("item %d: text is required", i+1)
		}

		status, _ := item["status"].(string)
		if status == "" {
			status = "pending"
		}
		switch status {
		case "pending", "in_progress", "completed":
		default:
			return "", fmt.Errorf("item %d: invalid status %q", i+1, status)
		}

		if status == "in_progress" {
			inProgressCount++
		}

		validated = append(validated, TodoItem{ID: id, Text: text, Status: status})
	}

	if inProgressCount > 1 {
		return "", fmt.Errorf("only one task can be in_progress at a time")
	}

	tm.Items = validated
	return tm.Render(), nil
}

func (tm *TodoManager) Render() string {
	if len(tm.Items) == 0 {
		return "(no tasks)"
	}

	markers := map[string]string{
		"pending":     "[ ]",
		"in_progress": "[>]",
		"completed":   "[x]",
	}

	var lines []string
	completed := 0
	for _, item := range tm.Items {
		marker := markers[item.Status]
		lines = append(lines, fmt.Sprintf("%s #%s: %s", marker, item.ID, item.Text))
		if item.Status == "completed" {
			completed++
		}
	}
	lines = append(lines, fmt.Sprintf("\n(%d/%d completed)", completed, len(tm.Items)))

	return strings.Join(lines, "\n")
}

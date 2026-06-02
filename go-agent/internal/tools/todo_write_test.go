package tools

import (
	"context"
	"strings"
	"testing"
)

func TestHandleTodoWrite_ReturnsSummaryAndParsedTodos(t *testing.T) {
	result, err := handleTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "content": "梳理需求", "status": "completed"},
		map[string]any{"id": "2", "content": "实现功能", "status": "in_progress"},
		map[string]any{"id": "3", "content": "运行测试", "status": "pending"},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "Todo list updated: 3 items (pending: 1, in_progress: 1, completed: 1)." {
		t.Fatalf("unexpected summary: %q", result.Output)
	}
	expected := []TodoItem{
		{ID: "1", Content: "梳理需求", Status: "completed"},
		{ID: "2", Content: "实现功能", Status: "in_progress"},
		{ID: "3", Content: "运行测试", Status: "pending"},
	}
	if len(result.Todos) != len(expected) {
		t.Fatalf("expected %d todos, got %#v", len(expected), result.Todos)
	}
	for i := range expected {
		if result.Todos[i] != expected[i] {
			t.Fatalf("todo[%d] expected %#v, got %#v", i, expected[i], result.Todos[i])
		}
	}
}

func TestHandleTodoWrite_AllowsEmptyTodoList(t *testing.T) {
	result, err := handleTodoWrite(context.Background(), map[string]any{"todos": []any{}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "Todo list updated: 0 items (pending: 0, in_progress: 0, completed: 0)." {
		t.Fatalf("unexpected summary: %q", result.Output)
	}
	if len(result.Todos) != 0 {
		t.Fatalf("expected empty todos, got %#v", result.Todos)
	}
}

func TestHandleTodoWrite_RejectsInvalidStatus(t *testing.T) {
	_, err := handleTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "content": "实现功能", "status": "blocked"},
	}})
	if err == nil {
		t.Fatalf("expected invalid status error")
	}
	if !strings.Contains(err.Error(), "todos[0].status must be one of pending, in_progress, completed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleTodoWrite_RejectsMissingContent(t *testing.T) {
	_, err := handleTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "status": "pending"},
	}})
	if err == nil {
		t.Fatalf("expected missing content error")
	}
	if !strings.Contains(err.Error(), "todos[0].content is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

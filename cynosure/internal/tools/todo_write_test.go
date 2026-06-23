package tools

import (
	"context"
	"strings"
	"testing"
)

func TestHandleTodoWrite_ReturnsSummaryAndParsedTodos(t *testing.T) {
	result, err := ExecuteTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "content": "梳理需求", "activeForm": "梳理需求中", "status": "completed"},
		map[string]any{"id": "2", "content": "实现功能", "activeForm": "实现功能中", "status": "in_progress"},
		map[string]any{"id": "3", "content": "运行测试", "activeForm": "运行测试中", "status": "pending"},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Output != "Todo list updated: 3 items (pending: 1, in_progress: 1, completed: 1)." {
		t.Fatalf("unexpected summary: %q", result.Output)
	}
	expected := []TodoItem{
		{ID: "1", Content: "梳理需求", ActiveForm: "梳理需求中", Status: "completed"},
		{ID: "2", Content: "实现功能", ActiveForm: "实现功能中", Status: "in_progress"},
		{ID: "3", Content: "运行测试", ActiveForm: "运行测试中", Status: "pending"},
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
	result, err := ExecuteTodoWrite(context.Background(), map[string]any{"todos": []any{}})
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
	_, err := ExecuteTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "content": "实现功能", "activeForm": "实现功能中", "status": "blocked"},
	}})
	if err == nil {
		t.Fatalf("expected invalid status error")
	}
	if !strings.Contains(err.Error(), "todos[0].status must be one of pending, in_progress, completed") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleTodoWrite_RejectsMissingContent(t *testing.T) {
	_, err := ExecuteTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "activeForm": "实现功能中", "status": "pending"},
	}})
	if err == nil {
		t.Fatalf("expected missing content error")
	}
	if !strings.Contains(err.Error(), "todos[0].content is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleTodoWrite_RejectsMissingActiveForm(t *testing.T) {
	_, err := ExecuteTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "content": "实现功能", "status": "pending"},
	}})
	if err == nil {
		t.Fatalf("expected missing activeForm error")
	}
	if !strings.Contains(err.Error(), "todos[0].activeForm is required") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleTodoList_ReturnsCurrentTodosFromContext(t *testing.T) {
	ctx := WithTodoSnapshot(context.Background(), []TodoItem{
		{ID: "1", Content: "梳理需求", Status: "completed"},
		{ID: "2", Content: "实现功能", Status: "in_progress"},
		{ID: "3", Content: "运行测试", Status: "pending"},
	})

	output, err := handleTodoList(ctx, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, want := range []string{
		"Todo list: 3 items (pending: 1, in_progress: 1, completed: 1).",
		"[completed] 1: 梳理需求",
		"[in_progress] 2: 实现功能",
		"[pending] 3: 运行测试",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q, got %q", want, output)
		}
	}
}

func TestHandleTodoList_ReturnsEmptyMessageWithoutTodos(t *testing.T) {
	output, err := handleTodoList(context.Background(), map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if output != "Todo list is empty. Use todo_write to create or update the current task plan." {
		t.Fatalf("unexpected output: %q", output)
	}
}

func TestHandleTodoList_RejectsArguments(t *testing.T) {
	ctx := WithTodoSnapshot(context.Background(), []TodoItem{{ID: "1", Content: "x", Status: "pending"}})
	_, err := handleTodoList(ctx, map[string]any{"unused": "value"})
	if err == nil {
		t.Fatalf("expected unexpected argument error")
	}
	if !strings.Contains(err.Error(), "todo_list does not accept arguments") {
		t.Fatalf("unexpected error: %v", err)
	}
}

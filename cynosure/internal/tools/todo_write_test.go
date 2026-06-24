package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

const wantTodoStatusPrefix = "当前任务状态信息为:"

func TestHandleTodoWrite_ReturnsJSONTodosAndParsedTodos(t *testing.T) {
	result, err := ExecuteTodoWrite(context.Background(), map[string]any{"todos": []any{
		map[string]any{"id": "1", "content": "梳理需求", "activeForm": "梳理需求中", "status": "completed"},
		map[string]any{"id": "2", "content": "实现功能", "activeForm": "实现功能中", "status": "in_progress"},
		map[string]any{"id": "3", "content": "运行测试", "activeForm": "运行测试中", "status": "pending"},
	}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []TodoItem{
		{ID: "1", Content: "梳理需求", ActiveForm: "梳理需求中", Status: "completed"},
		{ID: "2", Content: "实现功能", ActiveForm: "实现功能中", Status: "in_progress"},
		{ID: "3", Content: "运行测试", ActiveForm: "运行测试中", Status: "pending"},
	}
	payload := parseTodoStatusOutput(t, result.Output)
	if len(payload.Todos) != len(expected) {
		t.Fatalf("expected output todos %#v, got %#v", expected, payload.Todos)
	}
	for i := range expected {
		if payload.Todos[i] != expected[i] {
			t.Fatalf("output todo[%d] expected %#v, got %#v", i, expected[i], payload.Todos[i])
		}
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
	payload := parseTodoStatusOutput(t, result.Output)
	if payload.Todos == nil || len(payload.Todos) != 0 {
		t.Fatalf("expected empty output todos, got %#v", payload.Todos)
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
		{ID: "1", Content: "梳理需求", ActiveForm: "梳理需求中", Status: "completed"},
		{ID: "2", Content: "实现功能", ActiveForm: "实现功能中", Status: "in_progress"},
		{ID: "3", Content: "运行测试", ActiveForm: "运行测试中", Status: "pending"},
	})

	output, err := handleTodoList(ctx, map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload := parseTodoStatusOutput(t, output)
	expected := []TodoItem{
		{ID: "1", Content: "梳理需求", ActiveForm: "梳理需求中", Status: "completed"},
		{ID: "2", Content: "实现功能", ActiveForm: "实现功能中", Status: "in_progress"},
		{ID: "3", Content: "运行测试", ActiveForm: "运行测试中", Status: "pending"},
	}
	if len(payload.Todos) != len(expected) {
		t.Fatalf("expected output todos %#v, got %#v", expected, payload.Todos)
	}
	for i := range expected {
		if payload.Todos[i] != expected[i] {
			t.Fatalf("output todo[%d] expected %#v, got %#v", i, expected[i], payload.Todos[i])
		}
	}
}

func TestHandleTodoList_ReturnsEmptyMessageWithoutTodos(t *testing.T) {
	output, err := handleTodoList(context.Background(), map[string]any{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	payload := parseTodoStatusOutput(t, output)
	if payload.Todos == nil || len(payload.Todos) != 0 {
		t.Fatalf("expected empty output todos, got %#v", payload.Todos)
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

func parseTodoStatusOutput(t *testing.T, output string) struct {
	Todos []TodoItem `json:"todos"`
} {
	t.Helper()
	if !strings.HasPrefix(output, wantTodoStatusPrefix) {
		t.Fatalf("expected output prefix %q, got %q", wantTodoStatusPrefix, output)
	}
	jsonPart := strings.TrimSpace(strings.TrimPrefix(output, wantTodoStatusPrefix))
	var payload struct {
		Todos []TodoItem `json:"todos"`
	}
	if err := json.Unmarshal([]byte(jsonPart), &payload); err != nil {
		t.Fatalf("expected output JSON after prefix, got %q: %v", output, err)
	}
	return payload
}

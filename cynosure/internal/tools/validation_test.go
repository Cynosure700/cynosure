package tools

import (
	"encoding/json"
	"strings"
	"testing"
)

func schemaOf(t *testing.T, name string) json.RawMessage {
	t.Helper()
	for _, def := range AllToolDefs {
		if def.Function != nil && def.Function.Name == name {
			return RawSchemaFromParameters(def.Function.Parameters)
		}
	}
	t.Fatalf("tool %q not found in AllToolDefs", name)
	return nil
}

func TestValidateToolArgs_MissingRequired(t *testing.T) {
	err := ValidateToolArgs("read_file", schemaOf(t, "read_file"), map[string]any{})
	if err == nil {
		t.Fatal("expected error for missing required parameter path")
	}
	if !strings.Contains(err.Error(), "read_file") {
		t.Fatalf("error should mention tool name, got: %v", err)
	}
}

func TestValidateToolArgs_TypeMismatch(t *testing.T) {
	err := ValidateToolArgs("read_file", schemaOf(t, "read_file"), map[string]any{
		"path":  "main.go",
		"limit": "ten",
	})
	if err == nil {
		t.Fatal("expected error for limit type mismatch")
	}
}

func TestValidateToolArgs_IntegerAcceptsWholeFloat(t *testing.T) {
	err := ValidateToolArgs("read_file", schemaOf(t, "read_file"), map[string]any{
		"path":  "main.go",
		"limit": float64(20),
	})
	if err != nil {
		t.Fatalf("whole-number float should pass integer validation, got: %v", err)
	}
}

func TestValidateToolArgs_EnumViolation(t *testing.T) {
	err := ValidateToolArgs("todo_write", schemaOf(t, "todo_write"), map[string]any{
		"todos": []any{
			map[string]any{"id": "1", "content": "do it", "status": "not_a_status"},
		},
	})
	if err == nil {
		t.Fatal("expected error for invalid enum value in status")
	}
}

func TestValidateToolArgs_NestedRequired(t *testing.T) {
	err := ValidateToolArgs("todo_write", schemaOf(t, "todo_write"), map[string]any{
		"todos": []any{
			map[string]any{"id": "1", "content": "missing status"},
		},
	})
	if err == nil {
		t.Fatal("expected error for missing nested required field status")
	}
}

func TestValidateToolArgs_Valid(t *testing.T) {
	err := ValidateToolArgs("todo_write", schemaOf(t, "todo_write"), map[string]any{
		"todos": []any{
			map[string]any{"id": "1", "content": "do it", "status": TodoStatusPending},
		},
	})
	if err != nil {
		t.Fatalf("valid args should pass, got: %v", err)
	}
}

func TestValidateToolArgs_TodoListRejectsUnexpectedArgument(t *testing.T) {
	err := ValidateToolArgs("todo_list", schemaOf(t, "todo_list"), map[string]any{"todos": []any{}})
	if err == nil {
		t.Fatal("expected todo_list to reject unexpected arguments")
	}
}

func TestValidateToolArgs_EmptySchemaPassThrough(t *testing.T) {
	if err := ValidateToolArgs("x", nil, map[string]any{"a": 1}); err != nil {
		t.Fatalf("nil schema should pass through, got: %v", err)
	}
	if err := ValidateToolArgs("x", json.RawMessage("null"), map[string]any{"a": 1}); err != nil {
		t.Fatalf("null schema should pass through, got: %v", err)
	}
}

func TestValidateToolArgs_MalformedSchemaPassThrough(t *testing.T) {
	if err := ValidateToolArgs("x", json.RawMessage("{not json"), map[string]any{"a": 1}); err != nil {
		t.Fatalf("malformed schema should pass through, got: %v", err)
	}
}

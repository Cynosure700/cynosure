package tools

import (
	"strings"
	"testing"
)

func TestLoadSkillToolDescriptionRequiresExactNameBeforeUse(t *testing.T) {
	for _, tool := range AllToolDefs {
		if tool.Function == nil || tool.Function.Name != "load_skill" {
			continue
		}
		description := tool.Function.Description
		for _, want := range []string{"full instructions", "exact name", "before using or following"} {
			if !strings.Contains(description, want) {
				t.Fatalf("expected load_skill description to contain %q, got %q", want, description)
			}
		}
		if strings.Contains(strings.ToLower(description), "database") {
			t.Fatalf("expected load_skill description to avoid database wording, got %q", description)
		}
		return
	}
	t.Fatalf("expected AllToolDefs to include load_skill")
}

func TestTodoListToolDefinitionIsNoArgumentReadOnlyQuery(t *testing.T) {
	for _, tool := range AllToolDefs {
		if tool.Function == nil || tool.Function.Name != "todo_list" {
			continue
		}
		if !strings.Contains(tool.Function.Description, "current task plan") {
			t.Fatalf("expected todo_list description to mention current task plan, got %q", tool.Function.Description)
		}
		schema := string(RawSchemaFromParameters(tool.Function.Parameters))
		for _, want := range []string{`"type":"object"`, `"properties":{}`, `"additionalProperties":false`} {
			if !strings.Contains(schema, want) {
				t.Fatalf("expected todo_list schema to contain %s, got %s", want, schema)
			}
		}
		return
	}
	t.Fatalf("expected AllToolDefs to include todo_list")
}

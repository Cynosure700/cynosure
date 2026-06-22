package tools

import (
	"strings"
	"testing"
)

func TestAllToolSpecsExposeDefaultMaxResultSizeChars(t *testing.T) {
	if DefaultMaxResultSizeChars != 50000 {
		t.Fatalf("DefaultMaxResultSizeChars = %d, want 50000", DefaultMaxResultSizeChars)
	}
	if got := MaxResultSizeCharsForTool("bash"); got != 50000 {
		t.Fatalf("bash max result chars = %d, want 50000", got)
	}
	if got := MaxResultSizeCharsForTool("mcp__unknown__tool"); got != 50000 {
		t.Fatalf("unknown max result chars = %d, want 50000", got)
	}
}

func TestAllToolSpecsMatchAllToolDefs(t *testing.T) {
	specNames := map[string]struct{}{}
	for _, spec := range AllToolSpecs {
		if spec.Definition.Function == nil {
			t.Fatalf("tool spec has nil function: %#v", spec.Definition)
		}
		name := spec.Definition.Function.Name
		if name == "" {
			t.Fatalf("tool spec has empty name: %#v", spec.Definition)
		}
		if spec.MaxResultSizeChars != 50000 {
			t.Fatalf("%s max result chars = %d, want 50000", name, spec.MaxResultSizeChars)
		}
		specNames[name] = struct{}{}
	}
	defNames := map[string]struct{}{}
	for _, def := range AllToolDefs {
		if def.Function == nil {
			t.Fatalf("tool def has nil function: %#v", def)
		}
		defNames[def.Function.Name] = struct{}{}
	}
	if len(specNames) != len(defNames) {
		t.Fatalf("spec count = %d, def count = %d", len(specNames), len(defNames))
	}
	for name := range defNames {
		if _, ok := specNames[name]; !ok {
			t.Fatalf("AllToolDefs contains %s but AllToolSpecs does not", name)
		}
	}
}

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

func TestReadFileToolDescriptionRestrictsToExistingFiles(t *testing.T) {
	for _, tool := range AllToolDefs {
		if tool.Function == nil || tool.Function.Name != "read_file" {
			continue
		}
		description := tool.Function.Description
		for _, want := range []string{"existing regular file", "not for directories", "not for speculative"} {
			if !strings.Contains(description, want) {
				t.Fatalf("expected read_file description to contain %q, got %q", want, description)
			}
		}
		schema := string(RawSchemaFromParameters(tool.Function.Parameters))
		for _, want := range []string{"confirmed existing regular file", "Do not pass a directory", "do not use read_file to check whether a path exists"} {
			if !strings.Contains(schema, want) {
				t.Fatalf("expected read_file schema to contain %q, got %q", want, schema)
			}
		}
		return
	}
	t.Fatalf("expected AllToolDefs to include read_file")
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

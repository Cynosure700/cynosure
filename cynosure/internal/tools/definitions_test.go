package tools

import (
	"encoding/json"
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

func TestFileToolDescriptionsGuideReadBeforeWriteAndExactEdits(t *testing.T) {
	toolsByName := map[string]string{}
	schemasByName := map[string]string{}
	for _, tool := range AllToolDefs {
		if tool.Function == nil {
			continue
		}
		toolsByName[tool.Function.Name] = tool.Function.Description
		schemasByName[tool.Function.Name] = string(RawSchemaFromParameters(tool.Function.Parameters))
	}
	for name, wants := range map[string][]string{
		"read_file": {
			"Reads a file from the local filesystem",
			"Assume this tool is able to read all files on the machine",
			"If the user provides a path to a file, assume that path is valid",
			"It is okay to read a user-provided file path that does not exist",
			"For paths inferred by you rather than provided by the user, confirm the file exists before reading",
			"offset",
			"1-based starting line",
			"line number + tab",
		},
		"write_file": {
			"Writes a file to the local filesystem",
			"overwrite the existing file",
			"MUST use read_file first",
			"Prefer edit_file for modifying existing files",
			"NEVER create documentation files",
			"Only use emojis if the user explicitly requests it",
		},
		"edit_file": {
			"Performs exact string replacements in files",
			"MUST use read_file first",
			"You must use your `Read` tool at least once in the conversation before editing.",
			"preserve the exact indentation",
			"line number + tab",
			"old_text must be unique",
			"Prefer multi_edit",
		},
		"multi_edit": {
			"Performs multiple exact string replacements in one or more files",
			"MUST use read_file first",
			"You must use your `Read` tool at least once in the conversation before editing.",
			"Use file_path + edits for a single file, or files for multiple files",
			"preserve the exact indentation",
			"line number + tab",
			"old_string must be unique",
			"applied sequentially and atomically per file",
		},
	} {
		description, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected AllToolDefs to include %s", name)
		}
		schema := schemasByName[name]
		combined := description + "\n" + schema
		for _, want := range wants {
			if !strings.Contains(combined, want) {
				t.Fatalf("expected %s guidance to contain %q, got description=%q schema=%q", name, want, description, schema)
			}
		}
	}
	if strings.Contains(toolsByName["read_file"], "Only read confirmed existing files") ||
		strings.Contains(schemasByName["read_file"], "do not read user-provided missing files") {
		t.Fatalf("read_file guidance should allow direct reads of user-provided paths and missing-file errors, got description=%q schema=%q", toolsByName["read_file"], schemasByName["read_file"])
	}
}

func TestMultiEditSchemaSupportsSingleAndMultipleFiles(t *testing.T) {
	schema := schemaMapForTool(t, "multi_edit")
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("multi_edit schema properties missing or wrong type: %#v", schema["properties"])
	}
	for _, name := range []string{"file_path", "edits", "files"} {
		if _, ok := properties[name]; !ok {
			t.Fatalf("multi_edit schema should expose %s, got properties %#v", name, properties)
		}
	}
	files, ok := properties["files"].(map[string]any)
	if !ok {
		t.Fatalf("multi_edit files schema missing or wrong type: %#v", properties["files"])
	}
	items, ok := files["items"].(map[string]any)
	if !ok {
		t.Fatalf("multi_edit files items missing or wrong type: %#v", files["items"])
	}
	fileProperties, ok := items["properties"].(map[string]any)
	if !ok {
		t.Fatalf("multi_edit files item properties missing or wrong type: %#v", items["properties"])
	}
	for _, name := range []string{"file_path", "edits"} {
		if _, ok := fileProperties[name]; !ok {
			t.Fatalf("multi_edit files item should expose %s, got properties %#v", name, fileProperties)
		}
	}
	if required, ok := schema["required"].([]any); ok && len(required) > 0 {
		t.Fatalf("multi_edit schema should not require single-file fields when files is used, got %#v", required)
	}
}

func TestFileToolSchemasRequireFilePathNotPath(t *testing.T) {
	for _, name := range []string{"read_file", "write_file", "edit_file"} {
		schema := schemaMapForTool(t, name)
		properties, ok := schema["properties"].(map[string]any)
		if !ok {
			t.Fatalf("%s schema properties missing or wrong type: %#v", name, schema["properties"])
		}
		if _, ok := properties["file_path"]; !ok {
			t.Fatalf("%s schema should expose file_path, got properties %#v", name, properties)
		}
		if _, ok := properties["path"]; ok {
			t.Fatalf("%s schema should not expose path, got properties %#v", name, properties)
		}

		required, ok := schema["required"].([]any)
		if !ok {
			t.Fatalf("%s schema required missing or wrong type: %#v", name, schema["required"])
		}
		if !containsAnyString(required, "file_path") {
			t.Fatalf("%s schema should require file_path, got %#v", name, required)
		}
		if containsAnyString(required, "path") {
			t.Fatalf("%s schema should not require path, got %#v", name, required)
		}
	}
}

func TestReadFileSchemaExposesOffsetAndLimitAsOptionalLineControls(t *testing.T) {
	schema := schemaMapForTool(t, "read_file")
	properties, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Fatalf("read_file schema properties missing or wrong type: %#v", schema["properties"])
	}
	for _, name := range []string{"file_path", "offset", "limit"} {
		if _, ok := properties[name]; !ok {
			t.Fatalf("read_file schema should expose %s, got properties %#v", name, properties)
		}
	}
	required, ok := schema["required"].([]any)
	if !ok {
		t.Fatalf("read_file schema required missing or wrong type: %#v", schema["required"])
	}
	for _, name := range []string{"offset", "limit"} {
		if containsAnyString(required, name) {
			t.Fatalf("read_file schema should not require %s, got %#v", name, required)
		}
	}
}

func schemaMapForTool(t *testing.T, name string) map[string]any {
	t.Helper()
	for _, tool := range AllToolDefs {
		if tool.Function == nil || tool.Function.Name != name {
			continue
		}
		var schema map[string]any
		if err := json.Unmarshal(RawSchemaFromParameters(tool.Function.Parameters), &schema); err != nil {
			t.Fatalf("unmarshal %s schema: %v", name, err)
		}
		return schema
	}
	t.Fatalf("expected AllToolDefs to include %s", name)
	return nil
}

func containsAnyString(values []any, want string) bool {
	for _, value := range values {
		if s, ok := value.(string); ok && s == want {
			return true
		}
	}
	return false
}

func TestSearchToolDescriptionsMatchPromptStrategy(t *testing.T) {
	toolsByName := map[string]string{}
	schemasByName := map[string]string{}
	for _, tool := range AllToolDefs {
		if tool.Function == nil {
			continue
		}
		toolsByName[tool.Function.Name] = tool.Function.Description
		schemasByName[tool.Function.Name] = string(RawSchemaFromParameters(tool.Function.Parameters))
	}
	for name, wants := range map[string][]string{
		"grep": {
			"A powerful search tool",
			"NEVER invoke grep or rg as a bash command",
			"Go regular expression",
			"Output modes",
			"files_with_matches",
			"content",
			"count",
			"glob parameter",
		},
		"glob": {
			"Fast file pattern matching tool",
			"any codebase size",
			"**/*.js",
			"sorted by modification time",
			"find files by name patterns",
		},
	} {
		description, ok := toolsByName[name]
		if !ok {
			t.Fatalf("expected AllToolDefs to include %s", name)
		}
		combined := description + "\n" + schemasByName[name]
		for _, want := range wants {
			if !strings.Contains(combined, want) {
				t.Fatalf("expected %s guidance to contain %q, got %q", name, want, combined)
			}
		}
	}
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

func TestTodoWriteToolDefinitionGuidesSequentialIDs(t *testing.T) {
	for _, tool := range AllToolDefs {
		if tool.Function == nil || tool.Function.Name != "todo_write" {
			continue
		}
		schema := string(RawSchemaFromParameters(tool.Function.Parameters))
		for _, want := range []string{
			"Create a new id when adding a new todo",
			"When updating an existing todo",
			"reuse the existing todo's id",
			"1, 2, 3, 4",
		} {
			if !strings.Contains(schema, want) {
				t.Fatalf("expected todo_write schema to mention %q, got %s", want, schema)
			}
		}
		return
	}
	t.Fatalf("expected AllToolDefs to include todo_write")
}

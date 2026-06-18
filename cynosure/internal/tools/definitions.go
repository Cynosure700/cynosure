package tools

import (
	"encoding/json"

	openai "github.com/sashabaranov/go-openai"
)

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func toolDef(name, desc string, params any) openai.Tool {
	return openai.Tool{
		Type: "function",
		Function: &openai.FunctionDefinition{
			Name:        name,
			Description: desc,
			Parameters:  mustMarshal(params),
		},
	}
}

const DefaultMaxResultSizeChars = 50000

type ToolSpec struct {
	Definition         openai.Tool
	MaxResultSizeChars int
}

func toolSpec(name, desc string, params any) ToolSpec {
	return ToolSpec{Definition: toolDef(name, desc, params), MaxResultSizeChars: DefaultMaxResultSizeChars}
}

func toolDefsFromSpecs(specs []ToolSpec) []openai.Tool {
	defs := make([]openai.Tool, 0, len(specs))
	for _, spec := range specs {
		defs = append(defs, spec.Definition)
	}
	return defs
}

func MaxResultSizeCharsForTool(name string) int {
	for _, spec := range AllToolSpecs {
		if spec.Definition.Function == nil || spec.Definition.Function.Name != name {
			continue
		}
		if spec.MaxResultSizeChars > 0 {
			return spec.MaxResultSizeChars
		}
		return DefaultMaxResultSizeChars
	}
	return DefaultMaxResultSizeChars
}

func strParam(desc string) map[string]any {
	return map[string]any{
		"type":        "string",
		"description": desc,
	}
}

func intParam(desc string) map[string]any {
	return map[string]any{
		"type":        "integer",
		"description": desc,
	}
}

func boolParam(desc string) map[string]any {
	return map[string]any{
		"type":        "boolean",
		"description": desc,
	}
}

func stringArrayParam(desc string) map[string]any {
	return map[string]any{
		"type":        "array",
		"items":       map[string]any{"type": "string"},
		"description": desc,
	}
}

var baseToolSpecs = []ToolSpec{
	toolSpec("bash", "Execute a shell command via bash -c. Relative path arguments are interpreted under the current working directory. Mutating commands (write, delete, curl, etc.) require user approval before running.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": strParam("The shell command to execute"),
		},
		"required": []string{"command"},
	}),
	toolSpec("read_file", "Read a file from the filesystem", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":  strParam("Path to the file to read"),
			"limit": intParam("Maximum number of lines to read"),
		},
		"required": []string{"path"},
	}),
	toolSpec("write_file", "Write content to a file", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    strParam("Path to the file to write"),
			"content": strParam("Content to write to the file"),
		},
		"required": []string{"path", "content"},
	}),
	toolSpec("edit_file", "Replace text in a file by exact match", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     strParam("Path to the file to edit"),
			"old_text": strParam("Exact text to find and replace"),
			"new_text": strParam("Text to replace it with"),
		},
		"required": []string{"path", "old_text", "new_text"},
	}),
	toolSpec("load_skill", "Load the full instructions of a local skill by exact name before using or following that skill. Skills are loaded from the user's ~/.cynosure/skills and the workspace .cynosure/skills directories, with workspace skills taking precedence.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": strParam("Name of the skill to load"),
		},
		"required": []string{"name"},
	}),
	toolSpec("todo_write", "Create or update the current task plan. Use this tool to track progress on multi-step tasks.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"todos": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":      map[string]any{"type": "string"},
						"content": map[string]any{"type": "string"},
						"status":  map[string]any{"type": "string", "enum": []string{TodoStatusPending, TodoStatusInProgress, TodoStatusCompleted}},
					},
					"required": []string{"id", "content", "status"},
				},
			},
		},
		"required": []string{"todos"},
	}),
	toolSpec("todo_list", "Read the current task plan and todo status without modifying it. Use this when you need to recover or confirm task state after context compression.", map[string]any{
		"type":                 "object",
		"properties":           map[string]any{},
		"additionalProperties": false,
	}),
	toolSpec("grep", "A fast content search tool that works in any size codebase. Searches file contents using a Go regular expression. Always use this tool for search tasks; never invoke grep or rg via bash.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern":     strParam("The regular expression pattern to search for in file contents (Go regexp syntax)."),
			"path":        strParam("File or directory to search in. Defaults to the current working directory."),
			"glob":        strParam("Glob pattern to filter files by name, e.g. *.go."),
			"output_mode": map[string]any{"type": "string", "enum": []string{"content", "files_with_matches", "count"}, "description": "Output mode: content shows matching lines, files_with_matches shows file paths (default), count shows match counts."},
			"-i":          boolParam("Case insensitive search."),
			"-n":          boolParam("Show line numbers in content output mode."),
			"head_limit":  intParam("Limit output to the first N entries. Defaults to 100."),
		},
		"required": []string{"pattern"},
	}),
	toolSpec("glob", "Fast file pattern matching tool that works with any codebase size. Supports glob patterns like **/*.js or src/**/*.ts and returns matching file paths sorted by modification time.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern":    strParam("The glob pattern to match files against."),
			"path":       strParam("The directory to search in. Defaults to the current working directory."),
			"head_limit": intParam("Limit output to the first N entries. Defaults to 100."),
		},
		"required": []string{"pattern"},
	}),
	toolSpec("ls", "List files and directories in a given path. The path must be an absolute path. You can optionally provide an array of glob patterns to ignore. Prefer glob and grep when you know which directories to search.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":   strParam("The absolute path to the directory to list (must be absolute, not relative)."),
			"ignore": stringArrayParam("List of glob patterns to ignore."),
		},
		"required": []string{"path"},
	}),
	toolSpec("multi_edit", "Make multiple edits to a single file in one operation, built on the edit_file tool. Prefer this over edit_file when making several edits to the same file. Edits are applied sequentially and atomically: if any edit fails, none are applied.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": strParam("The absolute path to the file to modify (must be absolute, not relative)."),
			"edits": map[string]any{
				"type":        "array",
				"description": "Array of edit operations to perform sequentially on the file.",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"old_string":  strParam("The text to replace (must match the file contents exactly, including whitespace and indentation)."),
						"new_string":  strParam("The text to replace old_string with."),
						"replace_all": boolParam("Replace all occurrences of old_string. Optional, defaults to false."),
					},
					"required": []string{"old_string", "new_string"},
				},
				"minItems": 1,
			},
		},
		"required": []string{"file_path", "edits"},
	}),
	toolSpec("web_fetch", "Fetch content from a URL and process it with an AI model. Fetches the URL, converts HTML to text, and runs the prompt over the content. Use this to retrieve and analyze web content.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"url":    strParam("The URL to fetch content from. http URLs are upgraded to https."),
			"prompt": strParam("The prompt describing what to extract or analyze from the page."),
		},
		"required": []string{"url", "prompt"},
	}),
}

var baseToolDefs = toolDefsFromSpecs(baseToolSpecs)

// webSearchToolSpec is defined separately because it is not enabled by default;
// it is exposed in AllToolDefs so users can opt in via configuration.
var webSearchToolSpec = toolSpec("web_search", "Search the web and use the results to inform responses. Provides up-to-date information for current events and recent data.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"query":           strParam("The search query to use."),
		"allowed_domains": stringArrayParam("Only include search results from these domains."),
		"blocked_domains": stringArrayParam("Never include search results from these domains."),
	},
	"required": []string{"query"},
})
var webSearchToolDef = webSearchToolSpec.Definition

var spawnSubagentToolSpec = toolSpec("spawn_subagent", "Spawn a child agent with a fresh message list to complete an isolated task. The child agent may use workspace tools, but it cannot spawn another subagent. Only its final summary is returned to the parent agent.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"task": strParam("The task for the child agent to complete. Include all context it needs because parent conversation history is not shared."),
		"cwd":  strParam("Optional working directory for the child agent. Relative paths are resolved under the workspace root; absolute paths must remain inside the workspace."),
	},
	"required": []string{"task"},
})
var spawnSubagentToolDef = spawnSubagentToolSpec.Definition

// ReadPersistedOutputToolName is exposed automatically alongside context
// compression so the model can fetch the full content behind a
// <persisted-output> marker when the inline preview is insufficient.
const ReadPersistedOutputToolName = "read_persisted_output"

var ReadPersistedOutputToolSpec = toolSpec(ReadPersistedOutputToolName, "Read a chunk of a persisted tool output by id when a <persisted-output> marker preview is insufficient. Only outputs from the current conversation are accessible.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"id":     strParam("The persisted output id from the <persisted-output> marker, for example po_abc123."),
		"offset": intParam("Zero-based character offset to start reading from. Defaults to 0."),
		"limit":  intParam("Maximum characters to return. Defaults to 20000 and is capped by the runtime."),
	},
	"required": []string{"id"},
})
var ReadPersistedOutputToolDef = ReadPersistedOutputToolSpec.Definition

var AllToolSpecs = append(append([]ToolSpec(nil), baseToolSpecs...), spawnSubagentToolSpec, webSearchToolSpec, ReadPersistedOutputToolSpec)
var AllToolDefs = toolDefsFromSpecs(AllToolSpecs)
var ChildToolDefs = baseToolDefs

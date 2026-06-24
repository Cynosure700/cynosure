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
	toolSpec("bash", "Execute a shell command via bash -c. The command runs in the workspace root, and relative path arguments are interpreted under it.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": strParam("The shell command to execute"),
		},
		"required": []string{"command"},
	}),
	toolSpec("read_file", "Reads a file from the local filesystem. You can access any file directly by using this tool. Assume this tool is able to read all files on the machine. If the user provides a path to a file, assume that path is valid. It is okay to read a user-provided file path that does not exist; an error will be returned. For paths inferred by you rather than provided by the user, confirm the file exists before reading.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":  strParam("Path to the file to read. If the user provides a path to a file, assume that path is valid; missing user-provided files are okay and will return an error. For inferred paths, confirm the file exists before reading."),
			"limit": intParam("Maximum number of lines to read"),
		},
		"required": []string{"path"},
	}),
	toolSpec("write_file", "Writes a file to the local filesystem. This tool will overwrite the existing file if there is one at the provided path. If this is an existing file, you MUST use read_file first to read the file's contents. Prefer edit_file for modifying existing files because it only sends the diff. Only use write_file to create new files or for complete rewrites. NEVER create documentation files (*.md) or README files unless explicitly requested by the user. Only use emojis if the user explicitly requests it; avoid writing emojis to files unless asked.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    strParam("Path to the file to write. Existing files are overwritten; use read_file first before overwriting an existing file."),
			"content": strParam("Complete file content to write. Avoid documentation files and emojis unless the user explicitly requested them."),
		},
		"required": []string{"path", "content"},
	}),
	toolSpec("edit_file", "Performs exact string replacements in files. You MUST use read_file first in the conversation before editing. When editing text from read_file output, preserve the exact indentation after any line number prefix; never include the line number prefix in old_text or new_text. ALWAYS prefer editing existing files in the codebase. Prefer multi_edit when making several edits to the same file. NEVER write new files unless explicitly required. Only use emojis if the user explicitly requests it; avoid adding emojis to files unless asked. The edit will fail if old_text is not unique in the file; provide a larger string with more surrounding context to make old_text must be unique.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     strParam("Path to the file to edit"),
			"old_text": strParam("Exact text to find and replace. Must match file contents exactly, preserve indentation, exclude read_file line number prefixes, and old_text must be unique unless the tool explicitly supports replacing all occurrences."),
			"new_text": strParam("Text to replace old_text with. Preserve surrounding style and avoid emojis unless explicitly requested."),
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
	toolSpec("todo_write", "Update the todo list for the current session. To be used proactively and often to track progress and pending tasks. Make sure that at least one task is in_progress at all times. Always provide both content (imperative) and activeForm (present continuous) for each task. Before calling this tool, first read the task status in the current session and check whether an update is needed; if no update is needed, do not call this tool.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"todos": map[string]any{
				"type": "array",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"id":         map[string]any{"type": "string", "description": "Required stable id for this todo item. Create a new id when adding a new todo, using simple sequential numbers such as 1, 2, 3, 4. When updating an existing todo, reuse the existing todo's id instead of generating a new one."},
						"content":    map[string]any{"type": "string", "description": "The task description in imperative form (e.g. \"Run tests\")."},
						"activeForm": map[string]any{"type": "string", "description": "The task description in present continuous form (e.g. \"Running tests\")."},
						"status":     map[string]any{"type": "string", "enum": []string{TodoStatusPending, TodoStatusInProgress, TodoStatusCompleted}, "description": "Current status of the task. Only one task may be in_progress at any given time; mark a task completed immediately once it is done before starting the next one."},
					},
					"required": []string{"id", "content", "activeForm", "status"},
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
	toolSpec("grep", "A powerful search tool for content search in any size codebase. ALWAYS use grep for search tasks. NEVER invoke grep or rg as a bash command; this tool is optimized for Cynosure permissions and workspace access. Supports Go regular expression syntax, e.g. log.*Error or function\\s+\\w+. Filter files with the glob parameter, e.g. *.go or **/*.ts. Output modes: content shows matching lines, files_with_matches shows only file paths (default), and count shows match counts. Pattern syntax uses Go regular expression syntax in this project; literal braces and other regexp metacharacters need escaping when they should be matched literally.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern":     strParam("The regular expression pattern to search for in file contents (Go regexp syntax)."),
			"path":        strParam("File or directory to search in. Defaults to the workspace root."),
			"glob":        strParam("Glob parameter to filter files by name, e.g. *.go or **/*.ts."),
			"output_mode": map[string]any{"type": "string", "enum": []string{"content", "files_with_matches", "count"}, "description": "Output mode: content shows matching lines, files_with_matches shows file paths (default), count shows match counts."},
			"-i":          boolParam("Case insensitive search."),
			"-n":          boolParam("Show line numbers in content output mode."),
			"head_limit":  intParam("Limit output to the first N entries. Defaults to 100."),
		},
		"required": []string{"pattern"},
	}),
	toolSpec("glob", "Fast file pattern matching tool that works with any codebase size. Supports glob patterns like **/*.js or src/**/*.ts. Returns matching file paths sorted by modification time. Use this tool when you need to find files by name patterns. When you are doing an open-ended search that may require multiple rounds of globbing and grepping, use an explore subagent instead.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"pattern":    strParam("The glob pattern to match files against, such as **/*.go or src/**/*.ts."),
			"path":       strParam("The directory to search in. Defaults to the workspace root."),
			"head_limit": intParam("Limit output to the first N entries. Defaults to 100."),
		},
		"required": []string{"pattern"},
	}),
	toolSpec("ls", "List files and directories in a given path. You can optionally provide an array of glob patterns to ignore. Prefer glob and grep when you know which directories to search.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":   strParam("The absolute path to the directory to list."),
			"ignore": stringArrayParam("List of glob patterns to ignore."),
		},
		"required": []string{"path"},
	}),
	toolSpec("multi_edit", "Performs multiple exact string replacements in one file. You MUST use read_file first in the conversation before editing. Prefer multi_edit over edit_file when making several edits to the same file. When editing text from read_file output, preserve the exact indentation after any line number prefix; never include the line number prefix in old_string or new_string. Each old_string must be unique unless replace_all is true. Edits are applied sequentially and atomically: if any edit fails, none are applied. Only use emojis if the user explicitly requests it; avoid adding emojis to files unless asked.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"file_path": strParam("The absolute path to the file to modify (must be absolute, not relative)."),
			"edits": map[string]any{
				"type":        "array",
				"description": "Array of edit operations to perform sequentially on the file.",
				"items": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"old_string":  strParam("The text to replace. It must match file contents exactly, including whitespace and indentation, exclude read_file line number prefixes, and old_string must be unique unless replace_all is true."),
						"new_string":  strParam("The text to replace old_string with. Preserve surrounding style and avoid emojis unless explicitly requested."),
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

// webSearchToolSpec 单独定义，因为它默认不启用；
// 它会暴露在 AllToolDefs 中，以便用户通过配置选择启用。
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

var spawnSubagentToolSpec = toolSpec("spawn_subagent", "Spawn a typed child agent with a fresh message list to complete an isolated task. Use sub_type=explore for read-only code search, file discovery, implementation tracing, and evidence gathering. Use sub_type=general for non-search isolated analysis or execution tasks. The child agent cannot spawn another subagent. Only its final summary is returned to the parent agent.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"sub_type": map[string]any{"type": "string", "enum": []string{"general", "explore"}, "description": "Type of child agent to spawn. Use explore for read-only codebase search and file exploration. Use general for non-search isolated analysis or execution tasks."},
		"task":     strParam("The task for the child agent to complete. Include all context it needs because parent conversation history is not shared."),
	},
	"required": []string{"sub_type", "task"},
})
var spawnSubagentToolDef = spawnSubagentToolSpec.Definition

// ReadPersistedOutputToolName 会随上下文压缩自动一并暴露，
// 以便当内联预览不足时，模型能获取 <persisted-output> 标记背后的完整内容。
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

// UpdateMemoryToolName / DeleteMemoryToolName 让模型可以维护那些被证明
// 错误或过时的长期项目记忆。它们在运行时层执行（而非无状态的 Dispatch），
// 因为记忆文件位于工作区之外，且必须与 MEMORY.md 保持同步。
const (
	UpdateMemoryToolName = "update_memory"
	DeleteMemoryToolName = "delete_memory"
)

var updateMemoryToolSpec = toolSpec(UpdateMemoryToolName, "Update a long-term project memory that is wrong or outdated. Identify the memory by its file path as shown in the MEMORY.md index (for example foo.md). Provide at least one of name, description, or body; the file and the MEMORY.md index are updated together.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"path":        strParam("The memory file path relative to the memory directory, as shown in the MEMORY.md index, e.g. foo.md."),
		"name":        strParam("Optional new short title for the memory."),
		"description": strParam("Optional new one-sentence description for the memory."),
		"body":        strParam("Optional new full body content for the memory."),
	},
	"required": []string{"path"},
})

var deleteMemoryToolSpec = toolSpec(DeleteMemoryToolName, "Delete a long-term project memory that is wrong or no longer applicable. Identify the memory by its file path as shown in the MEMORY.md index (for example foo.md). The file and its MEMORY.md index entry are removed together.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"path": strParam("The memory file path relative to the memory directory, as shown in the MEMORY.md index, e.g. foo.md."),
	},
	"required": []string{"path"},
})

var AllToolSpecs = append(append([]ToolSpec(nil), baseToolSpecs...), spawnSubagentToolSpec, webSearchToolSpec, ReadPersistedOutputToolSpec, updateMemoryToolSpec, deleteMemoryToolSpec)
var AllToolDefs = toolDefsFromSpecs(AllToolSpecs)
var ChildToolDefs = baseToolDefs

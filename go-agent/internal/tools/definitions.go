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

func strParam(desc string, required bool) map[string]any {
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

var baseToolDefs = []openai.Tool{
	toolDef("bash", "Execute a shell command via bash -c. Relative path arguments are interpreted under the workspace root; absolute paths outside the workspace and dangerous commands are rejected unless explicitly allowed by configuration.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"command": strParam("The shell command to execute", true),
		},
		"required": []string{"command"},
	}),
	toolDef("read_file", "Read a file from the filesystem", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":  strParam("Path to the file to read", true),
			"limit": intParam("Maximum number of lines to read"),
		},
		"required": []string{"path"},
	}),
	toolDef("write_file", "Write content to a file", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":    strParam("Path to the file to write", true),
			"content": strParam("Content to write to the file", true),
		},
		"required": []string{"path", "content"},
	}),
	toolDef("edit_file", "Replace text in a file by exact match", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"path":     strParam("Path to the file to edit", true),
			"old_text": strParam("Exact text to find and replace", true),
			"new_text": strParam("Text to replace it with", true),
		},
		"required": []string{"path", "old_text", "new_text"},
	}),
	toolDef("load_skill", "Load the full information of a skill by name. It first looks up the current user's enabled database skills, then falls back to local builtin skills.", map[string]any{
		"type": "object",
		"properties": map[string]any{
			"name": strParam("Name of the skill to load", true),
		},
		"required": []string{"name"},
	}),
}

var spawnSubagentToolDef = toolDef("spawn_subagent", "Spawn a child agent with a fresh message list to complete an isolated task. The child agent may use workspace tools, but it cannot spawn another subagent. Only its final summary is returned to the parent agent.", map[string]any{
	"type": "object",
	"properties": map[string]any{
		"task": strParam("The task for the child agent to complete. Include all context it needs because parent conversation history is not shared.", true),
		"cwd":  strParam("Optional working directory for the child agent. Relative paths are resolved under the workspace root; absolute paths must remain inside the workspace.", false),
	},
	"required": []string{"task"},
})

var AllToolDefs = append(append([]openai.Tool(nil), baseToolDefs...), spawnSubagentToolDef)
var ChildToolDefs = baseToolDefs

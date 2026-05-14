package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	"nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

type ToolContext struct {
	User         storage.User
	Conversation storage.Conversation
	Loader       *sessions.SkillLoader
}

type ToolRegistry struct {
	store *storage.Store
	cfg   config.AppConfig
}

func NewToolRegistry(store *storage.Store, cfg config.AppConfig) *ToolRegistry {
	return &ToolRegistry{store: store, cfg: cfg}
}

func (r *ToolRegistry) Definitions(loader *sessions.SkillLoader) []openai.Tool {
	toolDefs := []openai.Tool{
		toolDef("read_file", "Read a file from the current user's workspace", map[string]any{"type": "object", "properties": map[string]any{"path": map[string]any{"type": "string"}}, "required": []string{"path"}}),
		toolDef("write_file", "Write a file inside the current user's workspace", map[string]any{"type": "object", "properties": map[string]any{"path": map[string]any{"type": "string"}, "content": map[string]any{"type": "string"}}, "required": []string{"path", "content"}}),
		toolDef("edit_file", "Replace exact text inside a file in the current user's workspace", map[string]any{"type": "object", "properties": map[string]any{"path": map[string]any{"type": "string"}, "old_text": map[string]any{"type": "string"}, "new_text": map[string]any{"type": "string"}}, "required": []string{"path", "old_text", "new_text"}}),
	}
	if loader != nil && loader.GetDescriptions() != "" {
		toolDefs = append(toolDefs, toolDef("load_skill", "Load one of the current user's enabled skills", map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}}, "required": []string{"name"}}))
	}
	return toolDefs
}

func (r *ToolRegistry) Execute(ctx context.Context, toolCtx ToolContext, name string, rawArgs string) (string, error) {
	var args map[string]any
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", fmt.Errorf("invalid tool arguments: %w", err)
	}
	if !r.isAllowed(name) {
		return "", fmt.Errorf("tool %s is not registered for web runtime", name)
	}
	if name == "load_skill" {
		skillName, _ := args["name"].(string)
		return toolCtx.Loader.GetContent(skillName)
	}
	path, _ := args["path"].(string)
	workspacePath, err := r.resolveUserPath(toolCtx.User.ID, path)
	if err != nil {
		return "", err
	}
	relPath := workspacePath
	switch name {
	case "read_file":
		return tools.RunRead(relPath, 0)
	case "write_file":
		content, _ := args["content"].(string)
		return tools.RunWrite(relPath, content)
	case "edit_file":
		oldText, _ := args["old_text"].(string)
		newText, _ := args["new_text"].(string)
		return tools.RunEdit(relPath, oldText, newText)
	default:
		return "", fmt.Errorf("tool %s is not implemented", name)
	}
}

func (r *ToolRegistry) resolveUserPath(userID, inputPath string) (string, error) {
	if strings.TrimSpace(inputPath) == "" {
		return "", fmt.Errorf("path is required")
	}
	workspaceRoot := filepath.Join(r.cfg.WorkspaceRoot, userID)
	if err := os.MkdirAll(workspaceRoot, 0o755); err != nil {
		return "", err
	}
	absolute, err := filepath.Abs(filepath.Join(workspaceRoot, inputPath))
	if err != nil {
		return "", err
	}
	rootAbs, err := filepath.Abs(workspaceRoot)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absolute, rootAbs) {
		return "", fmt.Errorf("path escapes user workspace")
	}
	return absolute, nil
}

func (r *ToolRegistry) isAllowed(name string) bool {
	allowed := []string{"read_file", "write_file", "edit_file", "load_skill"}
	for _, item := range allowed {
		if item == name {
			return true
		}
	}
	return false
}

func toolDef(name, desc string, params any) openai.Tool {
	return openai.Tool{Type: "function", Function: &openai.FunctionDefinition{Name: name, Description: desc, Parameters: mustMarshal(params)}}
}

func mustMarshal(v any) json.RawMessage {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return data
}

func RegisteredTools() []string {
	tools := []string{"edit_file", "load_skill", "read_file", "write_file"}
	sort.Strings(tools)
	return tools
}

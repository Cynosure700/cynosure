package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
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
	toolDefs := []openai.Tool{}
	if loader != nil && loader.GetDescriptions() != "" {
		toolDefs = append(toolDefs, toolDef("load_skill", "Load one of the current user's enabled capabilities", map[string]any{"type": "object", "properties": map[string]any{"name": map[string]any{"type": "string"}}, "required": []string{"name"}}))
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
		if toolCtx.Loader == nil {
			return "", fmt.Errorf("no capabilities are available in this conversation")
		}
		skillName, _ := args["name"].(string)
		return toolCtx.Loader.GetContent(skillName)
	}
	return "", fmt.Errorf("tool %s is not available in browser chat", name)
}

func (r *ToolRegistry) isAllowed(name string) bool {
	allowed := []string{"load_skill"}
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
	tools := []string{"load_skill"}
	sort.Strings(tools)
	return tools
}

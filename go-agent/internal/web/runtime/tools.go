package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
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

const defaultWebAllowedTool = "load_skill"

func NewToolRegistry(store *storage.Store, cfg config.AppConfig) *ToolRegistry {
	return &ToolRegistry{store: store, cfg: cfg}
}

func (r *ToolRegistry) Definitions(loader *sessions.SkillLoader) []openai.Tool {
	allowed := r.allowedToolNames()
	toolDefs := make([]openai.Tool, 0, len(allowed))
	for _, name := range allowed {
		if name == "load_skill" && (loader == nil || loader.GetDescriptions() == "") {
			continue
		}
		def, ok := lookupRegisteredTool(name)
		if !ok {
			continue
		}
		toolDefs = append(toolDefs, def)
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
	handler, ok := agenttools.Handlers[name]
	if !ok || handler == nil {
		return "", fmt.Errorf("tool %s has no handler", name)
	}
	return handler(ctx, args)
}

func (r *ToolRegistry) isAllowed(name string) bool {
	for _, toolName := range r.allowedToolNames() {
		if toolName == name {
			return true
		}
	}
	return false
}

func (r *ToolRegistry) allowedToolNames() []string {
	configured := r.cfg.WebAllowedTools
	if len(configured) == 0 {
		configured = []string{defaultWebAllowedTool}
	}

	names := make([]string, 0, len(configured))
	seen := make(map[string]struct{}, len(configured))
	for _, name := range configured {
		if _, ok := lookupRegisteredTool(name); !ok {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names
}

func lookupRegisteredTool(name string) (openai.Tool, bool) {
	for _, tool := range agenttools.ChildToolDefs {
		if tool.Function != nil && tool.Function.Name == name {
			return tool, true
		}
	}
	return openai.Tool{}, false
}

func RegisteredTools(cfg config.AppConfig) []string {
	registry := NewToolRegistry(nil, cfg)
	names := registry.allowedToolNames()
	sort.Strings(names)
	return names
}

package runtime

import (
	"sort"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/config"
	agenttools "nano_cc/internal/tools"
)

func loadAllowedToolNames(cfg config.AppConfig) []string {
	configured := cfg.WebAllowedTools
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

func runtimeEnvFromConfig(cfg config.AppConfig) agenttools.RuntimeEnv {
	return agenttools.RuntimeEnv{
		AppHome:                strings.TrimSpace(cfg.AppHome),
		CommandBinDir:          strings.TrimSpace(cfg.CommandBinDir),
		CommandScriptDir:       strings.TrimSpace(cfg.CommandScriptDir),
		WorkspaceRoot:          strings.TrimSpace(cfg.WorkspaceRoot),
		AllowOutsideWorkspace:  cfg.BashAllowOutsideWorkspace,
		AllowDangerousCommands: cfg.BashAllowDangerousCommands,
	}
}

func buildToolDefinitions(allowed []string) []openai.Tool {
	toolDefs := make([]openai.Tool, 0, len(allowed))
	for _, name := range allowed {
		def, ok := lookupRegisteredTool(name)
		if !ok {
			continue
		}
		toolDefs = append(toolDefs, def)
	}
	return toolDefs
}

func lookupRegisteredTool(name string) (openai.Tool, bool) {
	for _, tool := range agenttools.AllToolDefs {
		if tool.Function != nil && tool.Function.Name == name {
			return tool, true
		}
	}
	return openai.Tool{}, false
}

func RegisteredTools(cfg config.AppConfig) []string {
	names := loadAllowedToolNames(cfg)
	sort.Strings(names)
	return names
}

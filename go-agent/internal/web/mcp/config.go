package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"nano_cc/internal/logger"
	"nano_cc/internal/web/storage"
)

type builtinConfig struct {
	Servers *[]builtinServerConfig `json:"mcp_servers"`
}

type builtinServerConfig struct {
	Name      string            `json:"name"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	Enabled   *bool             `json:"enabled"`
	Transport *string           `json:"transport"`
	URL       *string           `json:"url"`
	Headers   map[string]string `json:"headers"`
}

func LoadBuiltinConfig(path string) ([]storage.MCPServer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info(fmt.Sprintf("mcp: builtin config not found path=%s, skip builtin servers", path))
			return nil, nil
		}
		return nil, err
	}

	var cfg builtinConfig
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse mcp_config.json: %w", err)
	}
	if cfg.Servers == nil {
		return nil, fmt.Errorf("mcp_config.json requires mcp_servers")
	}

	servers := make([]storage.MCPServer, 0, len(*cfg.Servers))
	seen := make(map[string]struct{}, len(*cfg.Servers))
	for i, item := range *cfg.Servers {
		server, key, err := builtinServer(i, item)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("builtin mcp server %d has duplicate sanitized name %q", i, key)
		}
		seen[key] = struct{}{}
		servers = append(servers, server)
	}
	logger.Info(fmt.Sprintf("mcp: loaded builtin stdio servers count=%d", len(servers)))
	return servers, nil
}

func builtinServer(index int, item builtinServerConfig) (storage.MCPServer, string, error) {
	if item.Transport != nil {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d must not set transport; builtin mcp only supports stdio", index)
	}
	if item.URL != nil {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d must not set url; builtin mcp only supports stdio", index)
	}
	if item.Headers != nil {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d must not set headers; builtin mcp only supports stdio", index)
	}

	name := strings.TrimSpace(item.Name)
	if name == "" {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %d name is required", index)
	}
	command := strings.TrimSpace(item.Command)
	if command == "" {
		return storage.MCPServer{}, "", fmt.Errorf("builtin mcp server %s command is required", name)
	}

	enabled := true
	if item.Enabled != nil {
		enabled = *item.Enabled
	}
	key := sanitizeName(name)
	return storage.MCPServer{
		ID:        "builtin:" + key,
		UserID:    "",
		Name:      "builtin_" + name,
		Transport: "stdio",
		Command:   command,
		Args:      item.Args,
		Env:       item.Env,
		Enabled:   enabled,
	}, key, nil
}

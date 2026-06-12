package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
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

type workspaceConfig struct {
	ServersArray *[]workspaceServerConfig         `json:"mcp_servers"`
	ServersMap   map[string]workspaceServerConfig `json:"mcpServers"`
}

type workspaceServerConfig struct {
	Name      string            `json:"name"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	Enabled   *bool             `json:"enabled"`
	Transport string            `json:"transport"`
	URL       string            `json:"url"`
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

func LoadWorkspaceConfig(path string) ([]storage.MCPServer, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Info(fmt.Sprintf("mcp: workspace config not found path=%s, skip workspace servers", path))
			return nil, nil
		}
		return nil, err
	}

	var cfg workspaceConfig
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("parse workspace mcp config %s: %w", path, err)
	}

	items := make([]workspaceServerConfig, 0)
	if cfg.ServersArray != nil {
		items = append(items, (*cfg.ServersArray)...)
	}
	if cfg.ServersMap != nil {
		keys := make([]string, 0, len(cfg.ServersMap))
		for key := range cfg.ServersMap {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			item := cfg.ServersMap[key]
			if strings.TrimSpace(item.Name) == "" {
				item.Name = key
			}
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("workspace mcp config %s requires mcp_servers or mcpServers", path)
	}

	servers := make([]storage.MCPServer, 0, len(items))
	seen := make(map[string]struct{}, len(items))
	for i, item := range items {
		server, key, err := workspaceServer(i, item)
		if err != nil {
			return nil, err
		}
		if _, ok := seen[key]; ok {
			return nil, fmt.Errorf("workspace mcp server %d has duplicate sanitized name %q", i, key)
		}
		seen[key] = struct{}{}
		servers = append(servers, server)
	}
	logger.Info(fmt.Sprintf("mcp: loaded workspace servers count=%d", len(servers)))
	return servers, nil
}

func workspaceServer(index int, item workspaceServerConfig) (storage.MCPServer, string, error) {
	name := strings.TrimSpace(item.Name)
	if name == "" {
		return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %d name is required", index)
	}
	transport := strings.TrimSpace(item.Transport)
	if transport == "" && strings.TrimSpace(item.Command) != "" {
		transport = "stdio"
	}
	if transport == "" {
		return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s transport is required", name)
	}
	enabled := true
	if item.Enabled != nil {
		enabled = *item.Enabled
	}
	command := strings.TrimSpace(item.Command)
	url := strings.TrimSpace(item.URL)
	switch transport {
	case "stdio":
		if command == "" {
			return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s command is required for stdio", name)
		}
	case "sse", "streamable":
		if url == "" {
			return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s url is required for %s", name, transport)
		}
	default:
		return storage.MCPServer{}, "", fmt.Errorf("workspace mcp server %s unsupported transport %q", name, transport)
	}
	key := sanitizeName(name)
	return storage.MCPServer{
		ID:        "workspace:" + key,
		Name:      name,
		Transport: transport,
		Command:   command,
		Args:      item.Args,
		Env:       item.Env,
		URL:       url,
		Headers:   item.Headers,
		Enabled:   enabled,
	}, key, nil
}

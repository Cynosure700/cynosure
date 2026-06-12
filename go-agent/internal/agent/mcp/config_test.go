package mcp

import (
	"os"
	"path/filepath"
	"testing"
)

func writeBuiltinConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "mcp_config.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write builtin config: %v", err)
	}
	return path
}

func TestLoadBuiltinConfigMissingFileReturnsEmpty(t *testing.T) {
	servers, err := LoadBuiltinConfig(filepath.Join(t.TempDir(), "mcp_config.json"))
	if err != nil {
		t.Fatalf("expected missing file to be ignored, got %v", err)
	}
	if len(servers) != 0 {
		t.Fatalf("expected no builtin servers, got %#v", servers)
	}
}

func TestLoadBuiltinConfigParsesStdioServers(t *testing.T) {
	path := writeBuiltinConfig(t, `{
		"mcp_servers": [{
			"name": "filesystem",
			"command": "npx",
			"args": ["-y", "@modelcontextprotocol/server-filesystem", "/tmp"],
			"env": {"NODE_ENV": "production"}
		}, {
			"name": "disabled",
			"command": "disabled-server",
			"enabled": false
		}]
	}`)

	servers, err := LoadBuiltinConfig(path)
	if err != nil {
		t.Fatalf("load builtin config: %v", err)
	}
	if len(servers) != 2 {
		t.Fatalf("expected two servers, got %d", len(servers))
	}
	first := servers[0]
	if first.ID != "builtin:filesystem" {
		t.Fatalf("unexpected id: %q", first.ID)
	}
	if first.UserID != "" {
		t.Fatalf("expected system user id to be empty, got %q", first.UserID)
	}
	if first.Name != "builtin_filesystem" {
		t.Fatalf("unexpected name: %q", first.Name)
	}
	if first.Transport != "stdio" {
		t.Fatalf("unexpected transport: %q", first.Transport)
	}
	if first.Command != "npx" {
		t.Fatalf("unexpected command: %q", first.Command)
	}
	if len(first.Args) != 3 || first.Args[2] != "/tmp" {
		t.Fatalf("unexpected args: %#v", first.Args)
	}
	if first.Env["NODE_ENV"] != "production" {
		t.Fatalf("unexpected env: %#v", first.Env)
	}
	if !first.Enabled {
		t.Fatalf("expected enabled to default true")
	}
	if servers[1].Enabled {
		t.Fatalf("expected explicit enabled=false to be preserved")
	}
}

func TestLoadBuiltinConfigRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty name", body: `{"mcp_servers":[{"name":" ","command":"cmd"}]}`},
		{name: "empty command", body: `{"mcp_servers":[{"name":"docs","command":" "}]}`},
		{name: "duplicate sanitized name", body: `{"mcp_servers":[{"name":"docs-api","command":"cmd"},{"name":"docs_api","command":"cmd"}]}`},
		{name: "transport rejected", body: `{"mcp_servers":[{"name":"docs","transport":"stdio","command":"cmd"}]}`},
		{name: "url rejected", body: `{"mcp_servers":[{"name":"docs","command":"cmd","url":"https://example.com"}]}`},
		{name: "headers rejected", body: `{"mcp_servers":[{"name":"docs","command":"cmd","headers":{}}]}`},
		{name: "missing mcp_servers rejected", body: `{}`},
		{name: "null mcp_servers rejected", body: `{"mcp_servers":null}`},
		{name: "unknown top-level field rejected", body: `{"mcpServers":[{"name":"docs","command":"cmd"}]}`},
		{name: "unknown server field rejected", body: `{"mcp_servers":[{"name":"docs","command":"cmd","extra":true}]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeBuiltinConfig(t, tt.body)
			if _, err := LoadBuiltinConfig(path); err == nil {
				t.Fatalf("expected invalid config to be rejected")
			}
		})
	}
}

func writeWorkspaceConfig(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".mcp.json")
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write workspace config: %v", err)
	}
	return path
}

func TestLoadWorkspaceConfigParsesMCPServersArray(t *testing.T) {
	path := writeWorkspaceConfig(t, `{
		"mcp_servers": [{
			"name": "filesystem",
			"transport": "stdio",
			"command": "npx",
			"args": ["-y", "server"],
			"env": {"A":"B"}
		}]
	}`)

	servers, err := LoadWorkspaceConfig(path)
	if err != nil {
		t.Fatalf("LoadWorkspaceConfig returned error: %v", err)
	}
	if len(servers) != 1 {
		t.Fatalf("expected one server, got %d", len(servers))
	}
	server := servers[0]
	if server.ID != "workspace:filesystem" || server.Name != "filesystem" || server.Transport != "stdio" {
		t.Fatalf("unexpected server identity: %#v", server)
	}
	if server.Command != "npx" || len(server.Args) != 2 || server.Env["A"] != "B" {
		t.Fatalf("unexpected stdio config: %#v", server)
	}
	if !server.Enabled {
		t.Fatalf("enabled should default true")
	}
}

func TestLoadWorkspaceConfigParsesMCPServersMap(t *testing.T) {
	path := writeWorkspaceConfig(t, `{
		"mcpServers": {
			"docs-api": {
				"url": "https://example.com/sse",
				"transport": "sse",
				"headers": {"Authorization":"Bearer token"},
				"enabled": false
			}
		}
	}`)

	servers, err := LoadWorkspaceConfig(path)
	if err != nil {
		t.Fatalf("LoadWorkspaceConfig returned error: %v", err)
	}
	server := servers[0]
	if server.ID != "workspace:docs_api" || server.Name != "docs-api" || server.Transport != "sse" {
		t.Fatalf("unexpected server identity: %#v", server)
	}
	if server.URL != "https://example.com/sse" || server.Headers["Authorization"] != "Bearer token" {
		t.Fatalf("unexpected sse config: %#v", server)
	}
	if server.Enabled {
		t.Fatalf("enabled=false should be preserved")
	}
}

func TestLoadWorkspaceConfigRejectsInvalidServers(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "missing name", body: `{"mcp_servers":[{"command":"cmd"}]}`},
		{name: "missing command", body: `{"mcp_servers":[{"name":"fs","transport":"stdio"}]}`},
		{name: "missing url", body: `{"mcp_servers":[{"name":"docs","transport":"sse"}]}`},
		{name: "duplicate sanitized", body: `{"mcp_servers":[{"name":"docs-api","command":"cmd"},{"name":"docs_api","command":"cmd"}]}`},
		{name: "unknown format", body: `{}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := LoadWorkspaceConfig(writeWorkspaceConfig(t, tt.body)); err == nil {
				t.Fatalf("expected invalid workspace config to fail")
			}
		})
	}
}

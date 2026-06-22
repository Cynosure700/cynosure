package mcp

import (
	"strings"
	"testing"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"cynosure/internal/agent/storage"
)

func TestBuildBuiltinStdioTransportRejectsEmptyCommand(t *testing.T) {
	_, err := buildBuiltinStdioTransport(storage.MCPServer{Command: " "})
	if err == nil {
		t.Fatalf("expected empty command to be rejected")
	}
}

func TestBuildBuiltinStdioTransportCreatesCommandTransport(t *testing.T) {
	transport, err := buildBuiltinStdioTransport(storage.MCPServer{
		Command: "cmd",
		Args:    []string{"arg1", "arg2"},
		Env:     map[string]string{"MCP_TEST_ENV": "enabled"},
	})
	if err != nil {
		t.Fatalf("build builtin stdio transport: %v", err)
	}
	cmdTransport, ok := transport.(*mcpsdk.CommandTransport)
	if !ok {
		t.Fatalf("expected CommandTransport, got %T", transport)
	}
	if cmdTransport.Command.Path != "cmd" {
		t.Fatalf("unexpected command path: %q", cmdTransport.Command.Path)
	}
	if len(cmdTransport.Command.Args) != 3 || cmdTransport.Command.Args[1] != "arg1" || cmdTransport.Command.Args[2] != "arg2" {
		t.Fatalf("unexpected command args: %#v", cmdTransport.Command.Args)
	}
	if !containsEnv(cmdTransport.Command.Env, "MCP_TEST_ENV=enabled") {
		t.Fatalf("expected custom env to be present in %#v", cmdTransport.Command.Env)
	}
}

func TestBuildTransportSupportsWorkspaceStdio(t *testing.T) {
	transport, err := buildTransport(storage.MCPServer{Transport: "stdio", Command: "cmd", Args: []string{"arg"}})
	if err != nil {
		t.Fatalf("buildTransport returned error: %v", err)
	}
	if _, ok := transport.(*mcpsdk.CommandTransport); !ok {
		t.Fatalf("expected CommandTransport, got %T", transport)
	}
}

func containsEnv(env []string, expected string) bool {
	for _, item := range env {
		if item == expected || strings.HasPrefix(item, expected) {
			return true
		}
	}
	return false
}

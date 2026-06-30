package mcp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

// headerRoundTripper 在每个 HTTP 请求上附加自定义请求头，用于 sse/streamable 传输鉴权。
type headerRoundTripper struct {
	base    http.RoundTripper
	headers map[string]string
}

func (t *headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		if strings.TrimSpace(key) == "" {
			continue
		}
		req.Header.Set(key, value)
	}
	return t.base.RoundTrip(req)
}

func httpClientWithHeaders(headers map[string]string) *http.Client {
	if len(headers) == 0 {
		return nil
	}
	return &http.Client{Transport: &headerRoundTripper{base: http.DefaultTransport, headers: headers}}
}

// buildTransport 根据配置的 transport 类型构造对应的 MCP 客户端传输层。
func buildTransport(server storage.MCPServer) (mcpsdk.Transport, error) {
	switch server.Transport {
	case "stdio":
		return buildStdioTransport(server)
	case "sse":
		if strings.TrimSpace(server.URL) == "" {
			return nil, fmt.Errorf("sse transport requires a url")
		}
		return &mcpsdk.SSEClientTransport{Endpoint: server.URL, HTTPClient: httpClientWithHeaders(server.Headers)}, nil
	case "streamable":
		if strings.TrimSpace(server.URL) == "" {
			return nil, fmt.Errorf("streamable transport requires a url")
		}
		return &mcpsdk.StreamableClientTransport{Endpoint: server.URL, HTTPClient: httpClientWithHeaders(server.Headers)}, nil
	default:
		return nil, fmt.Errorf("unsupported transport %q", server.Transport)
	}
}

func buildBuiltinStdioTransport(server storage.MCPServer) (mcpsdk.Transport, error) {
	return buildStdioTransport(server)
}

func buildStdioTransport(server storage.MCPServer) (mcpsdk.Transport, error) {
	command := strings.TrimSpace(server.Command)
	if command == "" {
		return nil, fmt.Errorf("stdio transport requires a command")
	}
	cmd := exec.Command(command, server.Args...)
	cmd.Env = append(os.Environ(), envPairs(server.Env)...)
	return &mcpsdk.CommandTransport{Command: cmd}, nil
}

func envPairs(env map[string]string) []string {
	if len(env) == 0 {
		return nil
	}
	pairs := make([]string, 0, len(env))
	for key, value := range env {
		if strings.TrimSpace(key) == "" {
			continue
		}
		pairs = append(pairs, key+"="+value)
	}
	return pairs
}

// serializeContent 把 MCP 工具调用返回的 content blocks 拼成纯文本，契合现有工具结果的字符串契约。
func serializeContent(result *mcpsdk.CallToolResult) string {
	if result == nil {
		return ""
	}
	var sb strings.Builder
	for _, content := range result.Content {
		switch c := content.(type) {
		case *mcpsdk.TextContent:
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(c.Text)
		default:
			data, err := json.Marshal(content)
			if err != nil {
				continue
			}
			if sb.Len() > 0 {
				sb.WriteString("\n")
			}
			sb.Write(data)
		}
	}
	if sb.Len() == 0 && result.StructuredContent != nil {
		if data, err := json.Marshal(result.StructuredContent); err == nil {
			sb.Write(data)
		}
	}
	return sb.String()
}

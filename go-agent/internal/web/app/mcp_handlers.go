package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/storage"
)

type mcpServerBody struct {
	Name      string            `json:"name"`
	Transport string            `json:"transport"`
	Command   string            `json:"command"`
	Args      []string          `json:"args"`
	Env       map[string]string `json:"env"`
	URL       string            `json:"url"`
	Headers   map[string]string `json:"headers"`
	Enabled   *bool             `json:"enabled"`
}

func normalizeTransport(transport string) string {
	switch transport {
	case "sse", "streamable":
		return transport
	default:
		return ""
	}
}

func validateMCPServer(server storage.MCPServer) error {
	if strings.TrimSpace(server.Name) == "" {
		return fmt.Errorf("mcp server name is required")
	}
	if server.Transport == "" {
		return fmt.Errorf("transport must be one of sse/streamable")
	}
	switch server.Transport {
	case "sse", "streamable":
		if strings.TrimSpace(server.URL) == "" {
			return fmt.Errorf("url is required for %s transport", server.Transport)
		}
	}
	return nil
}

func applyMCPServerBody(server *storage.MCPServer, body mcpServerBody) {
	server.Name = strings.TrimSpace(body.Name)
	server.Transport = normalizeTransport(body.Transport)
	server.Command = body.Command
	server.Args = body.Args
	server.Env = body.Env
	server.URL = body.URL
	server.Headers = body.Headers
	if body.Enabled != nil {
		server.Enabled = *body.Enabled
	}
}

func (s *Server) handleMCPServers(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		servers, err := s.store.ListMCPServersByUser(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"mcp_servers": servers})
	case http.MethodPost:
		var body mcpServerBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		server := storage.MCPServer{ID: newID("mcp"), UserID: user.ID, Enabled: true}
		applyMCPServerBody(&server, body)
		if err := validateMCPServer(server); err != nil {
			badRequest(w, err)
			return
		}
		if err := s.store.CreateMCPServer(r.Context(), server); err != nil {
			badRequest(w, err)
			return
		}
		s.mcpManager.Invalidate(user.ID)
		writeJSON(w, http.StatusCreated, map[string]any{"mcp_server": server})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleMCPServerByID(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	path := strings.TrimPrefix(r.URL.Path, "/api/mcp-servers/")
	parts := strings.Split(path, "/")
	serverID := parts[0]
	if serverID == "" {
		http.NotFound(w, r)
		return
	}
	server, err := s.store.GetMCPServerByID(r.Context(), serverID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if server.UserID != user.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// 子路径：/api/mcp-servers/{id}/test
	if len(parts) == 2 && parts[1] == "test" {
		if r.Method != http.MethodPost {
			methodNotAllowed(w)
			return
		}
		s.handleMCPServerTest(w, r, server)
		return
	}

	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"mcp_server": server})
	case http.MethodPut:
		var body mcpServerBody
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		applyMCPServerBody(&server, body)
		if err := validateMCPServer(server); err != nil {
			badRequest(w, err)
			return
		}
		if err := s.store.UpdateMCPServer(r.Context(), server); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.mcpManager.Invalidate(user.ID)
		writeJSON(w, http.StatusOK, map[string]any{"mcp_server": server})
	case http.MethodPatch:
		var body struct {
			Enabled bool `json:"enabled"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		server.Enabled = body.Enabled
		if err := s.store.UpdateMCPServer(r.Context(), server); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.mcpManager.Invalidate(user.ID)
		writeJSON(w, http.StatusOK, map[string]any{"mcp_server": server})
	case http.MethodDelete:
		if err := s.store.DeleteMCPServer(r.Context(), server.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.mcpManager.Invalidate(user.ID)
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleMCPServerTest(w http.ResponseWriter, r *http.Request, server storage.MCPServer) {
	ctx, cancel := context.WithTimeout(r.Context(), 35*time.Second)
	defer cancel()
	tools, err := s.mcpManager.TestServer(ctx, server)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{"ok": false, "error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "tools": tools})
}

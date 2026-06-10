package storage

import (
	"context"
	"database/sql"
	"encoding/json"
)

func encodeJSONField(value any) (string, error) {
	if value == nil {
		return "null", nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (s *Store) CreateMCPServer(ctx context.Context, server MCPServer) error {
	argsJSON, err := encodeJSONField(server.Args)
	if err != nil {
		return err
	}
	envJSON, err := encodeJSONField(server.Env)
	if err != nil {
		return err
	}
	headersJSON, err := encodeJSONField(server.Headers)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx, `
		INSERT INTO mcp_servers (id, user_id, name, transport, command, args_json, env_json, url, headers_json, enabled, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, server.ID, server.UserID, server.Name, server.Transport, server.Command, argsJSON, envJSON, server.URL, headersJSON, server.Enabled)
	return err
}

func (s *Store) ListMCPServersByUser(ctx context.Context, userID string) ([]MCPServer, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, name, transport, command, args_json, env_json, url, headers_json, enabled, created_at, updated_at
		FROM mcp_servers WHERE user_id = ? ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMCPServers(rows)
}

func (s *Store) ListEnabledMCPServersByUser(ctx context.Context, userID string) ([]MCPServer, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, name, transport, command, args_json, env_json, url, headers_json, enabled, created_at, updated_at
		FROM mcp_servers WHERE user_id = ? AND enabled = 1 ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanMCPServers(rows)
}

func (s *Store) GetMCPServerByID(ctx context.Context, id string) (MCPServer, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, name, transport, command, args_json, env_json, url, headers_json, enabled, created_at, updated_at
		FROM mcp_servers WHERE id = ?
	`, id)
	return scanMCPServer(row)
}

func (s *Store) UpdateMCPServer(ctx context.Context, server MCPServer) error {
	argsJSON, err := encodeJSONField(server.Args)
	if err != nil {
		return err
	}
	envJSON, err := encodeJSONField(server.Env)
	if err != nil {
		return err
	}
	headersJSON, err := encodeJSONField(server.Headers)
	if err != nil {
		return err
	}
	_, err = s.DB.ExecContext(ctx, `
		UPDATE mcp_servers
		SET name = ?, transport = ?, command = ?, args_json = ?, env_json = ?, url = ?, headers_json = ?, enabled = ?, updated_at = NOW()
		WHERE id = ?
	`, server.Name, server.Transport, server.Command, argsJSON, envJSON, server.URL, headersJSON, server.Enabled, server.ID)
	return err
}

func (s *Store) DeleteMCPServer(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM mcp_servers WHERE id = ?`, id)
	return err
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMCPServer(scanner rowScanner) (MCPServer, error) {
	var (
		server      MCPServer
		argsJSON    string
		envJSON     string
		headersJSON string
		enabled     int
	)
	if err := scanner.Scan(&server.ID, &server.UserID, &server.Name, &server.Transport, &server.Command, &argsJSON, &envJSON, &server.URL, &headersJSON, &enabled, &server.CreatedAt, &server.UpdatedAt); err != nil {
		return MCPServer{}, err
	}
	server.Enabled = enabled != 0
	if argsJSON != "" {
		_ = json.Unmarshal([]byte(argsJSON), &server.Args)
	}
	if envJSON != "" {
		_ = json.Unmarshal([]byte(envJSON), &server.Env)
	}
	if headersJSON != "" {
		_ = json.Unmarshal([]byte(headersJSON), &server.Headers)
	}
	return server, nil
}

func scanMCPServers(rows *sql.Rows) ([]MCPServer, error) {
	servers := make([]MCPServer, 0)
	for rows.Next() {
		server, err := scanMCPServer(rows)
		if err != nil {
			return nil, err
		}
		servers = append(servers, server)
	}
	return servers, rows.Err()
}

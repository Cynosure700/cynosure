package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"

	"nano_cc/internal/config"
)

type Store struct {
	DB    *sql.DB
	Redis *redis.Client
	Cfg   config.AppConfig
}

func NewStore(cfg config.AppConfig) (*Store, error) {
	db, err := sql.Open("mysql", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open mysql: %w", err)
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	return &Store{DB: db, Redis: rdb, Cfg: cfg}, nil
}

func (s *Store) HealthCheck(ctx context.Context) error {
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("mysql ping: %w", err)
	}
	if err := s.Redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}
	return nil
}

func (s *Store) Close() error {
	if s.Redis != nil {
		_ = s.Redis.Close()
	}
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}

func (s *Store) RunMigrations(ctx context.Context) error {
	data, err := os.ReadFile(filepath.Join("internal", "web", "storage", "migrations", "001_init.sql"))
	if err != nil {
		return fmt.Errorf("read migrations: %w", err)
	}
	if _, err := s.DB.ExecContext(ctx, string(data)); err != nil {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}

func (s *Store) CreateUser(ctx context.Context, user User) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO users (id, email, username, password_hash, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, user.ID, strings.ToLower(user.Email), user.Username, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (s *Store) GetUserByEmailOrUsername(ctx context.Context, login string) (User, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users
		WHERE LOWER(email) = LOWER(?) OR username = ?
	`, login, login)
	return scanUser(row)
}

func (s *Store) GetUserByID(ctx context.Context, userID string) (User, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, email, username, password_hash, created_at, updated_at
		FROM users WHERE id = ?
	`, userID)
	return scanUser(row)
}

func (s *Store) CreateSession(ctx context.Context, session AuthSession) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO auth_sessions (id, user_id, status, expires_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, NOW(), NOW())
	`, session.ID, session.UserID, session.Status, session.ExpiresAt)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return s.Redis.Set(ctx, sessionRedisKey(session.ID), session.UserID, time.Until(session.ExpiresAt)).Err()
}

func (s *Store) GetSession(ctx context.Context, sessionID string) (AuthSession, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, status, expires_at, revoked_at, created_at, updated_at
		FROM auth_sessions WHERE id = ?
	`, sessionID)
	var session AuthSession
	if err := row.Scan(&session.ID, &session.UserID, &session.Status, &session.ExpiresAt, &session.RevokedAt, &session.CreatedAt, &session.UpdatedAt); err != nil {
		return AuthSession{}, err
	}
	return session, nil
}

func (s *Store) RevokeSession(ctx context.Context, sessionID string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE auth_sessions
		SET status = 'revoked', revoked_at = NOW(), updated_at = NOW()
		WHERE id = ?
	`, sessionID)
	if err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return s.Redis.Del(ctx, sessionRedisKey(sessionID)).Err()
}

func (s *Store) CreateSkill(ctx context.Context, skill Skill) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO skills (id, user_id, name, slug, description, content, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`, skill.ID, skill.UserID, skill.Name, skill.Slug, skill.Description, skill.Content, skill.Status)
	return err
}

func (s *Store) ListSkillsByUser(ctx context.Context, userID string) ([]Skill, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, name, slug, description, content, status, created_at, updated_at
		FROM skills WHERE user_id = ? ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSkills(rows)
}

func (s *Store) ListEnabledSkillsByUser(ctx context.Context, userID string) ([]Skill, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, name, slug, description, content, status, created_at, updated_at
		FROM skills WHERE user_id = ? AND status = 'enabled' ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanSkills(rows)
}

func (s *Store) GetSkillByID(ctx context.Context, skillID string) (Skill, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, name, slug, description, content, status, created_at, updated_at
		FROM skills WHERE id = ?
	`, skillID)
	var skill Skill
	if err := row.Scan(&skill.ID, &skill.UserID, &skill.Name, &skill.Slug, &skill.Description, &skill.Content, &skill.Status, &skill.CreatedAt, &skill.UpdatedAt); err != nil {
		return Skill{}, err
	}
	return skill, nil
}

func (s *Store) UpdateSkill(ctx context.Context, skill Skill) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE skills
		SET name = ?, slug = ?, description = ?, content = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`, skill.Name, skill.Slug, skill.Description, skill.Content, skill.Status, skill.ID)
	return err
}

func (s *Store) DeleteSkill(ctx context.Context, skillID string) error {
	_, err := s.DB.ExecContext(ctx, `DELETE FROM skills WHERE id = ?`, skillID)
	return err
}

func (s *Store) CreateConversation(ctx context.Context, conversation Conversation) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO conversations (id, user_id, title, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`, conversation.ID, conversation.UserID, conversation.Title)
	return err
}

func (s *Store) ListConversationsByUser(ctx context.Context, userID string) ([]Conversation, error) {
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, user_id, title, created_at, updated_at
		FROM conversations WHERE user_id = ? ORDER BY updated_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var conversations []Conversation
	for rows.Next() {
		var c Conversation
		if err := rows.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, rows.Err()
}

func (s *Store) GetConversationByID(ctx context.Context, conversationID string) (Conversation, error) {
	row := s.DB.QueryRowContext(ctx, `
		SELECT id, user_id, title, created_at, updated_at
		FROM conversations WHERE id = ?
	`, conversationID)
	var c Conversation
	if err := row.Scan(&c.ID, &c.UserID, &c.Title, &c.CreatedAt, &c.UpdatedAt); err != nil {
		return Conversation{}, err
	}
	return c, nil
}

func (s *Store) TouchConversation(ctx context.Context, conversationID, title string) error {
	_, err := s.DB.ExecContext(ctx, `
		UPDATE conversations SET title = COALESCE(NULLIF(?, ''), title), updated_at = NOW() WHERE id = ?
	`, title, conversationID)
	return err
}

func (s *Store) CreateMessage(ctx context.Context, message Message) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO messages (id, conversation_id, user_id, role, content, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`, message.ID, message.ConversationID, message.UserID, message.Role, message.Content)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]Message, error) {
	if limit <= 0 {
		limit = 100
	}
	rows, err := s.DB.QueryContext(ctx, `
		SELECT id, conversation_id, user_id, role, content, created_at
		FROM messages WHERE conversation_id = ? ORDER BY created_at ASC LIMIT ?
	`, conversationID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var messages []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.ConversationID, &m.UserID, &m.Role, &m.Content, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}

func (s *Store) CreateToolCall(ctx context.Context, tc ToolCall) error {
	_, err := s.DB.ExecContext(ctx, `
		INSERT INTO tool_calls (id, conversation_id, user_id, tool_name, status, summary, created_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW())
	`, tc.ID, tc.ConversationID, tc.UserID, tc.ToolName, tc.Status, tc.Summary)
	return err
}

func (s *Store) SetConversationCache(ctx context.Context, conversationID string, messages []Message) error {
	encoded, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return s.Redis.Set(ctx, conversationCacheKey(conversationID), encoded, 30*time.Minute).Err()
}

func (s *Store) GetConversationCache(ctx context.Context, conversationID string) ([]Message, bool, error) {
	value, err := s.Redis.Get(ctx, conversationCacheKey(conversationID)).Result()
	if err == redis.Nil {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}
	var messages []Message
	if err := json.Unmarshal([]byte(value), &messages); err != nil {
		return nil, false, err
	}
	return messages, true, nil
}

func scanUser(row interface{ Scan(dest ...any) error }) (User, error) {
	var user User
	if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return User{}, err
	}
	return user, nil
}

func scanSkills(rows *sql.Rows) ([]Skill, error) {
	var skills []Skill
	for rows.Next() {
		var skill Skill
		if err := rows.Scan(&skill.ID, &skill.UserID, &skill.Name, &skill.Slug, &skill.Description, &skill.Content, &skill.Status, &skill.CreatedAt, &skill.UpdatedAt); err != nil {
			return nil, err
		}
		skills = append(skills, skill)
	}
	return skills, rows.Err()
}

func sessionRedisKey(sessionID string) string {
	return "session:" + sessionID
}

func conversationCacheKey(conversationID string) string {
	return "conversation-cache:" + conversationID
}

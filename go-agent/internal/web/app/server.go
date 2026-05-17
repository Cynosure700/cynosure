package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"nano_cc/internal/config"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/runtime"
	"nano_cc/internal/web/storage"
)

type serverStore interface {
	HealthCheck(ctx context.Context) error
	RunMigrations(ctx context.Context) error
	ListSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error)
	CreateSkill(ctx context.Context, skill storage.Skill) error
	GetSkillByID(ctx context.Context, skillID string) (storage.Skill, error)
	UpdateSkill(ctx context.Context, skill storage.Skill) error
	DeleteSkill(ctx context.Context, skillID string) error
	ListConversationsByUser(ctx context.Context, userID string) ([]storage.Conversation, error)
	CreateConversation(ctx context.Context, conversation storage.Conversation) error
	GetConversationByID(ctx context.Context, conversationID string) (storage.Conversation, error)
	ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error)
}

type Server struct {
	cfg           config.AppConfig
	store         serverStore
	authService   *auth.Service
	runtime       *runtime.Service
	builtinSkills *sessions.SkillLoader
	mux           *http.ServeMux
}

func NewServer() (*Server, error) {
	cfg, err := config.LoadWebConfig()
	if err != nil {
		return nil, err
	}
	if err := config.EnsureAppLayout(cfg); err != nil {
		return nil, err
	}
	if err := config.ValidateAppLayout(cfg); err != nil {
		return nil, err
	}
	config.InitLLM(cfg.LLM)
	store, err := storage.NewStore(cfg)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := store.HealthCheck(ctx); err != nil {
		return nil, err
	}
	if err := store.RunMigrations(ctx); err != nil {
		return nil, err
	}
	if err := logger.InitFileLoggerUnderWorkspaceRoot(cfg.WorkspaceRoot); err != nil {
		logger.Warn(fmt.Sprintf("failed to init file logger: %v", err))
	} else {
		logger.Info(fmt.Sprintf("LLM logs -> %s", logger.LogFilePath()))
	}
	builtinSkills, err := sessions.LoadBuiltinSkillsFromDir(cfg.BuiltinSkillsDir)
	if err != nil {
		return nil, fmt.Errorf("load builtin skills: %w", err)
	}
	runtimeService := runtime.NewService(store, cfg)
	runtimeService.SetBuiltinSkills(builtinSkills)
	server := &Server{
		cfg:           cfg,
		store:         store,
		authService:   auth.NewService(store, cfg),
		runtime:       runtimeService,
		builtinSkills: builtinSkills,
		mux:           http.NewServeMux(),
	}
	server.routes()
	return server, nil
}

func (s *Server) Run() error {
	return http.ListenAndServe(s.cfg.ServerAddr, s.withCORS(s.mux))
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/auth/register", s.handleRegister)
	s.mux.HandleFunc("/api/auth/login", s.handleLogin)
	s.mux.Handle("/api/auth/logout", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleLogout)))
	s.mux.Handle("/api/me", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleMe)))
	s.mux.Handle("/api/skills", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleSkills)))
	s.mux.Handle("/api/skills/", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleSkillByID)))
	s.mux.Handle("/api/conversations", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleConversations)))
	s.mux.Handle("/api/conversations/", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleConversationByID)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "tools": runtime.RegisteredTools(s.cfg)})
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	var body struct{ Email, Username, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, err)
		return
	}
	user, token, session, err := s.authService.Register(r.Context(), body.Email, body.Username, body.Password)
	if err != nil {
		badRequest(w, err)
		return
	}
	s.authService.SetSessionCookie(w, token, session.ExpiresAt)
	writeJSON(w, http.StatusCreated, map[string]any{"user": user})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	var body struct{ Login, Password string }
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, err)
		return
	}
	user, token, session, err := s.authService.Login(r.Context(), body.Login, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	s.authService.SetSessionCookie(w, token, session.ExpiresAt)
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	_, claims, err := s.authService.GetCurrentUser(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := s.authService.Logout(r.Context(), claims.SessionID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.authService.ClearSessionCookie(w)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	writeJSON(w, http.StatusOK, map[string]any{"user": user})
}

func (s *Server) handleSkills(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		skills, err := s.store.ListSkillsByUser(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"skills": appendBuiltinSkills(skills, s.builtinSkills)})
	case http.MethodPost:
		var body struct{ Name, Description, Content, Status string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		skill := storage.Skill{ID: newID("skill"), UserID: user.ID, Name: body.Name, Slug: slugify(body.Name), Description: body.Description, Content: body.Content, Status: normalizeSkillStatus(body.Status)}
		if err := validateSkill(skill); err != nil {
			badRequest(w, err)
			return
		}
		if err := validateNoBuiltinConflict(skill, s.builtinSkills); err != nil {
			badRequest(w, err)
			return
		}
		if err := s.store.CreateSkill(r.Context(), skill); err != nil {
			badRequest(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"skill": skill})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleSkillByID(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	skillID := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	if builtinSkill, ok := builtinSkillByID(skillID, s.builtinSkills); ok {
		s.handleBuiltinSkillByID(w, r, builtinSkill)
		return
	}
	skill, err := s.store.GetSkillByID(r.Context(), skillID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if skill.UserID != user.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"skill": skill})
	case http.MethodPut:
		var body struct{ Name, Description, Content, Status string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		skill.Name = body.Name
		skill.Slug = slugify(body.Name)
		skill.Description = body.Description
		skill.Content = body.Content
		skill.Status = normalizeSkillStatus(body.Status)
		if err := validateSkill(skill); err != nil {
			badRequest(w, err)
			return
		}
		if err := validateNoBuiltinConflict(skill, s.builtinSkills); err != nil {
			badRequest(w, err)
			return
		}
		if err := s.store.UpdateSkill(r.Context(), skill); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"skill": skill})
	case http.MethodDelete:
		if err := s.store.DeleteSkill(r.Context(), skill.ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"ok": true})
	case http.MethodPatch:
		var body struct{ Status string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		skill.Status = normalizeSkillStatus(body.Status)
		if err := s.store.UpdateSkill(r.Context(), skill); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"skill": skill})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleBuiltinSkillByID(w http.ResponseWriter, r *http.Request, skill storage.Skill) {
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, map[string]any{"skill": skill})
	case http.MethodPut, http.MethodPatch, http.MethodDelete:
		http.Error(w, "builtin skills are read-only", http.StatusForbidden)
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleConversations(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	switch r.Method {
	case http.MethodGet:
		conversations, err := s.store.ListConversationsByUser(r.Context(), user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"conversations": conversations})
	case http.MethodPost:
		var body struct{ Title string }
		_ = json.NewDecoder(r.Body).Decode(&body)
		conversation := storage.Conversation{ID: newID("conv"), UserID: user.ID, Title: defaultConversationTitle(body.Title)}
		if err := s.store.CreateConversation(r.Context(), conversation); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusCreated, map[string]any{"conversation": conversation})
	default:
		methodNotAllowed(w)
	}
}

func (s *Server) handleConversationByID(w http.ResponseWriter, r *http.Request) {
	user, _ := auth.UserFromContext(r.Context())
	path := strings.TrimPrefix(r.URL.Path, "/api/conversations/")
	parts := strings.Split(path, "/")
	conversationID := parts[0]
	conversation, err := s.store.GetConversationByID(r.Context(), conversationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if conversation.UserID != user.ID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	if len(parts) == 1 && r.Method == http.MethodGet {
		messages, err := s.store.ListMessagesByConversation(r.Context(), conversation.ID, 100)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"conversation": conversation, "messages": messages})
		return
	}

	if len(parts) == 2 && parts[1] == "stream" && r.Method == http.MethodPost {
		var body struct{ Content string }
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			badRequest(w, err)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		writer := runtime.SSEWriter{W: w}
		_ = writer.Event("conversation", map[string]any{"id": conversation.ID, "title": conversation.Title})
		assistant, err := s.runtime.RespondToConversation(r.Context(), conversation, user, body.Content, writer)
		if err != nil {
			_ = writer.Event("error", map[string]any{"message": err.Error()})
			return
		}
		_ = writer.Event("done", map[string]any{"message_id": assistant.ID})
		return
	}

	methodNotAllowed(w)
}

func (s *Server) withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", s.cfg.AllowedOrigin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func badRequest(w http.ResponseWriter, err error) { http.Error(w, err.Error(), http.StatusBadRequest) }
func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func slugify(input string) string {
	value := strings.ToLower(strings.TrimSpace(input))
	value = strings.ReplaceAll(value, " ", "-")
	value = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-':
			return r
		default:
			return '-'
		}
	}, value)
	value = strings.Trim(value, "-")
	if value == "" {
		return "skill"
	}
	return value
}

func normalizeSkillStatus(status string) string {
	switch status {
	case "enabled", "disabled", "draft":
		return status
	default:
		return "draft"
	}
}

func validateSkill(skill storage.Skill) error {
	if strings.TrimSpace(skill.Name) == "" {
		return fmt.Errorf("skill name is required")
	}
	if strings.TrimSpace(skill.Content) == "" {
		return fmt.Errorf("skill content is required")
	}
	if strings.TrimSpace(skill.Slug) == "" {
		return fmt.Errorf("skill slug is required")
	}
	return nil
}

func validateNoBuiltinConflict(skill storage.Skill, builtin *sessions.SkillLoader) error {
	if builtin == nil {
		return nil
	}
	if _, exists := builtin.Entries()[skill.Slug]; exists {
		return fmt.Errorf("skill slug conflicts with builtin skill")
	}
	return nil
}

func appendBuiltinSkills(skills []storage.Skill, builtin *sessions.SkillLoader) []storage.Skill {
	entries := builtinSkillEntries(builtin)
	if len(entries) == 0 {
		return skills
	}
	merged := make([]storage.Skill, 0, len(entries)+len(skills))
	merged = append(merged, entries...)
	merged = append(merged, skills...)
	return merged
}

func builtinSkillByID(skillID string, builtin *sessions.SkillLoader) (storage.Skill, bool) {
	for _, skill := range builtinSkillEntries(builtin) {
		if skill.ID == skillID {
			return skill, true
		}
	}
	return storage.Skill{}, false
}

func builtinSkillEntries(builtin *sessions.SkillLoader) []storage.Skill {
	if builtin == nil {
		return nil
	}
	entries := builtin.Entries()
	if len(entries) == 0 {
		return nil
	}
	names := make([]string, 0, len(entries))
	for name := range entries {
		names = append(names, name)
	}
	sort.Strings(names)
	skills := make([]storage.Skill, 0, len(names))
	for _, name := range names {
		entry := entries[name]
		description := ""
		if entry != nil {
			description = entry.Meta["description"]
		}
		content := ""
		if entry != nil {
			content = entry.Body
		}
		skills = append(skills, storage.Skill{
			ID:          builtinSkillID(name),
			Name:        name,
			Slug:        name,
			Description: description,
			Content:     content,
			Status:      "enabled",
			Source:      "builtin",
			ReadOnly:    true,
		})
	}
	return skills
}

func builtinSkillID(name string) string {
	return "builtin:" + name
}

func defaultConversationTitle(title string) string {
	if strings.TrimSpace(title) == "" {
		return "新对话"
	}
	return title
}

func newID(prefix string) string {
	return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
}

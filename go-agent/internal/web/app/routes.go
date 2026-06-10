package app

import (
	"net/http"

	"nano_cc/internal/web/runtime"
)

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/auth/register", s.handleRegister)
	s.mux.HandleFunc("/api/auth/login", s.handleLogin)
	s.mux.Handle("/api/auth/logout", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleLogout)))
	s.mux.Handle("/api/me", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleMe)))
	s.mux.Handle("/api/me/memory", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleUpdateMemoryPreference)))
	s.mux.Handle("/api/skills", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleSkills)))
	s.mux.Handle("/api/skills/", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleSkillByID)))
	s.mux.Handle("/api/conversations", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleConversations)))
	s.mux.Handle("/api/conversations/", s.authService.AuthenticateRequest(http.HandlerFunc(s.handleConversationByID)))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"status": "ok", "tools": runtime.RegisteredTools(s.cfg)})
}

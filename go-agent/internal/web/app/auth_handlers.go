package app

import (
	"encoding/json"
	"net/http"

	"nano_cc/internal/web/auth"
)

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

func (s *Server) handleUpdateMemoryPreference(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		methodNotAllowed(w)
		return
	}
	user, _ := auth.UserFromContext(r.Context())
	var body struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		badRequest(w, err)
		return
	}
	if err := s.store.UpdateUserMemoryEnabled(r.Context(), user.ID, body.Enabled); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"memory_enabled": body.Enabled})
}

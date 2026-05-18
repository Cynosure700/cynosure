package app

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/storage"
)

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

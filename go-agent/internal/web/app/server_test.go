package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"nano_cc/internal/sessions"
	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/storage"
)

type fakeServerStore struct {
	skillToReturn storage.Skill
	skills        []storage.Skill
	createCalled  bool
	updateCalled  bool
	deleteCalled  bool
	getCalled     bool
}

func (f *fakeServerStore) HealthCheck(ctx context.Context) error { return nil }

func (f *fakeServerStore) RunMigrations(ctx context.Context) error { return nil }

func (f *fakeServerStore) ListSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error) {
	return f.skills, nil
}

func (f *fakeServerStore) CreateSkill(ctx context.Context, skill storage.Skill) error {
	f.createCalled = true
	return nil
}

func (f *fakeServerStore) GetSkillByID(ctx context.Context, skillID string) (storage.Skill, error) {
	f.getCalled = true
	if f.skillToReturn.ID == "" {
		return storage.Skill{}, sql.ErrNoRows
	}
	return f.skillToReturn, nil
}

func (f *fakeServerStore) UpdateSkill(ctx context.Context, skill storage.Skill) error {
	f.updateCalled = true
	return nil
}

func (f *fakeServerStore) DeleteSkill(ctx context.Context, skillID string) error {
	f.deleteCalled = true
	return nil
}

func (f *fakeServerStore) ListConversationsByUser(ctx context.Context, userID string) ([]storage.Conversation, error) {
	return nil, nil
}

func (f *fakeServerStore) CreateConversation(ctx context.Context, conversation storage.Conversation) error {
	return nil
}

func (f *fakeServerStore) GetConversationByID(ctx context.Context, conversationID string) (storage.Conversation, error) {
	return storage.Conversation{}, sql.ErrNoRows
}

func (f *fakeServerStore) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error) {
	return nil, nil
}

func TestHandleSkills_RejectsBuiltinSlugConflictOnCreate(t *testing.T) {
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	store := &fakeServerStore{}
	server := &Server{store: store, builtinSkills: builtin}

	req := httptest.NewRequest(http.MethodPost, "/api/skills", strings.NewReader(`{"name":"Builtin Skill","description":"desc","content":"content"}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleSkills(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "skill slug conflicts with builtin skill") {
		t.Fatalf("expected builtin conflict error, got %q", resp.Body.String())
	}
	if store.createCalled {
		t.Fatalf("expected create skill not to be called on builtin conflict")
	}
}

func TestHandleSkillByID_RejectsBuiltinSlugConflictOnUpdate(t *testing.T) {
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	store := &fakeServerStore{skillToReturn: storage.Skill{ID: "skill_1", UserID: "usr_1", Name: "Original Skill", Slug: "original-skill", Content: "content", Status: "draft"}}
	server := &Server{store: store, builtinSkills: builtin}

	req := httptest.NewRequest(http.MethodPut, "/api/skills/skill_1", strings.NewReader(`{"name":"Builtin Skill","description":"desc","content":"updated content"}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleSkillByID(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
	if !strings.Contains(resp.Body.String(), "skill slug conflicts with builtin skill") {
		t.Fatalf("expected builtin conflict error, got %q", resp.Body.String())
	}
	if store.updateCalled {
		t.Fatalf("expected update skill not to be called on builtin conflict")
	}
}

func TestHandleSkills_IncludesBuiltinSkillsAsReadOnlyEntries(t *testing.T) {
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	store := &fakeServerStore{skills: []storage.Skill{{ID: "skill_1", UserID: "usr_1", Name: "Custom Skill", Slug: "custom-skill", Description: "Custom", Content: "custom body", Status: "enabled"}}}
	server := &Server{store: store, builtinSkills: builtin}

	req := httptest.NewRequest(http.MethodGet, "/api/skills", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleSkills(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	var body struct {
		Skills []storage.Skill `json:"skills"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Skills) != 2 {
		t.Fatalf("expected builtin and custom skills, got %#v", body.Skills)
	}
	builtinSkill := body.Skills[0]
	if builtinSkill.ID != builtinSkillID("builtin-skill") || builtinSkill.Source != "builtin" || !builtinSkill.ReadOnly {
		t.Fatalf("expected builtin read-only skill entry, got %#v", builtinSkill)
	}
	if builtinSkill.Status != "enabled" {
		t.Fatalf("expected builtin skill status enabled, got %#v", builtinSkill)
	}
	if body.Skills[1].ID != "skill_1" {
		t.Fatalf("expected custom skill to remain in list, got %#v", body.Skills[1])
	}
	if body.Skills[1].ReadOnly {
		t.Fatalf("expected custom skill to remain writable, got %#v", body.Skills[1])
	}
}

func TestHandleSkillByID_ReturnsBuiltinSkillWithoutStoreLookup(t *testing.T) {
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})
	store := &fakeServerStore{}
	server := &Server{store: store, builtinSkills: builtin}

	req := httptest.NewRequest(http.MethodGet, "/api/skills/"+builtinSkillID("builtin-skill"), nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleSkillByID(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if store.getCalled {
		t.Fatalf("expected builtin skill lookup to bypass database store")
	}
	var body struct {
		Skill storage.Skill `json:"skill"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Skill.ID != builtinSkillID("builtin-skill") || !body.Skill.ReadOnly || body.Skill.Source != "builtin" {
		t.Fatalf("expected builtin read-only skill payload, got %#v", body.Skill)
	}
	if body.Skill.Content != "builtin body" {
		t.Fatalf("expected builtin skill content, got %#v", body.Skill)
	}
}

func TestHandleSkillByID_RejectsBuiltinMutationRequests(t *testing.T) {
	builtin := sessions.NewSkillLoader()
	builtin.LoadFromEntries(map[string]*sessions.SkillEntry{
		"builtin-skill": {Meta: map[string]string{"description": "Builtin"}, Body: "builtin body", Path: "builtin://builtin-skill"},
	})

	methods := []struct {
		name   string
		method string
		body   string
	}{
		{name: "put", method: http.MethodPut, body: `{"name":"Updated Builtin","description":"desc","content":"updated"}`},
		{name: "patch", method: http.MethodPatch, body: `{"status":"disabled"}`},
		{name: "delete", method: http.MethodDelete},
	}

	for _, tc := range methods {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeServerStore{}
			server := &Server{store: store, builtinSkills: builtin}

			req := httptest.NewRequest(tc.method, "/api/skills/"+builtinSkillID("builtin-skill"), strings.NewReader(tc.body))
			req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
			resp := httptest.NewRecorder()

			server.handleSkillByID(resp, req)

			if resp.Code != http.StatusForbidden {
				t.Fatalf("expected status %d, got %d", http.StatusForbidden, resp.Code)
			}
			if !strings.Contains(resp.Body.String(), "builtin skills are read-only") {
				t.Fatalf("expected readonly error, got %q", resp.Body.String())
			}
			if store.getCalled || store.updateCalled || store.deleteCalled {
				t.Fatalf("expected builtin mutation to avoid store calls, got %+v", store)
			}
		})
	}
}

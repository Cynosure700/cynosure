package app

import (
	"context"
	"database/sql"
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
	createCalled  bool
	updateCalled  bool
}

func (f *fakeServerStore) HealthCheck(ctx context.Context) error { return nil }

func (f *fakeServerStore) RunMigrations(ctx context.Context) error { return nil }

func (f *fakeServerStore) ListSkillsByUser(ctx context.Context, userID string) ([]storage.Skill, error) {
	return nil, nil
}

func (f *fakeServerStore) CreateSkill(ctx context.Context, skill storage.Skill) error {
	f.createCalled = true
	return nil
}

func (f *fakeServerStore) GetSkillByID(ctx context.Context, skillID string) (storage.Skill, error) {
	if f.skillToReturn.ID == "" {
		return storage.Skill{}, sql.ErrNoRows
	}
	return f.skillToReturn, nil
}

func (f *fakeServerStore) UpdateSkill(ctx context.Context, skill storage.Skill) error {
	f.updateCalled = true
	return nil
}

func (f *fakeServerStore) DeleteSkill(ctx context.Context, skillID string) error { return nil }

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

package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"nano_cc/internal/sessions"
	"nano_cc/internal/web/auth"
	"nano_cc/internal/web/storage"
)

type fakeServerStore struct {
	skillToReturn        storage.Skill
	conversationToReturn storage.Conversation
	createdConversation  storage.Conversation
	messages             []storage.Message
	skills               []storage.Skill
	createCalled         bool
	updateCalled         bool
	deleteCalled         bool
	createConversation   bool
	deleteConversation   bool
	updateConversation   bool
	getCalled            bool
}

func (f *fakeServerStore) HealthCheck(ctx context.Context) error { return nil }

func (f *fakeServerStore) RunMigrations(ctx context.Context) error { return nil }

func (f *fakeServerStore) UpdateUserMemoryEnabled(ctx context.Context, userID string, enabled bool) error {
	return nil
}

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

func (f *fakeServerStore) ListMCPServersByUser(ctx context.Context, userID string) ([]storage.MCPServer, error) {
	return nil, nil
}

func (f *fakeServerStore) CreateMCPServer(ctx context.Context, server storage.MCPServer) error {
	return nil
}

func (f *fakeServerStore) GetMCPServerByID(ctx context.Context, id string) (storage.MCPServer, error) {
	return storage.MCPServer{}, nil
}

func (f *fakeServerStore) UpdateMCPServer(ctx context.Context, server storage.MCPServer) error {
	return nil
}

func (f *fakeServerStore) DeleteMCPServer(ctx context.Context, id string) error {
	return nil
}

func (f *fakeServerStore) ListConversationsByUser(ctx context.Context, userID string) ([]storage.Conversation, error) {
	return nil, nil
}

func (f *fakeServerStore) CreateConversation(ctx context.Context, conversation storage.Conversation) error {
	f.createConversation = true
	f.createdConversation = conversation
	return nil
}

func (f *fakeServerStore) GetConversationByID(ctx context.Context, conversationID string) (storage.Conversation, error) {
	if f.conversationToReturn.ID == "" {
		return storage.Conversation{}, sql.ErrNoRows
	}
	return f.conversationToReturn, nil
}

func (f *fakeServerStore) UpdateConversationTitle(ctx context.Context, conversationID, title string) error {
	f.updateConversation = true
	if f.conversationToReturn.ID == conversationID {
		f.conversationToReturn.Title = title
	}
	return nil
}

func (f *fakeServerStore) DeleteConversation(ctx context.Context, conversationID string) error {
	f.deleteConversation = true
	if f.conversationToReturn.ID == conversationID {
		f.conversationToReturn = storage.Conversation{}
	}
	return nil
}

func (f *fakeServerStore) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error) {
	return f.messages, nil
}

func loadWorkspaceBuiltinSkillsForTest(t *testing.T) *sessions.SkillLoader {
	t.Helper()
	workspaceRoot := t.TempDir()
	skillDir := filepath.Join(workspaceRoot, "skills", "builtin-skill")
	if err := os.MkdirAll(skillDir, 0o755); err != nil {
		t.Fatalf("mkdir skill dir: %v", err)
	}
	skillDoc := `---
name: builtin-skill
description: Builtin from workspace
---

Builtin body.`
	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(skillDoc), 0o644); err != nil {
		t.Fatalf("write skill doc: %v", err)
	}
	loader, err := sessions.LoadBuiltinSkillsFromDir(filepath.Join(workspaceRoot, "skills"))
	if err != nil {
		t.Fatalf("load builtin skills: %v", err)
	}
	return loader
}

func TestHandleSkills_RejectsBuiltinSlugConflictOnCreate(t *testing.T) {
	builtin := loadWorkspaceBuiltinSkillsForTest(t)
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
	builtin := loadWorkspaceBuiltinSkillsForTest(t)
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

func TestHandleSkills_ListsOnlyUserOwnedSkills(t *testing.T) {
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
	if len(body.Skills) != 1 {
		t.Fatalf("expected only custom skills, got %#v", body.Skills)
	}
	if body.Skills[0].ID != "skill_1" {
		t.Fatalf("expected custom skill to remain in list, got %#v", body.Skills[0])
	}
	if body.Skills[0].ReadOnly {
		t.Fatalf("expected custom skill to remain writable, got %#v", body.Skills[0])
	}
	if body.Skills[0].Source != "" {
		t.Fatalf("expected no builtin source marker on custom skill, got %#v", body.Skills[0])
	}
}

func TestHandleSkillByID_HidesBuiltinSkillDetails(t *testing.T) {
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

	if resp.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.Code)
	}
	if store.getCalled {
		t.Fatalf("expected builtin skill lookup to bypass database store")
	}
}

func TestHandleSkillByID_TreatsBuiltinMutationRequestsAsNotFound(t *testing.T) {
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

			if resp.Code != http.StatusNotFound {
				t.Fatalf("expected status %d, got %d", http.StatusNotFound, resp.Code)
			}
			if store.getCalled || store.updateCalled || store.deleteCalled {
				t.Fatalf("expected builtin mutation to avoid store calls, got %+v", store)
			}
		})
	}
}

func TestHandleConversationByID_DeletesOwnedConversation(t *testing.T) {
	store := &fakeServerStore{conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_1", Title: "test"}}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodDelete, "/api/conversations/conv_1", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversationByID(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if !store.deleteConversation {
		t.Fatalf("expected delete conversation to be called")
	}
	var body struct {
		OK bool `json:"ok"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.OK {
		t.Fatalf("expected ok response, got %#v", body)
	}
}

func TestHandleConversationByID_RejectsDeletingOtherUsersConversation(t *testing.T) {
	store := &fakeServerStore{conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_other", Title: "test"}}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodDelete, "/api/conversations/conv_1", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversationByID(resp, req)

	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, resp.Code)
	}
	if store.deleteConversation {
		t.Fatalf("expected delete conversation not to be called")
	}
}

func TestHandleConversationByID_ReturnsNotFoundAfterDelete(t *testing.T) {
	store := &fakeServerStore{conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_1", Title: "test"}}
	server := &Server{store: store}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/conversations/conv_1", nil)
	deleteReq = deleteReq.WithContext(context.WithValue(deleteReq.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	deleteResp := httptest.NewRecorder()
	server.handleConversationByID(deleteResp, deleteReq)

	if deleteResp.Code != http.StatusOK {
		t.Fatalf("expected delete status %d, got %d", http.StatusOK, deleteResp.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/conversations/conv_1", nil)
	getReq = getReq.WithContext(context.WithValue(getReq.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	getResp := httptest.NewRecorder()
	server.handleConversationByID(getResp, getReq)

	if getResp.Code != http.StatusNotFound {
		t.Fatalf("expected status %d after delete, got %d", http.StatusNotFound, getResp.Code)
	}
}

func TestHandleConversationByID_ReturnsOnlyDisplayMessages(t *testing.T) {
	store := &fakeServerStore{
		conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_1", Title: "test"},
		messages: []storage.Message{
			{ID: "msg_user", ConversationID: "conv_1", UserID: "usr_1", Role: "user", Content: "帮我加载技能"},
			{ID: "msg_assistant_tool", ConversationID: "conv_1", UserID: "usr_1", Role: "assistant", ToolCalls: []storage.MessageToolCall{{ID: "tool_1", Type: "function", Function: storage.MessageFunctionCall{Name: "load_skill", Arguments: `{"name":"builtin-skill"}`}}}},
			{ID: "msg_tool", ConversationID: "conv_1", UserID: "usr_1", Role: "tool", ToolCallID: "tool_1", Content: `{"status":"success","result":"loaded"}`},
			{ID: "msg_assistant", ConversationID: "conv_1", UserID: "usr_1", Role: "assistant", Content: "已经加载完成"},
		},
	}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodGet, "/api/conversations/conv_1", nil)
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversationByID(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	var body struct {
		Messages []storage.Message `json:"messages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(body.Messages) != 2 {
		t.Fatalf("expected only user and final assistant messages, got %#v", body.Messages)
	}
	if body.Messages[0].Role != "user" || body.Messages[1].Role != "assistant" || body.Messages[1].Content != "已经加载完成" {
		t.Fatalf("unexpected display messages: %#v", body.Messages)
	}
}

func TestHandleConversationByID_RenamesOwnedConversation(t *testing.T) {
	store := &fakeServerStore{conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_1", Title: "旧标题"}}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodPatch, "/api/conversations/conv_1", strings.NewReader(`{"title":"  新标题  "}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversationByID(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.Code)
	}
	if !store.updateConversation {
		t.Fatalf("expected update conversation title to be called")
	}
	var body struct {
		Conversation storage.Conversation `json:"conversation"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Conversation.Title != "新标题" {
		t.Fatalf("expected trimmed title, got %#v", body.Conversation)
	}
}

func TestHandleConversationByID_RejectsEmptyRename(t *testing.T) {
	store := &fakeServerStore{conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_1", Title: "旧标题"}}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodPatch, "/api/conversations/conv_1", strings.NewReader(`{"title":"   "}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversationByID(resp, req)

	if resp.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
	}
	if store.updateConversation {
		t.Fatalf("expected rename not to reach store")
	}
}

func TestHandleConversationByID_RejectsRenamingOtherUsersConversation(t *testing.T) {
	store := &fakeServerStore{conversationToReturn: storage.Conversation{ID: "conv_1", UserID: "usr_other", Title: "test"}}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodPatch, "/api/conversations/conv_1", strings.NewReader(`{"title":"新标题"}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversationByID(resp, req)

	if resp.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, resp.Code)
	}
	if store.updateConversation {
		t.Fatalf("expected rename not to reach store")
	}
}

func TestHandleConversations_CreatesConversationWithCustomTitle(t *testing.T) {
	store := &fakeServerStore{}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodPost, "/api/conversations", strings.NewReader(`{"title":"项目排期"}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversations(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.Code)
	}
	if !store.createConversation {
		t.Fatalf("expected create conversation to be called")
	}
	if store.createdConversation.UserID != "usr_1" || store.createdConversation.Title != "项目排期" {
		t.Fatalf("unexpected stored conversation: %#v", store.createdConversation)
	}
	if !strings.HasPrefix(store.createdConversation.ID, "conv_") {
		t.Fatalf("expected generated conversation id, got %#v", store.createdConversation)
	}
	if !strings.HasPrefix(store.createdConversation.RootMessageID, "msg_") {
		t.Fatalf("expected generated root message id, got %#v", store.createdConversation)
	}
}

func TestHandleConversations_UsesDefaultTitleWhenEmpty(t *testing.T) {
	store := &fakeServerStore{}
	server := &Server{store: store}

	req := httptest.NewRequest(http.MethodPost, "/api/conversations", strings.NewReader(`{"title":"   "}`))
	req = req.WithContext(context.WithValue(req.Context(), auth.UserContextKey, storage.User{ID: "usr_1"}))
	resp := httptest.NewRecorder()

	server.handleConversations(resp, req)

	if resp.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, resp.Code)
	}
	if !store.createConversation {
		t.Fatalf("expected create conversation to be called")
	}
	if store.createdConversation.Title != "新对话" {
		t.Fatalf("expected default title, got %#v", store.createdConversation)
	}
}

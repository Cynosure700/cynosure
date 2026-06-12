package local

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"nano_cc/internal/agent/storage"
)

type Store struct {
	mu              sync.RWMutex
	conversations   map[string]storage.Conversation
	messages        map[string][]storage.Message
	cache           map[string][]storage.Message
	toolCalls       []storage.ToolCall
	persisted       map[string]storage.PersistedOutput
	persistedByHash map[string]string
	summaries       map[string]storage.ContextSummary
	modelHistory    map[string][]storage.Message
	memories        []storage.Memory
	convMemories    map[string][]storage.ConversationMemory
	memory          *MarkdownMemoryStore
	sessionHistory  *SessionHistoryStore
	workspaceRoot   string
	locks           map[string]string
}

func NewStore() *Store {
	return &Store{
		conversations:   make(map[string]storage.Conversation),
		messages:        make(map[string][]storage.Message),
		cache:           make(map[string][]storage.Message),
		persisted:       make(map[string]storage.PersistedOutput),
		persistedByHash: make(map[string]string),
		summaries:       make(map[string]storage.ContextSummary),
		modelHistory:    make(map[string][]storage.Message),
		convMemories:    make(map[string][]storage.ConversationMemory),
		locks:           make(map[string]string),
	}
}

func NewStoreWithMemory(workspaceRoot string) (*Store, error) {
	store := NewStore()
	memory, err := NewMarkdownMemoryStore(workspaceRoot)
	if err != nil {
		return nil, err
	}
	sessionHistory, err := NewSessionHistoryStore()
	if err != nil {
		return nil, err
	}
	store.memory = memory
	store.sessionHistory = sessionHistory
	store.workspaceRoot = workspaceRoot
	return store, nil
}

func (s *Store) CreateConversation(ctx context.Context, conversation storage.Conversation) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if conversation.SessionID == "" {
		conversation.SessionID = conversation.ID
	}
	now := time.Now()
	if conversation.CreatedAt.IsZero() {
		conversation.CreatedAt = now
	}
	conversation.UpdatedAt = now
	s.conversations[conversation.ID] = conversation
	return nil
}

func (s *Store) GetConversationByID(ctx context.Context, conversationID string) (storage.Conversation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	conv, ok := s.conversations[conversationID]
	if !ok {
		return storage.Conversation{}, sql.ErrNoRows
	}
	return conv, nil
}

func (s *Store) UpdateConversationTitle(ctx context.Context, conversationID, title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	conv, ok := s.conversations[conversationID]
	if !ok {
		return sql.ErrNoRows
	}
	conv.Title = title
	conv.UpdatedAt = time.Now()
	s.conversations[conversationID] = conv
	return nil
}

func (s *Store) TouchConversationActivity(ctx context.Context, conversationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	conv, ok := s.conversations[conversationID]
	if !ok {
		return nil
	}
	conv.UpdatedAt = time.Now()
	s.conversations[conversationID] = conv
	return nil
}

func (s *Store) SetConversationHistory(ctx context.Context, conversationID string, messages []storage.Message) error {
	s.mu.Lock()
	conv := s.conversations[conversationID]
	s.messages[conversationID] = cloneMessages(messages)
	s.cache[conversationID] = cloneMessages(messages)
	s.mu.Unlock()
	if s.sessionHistory != nil && conv.ID != "" {
		return s.sessionHistory.SaveHistory(ctx, s.workspaceRoot, conv, messages)
	}
	return nil
}

func (s *Store) SetConversationCache(ctx context.Context, conversationID string, messages []storage.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.cache[conversationID] = cloneMessages(messages)
	return nil
}

func (s *Store) GetConversationCache(ctx context.Context, conversationID string) ([]storage.Message, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages, ok := s.cache[conversationID]
	return cloneMessages(messages), ok, nil
}

func (s *Store) ListMessagesByConversation(ctx context.Context, conversationID string, limit int) ([]storage.Message, error) {
	s.mu.RLock()
	conv := s.conversations[conversationID]
	messages := cloneMessages(s.messages[conversationID])
	s.mu.RUnlock()
	if len(messages) == 0 && s.sessionHistory != nil && conv.SessionID != "" {
		doc, err := s.sessionHistory.LoadHistory(ctx, conv.SessionID)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if err == nil {
			messages = cloneMessages(doc.Messages)
			s.mu.Lock()
			s.messages[conversationID] = cloneMessages(messages)
			s.cache[conversationID] = cloneMessages(messages)
			s.mu.Unlock()
		}
	}
	if limit > 0 && len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}
	return messages, nil
}

func (s *Store) CreateToolCall(ctx context.Context, tc storage.ToolCall) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if tc.CreatedAt.IsZero() {
		tc.CreatedAt = time.Now()
	}
	s.toolCalls = append(s.toolCalls, tc)
	return nil
}

func (s *Store) CreateSubagentMessage(ctx context.Context, message storage.SubagentMessage) error {
	return nil
}

func (s *Store) CreatePersistedOutput(ctx context.Context, output storage.PersistedOutput) error {
	s.mu.RLock()
	conv := s.conversations[output.ConversationID]
	workspaceRoot := s.workspaceRoot
	s.mu.RUnlock()
	if workspaceRoot != "" && conv.SessionID != "" {
		persisted, err := persistOutputToWorkspace(ctx, workspaceRoot, conv.SessionID, output)
		if err != nil {
			return err
		}
		output = persisted
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if output.CreatedAt.IsZero() {
		output.CreatedAt = time.Now()
	}
	s.persisted[output.ID] = output
	s.persistedByHash[persistedHashKey(output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)] = output.ID
	return nil
}

func (s *Store) GetPersistedOutputForConversation(ctx context.Context, id, userID, conversationID string) (storage.PersistedOutput, error) {
	s.mu.RLock()
	output, ok := s.persisted[id]
	if ok && output.UserID == userID && output.ConversationID == conversationID {
		s.mu.RUnlock()
		return output, nil
	}
	conv := s.conversations[conversationID]
	workspaceRoot := s.workspaceRoot
	s.mu.RUnlock()
	if workspaceRoot == "" || conv.SessionID == "" {
		return storage.PersistedOutput{}, sql.ErrNoRows
	}
	output, err := loadPersistedOutputFromWorkspace(ctx, workspaceRoot, conv.SessionID, id)
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	if output.UserID != userID || output.ConversationID != conversationID {
		return storage.PersistedOutput{}, sql.ErrNoRows
	}
	s.mu.Lock()
	s.persisted[output.ID] = output
	s.persistedByHash[persistedHashKey(output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)] = output.ID
	s.mu.Unlock()
	return output, nil
}

func (s *Store) GetPersistedOutputByMessageHash(ctx context.Context, conversationID, userID, messageID, toolCallID, strategy, contentSHA256 string) (storage.PersistedOutput, error) {
	s.mu.RLock()
	id, ok := s.persistedByHash[persistedHashKey(conversationID, userID, messageID, toolCallID, strategy, contentSHA256)]
	if ok {
		output := s.persisted[id]
		s.mu.RUnlock()
		return output, nil
	}
	conv := s.conversations[conversationID]
	workspaceRoot := s.workspaceRoot
	s.mu.RUnlock()
	if workspaceRoot == "" || conv.SessionID == "" {
		return storage.PersistedOutput{}, sql.ErrNoRows
	}
	output, err := findPersistedOutputByMessageHashInWorkspace(ctx, workspaceRoot, conv.SessionID, conversationID, userID, messageID, toolCallID, strategy, contentSHA256)
	if err != nil {
		return storage.PersistedOutput{}, err
	}
	s.mu.Lock()
	s.persisted[output.ID] = output
	s.persistedByHash[persistedHashKey(output.ConversationID, output.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)] = output.ID
	s.mu.Unlock()
	return output, nil
}

func (s *Store) CreateContextSummary(ctx context.Context, summary storage.ContextSummary) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.summaries[summary.ConversationID+"\x00"+summary.UserID+"\x00"+summary.SourceHistorySHA256] = summary
	return nil
}

func (s *Store) GetContextSummaryByHistoryHash(ctx context.Context, conversationID, userID, sourceHistorySHA256 string) (storage.ContextSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	summary, ok := s.summaries[conversationID+"\x00"+userID+"\x00"+sourceHistorySHA256]
	if !ok {
		return storage.ContextSummary{}, sql.ErrNoRows
	}
	return summary, nil
}

func (s *Store) ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	if s.memory != nil {
		return s.memory.ListRelevantMemories(ctx, userID)
	}
	return nil, nil
}

func (s *Store) ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error) {
	if s.memory != nil {
		return s.memory.ListMemoriesByUserAndType(ctx, userID, memType)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []storage.Memory
	for _, memory := range s.memories {
		if memory.UserID == userID && memory.Type == memType {
			result = append(result, memory)
		}
	}
	return result, nil
}

func (s *Store) ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	if s.memory != nil {
		return s.memory.ListProjectFactMemories(ctx, userID)
	}
	return s.ListMemoriesByUserAndType(ctx, userID, "project_fact")
}

func (s *Store) InsertMemory(ctx context.Context, m storage.Memory) error {
	if s.memory != nil {
		return s.memory.InsertMemory(ctx, m)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memories = append(s.memories, m)
	return nil
}

func (s *Store) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	if s.memory != nil {
		return s.memory.CountMemoriesByUserAndType(ctx, userID, memType)
	}
	items, err := s.ListMemoriesByUserAndType(ctx, userID, memType)
	return len(items), err
}

func (s *Store) CountProjectFactMemories(ctx context.Context, userID string) (int, error) {
	if s.memory != nil {
		return s.memory.CountProjectFactMemories(ctx, userID)
	}
	items, err := s.ListProjectFactMemories(ctx, userID)
	return len(items), err
}
func (s *Store) DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error {
	if s.memory != nil {
		return s.memory.DeleteOldestMemories(ctx, userID, memType, n)
	}
	return nil
}
func (s *Store) ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error {
	if s.memory != nil {
		return s.memory.ReplaceMemoriesByUserAndType(ctx, userID, memType, items)
	}
	return nil
}
func (s *Store) ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error {
	if s.memory != nil {
		return s.memory.ReplaceProjectFactMemories(ctx, userID, items)
	}
	return nil
}

func (s *Store) ListConversationMemories(ctx context.Context, conversationID string) ([]storage.ConversationMemory, error) {
	if s.memory != nil {
		sessionID := s.sessionIDForConversation(conversationID)
		return s.memory.ListConversationMemories(ctx, sessionID)
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	items := append([]storage.ConversationMemory(nil), s.convMemories[conversationID]...)
	sort.Slice(items, func(i, j int) bool { return items[i].CreatedAt.Before(items[j].CreatedAt) })
	return items, nil
}

func (s *Store) ReplaceConversationMemories(ctx context.Context, conversationID, userID string, items []storage.ConversationMemory) error {
	if s.memory != nil {
		sessionID := s.sessionIDForConversation(conversationID)
		return s.memory.ReplaceConversationMemories(ctx, sessionID, userID, items)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.convMemories[conversationID] = append([]storage.ConversationMemory(nil), items...)
	return nil
}

func (s *Store) GetConversationModelHistory(ctx context.Context, conversationID string) ([]storage.Message, bool, error) {
	s.mu.RLock()
	conv := s.conversations[conversationID]
	messages, ok := s.modelHistory[conversationID]
	s.mu.RUnlock()
	if (!ok || len(messages) == 0) && s.sessionHistory != nil && conv.SessionID != "" {
		doc, err := s.sessionHistory.LoadModelHistory(ctx, conv.SessionID)
		if err != nil {
			if os.IsNotExist(err) {
				return nil, false, nil
			}
			return nil, false, err
		}
		messages = doc.Messages
		ok = len(messages) > 0
		s.mu.Lock()
		s.modelHistory[conversationID] = cloneMessages(messages)
		s.mu.Unlock()
	}
	return cloneMessages(messages), ok, nil
}

func (s *Store) UpsertConversationModelHistory(ctx context.Context, conversationID, userID string, messages []storage.Message) error {
	s.mu.Lock()
	conv := s.conversations[conversationID]
	s.modelHistory[conversationID] = cloneMessages(messages)
	s.mu.Unlock()
	if s.sessionHistory != nil && conv.ID != "" {
		if conv.UserID == "" {
			conv.UserID = userID
		}
		return s.sessionHistory.SaveModelHistory(ctx, s.workspaceRoot, conv, messages)
	}
	return nil
}

func (s *Store) ListResumableSessions(ctx context.Context, currentWorkspace string) ([]storage.ResumableSession, error) {
	if s.sessionHistory == nil {
		return nil, nil
	}
	return s.sessionHistory.ListSessions(ctx, currentWorkspace)
}

func (s *Store) ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error) {
	if s.sessionHistory == nil {
		return storage.Conversation{}, nil, fmt.Errorf("session history store is not initialized")
	}
	doc, err := s.sessionHistory.LoadHistory(ctx, sessionID)
	if err != nil {
		return storage.Conversation{}, nil, err
	}
	current, err := normalizeWorkspacePath(currentWorkspace)
	if err != nil {
		return storage.Conversation{}, nil, err
	}
	stored, err := normalizeWorkspacePath(doc.WorkspaceRoot)
	if err != nil {
		return storage.Conversation{}, nil, err
	}
	if current != stored {
		return storage.Conversation{}, nil, fmt.Errorf("历史会话属于 %s，请 cd 到该目录或使用 --cwd %s 后再执行 /resume", doc.WorkspaceRoot, doc.WorkspaceRoot)
	}
	conv := storage.Conversation{ID: doc.ConversationID, SessionID: doc.SessionID, UserID: user.ID, Title: doc.Title, CreatedAt: doc.CreatedAt, UpdatedAt: doc.UpdatedAt}
	if conv.ID == "" {
		conv.ID = doc.SessionID
	}
	if conv.UserID == "" {
		conv.UserID = doc.UserID
	}
	if conv.UserID == "" {
		conv.UserID = LocalUserID
	}
	history := cloneMessages(doc.Messages)
	var modelHistory []storage.Message
	if modelDoc, err := s.sessionHistory.LoadModelHistory(ctx, sessionID); err == nil {
		modelHistory = cloneMessages(modelDoc.Messages)
	} else if !os.IsNotExist(err) {
		return storage.Conversation{}, nil, err
	}
	s.mu.Lock()
	s.conversations[conv.ID] = conv
	s.messages[conv.ID] = cloneMessages(history)
	s.cache[conv.ID] = cloneMessages(history)
	if len(modelHistory) > 0 {
		s.modelHistory[conv.ID] = cloneMessages(modelHistory)
	} else {
		delete(s.modelHistory, conv.ID)
	}
	s.mu.Unlock()
	return conv, history, nil
}

func (s *Store) AcquireConversationLock(ctx context.Context, conversationID, token string, ttl, waitTimeout time.Duration) (bool, error) {
	key := conversationID
	if s.memory != nil {
		key = s.memory.ProjectLockKey(s.sessionIDForConversation(conversationID))
	}
	deadline := time.Now().Add(waitTimeout)
	for {
		s.mu.Lock()
		if s.locks[key] == "" {
			s.locks[key] = token
			s.mu.Unlock()
			return true, nil
		}
		s.mu.Unlock()
		if waitTimeout <= 0 || time.Now().After(deadline) {
			return false, nil
		}
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(10 * time.Millisecond):
		}
	}
}
func (s *Store) RenewConversationLock(ctx context.Context, conversationID, token string, ttl time.Duration) (bool, error) {
	key := conversationID
	if s.memory != nil {
		key = s.memory.ProjectLockKey(s.sessionIDForConversation(conversationID))
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.locks[key] == token, nil
}
func (s *Store) ReleaseConversationLock(ctx context.Context, conversationID, token string) error {
	key := conversationID
	if s.memory != nil {
		key = s.memory.ProjectLockKey(s.sessionIDForConversation(conversationID))
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.locks[key] == token {
		delete(s.locks, key)
	}
	return nil
}

func (s *Store) sessionIDForConversation(conversationID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if conv, ok := s.conversations[conversationID]; ok && conv.SessionID != "" {
		return conv.SessionID
	}
	return conversationID
}

func cloneMessages(messages []storage.Message) []storage.Message {
	if len(messages) == 0 {
		return nil
	}
	cloned := make([]storage.Message, len(messages))
	copy(cloned, messages)
	for i := range cloned {
		cloned[i].ToolCalls = append([]storage.MessageToolCall(nil), cloned[i].ToolCalls...)
	}
	return cloned
}

func persistedHashKey(parts ...string) string {
	key := ""
	for _, part := range parts {
		key += part + "\x00"
	}
	return key
}

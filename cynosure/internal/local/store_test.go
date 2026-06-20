package local

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
)

func TestStoreMaintainsConversationHistory(t *testing.T) {
	store := NewStore()
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", UserID: "local", Title: "新对话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	messages := []storage.Message{{ID: "msg_1", ConversationID: conv.ID, UserID: conv.UserID, Role: "user", Content: "hello"}}
	if err := store.SetConversationHistory(ctx, conv.ID, messages); err != nil {
		t.Fatalf("SetConversationHistory returned error: %v", err)
	}

	got, err := store.ListMessagesByConversation(ctx, conv.ID, 100)
	if err != nil {
		t.Fatalf("ListMessagesByConversation returned error: %v", err)
	}
	if len(got) != 1 || got[0].Content != "hello" {
		t.Fatalf("messages = %#v, want one hello message", got)
	}
}

func TestStorePersistsToolResultOutputToWorkspaceFileAndRestoresAfterRestart(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	content := strings.Repeat("工具输出内容\n", 32)
	sum := sha256.Sum256([]byte(content))
	output := storage.PersistedOutput{
		ID:             "po_abc123",
		ConversationID: conv.ID,
		UserID:         conv.UserID,
		MessageID:      "msg_tool",
		ToolCallID:     "call_1",
		Kind:           "tool_result",
		Strategy:       "tool_result_compression",
		OriginalBytes:  len([]byte(content)),
		ContentSHA256:  hex.EncodeToString(sum[:]),
		Content:        content,
		Preview:        "工具输出内容",
	}

	if err := store.CreatePersistedOutput(ctx, output); err != nil {
		t.Fatalf("CreatePersistedOutput returned error: %v", err)
	}
	contentPath := filepath.Join(home, ".cynosure", "task_outputs", "tool-results", conv.SessionID+"-"+output.ID+".txt")
	metaPath := filepath.Join(home, ".cynosure", "task_outputs", "tool-results", conv.SessionID+"-"+output.ID+".json")
	contentBytes, err := os.ReadFile(contentPath)
	if err != nil {
		t.Fatalf("expected persisted content file: %v", err)
	}
	if string(contentBytes) != content {
		t.Fatalf("persisted content = %q, want original content", string(contentBytes))
	}
	metaBytes, err := os.ReadFile(metaPath)
	if err != nil {
		t.Fatalf("expected persisted metadata file: %v", err)
	}
	meta := string(metaBytes)
	if !strings.Contains(meta, `"session_id": "`+conv.SessionID+`"`) || !strings.Contains(meta, `"content_file": "`+conv.SessionID+`-`+output.ID+`.txt"`) {
		t.Fatalf("metadata missing session/content file fields: %s", meta)
	}

	freshStore, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory fresh returned error: %v", err)
	}
	if err := freshStore.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation fresh returned error: %v", err)
	}
	restored, err := freshStore.GetPersistedOutputForConversation(ctx, output.ID, conv.UserID, conv.ID)
	if err != nil {
		t.Fatalf("GetPersistedOutputForConversation returned error: %v", err)
	}
	if restored.Content != content || restored.ContentSHA256 != output.ContentSHA256 || restored.ToolCallID != output.ToolCallID {
		t.Fatalf("restored output = %#v, want original content/hash/tool call", restored)
	}
}

func TestStoreRestoresPersistedOutputByMessageHashFromWorkspace(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	content := strings.Repeat("跨重启工具输出\n", 16)
	sum := sha256.Sum256([]byte(content))
	output := storage.PersistedOutput{
		ID:             "po_hash_restore",
		ConversationID: conv.ID,
		UserID:         conv.UserID,
		MessageID:      "msg_tool",
		ToolCallID:     "call_1",
		Kind:           "tool_result",
		Strategy:       "tool_result_compression",
		OriginalBytes:  len([]byte(content)),
		ContentSHA256:  hex.EncodeToString(sum[:]),
		Content:        content,
		Preview:        "跨重启工具输出",
	}
	if err := store.CreatePersistedOutput(ctx, output); err != nil {
		t.Fatalf("CreatePersistedOutput returned error: %v", err)
	}

	freshStore, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory fresh returned error: %v", err)
	}
	if err := freshStore.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation fresh returned error: %v", err)
	}
	restored, err := freshStore.GetPersistedOutputByMessageHash(ctx, conv.ID, conv.UserID, output.MessageID, output.ToolCallID, output.Strategy, output.ContentSHA256)
	if err != nil {
		t.Fatalf("GetPersistedOutputByMessageHash returned error: %v", err)
	}
	if restored.ID != output.ID || restored.Content != content || restored.ContentSHA256 != output.ContentSHA256 {
		t.Fatalf("restored output = %#v, want original file-backed output", restored)
	}
}

func TestStoreRejectsCorruptPersistedToolResultFile(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	content := "完整工具输出"
	sum := sha256.Sum256([]byte(content))
	output := storage.PersistedOutput{ID: "po_bad", ConversationID: conv.ID, UserID: conv.UserID, MessageID: "msg_tool", ToolCallID: "call_1", Kind: "tool_result", Strategy: "tool_result_compression", OriginalBytes: len([]byte(content)), ContentSHA256: hex.EncodeToString(sum[:]), Content: content, Preview: "完整"}
	if err := store.CreatePersistedOutput(ctx, output); err != nil {
		t.Fatalf("CreatePersistedOutput returned error: %v", err)
	}
	contentPath := filepath.Join(home, ".cynosure", "task_outputs", "tool-results", conv.SessionID+"-"+output.ID+".txt")
	if err := os.WriteFile(contentPath, []byte("被篡改"), 0o644); err != nil {
		t.Fatalf("tamper content file: %v", err)
	}
	freshStore, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory fresh returned error: %v", err)
	}
	if err := freshStore.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation fresh returned error: %v", err)
	}
	if _, err := freshStore.GetPersistedOutputForConversation(ctx, output.ID, conv.UserID, conv.ID); err == nil {
		t.Fatal("expected corrupt persisted output file to be rejected")
	}
}

func TestStoreRejectsPersistedToolResultFromDifferentUser(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	content := "完整工具输出"
	sum := sha256.Sum256([]byte(content))
	output := storage.PersistedOutput{ID: "po_scope", ConversationID: conv.ID, UserID: conv.UserID, MessageID: "msg_tool", ToolCallID: "call_1", Kind: "tool_result", Strategy: "tool_result_compression", OriginalBytes: len([]byte(content)), ContentSHA256: hex.EncodeToString(sum[:]), Content: content, Preview: "完整"}
	if err := store.CreatePersistedOutput(ctx, output); err != nil {
		t.Fatalf("CreatePersistedOutput returned error: %v", err)
	}
	freshStore, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory fresh returned error: %v", err)
	}
	if err := freshStore.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation fresh returned error: %v", err)
	}
	if _, err := freshStore.GetPersistedOutputForConversation(ctx, output.ID, "other-user", conv.ID); err == nil {
		t.Fatal("expected persisted output from a different user to be rejected")
	}
}

func TestStoreAppendsToolResultLogUnderSessionDirectory(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	first := storage.ToolResultLogEntry{ConversationID: conv.ID, UserID: conv.UserID, ToolCallID: "call_1", ToolName: "bash", RawArgs: `{"command":"pwd"}`, Status: "success", Result: "line1\n```\nline2", AuditSummary: `{"outcome_summary":"line1"}`}
	second := storage.ToolResultLogEntry{ConversationID: conv.ID, UserID: conv.UserID, ToolCallID: "call_2", ToolName: "read_file", RawArgs: `{"path":"README.md"}`, Status: "error", Result: "not found"}

	if err := store.AppendToolResultLog(ctx, first); err != nil {
		t.Fatalf("AppendToolResultLog first returned error: %v", err)
	}
	if err := store.AppendToolResultLog(ctx, second); err != nil {
		t.Fatalf("AppendToolResultLog second returned error: %v", err)
	}
	logPath := filepath.Join(home, ".cynosure", "task_outputs", conv.SessionID, "tools.md")
	bodyBytes, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read tools.md: %v", err)
	}
	body := string(bodyBytes)
	for _, want := range []string{"bash · success", "read_file · error", "tool_call_id: call_1", "tool_call_id: call_2", "line1", "not found"} {
		if !strings.Contains(body, want) {
			t.Fatalf("tools.md missing %q: %s", want, body)
		}
	}
	if !strings.Contains(body, "````text\nline1\n```\nline2\n````") {
		t.Fatalf("tools.md should use a longer fence when result contains backticks: %s", body)
	}
}

func TestStoreConversationCacheIsIndependentCopy(t *testing.T) {
	store := NewStore()
	ctx := context.Background()
	messages := []storage.Message{{ID: "msg_1", Role: "user", Content: "hello"}}
	if err := store.SetConversationCache(ctx, "conv_1", messages); err != nil {
		t.Fatalf("SetConversationCache returned error: %v", err)
	}
	messages[0].Content = "mutated"

	got, ok, err := store.GetConversationCache(ctx, "conv_1")
	if err != nil {
		t.Fatalf("GetConversationCache returned error: %v", err)
	}
	if !ok {
		t.Fatal("cache not found")
	}
	if got[0].Content != "hello" {
		t.Fatalf("cached content = %q, want hello", got[0].Content)
	}
}

func TestMarkdownMemoryStoreWritesProjectScopedMemoryIndex(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	mem := storage.Memory{UserID: LocalUserID, Type: runtime.MemoryTypeProject, Name: "构建命令", Description: "使用 go test", Body: "当前项目使用 go test ./... 验证。"}
	if err := store.InsertMemory(ctx, mem); err != nil {
		t.Fatalf("InsertMemory returned error: %v", err)
	}

	memoryRoot := filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(workspace))
	indexPath := filepath.Join(memoryRoot, "memory.md")
	indexBytes, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("read memory index: %v", err)
	}
	index := string(indexBytes)
	if !strings.Contains(index, "# Memory Index") || !strings.Contains(index, "构建命令") || !strings.Contains(index, "使用 go test") {
		t.Fatalf("unexpected index content: %q", index)
	}

	items, err := store.ListRelevantMemories(ctx, LocalUserID)
	if err != nil {
		t.Fatalf("ListRelevantMemories returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected one memory, got %#v", items)
	}
	if items[0].Type != runtime.MemoryTypeProject || !strings.Contains(items[0].Body, "当前项目使用 go test") {
		t.Fatalf("unexpected memory item: %#v", items[0])
	}
}

func TestMarkdownMemoryStoreScanUpdateDeleteAndIndexLimits(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	if err := store.InsertMemory(ctx, storage.Memory{UserID: LocalUserID, Type: runtime.MemoryTypePreference, Name: "偏好", Description: "简洁中文", Body: "正文"}); err != nil {
		t.Fatalf("InsertMemory returned error: %v", err)
	}

	scanned, err := store.ScanRecentMemories(ctx)
	if err != nil {
		t.Fatalf("ScanRecentMemories returned error: %v", err)
	}
	if len(scanned) != 1 || scanned[0].Name != "偏好" || scanned[0].Type != runtime.MemoryTypePreference {
		t.Fatalf("unexpected scanned memories: %#v", scanned)
	}
	if scanned[0].Path == "memory.md" {
		t.Fatalf("scan should exclude memory.md")
	}

	newBody := "更新后的正文"
	if _, err := store.UpdateMemoryFile(ctx, scanned[0].Path, storage.MemoryUpdate{Body: &newBody}); err != nil {
		t.Fatalf("UpdateMemoryFile returned error: %v", err)
	}
	mem, err := store.ReadMemoryFile(ctx, scanned[0].Path)
	if err != nil {
		t.Fatalf("ReadMemoryFile returned error: %v", err)
	}
	if mem.Body != newBody {
		t.Fatalf("expected updated body %q, got %q", newBody, mem.Body)
	}

	if err := store.DeleteMemoryFile(ctx, scanned[0].Path); err != nil {
		t.Fatalf("DeleteMemoryFile returned error: %v", err)
	}
	memoryRoot := filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(workspace))
	if _, err := os.Stat(filepath.Join(memoryRoot, scanned[0].Path)); !os.IsNotExist(err) {
		t.Fatalf("expected memory file to be deleted, stat err = %v", err)
	}
	indexBytes, err := os.ReadFile(filepath.Join(memoryRoot, "memory.md"))
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	if strings.Contains(string(indexBytes), "偏好") {
		t.Fatalf("expected index entry removed after delete, got %q", string(indexBytes))
	}

	var b strings.Builder
	b.WriteString(memoryIndexHeader)
	for i := 0; i < memoryIndexMaxLines+50; i++ {
		b.WriteString("- entry line\n")
	}
	if err := atomicWriteFile(filepath.Join(memoryRoot, "memory.md"), []byte(b.String()), 0o644); err != nil {
		t.Fatalf("write big index: %v", err)
	}
	text, truncated, total := store.LoadMemoryIndexForPrompt(ctx)
	if !truncated {
		t.Fatalf("expected truncation, got truncated=false total=%d", total)
	}
	if total < memoryIndexMaxLines+50 {
		t.Fatalf("expected total to reflect real line count, got %d", total)
	}
	if lines := strings.Count(text, "\n") + 1; lines > memoryIndexMaxLines {
		t.Fatalf("expected at most %d lines, got %d", memoryIndexMaxLines, lines)
	}
}

func TestMarkdownMemoryStoreUpdateRenamesFileWhenNameChanges(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	if err := store.InsertMemory(ctx, storage.Memory{UserID: LocalUserID, Type: runtime.MemoryTypePreference, Name: "用户喜欢 Go 语言", Description: "d", Body: "正文"}); err != nil {
		t.Fatalf("InsertMemory returned error: %v", err)
	}
	scanned, err := store.ScanRecentMemories(ctx)
	if err != nil || len(scanned) != 1 {
		t.Fatalf("ScanRecentMemories returned err=%v scanned=%#v", err, scanned)
	}
	oldPath := scanned[0].Path

	newName := "用户喜欢 Python 语言"
	gotPath, err := store.UpdateMemoryFile(ctx, oldPath, storage.MemoryUpdate{Name: &newName})
	if err != nil {
		t.Fatalf("UpdateMemoryFile returned error: %v", err)
	}
	if gotPath == oldPath {
		t.Fatalf("expected renamed path, got same path %q", gotPath)
	}

	memoryRoot := filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(workspace))
	if _, err := os.Stat(filepath.Join(memoryRoot, oldPath)); !os.IsNotExist(err) {
		t.Fatalf("expected old file removed, stat err = %v", err)
	}
	mem, err := store.ReadMemoryFile(ctx, gotPath)
	if err != nil {
		t.Fatalf("ReadMemoryFile returned error: %v", err)
	}
	if mem.Name != newName {
		t.Fatalf("expected name %q, got %q", newName, mem.Name)
	}
	indexBytes, err := os.ReadFile(filepath.Join(memoryRoot, "memory.md"))
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	if strings.Contains(string(indexBytes), oldPath) {
		t.Fatalf("expected old path removed from index, got %q", string(indexBytes))
	}
	if !strings.Contains(string(indexBytes), gotPath) {
		t.Fatalf("expected new path in index, got %q", string(indexBytes))
	}
}

func TestMarkdownMemoryStoreConsolidationState(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	state, err := store.LoadConsolidationState(ctx)
	if err != nil {
		t.Fatalf("LoadConsolidationState returned error: %v", err)
	}
	if state.SessionCount != 0 || !state.LastRunAt.IsZero() {
		t.Fatalf("expected zero state, got %#v", state)
	}
	state.SessionCount = 3
	state.LastRunAt = time.Now()
	if err := store.SaveConsolidationState(ctx, state); err != nil {
		t.Fatalf("SaveConsolidationState returned error: %v", err)
	}
	reloaded, err := store.LoadConsolidationState(ctx)
	if err != nil {
		t.Fatalf("reload consolidation state: %v", err)
	}
	if reloaded.SessionCount != 3 || reloaded.LastRunAt.IsZero() {
		t.Fatalf("expected persisted state, got %#v", reloaded)
	}
}

func TestMarkdownMemoryStoreDoesNotListOtherWorkspaceMemoriesFromHomeLink(t *testing.T) {
	home := t.TempDir()
	workspace := filepath.Join(t.TempDir(), "project-a")
	otherWorkspace := filepath.Join(t.TempDir(), "project-b")
	t.Setenv("HOME", home)
	ctx := context.Background()

	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	if err := store.InsertMemory(ctx, storage.Memory{UserID: LocalUserID, Type: runtime.MemoryTypeProject, Name: "当前项目", Description: "current", Body: "只属于当前项目"}); err != nil {
		t.Fatalf("InsertMemory current returned error: %v", err)
	}
	otherStore, err := NewStoreWithMemory(otherWorkspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory other returned error: %v", err)
	}
	if err := otherStore.InsertMemory(ctx, storage.Memory{UserID: LocalUserID, Type: runtime.MemoryTypeProject, Name: "其他项目", Description: "other", Body: "不应被当前项目检索到"}); err != nil {
		t.Fatalf("InsertMemory other returned error: %v", err)
	}

	freshStore, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory fresh returned error: %v", err)
	}
	items, err := freshStore.ListRelevantMemories(ctx, LocalUserID)
	if err != nil {
		t.Fatalf("ListRelevantMemories returned error: %v", err)
	}
	if len(items) != 1 || !strings.Contains(items[0].Body, "只属于当前项目") {
		t.Fatalf("items = %#v, want only current workspace memory", items)
	}
	if _, err := os.Stat(filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(workspace), "memory.md")); err != nil {
		t.Fatalf("current workspace memory index not under home .cynosure: %v", err)
	}
	if _, err := os.Stat(filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(otherWorkspace), "memory.md")); err != nil {
		t.Fatalf("other workspace memory index not under home .cynosure: %v", err)
	}
}

func TestMarkdownConversationMemoryUsesSessionIDFile(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "041581e7-c3e7-46c8-afe7-7cdcc671e80e", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	first := []storage.ConversationMemory{{Name: "目标", Description: "第一轮", Body: "实现记忆"}}
	if err := store.ReplaceConversationMemories(ctx, conv.ID, LocalUserID, first); err != nil {
		t.Fatalf("ReplaceConversationMemories first returned error: %v", err)
	}
	second := []storage.ConversationMemory{{Name: "目标", Description: "第二轮", Body: "更新同一个文件"}}
	if err := store.ReplaceConversationMemories(ctx, conv.ID, LocalUserID, second); err != nil {
		t.Fatalf("ReplaceConversationMemories second returned error: %v", err)
	}

	memoryRoot := filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(workspace))
	sessionPath := filepath.Join(memoryRoot, "sessions", conv.SessionID+".md")
	bodyBytes, err := os.ReadFile(sessionPath)
	if err != nil {
		t.Fatalf("read session memory: %v", err)
	}
	body := string(bodyBytes)
	if strings.Contains(body, "第一轮") || !strings.Contains(body, "第二轮") || !strings.Contains(body, conv.SessionID) {
		t.Fatalf("session memory should be overwritten in one file, got %q", body)
	}
	entries, err := filepath.Glob(filepath.Join(memoryRoot, "sessions", "*.md"))
	if err != nil {
		t.Fatalf("glob session memories: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected one session memory file, got %v", entries)
	}
}

func TestMarkdownConversationMemoryBreakpointPersistsAcrossWrites(t *testing.T) {
	home := t.TempDir()
	workspace := t.TempDir()
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_bp", SessionID: "bbbbbbbb-c3e7-46c8-afe7-7cdcc671e80e", UserID: LocalUserID}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	// Save memories, then a breakpoint; both must coexist.
	if err := store.ReplaceConversationMemories(ctx, conv.ID, LocalUserID, []storage.ConversationMemory{{Name: "目标", Body: "实现"}}); err != nil {
		t.Fatalf("ReplaceConversationMemories: %v", err)
	}
	if err := store.SaveConversationMemoryBreakpoint(ctx, conv.ID, "msg_bp_1"); err != nil {
		t.Fatalf("SaveConversationMemoryBreakpoint: %v", err)
	}
	got, err := store.LoadConversationMemoryBreakpoint(ctx, conv.ID)
	if err != nil || got != "msg_bp_1" {
		t.Fatalf("expected breakpoint msg_bp_1, got %q err=%v", got, err)
	}
	// Rewriting memories must preserve the breakpoint.
	if err := store.ReplaceConversationMemories(ctx, conv.ID, LocalUserID, []storage.ConversationMemory{{Name: "目标", Body: "更新"}}); err != nil {
		t.Fatalf("ReplaceConversationMemories second: %v", err)
	}
	got, err = store.LoadConversationMemoryBreakpoint(ctx, conv.ID)
	if err != nil || got != "msg_bp_1" {
		t.Fatalf("expected breakpoint preserved after memory rewrite, got %q err=%v", got, err)
	}
	// Updating the breakpoint must preserve memories.
	if err := store.SaveConversationMemoryBreakpoint(ctx, conv.ID, "msg_bp_2"); err != nil {
		t.Fatalf("SaveConversationMemoryBreakpoint second: %v", err)
	}
	items, err := store.ListConversationMemories(ctx, conv.ID)
	if err != nil || len(items) != 1 || items[0].Name != "目标" {
		t.Fatalf("expected memories preserved after breakpoint update, got %#v err=%v", items, err)
	}
	got, _ = store.LoadConversationMemoryBreakpoint(ctx, conv.ID)
	if got != "msg_bp_2" {
		t.Fatalf("expected breakpoint msg_bp_2, got %q", got)
	}
}

func TestStorePersistsConversationAndModelHistoryUnderSessionDirectory(t *testing.T) {
	home := t.TempDir()
	workspace := filepath.Join(t.TempDir(), "project")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	ctx := context.Background()
	conv := storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: LocalUserID, Title: "TUI 会话"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	displayHistory := []storage.Message{{ID: "msg_1", ConversationID: conv.ID, UserID: conv.UserID, Role: "user", Content: "hello"}}
	modelHistory := []storage.Message{{ID: "msg_2", ConversationID: conv.ID, UserID: conv.UserID, Role: "assistant", Content: "hi"}}

	if err := store.SetConversationHistory(ctx, conv.ID, displayHistory); err != nil {
		t.Fatalf("SetConversationHistory returned error: %v", err)
	}
	if err := store.UpsertConversationModelHistory(ctx, conv.ID, conv.UserID, modelHistory); err != nil {
		t.Fatalf("UpsertConversationModelHistory returned error: %v", err)
	}

	historyPath := filepath.Join(home, ".cynosure", "session", conv.SessionID, "history")
	modelHistoryPath := filepath.Join(home, ".cynosure", "session", conv.SessionID, "model_history")
	for _, path := range []string{historyPath, modelHistoryPath} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected %s to exist: %v", path, err)
		}
	}
	loadedHistory, err := store.ListMessagesByConversation(ctx, conv.ID, 100)
	if err != nil {
		t.Fatalf("ListMessagesByConversation returned error: %v", err)
	}
	if len(loadedHistory) != 1 || loadedHistory[0].Content != "hello" {
		t.Fatalf("display history = %#v, want hello", loadedHistory)
	}
	loadedModelHistory, ok, err := store.GetConversationModelHistory(ctx, conv.ID)
	if err != nil {
		t.Fatalf("GetConversationModelHistory returned error: %v", err)
	}
	if !ok || len(loadedModelHistory) != 1 || loadedModelHistory[0].Content != "hi" {
		t.Fatalf("model history = %#v, ok=%v; want hi", loadedModelHistory, ok)
	}
}

func TestStoreListsAndResumesOnlyCurrentWorkspaceSessions(t *testing.T) {
	home := t.TempDir()
	workspace := filepath.Join(t.TempDir(), "project")
	otherWorkspace := filepath.Join(t.TempDir(), "other")
	if err := os.MkdirAll(workspace, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(otherWorkspace, 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HOME", home)
	ctx := context.Background()

	store, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory returned error: %v", err)
	}
	conv := storage.Conversation{ID: "conv_current", SessionID: "current-session", UserID: LocalUserID, Title: "当前项目"}
	if err := store.CreateConversation(ctx, conv); err != nil {
		t.Fatalf("CreateConversation returned error: %v", err)
	}
	if err := store.SetConversationHistory(ctx, conv.ID, []storage.Message{{Role: "user", Content: "current"}}); err != nil {
		t.Fatalf("SetConversationHistory current returned error: %v", err)
	}

	otherStore, err := NewStoreWithMemory(otherWorkspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory other returned error: %v", err)
	}
	otherConv := storage.Conversation{ID: "conv_other", SessionID: "other-session", UserID: LocalUserID, Title: "其他项目"}
	if err := otherStore.CreateConversation(ctx, otherConv); err != nil {
		t.Fatalf("CreateConversation other returned error: %v", err)
	}
	if err := otherStore.SetConversationHistory(ctx, otherConv.ID, []storage.Message{{Role: "user", Content: "other"}}); err != nil {
		t.Fatalf("SetConversationHistory other returned error: %v", err)
	}

	freshStore, err := NewStoreWithMemory(workspace)
	if err != nil {
		t.Fatalf("NewStoreWithMemory fresh returned error: %v", err)
	}
	sessions, err := freshStore.ListResumableSessions(ctx, workspace)
	if err != nil {
		t.Fatalf("ListResumableSessions returned error: %v", err)
	}
	if len(sessions) != 1 || sessions[0].SessionID != conv.SessionID || sessions[0].MessageCount != 1 {
		t.Fatalf("sessions = %#v, want only current-session", sessions)
	}

	resumed, history, err := freshStore.ResumeSession(ctx, conv.SessionID, workspace, storage.User{ID: LocalUserID})
	if err != nil {
		t.Fatalf("ResumeSession returned error: %v", err)
	}
	if resumed.SessionID != conv.SessionID || resumed.ID != conv.ID {
		t.Fatalf("resumed conversation = %#v, want original conversation", resumed)
	}
	if len(history) != 1 || history[0].Content != "current" {
		t.Fatalf("history = %#v, want current message", history)
	}
}

func TestBootstrapLoadsCynosureMarkdownIntoRuntime(t *testing.T) {
	tmp := t.TempDir()
	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldWD) })
	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "config.json"), []byte(`{
		"app_home":".",
		"system_prompt_path":"system_prompt.md",
		"builtin_skills_dir":"skills",
		"command_bin_dir":"bin",
		"command_script_dir":"cmd"
	}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "system_prompt.md"), []byte("Base prompt."), 0o644); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(tmp, "home")
	workspace := filepath.Join(tmp, "project")
	t.Setenv("HOME", home)
	writeFile(t, filepath.Join(home, ".cynosure", "settings.json"), `{"env":{"open_auth_token":"key","open_model":"model","open_base_url":"https://example.com"}}`)
	writeFile(t, filepath.Join(home, ".cynosure", "CYNOSURE.MD"), "# User Rule\n全局说明")
	writeFile(t, filepath.Join(workspace, ".cynosure", "CYNOSURE.MD"), "# Project Rule\n项目说明")

	bundle, err := Bootstrap(context.Background(), workspace)
	if err != nil {
		t.Fatalf("Bootstrap returned error: %v", err)
	}
	t.Cleanup(bundle.Close)
	if bundle.Runtime.CynosureMarkdown.UserContent != "# User Rule\n全局说明" {
		t.Fatalf("UserContent = %q", bundle.Runtime.CynosureMarkdown.UserContent)
	}
	if bundle.Runtime.CynosureMarkdown.WorkspaceContent != "# Project Rule\n项目说明" {
		t.Fatalf("WorkspaceContent = %q", bundle.Runtime.CynosureMarkdown.WorkspaceContent)
	}
	if bundle.Runtime.CynosureMarkdown.UserPath != filepath.Join(home, ".cynosure", "CYNOSURE.MD") {
		t.Fatalf("UserPath = %q", bundle.Runtime.CynosureMarkdown.UserPath)
	}
	if bundle.Runtime.CynosureMarkdown.WorkspacePath != filepath.Join(workspace, ".cynosure", "CYNOSURE.MD") {
		t.Fatalf("WorkspacePath = %q", bundle.Runtime.CynosureMarkdown.WorkspacePath)
	}
	if !bundle.Runtime.EnableMemory || !bundle.User.MemoryEnabled {
		t.Fatalf("memory should be enabled in TUI bootstrap")
	}
	if bundle.Conversation.SessionID == "" || bundle.Conversation.SessionID == bundle.Conversation.ID {
		t.Fatalf("SessionID = %q, ConversationID = %q", bundle.Conversation.SessionID, bundle.Conversation.ID)
	}
	if _, err := os.Stat(filepath.Join(home, ".cynosure", "memory", workspaceMemoryDirName(workspace), "memory.md")); err != nil {
		t.Fatalf("memory index not created: %v", err)
	}
}

func writeFile(t *testing.T, path string, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

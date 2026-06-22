package local

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"

	"cynosure/internal/agent/storage"
	"cynosure/internal/config"
	"cynosure/internal/idgen"
)

const memoryIndexHeader = "# Memory Index\n"

const (
	// FrontmatterMaxLines 限制扫描候选记忆时读取的 frontmatter 行数。
	FrontmatterMaxLines = 30
	// MaxMemoryFiles 限制确定性扫描保留的候选记忆文件数量。
	MaxMemoryFiles = 200
	// memoryIndexMaxLines 限制 memory.md 注入系统提示词时保留的行数。
	memoryIndexMaxLines = 200
	// memoryIndexMaxEntryBytes 限制 memory.md 单行（单条索引）注入的字节数。
	memoryIndexMaxEntryBytes = 25 * 1024
)

const consolidationStateFile = ".consolidation_state.json"
const injectedMemoriesStateFile = "injected_memories.json"

type MarkdownMemoryStore struct {
	mu          sync.Mutex
	rootDir     string
	indexPath   string
	sessionsDir string
	projectName string
}

type injectedMemoriesState map[string]time.Time

func NewMarkdownMemoryStore(workspaceRoot string) (*MarkdownMemoryStore, error) {
	projectName := sanitizeName(filepath.Base(filepath.Clean(workspaceRoot)))
	if projectName == "" {
		projectName = "project"
	}
	baseDir, err := config.CynosureMemoryDir()
	if err != nil {
		return nil, err
	}
	workspaceName, err := config.WorkspaceName(workspaceRoot)
	if err != nil {
		return nil, err
	}
	root := filepath.Join(baseDir, workspaceName)
	store := &MarkdownMemoryStore{
		rootDir:     root,
		indexPath:   filepath.Join(root, "memory.md"),
		sessionsDir: filepath.Join(root, "sessions"),
		projectName: projectName,
	}
	if err := store.EnsureLayout(); err != nil {
		return nil, err
	}
	return store, nil
}

func (m *MarkdownMemoryStore) EnsureLayout() error {
	if err := os.MkdirAll(m.sessionsDir, 0o755); err != nil {
		return err
	}
	if _, err := os.Stat(m.indexPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	return atomicWriteFile(m.indexPath, []byte(memoryIndexHeader+"\n"), 0o644)
}

func (m *MarkdownMemoryStore) ListRelevantMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.listRelevantMemoriesLocked()
}

func (m *MarkdownMemoryStore) InsertMemory(ctx context.Context, mem storage.Memory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if mem.ID == "" {
		mem.ID = idgen.New("mem")
	}
	if mem.Type == "" {
		mem.Type = "preference"
	}
	if mem.CreatedAt.IsZero() {
		mem.CreatedAt = time.Now()
	}
	mem.UpdatedAt = time.Now()
	filename := m.uniqueMemoryFilenameLocked(mem.Name)
	if err := m.writeMemoryFileLocked(filename, mem); err != nil {
		return err
	}
	return m.rewriteIndexLocked()
}

func (m *MarkdownMemoryStore) ListMemoriesByUserAndType(ctx context.Context, userID, memType string) ([]storage.Memory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	items, err := m.listRelevantMemoriesLocked()
	if err != nil {
		return nil, err
	}
	return filterMemories(items, memType), nil
}

func (m *MarkdownMemoryStore) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	items, err := m.ListMemoriesByUserAndType(ctx, userID, memType)
	return len(items), err
}

func (m *MarkdownMemoryStore) DeleteOldestMemories(ctx context.Context, userID, memType string, n int) error {
	if n <= 0 {
		return nil
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return err
	}
	deleted := 0
	kept := entries[:0]
	for _, entry := range entries {
		mem, err := m.readMemoryFileLocked(entry.Path)
		if err == nil && mem.Type == memType && deleted < n {
			_ = os.Remove(filepath.Join(m.rootDir, entry.Path))
			deleted++
			continue
		}
		kept = append(kept, entry)
	}
	return m.writeIndexEntriesLocked(kept)
}

func (m *MarkdownMemoryStore) ReplaceMemoriesByUserAndType(ctx context.Context, userID, memType string, items []storage.Memory) error {
	return m.replaceMemoriesByType(memType, items)
}

// LoadMemoryIndexForPrompt 读取 memory.md 用于注入系统提示词，受行数与单行字节
// 上限约束。仅当存在索引条目时返回文本；空索引不注入系统提示词。返回截断后
// 的文本、是否发生截断、以及参与注入的真实索引行数。
func (m *MarkdownMemoryStore) LoadMemoryIndexForPrompt() (string, bool, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		return "", false, 0
	}
	rawLines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	entryLines := make([]string, 0, len(rawLines))
	for _, line := range rawLines {
		if strings.HasPrefix(strings.TrimSpace(line), "- [") {
			entryLines = append(entryLines, line)
		}
	}
	if len(entryLines) == 0 {
		return "", false, 0
	}
	totalLines := len(entryLines)
	truncated := false
	entryLimit := memoryIndexMaxLines - 1
	if entryLimit < 1 {
		entryLimit = 1
	}
	if totalLines > entryLimit {
		entryLines = entryLines[:entryLimit]
		truncated = true
	}
	kept := make([]string, 0, len(entryLines)+1)
	kept = append(kept, strings.TrimRight(memoryIndexHeader, "\n"))
	for _, line := range entryLines {
		if len(line) > memoryIndexMaxEntryBytes {
			line = truncateBytesAtRune(line, memoryIndexMaxEntryBytes)
			truncated = true
		}
		kept = append(kept, line)
	}
	return strings.Join(kept, "\n"), truncated, totalLines
}

// ScanRecentMemories 扫描记忆目录下所有 .md（排除 memory.md），每个文件只读前
// FrontmatterMaxLines 行 frontmatter，按 mtime 降序保留最新 MaxMemoryFiles 个。
func (m *MarkdownMemoryStore) ScanRecentMemories() ([]storage.ScannedMemory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	files, err := filepath.Glob(filepath.Join(m.rootDir, "*.md"))
	if err != nil {
		return nil, err
	}
	result := make([]storage.ScannedMemory, 0, len(files))
	for _, file := range files {
		if filepath.Clean(file) == filepath.Clean(m.indexPath) {
			continue
		}
		rel, err := filepath.Rel(m.rootDir, file)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)
		if !safeRelativeMemoryPath(rel) {
			continue
		}
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		meta := readFrontmatterHead(file, FrontmatterMaxLines)
		result = append(result, storage.ScannedMemory{
			Path:        rel,
			Name:        meta["name"],
			Description: meta["description"],
			Type:        meta["metadata.type"],
			ModTime:     info.ModTime(),
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ModTime.After(result[j].ModTime)
	})
	if len(result) > MaxMemoryFiles {
		result = result[:MaxMemoryFiles]
	}
	return result, nil
}

// ReadMemoryFile 读取单条记忆文件的完整内容（name/description/type/body）。
func (m *MarkdownMemoryStore) ReadMemoryFile(path string) (storage.Memory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.readMemoryFileLocked(path)
}

// ShouldInjectMemory 实现文件级会话内去重：同一会话、同一记忆文件、同一 mtime
// 已注入过则跳过；mtime 变化则更新持久状态并允许重新注入。
func (m *MarkdownMemoryStore) ShouldInjectMemory(conversationID, path string, modTime time.Time) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if strings.TrimSpace(conversationID) == "" {
		return true, nil
	}
	if !safeRelativeMemoryPath(path) {
		return false, fmt.Errorf("unsafe memory path: %s", path)
	}
	state, err := m.readInjectedMemoriesStateLocked(conversationID)
	if err != nil {
		return false, err
	}
	prev, ok := state[path]
	if ok && prev.Equal(modTime) {
		return false, nil
	}
	state[path] = modTime
	if err := m.writeInjectedMemoriesStateLocked(conversationID, state); err != nil {
		return false, err
	}
	return true, nil
}

// ForgetInjectedMemory 清除所有会话中指定记忆文件的注入记录。记忆内容被更新、删除
// 或重命名后调用，确保后续轮次不被旧状态误去重。
func (m *MarkdownMemoryStore) ForgetInjectedMemory(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !safeRelativeMemoryPath(path) {
		return fmt.Errorf("unsafe memory path: %s", path)
	}
	entries, err := os.ReadDir(m.sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		conversationID := entry.Name()
		state, err := m.readInjectedMemoriesStateLocked(conversationID)
		if err != nil {
			return err
		}
		if _, ok := state[path]; !ok {
			continue
		}
		delete(state, path)
		if err := m.writeInjectedMemoriesStateLocked(conversationID, state); err != nil {
			return err
		}
	}
	return nil
}

// UpdateMemoryFile 部分更新指定记忆文件，并同步刷新 memory.md。当记忆标题（name）
// 发生变化时，文件名（即记忆的标题标识）会一并重命名，避免文件名与内容标题不一致
// 带来歧义。返回更新后真实的相对路径（重命名时与传入 path 不同）。
func (m *MarkdownMemoryStore) UpdateMemoryFile(path string, update storage.MemoryUpdate) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !safeRelativeMemoryPath(path) {
		return "", fmt.Errorf("unsafe memory path: %s", path)
	}
	mem, err := m.readMemoryFileLocked(path)
	if err != nil {
		return "", err
	}
	nameChanged := false
	if update.Name != nil {
		newName := strings.TrimSpace(*update.Name)
		nameChanged = newName != mem.Name
		mem.Name = newName
	}
	if update.Description != nil {
		mem.Description = strings.TrimSpace(*update.Description)
	}
	if update.Body != nil {
		mem.Body = strings.TrimSpace(*update.Body)
	}
	newPath := path
	if nameChanged {
		// 仅当标题对应的目标文件名与当前文件名不同才重命名，避免无意义改名。
		if desired := sanitizedMemoryFilename(mem.Name); desired != path {
			newPath = m.uniqueMemoryFilenameLocked(mem.Name)
			// 文件名是记忆的标题标识，重命名后令 ID 与新文件名保持一致。
			mem.ID = strings.TrimSuffix(newPath, filepath.Ext(newPath))
		}
	}
	if err := m.writeMemoryFileLocked(newPath, mem); err != nil {
		return "", err
	}
	if newPath != path {
		if err := os.Remove(filepath.Join(m.rootDir, path)); err != nil && !os.IsNotExist(err) {
			return "", err
		}
	}
	if err := m.rewriteIndexLocked(); err != nil {
		return "", err
	}
	return newPath, nil
}

// DeleteMemoryFile 删除指定记忆文件，并从 memory.md 移除其条目。
func (m *MarkdownMemoryStore) DeleteMemoryFile(path string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if !safeRelativeMemoryPath(path) {
		return fmt.Errorf("unsafe memory path: %s", path)
	}
	full := filepath.Join(m.rootDir, path)
	if err := os.Remove(full); err != nil && !os.IsNotExist(err) {
		return err
	}
	return m.rewriteIndexLocked()
}

// LoadConsolidationState 读取定时去重的状态元数据；文件缺失时返回零值。
func (m *MarkdownMemoryStore) LoadConsolidationState() (storage.ConsolidationState, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := os.ReadFile(filepath.Join(m.rootDir, consolidationStateFile))
	if err != nil {
		if os.IsNotExist(err) {
			return storage.ConsolidationState{}, nil
		}
		return storage.ConsolidationState{}, err
	}
	var state storage.ConsolidationState
	if err := json.Unmarshal(data, &state); err != nil {
		return storage.ConsolidationState{}, nil
	}
	return state, nil
}

// SaveConsolidationState 持久化定时去重的状态元数据。
func (m *MarkdownMemoryStore) SaveConsolidationState(state storage.ConsolidationState) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return atomicWriteFile(filepath.Join(m.rootDir, consolidationStateFile), data, 0o644)
}

func (m *MarkdownMemoryStore) ListConversationMemories(ctx context.Context, sessionID string) ([]storage.ConversationMemory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	items, _ := m.readSessionMemoryLocked(sessionID)
	return items, nil
}

func (m *MarkdownMemoryStore) ReplaceConversationMemories(ctx context.Context, sessionID, userID string, items []storage.ConversationMemory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// 保留已有断点，避免重写会话记忆时丢失。
	_, breakpoint := m.readSessionMemoryLocked(sessionID)
	return m.writeSessionMemoryLocked(sessionID, items, breakpoint)
}

// LoadConversationMemoryBreakpoint 读取会话记忆文件 frontmatter 中持久化的断点 ID。
func (m *MarkdownMemoryStore) LoadConversationMemoryBreakpoint(ctx context.Context, sessionID string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, breakpoint := m.readSessionMemoryLocked(sessionID)
	return breakpoint, nil
}

// SaveConversationMemoryBreakpoint 持久化断点 ID 到会话记忆文件 frontmatter，保留已有
// 会话记忆条目。
func (m *MarkdownMemoryStore) SaveConversationMemoryBreakpoint(ctx context.Context, sessionID, breakpointID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	items, _ := m.readSessionMemoryLocked(sessionID)
	return m.writeSessionMemoryLocked(sessionID, items, breakpointID)
}

func (m *MarkdownMemoryStore) ProjectLockKey(sessionID string) string {
	return m.projectName + ":" + sanitizeName(sessionID)
}

func (m *MarkdownMemoryStore) replaceMemoriesByType(memType string, items []storage.Memory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return err
	}
	kept := entries[:0]
	for _, entry := range entries {
		mem, err := m.readMemoryFileLocked(entry.Path)
		if err == nil && mem.Type == memType {
			_ = os.Remove(filepath.Join(m.rootDir, entry.Path))
			continue
		}
		kept = append(kept, entry)
	}
	if err := m.writeIndexEntriesLocked(kept); err != nil {
		return err
	}
	for _, item := range items {
		if item.ID == "" {
			item.ID = idgen.New("mem")
		}
		if item.Type == "" {
			item.Type = memType
		}
		filename := m.uniqueMemoryFilenameLocked(item.Name)
		if err := m.writeMemoryFileLocked(filename, item); err != nil {
			return err
		}
	}
	return m.rewriteIndexLocked()
}

type memoryIndexEntry struct {
	Name        string
	Path        string
	Description string
}

func (m *MarkdownMemoryStore) listRelevantMemoriesLocked() ([]storage.Memory, error) {
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return nil, err
	}
	result := make([]storage.Memory, 0, len(entries))
	for _, entry := range entries {
		mem, err := m.readMemoryFileLocked(entry.Path)
		if err != nil {
			continue
		}
		result = append(result, mem)
	}
	return result, nil
}

func (m *MarkdownMemoryStore) indexEntriesLocked() ([]memoryIndexEntry, error) {
	data, err := os.ReadFile(m.indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var entries []memoryIndexEntry
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "- [") {
			continue
		}
		closeName := strings.Index(line, "](")
		closePath := strings.Index(line, ")")
		if closeName < 3 || closePath <= closeName+2 {
			continue
		}
		name := line[3:closeName]
		path := line[closeName+2 : closePath]
		desc := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line[closePath+1:]), "—"))
		if !safeRelativeMemoryPath(path) {
			continue
		}
		entries = append(entries, memoryIndexEntry{Name: name, Path: filepath.ToSlash(filepath.Clean(path)), Description: desc})
	}
	return entries, scanner.Err()
}

func (m *MarkdownMemoryStore) rewriteIndexLocked() error {
	entries, err := m.indexEntriesLocked()
	if err != nil {
		return err
	}
	seen := make(map[string]memoryIndexEntry, len(entries))
	for _, entry := range entries {
		if mem, err := m.readMemoryFileLocked(entry.Path); err == nil {
			entry.Name = mem.Name
			entry.Description = mem.Description
			seen[entry.Path] = entry
		}
	}
	files, err := filepath.Glob(filepath.Join(m.rootDir, "*.md"))
	if err != nil {
		return err
	}
	for _, file := range files {
		if filepath.Clean(file) == filepath.Clean(m.indexPath) {
			continue
		}
		rel, err := filepath.Rel(m.rootDir, file)
		if err != nil {
			continue
		}
		rel = filepath.ToSlash(rel)
		mem, err := m.readMemoryFileLocked(rel)
		if err != nil {
			continue
		}
		seen[rel] = memoryIndexEntry{Name: mem.Name, Path: rel, Description: mem.Description}
	}
	entries = entries[:0]
	for _, entry := range seen {
		entries = append(entries, entry)
	}
	return m.writeIndexEntriesLocked(entries)
}

func (m *MarkdownMemoryStore) writeIndexEntriesLocked(entries []memoryIndexEntry) error {
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Name == entries[j].Name {
			return entries[i].Path < entries[j].Path
		}
		return entries[i].Name < entries[j].Name
	})
	var b strings.Builder
	b.WriteString(memoryIndexHeader)
	b.WriteString("\n")
	for _, entry := range entries {
		b.WriteString(fmt.Sprintf("- [%s](%s) — %s\n", entry.Name, entry.Path, entry.Description))
	}
	return atomicWriteFile(m.indexPath, []byte(b.String()), 0o644)
}

// sanitizedMemoryFilename 把记忆标题转换为规范的相对文件名（不保证唯一）。
func sanitizedMemoryFilename(name string) string {
	base := sanitizeName(name)
	if base == "" {
		base = "memory"
	}
	return base + ".md"
}

func (m *MarkdownMemoryStore) uniqueMemoryFilenameLocked(name string) string {
	path := sanitizedMemoryFilename(name)
	if _, err := os.Stat(filepath.Join(m.rootDir, path)); os.IsNotExist(err) {
		return path
	}
	base := strings.TrimSuffix(path, filepath.Ext(path))
	return base + "-" + idgen.Hex()[:8] + ".md"
}

func (m *MarkdownMemoryStore) writeMemoryFileLocked(path string, mem storage.Memory) error {
	if !safeRelativeMemoryPath(path) {
		return fmt.Errorf("unsafe memory path: %s", path)
	}
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("name: " + yamlScalar(mem.Name) + "\n")
	b.WriteString("description: " + yamlScalar(mem.Description) + "\n")
	b.WriteString("metadata:\n")
	b.WriteString("  node_type: memory\n")
	b.WriteString("  type: " + yamlScalar(mem.Type) + "\n")
	b.WriteString("  project: " + yamlScalar(m.projectName) + "\n")
	b.WriteString("  originSessionId: " + yamlScalar(mem.ID) + "\n")
	b.WriteString("---\n\n")
	b.WriteString(strings.TrimSpace(mem.Body))
	b.WriteString("\n")
	return atomicWriteFile(filepath.Join(m.rootDir, path), []byte(b.String()), 0o644)
}

func (m *MarkdownMemoryStore) readMemoryFileLocked(path string) (storage.Memory, error) {
	if !safeRelativeMemoryPath(path) {
		return storage.Memory{}, fmt.Errorf("unsafe memory path: %s", path)
	}
	data, err := os.ReadFile(filepath.Join(m.rootDir, path))
	if err != nil {
		return storage.Memory{}, err
	}
	meta, body := splitFrontMatter(string(data))
	return storage.Memory{ID: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), Type: meta["metadata.type"], Name: meta["name"], Description: meta["description"], Body: strings.TrimSpace(body)}, nil
}

func (m *MarkdownMemoryStore) sessionPath(sessionID string) string {
	return filepath.Join(m.sessionsDir, sanitizeName(sessionID)+".md")
}

func (m *MarkdownMemoryStore) injectedMemoriesStatePath(conversationID string) string {
	return filepath.Join(m.sessionsDir, sanitizeName(conversationID), injectedMemoriesStateFile)
}

func (m *MarkdownMemoryStore) readInjectedMemoriesStateLocked(conversationID string) (injectedMemoriesState, error) {
	data, err := os.ReadFile(m.injectedMemoriesStatePath(conversationID))
	if err != nil {
		if os.IsNotExist(err) {
			return injectedMemoriesState{}, nil
		}
		return nil, err
	}
	var state injectedMemoriesState
	if err := json.Unmarshal(data, &state); err != nil {
		return injectedMemoriesState{}, nil
	}
	if state == nil {
		state = injectedMemoriesState{}
	}
	return state, nil
}

func (m *MarkdownMemoryStore) writeInjectedMemoriesStateLocked(conversationID string, state injectedMemoriesState) error {
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return atomicWriteFile(m.injectedMemoriesStatePath(conversationID), data, 0o644)
}

// readSessionMemoryLocked 读取会话记忆文件，返回条目列表与持久化断点 ID。文件缺失时
// 返回空值。调用方须持有 m.mu。
func (m *MarkdownMemoryStore) readSessionMemoryLocked(sessionID string) ([]storage.ConversationMemory, string) {
	data, err := os.ReadFile(m.sessionPath(sessionID))
	if err != nil {
		return nil, ""
	}
	meta, body := splitFrontMatter(string(data))
	return parseConversationMemoryBody(sessionID, body), meta["metadata.breakpoint_id"]
}

func (m *MarkdownMemoryStore) writeSessionMemoryLocked(sessionID string, items []storage.ConversationMemory, breakpointID string) error {
	sessionID = sanitizeName(sessionID)
	var b strings.Builder
	b.WriteString("---\n")
	b.WriteString("name: current-session\n")
	b.WriteString("description: 当前会话主干信息\n")
	b.WriteString("metadata:\n")
	b.WriteString("  node_type: session_memory\n")
	b.WriteString("  project: " + yamlScalar(m.projectName) + "\n")
	b.WriteString("  session_id: " + yamlScalar(sessionID) + "\n")
	b.WriteString("  originSessionId: " + yamlScalar(sessionID) + "\n")
	b.WriteString("  breakpoint_id: " + yamlScalar(breakpointID) + "\n")
	b.WriteString("---\n\n")
	for _, item := range items {
		b.WriteString("## " + strings.TrimSpace(item.Name) + "\n\n")
		if strings.TrimSpace(item.Description) != "" {
			b.WriteString(strings.TrimSpace(item.Description) + "\n\n")
		}
		b.WriteString(strings.TrimSpace(item.Body) + "\n\n")
	}
	return atomicWriteFile(m.sessionPath(sessionID), []byte(b.String()), 0o644)
}

func parseConversationMemoryBody(sessionID, body string) []storage.ConversationMemory {
	sections := strings.Split(body, "\n## ")
	var result []storage.ConversationMemory
	for _, section := range sections {
		section = strings.TrimSpace(strings.TrimPrefix(section, "## "))
		if section == "" {
			continue
		}
		parts := strings.SplitN(section, "\n", 2)
		name := strings.TrimSpace(parts[0])
		content := ""
		if len(parts) > 1 {
			content = strings.TrimSpace(parts[1])
		}
		result = append(result, storage.ConversationMemory{ConversationID: sessionID, Name: name, Body: content})
	}
	return result
}

func readFrontmatterHead(file string, maxLines int) map[string]string {
	meta := make(map[string]string)
	f, err := os.Open(file)
	if err != nil {
		return meta
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var lines []string
	count := 0
	started := false
	for scanner.Scan() {
		if count >= maxLines {
			break
		}
		count++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if !started {
			if trimmed == "---" {
				started = true
			}
			continue
		}
		if trimmed == "---" {
			break
		}
		lines = append(lines, line)
	}
	var prefix string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, ":") {
			prefix = strings.TrimSuffix(trimmed, ":") + "."
			continue
		}
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if strings.HasPrefix(line, "  ") {
			key = prefix + key
		}
		meta[key] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return meta
}

// truncateBytesAtRune 把字符串截断到不超过 max 字节，且不切断多字节字符。
func truncateBytesAtRune(s string, max int) string {
	if len(s) <= max {
		return s
	}
	cut := max
	for cut > 0 && !utf8.RuneStart(s[cut]) {
		cut--
	}
	return s[:cut]
}

func splitFrontMatter(raw string) (map[string]string, string) {
	meta := make(map[string]string)
	if !strings.HasPrefix(raw, "---\n") {
		return meta, raw
	}
	end := strings.Index(raw[4:], "\n---")
	if end < 0 {
		return meta, raw
	}
	front := raw[4 : 4+end]
	body := raw[4+end+5:]
	var prefix string
	scanner := bufio.NewScanner(strings.NewReader(front))
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if strings.HasSuffix(trimmed, ":") {
			prefix = strings.TrimSuffix(trimmed, ":") + "."
			continue
		}
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if strings.HasPrefix(line, "  ") {
			key = prefix + key
		}
		meta[key] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return meta, body
}

func filterMemories(items []storage.Memory, memType string) []storage.Memory {
	var result []storage.Memory
	for _, item := range items {
		if item.Type == memType {
			result = append(result, item)
		}
	}
	return result
}

func sanitizeName(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.TrimSpace(s) {
		allowed := r == '.' || r == '_' || r == '-' || unicode.IsLetter(r) || unicode.IsDigit(r)
		if allowed {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-._")
}

func workspaceDirName(workspaceRoot string) string {
	name, err := config.WorkspaceName(workspaceRoot)
	if err != nil {
		return sanitizeName(filepath.Base(filepath.Clean(workspaceRoot)))
	}
	return name
}

func safeRelativeMemoryPath(path string) bool {
	if strings.TrimSpace(path) == "" || filepath.IsAbs(path) {
		return false
	}
	clean := filepath.Clean(path)
	return clean == filepath.Base(clean) && !strings.HasPrefix(clean, "..") && strings.HasSuffix(clean, ".md") && clean != "memory.md"
}

func yamlScalar(s string) string {
	s = strings.ReplaceAll(strings.TrimSpace(s), "\"", "\\\"")
	return "\"" + s + "\""
}

func atomicWriteFile(path string, data []byte, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".tmp-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, perm); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

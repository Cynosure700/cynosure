package local

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
	"nano_cc/internal/idgen"
)

const memoryIndexHeader = "# Memory Index\n"

type MarkdownMemoryStore struct {
	mu          sync.Mutex
	rootDir     string
	indexPath   string
	sessionsDir string
	projectName string
}

func NewMarkdownMemoryStore(workspaceRoot string) (*MarkdownMemoryStore, error) {
	projectName := sanitizeName(filepath.Base(filepath.Clean(workspaceRoot)))
	if projectName == "" {
		projectName = "project"
	}
	baseDir, err := config.CynosureMemoryDir()
	if err != nil {
		return nil, err
	}
	root := filepath.Join(baseDir, workspaceMemoryDirName(workspaceRoot))
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
		mem.Type = "user_preference"
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

func (m *MarkdownMemoryStore) ListProjectFactMemories(ctx context.Context, userID string) ([]storage.Memory, error) {
	return m.ListMemoriesByUserAndType(ctx, userID, "project_fact")
}

func (m *MarkdownMemoryStore) CountMemoriesByUserAndType(ctx context.Context, userID, memType string) (int, error) {
	items, err := m.ListMemoriesByUserAndType(ctx, userID, memType)
	return len(items), err
}

func (m *MarkdownMemoryStore) CountProjectFactMemories(ctx context.Context, userID string) (int, error) {
	items, err := m.ListProjectFactMemories(ctx, userID)
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

func (m *MarkdownMemoryStore) ReplaceProjectFactMemories(ctx context.Context, userID string, items []storage.Memory) error {
	return m.replaceMemoriesByType("project_fact", items)
}

func (m *MarkdownMemoryStore) ListConversationMemories(ctx context.Context, sessionID string) ([]storage.ConversationMemory, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	path := m.sessionPath(sessionID)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	_, body := splitFrontMatter(string(data))
	return parseConversationMemoryBody(sessionID, body), nil
}

func (m *MarkdownMemoryStore) ReplaceConversationMemories(ctx context.Context, sessionID, userID string, items []storage.ConversationMemory) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.writeSessionMemoryLocked(sessionID, items)
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

func (m *MarkdownMemoryStore) uniqueMemoryFilenameLocked(name string) string {
	base := sanitizeName(name)
	if base == "" {
		base = "memory"
	}
	path := base + ".md"
	if _, err := os.Stat(filepath.Join(m.rootDir, path)); os.IsNotExist(err) {
		return path
	}
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

func (m *MarkdownMemoryStore) writeSessionMemoryLocked(sessionID string, items []storage.ConversationMemory) error {
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

func workspaceMemoryDirName(workspaceRoot string) string {
	return config.WorkspaceKey(workspaceRoot)
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

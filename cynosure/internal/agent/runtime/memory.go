package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/idgen"
	"nano_cc/internal/logger"
)

const (
	MemoryTypeEpisodicMemory = "episodic_memory"
	MemoryTypeUserPreference = "user_preference"
	MemoryTypeProjectFact    = "project_fact"
)

const (
	maxPreferenceMemories    = 20
	maxSessionSummaries      = 20
	sessionSummaryPruneCount = 10
	maxProjectFactMemories   = 50
	maxInjectedMemories      = 10
	maxMemoryDialogueChars   = 4000
	maxMemoryNameRunes       = 80
	maxMemoryDescRunes       = 300
	maxMemoryBodyRunes       = 2000
)

func validMemoryType(t string) bool {
	switch t {
	case MemoryTypeEpisodicMemory, MemoryTypeUserPreference, MemoryTypeProjectFact:
		return true
	}
	return false
}

func (s *Service) memoryExtractionSystemPrompt() string {
	return s.Prompts.withDefaults().MemoryExtraction
}

func (s *Service) memorySelectionSystemPrompt() string {
	return s.Prompts.withDefaults().MemorySelection
}

func (s *Service) memoryConsolidationSystemPrompt(typeLabel, typeValue string) string {
	template := s.Prompts.withDefaults().MemoryConsolidation
	template = strings.ReplaceAll(template, "{{type_label}}", typeLabel)
	template = strings.ReplaceAll(template, "{{type_value}}", typeValue)
	return template
}

type extractedMemory struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// extractMemories runs once at the end of a conversation turn: it asks the LLM
// to extract the three kinds of memories from the dialogue, persists them, and
// triggers consolidation/pruning by type. It is best-effort: failures are
// logged and swallowed so the user-facing response is never affected.
func (s *Service) extractMemories(ctx context.Context, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderDialogueForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load existing memories failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.memoryExtractionSystemPrompt()},
			{Role: "user", Content: buildExtractionUserPrompt(existing, dialogue)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: extraction failed: %v", err))
		return
	}
	if len(resp.Choices) == 0 {
		return
	}
	items := parseExtractedMemories(resp.Choices[0].Message.Content)
	if len(items) == 0 {
		return
	}

	var touchedPref, touchedSession, touchedProjectFact bool
	for _, it := range items {
		uid := user.ID
		if err := s.Store.InsertMemory(ctx, storage.Memory{
			ID:          idgen.New("mem"),
			UserID:      uid,
			Type:        it.Type,
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		}); err != nil {
			logger.Warn(fmt.Sprintf("memory: insert failed: %v", err))
			continue
		}
		switch it.Type {
		case MemoryTypeUserPreference:
			touchedPref = true
		case MemoryTypeEpisodicMemory:
			touchedSession = true
		case MemoryTypeProjectFact:
			touchedProjectFact = true
		}
	}
	if touchedPref {
		s.maybeConsolidateUserPreferences(ctx, user.ID)
	}
	if touchedSession {
		s.maybePruneSessionSummaries(ctx, user.ID)
	}
	if touchedProjectFact {
		s.maybeConsolidateProjectFacts(ctx, user.ID)
	}
}

func (s *Service) maybeConsolidateUserPreferences(ctx context.Context, userID string) {
	items, err := s.Store.ListMemoriesByUserAndType(ctx, userID, MemoryTypeUserPreference)
	if err != nil || len(items) < maxPreferenceMemories {
		return
	}
	refined := s.consolidateViaLLM(ctx, "user_preference", MemoryTypeUserPreference, items)
	if refined == nil {
		return
	}
	if err := s.Store.ReplaceMemoriesByUserAndType(ctx, userID, MemoryTypeUserPreference, refined); err != nil {
		logger.Warn(fmt.Sprintf("memory: replace user preferences failed: %v", err))
	}
}

func (s *Service) maybeConsolidateProjectFacts(ctx context.Context, userID string) {
	items, err := s.Store.ListProjectFactMemories(ctx, userID)
	if err != nil || len(items) < maxProjectFactMemories {
		return
	}
	refined := s.consolidateViaLLM(ctx, "project_fact", MemoryTypeProjectFact, items)
	if refined == nil {
		return
	}
	if err := s.Store.ReplaceProjectFactMemories(ctx, userID, refined); err != nil {
		logger.Warn(fmt.Sprintf("memory: replace project fact memories failed: %v", err))
	}
}

func (s *Service) maybePruneSessionSummaries(ctx context.Context, userID string) {
	n, err := s.Store.CountMemoriesByUserAndType(ctx, userID, MemoryTypeEpisodicMemory)
	if err != nil || n < maxSessionSummaries {
		return
	}
	if err := s.Store.DeleteOldestMemories(ctx, userID, MemoryTypeEpisodicMemory, sessionSummaryPruneCount); err != nil {
		logger.Warn(fmt.Sprintf("memory: prune session summaries failed: %v", err))
	}
}

// consolidateViaLLM feeds the full memory list to the model and parses the
// refined complete list. Returns nil on failure so the caller leaves data
// untouched.
func (s *Service) consolidateViaLLM(ctx context.Context, typeLabel, typeValue string, items []storage.Memory) []storage.Memory {
	if s.LLM == nil {
		return nil
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.memoryConsolidationSystemPrompt(typeLabel, typeValue)},
			{Role: "user", Content: buildConsolidationUserPrompt(items)},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: consolidation failed: %v", err))
		return nil
	}
	if len(resp.Choices) == 0 {
		return nil
	}
	parsed := parseExtractedMemories(resp.Choices[0].Message.Content)
	refined := make([]storage.Memory, 0, len(parsed))
	for _, it := range parsed {
		if it.Type != typeValue {
			continue
		}
		refined = append(refined, storage.Memory{
			ID:          idgen.New("mem"),
			Type:        typeValue,
			Name:        it.Name,
			Description: it.Description,
			Body:        it.Body,
		})
	}
	return refined
}

// selectRelevantMemories runs once before the conversation loop: it asks the
// LLM to pick the memories most relevant to the current context and renders
// them into the MemorySection. Best-effort: returns "" on failure or no data.
func (s *Service) selectRelevantMemories(ctx context.Context, user storage.User, history []storage.Message) string {
	if !s.EnableMemory || s.LLM == nil {
		return ""
	}
	all, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load relevant memories failed: %v", err))
		return ""
	}
	if len(all) == 0 {
		return ""
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.memorySelectionSystemPrompt()},
			{Role: "user", Content: buildSelectionUserPrompt(all, renderDialogueForMemory(history))},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: selection failed: %v", err))
		return ""
	}
	if len(resp.Choices) == 0 {
		return ""
	}
	indices := parseSelectedIDs(resp.Choices[0].Message.Content)
	selected := pickMemoriesByIndex(all, indices, maxInjectedMemories)
	return renderMemorySection(selected)
}

// parseExtractedMemories extracts a JSON array of memory objects from model
// output, tolerating surrounding code fences or prose, then validates types
// and truncates fields.
func parseExtractedMemories(raw string) []extractedMemory {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []extractedMemory
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	result := make([]extractedMemory, 0, len(parsed))
	for _, it := range parsed {
		it.Type = strings.TrimSpace(it.Type)
		it.Name = strings.TrimSpace(it.Name)
		it.Description = strings.TrimSpace(it.Description)
		it.Body = strings.TrimSpace(it.Body)
		if !validMemoryType(it.Type) || it.Name == "" || it.Body == "" {
			continue
		}
		it.Name = truncateRunes(it.Name, maxMemoryNameRunes)
		it.Description = truncateRunes(it.Description, maxMemoryDescRunes)
		it.Body = truncateRunes(it.Body, maxMemoryBodyRunes)
		result = append(result, it)
	}
	return result
}

// parseSelectedIDs parses a JSON array of integer indices from model output.
func parseSelectedIDs(raw string) []int {
	trimmed := extractJSONArray(raw)
	if trimmed == "" {
		return nil
	}
	var parsed []int
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return nil
	}
	return parsed
}

// pickMemoriesByIndex resolves indices into memories, dropping out-of-range and
// duplicate indices, capped at max.
func pickMemoriesByIndex(all []storage.Memory, indices []int, max int) []storage.Memory {
	seen := make(map[int]struct{})
	result := make([]storage.Memory, 0, len(indices))
	for _, idx := range indices {
		if idx < 0 || idx >= len(all) {
			continue
		}
		if _, exists := seen[idx]; exists {
			continue
		}
		seen[idx] = struct{}{}
		result = append(result, all[idx])
		if len(result) >= max {
			break
		}
	}
	return result
}

// renderMemorySection groups selected project-scoped memories by type and
// renders a Markdown block. Returns "" when empty.
func renderMemorySection(memories []storage.Memory) string {
	if len(memories) == 0 {
		return ""
	}
	var userLines, projectFactLines []string
	for _, m := range memories {
		switch m.Type {
		case MemoryTypeUserPreference:
			userLines = append(userLines, memoryBlock("(喜好) ", m))
		case MemoryTypeEpisodicMemory:
			userLines = append(userLines, memoryBlock("(经历) ", m))
		case MemoryTypeProjectFact:
			projectFactLines = append(projectFactLines, memoryBlock("", m))
		}
	}
	sections := make([]string, 0, 3)
	sections = append(sections, "### 当前项目记忆\n以下记忆仅适用于当前项目；不要迁移到其他项目会话。")
	if len(userLines) > 0 {
		sections = append(sections, "#### 关于用户的长期记忆\n"+strings.Join(userLines, "\n"))
	}
	if len(projectFactLines) > 0 {
		sections = append(sections, "#### 当前项目事实\n"+strings.Join(projectFactLines, "\n"))
	}
	return strings.Join(sections, "\n\n")
}

func memoryBlock(prefix string, m storage.Memory) string {
	line := "- " + prefix + memoryLine(m)
	body := strings.TrimSpace(m.Body)
	if body == "" {
		return line
	}
	return line + "\n  " + strings.ReplaceAll(body, "\n", "\n  ")
}

func memoryLine(m storage.Memory) string {
	if strings.TrimSpace(m.Description) == "" {
		return m.Name
	}
	return m.Name + "：" + m.Description
}

func buildExtractionUserPrompt(existing []storage.Memory, dialogue string) string {
	return "Existing memories:\n" + renderMemoryListForPrompt(existing) + "\n\nDialogue:\n" + dialogue
}

func buildConsolidationUserPrompt(items []storage.Memory) string {
	var b strings.Builder
	for i, m := range items {
		b.WriteString(fmt.Sprintf("[%d] (%s) %s: %s\n", i, m.Type, m.Name, m.Description))
		if strings.TrimSpace(m.Body) != "" {
			b.WriteString("    body: " + m.Body + "\n")
		}
	}
	return b.String()
}

func buildSelectionUserPrompt(all []storage.Memory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Recent conversation:\n")
	b.WriteString(dialogue)
	b.WriteString("\n\nCandidate memories:\n")
	for i, m := range all {
		b.WriteString(fmt.Sprintf("[%d] (%s) %s: %s\n", i, m.Type, m.Name, m.Description))
	}
	return b.String()
}

func renderMemoryListForPrompt(items []storage.Memory) string {
	if len(items) == 0 {
		return "(none)"
	}
	lines := make([]string, 0, len(items))
	for _, m := range items {
		lines = append(lines, fmt.Sprintf("- [%s] %s: %s", m.Type, m.Name, m.Description))
	}
	return strings.Join(lines, "\n")
}

// renderDialogueForMemory builds a plain-text transcript of user and assistant
// messages only, dropping tool noise to keep the LLM calls cheap.
func renderDialogueForMemory(history []storage.Message) string {
	var b strings.Builder
	for _, msg := range history {
		var content string
		switch msg.Role {
		case "user":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[user] ")
		case "assistant":
			content = strings.TrimSpace(msg.Content)
			if content == "" {
				continue
			}
			b.WriteString("[assistant] ")
		default:
			continue
		}
		b.WriteString(content)
		b.WriteString("\n\n")
		if b.Len() > maxMemoryDialogueChars {
			break
		}
	}
	return b.String()
}

func extractJSONArray(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if start := strings.IndexByte(trimmed, '['); start >= 0 {
		if end := strings.LastIndexByte(trimmed, ']'); end > start {
			return trimmed[start : end+1]
		}
	}
	return ""
}

func truncateRunes(s string, max int) string {
	if runes := []rune(s); len(runes) > max {
		return string(runes[:max])
	}
	return s
}

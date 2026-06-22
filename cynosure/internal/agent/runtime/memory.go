package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/agent/storage"
	"cynosure/internal/idgen"
	"cynosure/internal/logger"
	"cynosure/internal/textutil"
)

const (
	// MemoryTypePreference 记录用户长期偏好、习惯与项目相关的稳定描述。
	MemoryTypePreference = "preference"
	// MemoryTypeFeedback 记录对 Agent 行为的纠正与肯定（行为指导）。
	MemoryTypeFeedback = "feedback"
	// MemoryTypeProject 记录项目进展、决策、截止日期等不可从代码推导的项目动态。
	MemoryTypeProject = "project"
	// MemoryTypeReference 记录外部系统中信息的定位信息。
	MemoryTypeReference = "reference"
)

const (
	maxInjectedMemories      = 5
	maxMemoryTranscriptChars = 24000
	maxMemoryNameRunes       = 80
	maxMemoryDescRunes       = 300
	maxMemoryBodyRunes       = 2000
	// memoryIndexMaxLines 与存储侧的上限保持一致，仅用于注入系统提示词中的
	// 截断警告文本。
	memoryIndexMaxLines = 200
)

// consolidationTypes 列出参与定时全量去重的四类长期记忆。
var consolidationTypes = []string{
	MemoryTypePreference,
	MemoryTypeFeedback,
	MemoryTypeProject,
	MemoryTypeReference,
}

func validMemoryType(t string) bool {
	switch t {
	case MemoryTypePreference, MemoryTypeFeedback, MemoryTypeProject, MemoryTypeReference:
		return true
	}
	return false
}

func (s *Service) memoryExtractionSystemPrompt() string {
	template := s.Prompts.withDefaults().MemoryExtraction
	return strings.ReplaceAll(template, "{{current_date}}", time.Now().Format("2006-01-02"))
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

// extractMemories 在一个会话回合结束时运行一次：它请求 LLM 从对话中提取四类
// 长期记忆，将其持久化，并按类型触发去重/淘汰。它是尽力而为的：任何失败都会
// 被记录并吞掉，从而绝不影响面向用户的响应。
func (s *Service) extractMemories(ctx context.Context, user storage.User, history []storage.Message) {
	if s.LLM == nil {
		return
	}
	dialogue := renderModelHistoryForMemory(history)
	if strings.TrimSpace(dialogue) == "" {
		return
	}
	existing, err := s.Store.ListRelevantMemories(ctx, user.ID)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load existing memories failed: %v", err))
	}
	existingFiles, err := s.Store.ScanRecentMemories(ctx)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: scan existing memory files failed: %v", err))
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.memoryExtractionSystemPrompt()},
			{Role: "user", Content: buildExtractionUserPrompt(existing, existingFiles, dialogue)},
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
	}
}

// maybeRunConsolidation 定期对四类长期记忆做全量去重/淘汰：默认累计 5+ 次会话且
// 距上次运行 >=24h 时触发，触发后整类喂给 LLM 合并并替换，同时刷新 memory.md，
// 最后落盘状态。best-effort：任何失败仅记录告警，不影响用户响应。
func (s *Service) maybeRunConsolidation(ctx context.Context, user storage.User) {
	if s.LLM == nil {
		return
	}
	state, err := s.Store.LoadConsolidationState(ctx)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: load consolidation state failed: %v", err))
		return
	}
	state.SessionCount++
	minSessions := s.Cfg.MemoryConsolidationMinSessions
	if minSessions <= 0 {
		minSessions = 5
	}
	interval := s.Cfg.MemoryConsolidationInterval
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	due := state.SessionCount >= minSessions && (state.LastRunAt.IsZero() || time.Since(state.LastRunAt) >= interval)
	if !due {
		if err := s.Store.SaveConsolidationState(ctx, state); err != nil {
			logger.Warn(fmt.Sprintf("memory: save consolidation state failed: %v", err))
		}
		return
	}
	for _, memType := range consolidationTypes {
		items, err := s.Store.ListMemoriesByUserAndType(ctx, user.ID, memType)
		if err != nil || len(items) == 0 {
			continue
		}
		refined := s.consolidateViaLLM(ctx, memType, memType, items)
		if refined == nil {
			continue
		}
		if err := s.Store.ReplaceMemoriesByUserAndType(ctx, user.ID, memType, refined); err != nil {
			logger.Warn(fmt.Sprintf("memory: replace %s memories failed: %v", memType, err))
		}
	}
	state.LastRunAt = time.Now()
	state.SessionCount = 0
	if err := s.Store.SaveConsolidationState(ctx, state); err != nil {
		logger.Warn(fmt.Sprintf("memory: save consolidation state failed: %v", err))
	}
}

// consolidateViaLLM 把完整的记忆列表喂给模型，并解析出精炼后的完整列表。
// 失败时返回 nil，使调用方保持数据不变。
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

// executeMemoryTool 处理 update_memory / delete_memory 工具：直接操作记忆文件并
// 同步刷新 memory.md（不进入无状态 Dispatch）。
func (s *Service) executeMemoryTool(ctx context.Context, name, rawArgs string) (string, error) {
	var args struct {
		Path        string  `json:"path"`
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Body        *string `json:"body"`
	}
	if err := json.Unmarshal([]byte(rawArgs), &args); err != nil {
		return "", fmt.Errorf("invalid %s arguments: %w", name, err)
	}
	path := strings.TrimSpace(args.Path)
	if path == "" {
		return "", fmt.Errorf("path is required")
	}
	switch name {
	case "delete_memory":
		if err := s.Store.DeleteMemoryFile(ctx, path); err != nil {
			return "", err
		}
		if err := s.Store.ForgetInjectedMemory(ctx, path); err != nil {
			logger.Warn(fmt.Sprintf("memory: forget injected memory failed: %v", err))
		}
		return fmt.Sprintf("Deleted memory %s", path), nil
	case "update_memory":
		if args.Name == nil && args.Description == nil && args.Body == nil {
			return "", fmt.Errorf("provide at least one of name, description, or body")
		}
		update := storage.MemoryUpdate{Name: args.Name, Description: args.Description, Body: args.Body}
		newPath, err := s.Store.UpdateMemoryFile(ctx, path, update)
		if err != nil {
			return "", err
		}
		// 更新标题会重命名文件，需清除旧路径与新路径的注入记录，下一轮按最新文件重读。
		if err := s.Store.ForgetInjectedMemory(ctx, path); err != nil {
			logger.Warn(fmt.Sprintf("memory: forget injected memory failed: %v", err))
		}
		if newPath != "" && newPath != path {
			if err := s.Store.ForgetInjectedMemory(ctx, newPath); err != nil {
				logger.Warn(fmt.Sprintf("memory: forget injected memory failed: %v", err))
			}
			return fmt.Sprintf("Updated memory %s (renamed to %s)", path, newPath), nil
		}
		return fmt.Sprintf("Updated memory %s", path), nil
	default:
		return "", fmt.Errorf("unknown memory tool %s", name)
	}
}

// buildMemorySection 构造系统提示词 <memory> 段：① memory.md 索引块（受行数/字节
// 上限约束，超限附警告）；② LLM 从确定性扫描候选中精选出的记忆完整内容块（带相对
// 时间与过期说明，并按会话去重/重读替换）。Best-effort：失败返回已得到的部分。
func (s *Service) buildMemorySection(ctx context.Context, conversationID string, user storage.User, history []storage.Message) string {
	if !s.EnableMemory {
		return ""
	}
	blocks := make([]string, 0, 2)
	if indexBlock := s.renderMemoryIndexBlock(ctx); indexBlock != "" {
		blocks = append(blocks, indexBlock)
	}
	if selectedBlock := s.renderSelectedMemoriesBlock(ctx, conversationID, user, history); selectedBlock != "" {
		blocks = append(blocks, selectedBlock)
	}
	return strings.Join(blocks, "\n\n")
}

// renderMemoryIndexBlock 渲染 memory.md 索引块（需求1）。
func (s *Service) renderMemoryIndexBlock(ctx context.Context) string {
	text, truncated, totalLines := s.Store.LoadMemoryIndexForPrompt(ctx)
	if strings.TrimSpace(text) == "" || !hasMemoryIndexEntries(text) {
		return ""
	}
	var b strings.Builder
	b.WriteString("### 过往记忆索引（memory.md）\n")
	b.WriteString("以下内容来自 memory.md，仅作为过往记忆文件的索引使用。仅用于 update_memory/delete_memory 定位、更新或删除对应记忆文件；索引条目不是有效记忆内容，不得当作用户偏好、项目事实或参考资料使用。\n\n")
	if truncated {
		b.WriteString(fmt.Sprintf("WARNING: MEMORY.md is %d lines (limit: %d). Only part of it was loaded.\n", totalLines, memoryIndexMaxLines))
		b.WriteString("Keep index entries to one line under ~200 chars; move detail into topic files.\n\n")
	}
	b.WriteString(text)
	return b.String()
}

func hasMemoryIndexEntries(text string) bool {
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "- [") {
			return true
		}
	}
	return false
}

// selectRelevantMemories 从确定性扫描得到的候选集中，请 LLM 精选最多
// maxInjectedMemories 条确定有帮助的记忆，返回选中的候选（保序）。无候选或失败时
// 返回 nil。
func (s *Service) selectRelevantMemories(ctx context.Context, candidates []storage.ScannedMemory, history []storage.Message) []storage.ScannedMemory {
	if s.LLM == nil || len(candidates) == 0 {
		return nil
	}
	resp, err := s.LLM.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: s.Cfg.LLM.ModelID,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: s.memorySelectionSystemPrompt()},
			{Role: "user", Content: buildSelectionUserPrompt(candidates, renderModelHistoryForMemory(history))},
		},
	})
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: selection failed: %v", err))
		return nil
	}
	if len(resp.Choices) == 0 {
		return nil
	}
	indices := parseSelectedIDs(resp.Choices[0].Message.Content)
	return pickScannedByIndex(candidates, indices, maxInjectedMemories)
}

// renderSelectedMemoriesBlock 扫描候选、LLM 精选、读取完整内容并渲染（需求4/5）。
func (s *Service) renderSelectedMemoriesBlock(ctx context.Context, conversationID string, user storage.User, history []storage.Message) string {
	candidates, err := s.Store.ScanRecentMemories(ctx)
	if err != nil {
		logger.Warn(fmt.Sprintf("memory: scan recent memories failed: %v", err))
		return ""
	}
	if len(candidates) == 0 {
		return ""
	}
	selected := s.selectRelevantMemories(ctx, candidates, history)
	if len(selected) == 0 {
		return ""
	}

	now := time.Now()
	type renderedMemory struct {
		mem     storage.Memory
		modTime time.Time
	}
	rendered := make([]renderedMemory, 0, len(selected))
	hasStale := false
	for _, cand := range selected {
		shouldInject, err := s.Store.ShouldInjectMemory(ctx, conversationID, cand.Path, cand.ModTime)
		if err != nil {
			logger.Warn(fmt.Sprintf("memory: injection dedup failed: %v", err))
			continue
		}
		if !shouldInject {
			continue
		}
		mem, err := s.Store.ReadMemoryFile(ctx, cand.Path)
		if err != nil {
			continue
		}
		rendered = append(rendered, renderedMemory{mem: mem, modTime: cand.ModTime})
		if now.Sub(cand.ModTime) > 24*time.Hour {
			hasStale = true
		}
	}
	if len(rendered) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("### 真实有效记忆\n以下内容来自被选中的具体记忆文件，不是 memory.md 索引；它们只是可能与当前会话相关的历史上下文，仅适用于当前项目。若与当前用户描述、当前会话上下文或当前项目事实不符，必须以当前描述和当前事实为准。")
	if hasStale {
		b.WriteString("\n")
		b.WriteString("Memories are point-in-time observations, not live state — claims about code behavior or file:line citations may be outdated. Verify against current code before asserting as fact.")
	}
	for _, item := range rendered {
		b.WriteString("\n\n")
		b.WriteString(renderSelectedMemory(item.mem, humanizeRelativeTime(item.modTime, now)))
	}
	return b.String()
}

// renderSelectedMemory 渲染单条被选中记忆的完整内容，带类型与相对时间。
func renderSelectedMemory(m storage.Memory, relativeTime string) string {
	header := "#### " + memoryLine(m)
	meta := make([]string, 0, 2)
	if strings.TrimSpace(m.Type) != "" {
		meta = append(meta, m.Type)
	}
	if strings.TrimSpace(relativeTime) != "" {
		meta = append(meta, relativeTime)
	}
	if len(meta) > 0 {
		header += " （" + strings.Join(meta, "，") + "）"
	}
	body := strings.TrimSpace(m.Body)
	if body == "" {
		return header
	}
	return header + "\n" + body
}

// humanizeRelativeTime 把时间渲染为相对当前的人类可读描述，如 "47 days ago"。
func humanizeRelativeTime(t, now time.Time) string {
	if t.IsZero() {
		return ""
	}
	d := now.Sub(t)
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%d minutes ago", int(d/time.Minute))
	case d < 24*time.Hour:
		return fmt.Sprintf("%d hours ago", int(d/time.Hour))
	default:
		return fmt.Sprintf("%d days ago", int(d/(24*time.Hour)))
	}
}

// parseExtractedMemories 从模型输出中提取一个记忆对象的 JSON 数组，容忍其周围
// 的代码围栏或散文，然后校验类型并截断各字段。
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

// parseSelectedIDs 从模型输出中解析一个整数下标的 JSON 数组。
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

// pickScannedByIndex 把下标解析为扫描得到的候选项，丢弃越界与重复的下标，
// 并以 max 为上限截断。
func pickScannedByIndex(all []storage.ScannedMemory, indices []int, max int) []storage.ScannedMemory {
	seen := make(map[int]struct{})
	result := make([]storage.ScannedMemory, 0, len(indices))
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

func memoryLine(m storage.Memory) string {
	if strings.TrimSpace(m.Description) == "" {
		return m.Name
	}
	return m.Name + "：" + m.Description
}

func buildExtractionUserPrompt(existing []storage.Memory, existingFiles []storage.ScannedMemory, dialogue string) string {
	var b strings.Builder
	b.WriteString("Existing memories:\n")
	b.WriteString(renderMemoryListForPrompt(existing))
	b.WriteString("\n\n## Existing memory files\n\n")
	b.WriteString(renderMemoryFilesForPrompt(existingFiles))
	b.WriteString("\n\nCheck this list before writing — update an existing file rather than creating a duplicate.")
	b.WriteString("\n\nDialogue:\n")
	b.WriteString(dialogue)
	return b.String()
}

func renderMemoryFilesForPrompt(files []storage.ScannedMemory) string {
	if len(files) == 0 {
		return "(none)"
	}
	lines := make([]string, 0, len(files))
	for _, f := range files {
		desc := strings.TrimSpace(f.Description)
		if desc == "" {
			desc = strings.TrimSpace(f.Name)
		}
		lines = append(lines, fmt.Sprintf("- %s: %s", f.Path, desc))
	}
	return strings.Join(lines, "\n")
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

func buildSelectionUserPrompt(all []storage.ScannedMemory, dialogue string) string {
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

// renderModelHistoryForMemory 为记忆提取构建模型历史的文字记录。与纯文本对话
// 不同，它包含完整的交互：user/assistant 文本、assistant 的工具调用（名称 +
// 参数）以及工具结果（状态 + 结果），从而让记忆扎根于实际发生的事情。它受
// maxMemoryTranscriptChars 限制其长度。
func renderModelHistoryForMemory(history []storage.Message) string {
	var b strings.Builder
	for _, msg := range history {
		switch msg.Role {
		case "user":
			content := strings.TrimSpace(msg.Content)
			if content != "" {
				b.WriteString("[user] ")
				b.WriteString(content)
				b.WriteString("\n\n")
			}
		case "assistant":
			content := strings.TrimSpace(msg.Content)
			if content != "" {
				b.WriteString("[assistant] ")
				b.WriteString(content)
				b.WriteString("\n\n")
			}
			for _, call := range msg.ToolCalls {
				name := strings.TrimSpace(call.Function.Name)
				if name == "" {
					continue
				}
				args := singleLine(call.Function.Arguments)
				b.WriteString("[tool_call] ")
				b.WriteString(name)
				b.WriteString("(")
				b.WriteString(args)
				b.WriteString(")\n\n")
			}
		case "tool":
			status, result, _ := textutil.ParseToolResult(msg.Content)
			result = strings.TrimSpace(result)
			if result == "" {
				continue
			}
			b.WriteString("[tool_result] ")
			if status != "" {
				b.WriteString(status)
				b.WriteString(": ")
			}
			b.WriteString(result)
			b.WriteString("\n\n")
		default:
			continue
		}
		if b.Len() > maxMemoryTranscriptChars {
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

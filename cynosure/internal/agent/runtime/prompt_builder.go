package runtime

import (
	"context"
	"fmt"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/assets"
	"cynosure/internal/agent/runtime/compression"
	"cynosure/internal/agent/storage"
	"cynosure/internal/assistant"
	"cynosure/internal/config"
	"cynosure/internal/sessions"
	agenttools "cynosure/internal/tools"
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {
	memoryIndex := ""
	memorySection := ""
	if memoryOn {
		memoryIndex, memorySection = s.buildMemorySection(ctx, conversation.ID, user, history)
	}
	return s.buildSystemPromptWithMemory(user, snapshot, memoryIndex, memorySection)
}

func (s *Service) buildSystemPromptWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memoryIndex string, memorySection string) string {
	return assistant.BuildSystemPrompt(s.promptOptionsWithMemory(user, snapshot, memoryIndex, memorySection))
}

func (s *Service) buildSystemReminderWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memoryIndex string, memorySection string) string {
	return assistant.BuildSystemReminder(s.promptOptionsWithMemory(user, snapshot, memoryIndex, memorySection))
}

func (s *Service) promptOptionsWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memoryIndex string, memorySection string) assistant.PromptOptions {
	toolNames := []string(nil)
	if s.Tools != nil {
		toolDefs := s.Tools.Definitions()
		toolNames = make([]string, 0, len(toolDefs))
		for _, def := range toolDefs {
			if def.Function == nil || strings.TrimSpace(def.Function.Name) == "" {
				continue
			}
			toolNames = append(toolNames, def.Function.Name)
		}
	}
	loader := sessions.NewSkillLoader()
	if snapshot != nil && snapshot.Merged != nil {
		loader = snapshot.Merged
	}
	return assistant.PromptOptions{
		BasePrompt:        s.BasePrompt,
		Surface:           fmt.Sprintf("the local TUI session for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		MemoryIndex:       memoryIndex,
		MemorySection:     memorySection,
		GitStatus:         s.GitStatus,
		CurrentDate:       time.Now().Format("2006-01-02"),
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		CynosureMarkdown: assistant.CynosureMarkdownContext{
			UserPath:         s.CynosureMarkdown.UserPath,
			UserContent:      s.CynosureMarkdown.UserContent,
			WorkspacePath:    s.CynosureMarkdown.WorkspacePath,
			WorkspaceContent: s.CynosureMarkdown.WorkspaceContent,
		},
		ToolNames: toolNames,
	}
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
	// 每轮从磁盘重载用户级 / 工作区级 skill，使对话期间新增/修改/删除即时生效；
	// 内置 skill 来自启动时固定的嵌入加载器（s.BuiltinSkills），无需重载。
	userAndWorkspace, err := sessions.LoadSkillsFromDirs([]sessions.SkillDir{
		{Path: s.UserSkillsDir, Source: "user"},
		{Path: config.WorkspaceCynosureSkillsDir(s.Cfg.WorkspaceRoot), Source: "workspace"},
	})
	if err != nil {
		return nil, err
	}
	merged := sessions.MergeSkillLoaders(s.BuiltinSkills, userAndWorkspace)
	snapshot := agenttools.NewSkillSnapshot(nil, merged)
	snapshot.BuiltinMaterializer = builtinSkillMaterializer()
	return snapshot, nil
}

// builtinSkillMaterializer 返回把内置 skill 整棵目录树落盘到 ~/.cynosure/system/skills/<name>/
// 的闭包，供 load_skill 加载 builtin 来源 skill 时调用。
func builtinSkillMaterializer() func(name string) (string, error) {
	return func(name string) (string, error) {
		destRoot, err := config.CynosureSystemSkillsDir()
		if err != nil {
			return "", err
		}
		return sessions.MaterializeBuiltinSkill(assets.BuiltinSkillsFS(), name, destRoot)
	}
}

func buildOpenAIMessages(systemPrompt string, history []storage.Message) []openai.ChatCompletionMessage {
	return buildOpenAIMessagesWithReminder(systemPrompt, "", history)
}

func buildOpenAIMessagesWithReminder(systemPrompt string, systemReminder string, history []storage.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}
	if reminder := strings.TrimSpace(systemReminder); reminder != "" {
		messages = append(messages, openai.ChatCompletionMessage{Role: "user", Content: reminder})
	}
	for _, msg := range compression.RepairToolCallBoundaries(history) {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: storageToolCallsToOpenAI(msg.ToolCalls)})
	}
	return messages
}

func estimatePromptWithReminder(systemPrompt string, systemReminder string) string {
	systemPrompt = strings.TrimSpace(systemPrompt)
	systemReminder = strings.TrimSpace(systemReminder)
	if systemReminder == "" {
		return systemPrompt
	}
	if systemPrompt == "" {
		return systemReminder
	}
	return systemPrompt + "\n\n" + systemReminder
}

func storageToolCallsToOpenAI(toolCalls []storage.MessageToolCall) []openai.ToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]openai.ToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, openai.ToolCall{ID: toolCall.ID, Type: openai.ToolType(toolCall.Type), Function: openai.FunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

func openAIToolCallsToStorage(toolCalls []openai.ToolCall) []storage.MessageToolCall {
	if len(toolCalls) == 0 {
		return nil
	}
	converted := make([]storage.MessageToolCall, 0, len(toolCalls))
	for _, toolCall := range toolCalls {
		converted = append(converted, storage.MessageToolCall{ID: toolCall.ID, Type: string(toolCall.Type), Function: storage.MessageFunctionCall{Name: toolCall.Function.Name, Arguments: toolCall.Function.Arguments}})
	}
	return converted
}

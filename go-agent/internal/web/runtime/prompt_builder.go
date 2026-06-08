package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/assistant"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message) string {
	return s.buildSystemPromptWithMemory(user, snapshot, s.selectRelevantMemories(ctx, user, history))
}

func (s *Service) buildSystemPromptWithMemory(user storage.User, snapshot *agenttools.SkillSnapshot, memorySection string) string {
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
	return assistant.BuildSystemPrompt(assistant.PromptOptions{
		BasePrompt:        s.BasePrompt,
		Surface:           fmt.Sprintf("the browser chat experience for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		MemorySection:     memorySection,
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		ToolNames:         toolNames,
	})
}

func (s *Service) buildConversationSkillSnapshot(skills []storage.Skill) *agenttools.SkillSnapshot {
	return agenttools.NewSkillSnapshot(buildDBSkillLoader(skills), s.BuiltinSkills)
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
	skills, err := s.Store.ListEnabledSkillsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.buildConversationSkillSnapshot(skills), nil
}

func buildDBSkillLoader(skills []storage.Skill) *sessions.SkillLoader {
	entries := make(map[string]*sessions.SkillEntry, len(skills))
	for _, skill := range skills {
		entries[skill.Slug] = &sessions.SkillEntry{
			Meta: map[string]string{"description": skill.Description},
			Body: skill.Content,
			Path: "db://skills/" + skill.ID,
		}
	}
	loader := sessions.NewSkillLoader()
	loader.LoadFromEntries(entries)
	return loader
}

func buildOpenAIMessages(systemPrompt string, history []storage.Message) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{{Role: "system", Content: systemPrompt}}
	for _, msg := range history {
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent, ToolCallID: msg.ToolCallID, ToolCalls: storageToolCallsToOpenAI(msg.ToolCalls)})
	}
	return messages
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

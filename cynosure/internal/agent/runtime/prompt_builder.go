package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/assistant"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot, history []storage.Message, memoryOn bool) string {
	memorySection := ""
	if memoryOn {
		memorySection = s.buildMemorySection(ctx, conversation.ID, user, history)
	}
	return s.buildSystemPromptWithMemory(user, snapshot, memorySection)
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
		Surface:           fmt.Sprintf("the local TUI session for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		MemorySection:     memorySection,
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		CynosureMarkdown: assistant.CynosureMarkdownContext{
			UserPath:         s.CynosureMarkdown.UserPath,
			UserContent:      s.CynosureMarkdown.UserContent,
			WorkspacePath:    s.CynosureMarkdown.WorkspacePath,
			WorkspaceContent: s.CynosureMarkdown.WorkspaceContent,
		},
		ToolNames: toolNames,
	})
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*agenttools.SkillSnapshot, error) {
	return agenttools.NewSkillSnapshot(nil, s.BuiltinSkills), nil
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

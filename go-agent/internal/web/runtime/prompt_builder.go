package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/assistant"
	"nano_cc/internal/logger"
	"nano_cc/internal/sessions"
	agenttools "nano_cc/internal/tools"
	"nano_cc/internal/web/storage"
)

const (
	maxProfileInjectionBytes  = 4096
	recentTopicsConversations = 5
)

func (s *Service) buildSystemPrompt(ctx context.Context, conversation storage.Conversation, user storage.User, snapshot *agenttools.SkillSnapshot) string {
	return s.buildSystemPromptWithMemory(user, snapshot, s.buildMemorySection(ctx, conversation, user))
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

// buildMemorySection assembles the user profile card and recent conversation
// topics into a Markdown block for the system prompt. All reads are
// best-effort: failures degrade to an empty section.
func (s *Service) buildMemorySection(ctx context.Context, conversation storage.Conversation, user storage.User) string {
	profileJSON := ""
	if profile, ok, err := s.Store.GetUserProfile(ctx, user.ID); err != nil {
		logger.Warn(fmt.Sprintf("memory: load user profile failed: %v", err))
	} else if ok {
		profileJSON = strings.TrimSpace(profile.ProfileJSON)
	}

	var topics []string
	if recent, err := s.Store.ListRecentTopicsByUser(ctx, user.ID, conversation.ID, recentTopicsConversations); err != nil {
		logger.Warn(fmt.Sprintf("memory: load recent topics failed: %v", err))
	} else {
		topics = aggregateTopics(recent)
	}

	return renderMemorySection(profileJSON, topics)
}

func renderMemorySection(profileJSON string, topics []string) string {
	sections := make([]string, 0, 2)
	if profileJSON = strings.TrimSpace(profileJSON); profileJSON != "" {
		if len(profileJSON) > maxProfileInjectionBytes {
			profileJSON = profileJSON[:maxProfileInjectionBytes] + "\n... (truncated)"
		}
		sections = append(sections, "### 用户档案卡\n以下是关于当前用户的已保存信息，回答时请参考，但不要照搬复述：\n```json\n"+profileJSON+"\n```")
	}
	if len(topics) > 0 {
		lines := make([]string, 0, len(topics))
		for _, topic := range topics {
			lines = append(lines, "- "+topic)
		}
		sections = append(sections, "### 近期聊过的话题\n仅供参考，帮助你理解用户近期关注点；你没有这些对话的原文，如需细节请让用户补充：\n"+strings.Join(lines, "\n"))
	}
	return strings.Join(sections, "\n\n")
}

func aggregateTopics(recent []storage.ConversationTopics) []string {
	seen := make(map[string]struct{})
	var result []string
	for _, item := range recent {
		var topics []string
		if err := json.Unmarshal([]byte(item.TopicsJSON), &topics); err != nil {
			continue
		}
		for _, topic := range topics {
			topic = strings.TrimSpace(topic)
			if topic == "" {
				continue
			}
			if _, exists := seen[topic]; exists {
				continue
			}
			seen[topic] = struct{}{}
			result = append(result, topic)
		}
	}
	return result
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

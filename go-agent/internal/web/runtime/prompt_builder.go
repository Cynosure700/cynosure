package runtime

import (
	"context"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/assistant"
	"nano_cc/internal/sessions"
	"nano_cc/internal/web/storage"
)

func (s *Service) buildSystemPrompt(user storage.User, loader *sessions.SkillLoader) string {
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
	return assistant.BuildSystemPrompt(assistant.PromptOptions{
		Surface:           fmt.Sprintf("the browser chat experience for user %s", user.Username),
		SkillDescriptions: loader.GetDescriptions(),
		WorkingDirectory:  strings.TrimSpace(s.Cfg.WorkspaceRoot),
		ToolNames:         toolNames,
	})
}

func (s *Service) buildConversationSkillLoader(skills []storage.Skill) *sessions.SkillLoader {
	return sessions.MergeSkillLoaders(s.BuiltinSkills, buildDBSkillLoader(skills))
}

func (s *Service) buildSkillSnapshot(ctx context.Context, userID string) (*sessions.SkillLoader, error) {
	skills, err := s.Store.ListEnabledSkillsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return s.buildConversationSkillLoader(skills), nil
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
		messages = append(messages, openai.ChatCompletionMessage{Role: msg.Role, Content: msg.Content})
	}
	return messages
}

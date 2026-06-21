package runtime

import (
	"fmt"
	"strings"

	"nano_cc/assets"
)

type FunctionalPrompts struct {
	MemoryExtraction         string
	MemorySelection          string
	MemoryConsolidation      string
	ConversationMemoryUpdate string
	ContextSummary           string
	GeneralSubagent          string
	ExploreSubagent          string
}

func LoadFunctionalPrompts() (FunctionalPrompts, error) {
	var prompts FunctionalPrompts
	load := func(name string, target *string) error {
		prompt, err := assets.FunctionalPrompt(name)
		if err != nil {
			return fmt.Errorf("load functional prompt %s: %w", name, err)
		}
		*target = prompt
		return nil
	}
	for _, item := range []struct {
		name   string
		target *string
	}{
		{name: "memory_extraction", target: &prompts.MemoryExtraction},
		{name: "memory_selection", target: &prompts.MemorySelection},
		{name: "memory_consolidation", target: &prompts.MemoryConsolidation},
		{name: "conversation_memory_update", target: &prompts.ConversationMemoryUpdate},
		{name: "context_summary", target: &prompts.ContextSummary},
		{name: "general_subagent", target: &prompts.GeneralSubagent},
		{name: "explore_subagent", target: &prompts.ExploreSubagent},
	} {
		if err := load(item.name, item.target); err != nil {
			return FunctionalPrompts{}, err
		}
	}
	return prompts, nil
}

func (p FunctionalPrompts) withDefaults() FunctionalPrompts {
	defaults := defaultFunctionalPrompts()
	if strings.TrimSpace(p.MemoryExtraction) == "" {
		p.MemoryExtraction = defaults.MemoryExtraction
	}
	if strings.TrimSpace(p.MemorySelection) == "" {
		p.MemorySelection = defaults.MemorySelection
	}
	if strings.TrimSpace(p.MemoryConsolidation) == "" {
		p.MemoryConsolidation = defaults.MemoryConsolidation
	}
	if strings.TrimSpace(p.ConversationMemoryUpdate) == "" {
		p.ConversationMemoryUpdate = defaults.ConversationMemoryUpdate
	}
	if strings.TrimSpace(p.ContextSummary) == "" {
		p.ContextSummary = defaults.ContextSummary
	}
	if strings.TrimSpace(p.GeneralSubagent) == "" {
		p.GeneralSubagent = defaults.GeneralSubagent
	}
	if strings.TrimSpace(p.ExploreSubagent) == "" {
		p.ExploreSubagent = defaults.ExploreSubagent
	}
	return p
}

func defaultFunctionalPrompts() FunctionalPrompts {
	prompts, err := LoadFunctionalPrompts()
	if err != nil {
		panic(err)
	}
	return prompts
}

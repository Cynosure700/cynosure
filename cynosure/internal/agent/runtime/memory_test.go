package runtime

import (
	"context"
	"strings"
	"testing"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
)

func TestParseExtractedMemories_ValidAndDropsInvalid(t *testing.T) {
	raw := "好的：\n```json\n[" +
		`{"name":"偏好简洁","type":"user_preference","description":"简洁中文","body":"细节"},` +
		`{"name":"x","type":"bogus","description":"d","body":"b"},` +
		`{"name":"","type":"project_fact","description":"d","body":"b"},` +
		`{"name":"构建命令","type":"project_fact","description":"d","body":"go test ./..."}` +
		"]\n```"
	got := parseExtractedMemories(raw)
	if len(got) != 2 {
		t.Fatalf("expected 2 valid memories, got %d (%v)", len(got), got)
	}
	if got[0].Type != MemoryTypeUserPreference || got[1].Type != MemoryTypeProjectFact {
		t.Fatalf("unexpected types: %v", got)
	}
}

func TestParseExtractedMemories_TruncatesLongFields(t *testing.T) {
	longName := strings.Repeat("名", maxMemoryNameRunes+10)
	longBody := strings.Repeat("体", maxMemoryBodyRunes+10)
	raw := `[{"name":"` + longName + `","type":"project_fact","description":"d","body":"` + longBody + `"}]`
	got := parseExtractedMemories(raw)
	if len(got) != 1 {
		t.Fatalf("expected 1 memory, got %v", got)
	}
	if runes := []rune(got[0].Name); len(runes) != maxMemoryNameRunes {
		t.Fatalf("expected name truncated to %d, got %d", maxMemoryNameRunes, len(runes))
	}
	if runes := []rune(got[0].Body); len(runes) != maxMemoryBodyRunes {
		t.Fatalf("expected body truncated to %d, got %d", maxMemoryBodyRunes, len(runes))
	}
}

func TestParseExtractedMemories_InvalidReturnsNil(t *testing.T) {
	if got := parseExtractedMemories("totally not json"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
	if got := parseExtractedMemories("[]"); len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestParseSelectedIDs(t *testing.T) {
	got := parseSelectedIDs("indices: [0, 2, 5]")
	want := []int{0, 2, 5}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
	if got := parseSelectedIDs("none"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestPickMemoriesByIndex_DedupBoundsAndCap(t *testing.T) {
	all := make([]storage.Memory, 15)
	for i := range all {
		all[i] = storage.Memory{Name: string(rune('a' + i))}
	}
	indices := []int{0, 0, 99, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	got := pickMemoriesByIndex(all, indices, maxInjectedMemories)
	if len(got) != maxInjectedMemories {
		t.Fatalf("expected cap at %d, got %d", maxInjectedMemories, len(got))
	}
	if got[0].Name != "a" || got[1].Name != "b" {
		t.Fatalf("unexpected dedup/order: %v", got[:2])
	}
}

func TestRenderMemorySection_EmptyWhenNoData(t *testing.T) {
	if got := renderMemorySection(nil); got != "" {
		t.Fatalf("expected empty section, got %q", got)
	}
}

func TestRenderMemorySection_GroupsByTypeWithNameAndDesc(t *testing.T) {
	got := renderMemorySection([]storage.Memory{
		{Type: MemoryTypeUserPreference, Name: "简洁", Description: "简洁中文", Body: "不应出现的body"},
		{Type: MemoryTypeEpisodicMemory, Name: "迁移", Description: "迁移资源"},
		{Type: MemoryTypeProjectFact, Name: "构建命令", Description: "使用 go test", Body: "当前项目使用 go test ./... 验证。"},
	})
	for _, want := range []string{"当前项目记忆", "仅适用于当前项目", "(喜好) 简洁：简洁中文", "(经历) 迁移：迁移资源", "当前项目事实", "构建命令：使用 go test", "当前项目使用 go test ./... 验证。"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected section to contain %q, got %q", want, got)
		}
	}
	if !strings.Contains(got, "不应出现的body") {
		t.Fatalf("selected memory body should be injected, got %q", got)
	}
}

func TestProjectFactPromptsAreProjectScoped(t *testing.T) {
	if MemoryTypeProjectFact != "project_fact" {
		t.Fatalf("MemoryTypeProjectFact = %q", MemoryTypeProjectFact)
	}
	service := &Service{}
	extractionPrompt := service.memoryExtractionSystemPrompt()
	selectionPrompt := service.memorySelectionSystemPrompt()
	for _, forbidden := range []string{
		strings.Join([]string{"seman", "tic"}, ""),
		strings.Join([]string{"shared", " across", " users"}, ""),
		strings.Join([]string{"通用", "知识"}, ""),
	} {
		if strings.Contains(extractionPrompt, forbidden) || strings.Contains(selectionPrompt, forbidden) {
			t.Fatalf("memory prompts should not contain %q", forbidden)
		}
	}
	for _, want := range []string{"project_fact", "current project", "ONLY for the current project"} {
		if !strings.Contains(extractionPrompt, want) {
			t.Fatalf("extraction prompt should contain %q", want)
		}
	}
}

func TestMemoryConsolidationPromptReplacesTemplateValues(t *testing.T) {
	service := &Service{}
	prompt := service.memoryConsolidationSystemPrompt("project_fact", MemoryTypeProjectFact)
	for _, want := range []string{`"project_fact" memories`, `type "project_fact"`} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected consolidation prompt to contain %q, got %q", want, prompt)
		}
	}
	for _, forbidden := range []string{"{{type_label}}", "{{type_value}}"} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("expected consolidation prompt to replace %q, got %q", forbidden, prompt)
		}
	}
}

func TestSelectRelevantMemories_DisabledReturnsEmpty(t *testing.T) {
	store := &fakeStore{memories: []storage.Memory{{UserID: "u1", Type: MemoryTypeUserPreference, Name: "n", Description: "d"}}}
	llm := &fakeLLMClient{}
	service := &Service{Store: store, LLM: llm, EnableMemory: false}
	if got := service.selectRelevantMemories(context.Background(), storage.User{ID: "u1"}, nil); got != "" {
		t.Fatalf("expected empty when disabled, got %q", got)
	}
	if len(llm.reqs) != 0 {
		t.Fatalf("expected no LLM calls when disabled, got %d", len(llm.reqs))
	}
}

func TestSelectRelevantMemories_NoDataReturnsEmpty(t *testing.T) {
	store := &fakeStore{}
	llm := &fakeLLMClient{}
	service := &Service{Store: store, LLM: llm, EnableMemory: true}
	if got := service.selectRelevantMemories(context.Background(), storage.User{ID: "u1"}, nil); got != "" {
		t.Fatalf("expected empty when no memories, got %q", got)
	}
	if len(llm.reqs) != 0 {
		t.Fatalf("expected no LLM calls when no memories, got %d", len(llm.reqs))
	}
}

func TestSelectRelevantMemories_PicksAndRenders(t *testing.T) {
	store := &fakeStore{memories: []storage.Memory{
		{UserID: "u1", Type: MemoryTypeUserPreference, Name: "简洁", Description: "简洁中文"},
		{UserID: "u1", Type: MemoryTypeUserPreference, Name: "无关", Description: "无关项"},
	}}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "[0]"}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}
	got := service.selectRelevantMemories(context.Background(), storage.User{ID: "u1"}, nil)
	if !strings.Contains(got, "简洁") || strings.Contains(got, "无关") {
		t.Fatalf("expected only selected memory rendered, got %q", got)
	}
}

func TestExtractMemories_PrunesSessionSummariesOverThreshold(t *testing.T) {
	store := &fakeStore{}
	for i := 0; i < maxSessionSummaries; i++ {
		store.memories = append(store.memories, storage.Memory{UserID: "u1", Type: MemoryTypeEpisodicMemory, Name: "s", Body: "b"})
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"新经历","type":"episodic_memory","description":"d","body":"b"}]`}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}
	service.extractMemories(context.Background(), storage.User{ID: "u1"}, []storage.Message{{Role: "user", Content: "hi"}})
	if len(store.deletedOldest) != 1 {
		t.Fatalf("expected one prune call, got %d", len(store.deletedOldest))
	}
	if store.deletedOldest[0][1] != MemoryTypeEpisodicMemory || store.deletedOldest[0][2] != "10" {
		t.Fatalf("unexpected prune args: %v", store.deletedOldest[0])
	}
}

func TestExtractMemories_ConsolidatesUserPreferencesOverThreshold(t *testing.T) {
	store := &fakeStore{}
	for i := 0; i < maxPreferenceMemories; i++ {
		store.memories = append(store.memories, storage.Memory{UserID: "u1", Type: MemoryTypeUserPreference, Name: "p", Body: "b"})
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"新偏好","type":"user_preference","description":"d","body":"b"}]`}}}},
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"合并偏好","type":"user_preference","description":"d","body":"b"}]`}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}
	service.extractMemories(context.Background(), storage.User{ID: "u1"}, []storage.Message{{Role: "user", Content: "hi"}})
	if len(store.replacedMemories) == 0 {
		t.Fatalf("expected user preferences consolidated/replaced")
	}
}

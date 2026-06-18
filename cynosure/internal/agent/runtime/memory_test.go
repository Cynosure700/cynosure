package runtime

import (
	"context"
	"strings"
	"testing"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"nano_cc/internal/agent/storage"
	"nano_cc/internal/config"
)

func TestParseExtractedMemories_ValidAndDropsInvalid(t *testing.T) {
	raw := "好的：\n```json\n[" +
		`{"name":"偏好简洁","type":"preference","description":"简洁中文","body":"细节"},` +
		`{"name":"x","type":"episodic_memory","description":"d","body":"b"},` +
		`{"name":"y","type":"project_fact","description":"d","body":"b"},` +
		`{"name":"","type":"feedback","description":"d","body":"b"},` +
		`{"name":"别用全局变量","type":"feedback","description":"d","body":"规则\nWhy: x\nHow to apply: y"}` +
		"]\n```"
	got := parseExtractedMemories(raw)
	if len(got) != 2 {
		t.Fatalf("expected 2 valid memories, got %d (%v)", len(got), got)
	}
	if got[0].Type != MemoryTypePreference || got[1].Type != MemoryTypeFeedback {
		t.Fatalf("unexpected types: %v", got)
	}
}

func TestParseExtractedMemories_AcceptsAllFourTypes(t *testing.T) {
	raw := `[` +
		`{"name":"a","type":"preference","description":"d","body":"b"},` +
		`{"name":"b","type":"feedback","description":"d","body":"b"},` +
		`{"name":"c","type":"project","description":"d","body":"b"},` +
		`{"name":"d","type":"reference","description":"d","body":"b"}` +
		`]`
	got := parseExtractedMemories(raw)
	if len(got) != 4 {
		t.Fatalf("expected 4 valid memories, got %d (%v)", len(got), got)
	}
	wantTypes := []string{MemoryTypePreference, MemoryTypeFeedback, MemoryTypeProject, MemoryTypeReference}
	for i, want := range wantTypes {
		if got[i].Type != want {
			t.Fatalf("memory[%d] type = %q, want %q", i, got[i].Type, want)
		}
	}
}

func TestParseExtractedMemories_TruncatesLongFields(t *testing.T) {
	longName := strings.Repeat("名", maxMemoryNameRunes+10)
	longBody := strings.Repeat("体", maxMemoryBodyRunes+10)
	raw := `[{"name":"` + longName + `","type":"project","description":"d","body":"` + longBody + `"}]`
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

func TestRenderMemorySection_SkipsLegacyTypes(t *testing.T) {
	got := renderMemorySection([]storage.Memory{
		{Type: "episodic_memory", Name: "旧经历", Description: "应被跳过", Body: "legacy"},
		{Type: "project_fact", Name: "旧事实", Description: "应被跳过", Body: "legacy"},
	})
	if got != "" {
		t.Fatalf("expected legacy-only section to be empty, got %q", got)
	}
}

func TestRenderMemorySection_GroupsByFourTypes(t *testing.T) {
	got := renderMemorySection([]storage.Memory{
		{Type: MemoryTypePreference, Name: "简洁", Description: "简洁中文", Body: "偏好正文"},
		{Type: MemoryTypeFeedback, Name: "别用全局", Description: "纠正", Body: "规则\nWhy: x\nHow to apply: y"},
		{Type: MemoryTypeProject, Name: "里程碑", Description: "截止 2026-06-24", Body: "事实\nWhy: x\nHow to apply: y"},
		{Type: MemoryTypeReference, Name: "监控面板", Description: "在 X 平台", Body: "外部引用正文"},
		{Type: "project_fact", Name: "旧事实", Description: "应被跳过", Body: "legacy"},
	})
	for _, want := range []string{
		"当前项目记忆", "仅适用于当前项目",
		"用户喜好与约束", "简洁：简洁中文", "偏好正文",
		"行为指导", "别用全局：纠正",
		"项目动态", "里程碑：截止 2026-06-24",
		"外部引用", "监控面板：在 X 平台", "外部引用正文",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected section to contain %q, got %q", want, got)
		}
	}
	if strings.Contains(got, "旧事实") || strings.Contains(got, "legacy") {
		t.Fatalf("legacy memory should not be injected, got %q", got)
	}
}

func TestMemoryExtractionPromptDescribesFourTypesAndExclusions(t *testing.T) {
	service := &Service{}
	extractionPrompt := service.memoryExtractionSystemPrompt()
	for _, want := range []string{
		"preference", "feedback", "project", "reference",
		"current project", "ONLY for the current project",
		"Why:", "How to apply:",
		"git log", "CYNOSURE.MD", "absolute date",
	} {
		if !strings.Contains(extractionPrompt, want) {
			t.Fatalf("extraction prompt should contain %q", want)
		}
	}
	// 旧类型不应再出现在抽取提示词的允许列表中。
	if strings.Contains(extractionPrompt, "\"user_preference\"") || strings.Contains(extractionPrompt, "\"episodic_memory\"") {
		t.Fatalf("extraction prompt should not advertise legacy types")
	}
}

func TestMemoryExtractionPromptInjectsCurrentDate(t *testing.T) {
	service := &Service{}
	prompt := service.memoryExtractionSystemPrompt()
	today := time.Now().Format("2006-01-02")
	if !strings.Contains(prompt, today) {
		t.Fatalf("expected extraction prompt to contain today's date %q, got %q", today, prompt)
	}
	if strings.Contains(prompt, "{{current_date}}") {
		t.Fatalf("expected current_date placeholder to be replaced, got %q", prompt)
	}
}

func TestMemoryConsolidationPromptReplacesTemplateValues(t *testing.T) {
	service := &Service{}
	prompt := service.memoryConsolidationSystemPrompt("feedback", MemoryTypeFeedback)
	for _, want := range []string{`"feedback" memories`, `type "feedback"`} {
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
	store := &fakeStore{memories: []storage.Memory{{UserID: "u1", Type: MemoryTypePreference, Name: "n", Description: "d"}}}
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
		{UserID: "u1", Type: MemoryTypePreference, Name: "简洁", Description: "简洁中文"},
		{UserID: "u1", Type: MemoryTypePreference, Name: "无关", Description: "无关项"},
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

func TestExtractMemories_ConsolidatesFeedbackOverThreshold(t *testing.T) {
	store := &fakeStore{}
	for i := 0; i < maxFeedbackMemories; i++ {
		store.memories = append(store.memories, storage.Memory{UserID: "u1", Type: MemoryTypeFeedback, Name: "f", Body: "b"})
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"新反馈","type":"feedback","description":"d","body":"b"}]`}}}},
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"合并反馈","type":"feedback","description":"d","body":"b"}]`}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}
	service.extractMemories(context.Background(), storage.User{ID: "u1"}, []storage.Message{{Role: "user", Content: "hi"}})
	if len(store.replacedMemories) == 0 {
		t.Fatalf("expected feedback memories consolidated/replaced")
	}
}

func TestExtractMemories_ConsolidatesUserPreferencesOverThreshold(t *testing.T) {
	store := &fakeStore{}
	for i := 0; i < maxPreferenceMemories; i++ {
		store.memories = append(store.memories, storage.Memory{UserID: "u1", Type: MemoryTypePreference, Name: "p", Body: "b"})
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"新偏好","type":"preference","description":"d","body":"b"}]`}}}},
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"合并偏好","type":"preference","description":"d","body":"b"}]`}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}
	service.extractMemories(context.Background(), storage.User{ID: "u1"}, []storage.Message{{Role: "user", Content: "hi"}})
	if len(store.replacedMemories) == 0 {
		t.Fatalf("expected user preferences consolidated/replaced")
	}
}

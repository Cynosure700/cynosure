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

func TestPickScannedByIndex_DedupBoundsAndCap(t *testing.T) {
	all := make([]storage.ScannedMemory, 15)
	for i := range all {
		all[i] = storage.ScannedMemory{Name: string(rune('a' + i)), Path: string(rune('a'+i)) + ".md"}
	}
	indices := []int{0, 0, 99, -1, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	got := pickScannedByIndex(all, indices, maxInjectedMemories)
	if len(got) != maxInjectedMemories {
		t.Fatalf("expected cap at %d, got %d", maxInjectedMemories, len(got))
	}
	if got[0].Name != "a" || got[1].Name != "b" {
		t.Fatalf("unexpected dedup/order: %v", got[:2])
	}
}

func TestHumanizeRelativeTime(t *testing.T) {
	now := time.Date(2026, 6, 18, 12, 0, 0, 0, time.UTC)
	cases := []struct {
		offset time.Duration
		want   string
	}{
		{30 * time.Second, "just now"},
		{5 * time.Minute, "5 minutes ago"},
		{3 * time.Hour, "3 hours ago"},
		{47 * 24 * time.Hour, "47 days ago"},
	}
	for _, c := range cases {
		if got := humanizeRelativeTime(now.Add(-c.offset), now); got != c.want {
			t.Fatalf("offset %v: got %q, want %q", c.offset, got, c.want)
		}
	}
	if got := humanizeRelativeTime(time.Time{}, now); got != "" {
		t.Fatalf("zero time should render empty, got %q", got)
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

func TestBuildExtractionUserPromptIncludesExistingFiles(t *testing.T) {
	files := []storage.ScannedMemory{{Path: "foo.md", Description: "项目偏好"}}
	prompt := buildExtractionUserPrompt(nil, files, "对话内容")
	for _, want := range []string{
		"## Existing memory files",
		"- foo.md: 项目偏好",
		"update an existing file rather than creating a duplicate",
		"对话内容",
	} {
		if !strings.Contains(prompt, want) {
			t.Fatalf("expected extraction user prompt to contain %q, got %q", want, prompt)
		}
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

func TestBuildMemorySection_DisabledReturnsEmpty(t *testing.T) {
	store := &fakeStore{scannedMemories: []storage.ScannedMemory{{Path: "n.md", Name: "n"}}}
	llm := &fakeLLMClient{}
	service := &Service{Store: store, LLM: llm, EnableMemory: false}
	if got := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil); got != "" {
		t.Fatalf("expected empty when disabled, got %q", got)
	}
	if len(llm.reqs) != 0 {
		t.Fatalf("expected no LLM calls when disabled, got %d", len(llm.reqs))
	}
}

func TestBuildMemorySection_NoCandidatesReturnsEmpty(t *testing.T) {
	store := &fakeStore{}
	llm := &fakeLLMClient{}
	service := &Service{Store: store, LLM: llm, EnableMemory: true}
	if got := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil); got != "" {
		t.Fatalf("expected empty when no candidates, got %q", got)
	}
	if len(llm.reqs) != 0 {
		t.Fatalf("expected no LLM calls when no candidates, got %d", len(llm.reqs))
	}
}

func TestBuildMemorySection_OmitsEmptyMemoryIndex(t *testing.T) {
	store := &fakeStore{memoryIndex: "# Memory Index\n\n"}
	llm := &fakeLLMClient{}
	service := &Service{Store: store, LLM: llm, EnableMemory: true}
	got := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil)
	if got != "" {
		t.Fatalf("expected empty memory section when memory.md has no entries, got %q", got)
	}
	if len(llm.reqs) != 0 {
		t.Fatalf("expected no LLM calls when no candidates, got %d", len(llm.reqs))
	}
}

func TestBuildMemorySection_InjectsIndexAndSelectedFullContentWithStaleNote(t *testing.T) {
	old := time.Now().Add(-47 * 24 * time.Hour)
	store := &fakeStore{
		memoryIndex: "# Memory Index\n\n- [简洁](pref.md) — 简洁中文",
		scannedMemories: []storage.ScannedMemory{
			{Path: "pref.md", Name: "简洁", Description: "简洁中文", Type: MemoryTypePreference, ModTime: old},
			{Path: "other.md", Name: "无关", Description: "无关项", Type: MemoryTypePreference, ModTime: old},
		},
		memoryFiles: map[string]storage.Memory{
			"pref.md": {Type: MemoryTypePreference, Name: "简洁", Description: "简洁中文", Body: "完整正文内容"},
		},
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "[0]"}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}
	got := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil)
	for _, want := range []string{
		"过往记忆索引（memory.md）",
		"仅用于 update_memory/delete_memory 定位、更新或删除对应记忆文件",
		"索引条目不是有效记忆内容",
		"真实有效记忆",
		"以下内容来自被选中的具体记忆文件，不是 memory.md 索引",
		"只是可能与当前会话相关的历史上下文",
		"若与当前用户描述、当前会话上下文或当前项目事实不符，必须以当前描述和当前事实为准",
		"- [简洁](pref.md) — 简洁中文",
		"简洁：简洁中文",
		"完整正文内容",
		"47 days ago",
		"point-in-time observations",
	} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected memory section to contain %q, got %q", want, got)
		}
	}
	if strings.Contains(got, "无关") {
		t.Fatalf("expected only selected memory injected, got %q", got)
	}
}

func TestBuildMemorySection_DedupesWithinConversationAndRereadsOnChange(t *testing.T) {
	t0 := time.Now().Add(-time.Hour)
	store := &fakeStore{
		scannedMemories: []storage.ScannedMemory{{Path: "pref.md", Name: "简洁", Description: "d", Type: MemoryTypePreference, ModTime: t0}},
		memoryFiles:     map[string]storage.Memory{"pref.md": {Type: MemoryTypePreference, Name: "简洁", Description: "d", Body: "正文v1"}},
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "[0]"}}}},
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "[0]"}}}},
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: "[0]"}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, LLM: llm, EnableMemory: true}

	first := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil)
	if !strings.Contains(first, "正文v1") {
		t.Fatalf("expected first injection to include body, got %q", first)
	}
	// Same mtime: should be deduped (no memory block).
	second := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil)
	if strings.Contains(second, "当前项目记忆") {
		t.Fatalf("expected dedup to skip already-injected memory, got %q", second)
	}
	// File changed: new mtime + new content should re-inject.
	store.scannedMemories[0].ModTime = t0.Add(time.Minute)
	store.memoryFiles["pref.md"] = storage.Memory{Type: MemoryTypePreference, Name: "简洁", Description: "d", Body: "正文v2"}
	third := service.buildMemorySection(context.Background(), "conv", storage.User{ID: "u1"}, nil)
	if !strings.Contains(third, "正文v2") {
		t.Fatalf("expected re-read updated body on mtime change, got %q", third)
	}
}

func TestExecuteMemoryTool_UpdateAndDeleteSyncStore(t *testing.T) {
	store := &fakeStore{memoryFiles: map[string]storage.Memory{"foo.md": {Name: "old", Body: "old body"}}}
	service := &Service{Store: store, Cfg: config.AppConfig{}, EnableMemory: true}

	out, err := service.executeMemoryTool(context.Background(), "update_memory", `{"path":"foo.md","body":"new body"}`)
	if err != nil {
		t.Fatalf("update_memory error: %v", err)
	}
	if !strings.Contains(out, "foo.md") || len(store.updatedMemoryFiles) != 1 || store.updatedMemoryFiles[0] != "foo.md" {
		t.Fatalf("expected update recorded, got out=%q updated=%v", out, store.updatedMemoryFiles)
	}
	if store.memoryFiles["foo.md"].Body != "new body" {
		t.Fatalf("expected body updated, got %#v", store.memoryFiles["foo.md"])
	}
	if len(store.forgottenInjectedPaths) != 1 || store.forgottenInjectedPaths[0] != "foo.md" {
		t.Fatalf("expected update to clear injected memory state, got %v", store.forgottenInjectedPaths)
	}

	if _, err := service.executeMemoryTool(context.Background(), "update_memory", `{"path":"foo.md"}`); err == nil {
		t.Fatalf("expected error when no fields provided")
	}

	if _, err := service.executeMemoryTool(context.Background(), "delete_memory", `{"path":"foo.md"}`); err != nil {
		t.Fatalf("delete_memory error: %v", err)
	}
	if len(store.deletedMemoryFiles) != 1 || store.deletedMemoryFiles[0] != "foo.md" {
		t.Fatalf("expected delete recorded, got %v", store.deletedMemoryFiles)
	}
	if len(store.forgottenInjectedPaths) != 2 || store.forgottenInjectedPaths[1] != "foo.md" {
		t.Fatalf("expected delete to clear injected memory state, got %v", store.forgottenInjectedPaths)
	}
}

func TestMaybeRunConsolidation_TriggersAfterIntervalAndSessions(t *testing.T) {
	store := &fakeStore{
		memories: []storage.Memory{
			{UserID: "u1", Type: MemoryTypeFeedback, Name: "f", Body: "b"},
		},
		consolidationState: storage.ConsolidationState{SessionCount: 4, LastRunAt: time.Now().Add(-48 * time.Hour)},
	}
	llm := &fakeLLMClient{responses: []openai.ChatCompletionResponse{
		{Choices: []openai.ChatCompletionChoice{{Message: openai.ChatCompletionMessage{Content: `[{"name":"合并反馈","type":"feedback","description":"d","body":"b"}]`}}}},
	}}
	service := &Service{Store: store, Cfg: config.AppConfig{MemoryConsolidationInterval: 24 * time.Hour, MemoryConsolidationMinSessions: 5}, LLM: llm, EnableMemory: true}
	service.maybeRunConsolidation(context.Background(), storage.User{ID: "u1"})
	if len(store.replacedMemories) == 0 {
		t.Fatalf("expected consolidation to replace memories")
	}
	last := store.savedConsolidationStates[len(store.savedConsolidationStates)-1]
	if last.SessionCount != 0 || last.LastRunAt.IsZero() {
		t.Fatalf("expected state reset after run, got %#v", last)
	}
}

func TestMaybeRunConsolidation_SkipsBeforeThreshold(t *testing.T) {
	store := &fakeStore{
		memories:           []storage.Memory{{UserID: "u1", Type: MemoryTypeFeedback, Name: "f", Body: "b"}},
		consolidationState: storage.ConsolidationState{SessionCount: 1, LastRunAt: time.Now().Add(-48 * time.Hour)},
	}
	llm := &fakeLLMClient{}
	service := &Service{Store: store, Cfg: config.AppConfig{MemoryConsolidationInterval: 24 * time.Hour, MemoryConsolidationMinSessions: 5}, LLM: llm, EnableMemory: true}
	service.maybeRunConsolidation(context.Background(), storage.User{ID: "u1"})
	if len(store.replacedMemories) != 0 {
		t.Fatalf("expected no consolidation before threshold, got %#v", store.replacedMemories)
	}
	if len(llm.reqs) != 0 {
		t.Fatalf("expected no LLM call before threshold, got %d", len(llm.reqs))
	}
	last := store.savedConsolidationStates[len(store.savedConsolidationStates)-1]
	if last.SessionCount != 2 {
		t.Fatalf("expected session count incremented to 2, got %#v", last)
	}
}

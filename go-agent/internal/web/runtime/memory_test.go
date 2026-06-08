package runtime

import (
	"context"
	"strings"
	"testing"

	"nano_cc/internal/web/storage"
)

func TestRenderMemorySection_EmptyWhenNoData(t *testing.T) {
	if got := renderMemorySection("", nil); got != "" {
		t.Fatalf("expected empty memory section, got %q", got)
	}
}

func TestRenderMemorySection_IncludesProfileAndTopics(t *testing.T) {
	got := renderMemorySection(`{"identity":"gopher"}`, []string{"话题一", "话题二"})
	for _, want := range []string{"用户档案卡", `{"identity":"gopher"}`, "近期聊过的话题", "- 话题一", "- 话题二"} {
		if !strings.Contains(got, want) {
			t.Fatalf("expected memory section to contain %q, got %q", want, got)
		}
	}
}

func TestRenderMemorySection_TruncatesOversizedProfile(t *testing.T) {
	big := strings.Repeat("a", maxProfileInjectionBytes+100)
	got := renderMemorySection(big, nil)
	if !strings.Contains(got, "(truncated)") {
		t.Fatalf("expected truncation marker, got len=%d", len(got))
	}
}

func TestAggregateTopics_DedupAndSkipInvalid(t *testing.T) {
	recent := []storage.ConversationTopics{
		{TopicsJSON: `["a","b"]`},
		{TopicsJSON: `not-json`},
		{TopicsJSON: `["b","c"," "]`},
	}
	got := aggregateTopics(recent)
	want := []string{"a", "b", "c"}
	if len(got) != len(want) {
		t.Fatalf("expected %v, got %v", want, got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("expected %v, got %v", want, got)
		}
	}
}

func TestParseTopics_ExtractsArrayAndCaps(t *testing.T) {
	raw := "好的，话题如下：\n```json\n[\"t1\",\"t1\",\"t2\",\"t3\",\"t4\",\"t5\",\"t6\",\"t7\"]\n```"
	got := parseTopics(raw)
	if len(got) != maxTopicCount {
		t.Fatalf("expected %d topics, got %d (%v)", maxTopicCount, len(got), got)
	}
	if got[0] != "t1" || got[1] != "t2" {
		t.Fatalf("unexpected dedup/order: %v", got)
	}
}

func TestParseTopics_TruncatesLongTopic(t *testing.T) {
	long := strings.Repeat("话", maxTopicRunes+10)
	got := parseTopics(`["` + long + `"]`)
	if len(got) != 1 {
		t.Fatalf("expected 1 topic, got %v", got)
	}
	if runes := []rune(got[0]); len(runes) != maxTopicRunes {
		t.Fatalf("expected topic truncated to %d runes, got %d", maxTopicRunes, len(runes))
	}
}

func TestParseTopics_InvalidReturnsNil(t *testing.T) {
	if got := parseTopics("totally not json"); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestRenderHistoryForTopics_OnlyUserAndAssistantText(t *testing.T) {
	history := []storage.Message{
		{Role: "user", Content: "问题"},
		{Role: "tool", Content: "工具结果"},
		{Role: "assistant", Content: "回答"},
		{Role: "assistant", Content: ""},
	}
	got := renderHistoryForTopics(history)
	if !strings.Contains(got, "[user] 问题") || !strings.Contains(got, "[assistant] 回答") {
		t.Fatalf("expected user/assistant text, got %q", got)
	}
	if strings.Contains(got, "工具结果") {
		t.Fatalf("expected tool content to be dropped, got %q", got)
	}
}

func TestUpdateUserProfile_RejectsInvalidJSON(t *testing.T) {
	service := &Service{Store: &fakeStore{}}
	_, err := service.updateUserProfile(context.Background(), ToolContext{User: storage.User{ID: "u1"}}, `{"profile":"not-json{"}`)
	if err == nil {
		t.Fatalf("expected error for invalid profile JSON")
	}
}

func TestUpdateUserProfile_PersistsValidProfile(t *testing.T) {
	store := &fakeStore{}
	service := &Service{Store: store}
	_, err := service.updateUserProfile(context.Background(), ToolContext{User: storage.User{ID: "u1"}}, `{"profile":"{\"identity\":\"gopher\"}"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.userProfile.UserID != "u1" || !strings.Contains(store.userProfile.ProfileJSON, "gopher") {
		t.Fatalf("expected profile persisted, got %#v", store.userProfile)
	}
}

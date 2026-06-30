package tui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/storage"
)

// TestEditFileToolMessagePrefersEventLineStartsOverRecompute 验证：当 tool_call_done
// 事件携带 exec 时计算好的真实行号时，展示以事件值为准，而不是按当前磁盘文件回算。
func TestEditFileToolMessagePrefersEventLineStartsOverRecompute(t *testing.T) {
	dir := t.TempDir()
	// 工作区文件已被后续改动，new_text 在当前文件中落到第 9 行；事件给出真实行号第 5 行。
	moved := "x1\nx2\nx3\nx4\nx5\nx6\nx7\nx8\nconst new = \"world\"\n"
	if err := os.WriteFile(filepath.Join(dir, "foo.ts"), []byte(moved), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	app := NewModel(nil, SessionInfo{CWD: dir})
	app.width = 120
	app.generation = 1
	app.running = true
	rawArgs := `{"file_path":"foo.ts","old_text":"const old = \"hello\"","new_text":"const new = \"world\""}`

	updated, _ := app.Update(Event{Generation: 1, Name: "tool_call_done", Data: map[string]any{
		"tool_call_id":     "tool_edit",
		"tool_name":        "edit_file",
		"raw_args":         rawArgs,
		"status":           "success",
		"edit_line_starts": [][]int{{5}},
	}})
	rendered := plainTerminalText(updated.(Model).renderMessages())

	for _, want := range []string{`- 5│ const old = "hello"`, `+ 5│ const new = "world"`} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("rendered = %q, want event-provided line number %q preferred over recompute", rendered, want)
		}
	}
	if strings.Contains(rendered, `+ 9│ const new = "world"`) {
		t.Fatalf("rendered = %q, should not recompute from the moved on-disk file when the event carries line starts", rendered)
	}
}

// TestResumeUsesPersistedEditLineStartsEvenAfterFileChanged 验证：进程重启后磁盘文件
// 已被修改，/resume 仍优先采用持久化在历史里的真实行号，而不是按当前文件回算。
func TestResumeUsesPersistedEditLineStartsEvenAfterFileChanged(t *testing.T) {
	dir := t.TempDir()
	// 当前文件里 new_text 落在第 4 行，与持久化行号(第 5 行)不一致。
	changed := "a\nb\nc\nconst new = \"world\"\nd\n"
	if err := os.WriteFile(filepath.Join(dir, "foo.ts"), []byte(changed), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	history := []storage.Message{
		{Role: "user", Content: "改一行"},
		{Role: "assistant", Content: "好的", ToolCalls: []storage.MessageToolCall{
			{ID: "call_1", Type: "function", Function: storage.MessageFunctionCall{Name: "edit_file", Arguments: `{"file_path":"foo.ts","old_text":"const old = \"hello\"","new_text":"const new = \"world\""}`}},
		}},
		{Role: "tool", ToolCallID: "call_1", Content: `{"status":"success","result":"The file foo.ts has been updated successfully."}`, EditLineStarts: [][]int{{5}}},
	}
	resumer := &fakeSessionResumer{
		sessions: []storage.ResumableSession{{SessionID: "session-1", Title: "会话", UpdatedAt: time.Now(), MessageCount: len(history)}},
		conv:     storage.Conversation{ID: "conv_1", SessionID: "session-1", UserID: "local-user", Title: "会话"},
		history:  history,
	}
	app := NewModel(nil, SessionInfo{CWD: dir, User: storage.User{ID: "local-user"}, Resumer: resumer})
	app.width = 120
	app.handleSlashCommand("/resume")
	if handled := app.handleResumeSelection("1"); !handled {
		t.Fatal("resume selection was not handled")
	}

	rendered := plainTerminalText(app.renderMessages())
	for _, want := range []string{`- 5│ const old = "hello"`, `+ 5│ const new = "world"`} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("resumed render = %q, want persisted line number %q used over recompute", rendered, want)
		}
	}
	if strings.Contains(rendered, `+ 4│ const new = "world"`) {
		t.Fatalf("resumed render = %q, should not recompute from the changed on-disk file when persisted starts exist", rendered)
	}
}

package tui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"nano_cc/internal/agent/runtime"
)

func TestApprovalRequestEntersSelectingAndRendersPanel(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.width = 100
	app.running = true
	reply := make(chan runtime.ApprovalDecision, 1)

	updated, _ := app.Update(Event{Name: "approval_request", Data: approvalRequestMsg{
		req:   runtime.ApprovalRequest{ToolName: "bash", Title: "Bash command", CommandText: "curl https://x.com", Rule: "curl *"},
		reply: reply,
	}})
	model := updated.(Model)

	if !model.approving {
		t.Fatal("expected approving state after approval_request")
	}
	rendered := plainTerminalText(model.renderApprovalPanel())
	for _, want := range []string{"Bash command", "curl https://x.com", "This command requires approval", "Do you want to proceed?", "1. Yes", "Yes, and don't ask again for: curl *", "3. No"} {
		if !strings.Contains(rendered, want) {
			t.Fatalf("panel = %q, want %q", rendered, want)
		}
	}
}

func TestApprovalKeySelectsYesAndRepliesDecision(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.running = true
	reply := make(chan runtime.ApprovalDecision, 1)
	updated, _ := app.Update(Event{Name: "approval_request", Data: approvalRequestMsg{
		req:   runtime.ApprovalRequest{ToolName: "bash", CommandText: "curl x", Rule: "curl *"},
		reply: reply,
	}})
	model := updated.(Model)

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("1")})
	model = updated.(Model)

	if model.approving {
		t.Fatal("expected approving cleared after decision")
	}
	select {
	case d := <-reply:
		if d != runtime.ApprovalYes {
			t.Fatalf("decision = %v, want ApprovalYes", d)
		}
	default:
		t.Fatal("expected decision sent to reply channel")
	}
}

func TestApprovalKeyRejectReplyNo(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.running = true
	reply := make(chan runtime.ApprovalDecision, 1)
	updated, _ := app.Update(Event{Name: "approval_request", Data: approvalRequestMsg{
		req:   runtime.ApprovalRequest{ToolName: "bash", CommandText: "curl x", Rule: "curl *"},
		reply: reply,
	}})
	model := updated.(Model)

	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("3")})
	model = updated.(Model)

	select {
	case d := <-reply:
		if d != runtime.ApprovalNo {
			t.Fatalf("decision = %v, want ApprovalNo", d)
		}
	default:
		t.Fatal("expected decision sent to reply channel")
	}
}

func TestApprovalEnterUsesCursorSelection(t *testing.T) {
	app := NewModel(nil, SessionInfo{})
	app.running = true
	reply := make(chan runtime.ApprovalDecision, 1)
	updated, _ := app.Update(Event{Name: "approval_request", Data: approvalRequestMsg{
		req:   runtime.ApprovalRequest{ToolName: "bash", CommandText: "curl x", Rule: "curl *"},
		reply: reply,
	}})
	model := updated.(Model)

	// 光标下移到选项 2（don't ask again），回车确认。
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyDown})
	model = updated.(Model)
	updated, _ = model.Update(tea.KeyMsg{Type: tea.KeyEnter})
	model = updated.(Model)

	select {
	case d := <-reply:
		if d != runtime.ApprovalYesAlways {
			t.Fatalf("decision = %v, want ApprovalYesAlways", d)
		}
	default:
		t.Fatal("expected decision sent to reply channel")
	}
}

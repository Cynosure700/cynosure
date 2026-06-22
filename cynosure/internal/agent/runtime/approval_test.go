package runtime

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/config"
)

type stubApprover struct {
	decision ApprovalDecision
	called   bool
	lastReq  ApprovalRequest
}

func (s *stubApprover) Decide(ctx context.Context, req ApprovalRequest) (ApprovalDecision, error) {
	s.called = true
	s.lastReq = req
	return s.decision, nil
}

func bashToolCall(command string) openai.ToolCall {
	return openai.ToolCall{Function: openai.FunctionCall{Name: "bash", Arguments: `{"command":"` + command + `"}`}}
}

func TestApproveToolCallReadOnlyToolSkipsApproval(t *testing.T) {
	approver := &stubApprover{decision: ApprovalNo}
	svc := &Service{Cfg: config.AppConfig{WorkspaceRoot: t.TempDir()}, Approver: approver}
	tc := openai.ToolCall{Function: openai.FunctionCall{Name: "read_file", Arguments: `{"path":"a.go"}`}}

	approved, _ := svc.approveToolCall(context.Background(), tc)
	if !approved {
		t.Fatal("read_file should be approved without prompting")
	}
	if approver.called {
		t.Fatal("approver should not be called for read-only tools")
	}
}

func TestApproveToolCallBypassMode(t *testing.T) {
	root := t.TempDir()
	if err := writePermissions(root, `{"permissions":{"defaultMode":"bypassPermissions"}}`); err != nil {
		t.Fatal(err)
	}
	approver := &stubApprover{decision: ApprovalNo}
	svc := &Service{Cfg: config.AppConfig{WorkspaceRoot: root}, Approver: approver}

	approved, _ := svc.approveToolCall(context.Background(), bashToolCall("curl https://x.com"))
	if !approved {
		t.Fatal("bypass mode should approve everything")
	}
	if approver.called {
		t.Fatal("approver should not be called in bypass mode")
	}
}

func TestApproveToolCallAllowedRuleSkipsPrompt(t *testing.T) {
	root := t.TempDir()
	if err := writePermissions(root, `{"permissions":{"allowedRules":["curl *"]}}`); err != nil {
		t.Fatal(err)
	}
	approver := &stubApprover{decision: ApprovalNo}
	svc := &Service{Cfg: config.AppConfig{WorkspaceRoot: root}, Approver: approver}

	approved, _ := svc.approveToolCall(context.Background(), bashToolCall("curl https://x.com"))
	if !approved {
		t.Fatal("allowed rule should approve without prompting")
	}
	if approver.called {
		t.Fatal("approver should not be called when rule already allowed")
	}
}

func TestApproveToolCallRejected(t *testing.T) {
	approver := &stubApprover{decision: ApprovalNo}
	svc := &Service{Cfg: config.AppConfig{WorkspaceRoot: t.TempDir()}, Approver: approver}

	approved, _ := svc.approveToolCall(context.Background(), bashToolCall("curl https://x.com"))
	if approved {
		t.Fatal("expected rejection")
	}
	if !approver.called {
		t.Fatal("approver should be called for bash command")
	}
	if approver.lastReq.Rule != "curl *" {
		t.Fatalf("rule = %q, want curl *", approver.lastReq.Rule)
	}
}

func TestApproveToolCallYesAlwaysPersistsRule(t *testing.T) {
	root := t.TempDir()
	approver := &stubApprover{decision: ApprovalYesAlways}
	svc := &Service{Cfg: config.AppConfig{WorkspaceRoot: root}, Approver: approver}

	approved, _ := svc.approveToolCall(context.Background(), bashToolCall("curl https://x.com"))
	if !approved {
		t.Fatal("expected approval")
	}
	perms, err := config.LoadWorkspacePermissions(root)
	if err != nil {
		t.Fatal(err)
	}
	if !perms.Allows("curl *") {
		t.Fatalf("expected curl * persisted, got %v", perms.AllowedRules)
	}
}

func TestApproveToolCallNoApproverFallsThrough(t *testing.T) {
	svc := &Service{Cfg: config.AppConfig{WorkspaceRoot: t.TempDir()}}
	approved, _ := svc.approveToolCall(context.Background(), bashToolCall("curl https://x.com"))
	if !approved {
		t.Fatal("no approver should fall through to approved")
	}
}

func writePermissions(root, body string) error {
	path := config.WorkspaceCynosureSettingsPath(root)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(body), 0o644)
}

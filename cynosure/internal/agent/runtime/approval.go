package runtime

import (
	"context"
	"encoding/json"
	"strings"

	openai "github.com/sashabaranov/go-openai"

	"cynosure/internal/config"
	"cynosure/internal/logger"
	agenttools "cynosure/internal/tools"
)

// ApprovalDecision 表示用户对一次需审批工具调用的决定。
type ApprovalDecision int

const (
	// ApprovalYes 放行本次执行。
	ApprovalYes ApprovalDecision = iota
	// ApprovalYesAlways 放行本次执行，并持久化放行规则（不再询问同类命令）。
	ApprovalYesAlways
	// ApprovalNo 拒绝本次执行。
	ApprovalNo
)

// ApprovalRequest 是向用户发起的一次审批请求。
type ApprovalRequest struct {
	ToolName    string // 工具名，如 "bash"
	Title       string // 面板标题，如 "Bash command"
	CommandText string // 完整命令 / 文件路径
	Description string // 操作摘要（bash 的 description 参数等）
	Rule        string // "don't ask again" 的放行规则候选，如 "curl *"
}

// ApprovalDecider 由交互前端（TUI）实现，阻塞等待用户作出审批决定。
type ApprovalDecider interface {
	// Decide 阻塞直到用户作出选择；ctx 取消（如 Ctrl+C）时应返回 ApprovalNo。
	Decide(ctx context.Context, req ApprovalRequest) (ApprovalDecision, error)
}

// approveToolCall 是主/子 agent 共用的审批闸门。返回 (approved, decision)。
//   - 不需审批 / bypass / 命中放行规则 / 无审批通道 → approved=true。
//   - 需审批时调用 Approver.Decide；ApprovalNo → approved=false。
//   - ApprovalYesAlways 时持久化放行规则到工作区 settings.json。
func (s *Service) approveToolCall(ctx context.Context, tc openai.ToolCall) (bool, ApprovalDecision) {
	var args map[string]any
	if strings.TrimSpace(tc.Function.Arguments) != "" {
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
	}
	need, rule := agenttools.RequiresApproval(tc.Function.Name, args)
	if !need {
		return true, ApprovalYes
	}

	// 每次审批时实时读取工作区权限配置。
	perms, err := config.LoadWorkspacePermissions(s.Cfg.WorkspaceRoot)
	if err != nil {
		logger.Warn("approval: load workspace permissions: " + err.Error())
	}
	if perms.IsBypass() {
		return true, ApprovalYes
	}
	if perms.Allows(rule) {
		return true, ApprovalYes
	}
	if s.Approver == nil {
		return true, ApprovalYes
	}

	decision, err := s.Approver.Decide(ctx, ApprovalRequest{
		ToolName:    tc.Function.Name,
		Title:       approvalTitle(tc.Function.Name),
		CommandText: approvalCommandText(tc.Function.Name, args),
		Description: approvalDescription(args),
		Rule:        rule,
	})
	if err != nil {
		logger.Warn("approval: decider error: " + err.Error())
		return false, ApprovalNo
	}
	if decision == ApprovalNo {
		return false, ApprovalNo
	}
	if decision == ApprovalYesAlways {
		if err := config.AppendWorkspaceApprovalRule(s.Cfg.WorkspaceRoot, rule); err != nil {
			logger.Warn("approval: persist rule: " + err.Error())
		}
	}
	return true, decision
}

func approvalTitle(toolName string) string {
	switch toolName {
	case "bash":
		return "Bash command"
	case "write_file":
		return "Write file"
	case "edit_file":
		return "Edit file"
	case "multi_edit":
		return "Multi edit"
	default:
		return toolName
	}
}

func approvalCommandText(toolName string, args map[string]any) string {
	if toolName == "bash" {
		if command, ok := args["command"].(string); ok {
			return command
		}
		return ""
	}
	for _, key := range []string{"file_path", "path"} {
		if value, ok := args[key].(string); ok && strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

func approvalDescription(args map[string]any) string {
	if value, ok := args["description"].(string); ok {
		return value
	}
	return ""
}

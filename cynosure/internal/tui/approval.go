package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Cynosure700/cynosure/cynosure/internal/agent/runtime"
)

// approvalRequestMsg 由 Decide 在 respond goroutine 中通过 events 通道发给主循环，
// 携带审批请求内容与回传决定用的 reply 通道。
type approvalRequestMsg struct {
	req   runtime.ApprovalRequest
	reply chan runtime.ApprovalDecision
}

// approvalView 是渲染审批面板所需的数据快照。
type approvalView struct {
	title       string
	commandText string
	description string
	rule        string
}

// Decide 实现 runtime.ApprovalDecider。它在 respond goroutine（独立于 TUI 主循环）
// 中被调用，向主循环投递审批请求并阻塞等待用户决定。ctx 取消时返回 ApprovalNo。
func (m Model) Decide(ctx context.Context, req runtime.ApprovalRequest) (runtime.ApprovalDecision, error) {
	reply := make(chan runtime.ApprovalDecision, 1)
	// Generation 置 0：审批事件不参与代际过滤（见 Update 的 Event 分支），
	// 避免注入时的 Model 拷贝与运行期 generation 不一致导致事件被丢弃。
	select {
	case m.events <- Event{Name: "approval_request", Data: approvalRequestMsg{req: req, reply: reply}}:
	case <-ctx.Done():
		return runtime.ApprovalNo, nil
	}
	select {
	case decision := <-reply:
		return decision, nil
	case <-ctx.Done():
		return runtime.ApprovalNo, nil
	}
}

// beginApproval 进入审批选择态。
func (m *Model) beginApproval(msg approvalRequestMsg) {
	m.approving = true
	m.autoFollow = true
	m.approvalCursor = 0
	m.approvalReplies = msg.reply
	m.approvalView = approvalView{
		title:       firstNonEmpty(msg.req.Title, msg.req.ToolName),
		commandText: msg.req.CommandText,
		description: msg.req.Description,
		rule:        msg.req.Rule,
	}
}

// resolveApproval 把用户决定回传给阻塞中的 Decide，并退出选择态。
func (m *Model) resolveApproval(decision runtime.ApprovalDecision) {
	if m.approvalReplies != nil {
		m.approvalReplies <- decision
		m.approvalReplies = nil
	}
	m.approving = false
	m.approvalCursor = 0
	switch decision {
	case runtime.ApprovalNo:
		m.appendMessage("system", "已拒绝该操作，结束本轮。")
	}
}

// handleApprovalKey 处理审批选择态下的按键，返回 true 表示已消费。
func (m *Model) handleApprovalKey(key string) bool {
	switch key {
	case "up", "k":
		if m.approvalCursor > 0 {
			m.approvalCursor--
		}
		return true
	case "down", "j":
		if m.approvalCursor < 2 {
			m.approvalCursor++
		}
		return true
	case "1":
		m.resolveApproval(runtime.ApprovalYes)
		return true
	case "2":
		m.resolveApproval(runtime.ApprovalYesAlways)
		return true
	case "3":
		m.resolveApproval(runtime.ApprovalNo)
		return true
	case "enter":
		m.resolveApproval(approvalDecisionForCursor(m.approvalCursor))
		return true
	case "esc":
		m.resolveApproval(runtime.ApprovalNo)
		return true
	}
	return false
}

func approvalDecisionForCursor(cursor int) runtime.ApprovalDecision {
	switch cursor {
	case 0:
		return runtime.ApprovalYes
	case 1:
		return runtime.ApprovalYesAlways
	default:
		return runtime.ApprovalNo
	}
}

// renderApprovalPanel 渲染类似截图的审批选择面板。
func (m Model) renderApprovalPanel() string {
	width := max(20, m.messageWidth()-4)
	var b strings.Builder
	b.WriteString(approvalTitleStyle().Render(m.approvalView.title))
	b.WriteString("\n")
	if cmd := strings.TrimSpace(m.approvalView.commandText); cmd != "" {
		b.WriteString("  " + wrapText(colorizeFileReferences(cmd), width-2))
		b.WriteString("\n")
	}
	if desc := strings.TrimSpace(m.approvalView.description); desc != "" {
		b.WriteString("  " + subtleStyle().Render(wrapText(desc, width-2)))
		b.WriteString("\n")
	}
	b.WriteString("\n")
	b.WriteString(systemStyle().Render("This command requires approval"))
	b.WriteString("\n\n")
	b.WriteString("Do you want to proceed?\n")
	options := []string{
		"Yes",
		fmt.Sprintf("Yes, and don't ask again for: %s", m.approvalView.rule),
		"No",
	}
	for i, opt := range options {
		line := fmt.Sprintf("  %d. %s", i+1, opt)
		if i == m.approvalCursor {
			line = approvalSelectedStyle().Render("❯ " + fmt.Sprintf("%d. %s", i+1, opt))
		}
		b.WriteString(line)
		if i < len(options)-1 {
			b.WriteString("\n")
		}
	}
	return approvalPanelStyle().Width(width).Render(b.String())
}

func approvalPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.butter).Padding(0, 1)
}

func approvalTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.blue).Bold(true)
}

func approvalSelectedStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

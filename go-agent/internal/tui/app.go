package tui

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/sessions"
)

type SessionInfo struct {
	User         storage.User
	Conversation storage.Conversation
	CWD          string
	Resumer      SessionResumer
	Skills       []sessions.SkillSummary
	MCPServers   []mcp.ServerStatus
	SkillCount   int
	MCPToolCount int
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
}

type Message struct {
	Role    string
	Content string
}

type Model struct {
	runtime          *runtime.Service
	session          SessionInfo
	messages         []Message
	input            textarea.Model
	viewport         viewport.Model
	width            int
	height           int
	running          bool
	events           chan Event
	cancel           context.CancelFunc
	renderer         *glamour.TermRenderer
	generation       int64
	resumeSelecting  bool
	resumeCandidates []storage.ResumableSession
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "输入消息，Enter 发送，/help 查看命令"
	input.Focus()
	input.SetHeight(3)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(100))
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithAltScreen(), tea.WithContext(ctx))
	_, err := program.Run()
	return err
}

func (m Model) Init() tea.Cmd { return textarea.Blink }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.viewport.Width = msg.Width
		m.viewport.Height = max(5, msg.Height-7)
		m.input.SetWidth(msg.Width - 4)
		m.refreshViewport()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "enter":
			text := strings.TrimSpace(m.input.Value())
			if text == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			if m.resumeSelecting && m.handleResumeSelection(text) {
				m.refreshViewport()
				return m, nil
			}
			if strings.HasPrefix(text, "/") && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			// 初版不逐字展开 reasoning，避免刷屏；状态栏仍显示生成中。
		case "assistant":
			if msg.Content != "" {
				m.replaceLastAssistant(msg.Content)
			}
		case "error":
			m.appendMessage("error", msg.Content)
		case "done":
			m.running = false
		}
		m.refreshViewport()
		if m.running {
			return m, m.waitEvent()
		}
		return m, nil
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	status := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render(
		fmt.Sprintf("go-agent TUI | cwd: %s | skills:%d | mcp:%d | %s", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount, runningText(m.running)),
	)
	border := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
	return status + "\n" + border.Width(max(10, m.width-2)).Height(max(5, m.height-7)).Render(m.viewport.View()) + "\n" + m.input.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			return Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			return Event{Generation: generation, Name: "error", Content: err.Error()}
		}
		return Event{Generation: generation, Name: "done"}
	}
}

func (m Model) waitEvent() tea.Cmd {
	return func() tea.Msg { return <-m.events }
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return true
	case "/cwd":
		m.appendMessage("system", "当前工作区："+m.session.CWD)
		return true
	case "/skills":
		m.appendMessage("system", renderSkillDetails(m.session.Skills, m.session.SkillCount))
		return true
	case "/mcp":
		m.appendMessage("system", renderMCPDetails(m.session.MCPServers, m.session.MCPToolCount))
		return true
	case "/resume":
		m.startResumeSelection()
		return true
	}
	m.appendMessage("system", "未知命令："+text)
	return true
}

func (m *Model) startResumeSelection() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /resume")
		return
	}
	if m.session.Resumer == nil {
		m.appendMessage("system", "当前运行环境不支持 /resume")
		return
	}
	sessions, err := m.session.Resumer.ListResumableSessions(context.Background(), m.session.CWD)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	if len(sessions) == 0 {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "当前目录暂无可恢复的历史会话")
		return
	}
	m.resumeSelecting = true
	m.resumeCandidates = sessions
	m.appendMessage("system", renderResumeCandidates(sessions))
}

func (m *Model) handleResumeSelection(text string) bool {
	text = strings.TrimSpace(text)
	if text == "/cancel" || text == "/clear" {
		m.resumeSelecting = false
		m.resumeCandidates = nil
		if text == "/clear" {
			m.messages = nil
			m.appendMessage("system", "已清空当前 TUI 显示上下文")
		} else {
			m.appendMessage("system", "已取消恢复历史会话")
		}
		return true
	}
	idx, err := strconv.Atoi(text)
	if err != nil || idx < 1 || idx > len(m.resumeCandidates) {
		m.appendMessage("system", fmt.Sprintf("请输入 1-%d 之间的序号，或输入 /cancel 取消", len(m.resumeCandidates)))
		return true
	}
	candidate := m.resumeCandidates[idx-1]
	conv, history, err := m.session.Resumer.ResumeSession(context.Background(), candidate.SessionID, m.session.CWD, m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return true
	}
	m.session.Conversation = conv
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.messages = messagesForDisplay(history)
	m.appendMessage("system", fmt.Sprintf("已恢复历史会话：%s", conv.SessionID))
	return true
}

func renderResumeCandidates(sessions []storage.ResumableSession) string {
	var b strings.Builder
	b.WriteString("可恢复的历史会话：")
	for i, session := range sessions {
		title := strings.TrimSpace(session.Title)
		if title == "" {
			title = "TUI 会话"
		}
		updated := "unknown"
		if !session.UpdatedAt.IsZero() {
			updated = session.UpdatedAt.Local().Format(time.RFC3339)
		}
		b.WriteString(fmt.Sprintf("\n%d. %s | %s | 消息:%d | %s", i+1, updated, title, session.MessageCount, session.SessionID))
	}
	b.WriteString("\n输入序号恢复，或输入 /cancel 取消。")
	return b.String()
}

func messagesForDisplay(history []storage.Message) []Message {
	if len(history) == 0 {
		return nil
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "assistant", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content})
		}
	}
	return messages
}

func renderSkillDetails(skills []sessions.SkillSummary, fallbackCount int) string {
	count := len(skills)
	if count == 0 && fallbackCount > 0 {
		count = fallbackCount
	}
	if len(skills) == 0 {
		return fmt.Sprintf("已加载 Skills：%d 个", count)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("已加载 Skills：%d 个", len(skills)))
	for _, skill := range skills {
		b.WriteString("\n- ")
		b.WriteString(skill.Name)
		if strings.TrimSpace(skill.Source) != "" {
			b.WriteString(" [")
			b.WriteString(skill.Source)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
		if strings.TrimSpace(skill.Path) != "" {
			b.WriteString("\n  path: ")
			b.WriteString(skill.Path)
		}
	}
	return b.String()
}

func renderMCPDetails(servers []mcp.ServerStatus, toolCount int) string {
	if len(servers) == 0 {
		return fmt.Sprintf("MCP Servers：0 个，工具：%d 个", toolCount)
	}
	var b strings.Builder
	b.WriteString(fmt.Sprintf("MCP Servers：%d 个，工具：%d 个", len(servers), toolCount))
	for _, server := range servers {
		state := "disabled"
		if server.Connected {
			state = "connected"
		} else if server.LastError != "" {
			state = "failed"
		} else if server.Enabled {
			state = "pending"
		}
		b.WriteString("\n- ")
		b.WriteString(server.Name)
		b.WriteString(" [")
		b.WriteString(server.Transport)
		b.WriteString("] ")
		b.WriteString(state)
		b.WriteString(fmt.Sprintf(", tools: %d", server.ToolCount))
		if server.Command != "" {
			b.WriteString("\n  command: ")
			b.WriteString(strings.Join(append([]string{server.Command}, server.Args...), " "))
		}
		if server.URL != "" {
			b.WriteString("\n  url: ")
			b.WriteString(server.URL)
		}
		if server.LastError != "" {
			b.WriteString("\n  error: ")
			b.WriteString(server.LastError)
		}
	}
	return b.String()
}

func (m *Model) appendMessage(role, content string) {
	m.messages = append(m.messages, Message{Role: role, Content: content})
}

func (m *Model) appendAssistantDelta(delta string) {
	if len(m.messages) == 0 || m.messages[len(m.messages)-1].Role != "assistant" {
		m.appendMessage("assistant", delta)
		return
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) replaceLastAssistant(content string) {
	if len(m.messages) == 0 || m.messages[len(m.messages)-1].Role != "assistant" {
		m.appendMessage("assistant", content)
		return
	}
	m.messages[len(m.messages)-1].Content = content
}

func (m *Model) refreshViewport() {
	m.viewport.SetContent(m.renderMessages())
	m.viewport.GotoBottom()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		prefix := lipgloss.NewStyle().Bold(true).Render(msg.Role + ":")
		content := msg.Content
		if msg.Role == "assistant" && m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = strings.TrimSpace(rendered)
			}
		}
		b.WriteString(prefix)
		b.WriteString("\n")
		b.WriteString(content)
		b.WriteString("\n\n")
	}
	return b.String()
}

func runningText(running bool) string {
	if running {
		return "generating"
	}
	return "ready"
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

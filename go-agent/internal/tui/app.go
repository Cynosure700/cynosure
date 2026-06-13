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
	"github.com/charmbracelet/x/ansi"

	"nano_cc/internal/agent/mcp"
	"nano_cc/internal/agent/runtime"
	"nano_cc/internal/agent/storage"
	"nano_cc/internal/logger"
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
	Role             string
	Content          string
	ReasoningContent string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	cyan     lipgloss.Color
	mint     lipgloss.Color
	lavender lipgloss.Color
	butter   lipgloss.Color
	coral    lipgloss.Color
}

var tuiPalette = palette{
	ink:      lipgloss.Color("255"),
	muted:    lipgloss.Color("244"),
	panel:    lipgloss.Color("238"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("210"),
}

const inputCursor = "█"

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
	toolCallCount    int
	contextTokens    int
	contextBudget    int
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	input := textarea.New()
	input.Placeholder = "问 go-agent 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	program := tea.NewProgram(NewModel(runtimeService, session), tea.WithContext(ctx), tea.WithAltScreen())
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
		m.viewport.Height = m.viewportHeight()
		m.input.SetWidth(max(20, msg.Width-4))
		m.renderer = newMarkdownRenderer(m.messageWidth())
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
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "assistant_delta":
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.updateMetaFromData(msg.Data)
			content := msg.Content
			if content == "" && msg.Data != nil {
				content = eventContent(msg.Data)
			}
			if content != "" {
				m.replaceLastAssistant(content, eventString(msg.Data, "reasoning_content"))
			}
		case "meta":
			m.updateMetaFromData(msg.Data)
		case "error":
			m.appendMessage("error", msg.Content)
			m.running = false
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
	m.refreshViewport()
	return m, cmd
}

func (m Model) View() string {
	return m.viewport.View()
}

func (m Model) respond(ctx context.Context, text string, generation int64) tea.Cmd {
	return func() tea.Msg {
		if m.runtime == nil {
			m.events <- Event{Generation: generation, Name: "error", Content: "runtime 未初始化"}
			return nil
		}
		_, err := m.runtime.RespondToConversation(ctx, m.session.Conversation, m.session.User, text, NewEventWriter(m.events, generation))
		if err != nil {
			m.events <- Event{Generation: generation, Name: "error", Content: err.Error()}
			return nil
		}
		m.events <- Event{Generation: generation, Name: "done"}
		return nil
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
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
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
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.appendMessage("assistant", delta)
		return
	}
	if m.messages[len(m.messages)-1].Role == "thinking" {
		m.messages[len(m.messages)-1].Role = "assistant"
	}
	m.messages[len(m.messages)-1].Content += delta
}

func (m *Model) appendThinkingDelta(delta string) {
	if delta == "" {
		return
	}
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", ReasoningContent: delta})
		return
	}
	m.messages[len(m.messages)-1].ReasoningContent += delta
}

func (m *Model) replaceLastAssistant(content, reasoning string) {
	if len(m.messages) == 0 || !isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = append(m.messages, Message{Role: "assistant", Content: content, ReasoningContent: reasoning})
		return
	}
	m.messages[len(m.messages)-1].Role = "assistant"
	m.messages[len(m.messages)-1].Content = content
	m.messages[len(m.messages)-1].ReasoningContent = reasoning
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	m.viewport.GotoBottom()
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(1, m.height)
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderWelcome())
	if len(m.messages) > 0 {
		b.WriteString("\n\n")
		b.WriteString(m.renderMessages())
	}
	b.WriteString("\n")
	b.WriteString(m.renderInput())
	b.WriteString("\n")
	b.WriteString(subtleStyle().Render(m.renderLiveStatus()))
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg Message) string {
	switch msg.Role {
	case "user":
		return promptLineStyle().Render("›") + " " + userStyle().Render(wrapText(msg.Content, m.messageWidth()-2))
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(strings.TrimSpace(rendered), m.messageWidth())
			}
		} else {
			content = wrapText(content, m.messageWidth())
		}
		if m.running && strings.TrimSpace(msg.ReasoningContent) != "" {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(strings.TrimSpace(msg.ReasoningContent), m.messageWidth()-2)) + "\n" + content
		}
		return assistantLeadStyle().Render("go-agent") + "\n" + content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(msg.Content, m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(msg.Content, m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(msg.Content, m.messageWidth()-2))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(msg.Content, m.messageWidth())
	}
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	renderer, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(max(10, width)))
	return renderer
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func (m Model) renderWelcome() string {
	art := strings.Join([]string{
		`   /\_/\\`,
		`  ( o.o )   nano, but cozy`,
		`   > ^ <    ask · think · build`,
	}, "\n")
	intro := "像 Claude Code 一样一问一答：在终端直接提问，回答会流式显示在下方。"
	quick := "快捷键：Enter 发送 · Ctrl+C 中断/退出 · /resume 恢复 · /clear 清屏"
	stats := fmt.Sprintf("Skills %d · MCP tools %d", m.session.SkillCount, m.session.MCPToolCount)
	return startupPanelStyle().Width(max(20, m.width-2)).Render(accentArtStyle().Render(art) + "\n\n" + intro + "\n" + subtleStyle().Render(quick) + "\n" + subtleStyle().Render(stats))
}

func (m Model) renderHeader() string {
	width := max(20, m.width)
	status := runningText(m.running)
	contextText := "上下文 --"
	if m.contextBudget > 0 {
		contextText = fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget))
	}
	left := titleStyle().Render("✦ go-agent")
	right := headerMetaStyle().Render(fmt.Sprintf("%s  ·  本轮工具 %d  ·  %s", status, m.toolCallCount, contextText))
	gap := max(1, width-lipgloss.Width(left)-lipgloss.Width(right)-4)
	line := left + strings.Repeat(" ", gap) + right
	workspace := subtleStyle().Render(fmt.Sprintf("cwd %s · skills %d · mcp tools %d", m.session.CWD, m.session.SkillCount, m.session.MCPToolCount))
	return headerStyle().Width(width).Render(line + "\n" + workspace)
}

func (m Model) renderConversationFrame() string {
	return conversationStyle().Width(max(10, m.width)).Render(m.viewport.View())
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := strings.TrimSpace(m.input.Value())
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		text = userStyle().Render(text) + inputPromptStyle().Render(inputCursor)
	}
	return inputLineStyle().Width(max(10, m.width)).Render(prompt + " " + text)
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	parts := []string{"Enter 发送", "Ctrl+C 中断/退出", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
	if m.contextBudget > 0 {
		parts = append(parts, fmt.Sprintf("上下文 %d%% · %s/%s", min(100, m.contextTokens*100/m.contextBudget), compactNumber(m.contextTokens), compactNumber(m.contextBudget)))
	} else {
		parts = append(parts, "上下文 --")
	}
	return strings.Join(parts, " · ")
}

func (m *Model) updateMetaFromData(data any) {
	if count, ok := eventInt(data, "tool_call_count"); ok {
		m.toolCallCount = count
	}
	if tokens, ok := eventInt(data, "context_tokens"); ok {
		m.contextTokens = tokens
	}
	if budget, ok := eventInt(data, "context_budget"); ok {
		m.contextBudget = budget
	}
}

func eventString(data any, key string) string {
	m, ok := data.(map[string]any)
	if !ok {
		return ""
	}
	value, _ := m[key].(string)
	return value
}

func eventInt(data any, key string) (int, bool) {
	m, ok := data.(map[string]any)
	if !ok {
		return 0, false
	}
	switch value := m[key].(type) {
	case int:
		return value, true
	case int64:
		return int(value), true
	case float64:
		return int(value), true
	default:
		return 0, false
	}
}

func roleLabel(label string, color lipgloss.Color) string {
	return lipgloss.NewStyle().Bold(true).Foreground(color).Render(label)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(tuiPalette.panel)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.mint)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle().Padding(1, 2, 0, 2)
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.ink)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func promptLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.cyan).Bold(true)
}

func assistantLeadStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func accentArtStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.lavender).Bold(true)
}

func inputPromptStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.mint).Bold(true)
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Background(lipgloss.Color("238")).Foreground(tuiPalette.ink)
}

func compactNumber(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%dk", n/1000)
	}
	return strconv.Itoa(n)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

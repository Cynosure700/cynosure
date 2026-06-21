package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
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
	ModelID      string
}

type SessionResumer interface {
	ListResumableSessions(ctx context.Context, workspaceRoot string) ([]storage.ResumableSession, error)
	ResumeSession(ctx context.Context, sessionID, currentWorkspace string, user storage.User) (storage.Conversation, []storage.Message, error)
	StartNewConversation(ctx context.Context, user storage.User) (storage.Conversation, error)
}

type Message struct {
	Role             string
	Content          string
	ReasoningContent string
	ToolCall         *ToolCallView
}

type ToolCallView struct {
	ID            string
	Name          string
	RawArgs       string
	ArgsPreview   string
	Status        string
	ResultPreview string
}

type palette struct {
	ink      lipgloss.Color
	muted    lipgloss.Color
	panel    lipgloss.Color
	panelDim lipgloss.Color
	blue     lipgloss.Color
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
	panelDim: lipgloss.Color("235"),
	blue:     lipgloss.Color("39"),
	cyan:     lipgloss.Color("81"),
	mint:     lipgloss.Color("120"),
	lavender: lipgloss.Color("183"),
	butter:   lipgloss.Color("229"),
	coral:    lipgloss.Color("209"),
}

const inputCursor = "█"

type Model struct {
	runtime           *runtime.Service
	session           SessionInfo
	messages          []Message
	input             textarea.Model
	viewport          viewport.Model
	width             int
	height            int
	running           bool
	events            chan Event
	cancel            context.CancelFunc
	renderer          *glamour.TermRenderer
	generation        int64
	resumeSelecting   bool
	resumeCandidates  []storage.ResumableSession
	approving         bool
	approvalView      approvalView
	approvalCursor    int
	approvalReplies   chan runtime.ApprovalDecision
	toolCallCount     int
	contextTokens     int
	contextBudget     int
	autoFollow        bool
	thinkingStartedAt time.Time
	thinkingNow       time.Time
	answerStarted     bool
}

type thinkingTickMsg struct {
	generation int64
	at         time.Time
}

func NewModel(runtimeService *runtime.Service, session SessionInfo) Model {
	if strings.TrimSpace(session.ModelID) == "" && runtimeService != nil {
		session.ModelID = runtimeService.Cfg.LLM.ModelID
	}
	input := textarea.New()
	input.Placeholder = "问 cynosure 一件事..."
	input.Focus()
	input.SetHeight(2)
	input.ShowLineNumbers = false
	vp := viewport.New(100, 20)
	renderer := newMarkdownRenderer(100)
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer, autoFollow: true}
}

func Run(ctx context.Context, runtimeService *runtime.Service, session SessionInfo) error {
	previousConsole := logger.SetConsoleEnabled(false)
	defer logger.SetConsoleEnabled(previousConsole)
	config := newRunConfig(ctx, NewModel(runtimeService, session))
	if runtimeService != nil {
		// 注入审批前端：Decide 通过共享的 events 通道与主循环交互（值拷贝共享 channel）。
		runtimeService.SetApprover(config.model)
	}
	program := tea.NewProgram(config.model, config.options...)
	_, err := program.Run()
	return err
}

type runConfig struct {
	model           Model
	options         []tea.ProgramOption
	altScreen       bool
	mouseCellMotion bool
}

func newRunConfig(ctx context.Context, model Model) runConfig {
	return runConfig{model: model, options: []tea.ProgramOption{tea.WithContext(ctx), tea.WithAltScreen()}, altScreen: true}
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
			if m.approving {
				m.resolveApproval(runtime.ApprovalNo)
			}
			if m.running && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				m.appendMessage("system", "已中断当前生成")
				return m, nil
			}
			return m, tea.Quit
		case "esc":
			if m.running && !m.approving && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
				return m, nil
			}
		}
		if m.approving {
			if m.handleApprovalKey(msg.String()) {
				m.refreshViewport()
				if m.running {
					return m, m.waitEvent()
				}
				return m, nil
			}
			return m, nil
		}
		switch msg.String() {
		case "enter":
			text := m.input.Value()
			if strings.TrimSpace(text) == "" || m.running {
				return m, nil
			}
			m.input.Reset()
			m.autoFollow = true
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
			m.thinkingStartedAt = time.Now()
			m.thinkingNow = m.thinkingStartedAt
			m.answerStarted = false
			m.toolCallCount = 0
			m.generation++
			generation := m.generation
			ctx, cancel := context.WithCancel(context.Background())
			m.cancel = cancel
			m.refreshViewport()
			return m, tea.Batch(m.waitEvent(), m.respond(ctx, text, generation), m.thinkingTick(generation))
		}
		if isTerminalProbeResponseInput(msg) {
			return m, nil
		}
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case tea.MouseMsg:
		if cmd, ok := m.updateViewportScroll(msg); ok {
			return m, cmd
		}
	case thinkingTickMsg:
		if msg.generation != m.generation || !m.running {
			return m, nil
		}
		m.thinkingNow = msg.at
		m.refreshViewport()
		return m, m.thinkingTick(msg.generation)
	case Event:
		if msg.Generation != 0 && msg.Generation != m.generation {
			if m.running {
				return m, m.waitEvent()
			}
			return m, nil
		}
		switch msg.Name {
		case "approval_request":
			if req, ok := msg.Data.(approvalRequestMsg); ok {
				m.beginApproval(req)
			}
		case "assistant_delta":
			m.answerStarted = true
			m.appendAssistantDelta(msg.Content)
		case "reasoning_delta":
			m.appendThinkingDelta(msg.Content)
		case "assistant":
			m.answerStarted = true
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
		case "tool_call_start":
			m.appendToolCallStart(msg.Data)
		case "tool_call_done":
			m.updateToolCallDone(msg.Data)
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

func (m *Model) updateViewportScroll(msg tea.Msg) (tea.Cmd, bool) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyPgDown:
			m.viewport.GotoBottom()
			m.autoFollow = true
			return nil, true
		case tea.KeyPgUp, tea.KeyUp, tea.KeyDown:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			m.autoFollow = m.viewport.AtBottom()
			return cmd, true
		}
	case tea.MouseMsg:
		if msg.Action == tea.MouseActionPress {
			switch msg.Button {
			case tea.MouseButtonWheelUp, tea.MouseButtonWheelDown, tea.MouseButtonWheelLeft, tea.MouseButtonWheelRight:
				var cmd tea.Cmd
				m.viewport, cmd = m.viewport.Update(msg)
				m.autoFollow = m.viewport.AtBottom()
				return cmd, true
			}
		}
	}
	return nil, false
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.renderConversationFrame(), m.renderInputArea())
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

func (m Model) thinkingTick(generation int64) tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return thinkingTickMsg{generation: generation, at: t}
	})
}

func (m *Model) handleSlashCommand(text string) bool {
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", "命令：/help /clear /cwd /skills /mcp /resume。Enter 发送，Ctrl+C 中断或退出。")
		return true
	case "/clear":
		m.startNewSession()
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

func (m *Model) startNewSession() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Ctrl+C 中断后再执行 /clear")
		return
	}
	if m.session.Resumer == nil {
		m.messages = nil
		m.resumeSelecting = false
		m.resumeCandidates = nil
		m.appendMessage("system", "已清空当前 TUI 显示上下文")
		return
	}
	conv, err := m.session.Resumer.StartNewConversation(context.Background(), m.session.User)
	if err != nil {
		m.appendMessage("error", err.Error())
		return
	}
	m.session.Conversation = conv
	m.messages = nil
	m.resumeSelecting = false
	m.resumeCandidates = nil
	m.toolCallCount = 0
	m.contextTokens = 0
	m.appendMessage("system", "已开启全新对话，上下文已清空")
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
			m.startNewSession()
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

func (m *Model) appendToolCallStart(data any) {
	tool := toolCallViewFromEvent(data)
	if tool.ID == "" && tool.Name == "" {
		return
	}
	if tool.Status == "" {
		tool.Status = "running"
	}
	m.messages = append(m.messages, Message{Role: "tool", ToolCall: &tool})
}

func (m *Model) updateToolCallDone(data any) {
	tool := toolCallViewFromEvent(data)
	if tool.ID == "" && tool.Name == "" {
		return
	}
	for i := len(m.messages) - 1; i >= 0; i-- {
		if m.messages[i].Role != "tool" || m.messages[i].ToolCall == nil {
			continue
		}
		if tool.ID != "" && m.messages[i].ToolCall.ID == tool.ID {
			m.messages[i].ToolCall.Name = firstNonEmpty(tool.Name, m.messages[i].ToolCall.Name)
			m.messages[i].ToolCall.RawArgs = firstNonEmpty(tool.RawArgs, m.messages[i].ToolCall.RawArgs)
			m.messages[i].ToolCall.ArgsPreview = firstNonEmpty(tool.ArgsPreview, m.messages[i].ToolCall.ArgsPreview)
			m.messages[i].ToolCall.Status = firstNonEmpty(tool.Status, m.messages[i].ToolCall.Status)
			m.messages[i].ToolCall.ResultPreview = tool.ResultPreview
			return
		}
	}
	m.messages = append(m.messages, Message{Role: "tool", ToolCall: &tool})
}

func toolCallViewFromEvent(data any) ToolCallView {
	return ToolCallView{
		ID:            eventString(data, "tool_call_id"),
		Name:          eventString(data, "tool_name"),
		RawArgs:       eventString(data, "raw_args"),
		ArgsPreview:   eventString(data, "args_preview"),
		Status:        eventString(data, "status"),
		ResultPreview: eventString(data, "result_preview"),
	}
}

func firstNonEmpty(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant" || role == "thinking"
}

func (m *Model) refreshViewport() {
	shouldFollow := m.autoFollow || m.viewport.AtBottom()
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.viewport.SetContent(m.renderTranscript())
	if shouldFollow {
		m.viewport.GotoBottom()
	}
}

func (m Model) viewportHeight() int {
	if m.height <= 0 {
		return 20
	}
	return max(0, m.height-m.inputAreaHeight())
}

func (m Model) headerHeight() int {
	return lipgloss.Height(m.renderHeader())
}

func (m Model) inputAreaHeight() int {
	return lipgloss.Height(m.renderInputArea())
}

func (m Model) renderTranscript() string {
	var b strings.Builder
	b.WriteString(m.renderHeader())
	if welcome := m.renderWelcome(); welcome != "" {
		b.WriteString("\n\n")
		b.WriteString(welcome)
	}
	if len(m.messages) > 0 {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(m.renderMessages())
	}
	if m.approving {
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(m.renderApprovalPanel())
	}
	return b.String()
}

func (m Model) renderMessages() string {
	var b strings.Builder
	lastUserIndex := m.lastUserMessageIndex()
	for i, msg := range m.messages {
		if m.shouldHideMessageAt(i, msg) {
			continue
		}
		b.WriteString(m.renderMessageAt(i, msg, m.shouldShowReasoningAt(i, msg, lastUserIndex)))
		b.WriteString("\n\n")
	}
	if indicator := m.renderThinkingIndicator(); indicator != "" {
		b.WriteString(indicator)
		b.WriteString("\n\n")
	}
	return b.String()
}

func (m Model) renderThinkingIndicator() string {
	if !m.running || m.answerStarted || m.thinkingStartedAt.IsZero() {
		return ""
	}
	now := m.thinkingNow
	if now.IsZero() || now.Before(m.thinkingStartedAt) {
		now = m.thinkingStartedAt
	}
	elapsedSeconds := int(now.Sub(m.thinkingStartedAt) / time.Second)
	if elapsedSeconds < 1 {
		elapsedSeconds = 1
	}
	return thinkingIndicatorStyle().Render(fmt.Sprintf("* Thinking... (%ds)", elapsedSeconds))
}

func (m Model) lastUserMessageIndex() int {
	for i := len(m.messages) - 1; i >= 0; i-- {
		if m.messages[i].Role == "user" {
			return i
		}
	}
	return -1
}

func (m Model) shouldShowReasoningAt(index int, msg Message, lastUserIndex int) bool {
	return m.running &&
		index > lastUserIndex &&
		isLiveAssistantRole(msg.Role) &&
		strings.TrimSpace(msg.Content) == "" &&
		strings.TrimSpace(msg.ReasoningContent) != "" &&
		!m.hasAssistantContentAfter(index)
}

func (m Model) shouldHideMessageAt(index int, msg Message) bool {
	if index < 0 || !m.hasAssistantContentAfter(index) {
		return false
	}
	if msg.Role == "tool" || msg.Role == "thinking" {
		return true
	}
	return msg.Role == "assistant" && strings.TrimSpace(msg.Content) == "" && strings.TrimSpace(msg.ReasoningContent) != ""
}

func (m Model) hasAssistantContentAfter(index int) bool {
	for i := index + 1; i < len(m.messages); i++ {
		if isLiveAssistantRole(m.messages[i].Role) && strings.TrimSpace(m.messages[i].Content) != "" {
			return true
		}
	}
	return false
}

func (m Model) renderMessage(msg Message) string {
	return m.renderMessageAt(-1, msg, m.running && strings.TrimSpace(msg.ReasoningContent) != "")
}

func (m Model) renderMessageAt(index int, msg Message, showReasoning bool) string {
	switch msg.Role {
	case "user":
		return renderSelectedUserMessage(msg.Content, m.messageWidth())
	case "assistant":
		content := msg.Content
		if m.renderer != nil {
			if rendered, err := m.renderer.Render(content); err == nil {
				content = wrapText(colorizeFileReferencesWithRestore(strings.TrimSpace(rendered), ansiForeground(tuiPalette.ink)), m.messageWidth())
			}
		} else {
			content = wrapText(colorizeFileReferencesWithRestore(content, ansiForeground(tuiPalette.ink)), m.messageWidth())
		}
		if showReasoning {
			content = thinkingStyle().Render("✽ 思考中\n"+wrapText(colorizeFileReferencesWithRestore(strings.TrimSpace(msg.ReasoningContent), ansiThinking()), m.messageWidth()-2)) + "\n" + content
		}
		return content
	case "thinking":
		return thinkingStyle().Render("✽ 思考中\n" + wrapText(colorizeFileReferencesWithRestore(msg.Content, ansiThinking()), m.messageWidth()-2))
	case "system":
		return systemStyle().Render("• " + wrapText(colorizeFileReferencesWithRestore(msg.Content, ansiForeground(tuiPalette.butter)), m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(colorizeFileReferencesWithRestore(msg.Content, ansiForeground(tuiPalette.coral)), m.messageWidth()-2))
	case "tool":
		return m.renderToolMessage(msg, index >= 0 && m.hasAssistantContentAfter(index))
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(colorizeFileReferences(msg.Content), m.messageWidth())
	}
}

func (m Model) renderToolMessage(msg Message, hideResult bool) string {
	if msg.ToolCall == nil {
		return renderToolBullet() + toolStyleForStatus("").Render("⏺ Tool")
	}
	tool := msg.ToolCall
	status := strings.TrimSpace(tool.Status)
	if status == "" {
		status = "running"
	}
	if isTodoWriteTool(tool.Name) {
		return m.renderTodoWriteToolMessage(tool, status)
	}
	icon := toolIcon(status)
	name := displayToolName(tool.Name, tool.RawArgs)
	line := icon + " " + name
	if strings.TrimSpace(tool.ArgsPreview) != "" {
		line += "(" + tool.ArgsPreview + ")"
	}
	result := status
	if !hideResult && strings.TrimSpace(tool.ResultPreview) != "" {
		resultPrefix := result + " · "
		result += " · " + alignToolResultPreview(tool.ResultPreview, lipgloss.Width("  ⎿ "+resultPrefix))
	}
	body := line + "\n  ⎿ " + result
	return renderToolBullet() + toolStyleForStatus(status).Render(wrapText(body, m.messageWidth()-3))
}

func alignToolResultPreview(preview string, continuationIndent int) string {
	lines := strings.Split(preview, "\n")
	if len(lines) <= 1 {
		return preview
	}
	indent := strings.Repeat(" ", max(0, continuationIndent))
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		lines[i] = indent + lines[i]
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderTodoWriteToolMessage(tool *ToolCallView, status string) string {
	body := toolIcon(status) + " Update Todos"
	if todos, ok := parseTodoWriteTodos(tool.RawArgs); ok {
		for i, todo := range todos {
			prefix := "    "
			if i == 0 {
				prefix = "  ⎿ "
			}
			body += "\n" + prefix + todo.checkbox(ansiForeground(toolTextColorForStatus(status))) + " " + todo.Content
		}
	} else if strings.TrimSpace(tool.ResultPreview) != "" {
		body += "\n  ⎿ " + status + " · " + tool.ResultPreview
	} else {
		body += "\n  ⎿ " + status
	}
	return renderToolBullet() + toolStyleForStatus(status).Render(wrapText(body, m.messageWidth()-3))
}

type todoDisplayItem struct {
	Content string `json:"content"`
	Status  string `json:"status"`
}

func (t todoDisplayItem) checkbox(restoreSequence string) string {
	switch t.Status {
	case "completed":
		return "[✓]"
	case "in_progress":
		return "[" + ansiForeground(tuiPalette.blue) + "•" + restoreSequence + "]"
	default:
		return "[ ]"
	}
}

func parseTodoWriteTodos(rawArgs string) ([]todoDisplayItem, bool) {
	var args struct {
		Todos []todoDisplayItem `json:"todos"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil {
		return nil, false
	}
	return args.Todos, true
}

func isTodoWriteTool(name string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "_", ""))
	return normalized == "todowrite"
}

func toolIcon(status string) string {
	switch status {
	case "success":
		return "✓"
	case "rejected", "error", "failed":
		return "✗"
	default:
		return "⏺"
	}
}

func displayToolName(name string, rawArgs string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return "Tool"
	}
	if strings.EqualFold(trimmed, "spawn_subagent") {
		if subType := spawnSubagentDisplayName(rawArgs); subType != "" {
			return subType
		}
	}
	if strings.HasPrefix(trimmed, "mcp__") {
		parts := strings.Split(trimmed, "__")
		if len(parts) >= 3 {
			return "MCP " + humanToolName(parts[len(parts)-1])
		}
	}
	return humanToolName(trimmed)
}

func spawnSubagentDisplayName(rawArgs string) string {
	var args struct {
		SubType string `json:"sub_type"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil {
		return ""
	}
	switch strings.TrimSpace(args.SubType) {
	case "explore":
		return "Explore"
	case "general":
		return "General"
	default:
		return ""
	}
}

func humanToolName(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool { return r == '_' || r == '-' })
	if len(parts) == 0 {
		return name
	}
	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + part[1:]
	}
	return strings.Join(parts, "")
}

func (m Model) messageWidth() int {
	return max(10, m.width)
}

func newMarkdownRenderer(width int) *glamour.TermRenderer {
	style := styles.DarkStyleConfig
	inlineReferenceColor := string(tuiPalette.butter)
	style.Code.Color = &inlineReferenceColor
	style.Code.BackgroundColor = nil
	style.CodeBlock.Chroma = nil
	style.CodeBlock.Theme = ""
	renderer, _ := glamour.NewTermRenderer(glamour.WithStyles(style), glamour.WithWordWrap(max(10, width)))
	return renderer
}

var fileReferencePattern = regexp.MustCompile(`(^|[\s\(\[\{\"'“‘，。；：、])((?:~?/|\./|\.\./)?[A-Za-z0-9._-]+(?:/[A-Za-z0-9._-]+)+(?:/)?|[A-Za-z0-9._-]+/|[A-Za-z0-9._-]+\.[A-Za-z][A-Za-z0-9]{0,7}|[A-Za-z_][A-Za-z0-9_]*(?:\.[A-Za-z_][A-Za-z0-9_]*)?\(\))([\s\)\]\}\"'”’。，；：、\x1b]|$)`)

func colorizeFileReferences(text string) string {
	return colorizeFileReferencesWithRestore(text, "\x1b[0m")
}

func colorizeFileReferencesForInput(text string) string {
	return colorizeFileReferencesWithRestore(text, userInputTextStart())
}

func colorizeFileReferencesForSelectedLine(text string) string {
	return colorizeFileReferencesWithRestore(text, selectedUserLineStart())
}

func colorizeFileReferencesWithRestore(text string, restoreSequence string) string {
	return fileReferencePattern.ReplaceAllStringFunc(text, func(match string) string {
		parts := fileReferencePattern.FindStringSubmatch(match)
		if len(parts) != 4 {
			return match
		}
		if strings.HasSuffix(parts[2], "/...") {
			return match
		}
		return parts[1] + renderFileReferenceWithRestore(parts[2], restoreSequence) + parts[3]
	})
}

func renderFileReference(text string) string {
	return renderFileReferenceWithRestore(text, "\x1b[0m")
}

func renderFileReferenceWithRestore(text string, restoreSequence string) string {
	if text == "" {
		return text
	}
	return "\x1b[38;5;" + string(tuiPalette.butter) + "m" + text + restoreSequence
}

func selectedUserLineStart() string {
	return "\x1b[48;5;" + string(tuiPalette.panel) + ";38;5;" + string(tuiPalette.ink) + "m"
}

func userInputTextStart() string {
	return "\x1b[0;38;5;" + string(tuiPalette.ink) + "m"
}

func ansiForeground(color lipgloss.Color) string {
	return "\x1b[0;38;5;" + string(color) + "m"
}

func ansiThinking() string {
	return "\x1b[0;3;38;5;242m"
}

func renderSelectedUserMessage(content string, width int) string {
	width = max(1, width)
	wrapped := wrapText("› "+content, width)
	lines := strings.Split(wrapped, "\n")
	for i, line := range lines {
		styled := colorizeFileReferencesForSelectedLine(line)
		padding := strings.Repeat(" ", max(0, width-lipgloss.Width(line)))
		lines[i] = selectedUserLineStart() + styled + padding + "\x1b[0m"
	}
	return strings.Join(lines, "\n")
}

func wrapText(text string, width int) string {
	width = max(1, width)
	return ansi.Hardwrap(text, width, true)
}

func truncateToWidth(text string, width int) string {
	if width <= 0 || lipgloss.Width(text) <= width {
		return text
	}
	if width == 1 {
		return "…"
	}
	var b strings.Builder
	for _, r := range text {
		candidate := b.String() + string(r) + "…"
		if lipgloss.Width(candidate) > width {
			break
		}
		b.WriteRune(r)
	}
	return b.String() + "…"
}

func (m Model) renderWelcome() string {
	return ""
}

func (m Model) renderHeader() string {
	outerWidth := max(30, m.width)
	panelWidth := max(20, outerWidth)
	contentWidth := max(10, panelWidth-4)
	if m.height > 0 && m.height < 6 {
		return m.renderHeaderLine(max(10, outerWidth))
	}
	if m.height > 0 && m.height < 10 {
		return compactHeaderStyle().Width(contentWidth).Render(m.renderHeaderLine(contentWidth))
	}
	if m.height > 0 && m.height < 18 {
		return compactHeaderStyle().Width(contentWidth).Render(m.renderCompactHeader(contentWidth))
	}
	return headerStyle().Width(contentWidth).Render(m.renderFullHeader(contentWidth))
}

func (m Model) renderHeaderLine(width int) string {
	return titleStyle().Render(centerHeaderLine("cynosure version: 0.0.1", width))
}

func (m Model) renderCompactHeader(width int) string {
	return titleStyle().Render(centerHeaderLine("cynosure version: 0.0.1", width)) + "\n" + titleStyle().Render(centerHeaderLine("Welcome back", width)) + "\n" + subtleStyle().Render(centerHeaderLine("model "+m.modelLabel()+" · "+renderFileReference(m.workspaceLabel()), width))
}

func (m Model) renderFullHeader(width int) string {
	lines := []string{
		titleStyle().Render(centerHeaderLine(`cynosure version: 0.0.1`, width)),
		mascotStyle().Render(centerHeaderLine(`/^ ^\`, width)),
		mascotStyle().Render(centerHeaderLine(`/ 0 0 \`, width)),
		mascotStyle().Render(centerHeaderLine(`V\ Y /V`, width)),
		mascotStyle().Render(centerHeaderLine(`  / - \`, width)),
		mascotStyle().Render(centerHeaderLine(` || (__V`, width)),
		titleStyle().Render(centerHeaderLine("Welcome back", width)),
		titleStyle().Render(centerHeaderLine("model "+m.modelLabel(), width)),
		subtleStyle().Render(centerHeaderLine(renderFileReference(m.workspaceLabel()), width)),
	}
	return strings.Join(lines, "\n")
}

func centerHeaderLine(text string, width int) string {
	return lipgloss.PlaceHorizontal(width, lipgloss.Center, truncateToWidth(text, width))
}

func (m Model) modelLabel() string {
	model := strings.TrimSpace(m.session.ModelID)
	if model == "" {
		model = "default model"
	}
	return model
}

func (m Model) workspaceLabel() string {
	cwd := strings.TrimSpace(m.session.CWD)
	if cwd == "" {
		return "."
	}
	return cwd
}

func (m Model) renderConversationFrame() string {
	return trimSelectablePadding(conversationStyle().Width(max(10, m.width)).Render(m.viewport.View()))
}

func trimSelectablePadding(view string) string {
	lines := strings.Split(view, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimRight(line, " ")
	}
	return strings.Join(lines, "\n")
}

func (m Model) renderInput() string {
	prompt := inputPromptStyle().Render("›")
	text := m.input.Value()
	if text == "" {
		text = inputPromptStyle().Render(inputCursor) + " " + subtleStyle().Render(m.input.Placeholder)
	} else {
		lineInfo := m.input.LineInfo()
		beforeCursor, afterCursor := splitInputAtCursor(text, lineInfo.StartColumn+lineInfo.ColumnOffset)
		text = userInputTextStart() +
			colorizeFileReferencesForInput(beforeCursor) +
			"\x1b[0m" +
			inputPromptStyle().Render(inputCursor) +
			userInputTextStart() +
			colorizeFileReferencesForInput(afterCursor) +
			"\x1b[0m"
	}
	return inputLineStyle().Width(max(10, m.width-4)).Render(prompt + " " + text)
}

func splitInputAtCursor(text string, cursorOffset int) (string, string) {
	runes := []rune(text)
	cursorOffset = min(max(0, cursorOffset), len(runes))
	return string(runes[:cursorOffset]), string(runes[cursorOffset:])
}

func (m Model) renderInputArea() string {
	if m.height > 0 && m.height < 7 {
		return m.renderInput()
	}
	return m.renderInput() + "\n" + inputStatusStyle().Width(max(10, m.width-2)).Render(m.renderLiveStatus())
}

func isTerminalProbeResponseInput(msg tea.KeyMsg) bool {
	if msg.Type != tea.KeyRunes {
		return false
	}
	value := string(msg.Runes)
	return strings.Contains(value, ";rgb:") && strings.Contains(value, "/")
}

func (m Model) renderLiveStatus() string {
	if m.width > 0 && m.width < 70 {
		parts := []string{"Enter", "Ctrl+C", "/help", fmt.Sprintf("工具 %d", m.toolCallCount)}
		if m.contextBudget > 0 {
			parts = append(parts, fmt.Sprintf("上下文 %d%%", min(100, m.contextTokens*100/m.contextBudget)))
		} else {
			parts = append(parts, "上下文 --")
		}
		return strings.Join(parts, " · ")
	}
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

func titleRuleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func headerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(headerAccentColor()).Padding(1, 1)
}

func compactHeaderStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(headerAccentColor()).Padding(0, 1)
}

func titleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(headerAccentColor()).Bold(true)
}

func headerMetaStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
}

func subtleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func conversationStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func startupPanelStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.coral).Padding(1, 2).Margin(1, 0)
}

func welcomeStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.ink)
}

func mascotStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(headerAccentColor()).Bold(true)
}

func sectionTitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().Bold(true).Foreground(tuiPalette.coral)
}

func sectionDividerStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func userStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted)
}

func thinkingStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("242")).Italic(true).PaddingLeft(2)
}

func thinkingIndicatorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral).Bold(true).PaddingLeft(1)
}

func systemStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.butter)
}

func errorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.coral)
}

func toolStyleForStatus(status string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(toolTextColorForStatus(status)).PaddingLeft(1)
}

func toolTextColorForStatus(status string) lipgloss.Color {
	switch status {
	case "success":
		return tuiPalette.mint
	case "rejected", "error", "failed":
		return tuiPalette.coral
	default:
		return lipgloss.Color("245")
	}
}

func renderToolBullet() string {
	return ansiForeground(tuiPalette.blue) + "●" + "\x1b[0m"
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

func headerAccentColor() lipgloss.Color {
	return tuiPalette.mint
}

func inputLineStyle() lipgloss.Style {
	return lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(tuiPalette.panel).Background(tuiPalette.panelDim).Foreground(tuiPalette.ink).Padding(0, 1)
}

func inputStatusStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(tuiPalette.muted).Padding(0, 1)
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

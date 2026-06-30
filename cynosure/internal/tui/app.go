package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"cynosure/internal/agent/mcp"
	"cynosure/internal/agent/runtime"
	"cynosure/internal/agent/storage"
	"cynosure/internal/logger"
	"cynosure/internal/sessions"
	"cynosure/internal/tools"
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
	ID               string
	Name             string
	RawArgs          string
	ArgsPreview      string
	Status           string
	ResultPreview    string
	Scope            string
	EphemeralGroupID string
	ParentToolCallID string
	SuppressResult   bool
	// EditLineStarts 缓存每个被编辑文件、每处改动在真实文件中的 1-based 起始行，
	// 仅用于编辑类工具（edit_file/multi_edit）的 TUI diff 行号展示。
	// 维度：[文件索引][该文件内改动索引] = 起始行；0 或缺失表示未知，渲染回退为 1。
	// 形状与 parseMultiEditDisplayFiles 返回的 []editDisplayFile 一一对应。
	EditLineStarts [][]int
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
	workingStarted    bool
	fileToolsExpanded bool
	renderCache       *messageRenderCache
}

// messageRenderCache 缓存每条消息的渲染结果，避免每个流式增量都重算全部历史。
// 渲染输出是 (消息内容, 宽度) 的纯函数，因此可安全按内容缓存；指针语义让
// Bubble Tea 的 Model 值拷贝之间共享同一份缓存。
type messageRenderCache struct {
	entries map[string]string
}

func newMessageRenderCache() *messageRenderCache {
	return &messageRenderCache{entries: make(map[string]string)}
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
	return Model{runtime: runtimeService, session: session, input: input, viewport: vp, width: 100, height: 20, events: make(chan Event, 128), renderer: renderer, autoFollow: true, renderCache: newMessageRenderCache()}
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
			if m.running && !m.approving && m.cancel != nil {
				m.cancel()
				m.generation++
				m.running = false
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
		case "ctrl+o":
			if m.toggleFileToolExpansion() {
				m.refreshViewport()
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
			if isSlashCommandInput(text) && m.handleSlashCommand(text) {
				m.refreshViewport()
				return m, nil
			}
			m.appendMessage("user", text)
			m.running = true
			m.thinkingStartedAt = time.Now()
			m.thinkingNow = m.thinkingStartedAt
			m.workingStarted = false
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
		m.applyEvent(msg)
		// 合并事件：非阻塞地把 channel 中已堆积的同批事件一次性消费完，
		// 一轮模型输出的大量连续 assistant_delta 会被压成一帧，渲染次数
		// 从「每个 token 一次」降到「每批一次」，从而让 Update 循环及时
		// 腾出来响应鼠标滚轮等输入，避免滚动延迟与卡顿。
		m.drainPendingEvents()
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

// drainPendingEvents 非阻塞地消费并应用 channel 中已就绪的事件。
func (m *Model) drainPendingEvents() {
	for {
		select {
		case next := <-m.events:
			m.applyEvent(next)
		default:
			return
		}
	}
}

// applyEvent 应用单个事件到 Model 状态，不触发重渲染，便于批量合并。
// 非当前 generation 的事件直接忽略（来自已被中断的旧轮次）。
func (m *Model) applyEvent(msg Event) {
	if msg.Generation != 0 && msg.Generation != m.generation {
		return
	}
	switch msg.Name {
	case "approval_request":
		if req, ok := msg.Data.(approvalRequestMsg); ok {
			m.beginApproval(req)
		}
	case "assistant_delta":
		m.workingStarted = true
		m.appendAssistantDelta(msg.Content)
	case "assistant_stream_reset":
		m.resetLiveAssistant()
	case "assistant":
		m.workingStarted = true
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
		m.workingStarted = true
		m.appendToolCallStart(msg.Data)
	case "tool_call_done":
		m.updateToolCallDone(msg.Data)
	case "tool_call_group_clear":
		m.clearToolCallGroup(msg.Data)
	case "error":
		m.appendMessage("error", msg.Content)
		m.running = false
	case "done":
		m.running = false
	}
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
	if !isSlashCommandInput(text) {
		return false
	}
	switch strings.TrimSpace(text) {
	case "/help":
		m.appendMessage("system", renderSlashCommandHelp())
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

func isSlashCommandInput(text string) bool {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "/") {
		return false
	}
	if strings.ContainsAny(text, " \t\r\n") {
		return false
	}
	for _, r := range text {
		if r > 127 {
			return false
		}
	}
	return true
}

func renderSlashCommandHelp() string {
	return strings.Join([]string{
		"/clear：开启全新对话并清空当前上下文。",
		"  /cwd：显示当前工作区。",
		"  /skills：显示已加载的 Skill 列表。",
		"  /mcp：显示 MCP server 状态和工具数量。",
		"  /resume：查看并恢复当前工作区的历史会话。",
		"  /cancel：在恢复历史会话选择中取消当前选择。",
	}, "\n")
}

func (m *Model) startNewSession() {
	if m.running {
		m.appendMessage("system", "当前正在生成，请先 Esc 中断后再执行 /clear")
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
		m.appendMessage("system", "当前正在生成，请先 Esc 中断后再执行 /resume")
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
	m.messages = messagesForDisplay(history, m.session.CWD)
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

// resumeToolResultPreviewMaxLines 与 runtime 在线流式时 previewToolResult 的
// 行数上限保持一致，确保 /resume 还原出的工具结果预览与实时渲染一致。
const resumeToolResultPreviewMaxLines = 3

// messagesForDisplay 把持久化的展示历史还原为 TUI 消息列表，原样重建实时会话：
// user/system/error 文本、assistant 文本，以及每次工具调用对应的工具行（含状态
// 图标与结果预览）。工具结果通过 tool_call_id 与 assistant 的 ToolCalls 关联，
// 因此即便 tool 结果消息排在 assistant 之后，也能在工具行上正确回填状态与预览。
func messagesForDisplay(history []storage.Message, workspaceRoot string) []Message {
	if len(history) == 0 {
		return nil
	}
	toolResults := make(map[string]string, len(history))
	// 持久化的 edit_file/multi_edit diff 真实行号：在工具执行时（文件内容最新）
	// 计算并随展示历史落库，因此 /resume 直接复用，无需依赖当前磁盘文件状态，
	// 即便文件后续被改动或进程重启也能还原准确行号。
	toolEditLineStarts := make(map[string][][]int, len(history))
	for _, msg := range history {
		if msg.Role == "tool" && strings.TrimSpace(msg.ToolCallID) != "" {
			toolResults[msg.ToolCallID] = msg.Content
			if len(msg.EditLineStarts) > 0 {
				toolEditLineStarts[msg.ToolCallID] = msg.EditLineStarts
			}
		}
	}
	messages := make([]Message, 0, len(history))
	for _, msg := range history {
		switch msg.Role {
		case "user", "system", "error":
			messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
		case "assistant":
			// 实时流式时空文本的 assistant 轮次不会产生文本气泡，只展示工具行，
			// 还原时保持一致，避免出现空白 assistant 行。
			if strings.TrimSpace(msg.Content) != "" || strings.TrimSpace(msg.ReasoningContent) != "" {
				messages = append(messages, Message{Role: msg.Role, Content: msg.Content, ReasoningContent: msg.ReasoningContent})
			}
			for _, call := range msg.ToolCalls {
				messages = append(messages, reconstructToolCallMessage(call, toolResults[call.ID], toolEditLineStarts[call.ID], workspaceRoot))
			}
		}
	}
	return messages
}

// reconstructToolCallMessage 由 assistant 的工具调用及其结果消息重建工具行视图，
// 字段映射与实时 emitToolCallStart/Done 一致：RawArgs 携带原始参数（渲染时解析
// 为参数预览），Status/ResultPreview 来源于持久化的工具结果消息。编辑类工具优先
// 使用持久化的真实行号 persistedLineStarts；缺失时（老会话）按当前 workspace 文件
// 尽力回算，作为兼容兜底。
func reconstructToolCallMessage(call storage.MessageToolCall, resultContent string, persistedLineStarts [][]int, workspaceRoot string) Message {
	view := ToolCallView{
		ID:      call.ID,
		Name:    call.Function.Name,
		RawArgs: call.Function.Arguments,
		Status:  "running",
	}
	if strings.TrimSpace(resultContent) != "" {
		status, preview := parseToolResultContent(resultContent)
		if status != "" {
			view.Status = status
		}
		view.ResultPreview = preview
	}
	if len(persistedLineStarts) > 0 {
		view.EditLineStarts = persistedLineStarts
	} else {
		view.EditLineStarts = editLineStartsForTool(workspaceRoot, view.Name, view.RawArgs, view.Status)
	}
	return Message{Role: "tool", ToolCall: &view}
}

// parseToolResultContent 解析持久化的工具结果消息内容（形如
// {"status":...,"result":...}），返回状态与裁剪后的结果预览；无法解析为 JSON 时
// 退化为对原始内容做行截断。
func parseToolResultContent(content string) (status, preview string) {
	var outcome struct {
		Status string `json:"status"`
		Result string `json:"result"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(content)), &outcome); err != nil {
		return "", truncateResultPreviewLines(content, resumeToolResultPreviewMaxLines)
	}
	return outcome.Status, truncateResultPreviewLines(outcome.Result, resumeToolResultPreviewMaxLines)
}

// truncateResultPreviewLines 复刻 runtime.truncatePreviewLines 的行截断规则，保证
// /resume 还原的工具结果预览与实时一致。
func truncateResultPreviewLines(text string, maxLines int) string {
	text = strings.TrimSpace(text)
	if text == "" || maxLines <= 0 {
		return ""
	}
	lines := strings.Split(text, "\n")
	if len(lines) <= maxLines {
		return strings.Join(lines, "\n")
	}
	omitted := len(lines) - maxLines
	kept := append([]string(nil), lines[:maxLines]...)
	kept = append(kept, fmt.Sprintf("... + %d lines", omitted))
	return strings.Join(kept, "\n")
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
		if label := skillSourceLabel(skill.Source); label != "" {
			b.WriteString(" [")
			b.WriteString(label)
			b.WriteString("]")
		}
		if strings.TrimSpace(skill.Description) != "" {
			b.WriteString(" ")
			b.WriteString(skill.Description)
		}
	}
	return b.String()
}

// skillSourceLabel 把 skill 来源标识映射为面向用户的中文来源类型，
// 不暴露磁盘路径。未知来源原样返回。
func skillSourceLabel(source string) string {
	switch strings.TrimSpace(source) {
	case "workspace":
		return "当前项目"
	case "user":
		return "家目录"
	case "builtin":
		return "系统内置"
	default:
		return strings.TrimSpace(source)
	}
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
	m.messages[len(m.messages)-1].Content += delta
}

// resetLiveAssistant 丢弃本轮正在流式显示的半截 assistant 内容。运行时在输出
// 被截断升级重试、或上下文溢出压缩重试时会发出 assistant_stream_reset：那段
// 已流式显示的内容会被废弃并重新生成，这里把它从界面移除，避免与重试结果叠加。
func (m *Model) resetLiveAssistant() {
	if len(m.messages) > 0 && isLiveAssistantRole(m.messages[len(m.messages)-1].Role) {
		m.messages = m.messages[:len(m.messages)-1]
	}
	// 重新进入“等待输出”状态：恢复 Thinking 提示，直到重试产生新的增量。
	m.workingStarted = false
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
	m.replaceSubagentEphemeralTool(tool)
	m.insertToolCallMessage(tool)
}

func (m *Model) replaceSubagentEphemeralTool(tool ToolCallView) {
	if tool.Scope != "subagent" || strings.TrimSpace(tool.EphemeralGroupID) == "" {
		return
	}
	filtered := m.messages[:0]
	for _, msg := range m.messages {
		if msg.Role == "tool" && msg.ToolCall != nil && msg.ToolCall.Scope == "subagent" && msg.ToolCall.EphemeralGroupID == tool.EphemeralGroupID {
			continue
		}
		filtered = append(filtered, msg)
	}
	m.messages = filtered
}

func (m *Model) insertToolCallMessage(tool ToolCallView) {
	msg := Message{Role: "tool", ToolCall: &tool}
	if tool.Scope != "subagent" || strings.TrimSpace(tool.ParentToolCallID) == "" {
		m.messages = append(m.messages, msg)
		return
	}
	parentIndex := -1
	for i := range m.messages {
		if m.messages[i].Role == "tool" && m.messages[i].ToolCall != nil && m.messages[i].ToolCall.ID == tool.ParentToolCallID {
			parentIndex = i
			break
		}
	}
	if parentIndex < 0 {
		m.messages = append(m.messages, msg)
		return
	}
	insertAt := parentIndex + 1
	m.messages = append(m.messages, Message{})
	copy(m.messages[insertAt+1:], m.messages[insertAt:])
	m.messages[insertAt] = msg
}

// editLineStartsForTool 为编辑类工具（edit_file/multi_edit）在 success 时计算每处
// 改动在真实文件中的起始行，供 diff 行号展示。它委托给 tools.EditFileLineStarts，
// 与运行时 exec 时的计算共用同一权威实现。仅作为缺少持久化/事件行号时的兜底
// （如老会话、或文件被外部改动后实时展示已无更准的来源）。
func editLineStartsForTool(workspaceRoot, name, rawArgs, status string) [][]int {
	if strings.TrimSpace(status) != "success" {
		return nil
	}
	return tools.EditFileLineStarts(workspaceRoot, name, rawArgs)
}

func (m *Model) updateToolCallDone(data any) {
	tool := toolCallViewFromEvent(data)
	if tool.ID == "" && tool.Name == "" {
		return
	}
	// 优先采用运行时在 exec 时下发的真实行号（文件内容最新、最准确）；事件未携带时
	// 才回退到按当前 workspace 文件回算。
	eventLineStarts := eventEditLineStarts(data)
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
			m.messages[i].ToolCall.Scope = firstNonEmpty(tool.Scope, m.messages[i].ToolCall.Scope)
			m.messages[i].ToolCall.EphemeralGroupID = firstNonEmpty(tool.EphemeralGroupID, m.messages[i].ToolCall.EphemeralGroupID)
			m.messages[i].ToolCall.ParentToolCallID = firstNonEmpty(tool.ParentToolCallID, m.messages[i].ToolCall.ParentToolCallID)
			m.messages[i].ToolCall.SuppressResult = tool.SuppressResult || m.messages[i].ToolCall.SuppressResult
			m.messages[i].ToolCall.EditLineStarts = m.resolveEditLineStarts(eventLineStarts, m.messages[i].ToolCall.Name, m.messages[i].ToolCall.RawArgs, m.messages[i].ToolCall.Status)
			return
		}
	}
	if tool.Scope == "subagent" && strings.TrimSpace(tool.EphemeralGroupID) != "" {
		return
	}
	tool.EditLineStarts = m.resolveEditLineStarts(eventLineStarts, tool.Name, tool.RawArgs, tool.Status)
	m.messages = append(m.messages, Message{Role: "tool", ToolCall: &tool})
}

// resolveEditLineStarts 优先返回事件携带的真实行号；缺失时按当前 workspace 文件回算。
func (m *Model) resolveEditLineStarts(eventLineStarts [][]int, name, rawArgs, status string) [][]int {
	if len(eventLineStarts) > 0 {
		return eventLineStarts
	}
	return editLineStartsForTool(m.session.CWD, name, rawArgs, status)
}

func (m *Model) toggleFileToolExpansion() bool {
	for i := range m.messages {
		if m.isExpandableToolMessage(i) {
			m.fileToolsExpanded = !m.fileToolsExpanded
			return true
		}
	}
	return false
}

func (m Model) isExpandableToolMessage(index int) bool {
	if index < 0 || index >= len(m.messages) {
		return false
	}
	msg := m.messages[index]
	if msg.Role != "tool" || msg.ToolCall == nil {
		return false
	}
	tool := msg.ToolCall
	if strings.TrimSpace(tool.Status) != "success" {
		return false
	}
	if isWriteFileTool(tool.Name) {
		args, ok := parseWriteFileArgs(tool.RawArgs)
		if !ok {
			return false
		}
		return len(splitDisplayLines(args.Content)) > fileToolPreviewMaxLines
	}
	if isEditFileTool(tool.Name) {
		args, ok := parseEditFileArgs(tool.RawArgs)
		if !ok {
			return false
		}
		// 仅用于行数判断，行号不影响行数，故 startLine 传 0（回退 1）。
		return len(buildEditDiffLines(args.OldText, args.NewText, 0)) > fileToolPreviewMaxLines
	}
	if isMultiEditTool(tool.Name) {
		files, ok := parseMultiEditDisplayFiles(tool.RawArgs)
		if !ok {
			return false
		}
		lineCount := 0
		for _, file := range files {
			for _, change := range file.Changes {
				lineCount += len(buildEditDiffLines(change.Old, change.New, 0))
				if lineCount > fileToolPreviewMaxLines {
					return true
				}
			}
		}
		return false
	}
	return false
}

func (m Model) isFileToolExpanded() bool {
	return m.fileToolsExpanded
}

func (m *Model) clearToolCallGroup(data any) {
	groupID := eventString(data, "ephemeral_group_id")
	if strings.TrimSpace(groupID) == "" {
		return
	}
	filtered := m.messages[:0]
	for _, msg := range m.messages {
		if msg.Role == "tool" && msg.ToolCall != nil && msg.ToolCall.EphemeralGroupID == groupID {
			continue
		}
		filtered = append(filtered, msg)
	}
	m.messages = filtered
}

func toolCallViewFromEvent(data any) ToolCallView {
	return ToolCallView{
		ID:               eventString(data, "tool_call_id"),
		Name:             eventString(data, "tool_name"),
		RawArgs:          eventString(data, "raw_args"),
		ArgsPreview:      eventString(data, "args_preview"),
		Status:           eventString(data, "status"),
		ResultPreview:    eventString(data, "result_preview"),
		Scope:            eventString(data, "scope"),
		EphemeralGroupID: eventString(data, "ephemeral_group_id"),
		ParentToolCallID: eventString(data, "parent_tool_call_id"),
		SuppressResult:   eventBool(data, "suppress_result"),
	}
}

func firstNonEmpty(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return primary
	}
	return fallback
}

func isLiveAssistantRole(role string) bool {
	return role == "assistant"
}

func (m *Model) refreshViewport() {
	shouldFollow := m.autoFollow || m.viewport.AtBottom()
	m.viewport.Width = max(10, m.width)
	m.viewport.Height = m.viewportHeight()
	m.pruneRenderCache()
	m.viewport.SetContent(m.renderTranscript())
	if shouldFollow {
		m.viewport.GotoBottom()
	}
}

// pruneRenderCache 仅保留当前消息列表对应的缓存项，避免流式增量（每次内容
// 增长都会产生新 key）和窗口宽度变化导致缓存无限膨胀。
func (m *Model) pruneRenderCache() {
	if m.renderCache == nil {
		return
	}
	live := make(map[string]struct{}, len(m.messages))
	width := m.messageWidth()
	for _, msg := range m.messages {
		live[m.messageRenderKey(msg, width)] = struct{}{}
	}
	for key := range m.renderCache.entries {
		if _, ok := live[key]; !ok {
			delete(m.renderCache.entries, key)
		}
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
	for _, msg := range m.messages {
		if b.Len() > 0 && isSubagentToolMessage(msg) {
			trimTrailingBlankLines(&b)
			b.WriteString("\n")
		}
		b.WriteString(m.renderCachedMessage(msg))
		b.WriteString("\n\n")
	}
	if indicator := m.renderThinkingIndicator(); indicator != "" {
		b.WriteString(indicator)
		b.WriteString("\n\n")
	}
	return b.String()
}

// renderCachedMessage 复用缓存渲染结果，仅在 (内容, 宽度) 变化时重算。
// 流式输出时只有最后一条 assistant 消息在变化，因此单帧渲染成本从
// O(历史长度) 降到 O(1)。
func (m Model) renderCachedMessage(msg Message) string {
	if m.renderCache == nil {
		return m.renderMessageAt(msg)
	}
	key := m.messageRenderKey(msg, m.messageWidth())
	if cached, ok := m.renderCache.entries[key]; ok {
		return cached
	}
	rendered := m.renderMessageAt(msg)
	m.renderCache.entries[key] = rendered
	return rendered
}

// messageRenderKey 捕获所有会影响一条消息渲染输出的字段。
func (m Model) messageRenderKey(msg Message, width int) string {
	var b strings.Builder
	b.WriteString(strconv.Itoa(width))
	b.WriteByte('\x00')
	b.WriteString(msg.Role)
	b.WriteByte('\x00')
	b.WriteString(msg.Content)
	b.WriteByte('\x00')
	b.WriteString(msg.ReasoningContent)
	if tool := msg.ToolCall; tool != nil {
		b.WriteByte('\x00')
		fields := []string{
			tool.ID, tool.Name, tool.RawArgs, tool.ArgsPreview, tool.Status,
			tool.ResultPreview, tool.Scope, tool.EphemeralGroupID, tool.ParentToolCallID,
		}
		for _, field := range fields {
			b.WriteString(field)
			b.WriteByte('\x00')
		}
		if tool.SuppressResult {
			b.WriteByte('1')
		}
		b.WriteByte('\x00')
		if m.fileToolsExpanded {
			b.WriteByte('E')
		}
	}
	return b.String()
}

func trimTrailingBlankLines(b *strings.Builder) {
	text := strings.TrimRight(b.String(), "\n")
	b.Reset()
	b.WriteString(text)
}

func isSubagentToolMessage(msg Message) bool {
	return msg.Role == "tool" && msg.ToolCall != nil && msg.ToolCall.Scope == "subagent"
}

func (m Model) renderThinkingIndicator() string {
	if !m.running || m.thinkingStartedAt.IsZero() {
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
	label := "Thinking"
	if m.workingStarted {
		label = "Working"
	}
	return thinkingIndicatorStyle().Render(fmt.Sprintf("* %s... (%ds)", label, elapsedSeconds))
}

func (m Model) renderMessage(msg Message) string {
	return m.renderMessageAt(msg)
}

func (m Model) renderMessageAt(msg Message) string {
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
		return content
	case "system":
		return systemStyle().Render("• " + wrapText(colorizeFileReferencesWithRestore(msg.Content, ansiForeground(tuiPalette.butter)), m.messageWidth()-2))
	case "error":
		return errorStyle().Render("✗ " + wrapText(colorizeFileReferencesWithRestore(msg.Content, ansiForeground(tuiPalette.coral)), m.messageWidth()-2))
	case "tool":
		return m.renderToolMessage(msg)
	default:
		return roleLabel(msg.Role, lipgloss.Color("245")) + "\n" + wrapText(colorizeFileReferences(msg.Content), m.messageWidth())
	}
}

func (m Model) renderToolMessage(msg Message) string {
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
	if isWriteFileTool(tool.Name) {
		if rendered, ok := m.renderWriteFileToolMessage(tool, status); ok {
			return rendered
		}
	}
	if isEditFileTool(tool.Name) {
		if rendered, ok := m.renderEditFileToolMessage(tool, status); ok {
			return rendered
		}
	}
	if isMultiEditTool(tool.Name) {
		if rendered, ok := m.renderMultiEditToolMessage(tool, status); ok {
			return rendered
		}
	}
	name := displayToolName(tool.Name, tool.RawArgs)
	line := name
	if args := toolArgsDisplay(tool.RawArgs, tool.ArgsPreview); args != "" {
		line += "(" + args + ")"
	}
	if tool.Scope == "subagent" {
		return ansiForeground(tuiPalette.muted) + renderSubagentToolStatus(line, m.messageWidth()) + "\x1b[0m"
	}
	icon := toolIcon(status)
	line = limitToWrappedLines(icon+" "+line, m.messageWidth()-3, 2)
	result := status
	if !tool.SuppressResult && strings.TrimSpace(tool.ResultPreview) != "" {
		resultPrefix := result + " · "
		result += " · " + alignToolResultPreview(tool.ResultPreview, lipgloss.Width("  ⎿ "+resultPrefix))
	}
	body := line + "\n  ⎿ " + result
	return renderToolBullet() + toolStyleForStatus(status).Render(wrapText(body, m.messageWidth()-3))
}

func renderSubagentToolStatus(body string, width int) string {
	width = max(1, width)
	lines := strings.Split(body, "\n")
	if len(lines) == 0 {
		return ""
	}
	prefix := strings.Repeat(" ", lipgloss.Width("   ⎿ "))
	wrapped := limitToWrappedLines(lines[0], max(1, width-lipgloss.Width(prefix)), 2)
	wrappedLines := strings.Split(wrapped, "\n")
	for i := range wrappedLines {
		wrappedLines[i] = prefix + strings.TrimLeft(wrappedLines[i], " ")
	}
	if len(lines) > 1 {
		wrappedLines = append(wrappedLines, lines[1:]...)
	}
	return strings.Join(wrappedLines, "\n")
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

const fileToolPreviewMaxLines = 10

type writeFileDisplayArgs struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

func (m Model) renderWriteFileToolMessage(tool *ToolCallView, status string) (string, bool) {
	if status != "success" {
		return "", false
	}
	args, ok := parseWriteFileArgs(tool.RawArgs)
	if !ok {
		return "", false
	}
	lines := splitDisplayLines(args.Content)
	highlightedLines := highlightCodeLines(args.FilePath, lines)
	bodyLines := []string{fmt.Sprintf("%s Wrote %d %s to %s", toolIcon(status), len(lines), pluralizeLine(len(lines)), args.FilePath)}
	previewLineCount := min(len(lines), fileToolPreviewMaxLines)
	if m.isFileToolExpanded() {
		previewLineCount = len(lines)
	}
	for i := 0; i < previewLineCount; i++ {
		bodyLines = append(bodyLines, formatWritePreviewLine(i+1, highlightedLines[i]))
	}
	if omitted := len(lines) - previewLineCount; omitted > 0 {
		bodyLines = append(bodyLines, fmt.Sprintf("... +%d lines  [Ctrl+O to expand]", omitted))
	} else if m.isFileToolExpanded() && len(lines) > fileToolPreviewMaxLines {
		bodyLines = append(bodyLines, "[Ctrl+O to collapse]")
	}
	body := strings.Join(bodyLines, "\n")
	return renderToolBullet() + toolStyleForStatus(status).Render(wrapText(body, m.messageWidth()-3)), true
}

func parseWriteFileArgs(rawArgs string) (writeFileDisplayArgs, bool) {
	var args writeFileDisplayArgs
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil {
		return args, false
	}
	return args, strings.TrimSpace(args.FilePath) != ""
}

func formatWritePreviewLine(lineNumber int, content string) string {
	return fmt.Sprintf("%3d│ %s", lineNumber, content)
}

func highlightCodeLines(filePath string, lines []string) []string {
	source := strings.Join(lines, "\n")
	var b strings.Builder
	lexerName := ""
	if lexer := lexers.Match(filePath); lexer != nil {
		lexerName = lexer.Config().Name
	}
	if err := quick.Highlight(&b, source, lexerName, "terminal16m", "monokai"); err != nil {
		return lines
	}
	rendered := strings.TrimSuffix(b.String(), "\n")
	if rendered == "" && len(lines) == 1 && lines[0] == "" {
		return []string{""}
	}
	highlighted := strings.Split(rendered, "\n")
	if len(highlighted) < len(lines) {
		highlighted = append(highlighted, lines[len(highlighted):]...)
	}
	if len(highlighted) > len(lines) {
		highlighted = highlighted[:len(lines)]
	}
	return highlighted
}

type editFileDisplayArgs struct {
	FilePath string `json:"file_path"`
	OldText  string `json:"old_text"`
	NewText  string `json:"new_text"`
}

type multiEditDisplayArgs struct {
	FilePath string                 `json:"file_path"`
	Edits    []multiEditDisplayEdit `json:"edits"`
	Files    []multiEditDisplayFile `json:"files"`
}

type multiEditDisplayFile struct {
	FilePath string                 `json:"file_path"`
	Edits    []multiEditDisplayEdit `json:"edits"`
}

type multiEditDisplayEdit struct {
	OldString string `json:"old_string"`
	NewString string `json:"new_string"`
}

type editDiffLine struct {
	Kind   rune
	Number int
	Text   string
}

func (m Model) renderEditFileToolMessage(tool *ToolCallView, status string) (string, bool) {
	if status != "success" {
		return "", false
	}
	args, ok := parseEditFileArgs(tool.RawArgs)
	if !ok {
		return "", false
	}
	bodyLines := buildEditPreviewBodyLines(toolIcon(status), "", []editDisplayChange{{Old: args.OldText, New: args.NewText}}, editLineStartsForFile(tool.EditLineStarts, 0), status, m.isFileToolExpanded())
	return m.renderFileToolBody(bodyLines, status), true
}

func (m Model) renderMultiEditToolMessage(tool *ToolCallView, status string) (string, bool) {
	if status != "success" {
		return "", false
	}
	files, ok := parseMultiEditDisplayFiles(tool.RawArgs)
	if !ok {
		return "", false
	}
	blocks := make([]string, 0, len(files))
	for fileIdx, file := range files {
		bodyLines := buildEditPreviewBodyLines(toolIcon(status), file.FilePath, file.Changes, editLineStartsForFile(tool.EditLineStarts, fileIdx), status, m.isFileToolExpanded())
		blocks = append(blocks, m.renderFileToolBody(bodyLines, status))
	}
	return strings.Join(blocks, "\n\n"), true
}

// editLineStartsForFile 取出第 fileIdx 个文件的每处改动起始行；越界返回 nil，
// 由 buildEditDiffLines 回退为 1。
func editLineStartsForFile(starts [][]int, fileIdx int) []int {
	if fileIdx < 0 || fileIdx >= len(starts) {
		return nil
	}
	return starts[fileIdx]
}

func (m Model) renderFileToolBody(bodyLines []string, status string) string {
	body := strings.Join(bodyLines, "\n")
	return renderToolBullet() + toolStyleForStatus(status).Render(wrapText(body, m.messageWidth()-3))
}

func parseEditFileArgs(rawArgs string) (editFileDisplayArgs, bool) {
	var args editFileDisplayArgs
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil {
		return args, false
	}
	return args, strings.TrimSpace(args.FilePath) != ""
}

type editDisplayFile struct {
	FilePath string
	Changes  []editDisplayChange
}

type editDisplayChange struct {
	Old string
	New string
}

func parseMultiEditDisplayFiles(rawArgs string) ([]editDisplayFile, bool) {
	var args multiEditDisplayArgs
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err != nil {
		return nil, false
	}
	if len(args.Files) > 0 {
		files := make([]editDisplayFile, 0, len(args.Files))
		for _, file := range args.Files {
			parsed, ok := multiEditDisplayFileToEditDisplay(file.FilePath, file.Edits)
			if !ok {
				return nil, false
			}
			files = append(files, parsed)
		}
		return files, true
	}
	file, ok := multiEditDisplayFileToEditDisplay(args.FilePath, args.Edits)
	if !ok {
		return nil, false
	}
	return []editDisplayFile{file}, true
}

func multiEditDisplayFileToEditDisplay(filePath string, edits []multiEditDisplayEdit) (editDisplayFile, bool) {
	if strings.TrimSpace(filePath) == "" || len(edits) == 0 {
		return editDisplayFile{}, false
	}
	changes := make([]editDisplayChange, 0, len(edits))
	for _, edit := range edits {
		changes = append(changes, editDisplayChange{Old: edit.OldString, New: edit.NewString})
	}
	return editDisplayFile{FilePath: filePath, Changes: changes}, true
}

func buildEditPreviewBodyLines(icon, filePath string, changes []editDisplayChange, startLines []int, status string, expanded bool) []string {
	diffLines := make([]editDiffLine, 0)
	for i, change := range changes {
		startLine := 0
		if i < len(startLines) {
			startLine = startLines[i]
		}
		diffLines = append(diffLines, buildEditDiffLines(change.Old, change.New, startLine)...)
	}
	added, removed := countEditDiffLines(diffLines)
	summary := fmt.Sprintf("Added %d %s, removed %d %s", added, pluralizeLine(added), removed, pluralizeLine(removed))
	if strings.TrimSpace(filePath) != "" {
		summary = filePath + ": " + summary
	}
	if icon != "" {
		summary = icon + " " + summary
	}
	bodyLines := []string{summary}
	previewLines := diffLines
	if !expanded && len(diffLines) > fileToolPreviewMaxLines {
		previewLines = diffLines[:fileToolPreviewMaxLines]
	}
	for _, line := range previewLines {
		bodyLines = append(bodyLines, renderEditDiffLine(line, status))
	}
	if omitted := len(diffLines) - len(previewLines); omitted > 0 {
		bodyLines = append(bodyLines, fmt.Sprintf("... +%d lines  [Ctrl+O to expand]", omitted))
	} else if expanded && len(diffLines) > fileToolPreviewMaxLines {
		bodyLines = append(bodyLines, "[Ctrl+O to collapse]")
	}
	return bodyLines
}

// buildEditDiffLines 构建一处编辑的 diff 行，并按两套参照系编号：
// 删除行（-）使用原文件行号，新增行（+）使用新文件行号，上下文行使用新文件行号。
// startLine 是该编辑匹配区域在文件中的 1-based 起始行；<=0 时回退为 1（行号未知）。
// 因前缀在新旧文件中完全一致，删除块与新增块的起始行相同，均为 startLine，
// 之后各自按行数自增，从而在尾段因 old/new 行数不同而自然分叉。
func buildEditDiffLines(oldText, newText string, startLine int) []editDiffLine {
	oldLines := splitDisplayLines(oldText)
	newLines := splitDisplayLines(newText)
	prefix := commonPrefixLineCount(oldLines, newLines)
	suffix := commonSuffixLineCount(oldLines, newLines, prefix)

	if startLine <= 0 {
		startLine = 1
	}
	oldNo := startLine // 原文件行号游标
	newNo := startLine // 新文件行号游标

	lines := make([]editDiffLine, 0, len(oldLines)+len(newLines)+1)
	for i := 0; i < prefix; i++ {
		lines = append(lines, editDiffLine{Kind: ' ', Number: newNo, Text: oldLines[i]})
		oldNo++
		newNo++
	}
	for i := prefix; i < len(oldLines)-suffix; i++ {
		lines = append(lines, editDiffLine{Kind: '-', Number: oldNo, Text: oldLines[i]})
		oldNo++
	}
	if prefix < len(oldLines)-suffix && prefix < len(newLines)-suffix {
		lines = append(lines, editDiffLine{Kind: '…'})
	}
	for i := prefix; i < len(newLines)-suffix; i++ {
		lines = append(lines, editDiffLine{Kind: '+', Number: newNo, Text: newLines[i]})
		newNo++
	}
	for i := len(oldLines) - suffix; i < len(oldLines); i++ {
		lines = append(lines, editDiffLine{Kind: ' ', Number: newNo, Text: oldLines[i]})
		oldNo++
		newNo++
	}
	return lines
}

func commonPrefixLineCount(a, b []string) int {
	n := min(len(a), len(b))
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return i
		}
	}
	return n
}

func commonSuffixLineCount(a, b []string, prefix int) int {
	maxSuffix := min(len(a), len(b)) - prefix
	for i := 0; i < maxSuffix; i++ {
		if a[len(a)-1-i] != b[len(b)-1-i] {
			return i
		}
	}
	return maxSuffix
}

func countEditDiffLines(lines []editDiffLine) (int, int) {
	added, removed := 0, 0
	for _, line := range lines {
		switch line.Kind {
		case '+':
			added++
		case '-':
			removed++
		}
	}
	return added, removed
}

func renderEditDiffLine(line editDiffLine, status string) string {
	if line.Kind == '…' {
		return "..."
	}
	text := fmt.Sprintf("%c%2d│ %s", line.Kind, line.Number, line.Text)
	restore := ansiForeground(toolTextColorForStatus(status))
	switch line.Kind {
	case '+':
		return ansiForeground(tuiPalette.mint) + text + restore
	case '-':
		return ansiForeground(tuiPalette.coral) + text + restore
	case ' ':
		return ansiForeground(tuiPalette.muted) + text + restore
	default:
		return text
	}
}

func splitDisplayLines(text string) []string {
	lines := strings.Split(strings.TrimSuffix(text, "\n"), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{""}
	}
	return lines
}

func pluralizeLine(n int) string {
	if n == 1 {
		return "line"
	}
	return "lines"
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

func isWriteFileTool(name string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "_", ""))
	return normalized == "writefile"
}

func isEditFileTool(name string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "_", ""))
	return normalized == "editfile"
}

func isMultiEditTool(name string) bool {
	normalized := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), "_", ""))
	return normalized == "multiedit"
}

func toolArgsDisplay(rawArgs, fallback string) string {
	var args map[string]any
	if err := json.Unmarshal([]byte(strings.TrimSpace(rawArgs)), &args); err == nil && len(args) > 0 {
		keys := make([]string, 0, len(args))
		for key := range args {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		parts := make([]string, 0, len(keys))
		for _, key := range keys {
			parts = append(parts, key+": "+formatToolArgValue(args[key]))
		}
		return strings.Join(parts, ", ")
	}
	return strings.TrimSpace(fallback)
}

func formatToolArgValue(value any) string {
	switch v := value.(type) {
	case nil:
		return "null"
	case string:
		return v
	case float64, bool:
		return fmt.Sprint(v)
	default:
		encoded, err := json.Marshal(v)
		if err != nil {
			return fmt.Sprint(v)
		}
		return string(encoded)
	}
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
	switch strings.ToLower(trimmed) {
	case "write_file":
		return "write"
	case "read_file":
		return "read"
	case "glob":
		return "glob"
	case "grep":
		return "grep"
	case "edit_file":
		return "edit"
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
	value := string(color)
	if strings.HasPrefix(value, "#") {
		if r, g, b, ok := parseHexColor(value); ok {
			return fmt.Sprintf("\x1b[0;38;2;%d;%d;%dm", r, g, b)
		}
	}
	return "\x1b[0;38;5;" + value + "m"
}

func parseHexColor(value string) (int, int, int, bool) {
	hex := strings.TrimPrefix(value, "#")
	if len(hex) != 6 {
		return 0, 0, 0, false
	}
	parsed, err := strconv.ParseInt(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, false
	}
	return int(parsed>>16) & 0xFF, int(parsed>>8) & 0xFF, int(parsed) & 0xFF, true
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

// limitToWrappedLines 将文本按 width 折行后，最多保留 maxLines 行，
// 超出部分以省略号结尾，用于截断过长的工具参数展示。
func limitToWrappedLines(text string, width, maxLines int) string {
	width = max(1, width)
	wrapped := wrapText(text, width)
	if maxLines <= 0 {
		return wrapped
	}
	lines := strings.Split(wrapped, "\n")
	if len(lines) <= maxLines {
		return wrapped
	}
	kept := lines[:maxLines]
	last := strings.TrimRight(kept[maxLines-1], " ")
	kept[maxLines-1] = truncateToWidth(last+"…", width)
	return strings.Join(kept, "\n")
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
		beforeCursor, afterCursor := splitInputAtCursor(text, m.inputCursorOffset())
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

// inputCursorOffset 返回光标在整段输入文本中的绝对 rune 偏移。
// textarea 的 LineInfo().StartColumn+ColumnOffset 是「当前逻辑行内」的列号，
// 并不包含光标所在行之前各行的长度。对单行输入二者相等，但多行（粘贴）输入
// 时必须把光标行之前各逻辑行的长度（每行含一个换行符）累加进来，否则可见
// 光标会被渲染到全文很靠前的位置。
func (m Model) inputCursorOffset() int {
	li := m.input.LineInfo()
	colWithinRow := li.StartColumn + li.ColumnOffset
	row := m.input.Line()
	rows := strings.Split(m.input.Value(), "\n")
	offset := 0
	for i := 0; i < row && i < len(rows); i++ {
		offset += len([]rune(rows[i])) + 1 // +1 计入该行末尾的换行符
	}
	return offset + colWithinRow
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
		parts := []string{"Enter", "Esc 中断", "/help"}
		if m.contextBudget > 0 {
			parts = append(parts, fmt.Sprintf("上下文 %d%%", min(100, m.contextTokens*100/m.contextBudget)))
		} else {
			parts = append(parts, "上下文 --")
		}
		return strings.Join(parts, " · ")
	}
	parts := []string{"Enter 发送", "Esc 中断", "Ctrl+C 退出", "/help"}
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

func eventBool(data any, key string) bool {
	m, ok := data.(map[string]any)
	if !ok {
		return false
	}
	value, _ := m[key].(bool)
	return value
}

// eventEditLineStarts 从 tool_call_done 事件读取运行时在 exec 时计算好的 diff 真实
// 行号。事件经内存 channel 传递（非 JSON），因此值就是 [][]int 本身。
func eventEditLineStarts(data any) [][]int {
	m, ok := data.(map[string]any)
	if !ok {
		return nil
	}
	starts, _ := m["edit_line_starts"].([][]int)
	return starts
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
		return lipgloss.Color("#8FB89A")
	case "rejected", "error", "failed":
		return lipgloss.Color("#C79B92")
	default:
		return lipgloss.Color("#9A9A9A")
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

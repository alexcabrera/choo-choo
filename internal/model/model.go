package model

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"

	"github.com/alexcabrera/choo-choo/internal/crush"
	"github.com/alexcabrera/choo-choo/internal/orchestrator"
	"github.com/alexcabrera/choo-choo/internal/state"
	"github.com/alexcabrera/choo-choo/internal/ticket"
	"github.com/alexcabrera/choo-choo/internal/ui/chat"
	"github.com/alexcabrera/choo-choo/internal/ui/kanban"
	"github.com/alexcabrera/choo-choo/internal/ui/popup"
)

type FocusArea int

const (
	FocusChat FocusArea = iota
	FocusKanban
	FocusPreview
	FocusPopup
)

type initCompleteMsg struct {
	err error
}

type spinnerTickMsg struct{}

type streamChunkMsg struct {
	content string
	done    bool
}

type phaseCompleteMsg struct {
	phase state.Phase
	err   error
}

type ticketEventMsg struct {
	ticketID string
	status   string
}

type ticketEventsStartMsg struct {
	events <-chan orchestrator.TicketEvent
	err    error
}

type crushStreamMsg struct {
	event crush.StreamEvent
}

type crushStreamStartMsg struct {
	events <-chan crush.StreamEvent
	event  crush.StreamEvent
	err    error
}

type orchestratorTicketEventMsg struct {
	event orchestrator.TicketEvent
}

var (
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#6B7280")
	successColor   = lipgloss.Color("#10B981")
	errorColor     = lipgloss.Color("#EF4444")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1)

	phaseStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(secondaryColor).
			Padding(0, 1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(secondaryColor).
			Padding(1, 2)

	focusedBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)
)

type Model struct {
	orchestrator *orchestrator.Orchestrator
	chat         *chat.ChatModel
	kanban       *kanban.KanbanModel
	popup        *popup.PopupModel
	phase        state.Phase
	tickets      []ticket.Ticket
	width        int
	height       int
	focus        FocusArea
	loading      bool
	runningPhase bool
	errors       []string
	helpVisible  bool
	artifactPath string
	artifactText string
	ctx          context.Context
	cancelCtx    context.CancelFunc
	crushEvents  <-chan crush.StreamEvent
	ticketEvents <-chan orchestrator.TicketEvent
}

func New() Model {
	ctx, cancel := context.WithCancel(context.Background())
	return Model{
		phase:     state.PhaseInit,
		focus:     FocusChat,
		loading:   false,
		errors:    []string{},
		chat:      chat.NewChatModel(),
		kanban:    kanban.NewKanbanModel(),
		popup:     popup.NewPopupModel(),
		ctx:       ctx,
		cancelCtx: cancel,
	}
}

func InitialModel(orch *orchestrator.Orchestrator) Model {
	m := New()
	m.orchestrator = orch
	m.phase = orch.GetPhase()
	m.chat = chat.NewChatModel()
	m.kanban = kanban.NewKanbanModel()
	m.popup = popup.NewPopupModel()
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			if m.orchestrator != nil {
				err := m.orchestrator.Init()
				return initCompleteMsg{err: err}
			}
			return initCompleteMsg{err: nil}
		},
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case initCompleteMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err.Error())
		}
		m.loading = false
		if m.orchestrator != nil {
			tickets, _ := m.orchestrator.GetTickets()
			m.tickets = tickets
			if m.kanban != nil {
				m.kanban.SetTickets(tickets)
			}
		}
		return m, nil

	case phaseCompleteMsg:
		m.runningPhase = false
		if msg.err != nil {
			m.errors = append(m.errors, msg.err.Error())
		} else {
			m.phase = msg.phase
			if m.orchestrator != nil {
				m.orchestrator.SetPhase(msg.phase)
				m.orchestrator.PersistState()
			}
		}
		// Reload tickets after phase changes that may create them
		if m.orchestrator != nil {
			tickets, _ := m.orchestrator.GetTickets()
			m.tickets = tickets
			if m.kanban != nil {
				m.kanban.SetTickets(tickets)
			}
		}
		return m, nil

	case streamChunkMsg:
		if m.chat != nil && !msg.done {
			m.chat.AppendAssistantChunk(msg.content)
		}
		if msg.done {
			m.chat.SetStreaming(false)
		}
		return m, nil

	case crushStreamMsg:
		return m.handleCrushStreamEvent(msg.event)

	case crushStreamStartMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err.Error())
			m.chat.SetStreaming(false)
			return m, nil
		}
		m.crushEvents = msg.events
		m.chat.SetStreaming(true)
		return m.handleCrushStreamEvent(msg.event)

	case ticketEventsStartMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err.Error())
			return m, nil
		}
		m.ticketEvents = msg.events
		return m, m.waitForNextTicketEvent()

	case orchestratorTicketEventMsg:
		return m.handleOrchestratorTicketEvent(msg.event)

	case ticketEventMsg:
		if m.orchestrator != nil {
			tickets, _ := m.orchestrator.GetTickets()
			m.tickets = tickets
			if m.kanban != nil {
				m.kanban.SetTickets(tickets)
			}
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		return m.handleWindowResize(msg)

	case spinnerTickMsg:
		if m.chat != nil && m.chat.IsStreaming() {
			return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return spinnerTickMsg{}
			})
		}
		return m, nil
	}

	return m, nil
}

func (m Model) handleWindowResize(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	m.width = msg.Width
	m.height = msg.Height
	if m.kanban != nil {
		m.kanban.SetSize(m.width, m.height)
	}
	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.popup != nil && m.popup.IsOpen() {
		newPopup, cmd := m.popup.Update(msg)
		m.popup = newPopup
		return m, cmd
	}

	keyStr := msg.String()
	switch keyStr {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "?":
		m.helpVisible = !m.helpVisible
	case "tab":
		m.cycleFocus()
	case "shift+tab":
		m.cycleFocusBack()
	case "1":
		m.focus = FocusChat
	case "2":
		m.focus = FocusKanban
	case "3":
		m.focus = FocusPreview
	case "enter":
		if m.focus == FocusChat && m.phase == state.PhaseInit {
			return m, m.startDesignPhase()
		}
		if m.focus == FocusChat && m.chat != nil {
			inputText := m.chat.GetInput()
			if inputText != "" && m.orchestrator != nil && m.orchestrator.GetSessionID() != "" {
				newChat, cmd := m.chat.Update(msg)
				m.chat = newChat
				m.chat.SetStreaming(true)
				return m, tea.Batch(cmd, m.waitForCrushStream(func(ctx context.Context) (<-chan crush.StreamEvent, error) {
					return m.orchestrator.SendMessage(ctx, inputText)
				}))
			}
			newChat, cmd := m.chat.Update(msg)
			m.chat = newChat
			return m, cmd
		}
		if m.focus == FocusKanban && m.kanban != nil {
			tk := m.kanban.SelectedTicket()
			if tk != nil && m.popup != nil {
				m.popup.Open(tk)
			}
		}
	case "d":
		if !m.runningPhase && m.phase == state.PhaseInit {
			m.phase = state.PhaseDesign
			m.runningPhase = true
			return m, m.startDesignPhase()
		}
	case "p":
		if !m.runningPhase && m.phase == state.PhaseDesign {
			m.phase = state.PhasePlan
			m.runningPhase = true
			return m, m.startPlanPhase()
		}
	case "v":
		if !m.runningPhase && m.phase == state.PhasePlan {
			m.phase = state.PhaseValidate
			m.runningPhase = true
			return m, m.startValidatePhase()
		}
	case "x":
		if !m.runningPhase && m.phase == state.PhaseValidate {
			m.phase = state.PhaseExecution
			m.runningPhase = true
			return m, m.startExecutePhase()
		}
	case "y":
		if !m.runningPhase && m.phase == state.PhaseExecution {
			m.phase = state.PhaseVerify
			m.runningPhase = true
			return m, m.startVerifyPhase()
		}
	case "c":
		if !m.runningPhase && m.phase == state.PhaseVerify {
			m.phase = state.PhaseCloseGaps
			m.runningPhase = true
			return m, m.startCloseGapsPhase()
		}
	case "g":
		if !m.runningPhase && m.phase == state.PhaseCloseGaps {
			m.phase = state.PhaseDone
			m.runningPhase = true
			return m, m.startDonePhase()
		}
	case "esc":
		if m.popup != nil && m.popup.IsOpen() {
			m.popup.Close()
		} else {
			m.helpVisible = false
		}
	default:
		switch m.focus {
		case FocusChat:
			if m.chat != nil {
				newChat, cmd := m.chat.Update(msg)
				m.chat = newChat
				return m, cmd
			}
		case FocusKanban:
			if m.kanban != nil {
				newKanban, cmd := m.kanban.Update(msg)
				m.kanban = newKanban
				return m, cmd
			}
		}
	}

	return m, nil
}

func (m Model) startDesignPhase() tea.Cmd {
	return m.waitForCrushStream(m.orchestrator.RunDesign)
}

func (m Model) startPlanPhase() tea.Cmd {
	return m.waitForCrushStream(m.orchestrator.RunPlan)
}

func (m Model) startValidatePhase() tea.Cmd {
	return m.waitForCrushStream(m.orchestrator.RunValidate)
}

func (m Model) startExecutePhase() tea.Cmd {
	return func() tea.Msg {
		if m.orchestrator == nil {
			return phaseCompleteMsg{phase: state.PhaseExecution, err: nil}
		}

		events, err := m.orchestrator.RunExecute(m.ctx)
		if err != nil {
			return phaseCompleteMsg{phase: state.PhaseExecution, err: err}
		}

		return ticketEventsStartMsg{events: events, err: nil}
	}
}

func (m Model) startVerifyPhase() tea.Cmd {
	return m.waitForCrushStream(m.orchestrator.RunVerify)
}

func (m Model) startCloseGapsPhase() tea.Cmd {
	return m.waitForCrushStream(m.orchestrator.RunCloseGaps)
}

func (m Model) startDonePhase() tea.Cmd {
	return func() tea.Msg {
		if m.orchestrator != nil {
			m.orchestrator.SetPhase(state.PhaseDone)
		}
		return phaseCompleteMsg{phase: state.PhaseDone, err: nil}
	}
}

func (m Model) waitForNextTicketEvent() tea.Cmd {
	return func() tea.Msg {
		if m.ticketEvents == nil {
			return nil
		}
		event, ok := <-m.ticketEvents
		if !ok {
			return nil
		}
		return orchestratorTicketEventMsg{event: event}
	}
}

func (m Model) waitForCrushStream(runFunc func(context.Context) (<-chan crush.StreamEvent, error)) tea.Cmd {
	return func() tea.Msg {
		if m.orchestrator == nil {
			return crushStreamStartMsg{err: nil, event: crush.StreamEvent{Type: crush.EventTypeDone}}
		}

		events, err := runFunc(m.ctx)
		if err != nil {
			return crushStreamStartMsg{err: err}
		}

		event, ok := <-events
		if !ok {
			return crushStreamStartMsg{events: events, event: crush.StreamEvent{Type: crush.EventTypeDone}}
		}
		return crushStreamStartMsg{events: events, event: event}
	}
}

func (m Model) waitForNextCrushEvent() tea.Cmd {
	return func() tea.Msg {
		if m.crushEvents == nil {
			return crushStreamMsg{event: crush.StreamEvent{Type: crush.EventTypeDone}}
		}
		event, ok := <-m.crushEvents
		if !ok {
			return crushStreamMsg{event: crush.StreamEvent{Type: crush.EventTypeDone}}
		}
		return crushStreamMsg{event: event}
	}
}

func (m Model) handleCrushStreamEvent(event crush.StreamEvent) (tea.Model, tea.Cmd) {
	switch event.Type {
	case crush.EventTypeStdout:
		if m.chat != nil {
			m.chat.AppendAssistantChunk(event.Content)
		}
		return m, m.waitForNextCrushEvent()
	case crush.EventTypeStderr:
		if m.chat != nil {
			m.chat.AppendErrorChunk(event.Content)
		}
		m.errors = append(m.errors, event.Content)
		return m, m.waitForNextCrushEvent()
	case crush.EventTypeError:
		if m.chat != nil {
			m.chat.AppendErrorChunk(event.Content)
		}
		m.errors = append(m.errors, event.Content)
		m.chat.SetStreaming(false)
		return m, nil
	case crush.EventTypeDone:
		m.chat.SetStreaming(false)
		m.crushEvents = nil
		if m.orchestrator != nil {
			m.orchestrator.CaptureSessionID()
		}
		return m, func() tea.Msg {
			return phaseCompleteMsg{phase: m.phase, err: nil}
		}
	default:
		return m, m.waitForNextCrushEvent()
	}
}

func (m Model) handleOrchestratorTicketEvent(event orchestrator.TicketEvent) (tea.Model, tea.Cmd) {
	if m.kanban != nil {
		m.kanban.UpdateTicketStatus(event.TicketID, event.Status)
	}
	return m, m.waitForNextTicketEvent()
}

func (m *Model) cycleFocus() {
	switch m.focus {
	case FocusChat:
		m.focus = FocusKanban
	case FocusKanban:
		m.focus = FocusPreview
	case FocusPreview:
		m.focus = FocusChat
	default:
		m.focus = FocusChat
	}
}

func (m *Model) cycleFocusBack() {
	switch m.focus {
	case FocusChat:
		m.focus = FocusPreview
	case FocusKanban:
		m.focus = FocusChat
	case FocusPreview:
		m.focus = FocusKanban
	default:
		m.focus = FocusChat
	}
}

func (m Model) View() tea.View {
	if m.popup != nil && m.popup.IsOpen() {
		baseContent := m.viewForPhase()
		return tea.NewView(m.renderPopupOverlay(baseContent.Content))
	}
	return m.viewForPhase()
}

func (m Model) viewForPhase() tea.View {
	var b strings.Builder

	b.WriteString(m.renderHeader())
	b.WriteString("\n")

	contentHeight := max(m.height-4, 10)


	switch m.phase {
	case state.PhaseInit:
		b.WriteString(m.viewInit(contentHeight))
	case state.PhaseDesign:
		b.WriteString(m.viewDesign(contentHeight))
	case state.PhasePlan:
		b.WriteString(m.viewPlan(contentHeight))
	case state.PhaseValidate:
		b.WriteString(m.viewValidate(contentHeight))
	case state.PhaseExecution:
		b.WriteString(m.viewExecution(contentHeight))
	case state.PhaseVerify:
		b.WriteString(m.viewVerify(contentHeight))
	case state.PhaseCloseGaps:
		b.WriteString(m.viewCloseGaps(contentHeight))
	case state.PhaseDone:
		b.WriteString(m.viewDone(contentHeight))
	default:
		b.WriteString(m.viewInit(contentHeight))
	}

	b.WriteString("\n")
	b.WriteString(m.renderFooter())

	if m.helpVisible {
		b.WriteString("\n")
		b.WriteString(m.renderHelp())
	}

	return tea.NewView(b.String())
}

func (m Model) renderHeader() string {
	phaseText := m.phaseDisplayName()
	phaseDisplay := phaseStyle.Render(phaseText)

	title := titleStyle.Render("🚂 choo-choo")

	spacer := ""
	available := m.width - lipgloss.Width(title) - lipgloss.Width(phaseDisplay) - 4
	if available > 0 {
		spacer = strings.Repeat(" ", available)
	}

	return fmt.Sprintf("%s%s%s", title, spacer, phaseDisplay)
}

func (m Model) phaseDisplayName() string {
	switch m.phase {
	case state.PhaseInit:
		return "⚡ Initialize"
	case state.PhaseDesign:
		return "🎨 Design"
	case state.PhasePlan:
		return "📋 Plan"
	case state.PhaseValidate:
		return "✓ Validate"
	case state.PhaseExecution:
		return "⚡ Execute"
	case state.PhaseVerify:
		return "🔍 Verify"
	case state.PhaseCloseGaps:
		return "🔧 Close Gaps"
	case state.PhaseDone:
		return "✅ Done"
	default:
		return string(m.phase)
	}
}

func (m Model) renderFooter() string {
	var focusArea, focusHint string
	switch m.focus {
	case FocusChat:
		focusArea = "Chat"
		focusHint = "type message"
	case FocusKanban:
		focusArea = "Kanban"
		focusHint = "j/k navigate"
	case FocusPreview:
		focusArea = "Preview"
		focusHint = "read-only"
	}

	focusAreaStyled := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render(focusArea)
	focusIndicator := fmt.Sprintf("[%s - %s]", focusAreaStyled, focusHint)

	help := helpStyle.Render("?: help | tab: switch | q: quit | " + focusIndicator)
	return help
}

func (m Model) renderHelp() string {
	var b strings.Builder
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("                         KEYBOARD SHORTCUTS\n")
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	b.WriteString("  1/2/3        Switch focus (Chat/Kanban/Preview)\n")
	b.WriteString("  Tab          Cycle focus forward\n")
	b.WriteString("  Shift+Tab    Cycle focus backward\n")
	b.WriteString("  Enter        Confirm / Open ticket details / Send message\n")
	b.WriteString("  Escape       Close popup / help\n")
	b.WriteString("  ?            Toggle this help\n")
	b.WriteString("  --- Phase shortcuts ---\n")
	b.WriteString("  d            Start Design phase\n")
	b.WriteString("  p            Start Plan phase\n")
	b.WriteString("  v            Start Validate phase\n")
	b.WriteString("  x            Start Execute phase\n")
	b.WriteString("  y            Start Verify phase\n")
	b.WriteString("  c            Start Close-Gaps phase\n")
	b.WriteString("  g            Complete to Done phase\n")
	b.WriteString("  --- Navigation ---\n")
	b.WriteString("  j/k          Navigate kanban (down/up)\n")
	b.WriteString("  h/l          Navigate kanban (left/right)\n")
	b.WriteString("  q/Ctrl+C     Quit\n")
	b.WriteString("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	return boxStyle.Render(b.String())
}

func (m Model) viewInit(_ int) string {
	var b strings.Builder
	b.WriteString("\n\n")

	boxWidth := max(min(m.width-4, 60), 20)

	welcomeBox := boxStyle.
		Width(boxWidth).
		Render(`
   Welcome to choo-choo!
   
   The workflow orchestration TUI for Crush.
   
   Phases:
   1. Design    - Chat with AI to design your feature
   2. Plan      - Break down into hierarchical tickets
   3. Validate  - Verify plan is executable
   4. Execute   - Run tickets in parallel or sequence
   5. Verify    - Confirm implementation matches criteria
   
   Press 'd' to start Design phase
   or press '?' for keyboard shortcuts
`)
	b.WriteString(welcomeBox)

	if len(m.errors) > 0 {
		b.WriteString("\n")
		errBox := boxStyle.
			BorderForeground(errorColor).
			Render("Errors:\n" + strings.Join(m.errors, "\n"))
		b.WriteString(errBox)
	}

	return b.String()
}

func (m Model) viewDesign(height int) string {
	chatWidth := max(m.width*2/3, 20)
	previewWidth := max(m.width-chatWidth-4, 10)


	chatContent := ""
	if m.chat != nil {
		m.chat.SetSize(chatWidth-4, height-4)
		chatContent = m.chat.View().Content
	}

	chatStyle := boxStyle
	if m.focus == FocusChat {
		chatStyle = focusedBoxStyle
	}
	chatBox := chatStyle.
		Width(chatWidth).
		Height(height).
		Render(chatContent)

	previewStyle := boxStyle
	if m.focus == FocusPreview {
		previewStyle = focusedBoxStyle
	}
	previewContent := m.renderPreviewContent()
	previewBox := previewStyle.
		Width(previewWidth).
		Height(height).
		Render(previewContent)

	return lipgloss.JoinHorizontal(lipgloss.Top, chatBox, " ", previewBox)
}

func (m Model) viewPlan(height int) string {
	maxLines := height - 4
	if maxLines < 10 {
		maxLines = 10
	}
	treeContent := m.renderTicketTree(maxLines)
	width := max(m.width-4, 20)
	style := boxStyle
	if m.focus == FocusKanban {
		style = focusedBoxStyle
	}
	treeBox := style.
		Width(width).
		Height(height).
		Render(treeContent)

	return treeBox
}

func (m Model) viewValidate(height int) string {
	var b strings.Builder

	b.WriteString("Validation Results\n\n")

	stats := fmt.Sprintf("Total tickets: %d\n", len(m.tickets))
	b.WriteString(stats)

	if len(m.tickets) > 0 {
		ready := 0
		blocked := 0
		done := 0
		for _, t := range m.tickets {
			switch t.Status {
			case ticket.StatusClosed:
				done++
			case ticket.StatusOpen:
				if len(t.Dependencies) == 0 {
					ready++
				} else {
					blocked++
				}
			}
		}
		b.WriteString(fmt.Sprintf("  Ready: %d\n", ready))
		b.WriteString(fmt.Sprintf("  Blocked: %d\n", blocked))
		b.WriteString(fmt.Sprintf("  Done: %d\n", done))
	}

	b.WriteString("\nPress 'x' to start Execute phase\n")

	width := max(m.width-4, 20)
	style := boxStyle
	if m.focus == FocusKanban {
		style = focusedBoxStyle
	}
	return style.
		Width(width).
		Height(height).
		Render(b.String())
}

func (m Model) viewExecution(_ int) string {
	style := boxStyle
	if m.focus == FocusKanban {
		style = focusedBoxStyle
	}
	if m.kanban != nil {
		var b strings.Builder
		
		total, done := m.kanban.GetProgress()
		if total > 0 {
			percentage := float64(done) / float64(total) * 100
			progressBarWidth := 20
			filled := int(float64(progressBarWidth) * float64(done) / float64(total))
			
			progressBar := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#10B981")).
				Render(strings.Repeat("█", filled)) +
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("#3B3B5C")).
					Render(strings.Repeat("░", progressBarWidth-filled))
			
			progressText := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#9CA3AF")).
				Render(fmt.Sprintf(" %d/%d tickets (%.0f%%)", done, total, percentage))
			
			b.WriteString(progressBar + progressText + "\n\n")
		}
		
		kanbanView := m.kanban.View()
		b.WriteString(kanbanView.Content)
		
		return style.Render(b.String())
	}
	return style.Render("No tickets loaded")
}

func (m Model) viewVerify(height int) string {
	width := max(m.width-4, 20)
	style := boxStyle
	if m.focus == FocusKanban {
		style = focusedBoxStyle
	}
	return style.
		Width(width).
		Height(height).
		Render("Verification Phase\n\nReviewing completed tickets...")
}

func (m Model) viewCloseGaps(height int) string {
	width := max(m.width-4, 20)
	style := boxStyle
	if m.focus == FocusKanban {
		style = focusedBoxStyle
	}
	return style.
		Width(width).
		Height(height).
		Render("Gap Closure Phase\n\nFixing discrepancies found during verification...")
}

func (m Model) viewDone(height int) string {
	width := max(m.width-4, 20)
	return boxStyle.
		BorderForeground(successColor).
		Width(width).
		Height(height).
		Render("✅ All Done!\n\nThe choo-choo workflow has completed successfully.")
}

func (m Model) renderPreviewContent() string {
	var b strings.Builder
	b.WriteString("Artifacts\n\n")

	if len(m.tickets) == 0 {
		b.WriteString("No artifacts yet.\n")
		b.WriteString("Complete the design phase\n")
		b.WriteString("to generate artifacts.")
	} else {
		b.WriteString(fmt.Sprintf("Tickets: %d\n", len(m.tickets)))

		epics := 0
		stories := 0
		tasks := 0
		for _, t := range m.tickets {
			switch t.Type {
			case ticket.TypeEpic:
				epics++
			case ticket.TypeStory:
				stories++
			case ticket.TypeTask:
				tasks++
			}
		}
		b.WriteString(fmt.Sprintf("  Epics: %d\n", epics))
		b.WriteString(fmt.Sprintf("  Stories: %d\n", stories))
		b.WriteString(fmt.Sprintf("  Tasks: %d\n", tasks))
	}

	return b.String()
}

func (m Model) renderTicketTree(maxLines int) string {
	var b strings.Builder
	b.WriteString("Ticket Tree\n\n")
	lineCount := 2

	if len(m.tickets) == 0 {
		b.WriteString("No tickets yet.\n")
		b.WriteString("Press 'p' to generate tickets from design.")
		return b.String()
	}

	epics := []ticket.Ticket{}
	stories := []ticket.Ticket{}
	tasks := []ticket.Ticket{}

	for _, t := range m.tickets {
		switch t.Type {
		case ticket.TypeEpic:
			epics = append(epics, t)
		case ticket.TypeStory:
			stories = append(stories, t)
		case ticket.TypeTask:
			tasks = append(tasks, t)
		}
	}

	totalLines := 0
	for _, epic := range epics {
		totalLines++
		for _, story := range stories {
			if story.Parent == epic.ID {
				totalLines++
				for _, task := range tasks {
					if task.Parent == story.ID {
						totalLines++
					}
				}
			}
		}
	}

	renderedLines := 0
	for _, epic := range epics {
		if lineCount >= maxLines-1 {
			remaining := totalLines - renderedLines
			b.WriteString(fmt.Sprintf("\n... and %d more tickets\n", remaining))
			return b.String()
		}
		b.WriteString(fmt.Sprintf("📦 %s: %s\n", epic.ID, epic.Title))
		lineCount++
		renderedLines++

		for _, story := range stories {
			if story.Parent == epic.ID {
				if lineCount >= maxLines-1 {
					remaining := totalLines - renderedLines
					b.WriteString(fmt.Sprintf("\n... and %d more tickets\n", remaining))
					return b.String()
				}
				b.WriteString(fmt.Sprintf("  ├─ 📋 %s: %s\n", story.ID, story.Title))
				lineCount++
				renderedLines++

				for i, task := range tasks {
					if task.Parent == story.ID {
						if lineCount >= maxLines-1 {
							remaining := totalLines - renderedLines
							b.WriteString(fmt.Sprintf("\n... and %d more tickets\n", remaining))
							return b.String()
						}
						prefix := "  │  ├─"
						if i == len(tasks)-1 {
							prefix = "  │  └─"
						}
						status := "⬜"
						switch task.Status {
						case ticket.StatusInProgress:
							status = "🔄"
						case ticket.StatusClosed:
							status = "✅"
						}
						b.WriteString(fmt.Sprintf("%s %s %s: %s\n", prefix, status, task.ID, task.Title))
						lineCount++
						renderedLines++
					}
				}
			}
		}
	}

	return b.String()
}

func (m Model) renderPopupOverlay(baseContent string) string {
	if m.popup == nil {
		return baseContent
	}

	popupWidth := max(m.width*2/3, 40)
	popupHeight := max(m.height*2/3, 10)

	popupContent := m.popup.View(popupWidth, popupHeight).Content

	lines := strings.Split(baseContent, "\n")
	contentLineCount := len(lines)

	popupLines := strings.Split(popupContent, "\n")
	popupLineCount := len(popupLines)

	startRow := max((contentLineCount-popupLineCount)/2, 0)

	var result strings.Builder
	for i, line := range lines {
		if i >= startRow && i < startRow+popupLineCount {
			popupLineIdx := i - startRow
			if popupLineIdx < len(popupLines) {
				result.WriteString(popupLines[popupLineIdx])
			}
		} else {
			result.WriteString(line)
		}
		result.WriteString("\n")
	}

	return result.String()
}

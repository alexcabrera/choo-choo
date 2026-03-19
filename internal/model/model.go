package model

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/alexcabrera/choo-choo/internal/orchestrator"
	"github.com/alexcabrera/choo-choo/internal/state"
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

type Model struct {
	orchestrator *orchestrator.Orchestrator
	chat         *chat.ChatModel
	kanban       *kanban.KanbanModel
	popup        *popup.PopupModel
	phase        state.Phase
	width        int
	height       int
	focus        FocusArea
	loading      bool
	errors       []string
}

func New() Model {
	return Model{
		phase:   state.PhaseInit,
		focus:   FocusChat,
		loading: false,
		errors:  []string{},
	}
}

func InitialModel(orch *orchestrator.Orchestrator) Model {
	m := New()
	m.orchestrator = orch
	m.chat = chat.NewChatModel()
	m.kanban = kanban.NewKanbanModel()
	m.popup = popup.NewPopupModel()
	m.phase = orch.GetPhase()
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
		return m, nil

	case PhaseChangeMsg:
		return m.handlePhaseChange(msg)

	case FocusChangeMsg:
		m.focus = msg.Focus
		return m, nil

	case TicketUpdateMsg:
		return m.handleTicketUpdate(msg)

	case PopupOpenMsg:
		return m.handlePopupOpen(msg)

	case PopupCloseMsg:
		if m.popup != nil {
			m.popup.Close()
		}
		return m, nil

	case ErrorMsg:
		m.errors = append(m.errors, msg.Err.Error())
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

func (m Model) handlePhaseChange(msg PhaseChangeMsg) (tea.Model, tea.Cmd) {
	newPhase := state.Phase(msg.Phase)
	if !newPhase.IsValid() {
		return m, nil
	}
	m.phase = newPhase
	return m, nil
}

func (m Model) handleTicketUpdate(msg TicketUpdateMsg) (tea.Model, tea.Cmd) {
	if m.orchestrator != nil {
		tickets, _ := m.orchestrator.GetTickets()
		if m.kanban != nil {
			m.kanban.SetTickets(tickets)
		}
	}
	return m, nil
}

func (m Model) handlePopupOpen(msg PopupOpenMsg) (tea.Model, tea.Cmd) {
	if m.orchestrator != nil && m.popup != nil {
		tickets, err := m.orchestrator.GetTickets()
		if err != nil {
			return m, nil
		}
		for _, t := range tickets {
			if t.ID == msg.TicketID {
				m.popup.Open(&t)
				break
			}
		}
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

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "tab":
		m.cycleFocus()
	case "1":
		m.focus = FocusChat
	case "2":
		m.focus = FocusKanban
	}

	return m, nil
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

func (m Model) View() tea.View {
	if m.popup != nil && m.popup.IsOpen() {
		baseContent := m.viewForPhase()
		return tea.NewView(m.renderPopupOverlay(baseContent.Content))
	}
	return m.viewForPhase()
}

func (m Model) viewForPhase() tea.View {
	switch m.phase {
	case state.PhaseInit:
		return m.viewInit()
	case state.PhaseDesign:
		return m.viewDesign()
	case state.PhasePlan:
		return m.viewPlan()
	case state.PhaseExecution:
		return m.viewExecution()
	default:
		return tea.NewView("choo-choo\n\nUnknown phase")
	}
}

func (m Model) viewInit() tea.View {
	var b strings.Builder
	b.WriteString("choo-choo - Initialization\n\n")
	if m.loading {
		b.WriteString("Loading...")
	} else if len(m.errors) > 0 {
		b.WriteString("Errors:\n")
		for _, e := range m.errors {
			b.WriteString(fmt.Sprintf("  - %s\n", e))
		}
	} else {
		b.WriteString("Ready. Press '1' for chat, '2' for kanban.")
	}
	return tea.NewView(b.String())
}

func (m Model) viewDesign() tea.View {
	var b strings.Builder
	b.WriteString("choo-choo - Design Phase\n\n")
	b.WriteString("Focus: ")
	if m.focus == FocusChat && m.chat != nil {
		b.WriteString("Chat\n")
		b.WriteString(fmt.Sprintf("Messages: %d\n", len(m.chat.GetInput())))
	} else {
		b.WriteString("Other\n")
	}
	return tea.NewView(b.String())
}

func (m Model) viewPlan() tea.View {
	var b strings.Builder
	b.WriteString("choo-choo - Plan Phase\n\n")
	b.WriteString("Focus: ")
	switch m.focus {
	case FocusChat:
		b.WriteString("Chat")
	case FocusKanban:
		b.WriteString("Kanban")
	case FocusPreview:
		b.WriteString("Preview")
	}
	b.WriteString("\n")
	return tea.NewView(b.String())
}

func (m Model) viewExecution() tea.View {
	if m.kanban != nil {
		return m.kanban.View()
	}
	return tea.NewView("choo-choo - Execution Phase\n\nNo kanban loaded")
}

func (m Model) renderPopupOverlay(baseContent string) string {
	if m.popup == nil {
		return baseContent
	}

	popupWidth := m.width * 2 / 3
	if popupWidth < 40 {
		popupWidth = 40
	}
	popupHeight := m.height * 2 / 3
	if popupHeight < 10 {
		popupHeight = 10
	}

	var b strings.Builder
	b.WriteString(baseContent)
	b.WriteString("\n\n")

	header := fmt.Sprintf("=== %s ===", m.popup.GetActiveTab().String())
	b.WriteString(header + "\n")
	b.WriteString(m.popup.GetContent())

	return b.String()
}

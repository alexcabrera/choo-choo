package model

import (
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
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	return m, nil
}

func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	}

	return m, nil
}

func (m Model) View() tea.View {
	switch m.phase {
	case state.PhaseInit:
		return tea.NewView("choo-choo - Press 'q' to quit\n\nNo project initialized.")
	case state.PhaseDesign:
		return tea.NewView("choo-choo - Design Phase\n\nTODO: Implement design UI")
	case state.PhasePlan:
		return tea.NewView("choo-choo - Plan Phase\n\nTODO: Implement plan UI")
	case state.PhaseExecution:
		return tea.NewView("choo-choo - Execute Phase\n\nTODO: Implement kanban UI")
	default:
		return tea.NewView("choo-choo\n\nTODO: Implement UI")
	}
}

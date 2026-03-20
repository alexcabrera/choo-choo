package kanban

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"

	"github.com/alexcabrera/choo-choo/internal/ticket"
)

const (
	ColumnTodo = iota
	ColumnDoing
	ColumnDone
)

var (
	todoColor    = lipgloss.Color("#6B7280")
	doingColor   = lipgloss.Color("#F59E0B")
	doneColor    = lipgloss.Color("#10B981")

	columnColors = [3]lipgloss.Color{todoColor, doingColor, doneColor}
	columnNames  = [3]string{"TODO", "DOING", "DONE"}
)

type KanbanModel struct {
	columns   [3][]ticket.Ticket
	cursorCol int
	cursorRow int
	width     int
	height    int
}

func NewKanbanModel() *KanbanModel {
	return &KanbanModel{
		columns:   [3][]ticket.Ticket{{}, {}, {}},
		cursorCol: 0,
		cursorRow: 0,
	}
}

func (m *KanbanModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *KanbanModel) Update(msg tea.Msg) (*KanbanModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m *KanbanModel) handleKeyPress(msg tea.KeyMsg) (*KanbanModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.MoveCursorUp()
	case "down", "j":
		m.MoveCursorDown()
	case "left", "h":
		m.MoveCursorLeft()
	case "right", "l":
		m.MoveCursorRight()
	case "enter":
		return m, nil
	}
	return m, nil
}

func (m *KanbanModel) SetTickets(tickets []ticket.Ticket) {
	m.columns = [3][]ticket.Ticket{{}, {}, {}}
	for _, t := range tickets {
		switch t.Status {
		case ticket.StatusOpen:
			m.columns[ColumnTodo] = append(m.columns[ColumnTodo], t)
		case ticket.StatusInProgress:
			m.columns[ColumnDoing] = append(m.columns[ColumnDoing], t)
		case ticket.StatusClosed:
			m.columns[ColumnDone] = append(m.columns[ColumnDone], t)
		}
	}
}

func (m *KanbanModel) UpdateTicketStatus(ticketID string, status string) {
	var targetCol int
	var newStatus ticket.Status

	switch status {
	case "executing":
		targetCol = ColumnDoing
		newStatus = ticket.StatusInProgress
	case "completed", "failed":
		targetCol = ColumnDone
		newStatus = ticket.StatusClosed
	default:
		return
	}

	for fromCol := range 3 {
		for i, t := range m.columns[fromCol] {
			if t.ID == ticketID {
				if fromCol == targetCol {
					m.columns[fromCol][i].Status = newStatus
					return
				}
				m.columns[fromCol] = append(m.columns[fromCol][:i], m.columns[fromCol][i+1:]...)
				t.Status = newStatus
				m.columns[targetCol] = append(m.columns[targetCol], t)
				if m.cursorCol == fromCol && m.cursorRow >= len(m.columns[fromCol]) && len(m.columns[fromCol]) > 0 {
					m.cursorRow = len(m.columns[fromCol]) - 1
				}
				return
			}
		}
	}
}

func (m *KanbanModel) SelectedTicket() *ticket.Ticket {
	col := m.columns[m.cursorCol]
	if m.cursorRow < 0 || m.cursorRow >= len(col) {
		return nil
	}
	return &col[m.cursorRow]
}

func (m *KanbanModel) MoveTicket(toCol int) bool {
	if toCol < 0 || toCol > 2 {
		return false
	}
	if m.cursorCol == toCol {
		return false
	}

	tk := m.SelectedTicket()
	if tk == nil {
		return false
	}

	fromCol := m.cursorCol
	row := m.cursorRow

	m.columns[fromCol] = append(m.columns[fromCol][:row], m.columns[fromCol][row+1:]...)

	switch toCol {
	case ColumnTodo:
		tk.Status = ticket.StatusOpen
	case ColumnDoing:
		tk.Status = ticket.StatusInProgress
	case ColumnDone:
		tk.Status = ticket.StatusClosed
	}

	m.columns[toCol] = append(m.columns[toCol], *tk)

	if m.cursorRow >= len(m.columns[fromCol]) && len(m.columns[fromCol]) > 0 {
		m.cursorRow = len(m.columns[fromCol]) - 1
	}

	return true
}

func (m *KanbanModel) MoveCursorUp() {
	if m.cursorRow > 0 {
		m.cursorRow--
	}
}

func (m *KanbanModel) MoveCursorDown() {
	col := m.columns[m.cursorCol]
	if m.cursorRow < len(col)-1 {
		m.cursorRow++
	}
}

func (m *KanbanModel) MoveCursorLeft() {
	if m.cursorCol > 0 {
		m.cursorCol--
		m.cursorRow = 0
	}
}

func (m *KanbanModel) MoveCursorRight() {
	if m.cursorCol < 2 {
		m.cursorCol++
		m.cursorRow = 0
	}
}

func (m *KanbanModel) GetSelectedTicketID() string {
	tk := m.SelectedTicket()
	if tk == nil {
		return ""
	}
	return tk.ID
}

func (m *KanbanModel) GetProgress() (total int, done int) {
	total = len(m.columns[ColumnTodo]) + len(m.columns[ColumnDoing]) + len(m.columns[ColumnDone])
	done = len(m.columns[ColumnDone])
	return total, done
}

func (m *KanbanModel) View() tea.View {
	if m.width < 60 {
		return tea.NewView("Terminal too narrow. Please resize to at least 60 columns.")
	}

	colWidth := min((m.width-6)/3, 30)


	columnViews := make([]string, 3)
	for i := range m.columns {
		columnViews[i] = m.renderColumn(i, colWidth)
	}

	return tea.NewView(lipgloss.JoinHorizontal(lipgloss.Top, columnViews...))
}

func (m *KanbanModel) renderColumn(colIdx int, width int) string {
	var b strings.Builder

	color := columnColors[colIdx]

	name := columnNames[colIdx]
	count := len(m.columns[colIdx])
	headerText := fmt.Sprintf("%s (%d)", name, count)

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(color).
		Padding(0, 1).
		Render(headerText)

	b.WriteString(header)
	b.WriteString("\n")

	sep := lipgloss.NewStyle().
		Foreground(color).
		Render(strings.Repeat("─", width-2))
	b.WriteString(sep)
	b.WriteString("\n")

	contentHeight := max(m.height-6, 5)


	lines := 0
	for rowIdx, t := range m.columns[colIdx] {
		if lines >= contentHeight {
			break
		}

		cardText := m.renderTicketCard(t, colIdx, rowIdx, width-4)
		b.WriteString(cardText)
		b.WriteString("\n")
		lines++
	}

	if len(m.columns[colIdx]) == 0 {
		empty := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#4A4A4A")).
			Italic(true).
			Padding(1, 1).
			Render("No tickets")
		b.WriteString(empty)
	}

	columnBox := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color).
		Padding(0, 1).
		Render(b.String())

	return columnBox
}

func (m *KanbanModel) renderTicketCard(t ticket.Ticket, colIdx, rowIdx, maxWidth int) string {
	isSelected := colIdx == m.cursorCol && rowIdx == m.cursorRow

	var typeIcon string
	switch t.Type {
	case ticket.TypeEpic:
		typeIcon = "📦"
	case ticket.TypeStory:
		typeIcon = "📋"
	case ticket.TypeTask:
		typeIcon = "✓"
	default:
		typeIcon = "•"
	}

	title := t.Title
	if len(title) > maxWidth-8 {
		title = title[:maxWidth-11] + "..."
	}

	cardText := fmt.Sprintf("%s %s: %s", typeIcon, t.ID, title)

	if isSelected {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#3D3D5C")).
			Bold(true).
			Padding(0, 1).
			Render("▶ " + cardText)
	}

	return lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Padding(0, 1).
		Render("  " + cardText)
}

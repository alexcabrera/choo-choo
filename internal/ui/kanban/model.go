package kanban

import "github.com/alexcabrera/choo-choo/internal/ticket"

const (
	ColumnTodo = iota
	ColumnDoing
	ColumnDone
)

type KanbanModel struct {
	columns   [3][]ticket.Ticket
	cursorCol int
	cursorRow int
}

func NewKanbanModel() *KanbanModel {
	return &KanbanModel{
		columns:   [3][]ticket.Ticket{{}, {}, {}},
		cursorCol: 0,
		cursorRow: 0,
	}
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

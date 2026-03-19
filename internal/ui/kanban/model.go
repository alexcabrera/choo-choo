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

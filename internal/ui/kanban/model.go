package kanban

import "github.com/alexcabrera/choo-choo/internal/ticket"

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

package kanban

import (
	"testing"

	"github.com/alexcabrera/choo-choo/internal/ticket"
)

func TestNewKanbanModel(t *testing.T) {
	m := NewKanbanModel()
	if m == nil {
		t.Fatal("NewKanbanModel() returned nil")
	}
	if m.cursorCol != 0 {
		t.Errorf("cursorCol = %d, want 0", m.cursorCol)
	}
	if m.cursorRow != 0 {
		t.Errorf("cursorRow = %d, want 0", m.cursorRow)
	}
	for i, col := range m.columns {
		if col == nil {
			t.Errorf("columns[%d] is nil", i)
		}
	}
}

func TestSetTickets(t *testing.T) {
	m := NewKanbanModel()
	tickets := []ticket.Ticket{
		{ID: "T-001", Status: ticket.StatusOpen},
		{ID: "T-002", Status: ticket.StatusInProgress},
		{ID: "T-003", Status: ticket.StatusClosed},
		{ID: "T-004", Status: ticket.StatusOpen},
	}
	m.SetTickets(tickets)

	if len(m.columns[ColumnTodo]) != 2 {
		t.Errorf("TODO column has %d tickets, want 2", len(m.columns[ColumnTodo]))
	}
	if len(m.columns[ColumnDoing]) != 1 {
		t.Errorf("DOING column has %d tickets, want 1", len(m.columns[ColumnDoing]))
	}
	if len(m.columns[ColumnDone]) != 1 {
		t.Errorf("DONE column has %d tickets, want 1", len(m.columns[ColumnDone]))
	}
}

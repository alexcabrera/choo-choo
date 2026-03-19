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

func TestSelectedTicket(t *testing.T) {
	m := NewKanbanModel()
	tickets := []ticket.Ticket{
		{ID: "T-001", Status: ticket.StatusOpen},
		{ID: "T-002", Status: ticket.StatusOpen},
	}
	m.SetTickets(tickets)

	selected := m.SelectedTicket()
	if selected == nil {
		t.Fatal("SelectedTicket() returned nil")
	}
	if selected.ID != "T-001" {
		t.Errorf("SelectedTicket().ID = %q, want %q", selected.ID, "T-001")
	}

	m.cursorRow = 1
	selected = m.SelectedTicket()
	if selected == nil {
		t.Fatal("SelectedTicket() returned nil for row 1")
	}
	if selected.ID != "T-002" {
		t.Errorf("SelectedTicket().ID = %q, want %q", selected.ID, "T-002")
	}

	m.cursorRow = 99
	selected = m.SelectedTicket()
	if selected != nil {
		t.Error("SelectedTicket() should return nil for out of bounds row")
	}
}

func TestMoveTicket(t *testing.T) {
	m := NewKanbanModel()
	tickets := []ticket.Ticket{
		{ID: "T-001", Status: ticket.StatusOpen},
	}
	m.SetTickets(tickets)

	if !m.MoveTicket(ColumnDoing) {
		t.Error("MoveTicket() should return true")
	}

	if len(m.columns[ColumnTodo]) != 0 {
		t.Errorf("TODO column has %d tickets, want 0", len(m.columns[ColumnTodo]))
	}
	if len(m.columns[ColumnDoing]) != 1 {
		t.Errorf("DOING column has %d tickets, want 1", len(m.columns[ColumnDoing]))
	}
	if m.columns[ColumnDoing][0].Status != ticket.StatusInProgress {
		t.Errorf("Moved ticket status = %v, want %v", m.columns[ColumnDoing][0].Status, ticket.StatusInProgress)
	}
}

func TestMoveTicketInvalidColumn(t *testing.T) {
	m := NewKanbanModel()
	tickets := []ticket.Ticket{{ID: "T-001", Status: ticket.StatusOpen}}
	m.SetTickets(tickets)

	if m.MoveTicket(99) {
		t.Error("MoveTicket(99) should return false")
	}
	if m.MoveTicket(ColumnTodo) {
		t.Error("MoveTicket to same column should return false")
	}
}

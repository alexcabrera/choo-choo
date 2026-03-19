package ticket

import "testing"

func TestTicketTypeValues(t *testing.T) {
	types := []TicketType{TypeEpic, TypeStory, TypeTask, TypeChore, TypeBug, TypeFeature}
	for _, tt := range types {
		if tt == "" {
			t.Error("TicketType should not be empty")
		}
	}
}

func TestStatusValues(t *testing.T) {
	statuses := []Status{StatusOpen, StatusInProgress, StatusClosed}
	for _, s := range statuses {
		if s == "" {
			t.Error("Status should not be empty")
		}
	}
}

func TestTicketFields(t *testing.T) {
	ticket := Ticket{
		ID:           "T-001",
		Type:         TypeTask,
		Title:        "Test Ticket",
		Description:  "A test ticket",
		Status:       StatusOpen,
		Parent:       "E-001",
		Dependencies: []string{"T-000"},
		Accepts:      []string{"AC1", "AC2"},
		Notes:        []string{"Note 1"},
	}

	if ticket.ID != "T-001" {
		t.Errorf("Ticket.ID = %q, want %q", ticket.ID, "T-001")
	}
	if ticket.Type != TypeTask {
		t.Errorf("Ticket.Type = %v, want %v", ticket.Type, TypeTask)
	}
	if ticket.Status != StatusOpen {
		t.Errorf("Ticket.Status = %v, want %v", ticket.Status, StatusOpen)
	}
	if len(ticket.Dependencies) != 1 {
		t.Errorf("Ticket.Dependencies length = %d, want 1", len(ticket.Dependencies))
	}
	if len(ticket.Accepts) != 2 {
		t.Errorf("Ticket.Accepts length = %d, want 2", len(ticket.Accepts))
	}
}

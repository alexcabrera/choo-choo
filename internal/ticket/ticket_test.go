package ticket

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

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

func TestNewTicketManager(t *testing.T) {
	tm := NewTicketManager("/usr/local/bin/tk", "/workspace/.tickets")
	if tm == nil {
		t.Fatal("NewTicketManager() returned nil")
	}
	if tm.tkPath != "/usr/local/bin/tk" {
		t.Errorf("tkPath = %q, want %q", tm.tkPath, "/usr/local/bin/tk")
	}
	if tm.ticketsDir != "/workspace/.tickets" {
		t.Errorf("ticketsDir = %q, want %q", tm.ticketsDir, "/workspace/.tickets")
	}
}

func TestTicketManagerListEmpty(t *testing.T) {
	tmpDir := t.TempDir()
	tm := NewTicketManager("tk", tmpDir)
	tickets, err := tm.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("List() returned %d tickets, want 0", len(tickets))
	}
}

func TestTicketManagerGetNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	tm := NewTicketManager("tk", tmpDir)
	_, err := tm.Get("nonexistent")
	if err == nil {
		t.Error("Get() should return error for nonexistent ticket")
	}
}

func TestTicketManagerGetAndList(t *testing.T) {
	tmpDir := t.TempDir()
	ticketContent := `---
id: T-001
status: open
deps: []
links: []
created: 2026-03-18T12:00:00Z
type: task
priority: 1
---
# Test Ticket
Some description here.
`
	path := filepath.Join(tmpDir, "T-001.md")
	if err := os.WriteFile(path, []byte(ticketContent), 0644); err != nil {
		t.Fatal(err)
	}

	tm := NewTicketManager("tk", tmpDir)

	ticket, err := tm.Get("T-001")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if ticket.ID != "T-001" {
		t.Errorf("Get() ID = %q, want %q", ticket.ID, "T-001")
	}

	tickets, err := tm.List()
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(tickets) != 1 {
		t.Fatalf("List() returned %d tickets, want 1", len(tickets))
	}
	if tickets[0].ID != "T-001" {
		t.Errorf("List()[0].ID = %q, want %q", tickets[0].ID, "T-001")
	}
}

func TestTicketManagerSetStatus(t *testing.T) {
	tmpDir := t.TempDir()
	ticketContent := `---
id: T-001
status: open
deps: []
---
# Test
`
	path := filepath.Join(tmpDir, "T-001.md")
	if err := os.WriteFile(path, []byte(ticketContent), 0644); err != nil {
		t.Fatal(err)
	}

	tm := NewTicketManager("tk", tmpDir)
	if err := tm.SetStatus("T-001", StatusInProgress); err != nil {
		t.Fatalf("SetStatus() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "status: in_progress") {
		t.Errorf("SetStatus did not update status, content:\n%s", string(data))
	}
}

func TestTicketManagerAddNote(t *testing.T) {
	tmpDir := t.TempDir()
	ticketContent := `---
id: T-001
status: open
---
# Test
`
	path := filepath.Join(tmpDir, "T-001.md")
	if err := os.WriteFile(path, []byte(ticketContent), 0644); err != nil {
		t.Fatal(err)
	}

	tm := NewTicketManager("tk", tmpDir)
	if err := tm.AddNote("T-001", "added a note"); err != nil {
		t.Fatalf("AddNote() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), "added a note") {
		t.Errorf("AddNote did not append note, content:\n%s", string(data))
	}
}

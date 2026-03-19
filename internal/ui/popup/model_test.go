package popup

import (
	"testing"

	"github.com/alexcabrera/choo-choo/internal/ticket"
)

func TestTabString(t *testing.T) {
	tests := []struct {
		tab      Tab
		expected string
	}{
		{TabDetails, "Details"},
		{TabLog, "Log"},
		{TabDiff, "Diff"},
		{Tab(99), "Unknown"},
	}
	for _, tt := range tests {
		if got := tt.tab.String(); got != tt.expected {
			t.Errorf("Tab(%d).String() = %q, want %q", tt.tab, got, tt.expected)
		}
	}
}

func TestTabContentFields(t *testing.T) {
	tc := TabContent{
		Title:   "Test",
		Content: "Some content",
	}
	if tc.Title != "Test" {
		t.Errorf("TabContent.Title = %q, want %q", tc.Title, "Test")
	}
	if tc.Content != "Some content" {
		t.Errorf("TabContent.Content = %q, want %q", tc.Content, "Some content")
	}
}

func TestNewPopupModel(t *testing.T) {
	m := NewPopupModel()
	if m == nil {
		t.Fatal("NewPopupModel() returned nil")
	}
	if m.open {
		t.Error("NewPopupModel().open should be false")
	}
	if m.ticket != nil {
		t.Error("NewPopupModel().ticket should be nil")
	}
	if m.activeTab != TabDetails {
		t.Errorf("NewPopupModel().activeTab = %v, want %v", m.activeTab, TabDetails)
	}
	if len(m.tabs) != 3 {
		t.Errorf("NewPopupModel().tabs should have 3 elements, got %d", len(m.tabs))
	}
}

func TestPopupModelOpen(t *testing.T) {
	m := NewPopupModel()
	tk := &ticket.Ticket{
		ID:     "T-001",
		Title:  "Test Ticket",
		Status: ticket.StatusOpen,
	}
	m.Open(tk)

	if !m.open {
		t.Error("Open() should set open to true")
	}
	if m.ticket != tk {
		t.Error("Open() should set ticket")
	}
	if m.activeTab != TabDetails {
		t.Errorf("Open() activeTab = %v, want %v", m.activeTab, TabDetails)
	}
	if m.tabs[0].Content == "" {
		t.Error("Open() should populate Details tab content")
	}
}

func TestPopupModelClose(t *testing.T) {
	m := NewPopupModel()
	tk := &ticket.Ticket{ID: "T-001"}
	m.Open(tk)
	m.Close()

	if m.open {
		t.Error("Close() should set open to false")
	}
	if m.ticket != nil {
		t.Error("Close() should clear ticket")
	}
	if m.activeTab != TabDetails {
		t.Errorf("Close() activeTab = %v, want %v", m.activeTab, TabDetails)
	}
}

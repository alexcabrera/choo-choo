package popup

import "testing"

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

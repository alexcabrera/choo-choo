package kanban

import "testing"

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

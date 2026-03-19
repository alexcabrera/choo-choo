package state

import "testing"

func TestNewState(t *testing.T) {
	s := NewState()
	if s.Phase != PhaseInit {
		t.Errorf("NewState().Phase = %v, want %v", s.Phase, PhaseInit)
	}
	if s.Epic != "" {
		t.Errorf("NewState().Epic should be empty, got %q", s.Epic)
	}
	if s.Focus != "" {
		t.Errorf("NewState().Focus should be empty, got %q", s.Focus)
	}
	if s.Learnings == nil {
		t.Error("NewState().Learnings should not be nil")
	}
	if s.Decisions == nil {
		t.Error("NewState().Decisions should not be nil")
	}
}

func TestStateYAMLTags(t *testing.T) {
	s := State{
		Phase:     PhaseDesign,
		Epic:      "E-001",
		Focus:     "T-003",
		Learnings: []string{"learning1"},
		Decisions: []string{"decision1"},
	}
	if s.Phase != PhaseDesign {
		t.Errorf("State.Phase = %v, want %v", s.Phase, PhaseDesign)
	}
}

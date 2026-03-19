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

func TestStateSetters(t *testing.T) {
	s := NewState()

	s.SetPhase(PhaseExecution)
	if s.Phase != PhaseExecution {
		t.Errorf("SetPhase: Phase = %v, want %v", s.Phase, PhaseExecution)
	}

	s.SetFocus("T-123")
	if s.Focus != "T-123" {
		t.Errorf("SetFocus: Focus = %q, want %q", s.Focus, "T-123")
	}

	s.AddLearning("learned something")
	if len(s.Learnings) != 1 || s.Learnings[0] != "learned something" {
		t.Errorf("AddLearning: Learnings = %v, want %v", s.Learnings, []string{"learned something"})
	}

	s.AddDecision("decided something")
	if len(s.Decisions) != 1 || s.Decisions[0] != "decided something" {
		t.Errorf("AddDecision: Decisions = %v, want %v", s.Decisions, []string{"decided something"})
	}
}

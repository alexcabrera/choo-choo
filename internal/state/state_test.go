package state

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestLoadStateValid(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "state.md")
	content := `---
phase: execution
epic: E-001
focus: T-123
started: 2026-03-18T12:00:00Z
learnings: []
decisions: []
---
# Session State
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	s, err := LoadState(path)
	if err != nil {
		t.Fatalf("LoadState() error = %v", err)
	}
	if s.Phase != PhaseExecution {
		t.Errorf("Phase = %v, want %v", s.Phase, PhaseExecution)
	}
	if s.Epic != "E-001" {
		t.Errorf("Epic = %q, want %q", s.Epic, "E-001")
	}
	if s.Focus != "T-123" {
		t.Errorf("Focus = %q, want %q", s.Focus, "T-123")
	}
}

func TestLoadStateMissingFile(t *testing.T) {
	_, err := LoadState("/nonexistent/path/state.md")
	if err == nil {
		t.Error("LoadState() should return error for missing file")
	}
}

func TestLoadStateInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "state.md")
	content := `not yaml frontmatter at all
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadState(path)
	if err == nil {
		t.Error("LoadState() should return error for invalid YAML")
	}
}

func TestLoadStateInvalidPhase(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "state.md")
	content := `---
phase: invalid-phase
epic: ""
focus: ""
started: 2026-03-18T12:00:00Z
learnings: []
decisions: []
---
`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := LoadState(path)
	if err == nil {
		t.Error("LoadState() should return error for invalid phase")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "state.md")

	s := NewState()
	s.SetPhase(PhaseDesign)
	s.SetFocus("T-456")
	s.AddLearning("test learning")
	s.AddDecision("test decision")

	if err := s.Save(path); err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	loaded, err := LoadState(path)
	if err != nil {
		t.Fatalf("LoadState() error = %v", err)
	}
	if loaded.Phase != PhaseDesign {
		t.Errorf("loaded Phase = %v, want %v", loaded.Phase, PhaseDesign)
	}
	if loaded.Focus != "T-456" {
		t.Errorf("loaded Focus = %q, want %q", loaded.Focus, "T-456")
	}
	if len(loaded.Learnings) != 1 || loaded.Learnings[0] != "test learning" {
		t.Errorf("loaded Learnings = %v", loaded.Learnings)
	}
	if len(loaded.Decisions) != 1 || loaded.Decisions[0] != "test decision" {
		t.Errorf("loaded Decisions = %v", loaded.Decisions)
	}
}

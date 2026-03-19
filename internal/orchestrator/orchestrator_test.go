package orchestrator

import (
	"testing"

	"github.com/alexcabrera/choo-choo/internal/state"
)

func TestNewOrchestrator(t *testing.T) {
	o := NewOrchestrator("/workspace/myproject")
	if o == nil {
		t.Fatal("NewOrchestrator() returned nil")
	}
	if o.projectDir != "/workspace/myproject" {
		t.Errorf("projectDir = %q, want %q", o.projectDir, "/workspace/myproject")
	}
	if o.crushPath != "crush" {
		t.Errorf("crushPath = %q, want %q", o.crushPath, "crush")
	}
	if o.tkPath != "ticket" {
		t.Errorf("tkPath = %q, want %q", o.tkPath, "ticket")
	}
}

func TestGetPhaseNilState(t *testing.T) {
	o := NewOrchestrator("/workspace")
	phase := o.GetPhase()
	if phase != state.PhaseInit {
		t.Errorf("GetPhase() = %v, want %v", phase, state.PhaseInit)
	}
}

func TestGetPhaseWithState(t *testing.T) {
	o := NewOrchestrator("/workspace")
	o.state = state.NewState()
	o.state.SetPhase(state.PhaseExecution)
	phase := o.GetPhase()
	if phase != state.PhaseExecution {
		t.Errorf("GetPhase() = %v, want %v", phase, state.PhaseExecution)
	}
}

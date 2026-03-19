package orchestrator

import (
	"github.com/alexcabrera/choo-choo/internal/state"
)

type Orchestrator struct {
	projectDir string
	state      *state.State
	crushPath  string
	tkPath     string
}

func NewOrchestrator(projectDir string) *Orchestrator {
	return &Orchestrator{
		projectDir: projectDir,
		crushPath:  "crush",
		tkPath:     "ticket",
	}
}

func (o *Orchestrator) GetPhase() state.Phase {
	if o.state == nil {
		return state.PhaseInit
	}
	return o.state.Phase
}

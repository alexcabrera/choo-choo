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

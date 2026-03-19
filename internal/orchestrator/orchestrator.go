package orchestrator

import (
	"github.com/alexcabrera/choo-choo/internal/state"
	"github.com/alexcabrera/choo-choo/internal/ticket"
)

type Orchestrator struct {
	projectDir     string
	state          *state.State
	crushPath      string
	tkPath         string
	ticketManager  *ticket.TicketManager
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

func (o *Orchestrator) GetTickets() ([]ticket.Ticket, error) {
	if o.ticketManager == nil {
		ticketsDir := o.projectDir + "/.tickets"
		o.ticketManager = ticket.NewTicketManager(o.tkPath, ticketsDir)
	}
	return o.ticketManager.List()
}

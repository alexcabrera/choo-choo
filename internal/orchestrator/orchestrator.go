package orchestrator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alexcabrera/choo-choo/internal/state"
	"github.com/alexcabrera/choo-choo/internal/ticket"
)

type Orchestrator struct {
	projectDir    string
	state         *state.State
	crushPath     string
	tkPath        string
	ticketManager *ticket.TicketManager
}

func NewOrchestrator(projectDir string) *Orchestrator {
	return &Orchestrator{
		projectDir: projectDir,
		crushPath:  "crush",
		tkPath:     "ticket",
	}
}

func (o *Orchestrator) Init() error {
	if err := o.verifyTools(); err != nil {
		return err
	}

	if err := o.createDirectories(); err != nil {
		return err
	}

	if err := o.initState(); err != nil {
		return err
	}

	return nil
}

func (o *Orchestrator) verifyTools() error {
	if _, err := exec.LookPath(o.crushPath); err != nil {
		return fmt.Errorf("crush not found in PATH: %w", err)
	}
	if _, err := exec.LookPath(o.tkPath); err != nil {
		return fmt.Errorf("tk not found in PATH: %w", err)
	}
	return nil
}

func (o *Orchestrator) createDirectories() error {
	dirs := []string{
		filepath.Join(o.projectDir, ".tickets"),
		filepath.Join(o.projectDir, "specs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (o *Orchestrator) initState() error {
	statePath := filepath.Join(o.projectDir, "STATE.md")

	if _, err := os.Stat(statePath); err == nil {
		loadedState, err := state.LoadState(statePath)
		if err != nil {
			return fmt.Errorf("failed to load existing state: %w", err)
		}
		o.state = loadedState
		return nil
	}

	o.state = state.NewState()
	if err := o.state.Save(statePath); err != nil {
		return fmt.Errorf("failed to save initial state: %w", err)
	}

	return nil
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

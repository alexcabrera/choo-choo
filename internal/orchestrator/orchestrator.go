package orchestrator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alexcabrera/choo-choo/internal/crush"
	"github.com/alexcabrera/choo-choo/internal/state"
	"github.com/alexcabrera/choo-choo/internal/ticket"
)

type Orchestrator struct {
	projectDir    string
	state         *state.State
	crushPath     string
	tkPath        string
	ticketManager *ticket.TicketManager
	crushRunner   *crush.CrushRunner
	sequential    bool
}

func NewOrchestrator(projectDir string) *Orchestrator {
	return &Orchestrator{
		projectDir: projectDir,
		crushPath:  "crush",
		tkPath:     "ticket",
	}
}

func (o *Orchestrator) SetSequential(sequential bool) {
	o.sequential = sequential
}

func (o *Orchestrator) Init() error {
	if err := o.createDirectories(); err != nil {
		return err
	}

	if err := o.initState(); err != nil {
		return err
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

func (o *Orchestrator) SetPhase(phase state.Phase) {
	if o.state != nil {
		o.state.SetPhase(phase)
	}
}

func (o *Orchestrator) PersistState() error {
	if o.state == nil {
		return nil
	}
	statePath := filepath.Join(o.projectDir, "STATE.md")
	return o.state.Save(statePath)
}

func (o *Orchestrator) GetTickets() ([]ticket.Ticket, error) {
	if o.ticketManager == nil {
		ticketsDir := o.projectDir + "/.tickets"
		o.ticketManager = ticket.NewTicketManager(o.tkPath, ticketsDir)
	}
	return o.ticketManager.List()
}

func (o *Orchestrator) ReadArtifact(relativePath string) (string, error) {
	fullPath := filepath.Join(o.projectDir, relativePath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to read artifact %s: %w", relativePath, err)
	}

	return string(data), nil
}

func (o *Orchestrator) RunDesign(ctx context.Context) (<-chan crush.StreamEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	if o.state == nil {
		o.state = state.NewState()
	}
	o.state.SetPhase(state.PhaseDesign)

	prompt := "use the design skill to create a specification for this project"

	events, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Quiet: false})
	if err != nil {
		return nil, fmt.Errorf("failed to start design phase: %w", err)
	}

	return events, nil
}

func (o *Orchestrator) RunPlan(ctx context.Context) (<-chan crush.StreamEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	if o.state == nil {
		o.state = state.NewState()
	}
	o.state.SetPhase(state.PhasePlan)

	prompt := "use the plan skill to create tickets from the specification"

	events, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Quiet: false})
	if err != nil {
		return nil, fmt.Errorf("failed to start plan phase: %w", err)
	}

	return events, nil
}

type ValidationResult struct {
	IsValid      bool
	ExecutionOrder []string
	Errors       []string
	Warnings     []string
}

func (o *Orchestrator) RunValidate(ctx context.Context) (<-chan crush.StreamEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	prompt := "use the validate skill to check the ticket plan"

	events, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Quiet: false})
	if err != nil {
		return nil, fmt.Errorf("failed to run validation: %w", err)
	}

	return events, nil
}

func (o *Orchestrator) GetValidationResult() (*ValidationResult, error) {
	tickets, err := o.GetTickets()
	if err != nil {
		return nil, err
	}

	result := &ValidationResult{
		IsValid:        true,
		ExecutionOrder: []string{},
		Errors:         []string{},
		Warnings:       []string{},
	}

	order, err := o.ticketManager.GetExecutionOrder()
	if err != nil {
		result.Warnings = append(result.Warnings, err.Error())
	} else {
		for _, level := range order {
			result.ExecutionOrder = append(result.ExecutionOrder, level...)
		}
	}

	if len(tickets) == 0 {
		result.Warnings = append(result.Warnings, "no tickets found")
	}

	return result, nil
}

type TicketEvent struct {
	TicketID string
	Status   string
	Error    error
}

func (o *Orchestrator) RunExecute(ctx context.Context) (<-chan TicketEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	if o.state == nil {
		o.state = state.NewState()
	}
	o.state.SetPhase(state.PhaseExecution)

	if o.ticketManager == nil {
		ticketsDir := o.projectDir + "/.tickets"
		o.ticketManager = ticket.NewTicketManager(o.tkPath, ticketsDir)
	}

	order, err := o.ticketManager.GetExecutionOrder()
	if err != nil {
		return nil, fmt.Errorf("failed to get execution order: %w", err)
	}

	events := make(chan TicketEvent, 100)

	go func() {
		defer close(events)

		for _, level := range order {
			if o.sequential {
				o.executeLevelSequential(ctx, level, events)
			} else {
				o.executeLevelParallel(ctx, level, events)
			}
		}
	}()

	return events, nil
}

func (o *Orchestrator) executeLevelSequential(ctx context.Context, level []string, events chan<- TicketEvent) {
	for _, ticketID := range level {
		select {
		case <-ctx.Done():
			return
		default:
			events <- TicketEvent{TicketID: ticketID, Status: "executing"}
			prompt := fmt.Sprintf("use the execute skill to implement ticket %s", ticketID)
			_, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Yolo: true})
			if err != nil {
				events <- TicketEvent{TicketID: ticketID, Status: "failed", Error: err}
			} else {
				events <- TicketEvent{TicketID: ticketID, Status: "completed"}
			}
		}
	}
}

func (o *Orchestrator) executeLevelParallel(ctx context.Context, level []string, events chan<- TicketEvent) {
	type result struct {
		ticketID string
		err      error
	}

	results := make(chan result, len(level))

	for _, ticketID := range level {
		go func(id string) {
			events <- TicketEvent{TicketID: id, Status: "executing"}
			prompt := fmt.Sprintf("use the execute skill to implement ticket %s", id)
			_, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Yolo: true})
			results <- result{ticketID: id, err: err}
		}(ticketID)
	}

	for range level {
		select {
		case <-ctx.Done():
			return
		case r := <-results:
			if r.err != nil {
				events <- TicketEvent{TicketID: r.ticketID, Status: "failed", Error: r.err}
			} else {
				events <- TicketEvent{TicketID: r.ticketID, Status: "completed"}
			}
		}
	}
}

func (o *Orchestrator) RunVerify(ctx context.Context) (<-chan crush.StreamEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	if o.state == nil {
		o.state = state.NewState()
	}
	o.state.SetPhase(state.PhaseVerify)

	prompt := "use the verify skill to verify all completed tickets"

	events, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Quiet: false})
	if err != nil {
		return nil, fmt.Errorf("failed to start verify phase: %w", err)
	}

	return events, nil
}

func (o *Orchestrator) RunCloseGaps(ctx context.Context) (<-chan crush.StreamEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	if o.state == nil {
		o.state = state.NewState()
	}
	o.state.SetPhase(state.PhaseCloseGaps)

	prompt := "use the close-gaps skill to fix any gaps found during verification"

	events, err := o.crushRunner.Run(ctx, prompt, crush.RunOptions{Quiet: false})
	if err != nil {
		return nil, fmt.Errorf("failed to start close-gaps phase: %w", err)
	}

	return events, nil
}

func (o *Orchestrator) SendMessage(ctx context.Context, message string) (<-chan crush.StreamEvent, error) {
	if o.crushRunner == nil {
		o.crushRunner = crush.NewRunner(o.crushPath, o.projectDir)
	}

	if o.crushRunner.GetSessionID() == "" {
		if err := o.crushRunner.CaptureLastSessionID(); err != nil {
			return nil, fmt.Errorf("no session to continue: %w", err)
		}
	}

	events, err := o.crushRunner.RunWithSession(ctx, message, crush.RunOptions{Quiet: false})
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return events, nil
}

func (o *Orchestrator) CaptureSessionID() error {
	if o.crushRunner == nil {
		return fmt.Errorf("no crush runner available")
	}
	return o.crushRunner.CaptureLastSessionID()
}

func (o *Orchestrator) GetSessionID() string {
	if o.crushRunner == nil {
		return ""
	}
	return o.crushRunner.GetSessionID()
}

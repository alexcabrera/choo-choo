package crush

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

type EventType int

const (
	EventTypeStdout EventType = iota
	EventTypeStderr
	EventTypeError
	EventTypeDone
)

func (e EventType) String() string {
	switch e {
	case EventTypeStdout:
		return "stdout"
	case EventTypeStderr:
		return "stderr"
	case EventTypeError:
		return "error"
	case EventTypeDone:
		return "done"
	default:
		return "unknown"
	}
}

type StreamEvent struct {
	Type    EventType
	Content string
	Time    time.Time
}

type RunOptions struct {
	Quiet bool
	Yolo  bool
	Model string
}

type CrushRunner struct {
	crushPath string
	workDir   string
	sessionID string
}

func NewRunner(crushPath, workDir string) *CrushRunner {
	return &CrushRunner{
		crushPath: crushPath,
		workDir:   workDir,
	}
}

func (r *CrushRunner) Run(ctx context.Context, prompt string, opts RunOptions) (<-chan StreamEvent, error) {
	args := []string{"run", prompt}
	if opts.Quiet {
		args = append(args, "--quiet")
	}
	if opts.Yolo {
		args = append(args, "--yolo")
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, r.crushPath, args...)
	cmd.Dir = r.workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start crush: %w", err)
	}

	events := make(chan StreamEvent, 100)
	var wg sync.WaitGroup

	wg.Add(2)
	go r.readStream(&wg, events, stdout, EventTypeStdout)
	go r.readStream(&wg, events, stderr, EventTypeStderr)

	go func() {
		wg.Wait()
		err := cmd.Wait()
		if err != nil {
			events <- StreamEvent{Type: EventTypeError, Content: err.Error(), Time: time.Now()}
		}
		events <- StreamEvent{Type: EventTypeDone, Time: time.Now()}
		close(events)
	}()

	return events, nil
}

func (r *CrushRunner) readStream(wg *sync.WaitGroup, events chan<- StreamEvent, reader io.Reader, eventType EventType) {
	defer wg.Done()
	buf := make([]byte, 4096)
	for {
		n, err := reader.Read(buf)
		if n > 0 {
			events <- StreamEvent{
				Type:    eventType,
				Content: string(buf[:n]),
				Time:    time.Now(),
			}
		}
		if err != nil {
			if err != io.EOF {
				events <- StreamEvent{
					Type:    EventTypeError,
					Content: err.Error(),
					Time:    time.Now(),
				}
			}
			return
		}
	}
}

func (r *CrushRunner) SetSessionID(id string) {
	r.sessionID = id
}

func (r *CrushRunner) GetSessionID() string {
	return r.sessionID
}

func (r *CrushRunner) RunWithSession(ctx context.Context, opts RunOptions) (<-chan StreamEvent, error) {
	if r.sessionID == "" {
		return nil, fmt.Errorf("no session ID set")
	}

	args := []string{"run", "--continue", r.sessionID}
	if opts.Quiet {
		args = append(args, "--quiet")
	}
	if opts.Yolo {
		args = append(args, "--yolo")
	}
	if opts.Model != "" {
		args = append(args, "--model", opts.Model)
	}

	cmd := exec.CommandContext(ctx, r.crushPath, args...)
	cmd.Dir = r.workDir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start crush: %w", err)
	}

	events := make(chan StreamEvent, 100)
	var wg sync.WaitGroup

	wg.Add(2)
	go r.readStream(&wg, events, stdout, EventTypeStdout)
	go r.readStream(&wg, events, stderr, EventTypeStderr)

	go func() {
		wg.Wait()
		err := cmd.Wait()
		if err != nil {
			events <- StreamEvent{Type: EventTypeError, Content: err.Error(), Time: time.Now()}
		}
		events <- StreamEvent{Type: EventTypeDone, Time: time.Now()}
		close(events)
	}()

	return events, nil
}

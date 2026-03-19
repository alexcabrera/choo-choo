package crush

import (
	"context"
	"testing"
	"time"
)

func TestEventTypeValues(t *testing.T) {
	tests := []struct {
		eventType EventType
		expected  int
	}{
		{EventTypeStdout, 0},
		{EventTypeStderr, 1},
		{EventTypeError, 2},
		{EventTypeDone, 3},
	}

	for _, tt := range tests {
		t.Run(tt.eventType.String(), func(t *testing.T) {
			if int(tt.eventType) != tt.expected {
				t.Errorf("EventType value = %d, want %d", tt.eventType, tt.expected)
			}
		})
	}
}

func TestStreamEventFields(t *testing.T) {
	evt := StreamEvent{
		Type:    EventTypeStdout,
		Content: "test output",
	}
	if evt.Type != EventTypeStdout {
		t.Errorf("StreamEvent.Type = %v, want %v", evt.Type, EventTypeStdout)
	}
	if evt.Content != "test output" {
		t.Errorf("StreamEvent.Content = %q, want %q", evt.Content, "test output")
	}
}

func TestRunOptionsDefaults(t *testing.T) {
	opts := RunOptions{}
	if opts.Quiet != false {
		t.Error("RunOptions.Quiet should default to false")
	}
	if opts.Yolo != false {
		t.Error("RunOptions.Yolo should default to false")
	}
	if opts.Model != "" {
		t.Error("RunOptions.Model should default to empty")
	}
}

func TestNewRunner(t *testing.T) {
	runner := NewRunner("/usr/local/bin/crush", "/workspace")
	if runner == nil {
		t.Fatal("NewRunner() returned nil")
	}
	if runner.crushPath != "/usr/local/bin/crush" {
		t.Errorf("crushPath = %q, want %q", runner.crushPath, "/usr/local/bin/crush")
	}
	if runner.workDir != "/workspace" {
		t.Errorf("workDir = %q, want %q", runner.workDir, "/workspace")
	}
	if runner.sessionID != "" {
		t.Errorf("sessionID should be empty, got %q", runner.sessionID)
	}
}

func TestSessionIDOperations(t *testing.T) {
	runner := NewRunner("crush", "/workspace")

	runner.SetSessionID("test-session-123")
	if runner.GetSessionID() != "test-session-123" {
		t.Errorf("GetSessionID() = %q, want %q", runner.GetSessionID(), "test-session-123")
	}
}

func TestRunWithSessionNoID(t *testing.T) {
	runner := NewRunner("crush", "/workspace")
	ctx := context.Background()

	_, err := runner.RunWithSession(ctx, RunOptions{})
	if err == nil {
		t.Error("RunWithSession() should return error when sessionID is empty")
	}
}

func TestRunWithContextCancellation(t *testing.T) {
	runner := NewRunner("echo", "/tmp")
	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	_, err := runner.Run(ctx, "test", RunOptions{})
	if err != nil {
		t.Logf("Run with cancelled context returned: %v", err)
	}
}

func TestProcessLifecycle(t *testing.T) {
	runner := NewRunner("echo", "/tmp")
	ctx := context.Background()

	events, err := runner.Run(ctx, "hello", RunOptions{})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	select {
	case evt := <-events:
		if evt.Type == EventTypeError {
			t.Errorf("Received error event: %s", evt.Content)
		}
	case <-time.After(5 * time.Second):
		t.Error("Timeout waiting for event")
	}
}

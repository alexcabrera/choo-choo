package crush

import "testing"

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

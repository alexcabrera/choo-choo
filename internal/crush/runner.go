package crush

import "time"

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

package ticket

import "time"

type TicketType string

const (
	TypeEpic   TicketType = "epic"
	TypeStory  TicketType = "story"
	TypeTask   TicketType = "task"
	TypeChore  TicketType = "chore"
	TypeBug    TicketType = "bug"
	TypeFeature TicketType = "feature"
)

type Status string

const (
	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusClosed     Status = "closed"
)

type Ticket struct {
	ID           string
	Type         TicketType
	Title        string
	Description  string
	Status       Status
	Parent       string
	Dependencies []string
	Accepts      []string
	Notes        []string
	Created      time.Time
	Updated      time.Time
}

package model

type PhaseChangeMsg struct {
	Phase string
}

type TicketUpdateMsg struct {
	TicketID string
	Status   string
}

type StreamChunkMsg struct {
	Content string
	Done    bool
}

type ErrorMsg struct {
	Err error
}

type PopupOpenMsg struct {
	TicketID string
}

type PopupCloseMsg struct{}

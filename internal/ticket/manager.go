package ticket

type TicketManager struct {
	tkPath     string
	ticketsDir string
}

func NewTicketManager(tkPath, ticketsDir string) *TicketManager {
	return &TicketManager{
		tkPath:     tkPath,
		ticketsDir: ticketsDir,
	}
}

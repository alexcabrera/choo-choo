package model

import "github.com/charmbracelet/lipgloss"

var (
	ColorPrimary   = lipgloss.Color("#7D56F4")
	ColorSecondary = lipgloss.Color("#4A4A4A")
	ColorError     = lipgloss.Color("#FF6B6B")
	ColorSuccess   = lipgloss.Color("#4ECDC4")

	BorderStyle = lipgloss.RoundedBorder()

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 1)

	BodyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA"))

	EmphasisStyle = lipgloss.NewStyle().
			Italic(true).
			Foreground(ColorSecondary)

	ChatBoxStyle = lipgloss.NewStyle().
			Border(BorderStyle).
			BorderForeground(ColorSecondary).
			Padding(1, 2)

	KanbanColumnStyle = lipgloss.NewStyle().
				Border(BorderStyle).
				BorderForeground(ColorSecondary).
				Padding(0, 1).
				Margin(0, 1)

	TicketCardStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(ColorSecondary).
				Padding(0, 1).
				Margin(1, 0)

	PopupStyle = lipgloss.NewStyle().
			Border(BorderStyle).
			BorderForeground(ColorPrimary).
			Padding(1, 2)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			Padding(0, 2)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(ColorSecondary).
				Padding(0, 2)
)

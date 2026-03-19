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

func HorizontalJoin(width int, elements ...string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top, elements...)
}

func VerticalJoin(height int, elements ...string) string {
	return lipgloss.JoinVertical(lipgloss.Left, elements...)
}

func WithBorder(style lipgloss.Style, content string, width, height int) string {
	return style.Width(width).Height(height).Render(content)
}

func WithMargin(style lipgloss.Style, content string) string {
	return style.Render(content)
}

func WithPadding(style lipgloss.Style, content string) string {
	return style.Render(content)
}

func LayoutColumns(totalWidth int, columns ...string) string {
	if len(columns) == 0 {
		return ""
	}
	colWidth := totalWidth / len(columns)
	styledCols := make([]string, len(columns))
	for i, col := range columns {
		styledCols[i] = lipgloss.NewStyle().Width(colWidth).Render(col)
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, styledCols...)
}

func LayoutRows(totalHeight int, rows ...string) string {
	if len(rows) == 0 {
		return ""
	}
	rowHeight := totalHeight / len(rows)
	styledRows := make([]string, len(rows))
	for i, row := range rows {
		styledRows[i] = lipgloss.NewStyle().Height(rowHeight).Render(row)
	}
	return lipgloss.JoinVertical(lipgloss.Left, styledRows...)
}

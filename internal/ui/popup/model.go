package popup

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"

	"github.com/alexcabrera/choo-choo/internal/ticket"
)

type Tab int

const (
	TabDetails Tab = iota
	TabLog
	TabDiff
)

func (t Tab) String() string {
	switch t {
	case TabDetails:
		return "Details"
	case TabLog:
		return "Log"
	case TabDiff:
		return "Diff"
	default:
		return "Unknown"
	}
}

var (
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#6B7280")

	popupBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(primaryColor).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(secondaryColor).
				Padding(0, 2)
)

type TabContent struct {
	Title   string
	Content string
}

type PopupModel struct {
	open      bool
	ticket    *ticket.Ticket
	activeTab Tab
	tabs      [3]TabContent
	scrollY   int
}

func NewPopupModel() *PopupModel {
	return &PopupModel{
		open:      false,
		activeTab: TabDetails,
		tabs: [3]TabContent{
			{Title: "Details"},
			{Title: "Log"},
			{Title: "Diff"},
		},
		scrollY: 0,
	}
}

func (m *PopupModel) Open(t *ticket.Ticket) {
	m.ticket = t
	m.activeTab = TabDetails
	m.open = true
	m.scrollY = 0
	m.tabs[0].Content = formatTicketDetails(t)
	m.tabs[1].Content = "No log available yet."
	m.tabs[2].Content = "No diff available yet."
}

func (m *PopupModel) Close() {
	m.open = false
	m.ticket = nil
	m.activeTab = TabDetails
	m.scrollY = 0
	for i := range m.tabs {
		m.tabs[i].Content = ""
	}
}

func formatTicketDetails(t *ticket.Ticket) string {
	if t == nil {
		return "No ticket selected"
	}

	var b strings.Builder

	var typeEmoji string
	switch t.Type {
	case ticket.TypeEpic:
		typeEmoji = "📦"
	case ticket.TypeStory:
		typeEmoji = "📋"
	case ticket.TypeTask:
		typeEmoji = "✓"
	case ticket.TypeBug:
		typeEmoji = "🐛"
	case ticket.TypeFeature:
		typeEmoji = "✨"
	case ticket.TypeChore:
		typeEmoji = "🔧"
	default:
		typeEmoji = "•"
	}

	b.WriteString(fmt.Sprintf("ID: %s\n", t.ID))
	b.WriteString(fmt.Sprintf("Type: %s %s\n", typeEmoji, t.Type))
	b.WriteString(fmt.Sprintf("Status: %s\n", t.Status))

	if t.Parent != "" {
		b.WriteString(fmt.Sprintf("Parent: %s\n", t.Parent))
	}

	b.WriteString(fmt.Sprintf("Title: %s\n\n", t.Title))

	if t.Description != "" {
		b.WriteString("Description:\n")
		b.WriteString(t.Description)
		b.WriteString("\n\n")
	}

	if len(t.Dependencies) > 0 {
		b.WriteString(fmt.Sprintf("Dependencies: %s\n", strings.Join(t.Dependencies, ", ")))
	}

	if len(t.Accepts) > 0 {
		b.WriteString("\nAcceptance Criteria:\n")
		for _, a := range t.Accepts {
			b.WriteString(fmt.Sprintf("  • %s\n", a))
		}
	}

	return b.String()
}

func (m *PopupModel) SetLog(log string) {
	m.tabs[TabLog].Content = log
}

func (m *PopupModel) SetDiff(diff string) {
	m.tabs[TabDiff].Content = diff
}

func (m *PopupModel) NextTab() {
	m.activeTab = (m.activeTab + 1) % 3
	m.scrollY = 0
}

func (m *PopupModel) PrevTab() {
	m.activeTab = (m.activeTab + 2) % 3
	m.scrollY = 0
}

func (m *PopupModel) IsOpen() bool {
	return m.open
}

func (m *PopupModel) HandleScroll(direction int) {
	if direction < 0 && m.scrollY > 0 {
		m.scrollY--
	} else if direction > 0 {
		m.scrollY++
	}
}

func (m *PopupModel) GetActiveTab() Tab {
	return m.activeTab
}

func (m *PopupModel) GetContent() string {
	return m.tabs[m.activeTab].Content
}

func (m *PopupModel) Update(msg tea.Msg) (*PopupModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	}
	return m, nil
}

func (m *PopupModel) handleKeyPress(msg tea.KeyMsg) (*PopupModel, tea.Cmd) {
	switch msg.String() {
	case "tab":
		m.NextTab()
	case "shift+tab":
		m.PrevTab()
	case "up", "k":
		m.HandleScroll(-1)
	case "down", "j":
		m.HandleScroll(1)
	case "esc", "q":
		m.Close()
	}
	return m, nil
}

func (m *PopupModel) View(width, height int) tea.View {
	popupWidth := max(width*2/3, 40)
	popupHeight := max(height*2/3, 10)

	var b strings.Builder

	b.WriteString(m.renderTabBar(popupWidth))
	b.WriteString("\n")
	b.WriteString(lipgloss.NewStyle().Foreground(secondaryColor).Render(strings.Repeat("─", popupWidth-4)))
	b.WriteString("\n\n")

	content := m.GetContent()
	lines := strings.Split(content, "\n")

	maxLines := popupHeight - 6
	endIdx := min(m.scrollY+maxLines, len(lines))
	if m.scrollY > len(lines) {
		m.scrollY = 0
	}

	for i := m.scrollY; i < endIdx; i++ {
		line := lines[i]
		if len(line) > popupWidth-6 {
			line = line[:popupWidth-9] + "..."
		}
		b.WriteString(line + "\n")
	}

	if m.scrollY > 0 {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("↑ More above\n"))
	}
	if endIdx < len(lines) {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render("↓ More below\n"))
	}

	footer := "\n" + lipgloss.NewStyle().Foreground(secondaryColor).Render("tab: switch | ↑↓: scroll | esc/q: close")
	b.WriteString(footer)

	contentBox := popupBoxStyle.
		Width(popupWidth).
		Height(popupHeight).
		Render(b.String())

	return tea.NewView(contentBox)
}

func (m *PopupModel) renderTabBar(width int) string {
	tabs := make([]string, 3)
	for i, tab := range m.tabs {
		if Tab(i) == m.activeTab {
			tabs[i] = tabActiveStyle.Render(tab.Title)
		} else {
			tabs[i] = tabInactiveStyle.Render(tab.Title)
		}
	}

	separator := lipgloss.NewStyle().Foreground(secondaryColor).Render(" │ ")

	tabBar := lipgloss.NewStyle().
		Width(width - 4).
		Render(lipgloss.JoinHorizontal(lipgloss.Left, tabs[0], separator, tabs[1], separator, tabs[2]))

	return tabBar
}

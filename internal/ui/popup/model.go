package popup

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"

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

type TabContent struct {
	Title   string
	Content string
}

type PopupModel struct {
	open      bool
	ticket    *ticket.Ticket
	activeTab Tab
	tabs      [3]TabContent
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
	}
}

func (m *PopupModel) Open(t *ticket.Ticket) {
	m.ticket = t
	m.activeTab = TabDetails
	m.open = true
	m.tabs[0].Content = formatTicketDetails(t)
}

func (m *PopupModel) Close() {
	m.open = false
	m.ticket = nil
	m.activeTab = TabDetails
	for i := range m.tabs {
		m.tabs[i].Content = ""
	}
}

func formatTicketDetails(t *ticket.Ticket) string {
	if t == nil {
		return ""
	}
	return "ID: " + t.ID + "\nTitle: " + t.Title + "\nStatus: " + string(t.Status)
}

func (m *PopupModel) SetLog(log string) {
	m.tabs[TabLog].Content = log
}

func (m *PopupModel) SetDiff(diff string) {
	m.tabs[TabDiff].Content = diff
}

func (m *PopupModel) NextTab() {
	m.activeTab = (m.activeTab + 1) % 3
}

func (m *PopupModel) PrevTab() {
	m.activeTab = (m.activeTab + 2) % 3
}

func (m *PopupModel) IsOpen() bool {
	return m.open
}

func (m *PopupModel) HandleScroll(direction int) {
	// Scroll handling for content viewport
	// Positive direction = scroll down, negative = scroll up
	// Actual scroll offset managed by rendering layer
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
	case "escape", "q":
		m.Close()
	}
	return m, nil
}

func (m *PopupModel) View(width, height int) tea.View {
	var b strings.Builder

	popupWidth := width * 2 / 3
	if popupWidth < 40 {
		popupWidth = 40
	}
	popupHeight := height * 2 / 3
	if popupHeight < 10 {
		popupHeight = 10
	}

	tabBar := m.renderTabBar(popupWidth)
	b.WriteString(tabBar + "\n")
	b.WriteString(strings.Repeat("-", popupWidth) + "\n")

	content := m.GetContent()
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if i >= popupHeight-4 {
			break
		}
		if len(line) > popupWidth-2 {
			line = line[:popupWidth-5] + "..."
		}
		b.WriteString(line + "\n")
	}

	border := "+" + strings.Repeat("-", popupWidth-2) + "+"
	contentView := b.String()

	var fullView strings.Builder
	for i := 0; i < (height-popupHeight)/2; i++ {
		fullView.WriteString("\n")
	}
	fullView.WriteString(border + "\n")
	for _, line := range strings.Split(contentView, "\n") {
		padded := fmt.Sprintf("|%-*s|", popupWidth-2, line)
		fullView.WriteString(padded + "\n")
	}
	fullView.WriteString(border)

	return tea.NewView(fullView.String())
}

func (m *PopupModel) renderTabBar(width int) string {
	var tabs []string
	for i, tab := range m.tabs {
		style := "  %s  "
		if Tab(i) == m.activeTab {
			style = "[%s]"
		}
		tabs = append(tabs, fmt.Sprintf(style, tab.Title))
	}
	return strings.Join(tabs, " | ")
}

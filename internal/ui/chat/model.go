package chat

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type Role int

const (
	RoleUser Role = iota
	RoleAssistant
	RoleError
)

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAssistant:
		return "assistant"
	case RoleError:
		return "error"
	default:
		return "unknown"
	}
}

func (r Role) DisplayName() string {
	switch r {
	case RoleUser:
		return "You"
	case RoleAssistant:
		return "Crush"
	case RoleError:
		return "Error"
	default:
		return "Unknown"
	}
}

type ChatMessage struct {
	Role      Role
	Content   string
	Timestamp time.Time
}

type spinnerTickMsg struct{}

var (
	userStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

	assistantStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#10B981")).
			Bold(true)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EF4444")).
			Bold(true)

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7D56F4")).
				Bold(true)

	inputBoxStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#2D2D44")).
			Padding(0, 1)

	spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
)

type ChatModel struct {
	messages     []ChatMessage
	streaming    bool
	input        string
	viewportY    int
	spinnerFrame int
	width        int
	height       int
}

func NewChatModel() *ChatModel {
	return &ChatModel{
		messages:     []ChatMessage{},
		streaming:    false,
		input:        "",
		viewportY:    0,
		spinnerFrame: 0,
	}
}

func (m *ChatModel) AddMessage(role Role, content string) {
	msg := ChatMessage{
		Role:      role,
		Content:   content,
		Timestamp: time.Now().UTC(),
	}
	m.messages = append(m.messages, msg)
}

func (m *ChatModel) SetStreaming(streaming bool) {
	m.streaming = streaming
}

func (m *ChatModel) GetInput() string {
	return m.input
}

func (m *ChatModel) HandleTextInput(key string) {
	switch key {
	case "backspace":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	case "delete":
		if len(m.input) > 0 {
			m.input = m.input[:len(m.input)-1]
		}
	default:
		if len(key) == 1 {
			m.input += key
		}
	}
}

func (m *ChatModel) ClearInput() {
	m.input = ""
}

func (m *ChatModel) HandleViewportScroll(direction int) {
	if direction < 0 && m.viewportY > 0 {
		m.viewportY--
	} else if direction > 0 {
		m.viewportY++
	}
}

func (m *ChatModel) GetViewportY() int {
	return m.viewportY
}

func (m *ChatModel) HandleSpinnerTick() int {
	m.spinnerFrame = (m.spinnerFrame + 1) % len(spinnerFrames)
	return m.spinnerFrame
}

func (m *ChatModel) GetSpinnerFrame() int {
	return m.spinnerFrame
}

func (m *ChatModel) SendMessage() string {
	if m.input == "" {
		return ""
	}
	content := m.input
	m.AddMessage(RoleUser, content)
	m.input = ""
	return content
}

func (m *ChatModel) IsStreaming() bool {
	return m.streaming
}

func (m *ChatModel) AppendAssistantChunk(content string) {
	if len(m.messages) > 0 && m.messages[len(m.messages)-1].Role == RoleAssistant {
		m.messages[len(m.messages)-1].Content += content
	} else {
		m.AddMessage(RoleAssistant, content)
	}
}

func (m *ChatModel) AppendErrorChunk(content string) {
	if len(m.messages) > 0 && m.messages[len(m.messages)-1].Role == RoleError {
		m.messages[len(m.messages)-1].Content += content
	} else {
		m.AddMessage(RoleError, content)
	}
}

func (m *ChatModel) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *ChatModel) Update(msg tea.Msg) (*ChatModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)
	case spinnerTickMsg:
		m.HandleSpinnerTick()
		if m.streaming {
			return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
				return spinnerTickMsg{}
			})
		}
	}
	return m, nil
}

func (m *ChatModel) handleKeyPress(msg tea.KeyMsg) (*ChatModel, tea.Cmd) {
	switch msg.String() {
	case "up":
		m.HandleViewportScroll(-1)
	case "down":
		m.HandleViewportScroll(1)
	case "enter":
		m.SendMessage()
	case "backspace":
		m.HandleTextInput("backspace")
	default:
		m.HandleTextInput(msg.String())
	}
	return m, nil
}

func (m *ChatModel) View() tea.View {
	var b strings.Builder

	b.WriteString("Chat\n")
	b.WriteString(strings.Repeat("─", m.width-4))
	b.WriteString("\n\n")

	messageHeight := max(m.height-6, 5)


	linesRendered := 0
	for _, msg := range m.messages {
		if linesRendered >= messageHeight {
			break
		}

		var rolePrefix string
		var content string

		switch msg.Role {
		case RoleUser:
			rolePrefix = userStyle.Render("▶ You:")
			content = msg.Content
		case RoleError:
			rolePrefix = errorStyle.Render("⚠ Error:")
			content = msg.Content
		default:
			rolePrefix = assistantStyle.Render("◆ Crush:")
			content = msg.Content
		}

		b.WriteString(rolePrefix + "\n")

		contentLines := strings.Split(content, "\n")
		for _, line := range contentLines {
			if linesRendered >= messageHeight {
				break
			}
			wrapped := m.wrapText(line, m.width-8)
			for _, wrappedLine := range strings.Split(wrapped, "\n") {
				if linesRendered >= messageHeight {
					break
				}
				b.WriteString("  " + wrappedLine + "\n")
				linesRendered++
			}
		}
		b.WriteString("\n")
		linesRendered++
	}

	if m.streaming {
		spinner := spinnerFrames[m.spinnerFrame]
		b.WriteString(assistantStyle.Render(spinner+" Crush is thinking...") + "\n")
	}

	b.WriteString("\n")
	b.WriteString(m.renderInput())

	return tea.NewView(b.String())
}

func (m *ChatModel) renderInput() string {
	prompt := inputPromptStyle.Render(">")
	input := m.input + "│"

	inputLine := fmt.Sprintf("%s %s", prompt, input)
	return inputBoxStyle.Width(m.width - 4).Render(inputLine)
}

func (m *ChatModel) wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	currentLen := 0

	for _, word := range words {
		wordLen := len(word)
		if currentLen+wordLen+1 > maxWidth {
			if currentLen > 0 {
				result.WriteString("\n")
			}
			result.WriteString(word)
			currentLen = wordLen
		} else {
			if currentLen > 0 {
				result.WriteString(" ")
				currentLen++
			}
			result.WriteString(word)
			currentLen += wordLen
		}
	}

	return result.String()
}

func (m *ChatModel) GetMessages() []ChatMessage {
	return m.messages
}

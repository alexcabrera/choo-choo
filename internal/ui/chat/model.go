package chat

import (
	"fmt"
	"strings"
	tea "charm.land/bubbletea/v2"
	"time"
)

type Role int

const (
	RoleUser Role = iota
	RoleAssistant
)

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAssistant:
		return "assistant"
	default:
		return "unknown"
	}
}

type ChatMessage struct {
	Role      Role
	Content   string
	Timestamp time.Time
}

type spinnerTickMsg struct{}

type ChatModel struct {
	messages     []ChatMessage
	streaming    bool
	input        string
	viewportY    int
	spinnerFrame int
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
	m.spinnerFrame = (m.spinnerFrame + 1) % 4
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

	b.WriteString("=== Chat ===\n\n")

	for _, msg := range m.messages {
		prefix := "[User]"
		if msg.Role == RoleAssistant {
			prefix = "[Assistant]"
		}
		b.WriteString(fmt.Sprintf("%s %s\n", prefix, msg.Content))
	}

	if m.streaming {
		spinner := []string{"|", "/", "-", "\\"}
		b.WriteString(fmt.Sprintf("\n[Assistant] %s Thinking...", spinner[m.spinnerFrame]))
	}

	b.WriteString("\n\n> " + m.input + "_")

	return tea.NewView(b.String())
}

func (m *ChatModel) GetMessages() []ChatMessage {
	return m.messages
}

package chat

import "time"

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

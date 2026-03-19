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
	messages  []ChatMessage
	streaming bool
}

func NewChatModel() *ChatModel {
	return &ChatModel{
		messages:  []ChatMessage{},
		streaming: false,
	}
}

package chat

import "testing"

func TestRoleString(t *testing.T) {
	tests := []struct {
		role     Role
		expected string
	}{
		{RoleUser, "user"},
		{RoleAssistant, "assistant"},
		{Role(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.role.String(); got != tt.expected {
			t.Errorf("Role(%d).String() = %q, want %q", tt.role, got, tt.expected)
		}
	}
}

func TestChatMessageFields(t *testing.T) {
	msg := ChatMessage{
		Role:    RoleUser,
		Content: "Hello",
	}
	if msg.Role != RoleUser {
		t.Errorf("ChatMessage.Role = %v, want %v", msg.Role, RoleUser)
	}
	if msg.Content != "Hello" {
		t.Errorf("ChatMessage.Content = %q, want %q", msg.Content, "Hello")
	}
}

func TestChatModelFields(t *testing.T) {
	m := ChatModel{}
	if m.messages != nil {
		t.Error("ChatModel.messages should be nil initially")
	}
	if m.streaming != false {
		t.Error("ChatModel.streaming should be false initially")
	}
}

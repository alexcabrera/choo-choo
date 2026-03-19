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

func TestNewChatModel(t *testing.T) {
	m := NewChatModel()
	if m == nil {
		t.Fatal("NewChatModel() returned nil")
	}
	if m.messages == nil {
		t.Error("NewChatModel().messages should not be nil")
	}
	if len(m.messages) != 0 {
		t.Errorf("NewChatModel().messages should be empty, got %d", len(m.messages))
	}
	if m.streaming {
		t.Error("NewChatModel().streaming should be false")
	}
}

func TestChatModelAddMessage(t *testing.T) {
	m := NewChatModel()
	m.AddMessage(RoleUser, "Hello")
	if len(m.messages) != 1 {
		t.Fatalf("AddMessage: len(messages) = %d, want 1", len(m.messages))
	}
	if m.messages[0].Role != RoleUser {
		t.Errorf("AddMessage: Role = %v, want %v", m.messages[0].Role, RoleUser)
	}
	if m.messages[0].Content != "Hello" {
		t.Errorf("AddMessage: Content = %q, want %q", m.messages[0].Content, "Hello")
	}
	if m.messages[0].Timestamp.IsZero() {
		t.Error("AddMessage: Timestamp should be set")
	}
}

func TestChatModelSetStreaming(t *testing.T) {
	m := NewChatModel()
	if m.streaming {
		t.Error("streaming should start as false")
	}
	m.SetStreaming(true)
	if !m.streaming {
		t.Error("SetStreaming(true): streaming should be true")
	}
	m.SetStreaming(false)
	if m.streaming {
		t.Error("SetStreaming(false): streaming should be false")
	}
}

package model

import "testing"

func TestErrorCodeValues(t *testing.T) {
	codes := []ErrorCode{
		ErrCrushFailed,
		ErrValidation,
		ErrFileCorrupted,
		ErrMaxRetries,
		ErrProcessCrash,
	}
	for _, code := range codes {
		if code == "" {
			t.Error("ErrorCode should not be empty")
		}
	}
}

func TestAppErrorError(t *testing.T) {
	err := AppError{
		Code:    ErrCrushFailed,
		Message: "crush process exited with code 1",
	}
	expected := "crush_failed: crush process exited with code 1"
	if got := err.Error(); got != expected {
		t.Errorf("AppError.Error() = %q, want %q", got, expected)
	}
}

func TestNewAppError(t *testing.T) {
	err := NewAppError(ErrValidation, "invalid input", "check your input")
	if err == nil {
		t.Fatal("NewAppError() returned nil")
	}
	if err.Code != ErrValidation {
		t.Errorf("Code = %v, want %v", err.Code, ErrValidation)
	}
	if err.Message != "invalid input" {
		t.Errorf("Message = %q, want %q", err.Message, "invalid input")
	}
	if err.Recovery != "check your input" {
		t.Errorf("Recovery = %q, want %q", err.Recovery, "check your input")
	}
}

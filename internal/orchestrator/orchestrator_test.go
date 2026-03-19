package orchestrator

import "testing"

func TestNewOrchestrator(t *testing.T) {
	o := NewOrchestrator("/workspace/myproject")
	if o == nil {
		t.Fatal("NewOrchestrator() returned nil")
	}
	if o.projectDir != "/workspace/myproject" {
		t.Errorf("projectDir = %q, want %q", o.projectDir, "/workspace/myproject")
	}
	if o.crushPath != "crush" {
		t.Errorf("crushPath = %q, want %q", o.crushPath, "crush")
	}
	if o.tkPath != "ticket" {
		t.Errorf("tkPath = %q, want %q", o.tkPath, "ticket")
	}
}

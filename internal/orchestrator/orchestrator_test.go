package orchestrator

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/alexcabrera/choo-choo/internal/state"
)

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

func TestGetPhaseNilState(t *testing.T) {
	o := NewOrchestrator("/workspace")
	phase := o.GetPhase()
	if phase != state.PhaseInit {
		t.Errorf("GetPhase() = %v, want %v", phase, state.PhaseInit)
	}
}

func TestGetPhaseWithState(t *testing.T) {
	o := NewOrchestrator("/workspace")
	o.state = state.NewState()
	o.state.SetPhase(state.PhaseExecution)
	phase := o.GetPhase()
	if phase != state.PhaseExecution {
		t.Errorf("GetPhase() = %v, want %v", phase, state.PhaseExecution)
	}
}

func TestInitCreatesDirectories(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "choo-choo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	o := NewOrchestrator(tmpDir)

	o.crushPath = "echo"
	o.tkPath = "echo"

	err = o.Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	ticketsDir := filepath.Join(tmpDir, ".tickets")
	if _, err := os.Stat(ticketsDir); os.IsNotExist(err) {
		t.Error(".tickets directory was not created")
	}

	specsDir := filepath.Join(tmpDir, "specs")
	if _, err := os.Stat(specsDir); os.IsNotExist(err) {
		t.Error("specs directory was not created")
	}

	stateFile := filepath.Join(tmpDir, "STATE.md")
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		t.Error("STATE.md was not created")
	}
}

func TestReadArtifact(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "choo-choo-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	content := "test content"
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	o := NewOrchestrator(tmpDir)
	result, err := o.ReadArtifact("test.txt")
	if err != nil {
		t.Fatalf("ReadArtifact() error = %v", err)
	}
	if result != content {
		t.Errorf("ReadArtifact() = %q, want %q", result, content)
	}
}

func TestReadArtifactNotFound(t *testing.T) {
	o := NewOrchestrator("/nonexistent")
	_, err := o.ReadArtifact("missing.txt")
	if err == nil {
		t.Error("ReadArtifact() should return error for missing file")
	}
}

func TestRunDesignRequiresCrush(t *testing.T) {
	o := NewOrchestrator("/workspace")
	o.crushPath = "nonexistent-crush-command"
	ctx := context.Background()

	_, err := o.RunDesign(ctx)
	if err == nil {
		t.Error("RunDesign() should return error when crush not found")
	}
}

func TestRunPlanRequiresCrush(t *testing.T) {
	o := NewOrchestrator("/workspace")
	o.crushPath = "nonexistent-crush-command"
	ctx := context.Background()

	_, err := o.RunPlan(ctx)
	if err == nil {
		t.Error("RunPlan() should return error when crush not found")
	}
}

func TestValidationResult(t *testing.T) {
	o := NewOrchestrator("/workspace")
	o.ticketManager = nil

	_, err := o.GetValidationResult()
	if err == nil {
		t.Error("GetValidationResult() should return error when ticketManager is nil")
	}
}

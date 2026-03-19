package state

import "testing"

func TestPhaseIsValid(t *testing.T) {
	tests := []struct {
		phase    Phase
		expected bool
	}{
		{PhaseInit, true},
		{PhaseDesign, true},
		{PhasePlan, true},
		{PhaseValidate, true},
		{PhaseExecution, true},
		{PhaseVerify, true},
		{PhaseCloseGaps, true},
		{PhaseDone, true},
		{Phase("invalid"), false},
		{Phase(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.phase), func(t *testing.T) {
			if got := tt.phase.IsValid(); got != tt.expected {
				t.Errorf("Phase(%q).IsValid() = %v, want %v", tt.phase, got, tt.expected)
			}
		})
	}
}

func TestPhaseString(t *testing.T) {
	tests := []struct {
		phase    Phase
		expected string
	}{
		{PhaseInit, "init"},
		{PhaseDesign, "design"},
		{PhasePlan, "plan"},
		{PhaseValidate, "validate"},
		{PhaseExecution, "execution"},
		{PhaseVerify, "verification"},
		{PhaseCloseGaps, "gap-closure"},
		{PhaseDone, "done"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.phase.String(); got != tt.expected {
				t.Errorf("Phase.String() = %q, want %q", got, tt.expected)
			}
		})
	}
}

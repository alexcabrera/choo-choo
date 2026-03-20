package state

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type State struct {
	Phase     Phase     `yaml:"phase"`
	Epic      string    `yaml:"epic"`
	Focus     string    `yaml:"focus"`
	Started   time.Time `yaml:"started"`
	Learnings []string  `yaml:"learnings"`
	Decisions []string  `yaml:"decisions"`
}

func NewState() *State {
	return &State{
		Phase:     PhaseInit,
		Started:   time.Now().UTC(),
		Learnings: []string{},
		Decisions: []string{},
	}
}

func LoadState(path string) (*State, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %w", err)
	}

	content := string(data)

	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("invalid state file: missing YAML frontmatter")
	}

	endIdx := strings.Index(content[4:], "\n---")
	if endIdx == -1 {
		return nil, fmt.Errorf("invalid state file: unterminated frontmatter")
	}

	frontmatter := content[4 : 4+endIdx]

	var state State
	if err := yaml.Unmarshal([]byte(frontmatter), &state); err != nil {
		return nil, fmt.Errorf("failed to parse state YAML: %w", err)
	}

	if !state.Phase.IsValid() {
		return nil, fmt.Errorf("invalid phase: %s", state.Phase)
	}

	return &state, nil
}

func (s *State) Save(path string) error {
	var sb strings.Builder
	sb.WriteString("---\n")

	encoder := yaml.NewEncoder(&sb)
	if err := encoder.Encode(map[string]any{
		"phase":     s.Phase,
		"epic":      s.Epic,
		"focus":     s.Focus,
		"started":   s.Started,
		"learnings": s.Learnings,
		"decisions": s.Decisions,
	}); err != nil {
		return fmt.Errorf("failed to encode state: %w", err)
	}

	sb.WriteString("---\n")
	sb.WriteString("\n# choo-choo Session State\n")

	content := sb.String()

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write state file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to rename state file: %w", err)
	}

	return nil
}

func (s *State) SetPhase(p Phase) {
	s.Phase = p
}

func (s *State) SetFocus(ticketID string) {
	s.Focus = ticketID
}

func (s *State) AddLearning(l string) {
	s.Learnings = append(s.Learnings, l)
}

func (s *State) AddDecision(d string) {
	s.Decisions = append(s.Decisions, d)
}

package ticket

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type TicketManager struct {
	tkPath     string
	ticketsDir string
}

func NewTicketManager(tkPath, ticketsDir string) *TicketManager {
	return &TicketManager{
		tkPath:     tkPath,
		ticketsDir: ticketsDir,
	}
}

func (tm *TicketManager) List() ([]Ticket, error) {
	entries, err := os.ReadDir(tm.ticketsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tickets directory: %w", err)
	}

	var tickets []Ticket
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		path := filepath.Join(tm.ticketsDir, entry.Name())
		ticket, err := tm.loadTicketFile(path)
		if err != nil {
			continue
		}
		tickets = append(tickets, *ticket)
	}

	return tickets, nil
}

func (tm *TicketManager) Get(id string) (*Ticket, error) {
	path := filepath.Join(tm.ticketsDir, id+".md")
	return tm.loadTicketFile(path)
}

func (tm *TicketManager) loadTicketFile(path string) (*Ticket, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read ticket file: %w", err)
	}

	content := string(data)
	if !strings.HasPrefix(content, "---\n") {
		return nil, fmt.Errorf("invalid ticket file: missing frontmatter")
	}

	endIdx := strings.Index(content[4:], "\n---")
	if endIdx == -1 {
		return nil, fmt.Errorf("invalid ticket file: unterminated frontmatter")
	}

	frontmatter := content[4 : 4+endIdx]

	var t Ticket
	if err := yaml.Unmarshal([]byte(frontmatter), &t); err != nil {
		return nil, fmt.Errorf("failed to parse ticket YAML: %w", err)
	}

	t.Updated = time.Now()

	return &t, nil
}

func (tm *TicketManager) SetStatus(id string, status Status) error {
	path := filepath.Join(tm.ticketsDir, id+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read ticket file: %w", err)
	}

	content := string(data)
	lines := strings.Split(content, "\n")
	var newLines []string
	inFrontmatter := false
	frontmatterDone := false

	for _, line := range lines {
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				newLines = append(newLines, line)
				continue
			}
			frontmatterDone = true
		}

		if inFrontmatter && !frontmatterDone {
			if strings.HasPrefix(line, "status:") {
				newLines = append(newLines, "status: "+string(status))
				continue
			}
		}

		newLines = append(newLines, line)
	}

	newContent := strings.Join(newLines, "\n")
	return os.WriteFile(path, []byte(newContent), 0644)
}

func (tm *TicketManager) AddNote(id, note string) error {
	path := filepath.Join(tm.ticketsDir, id+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read ticket file: %w", err)
	}

	content := string(data)
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	newNote := fmt.Sprintf("\n**%s**: %s", timestamp, note)

	newContent := content + newNote
	return os.WriteFile(path, []byte(newContent), 0644)
}

func (tm *TicketManager) GetChildren(parentID string) ([]Ticket, error) {
	all, err := tm.List()
	if err != nil {
		return nil, err
	}

	var children []Ticket
	for _, t := range all {
		if t.Parent == parentID {
			children = append(children, t)
		}
	}
	return children, nil
}

func (tm *TicketManager) GetDependencies(id string) ([]Ticket, error) {
	tk, err := tm.Get(id)
	if err != nil {
		return nil, err
	}

	var deps []Ticket
	for _, depID := range tk.Dependencies {
		dep, err := tm.Get(depID)
		if err != nil {
			continue
		}
		deps = append(deps, *dep)
	}
	return deps, nil
}

func (tm *TicketManager) GetExecutionOrder() ([][]string, error) {
	all, err := tm.List()
	if err != nil {
		return nil, err
	}

	depMap := make(map[string][]string)
	for _, t := range all {
		depMap[t.ID] = t.Dependencies
	}

	inDegree := make(map[string]int)
	for id := range depMap {
		inDegree[id] = 0
	}
	for _, deps := range depMap {
		for _, dep := range deps {
			if _, exists := inDegree[dep]; exists {
				continue
			}
		}
		for _, dep := range deps {
			if _, exists := depMap[dep]; exists {
				inDegree[dep]++
			}
		}
	}

	var levels [][]string
	remaining := make(map[string]bool)
	for id := range depMap {
		remaining[id] = true
	}

	for len(remaining) > 0 {
		var level []string
		for id := range remaining {
			ready := true
			for _, dep := range depMap[id] {
				if remaining[dep] {
					ready = false
					break
				}
			}
			if ready {
				level = append(level, id)
			}
		}

		if len(level) == 0 {
			break
		}

		levels = append(levels, level)
		for _, id := range level {
			delete(remaining, id)
		}
	}

	return levels, nil
}

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

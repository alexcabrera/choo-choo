package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	tea "charm.land/bubbletea/v2"

	"github.com/alexcabrera/choo-choo/internal/model"
	"github.com/alexcabrera/choo-choo/internal/orchestrator"
)

func main() {
	sequential := flag.Bool("sequential", false, "Run execution sequentially instead of in parallel")
	flag.Parse()

	projectDir := "."
	if args := flag.Args(); len(args) > 0 {
		projectDir = args[0]
	}

	absDir, err := filepath.Abs(projectDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving project directory: %v\n", err)
		os.Exit(1)
	}

	orch := orchestrator.NewOrchestrator(absDir)
	orch.SetSequential(*sequential)

	p := tea.NewProgram(model.InitialModel(orch))

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

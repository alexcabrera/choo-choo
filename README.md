# choo-choo

> "I choo-choo-choose you" - Ralph Wiggum

A TUI orchestration layer for the Crush coding agent with structured design → plan → execute → verify workflow.

## Overview

choo-choo wraps Crush in a beautiful terminal UI that provides:

- **Design phase**: Interactive chat with artifact preview
- **Plan phase**: Visual ticket tree exploration
- **Execute phase**: Kanban board with parallel execution
- **Verify phase**: Verification status per ticket
- **Gap closure**: Automatic fix loop with escalation

Built with [bubbletea](https://github.com/charmbracelet/bubbletea) from the [Charm](https://charm.sh) ecosystem.

## Installation

```bash
go install github.com/alexcabrera/choo-choo/cmd/choo-choo@latest
```

## Usage

```bash
# In an empty directory - initializes a new project
choo-choo

# In an existing project - continues from saved state
choo-choo

# Sequential execution (no parallel)
choo-choo --sequential
```

## Requirements

- [Crush](https://github.com/charmbracelet/crush) - The AI coding agent
- [choo-choo-skills](../agent-skills/choo-choo-skills/) - Skills symlinked to `~/.config/crush/skills/`

## Skills

choo-choo orchestrates these skills (can also be used directly in Crush):

| Skill | Phase | Purpose |
|-------|-------|---------|
| `design` | Design | Transform ideas into requirements + architecture |
| `plan` | Planning | Decompose design into ticket tree |
| `validate` | Validation | Verify plan is executable |
| `execute` | Execution | Implement a single ticket |
| `verify` | Verification | Confirm implementation matches criteria |
| `close-gaps` | Gap Closure | Fix discrepancies found by verify |

## Development

```bash
# Run in development
go run ./cmd/choo-choo

# Build
go build -o choo-choo ./cmd/choo-choo

# Test
go test ./...
```

## Architecture

```
┌─────────────────────────────────────────┐
│           choo-choo TUI                  │
│  ┌─────────────────────────────────────┐│
│  │         bubbletea Model             ││
│  │  ┌─────────┐ ┌─────────┐ ┌────────┐││
│  │  │  Chat   │ │ Kanban  │ │ Popup  │││
│  │  └─────────┘ └─────────┘ └────────┘││
│  └─────────────────────────────────────┘│
│  ┌─────────────────────────────────────┐│
│  │         Orchestrator                ││
│  │  - Phase management                 ││
│  │  - Crush process spawning           ││
│  │  - STATE.md read/write              ││
│  └─────────────────────────────────────┘│
└─────────────────────────────────────────┘
              │
              ▼
        ┌──────────┐
        │  Crush   │
        │(headless)│
        └──────────┘
              │
              ▼
        ┌──────────────────┐
        │ choo-choo-skills │
        │  design, plan,   │
        │  execute, etc.   │
        └──────────────────┘
```

## Related

- [choo-choo-skills](../agent-skills/choo-choo-skills/) - The skills used by this TUI
- [Crush](https://github.com/charmbracelet/crush) - The underlying AI coding agent
- [Ralph TUI](https://github.com/subsy/ralph-tui) - Inspiration for agent loop orchestration

## License

MIT

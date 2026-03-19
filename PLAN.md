# choo-choo TUI Plan

## Overview

This plan decomposes the choo-choo TUI implementation into a hierarchical ticket structure with clear dependencies enabling parallel execution where possible.

## Statistics

- **Total tickets**: 58
- **Epic**: 1
- **Stories**: 8
- **Tasks**: 49

## Hierarchy

```
Epic: cc-ef7f - choo-choo TUI Implementation
├── Story: cc-uog1 - State Manager (5 tasks)
├── Story: cc-5g3r - Crush Runner (5 tasks)
├── Story: cc-eryx - Ticket Manager (6 tasks)
├── Story: cc-z0mt - Orchestrator (8 tasks)
├── Story: cc-baea - TUI Model (5 tasks)
├── Story: cc-77cj - Chat Model (5 tasks)
├── Story: cc-2yy2 - Kanban Model (6 tasks)
├── Story: cc-vux3 - Popup Model (6 tasks)
└── Integration Tests (3 tasks)
```

## Dependency Levels

### Level 0 (No dependencies - can start immediately)
- cc-mlnx: Define Phase type and constants
- cc-gcfk: Define StreamEvent and RunOptions types
- cc-gjz5: Define ChatModel and ChatMessage types
- cc-zoev: Define Ticket types

### Level 1
- cc-mv1n: Implement State struct
- cc-245p: Implement CrushRunner struct
- cc-0aim: Implement ChatModel constructor
- cc-ik7b: Implement TicketManager struct
- cc-ye9o: Define KanbanModel struct
- cc-3mtz: Define PopupModel and Tab types

### Level 2
- cc-l2w6: Implement LoadState function
- cc-xmnh: Implement State.Save method
- cc-dlgr: Implement State setters
- cc-qej4: Implement Run method
- cc-y1th: Implement process lifecycle
- cc-0loh: Implement ChatModel Update
- cc-fyja: Implement AddMessage and SetStreaming
- cc-w23o: Implement List and Get methods
- cc-4rcv: Implement KanbanModel constructor
- cc-lpkj: Implement PopupModel constructor

### Level 3+
Continues with Orchestrator tasks, then TUI Model, then integration tests.

See `.tickets/` for full dependency details.

## Critical Path

```
cc-mlnx → cc-mv1n → cc-dlgr → cc-j9ac → cc-0xcq → cc-x0nj → cc-k4cx → cc-rc6w → cc-jhsl → cc-8cey → cc-1t30 → cc-c1fl → cc-rca8 → cc-v0sc → cc-n2aw
```

## Risks

1. **Crush API changes**: Headless mode flags may change - verify with Crush docs
2. **bubbletea v2 import**: Must use `charm.land/bubbletea/v2` not github.com
3. **Parallel execution git conflicts**: Multiple Crush processes may conflict on git operations - use `--sequential` if needed

## Success Criteria

- All 12 acceptance criteria from design.md pass
- All integration tests pass
- TUI renders correctly at 80x24 minimum size
- Crash recovery works via STATE.md

## Testing

Uses `tmux-test` skill for TUI testing:
- Start choo-choo in isolated tmux session
- Send keystrokes programmatically
- Capture screen for verification
- Detect hangs with timeout

See `choo-choo-skills/tmux-test/` for details.

## Reference

- Design doc: `~/Code/agent-skills/specs/choo-choo/design.md`
- Skills: `~/Code/agent-skills/choo-choo-skills/`
- TUI code: `~/Code/choo-choo/`

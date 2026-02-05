# Bubble Tea Implementation Plan

## Phase 1: Foundation & List TUI

### 1.1 Add Bubble Tea Dependencies
- Add `github.com/charmbracelet/bubbletea`
- Add `github.com/charmbracelet/bubbles`
- Add `github.com/charmbracelet/lipgloss`
- Already have: `github.com/charmbracelet/huh`

### 1.2 Create TUI Package Structure
```
backend/internal/tui/
├── common/
│   ├── keys.go         # Shared keybindings
│   ├── styles.go       # Lip Gloss styles
│   └── messages.go     # Common message types
├── components/
│   ├── issuelist/      # Reusable issue list
│   ├── issuedetail/    # Issue detail view
│   ├── statusbar/      # Bottom status/help bar
│   └── filterinput/    # Filter text input
└── commands/
    ├── list/           # bd list --tui
    ├── graph/          # bd graph --tui
    └── status/         # bd status --tui
```

### 1.3 Implement `bd list --tui`
- Add `--tui` flag to existing `listCmd`
- Create `listModel` with:
  - `list.Model` for issue browsing
  - `textinput.Model` for filtering
  - `help.Model` for keybindings
- Keybindings:
  - `j/k` or arrows: navigate
  - `/`: focus filter
  - `Enter`: view details
  - `c`: close issue
  - `e`: edit (opens $EDITOR)
  - `p`: cycle priority
  - `a`: assign
  - `q`: quit

### 1.4 Implement Issue Detail View
- `viewport.Model` for scrollable content
- Show: title, status, priority, description, deps, comments
- Actions: close, edit, assign, add dep

---

## Phase 2: Graph TUI

### 2.1 Implement `bd graph --tui`
- Add `--tui` flag to existing `graphCmd`
- Create tree model with expand/collapse
- Keybindings:
  - `j/k`: navigate nodes
  - `h/l` or `←/→`: collapse/expand
  - `Enter`: view issue details
  - `f`: focus on subtree
  - `b`: jump to blockers
  - `Esc`: unfocus/back

### 2.2 Graph Rendering
- Use box-drawing characters for tree
- Color by status
- Show expand indicators (`▶`/`▼`)
- Highlight selected node

---

## Phase 3: Status Dashboard

### 3.1 Implement `bd status --tui`
- Add `--tui` flag to existing `statusCmd`
- Create dashboard model with:
  - `progress.Model` for distribution bars
  - `timer` for auto-refresh
  - Alert list for actionable items

### 3.2 Dashboard Features
- Summary counts with visual bars
- Activity sparkline (24h)
- Alerts: overdue, blocked, ready-to-close epics
- Keybindings:
  - `r`: manual refresh
  - `l`: jump to list (ready work)
  - `o`: show overdue
  - `q`: quit

---

## Phase 4: Enhanced Create Form

### 4.1 Dependency Picker Component
- Integrate `list.Model` into form flow
- Fuzzy search existing issues
- Select dependency type (blocks, related, etc.)
- Multi-select support

### 4.2 Preview Step
- Add final step showing issue preview
- Confirm or go back to edit
- Show rendered markdown for description

### 4.3 Template Selector
- List available templates
- Pre-fill form from template
- Show template preview

---

## Phase 5: Full TUI Mode

### 5.1 Create `bd tui` Command
- New command: combined interface
- Split view: list + detail panes
- Tab bar: Issues | Graph | Status
- Command palette (⌘K style)

### 5.2 Split View Layout
- Resizable panes
- Focus switching with Tab
- Synchronized selection

### 5.3 Command Palette
- Fuzzy command search
- Recent commands
- Contextual actions based on selection

---

## Phase 6: Daemon Monitor

### 6.1 Implement `bd daemon --monitor`
- Real-time stats display
- Log tail with viewport
- Recent activity list

### 6.2 Monitor Actions
- Trigger sync
- Restart daemon
- Compact database
- Toggle log verbosity

---

## Integration Checklist

### CLI Flags
- [ ] `bd list --tui`
- [ ] `bd graph --tui`
- [ ] `bd status --tui`
- [ ] `bd tui` (new command)
- [ ] `bd daemon --monitor`

### Shared Components
- [ ] Issue list with filtering
- [ ] Issue detail view
- [ ] Keybinding help bar
- [ ] Status/progress bars
- [ ] Spinner for async ops

### Daemon Integration
- [ ] RPC calls return tea.Cmd
- [ ] Async operation handling
- [ ] Error display in TUI

### Web Terminal Integration
- [ ] Test with xterm.js
- [ ] Resize handling
- [ ] Mouse support (optional)

---

## File Changes Summary

### New Files
```
backend/internal/tui/common/keys.go
backend/internal/tui/common/styles.go
backend/internal/tui/common/messages.go
backend/internal/tui/components/issuelist/model.go
backend/internal/tui/components/issuedetail/model.go
backend/internal/tui/components/statusbar/model.go
backend/internal/tui/commands/list/model.go
backend/internal/tui/commands/graph/model.go
backend/internal/tui/commands/status/model.go
backend/internal/tui/commands/full/model.go
backend/cmd/bd/tui.go
```

### Modified Files
```
backend/cmd/bd/list.go      # Add --tui flag
backend/cmd/bd/graph.go     # Add --tui flag
backend/cmd/bd/status.go    # Add --tui flag
backend/cmd/bd/daemon.go    # Add --monitor flag
backend/cmd/bd/create_form.go  # Add dep picker, preview
backend/go.mod              # Add bubbletea, bubbles, lipgloss
```

---

## Testing Strategy

### Unit Tests
- Model state transitions
- Update function with mock messages
- View rendering

### Integration Tests
- TUI with mock store
- Daemon RPC integration
- Keyboard input sequences

### Manual Testing
- Native terminal (iTerm, Terminal.app)
- xterm.js in browser
- SSH sessions
- tmux/screen

---

## Timeline Estimate

| Phase | Scope |
|-------|-------|
| Phase 1 | Foundation + List TUI |
| Phase 2 | Graph TUI |
| Phase 3 | Status Dashboard |
| Phase 4 | Enhanced Create Form |
| Phase 5 | Full TUI Mode |
| Phase 6 | Daemon Monitor |

---

## Dependencies

```go
require (
    github.com/charmbracelet/bubbletea v1.x
    github.com/charmbracelet/bubbles v0.x
    github.com/charmbracelet/lipgloss v1.x
    github.com/charmbracelet/huh v0.x  // already have
)
```

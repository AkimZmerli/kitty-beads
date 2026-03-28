# Bubble Tea Implementation Plan

## Phase 1: Foundation & List TUI

### 1.1 Add Bubble Tea Dependencies
- Add `github.com/charmbracelet/bubbletea`
- Add `github.com/charmbracelet/bubbles`
- Add `github.com/charmbracelet/lipgloss`
- Add `github.com/charmbracelet/glamour` (markdown rendering)
- Add `github.com/charmbracelet/log` (structured logging)
- Add `github.com/charmbracelet/harmonica` (spring animations)
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

## Phase 7: SSH Access (Wish)

### 7.1 Implement `bd serve` Command
- Start SSH server with Wish
- Serve full TUI over SSH
- Auto-generate host keys on first run

```go
// backend/cmd/bd/serve.go
import (
    "github.com/charmbracelet/wish"
    "github.com/charmbracelet/wish/bubbletea"
)

func runServe(cmd *cobra.Command, args []string) error {
    s, err := wish.NewServer(
        wish.WithAddress(":2222"),
        wish.WithHostKeyPath(".ssh/beads_host_key"),
        wish.WithMiddleware(
            bubbletea.Middleware(teaHandler),
            logging.Middleware(),
        ),
    )
    return s.ListenAndServe()
}
```

### 7.2 SSH Features
- Public key authentication
- Per-user sessions
- Shared team access to beads
- Read-only guest mode

### 7.3 Use Cases
- Access beads from any machine: `ssh beads.local -p 2222`
- Team collaboration without local install
- Headless server management
- Remote pairing sessions

---

## Integration Checklist

### CLI Flags
- [ ] `bd list --tui`
- [ ] `bd graph --tui`
- [ ] `bd status --tui`
- [ ] `bd tui` (new command)
- [ ] `bd daemon --monitor`
- [ ] `bd serve` (SSH server)

### Shared Components
- [ ] Issue list with filtering
- [ ] Issue detail view
- [ ] Keybinding help bar
- [ ] Status/progress bars
- [ ] Spinner for async ops
- [ ] Markdown renderer (Glamour)
- [ ] Structured logger (Log)
- [ ] Animation helpers (Harmonica)

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
backend/internal/tui/common/markdown.go    # Glamour markdown rendering
backend/internal/tui/common/logger.go      # Charm log integration
backend/internal/tui/common/animation.go   # Harmonica spring animations
backend/internal/tui/components/issuelist/model.go
backend/internal/tui/components/issuedetail/model.go
backend/internal/tui/components/statusbar/model.go
backend/internal/tui/commands/list/model.go
backend/internal/tui/commands/graph/model.go
backend/internal/tui/commands/status/model.go
backend/internal/tui/commands/full/model.go
backend/cmd/bd/tui.go
backend/cmd/bd/serve.go                    # Wish SSH server
```

### Modified Files
```
backend/cmd/bd/list.go      # Add --tui flag
backend/cmd/bd/graph.go     # Add --tui flag
backend/cmd/bd/status.go    # Add --tui flag
backend/cmd/bd/daemon.go    # Add --monitor flag
backend/cmd/bd/create_form.go  # Add dep picker, preview
backend/go.mod              # Add bubbletea, bubbles, lipgloss, glamour, log, harmonica, wish
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
- SSH sessions (including via `bd serve`)
- tmux/screen
- Remote access via Wish SSH server

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
| Phase 7 | SSH Access (Wish) |

---

## Dependencies

```go
require (
    github.com/charmbracelet/bubbletea v1.x
    github.com/charmbracelet/bubbles v0.x
    github.com/charmbracelet/lipgloss v1.x
    github.com/charmbracelet/glamour v0.x   // markdown rendering
    github.com/charmbracelet/log v0.4.x     // structured logging
    github.com/charmbracelet/harmonica v0.x // spring animations
    github.com/charmbracelet/huh v0.x       // already have
    github.com/charmbracelet/wish v1.x      // SSH server (Phase 7)
)
```

---

## Charmbracelet Library Usage

### Glamour - Markdown Rendering
Used for rendering issue descriptions, comments, and bead content in the TUI.

```go
// backend/internal/tui/common/markdown.go
import "github.com/charmbracelet/glamour"

func RenderMarkdown(content string) (string, error) {
    return glamour.Render(content, "tokyo-night")
}
```

**Integration points:**
- Issue detail view: render description
- `bd show <issue>` output
- Comment rendering
- Help text display

### Log - Structured Logging
Consistent, styled logging across all TUI operations.

```go
// backend/internal/tui/common/logger.go
import "github.com/charmbracelet/log"

var Logger = log.NewWithOptions(os.Stderr, log.Options{
    ReportTimestamp: true,
    Level:           log.InfoLevel,
})
```

**Integration points:**
- Debug TUI state transitions
- PTY operations logging
- Daemon monitor logs

### Harmonica - Spring Animations
Smooth transitions for pane focus, resizing, and UI feedback.

```go
// backend/internal/tui/common/animation.go
import "github.com/charmbracelet/harmonica"

// Create spring for smooth pane transitions
spring := harmonica.NewSpring(harmonica.FPS(60), 6.0, 0.5)
```

**Integration points:**
- Focus indicator transitions
- Pane resize animations
- List scroll smoothing
- Progress bar animations

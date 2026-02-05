# Splitty: Bubble Tea Split Terminal Library

## Vision

**Splitty** is a standalone, reusable Go library for iTerm2-style split pane terminal multiplexing in Bubble Tea applications. It will be developed as an independent package and imported into kitty-beads.

```
github.com/AkimZmerli/splitty  →  imported by  →  kitty-beads
```

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          splitty.Manager                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                          │
│   ┌─────────────────────┐    │    ┌─────────────────────────────────┐   │
│   │     Pane 1          │    │    │         Pane 2                  │   │
│   │  ┌───────────────┐  │    │    │  ┌───────────────────────────┐  │   │
│   │  │   PTY 1       │  │  ──┼──  │  │         PTY 2             │  │   │
│   │  │  (bd list)    │  │ divider │  │       (bash)              │  │   │
│   │  └───────────────┘  │    │    │  └───────────────────────────┘  │   │
│   │  [FOCUSED]          │    │    │                                 │   │
│   └─────────────────────┘    │    └─────────────────────────────────┘   │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│  Ctrl+\ split-v  Ctrl+- split-h  Ctrl+w close  Ctrl+hjkl nav  q quit   │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Repository Structure

```
github.com/AkimZmerli/splitty/
├── splitty.go            # Manager, New(), tea.Model implementation
├── pane.go               # Pane type with PTY
├── node.go               # SplitNode, SplitContainer (binary tree)
├── layout.go             # Layout calculation algorithms
├── navigation.go         # Focus navigation logic
├── resize.go             # Resize handling
├── options.go            # Functional options pattern
├── keys.go               # KeyMap with configurable bindings
├── theme.go              # Theme type + built-in themes
├── messages.go           # Public tea.Msg types
├── persist.go            # Layout save/restore
├── presets.go            # Named layout presets
├── logger.go             # Logging with charmbracelet/log
├── terminal/
│   ├── pty.go            # PTY management (creack/pty)
│   ├── buffer.go         # ANSI-aware terminal buffer
│   ├── parser.go         # ANSI escape sequence parser
│   └── screen.go         # Virtual screen representation
├── examples/
│   ├── basic/main.go           # Minimal usage
│   ├── custom-theme/main.go    # Custom styling
│   ├── custom-keys/main.go     # Custom keybindings
│   ├── embedded/main.go        # Embed in larger app
│   └── presets/main.go         # Using layout presets
├── _testdata/                  # Test fixtures
├── splitty_test.go
├── layout_test.go
├── navigation_test.go
├── README.md
├── LICENSE                     # MIT
├── go.mod
├── go.sum
└── .github/
    └── workflows/
        └── ci.yml              # Tests, lint, release
```

---

## Public API Design

### Core Types

```go
package splitty

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    "github.com/charmbracelet/log"
)

// Manager is the main Bubble Tea model for split pane management
type Manager struct {
    // unexported fields
}

// Pane represents a single terminal pane with its own PTY
type Pane struct {
    ID     string
    Title  string
    CWD    string
    Width  int
    Height int
}

// Direction specifies split orientation
type Direction int

const (
    Vertical   Direction = iota // Side by side (│)
    Horizontal                   // Stacked (─)
)

// Ensure Manager implements tea.Model
var _ tea.Model = (*Manager)(nil)
```

### Constructor & Options

```go
// New creates a new split pane manager with the given options
func New(opts ...Option) *Manager

// Option configures the Manager
type Option func(*Manager)

// WithShell sets the shell command for new panes (default: $SHELL or /bin/bash)
func WithShell(shell string) Option

// WithTheme sets the visual theme
func WithTheme(theme Theme) Option

// WithKeyMap sets custom keybindings
func WithKeyMap(km KeyMap) Option

// WithMinSize sets minimum pane dimensions (default: 10x5)
func WithMinSize(width, height int) Option

// WithStatusBar enables/disables the bottom status bar (default: true)
func WithStatusBar(enabled bool) Option

// WithMouse enables/disables mouse support (default: true)
func WithMouse(enabled bool) Option

// WithPreset initializes with a named layout preset
func WithPreset(name string) Option

// WithEnv sets additional environment variables for PTY sessions
func WithEnv(env []string) Option

// WithLogger sets a charmbracelet/log logger for debugging (default: nil/disabled)
func WithLogger(logger *log.Logger) Option
```

### Manager Methods

```go
// Bubble Tea interface
func (m *Manager) Init() tea.Cmd
func (m *Manager) Update(msg tea.Msg) (tea.Model, tea.Cmd)
func (m *Manager) View() string

// Split operations
func (m *Manager) Split(dir Direction) tea.Cmd
func (m *Manager) Close() tea.Cmd
func (m *Manager) ClosePane(id string) tea.Cmd

// Focus operations
func (m *Manager) Focus(dir Direction)
func (m *Manager) FocusPane(id string)
func (m *Manager) FocusedPane() *Pane
func (m *Manager) Panes() []*Pane

// Layout operations
func (m *Manager) Zoom() tea.Cmd
func (m *Manager) Unzoom() tea.Cmd
func (m *Manager) IsZoomed() bool
func (m *Manager) Swap()
func (m *Manager) Resize(dir Direction, delta float64)

// Persistence
func (m *Manager) SaveLayout(path string) error
func (m *Manager) LoadLayout(path string) error

// Broadcast
func (m *Manager) SetBroadcast(enabled bool)
func (m *Manager) IsBroadcasting() bool

// Send input to focused pane (or all if broadcasting)
func (m *Manager) SendInput(data []byte)
```

### Themes

```go
// Theme defines the visual styling for split panes
type Theme struct {
    FocusedBorder   lipgloss.Style
    UnfocusedBorder lipgloss.Style
    Divider         lipgloss.Style
    DividerChar     string
    StatusBar       lipgloss.Style
    StatusText      lipgloss.Style
    ZoomIndicator   string
    BroadcastIndicator string
}

// Built-in themes
var (
    DefaultTheme = Theme{...}   // Works everywhere
    TokyoNight   = Theme{...}   // Tokyo Night colors
    Dracula      = Theme{...}   // Dracula colors
    Nord         = Theme{...}   // Nord colors
    Catppuccin   = Theme{...}   // Catppuccin Mocha
)
```

### KeyMap

```go
// KeyMap defines all keybindings for split operations
type KeyMap struct {
    SplitVertical   key.Binding
    SplitHorizontal key.Binding
    Close           key.Binding
    FocusLeft       key.Binding
    FocusDown       key.Binding
    FocusUp         key.Binding
    FocusRight      key.Binding
    FocusCycle      key.Binding
    Zoom            key.Binding
    Swap            key.Binding
    ResizeLeft      key.Binding
    ResizeRight     key.Binding
    ResizeUp        key.Binding
    ResizeDown      key.Binding
    Broadcast       key.Binding
}

// DefaultKeyMap returns the default keybindings
func DefaultKeyMap() KeyMap
```

### Messages (for embedding)

```go
// PaneSplitMsg is sent when a pane is split
type PaneSplitMsg struct {
    ParentID  string
    NewPaneID string
    Direction Direction
}

// PaneClosedMsg is sent when a pane is closed
type PaneClosedMsg struct {
    PaneID string
}

// PaneFocusedMsg is sent when focus changes
type PaneFocusedMsg struct {
    PaneID string
}

// PaneOutputMsg contains output from a pane's PTY
type PaneOutputMsg struct {
    PaneID string
    Data   []byte
}
```

### Presets

```go
// Built-in layout presets
const (
    PresetSingle = "single"  // One pane (default)
    PresetDev    = "dev"     // 60/40 vertical split
    PresetTriple = "triple"  // Three columns
    PresetQuad   = "quad"    // 2x2 grid
)

// RegisterPreset adds a custom preset
func RegisterPreset(name string, builder func() *node)
```

---

## Usage Examples

### Basic Usage

```go
package main

import (
    "log"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/AkimZmerli/splitty"
)

func main() {
    m := splitty.New()

    p := tea.NewProgram(m,
        tea.WithAltScreen(),
        tea.WithMouseCellMotion(),
    )

    if _, err := p.Run(); err != nil {
        log.Fatal(err)
    }
}
```

### Custom Configuration

```go
func main() {
    m := splitty.New(
        splitty.WithShell("/bin/zsh"),
        splitty.WithTheme(splitty.TokyoNight),
        splitty.WithPreset(splitty.PresetDev),
        splitty.WithMinSize(20, 10),
        splitty.WithEnv([]string{
            "EDITOR=nvim",
            "PAGER=less",
        }),
    )

    p := tea.NewProgram(m, tea.WithAltScreen())
    p.Run()
}
```

### Custom Keybindings

```go
func main() {
    keys := splitty.DefaultKeyMap()
    keys.SplitVertical = key.NewBinding(
        key.WithKeys("ctrl+v"),
        key.WithHelp("ctrl+v", "split vertical"),
    )
    keys.SplitHorizontal = key.NewBinding(
        key.WithKeys("ctrl+s"),
        key.WithHelp("ctrl+s", "split horizontal"),
    )

    m := splitty.New(splitty.WithKeyMap(keys))

    tea.NewProgram(m, tea.WithAltScreen()).Run()
}
```

### Debug Logging

```go
func main() {
    // Create a file logger for debugging
    f, _ := os.OpenFile("splitty.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    logger := log.NewWithOptions(f, log.Options{
        ReportTimestamp: true,
        Level:           log.DebugLevel,
    })

    m := splitty.New(
        splitty.WithLogger(logger),
        splitty.WithTheme(splitty.TokyoNight),
    )

    tea.NewProgram(m, tea.WithAltScreen()).Run()
}
```

### Embedding in Larger App

```go
type MyApp struct {
    splits  *splitty.Manager
    sidebar sidebar.Model
    width   int
    height  int
}

func NewApp() *MyApp {
    return &MyApp{
        splits: splitty.New(
            splitty.WithStatusBar(false), // We'll render our own
        ),
        sidebar: sidebar.New(),
    }
}

func (m *MyApp) Init() tea.Cmd {
    return m.splits.Init()
}

func (m *MyApp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height

    case tea.KeyMsg:
        if msg.String() == "ctrl+b" {
            m.sidebar.Toggle()
            return m, nil
        }

    case splitty.PaneFocusedMsg:
        // React to focus changes
        m.sidebar.SetContext(msg.PaneID)
    }

    // Update splits
    var cmd tea.Cmd
    newSplits, cmd := m.splits.Update(msg)
    m.splits = newSplits.(*splitty.Manager)
    cmds = append(cmds, cmd)

    return m, tea.Batch(cmds...)
}

func (m *MyApp) View() string {
    sidebarWidth := 30
    if !m.sidebar.Visible() {
        sidebarWidth = 0
    }

    splitsView := lipgloss.NewStyle().
        Width(m.width - sidebarWidth).
        Height(m.height).
        Render(m.splits.View())

    if sidebarWidth > 0 {
        return lipgloss.JoinHorizontal(
            lipgloss.Top,
            m.sidebar.View(),
            splitsView,
        )
    }
    return splitsView
}
```

---

## Integration with kitty-beads

### Import in go.mod

```go
// kitty-beads/backend/go.mod
require (
    github.com/AkimZmerli/splitty v0.1.0
    // ... other deps
)
```

### Usage in bd tui Command

```go
// kitty-beads/backend/cmd/bd/tui.go
package main

import (
    "github.com/AkimZmerli/splitty"
    "github.com/spf13/cobra"
    tea "github.com/charmbracelet/bubbletea"
)

var tuiCmd = &cobra.Command{
    Use:   "tui",
    Short: "Launch interactive TUI with split panes",
    RunE: func(cmd *cobra.Command, args []string) error {
        preset, _ := cmd.Flags().GetString("preset")

        opts := []splitty.Option{
            splitty.WithTheme(splitty.TokyoNight),
        }
        if preset != "" {
            opts = append(opts, splitty.WithPreset(preset))
        }

        m := splitty.New(opts...)

        p := tea.NewProgram(m,
            tea.WithAltScreen(),
            tea.WithMouseCellMotion(),
        )

        _, err := p.Run()
        return err
    },
}

func init() {
    tuiCmd.Flags().String("preset", "", "Layout preset (dev, triple, quad)")
    rootCmd.AddCommand(tuiCmd)
}
```

### Custom kitty-beads Theme

```go
// kitty-beads/backend/internal/tui/theme.go
package tui

import (
    "github.com/AkimZmerli/splitty"
    "github.com/charmbracelet/lipgloss"
)

// KittyBeadsTheme extends Tokyo Night with beads-specific styling
var KittyBeadsTheme = splitty.Theme{
    FocusedBorder: lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#7AA2F7")),
    UnfocusedBorder: lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(lipgloss.Color("#414868")),
    Divider: lipgloss.NewStyle().
        Foreground(lipgloss.Color("#565F89")),
    DividerChar: "│",
    StatusBar: lipgloss.NewStyle().
        Background(lipgloss.Color("#1A1B26")).
        Foreground(lipgloss.Color("#A9B1D6")),
    ZoomIndicator:      "◉ ZOOM",
    BroadcastIndicator: "📡 BROADCAST",
}
```

---

## Implementation Phases

### Phase 1: Core Library Setup
**Repo:** `github.com/AkimZmerli/splitty`

- [ ] Initialize Go module
- [ ] Create basic file structure
- [ ] Implement `SplitNode` and `SplitContainer` types
- [ ] Implement `Manager` with `tea.Model` interface
- [ ] Basic layout calculation
- [ ] Single pane rendering
- [ ] Add MIT LICENSE and README

**Deliverable:** `splitty.New()` creates a single-pane terminal

### Phase 2: Split Operations

- [ ] Implement vertical split (`Ctrl+\`)
- [ ] Implement horizontal split (`Ctrl+-`)
- [ ] Implement close pane (`Ctrl+w`)
- [ ] Binary tree manipulation for nested splits
- [ ] Recalculate layout on split/close
- [ ] Focus tracking

**Deliverable:** Can create and close splits, focus follows

### Phase 3: PTY Management

- [ ] Create `terminal/pty.go` with PTY wrapper
- [ ] Implement `terminal/buffer.go` for ANSI buffer
- [ ] Create `terminal/parser.go` for escape sequences
- [ ] Spawn shell per pane
- [ ] Route input to focused pane
- [ ] Handle PTY resize (SIGWINCH)
- [ ] Render PTY output

**Deliverable:** Each pane runs independent shell

### Phase 4: Navigation

- [ ] Directional focus (Ctrl+hjkl)
- [ ] Tab cycling
- [ ] Focus indicator styling
- [ ] Pane ID/title in border

**Deliverable:** Intuitive pane navigation

### Phase 5: Theming & Configuration

- [ ] Implement `Theme` type
- [ ] Built-in themes (Default, TokyoNight, Dracula, Nord)
- [ ] Implement `KeyMap` type
- [ ] Functional options pattern
- [ ] `WithShell`, `WithTheme`, `WithKeyMap`, etc.

**Deliverable:** Fully configurable library

### Phase 6: Resizing

- [ ] Keyboard resize (Ctrl+Shift+arrows)
- [ ] Mouse divider detection
- [ ] Mouse drag to resize
- [ ] Minimum size constraints
- [ ] Ratio clamping

**Deliverable:** Resize panes with keyboard and mouse

### Phase 7: Advanced Features

- [ ] Zoom/maximize toggle (Ctrl+z)
- [ ] Layout persistence (save/load JSON)
- [ ] Broadcast mode (Ctrl+b)
- [ ] Pane swap (Ctrl+x)
- [ ] Named presets (dev, triple, quad)
- [ ] Custom preset registration

**Deliverable:** Power-user features complete

### Phase 8: Polish & Release

- [ ] Comprehensive tests
- [ ] Examples directory
- [ ] Full README with badges
- [ ] godoc comments
- [ ] GitHub Actions CI
- [ ] Tag v0.1.0
- [ ] Submit to awesome-bubbletea

**Deliverable:** Published, production-ready library

### Phase 9: kitty-beads Integration

- [ ] Add `github.com/AkimZmerli/splitty` to go.mod
- [ ] Create `bd tui` command
- [ ] Create KittyBeadsTheme
- [ ] Integrate with daemon monitor
- [ ] Test with xterm.js frontend

**Deliverable:** Split terminals in kitty-beads

---

## Default Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+\` | Split vertical (side by side) |
| `Ctrl+-` | Split horizontal (stacked) |
| `Ctrl+w` | Close focused pane |
| `Ctrl+h` | Focus left |
| `Ctrl+j` | Focus down |
| `Ctrl+k` | Focus up |
| `Ctrl+l` | Focus right |
| `Tab` | Cycle focus forward |
| `Shift+Tab` | Cycle focus backward |
| `Ctrl+z` | Toggle zoom |
| `Ctrl+x` | Swap with sibling |
| `Ctrl+b` | Toggle broadcast mode |
| `Ctrl+Shift+←` | Resize left |
| `Ctrl+Shift+→` | Resize right |
| `Ctrl+Shift+↑` | Resize up |
| `Ctrl+Shift+↓` | Resize down |

---

## Dependencies

```go
// splitty/go.mod
module github.com/AkimZmerli/splitty

go 1.21

require (
    github.com/charmbracelet/bubbletea v1.3.0
    github.com/charmbracelet/bubbles v0.20.0
    github.com/charmbracelet/lipgloss v1.0.0
    github.com/charmbracelet/log v0.4.0
    github.com/creack/pty v1.1.21
    golang.org/x/term v0.27.0
)
```

---

## Testing Strategy

### Unit Tests
```go
func TestLayoutCalculation(t *testing.T) { ... }
func TestSplitVertical(t *testing.T) { ... }
func TestFocusNavigation(t *testing.T) { ... }
func TestResizeBoundaries(t *testing.T) { ... }
func TestThemeApplication(t *testing.T) { ... }
```

### Integration Tests
```go
func TestFullSplitWorkflow(t *testing.T) { ... }
func TestLayoutPersistence(t *testing.T) { ... }
func TestPTYSpawnAndIO(t *testing.T) { ... }
```

### Manual Testing Matrix
- [ ] macOS Terminal.app
- [ ] macOS iTerm2
- [ ] macOS Alacritty
- [ ] Linux GNOME Terminal
- [ ] Linux Kitty
- [ ] Windows Terminal (WSL)
- [ ] SSH session
- [ ] tmux (nested)
- [ ] xterm.js (browser)

---

## UI Mockups

### Single Pane (Initial State)
```
┌─────────────────────────────────────────────────────────────────────────┐
│ ~/project $                                                              │
│                                                                          │
│                                                                          │
│                                                                          │
│                                                                          │
│                                                                          │
│                                                                          │
│                                                                          │
├─────────────────────────────────────────────────────────────────────────┤
│ Ctrl+\ split  Ctrl+- hsplit  Ctrl+w close                      Pane 1/1 │
└─────────────────────────────────────────────────────────────────────────┘
```

### Vertical Split
```
┌────────────────────────────────┬────────────────────────────────────────┐
│ ~/project $ bd list            │ ~/project $ npm run dev                │
│                                │                                        │
│  > ● P1  issue-1         epic  │ > kitty-beads@1.0.0 dev                │
│    ○ P1  issue-2         task  │ > vite                                 │
│    ◐ P2  issue-3         epic  │                                        │
│                                │   VITE v5.0.0  ready in 234 ms         │
│                                │                                        │
│  [FOCUSED]                     │                                        │
├────────────────────────────────┴────────────────────────────────────────┤
│ Ctrl+hjkl nav  Ctrl+z zoom  Ctrl+w close                       Pane 1/2 │
└─────────────────────────────────────────────────────────────────────────┘
```

### Quad Preset
```
┌─────────────────────────────────┬───────────────────────────────────────┐
│ ~/project $ bd status           │ ~/project $ docker logs -f api        │
│                                 │                                       │
│  Summary: 127 total             │ [INFO] Server started on :8080        │
│  Open: 23  Blocked: 3           │ [INFO] Request: GET /api/health       │
│                                 │                                       │
├─────────────────────────────────┼───────────────────────────────────────┤
│ ~/project $ bd graph            │ ~/project $ htop                      │
│                                 │                                       │
│  ▼ ● epic-1 [epic]             │   CPU [████████░░] 42%                │
│    ├── ✓ task-1  [FOCUSED]     │   MEM [██████░░░░] 31%                │
│    └── ○ task-2                │                                       │
├─────────────────────────────────┴───────────────────────────────────────┤
│ Ctrl+hjkl nav  Ctrl+z zoom                                     Pane 3/4 │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## Related Work

| Project | Language | Notes |
|---------|----------|-------|
| tmux | C | Gold standard, server-client |
| iTerm2 | Obj-C | macOS native, profiles |
| Zellij | Rust | Plugins, WASM |
| WezTerm | Rust | Lua config, GPU |
| [teacup](https://github.com/knipferrc/teacup) | Go | Bubble Tea layout utils |

**Splitty differentiator:** Purpose-built for Bubble Tea embedding, not a standalone terminal emulator.

     Splitty: Build Plan

     Context

     Splitty is a standalone, open-source Go library for iTerm2-style split pane terminal
     multiplexing in Bubble Tea applications. The kitty-beads repo already has a full epic
     tracking this: kitty-beads-hq5 with 8 phase sub-tasks (hq5.1 through hq5.8). We will update
     the epic status as we build, not create a new one.

     The repo at /Users/webdev4life/git/splitty is empty (just git init'd). We need to build the
     entire library and push the first commit to github.com/AkimZmerli/splitty.

     ---
     Architecture: Vertical Slice + Bounded Domains

     Single-import API: import "github.com/AkimZmerli/splitty" gets everything.

     Feature slices (each file is a complete vertical slice through the feature):


     ┌──────────────┬──────────────────────┬───────────────────────────────────────────────────┐
     │     File     │    Feature Slice     │                    Key Exports                    │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ split.go     │ Split/Close panes    │ Split(), Close(), ClosePane()                     │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ navigate.go  │ Focus navigation     │ Focus(), FocusPane(), FocusedPane(), Panes()      │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ resize.go    │ Resize + layout calc │ Resize(), calculateLayout()                       │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ zoom.go      │ Zoom/maximize + swap │ Zoom(), Unzoom(), IsZoomed(), Swap()              │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ broadcast.go │ Broadcast input      │ SetBroadcast(), IsBroadcasting(), SendInput()     │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ persist.go   │ Save/load layouts    │ SaveLayout(), LoadLayout()                        │
     ├──────────────┼──────────────────────┼───────────────────────────────────────────────────┤
     │ presets.go   │ Named layout presets │ PresetSingle, PresetDev, PresetTriple, PresetQuad │
     └──────────────┴──────────────────────┴───────────────────────────────────────────────────┘
     Core infrastructure (shared by all slices):
     ┌─────────────┬──────────────────────────────────────────────────────────────┐
     │    File     │                        Responsibility                        │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ splitty.go  │ Manager struct, New(), Init(), Update(), View() orchestrator │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ node.go     │ Binary tree: leafNode, splitNode (algorithmic core)          │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ pane.go     │ Pane type with PTY lifecycle                                 │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ messages.go │ All public tea.Msg types                                     │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ keys.go     │ KeyMap, DefaultKeyMap()                                      │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ theme.go    │ Theme + 5 built-in themes                                    │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ options.go  │ Functional Option type + all With* functions                 │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ logger.go   │ Optional charmbracelet/log wrapper                           │
     ├─────────────┼──────────────────────────────────────────────────────────────┤
     │ doc.go      │ Package-level godoc                                          │
     └─────────────┴──────────────────────────────────────────────────────────────┘
     Bounded domain (internal subpackage):
     File: terminal/screen.go
     Responsibility: Virtual terminal: cell grid + ANSI parser (using charmbracelet/x/ansi) +
     Render()
     ────────────────────────────────────────
     File: terminal/cell.go
     Responsibility: Cell and Style types
     ────────────────────────────────────────
     File: terminal/cursor.go
     Responsibility: Cursor state
     ────────────────────────────────────────
     File: terminal/pty.go
     Responsibility: PTY wrapper using creack/pty
     ---
     Key Design Decisions

     1. Binary tree for layout: leaves = panes, internal nodes = splits with direction + ratio
     2. Self-re-invoking tea.Cmd for PTY reads (idiomatic Bubble Tea, no raw goroutines)
     3. charmbracelet/x/ansi for ANSI parsing (stays in Charm ecosystem, lighter than midterm)
     4. sync.RWMutex on Screen for thread safety between Cmd goroutines and View()
     5. Lipgloss for all rendering/styling (borders, joining, themes)

     ---
     Implementation Order

     Slice 0: Foundation

     Files: go.mod, doc.go, messages.go, keys.go, theme.go, options.go, logger.go, node.go

     Slice 1: Terminal Emulation

     Files: terminal/cell.go, terminal/cursor.go, terminal/screen.go, terminal/pty.go

     Slice 2: Single Pane (first working UI)

     Files: pane.go, splitty.go (initial - single pane only)

     Slice 3: Split & Close

     Files: split.go, update splitty.go

     Slice 4: Navigate

     Files: navigate.go, update splitty.go

     Slice 5: Resize

     Files: resize.go, update splitty.go

     Slice 6: Zoom & Swap

     Files: zoom.go, update splitty.go

     Slice 7: Broadcast

     Files: broadcast.go, update splitty.go

     Slice 8: Persist & Presets

     Files: persist.go, presets.go

     Slice 9: Polish & Ship

     Files: README.md, ARCHITECTURE.md, LICENSE, Makefile, .github/workflows/ci.yml, examples/*,
     .gitignore

     ---
     Documentation Plan

     - README.md: Playful, engaging, with shield badges (Go Reference, Go Report Card, CI,
     License, Release). ASCII art header, demo GIF placeholder, quick start, feature list.
     Icons/badges allowed here.
     - ARCHITECTURE.md: Clean technical doc, NO icons/emojis. Binary tree model, terminal
     pipeline, message flow, file organization rationale.
     - doc.go: Package-level godoc with minimal example
     - Examples: 5 example programs (basic, custom-theme, custom-keys, embedded, presets)

     Sub-agents will be spawned to write README.md and ARCHITECTURE.md in parallel.

     ---
     Bead Epic Update

     Epic kitty-beads-hq5 already exists with 8 phases. We will update phase statuses via bd
     update in the kitty-beads repo as slices are completed:
     - hq5.1 (Core Library Setup) -> in_progress then closed after Slice 0-2
     - hq5.2 (Split Operations) -> after Slice 3
     - hq5.3 (PTY Management) -> covered by Slice 1-2
     - hq5.4 (Navigation & Theming) -> after Slice 4 + themes in Slice 0
     - hq5.5 (Resizing) -> after Slice 5
     - hq5.6 (Advanced Features) -> after Slice 6-8
     - hq5.7 (Polish & Release) -> after Slice 9
     - hq5.8 (kitty-beads Integration) -> deferred (separate PR in kitty-beads)

     ---
     Dependencies

     github.com/charmbracelet/bubbletea  v1.3.x
     github.com/charmbracelet/bubbles    v0.20.x
     github.com/charmbracelet/lipgloss   v1.x
     github.com/charmbracelet/log        v0.4.x
     github.com/charmbracelet/x/ansi     latest
     github.com/creack/pty               v1.1.x
     golang.org/x/term                   latest

     ---
     Remote Push

     After all slices are complete:
     git remote add origin https://github.com/AkimZmerli/splitty.git
     git branch -M main
     git add -A && git commit
     git push -u origin main

     ---
     Verification

     1. go build ./... compiles clean
     2. go vet ./... passes
     3. go test ./... -short passes (unit tests)
     4. go run examples/basic/main.go launches interactive terminal with working shell
     5. Split (Ctrl+), navigate (Ctrl+hjkl), close (Ctrl+w) all work
     6. Themes render correctly
     7. README renders properly on GitHub

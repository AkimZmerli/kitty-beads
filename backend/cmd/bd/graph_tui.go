package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/tui/commands/graphtui"
)

// runGraphTUI launches the Bubble Tea TUI for the dependency graph browser.
func runGraphTUI(rootID string) {
	ctx := rootCtx
	activeStore := store
	if activeStore == nil && dbPath != "" {
		roStore, err := sqlite.NewReadOnlyWithTimeout(ctx, dbPath, lockTimeout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening database: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = roStore.Close() }()
		activeStore = roStore
	}
	if activeStore == nil {
		fmt.Fprintf(os.Stderr, "Error: no database connection\n")
		os.Exit(1)
	}
	cfg := graphtui.Config{
		Ctx:          ctx,
		Store:        activeStore,
		DaemonClient: daemonClient,
		RootID:       rootID,
		Actor:        actor,
		DbPath:       dbPath,
		LockTimeout:  lockTimeout,
	}
	model := graphtui.New(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

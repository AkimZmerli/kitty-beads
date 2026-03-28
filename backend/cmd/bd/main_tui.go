package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/tui/commands/maintui"
)

// runBeadsTUI launches the Bubble Tea hub TUI with tabbed Issues/Graph/Status views.
func runBeadsTUI() {
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

	cfg := maintui.Config{
		Ctx:          ctx,
		Store:        activeStore,
		DaemonClient: daemonClient,
		Actor:        actor,
		DbPath:       dbPath,
		LockTimeout:  lockTimeout,
	}

	model := maintui.New(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}

	if m, ok := finalModel.(maintui.Model); ok && m.ClosedCount() > 0 {
		markDirtyAndScheduleFlush()
	}
}

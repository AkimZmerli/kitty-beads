package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/tui/commands/listtui"
	"github.com/steveyegge/beads/internal/types"
)

// runListTUI launches the Bubble Tea TUI for issue browsing.
func runListTUI(filter types.IssueFilter) {
	ctx := rootCtx

	// In daemon mode without a direct store, open a read-only store for
	// label fetching, detail views, and dependency lookups.
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

	cfg := listtui.Config{
		Ctx:          ctx,
		Store:        activeStore,
		DaemonClient: daemonClient,
		Filter:       filter,
		Actor:        actor,
		DbPath:       dbPath,
		LockTimeout:  lockTimeout,
	}

	model := listtui.New(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen(), tea.WithMouseCellMotion())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}

	// If any issues were closed, trigger auto-flush
	if m, ok := finalModel.(listtui.Model); ok && m.ClosedCount() > 0 {
		markDirtyAndScheduleFlush()
	}
}

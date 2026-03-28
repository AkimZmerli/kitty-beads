package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/tui/commands/daemontui"
)

func runDaemonMonitorTUI() {
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
	cfg := daemontui.Config{
		Ctx:          ctx,
		Store:        activeStore,
		DaemonClient: daemonClient,
		DbPath:       dbPath,
		LockTimeout:  lockTimeout,
	}
	model := daemontui.New(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

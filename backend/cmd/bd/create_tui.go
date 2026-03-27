package main

import (
	"encoding/json"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/tui/commands/createtui"
	"github.com/steveyegge/beads/internal/types"
)

// runCreateFormTUI launches the enhanced Bubble Tea create-form.
// It mirrors the pattern of runListTUI in list_tui.go.
func runCreateFormTUI() {
	ctx := rootCtx

	// In daemon mode without a direct store, open a read-only store so the
	// dependency picker can enumerate open issues.
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

	cfg := createtui.Config{
		Ctx:          ctx,
		Store:        activeStore,
		DaemonClient: daemonClient,
		Actor:        actor,
		DbPath:       dbPath,
		LockTimeout:  lockTimeout,
	}

	model := createtui.New(cfg)
	p := tea.NewProgram(model, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}

	m, ok := finalModel.(createtui.Model)
	if !ok || !m.Confirmed() {
		fmt.Fprintln(os.Stderr, "Issue creation canceled.")
		return
	}

	fv := m.ToFormValues()

	// Daemon path: use RPC
	if daemonClient != nil {
		createArgs := &rpc.CreateArgs{
			Title:        fv.Title,
			Description:  fv.Description,
			IssueType:    fv.IssueType,
			Priority:     fv.Priority,
			Assignee:     fv.Assignee,
			Labels:       fv.Labels,
			Dependencies: fv.Dependencies,
		}
		resp, err := daemonClient.Create(createArgs)
		if err != nil {
			FatalError("%v", err)
		}
		if jsonOutput {
			fmt.Println(string(resp.Data))
		} else {
			var issue types.Issue
			if err := json.Unmarshal(resp.Data, &issue); err != nil {
				FatalError("parsing response: %v", err)
			}
			printCreatedIssue(&issue)
		}
		return
	}

	// Direct path: use store
	if store == nil {
		FatalError("no database connection")
	}

	cfv := &createFormValues{
		Title:        fv.Title,
		Description:  fv.Description,
		IssueType:    fv.IssueType,
		Priority:     fv.Priority,
		Assignee:     fv.Assignee,
		Labels:       fv.Labels,
		Dependencies: fv.Dependencies,
	}

	issue, err := CreateIssueFromFormValues(rootCtx, store, cfv, actor)
	if err != nil {
		FatalError("%v", err)
	}

	markDirtyAndScheduleFlush()

	if jsonOutput {
		outputJSON(issue)
	} else {
		printCreatedIssue(issue)
	}
}

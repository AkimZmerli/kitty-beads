// Package listtui provides the top-level Bubble Tea model for `bd list --tui`.
package listtui

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/tui/components/issuedetail"
	"github.com/steveyegge/beads/internal/tui/components/issuelist"
	"github.com/steveyegge/beads/internal/tui/components/statusbar"
	"github.com/steveyegge/beads/internal/types"
)

type viewState int

const (
	stateLoading viewState = iota
	stateList
	stateDetail
)

// Config holds the external dependencies injected from the CLI layer.
type Config struct {
	Ctx          context.Context
	Store        storage.Storage
	DaemonClient *rpc.Client
	Filter       types.IssueFilter
	Actor        string
	DbPath       string
	LockTimeout  time.Duration
}

// Model is the top-level Bubble Tea model.
type Model struct {
	cfg       Config
	state     viewState
	list      issuelist.Model
	detail    issuedetail.Model
	statusbar statusbar.Model
	listKeys  common.ListKeys
	detailKeys common.DetailKeys
	width     int
	height    int
	closed    int // count of issues closed during this session
}

// New creates a new list TUI model.
func New(cfg Config) Model {
	listKeys := common.NewListKeys()
	detailKeys := common.NewDetailKeys()

	m := Model{
		cfg:        cfg,
		state:      stateLoading,
		list:       issuelist.New(0, 0),
		detail:     issuedetail.New(0, 0),
		statusbar:  statusbar.New(listKeys),
		listKeys:   listKeys,
		detailKeys: detailKeys,
	}
	return m
}

// ClosedCount returns how many issues were closed during this TUI session.
func (m Model) ClosedCount() int {
	return m.closed
}

// Init returns the initial command to load issues.
func (m Model) Init() tea.Cmd {
	return m.loadIssuesCmd()
}

// Update handles all messages and routes to the appropriate component.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		listH := m.height - 2 // leave room for status bar
		m.list.SetSize(m.width, listH)
		m.detail.SetSize(m.width, listH)
		m.statusbar.SetWidth(m.width)
		return m, nil

	case common.IssuesLoadedMsg:
		m.list.SetItems(msg.Issues, msg.Labels)
		m.state = stateList
		return m, nil

	case common.IssueDetailLoadedMsg:
		m.detail.SetContent(msg.Issue, msg.Labels, msg.Deps, msg.Dependents)
		m.state = stateDetail
		m.statusbar.SetKeyMap(m.detailKeys)
		return m, nil

	case common.IssueClosedMsg:
		m.closed++
		statusText := fmt.Sprintf("Closed %s", msg.ID)
		if len(msg.NewlyUnblocked) > 0 {
			ids := ""
			for i, u := range msg.NewlyUnblocked {
				if i > 0 {
					ids += ", "
				}
				ids += u.ID
			}
			statusText += fmt.Sprintf(" (unblocked: %s)", ids)
		}
		// Reload the list and show status message
		return m, tea.Batch(
			m.loadIssuesCmd(),
			func() tea.Msg { return common.StatusMsg{Text: statusText} },
		)

	case common.ErrMsg:
		return m, func() tea.Msg {
			return common.StatusMsg{Text: fmt.Sprintf("Error: %v", msg.Err)}
		}

	case common.StatusMsg:
		var cmd tea.Cmd
		m.statusbar, cmd = m.statusbar.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Delegate remaining messages to the active component
	return m.delegateUpdate(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, m.listKeys.Quit) {
		return m, tea.Quit
	}

	switch m.state {
	case stateList:
		// Don't intercept keys when the list is filtering
		if m.list.FilterState() == list.Filtering {
			return m.delegateUpdate(msg)
		}

		switch {
		case key.Matches(msg, m.listKeys.Enter):
			if item := m.list.SelectedItem(); item != nil {
				return m, m.loadDetailCmd(item.Issue.ID)
			}
		case key.Matches(msg, m.listKeys.Close):
			if item := m.list.SelectedItem(); item != nil {
				return m, m.closeIssueCmd(item.Issue.ID)
			}
		default:
			return m.delegateUpdate(msg)
		}

	case stateDetail:
		switch {
		case key.Matches(msg, m.detailKeys.Back):
			m.state = stateList
			m.statusbar.SetKeyMap(m.listKeys)
			return m, nil
		default:
			return m.delegateUpdate(msg)
		}
	}

	return m, nil
}

func (m Model) delegateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case stateList:
		m.list, cmd = m.list.Update(msg)
	case stateDetail:
		m.detail, cmd = m.detail.Update(msg)
	}
	return m, cmd
}

// View renders the current view.
func (m Model) View() string {
	if m.width == 0 {
		return "" // Wait for WindowSizeMsg
	}

	var content string
	switch m.state {
	case stateLoading:
		content = lipgloss.Place(
			m.width, m.height-2,
			lipgloss.Center, lipgloss.Center,
			common.MutedTextStyle.Render("Loading issues..."),
		)
	case stateList:
		content = m.list.View()
	case stateDetail:
		content = m.detail.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		m.statusbar.View(),
	)
}

// loadIssuesCmd returns a command that loads issues asynchronously.
func (m Model) loadIssuesCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx

		// Daemon path: use RPC
		if m.cfg.DaemonClient != nil {
			issues, err := m.loadIssuesViaRPC()
			if err != nil {
				return common.ErrMsg{Err: err}
			}
			// Batch-fetch labels via store (read-only)
			labels, err := m.loadLabelsForIssues(ctx, issues)
			if err != nil {
				// Non-fatal: show issues without labels
				return common.IssuesLoadedMsg{Issues: issues, Labels: nil}
			}
			return common.IssuesLoadedMsg{Issues: issues, Labels: labels}
		}

		// Direct path: use store
		if m.cfg.Store != nil {
			issues, err := m.cfg.Store.SearchIssues(ctx, "", m.cfg.Filter)
			if err != nil {
				return common.ErrMsg{Err: err}
			}
			labels, err := m.loadLabelsForIssues(ctx, issues)
			if err != nil {
				return common.IssuesLoadedMsg{Issues: issues, Labels: nil}
			}
			return common.IssuesLoadedMsg{Issues: issues, Labels: labels}
		}

		return common.ErrMsg{Err: fmt.Errorf("no storage available")}
	}
}

func (m Model) loadIssuesViaRPC() ([]*types.Issue, error) {
	f := m.cfg.Filter
	listArgs := &rpc.ListArgs{}

	if f.Status != nil {
		listArgs.Status = string(*f.Status)
	}
	if f.IssueType != nil {
		listArgs.IssueType = string(*f.IssueType)
	}
	if f.Assignee != nil {
		listArgs.Assignee = *f.Assignee
	}
	if f.Priority != nil {
		listArgs.Priority = f.Priority
	}
	if f.Limit > 0 {
		listArgs.Limit = f.Limit
	}
	if len(f.Labels) > 0 {
		listArgs.Labels = f.Labels
	}
	if len(f.LabelsAny) > 0 {
		listArgs.LabelsAny = f.LabelsAny
	}
	if len(f.IDs) > 0 {
		listArgs.IDs = f.IDs
	}
	if f.TitleSearch != "" {
		listArgs.Query = f.TitleSearch
	}
	listArgs.TitleContains = f.TitleContains
	listArgs.DescriptionContains = f.DescriptionContains
	listArgs.NotesContains = f.NotesContains
	if f.CreatedAfter != nil {
		listArgs.CreatedAfter = f.CreatedAfter.Format(time.RFC3339)
	}
	if f.CreatedBefore != nil {
		listArgs.CreatedBefore = f.CreatedBefore.Format(time.RFC3339)
	}
	if f.UpdatedAfter != nil {
		listArgs.UpdatedAfter = f.UpdatedAfter.Format(time.RFC3339)
	}
	if f.UpdatedBefore != nil {
		listArgs.UpdatedBefore = f.UpdatedBefore.Format(time.RFC3339)
	}
	if f.ClosedAfter != nil {
		listArgs.ClosedAfter = f.ClosedAfter.Format(time.RFC3339)
	}
	if f.ClosedBefore != nil {
		listArgs.ClosedBefore = f.ClosedBefore.Format(time.RFC3339)
	}
	listArgs.EmptyDescription = f.EmptyDescription
	listArgs.NoAssignee = f.NoAssignee
	listArgs.NoLabels = f.NoLabels
	listArgs.PriorityMin = f.PriorityMin
	listArgs.PriorityMax = f.PriorityMax
	listArgs.Pinned = f.Pinned
	if f.ParentID != nil {
		listArgs.ParentID = *f.ParentID
	}
	if len(f.ExcludeStatus) > 0 {
		for _, s := range f.ExcludeStatus {
			listArgs.ExcludeStatus = append(listArgs.ExcludeStatus, string(s))
		}
	}
	if len(f.ExcludeTypes) > 0 {
		for _, t := range f.ExcludeTypes {
			listArgs.ExcludeTypes = append(listArgs.ExcludeTypes, string(t))
		}
	}
	listArgs.Deferred = f.Deferred
	if f.DeferAfter != nil {
		listArgs.DeferAfter = f.DeferAfter.Format(time.RFC3339)
	}
	if f.DeferBefore != nil {
		listArgs.DeferBefore = f.DeferBefore.Format(time.RFC3339)
	}
	if f.DueAfter != nil {
		listArgs.DueAfter = f.DueAfter.Format(time.RFC3339)
	}
	if f.DueBefore != nil {
		listArgs.DueBefore = f.DueBefore.Format(time.RFC3339)
	}
	listArgs.Overdue = f.Overdue
	if f.Ephemeral != nil {
		listArgs.Ephemeral = f.Ephemeral
	}
	if f.MolType != nil {
		listArgs.MolType = string(*f.MolType)
	}

	resp, err := m.cfg.DaemonClient.List(listArgs)
	if err != nil {
		return nil, fmt.Errorf("daemon list: %w", err)
	}
	if !resp.Success {
		return nil, fmt.Errorf("daemon list failed: %s", resp.Error)
	}

	var issues []*types.Issue
	if err := json.Unmarshal(resp.Data, &issues); err != nil {
		return nil, fmt.Errorf("parsing list response: %w", err)
	}
	return issues, nil
}

func (m Model) loadLabelsForIssues(ctx context.Context, issues []*types.Issue) (map[string][]string, error) {
	s := m.cfg.Store
	if s == nil {
		return nil, nil
	}
	ids := make([]string, len(issues))
	for i, iss := range issues {
		ids[i] = iss.ID
	}
	return s.GetLabelsForIssues(ctx, ids)
}

// loadDetailCmd returns a command that loads a single issue's detail.
func (m Model) loadDetailCmd(id string) tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx
		s := m.cfg.Store
		if s == nil {
			return common.ErrMsg{Err: fmt.Errorf("detail view requires direct store access")}
		}

		issue, err := s.GetIssue(ctx, id)
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("loading issue %s: %w", id, err)}
		}

		labels, err := s.GetLabelsForIssues(ctx, []string{id})
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("loading labels for %s: %w", id, err)}
		}

		deps, err := s.GetDependenciesWithMetadata(ctx, id)
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("loading deps for %s: %w", id, err)}
		}

		dependents, err := s.GetDependentsWithMetadata(ctx, id)
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("loading dependents for %s: %w", id, err)}
		}

		return common.IssueDetailLoadedMsg{
			Issue:      issue,
			Labels:     labels[id],
			Deps:       deps,
			Dependents: dependents,
		}
	}
}

// closeIssueCmd returns a command that closes an issue.
func (m Model) closeIssueCmd(id string) tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx

		// Daemon path: use RPC
		if m.cfg.DaemonClient != nil {
			resp, err := m.cfg.DaemonClient.CloseIssue(&rpc.CloseArgs{
				ID:          id,
				SuggestNext: true,
			})
			if err != nil {
				return common.ErrMsg{Err: fmt.Errorf("closing %s: %w", id, err)}
			}
			if !resp.Success {
				return common.ErrMsg{Err: fmt.Errorf("close %s: %s", id, resp.Error)}
			}
			return common.IssueClosedMsg{ID: id}
		}

		// Direct path: use store
		s := m.cfg.Store
		if s == nil {
			return common.ErrMsg{Err: fmt.Errorf("close requires storage access")}
		}

		// Check if blocked
		blocked, blockers, err := s.IsBlocked(ctx, id)
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("checking blockers for %s: %w", id, err)}
		}
		if blocked && len(blockers) > 0 {
			return common.ErrMsg{Err: fmt.Errorf("cannot close %s: blocked by %v", id, blockers)}
		}

		session := os.Getenv("CLAUDE_SESSION_ID")
		if err := s.CloseIssue(ctx, id, "", m.cfg.Actor, session); err != nil {
			return common.ErrMsg{Err: fmt.Errorf("closing %s: %w", id, err)}
		}

		unblocked, err := s.GetNewlyUnblockedByClose(ctx, id)
		if err != nil {
			// Non-fatal: close succeeded, just no unblocked info
			return common.IssueClosedMsg{ID: id}
		}

		return common.IssueClosedMsg{ID: id, NewlyUnblocked: unblocked}
	}
}

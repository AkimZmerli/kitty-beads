// Package maintui provides the top-level "hub" Bubble Tea model for `bd tui --beads`.
// It renders a tab bar at the top and delegates to the Issues, Graph, or Status sub-models.
package maintui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/tui/commands/graphtui"
	"github.com/steveyegge/beads/internal/tui/commands/listtui"
	"github.com/steveyegge/beads/internal/tui/commands/statustui"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/types"
)

type activeTab int

const (
	tabIssues activeTab = iota
	tabGraph
	tabStatus
)

// Config holds the external dependencies injected from the CLI layer.
type Config struct {
	Ctx          context.Context
	Store        storage.Storage
	DaemonClient *rpc.Client
	Actor        string
	DbPath       string
	LockTimeout  time.Duration
}

// Model is the top-level hub Bubble Tea model.
type Model struct {
	cfg       Config
	activeTab activeTab
	list      listtui.Model
	graph     graphtui.Model
	status    statustui.Model
	width     int
	height    int
}

// New creates a new hub model wiring up all three sub-models.
func New(cfg Config) Model {
	listCfg := listtui.Config{
		Ctx:          cfg.Ctx,
		Store:        cfg.Store,
		DaemonClient: cfg.DaemonClient,
		Filter:       types.IssueFilter{},
		Actor:        cfg.Actor,
		DbPath:       cfg.DbPath,
		LockTimeout:  cfg.LockTimeout,
	}

	graphCfg := graphtui.Config{
		Ctx:          cfg.Ctx,
		Store:        cfg.Store,
		DaemonClient: cfg.DaemonClient,
		RootID:       "",
		Actor:        cfg.Actor,
		DbPath:       cfg.DbPath,
		LockTimeout:  cfg.LockTimeout,
	}

	statusCfg := statustui.Config{
		Ctx:          cfg.Ctx,
		Store:        cfg.Store,
		DaemonClient: cfg.DaemonClient,
		DbPath:       cfg.DbPath,
		LockTimeout:  cfg.LockTimeout,
	}

	return Model{
		cfg:       cfg,
		activeTab: tabIssues,
		list:      listtui.New(listCfg),
		graph:     graphtui.New(graphCfg),
		status:    statustui.New(statusCfg),
	}
}

// ClosedCount returns how many issues were closed during this session (from the Issues tab).
func (m Model) ClosedCount() int {
	return m.list.ClosedCount()
}

// Init returns a batch of all three sub-model Init commands.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.list.Init(),
		m.graph.Init(),
		m.status.Init(),
	)
}

// Update handles all messages, delegating to the active sub-model for most messages
// and forwarding WindowSizeMsg to all three.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Subtract 1 for the tab bar line
		subMsg := tea.WindowSizeMsg{Width: m.width, Height: m.height - 1}
		var cmd1, cmd2, cmd3 tea.Cmd
		var tm1, tm2, tm3 tea.Model
		tm1, cmd1 = m.list.Update(subMsg)
		tm2, cmd2 = m.graph.Update(subMsg)
		tm3, cmd3 = m.status.Update(subMsg)
		m.list = tm1.(listtui.Model)
		m.graph = tm2.(graphtui.Model)
		m.status = tm3.(statustui.Model)
		return m, tea.Batch(cmd1, cmd2, cmd3)

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Delegate all other messages to the active sub-model only.
	return m.delegateUpdate(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit

	case "1":
		m.activeTab = tabIssues
		return m, nil

	case "2":
		m.activeTab = tabGraph
		return m, nil

	case "3":
		m.activeTab = tabStatus
		return m, nil

	case "tab":
		m.activeTab = (m.activeTab + 1) % 3
		return m, nil
	}

	return m.delegateUpdate(msg)
}

func (m Model) delegateUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var tm tea.Model

	switch m.activeTab {
	case tabIssues:
		tm, cmd = m.list.Update(msg)
		m.list = tm.(listtui.Model)
	case tabGraph:
		tm, cmd = m.graph.Update(msg)
		m.graph = tm.(graphtui.Model)
	case tabStatus:
		tm, cmd = m.status.Update(msg)
		m.status = tm.(statustui.Model)
	}

	return m, cmd
}

// View renders the tab bar on the first line followed by the active sub-model view.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}
	return lipgloss.JoinVertical(lipgloss.Left,
		m.renderTabBar(),
		m.renderActiveTab(),
	)
}

// renderTabBar renders the single-line tab bar at the top.
func (m Model) renderTabBar() string {
	activeStyle := lipgloss.NewStyle().
		Foreground(common.NeonCyan).
		Bold(true).
		Underline(true)

	inactiveStyle := lipgloss.NewStyle().
		Foreground(common.TextMuted)

	hintStyle := lipgloss.NewStyle().
		Foreground(common.TextDark)

	type tabDef struct {
		label string
		tab   activeTab
	}

	tabs := []tabDef{
		{"[1] Issues", tabIssues},
		{"[2] Graph", tabGraph},
		{"[3] Status", tabStatus},
	}

	var tabParts []string
	for _, t := range tabs {
		if m.activeTab == t.tab {
			tabParts = append(tabParts, activeStyle.Render(t.label))
		} else {
			tabParts = append(tabParts, inactiveStyle.Render(t.label))
		}
	}

	hint := hintStyle.Render("q:quit  Tab:switch")

	left := "  " + tabParts[0] + "   " + tabParts[1] + "   " + tabParts[2]
	// Right-align the hint by padding in between
	gap := m.width - lipgloss.Width(left) - lipgloss.Width(hint) - 2
	if gap < 1 {
		gap = 1
	}
	return fmt.Sprintf("%s%s%s", left, lipgloss.NewStyle().Width(gap).Render(""), hint)
}

// renderActiveTab returns the view of the currently-active sub-model.
func (m Model) renderActiveTab() string {
	switch m.activeTab {
	case tabIssues:
		return m.list.View()
	case tabGraph:
		return m.graph.View()
	case tabStatus:
		return m.status.View()
	}
	return ""
}

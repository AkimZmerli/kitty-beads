// Package statustui provides the Bubble Tea model for `bd status --tui`.
package statustui

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/types"
)

const (
	barWidth        = 20
	refreshInterval = 5 * time.Second
)

// tickMsg is sent on each auto-refresh tick.
type tickMsg time.Time

// statsLoadedMsg carries freshly loaded statistics.
type statsLoadedMsg struct {
	stats *types.Statistics
}

// errMsg carries a load error.
type errMsg struct {
	err error
}

// Config holds the external dependencies injected from the CLI layer.
type Config struct {
	Ctx          context.Context
	Store        storage.Storage
	DaemonClient *rpc.Client
	DbPath       string
	LockTimeout  time.Duration
}

// Model is the Bubble Tea model for the status dashboard.
type Model struct {
	cfg    Config
	stats  *types.Statistics
	err    error
	width  int
	height int
}

// New creates a new status TUI model.
func New(cfg Config) Model {
	return Model{cfg: cfg}
}

// Init returns the initial load command and starts the refresh ticker.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadCmd(),
		tea.Tick(refreshInterval, func(t time.Time) tea.Msg { return tickMsg(t) }),
	)
}

// Update handles all incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case statsLoadedMsg:
		m.stats = msg.stats
		m.err = nil
		return m, nil

	case errMsg:
		m.err = msg.err
		return m, nil

	case tickMsg:
		return m, tea.Batch(
			m.loadCmd(),
			tea.Tick(refreshInterval, func(t time.Time) tea.Msg { return tickMsg(t) }),
		)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q", "ctrl+c":
			return m, tea.Quit
		case "r", "R":
			return m, m.loadCmd()
		}
	}

	return m, nil
}

// View renders the status dashboard.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var sb strings.Builder

	// Header
	headerStyle := lipgloss.NewStyle().
		Foreground(common.NeonCyan).
		Bold(true).
		Padding(0, 1)
	sb.WriteString(headerStyle.Render("Issue Database Status"))
	sb.WriteString("\n\n")

	if m.err != nil {
		errStyle := lipgloss.NewStyle().Foreground(common.NeonPink)
		sb.WriteString(errStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		sb.WriteString("\n\n")
		sb.WriteString(m.renderHelp())
		return sb.String()
	}

	if m.stats == nil {
		sb.WriteString(common.MutedTextStyle.Render("Loading..."))
		sb.WriteString("\n")
		return sb.String()
	}

	s := m.stats
	total := s.TotalIssues

	// Label column width (fixed, for alignment)
	labelStyle := lipgloss.NewStyle().
		Foreground(common.TextNormal).
		Width(14)
	countStyle := lipgloss.NewStyle().
		Foreground(common.TextBright).
		Width(5).
		Align(lipgloss.Right)
	mutedStyle := lipgloss.NewStyle().Foreground(common.TextMuted)

	// Section title
	sectionStyle := lipgloss.NewStyle().
		Foreground(common.TextMuted).
		Bold(true).
		MarginBottom(0)
	sb.WriteString(sectionStyle.Render("By Status"))
	sb.WriteString("\n")

	type row struct {
		label string
		value int
		color lipgloss.Color
	}

	rows := []row{
		{"Open", s.OpenIssues, common.NeonGreen},
		{"In Progress", s.InProgressIssues, common.NeonYellow},
		{"Blocked", s.BlockedIssues, common.NeonPink},
		{"Deferred", s.DeferredIssues, common.NeonCyan},
		{"Closed", s.ClosedIssues, common.TextMuted},
		{"Ready", s.ReadyIssues, common.NeonBlue},
	}

	for _, r := range rows {
		bar := renderBar(r.value, total, barWidth, r.color)
		line := fmt.Sprintf("  %s %s %s",
			labelStyle.Render(r.label),
			bar,
			countStyle.Render(fmt.Sprintf("%d", r.value)),
		)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Totals line
	sb.WriteString("\n")
	totalLine := fmt.Sprintf("  %s %s",
		mutedStyle.Render("Total:"),
		common.BrightTextStyle.Render(fmt.Sprintf("%d issues", total)),
	)
	sb.WriteString(totalLine)
	sb.WriteString("\n")

	// Extended stats (non-zero only)
	hasExtended := s.TombstoneIssues > 0 || s.AverageLeadTime > 0
	if hasExtended {
		sb.WriteString("\n")
		sb.WriteString(sectionStyle.Render("Extended"))
		sb.WriteString("\n")

		if s.TombstoneIssues > 0 {
			sb.WriteString(fmt.Sprintf("  %s %s\n",
				labelStyle.Render("Tombstones"),
				countStyle.Render(fmt.Sprintf("%d", s.TombstoneIssues)),
			))
		}
		if s.AverageLeadTime > 0 {
			leadStyle := lipgloss.NewStyle().Foreground(common.NeonMagenta)
			sb.WriteString(fmt.Sprintf("  %s %s\n",
				labelStyle.Render("Avg Lead Time"),
				leadStyle.Render(fmt.Sprintf("%.1fh", s.AverageLeadTime)),
			))
		}
	}

	// Alerts section
	hasAlerts := s.EpicsEligibleForClosure > 0 || s.PinnedIssues > 0
	if hasAlerts {
		sb.WriteString("\n")
		alertHeaderStyle := lipgloss.NewStyle().
			Foreground(common.NeonOrange).
			Bold(true)
		sb.WriteString(alertHeaderStyle.Render("Alerts"))
		sb.WriteString("\n")

		alertStyle := lipgloss.NewStyle().Foreground(common.NeonYellow)

		if s.EpicsEligibleForClosure > 0 {
			sb.WriteString(fmt.Sprintf("  %s %s\n",
				alertStyle.Render("Epics ready to close:"),
				common.BrightTextStyle.Render(fmt.Sprintf("%d", s.EpicsEligibleForClosure)),
			))
		}
		if s.PinnedIssues > 0 {
			pinnedStyle := lipgloss.NewStyle().Foreground(common.NeonMagenta)
			sb.WriteString(fmt.Sprintf("  %s %s\n",
				pinnedStyle.Render("Pinned issues:"),
				common.BrightTextStyle.Render(fmt.Sprintf("%d", s.PinnedIssues)),
			))
		}
	}

	// Footer / help
	sb.WriteString("\n")
	sb.WriteString(m.renderHelp())

	return sb.String()
}

// renderHelp renders the keybinding hint line.
func (m Model) renderHelp() string {
	sep := common.HelpSepStyle.Render(" • ")
	r := common.HelpKeyStyle.Render("r") + common.HelpDescStyle.Render(" refresh")
	q := common.HelpKeyStyle.Render("q") + common.HelpDescStyle.Render(" quit")
	auto := common.HelpDescStyle.Render(fmt.Sprintf("(auto-refresh every %ds)", int(refreshInterval.Seconds())))
	return "  " + r + sep + q + "  " + auto
}

// renderBar produces a filled/empty ASCII progress bar coloured with lipgloss.
func renderBar(value, total, width int, color lipgloss.Color) string {
	if total == 0 {
		return lipgloss.NewStyle().Foreground(common.TextDark).Render(strings.Repeat("░", width))
	}
	filled := value * width / total
	if filled > width {
		filled = width
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return lipgloss.NewStyle().Foreground(color).Render(bar)
}

// loadCmd returns a tea.Cmd that fetches statistics asynchronously.
func (m Model) loadCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx

		if m.cfg.DaemonClient != nil {
			resp, err := m.cfg.DaemonClient.Stats()
			if err != nil {
				return errMsg{err: fmt.Errorf("daemon stats: %w", err)}
			}
			if !resp.Success {
				return errMsg{err: fmt.Errorf("daemon stats failed: %s", resp.Error)}
			}
			var stats types.Statistics
			if err := json.Unmarshal(resp.Data, &stats); err != nil {
				return errMsg{err: fmt.Errorf("parsing stats response: %w", err)}
			}
			return statsLoadedMsg{stats: &stats}
		}

		if m.cfg.Store != nil {
			stats, err := m.cfg.Store.GetStatistics(ctx)
			if err != nil {
				return errMsg{err: fmt.Errorf("getting statistics: %w", err)}
			}
			return statsLoadedMsg{stats: stats}
		}

		return errMsg{err: fmt.Errorf("no storage available")}
	}
}

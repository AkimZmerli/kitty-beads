// Package daemontui provides the Bubble Tea model for `bd daemon --monitor`.
package daemontui

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/types"
)

// refreshMsg is sent by the auto-refresh ticker.
type refreshMsg struct{}

// statsLoadedMsg carries the result of a stats fetch.
type statsLoadedMsg struct {
	health *rpc.HealthResponse
	stats  *types.Statistics
	err    error
}

// Config holds the external dependencies injected from the CLI layer.
type Config struct {
	Ctx          context.Context
	Store        storage.Storage
	DaemonClient *rpc.Client
	DbPath       string
	LockTimeout  time.Duration
}

// Model is the top-level Bubble Tea model for the daemon monitor TUI.
type Model struct {
	cfg     Config
	width   int
	height  int
	health  *rpc.HealthResponse
	stats   *types.Statistics
	lastErr error
	loading bool
}

// New creates a new daemon monitor TUI model.
func New(cfg Config) Model {
	return Model{
		cfg:     cfg,
		loading: true,
	}
}

// Init starts the first data fetch and the auto-refresh ticker.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.fetchStatsCmd(),
		tickCmd(),
	)
}

func tickCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return refreshMsg{}
	})
}

// Update handles all incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case refreshMsg:
		return m, tea.Batch(m.fetchStatsCmd(), tickCmd())

	case statsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.lastErr = msg.err
			m.health = nil
			m.stats = msg.stats // may still have store stats even on daemon error
		} else {
			m.lastErr = nil
			m.health = msg.health
			m.stats = msg.stats
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, m.fetchStatsCmd()
		}
	}

	return m, nil
}

// View renders the monitor panel.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	title := renderTitle()
	daemonPanel := m.renderDaemonPanel()
	statsPanel := m.renderStatsPanel()
	helpBar := renderHelpBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		daemonPanel,
		"",
		statsPanel,
		"",
		helpBar,
	)
}

// -- rendering helpers -------------------------------------------------------

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(common.NeonCyan).
			Bold(true).
			PaddingLeft(1)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(common.NightBorder).
			Padding(0, 2).
			MarginLeft(1)

	panelHeaderStyle = lipgloss.NewStyle().
				Foreground(common.TextBright).
				Bold(true).
				MarginBottom(1)

	labelStyle = lipgloss.NewStyle().
			Foreground(common.TextMuted).
			Width(22)

	valueStyle = lipgloss.NewStyle().
			Foreground(common.TextNormal)

	goodStyle = lipgloss.NewStyle().
			Foreground(common.NeonGreen)

	warnStyle = lipgloss.NewStyle().
			Foreground(common.NeonYellow)

	badStyle = lipgloss.NewStyle().
			Foreground(common.NeonPink)
)

func renderTitle() string {
	return titleStyle.Render("bd daemon --monitor")
}

func (m Model) renderDaemonPanel() string {
	var rows []string
	rows = append(rows, panelHeaderStyle.Render("Daemon"))

	if m.loading && m.health == nil {
		rows = append(rows, common.MutedTextStyle.Render("  Fetching daemon status..."))
		return panelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
	}

	if m.cfg.DaemonClient == nil {
		rows = append(rows, row("Status", badStyle.Render("not running")))
		rows = append(rows, row("RPC", common.MutedTextStyle.Render("unavailable")))
	} else if m.health != nil {
		statusStr := formatHealthStatus(m.health.Status)
		rows = append(rows, row("Status", statusStr))
		rows = append(rows, row("Version", valueStyle.Render(m.health.Version)))
		rows = append(rows, row("Uptime", valueStyle.Render(formatUptime(m.health.Uptime))))
		rows = append(rows, row("Memory", valueStyle.Render(fmt.Sprintf("%d MB", m.health.MemoryAllocMB))))
		rows = append(rows, row("Active connections", valueStyle.Render(fmt.Sprintf("%d / %d", m.health.ActiveConns, m.health.MaxConns))))
		rows = append(rows, row("DB response", valueStyle.Render(fmt.Sprintf("%.1f ms", m.health.DBResponseTime))))
	} else if m.lastErr != nil {
		rows = append(rows, row("Status", badStyle.Render("error")))
		rows = append(rows, row("Error", warnStyle.Render(truncate(m.lastErr.Error(), 60))))
	}

	return panelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m Model) renderStatsPanel() string {
	var rows []string
	rows = append(rows, panelHeaderStyle.Render("Issue Statistics"))

	if m.loading && m.stats == nil {
		rows = append(rows, common.MutedTextStyle.Render("  Fetching statistics..."))
		return panelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
	}

	if m.stats == nil {
		rows = append(rows, common.MutedTextStyle.Render("  No statistics available"))
		return panelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
	}

	s := m.stats
	rows = append(rows, row("Total", valueStyle.Render(fmt.Sprintf("%d", s.TotalIssues))))
	rows = append(rows, row("Open", goodStyle.Render(fmt.Sprintf("%d", s.OpenIssues))))
	rows = append(rows, row("In Progress", warnStyle.Render(fmt.Sprintf("%d", s.InProgressIssues))))
	rows = append(rows, row("Blocked", badStyle.Render(fmt.Sprintf("%d", s.BlockedIssues))))
	rows = append(rows, row("Deferred", common.MutedTextStyle.Render(fmt.Sprintf("%d", s.DeferredIssues))))
	rows = append(rows, row("Ready", lipgloss.NewStyle().Foreground(common.NeonCyan).Render(fmt.Sprintf("%d", s.ReadyIssues))))
	rows = append(rows, row("Closed", common.MutedTextStyle.Render(fmt.Sprintf("%d", s.ClosedIssues))))
	if s.AverageLeadTime > 0 {
		rows = append(rows, row("Avg lead time", valueStyle.Render(fmt.Sprintf("%.1f h", s.AverageLeadTime))))
	}

	return panelStyle.Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func renderHelpBar() string {
	sep := common.HelpSepStyle.Render(" • ")
	parts := []string{
		common.HelpKeyStyle.Render("r") + common.HelpDescStyle.Render(" refresh"),
		common.HelpKeyStyle.Render("q") + common.HelpDescStyle.Render(" quit"),
	}
	bar := ""
	for i, p := range parts {
		if i > 0 {
			bar += sep
		}
		bar += p
	}
	return lipgloss.NewStyle().PaddingLeft(2).Render(bar)
}

// row renders a label/value pair.
func row(label, value string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		labelStyle.Render(label+":"),
		value,
	)
}

func formatHealthStatus(status string) string {
	switch status {
	case "healthy":
		return goodStyle.Render("healthy")
	case "degraded":
		return warnStyle.Render("degraded")
	case "unhealthy":
		return badStyle.Render("unhealthy")
	default:
		return common.MutedTextStyle.Render(status)
	}
}

func formatUptime(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.0fs", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%.0fm %.0fs", math.Floor(seconds/60), math.Mod(seconds, 60))
	}
	h := math.Floor(seconds / 3600)
	m := math.Floor(math.Mod(seconds, 3600) / 60)
	return fmt.Sprintf("%.0fh %.0fm", h, m)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// -- commands ----------------------------------------------------------------

func (m Model) fetchStatsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx
		result := statsLoadedMsg{}

		// Try daemon health if client is available.
		if m.cfg.DaemonClient != nil {
			health, err := m.cfg.DaemonClient.Health()
			if err != nil {
				result.err = fmt.Errorf("health check: %w", err)
			} else {
				result.health = health
			}

			// Fetch stats via daemon RPC.
			statsResp, statsErr := m.cfg.DaemonClient.Stats()
			if statsErr == nil && statsResp != nil {
				var stats types.Statistics
				if jsonErr := json.Unmarshal(statsResp.Data, &stats); jsonErr == nil {
					result.stats = &stats
				}
			}
			// If we got health data, return it even without stats.
			if result.health != nil || result.err != nil {
				return result
			}
		}

		// Fall back to direct store access.
		if m.cfg.Store != nil {
			stats, err := m.cfg.Store.GetStatistics(ctx)
			if err != nil {
				if result.err == nil {
					result.err = fmt.Errorf("store stats: %w", err)
				}
			} else {
				result.stats = stats
			}
		}

		return result
	}
}

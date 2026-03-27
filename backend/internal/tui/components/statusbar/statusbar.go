// Package statusbar provides a bottom help/status bar for the TUI.
package statusbar

import (
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/tui/common"
)

const statusTimeout = 3 * time.Second

type clearStatusMsg struct{}

// Model is the status bar component.
type Model struct {
	help      help.Model
	keyMap    help.KeyMap
	statusMsg string
	width     int
}

// New creates a new status bar.
func New(km help.KeyMap) Model {
	h := help.New()
	h.Styles.ShortKey = common.HelpKeyStyle
	h.Styles.ShortDesc = common.HelpDescStyle
	h.Styles.ShortSeparator = common.HelpSepStyle
	h.Styles.FullKey = common.HelpKeyStyle
	h.Styles.FullDesc = common.HelpDescStyle
	h.Styles.FullSeparator = common.HelpSepStyle
	return Model{
		help:   h,
		keyMap: km,
	}
}

// SetKeyMap changes the active keybindings displayed.
func (m *Model) SetKeyMap(km help.KeyMap) {
	m.keyMap = km
}

// SetWidth sets the available width for the status bar.
func (m *Model) SetWidth(w int) {
	m.width = w
	m.help.Width = w
}

// Update handles messages for the status bar.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case common.StatusMsg:
		m.statusMsg = msg.Text
		return m, tea.Tick(statusTimeout, func(time.Time) tea.Msg {
			return clearStatusMsg{}
		})
	case clearStatusMsg:
		m.statusMsg = ""
	case tea.KeyMsg:
		if key.Matches(msg, key.NewBinding(key.WithKeys("?"))) {
			m.help.ShowAll = !m.help.ShowAll
		}
	}
	return m, nil
}

// View renders the status bar.
func (m Model) View() string {
	if m.statusMsg != "" {
		style := lipgloss.NewStyle().
			Foreground(common.NeonGreen).
			Width(m.width)
		return style.Render(m.statusMsg)
	}
	return m.help.View(m.keyMap)
}

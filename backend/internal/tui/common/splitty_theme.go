package common

import (
	"github.com/AkimZmerli/splitty"
	"github.com/charmbracelet/lipgloss"
)

// KittyBeadsTheme is a splitty Theme using the Tokyo Night colors
// from the kitty-beads TUI palette.
var KittyBeadsTheme = splitty.Theme{
	BorderActive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonBlue),
	BorderInactive: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NightBorder),
	BorderScrollback: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonYellow),
	BorderScrollbackFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonOrange),
	BorderResize: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonCyan),
	BorderCopyMode: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NeonGreen),
	BorderCopyModeFocused: lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#B4F9A0")),
	Divider: lipgloss.NewStyle().
		Foreground(NightBorderHi),
	DividerChar: "│",
	StatusBar: lipgloss.NewStyle().
		Background(NightBgDark).
		Foreground(TextNormal).
		Padding(0, 1),
	StatusText: lipgloss.NewStyle().
		Foreground(TextNormal),
	ZoomIndicator:      "◉ ZOOM",
	BroadcastIndicator: "◉ BROADCAST",
}

// Package common provides shared styles, keybindings, and messages for the TUI.
// Uses the Tokyo Night theme, independent from the Ayu theme in internal/ui.
package common

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

// Tokyo Night - Backgrounds
var (
	NightBg            = lipgloss.Color("#1a1b26")
	NightBgDark        = lipgloss.Color("#16161e")
	NightBgHighlight   = lipgloss.Color("#24283b")
	NightSurface       = lipgloss.Color("#1f2335")
	NightSurfaceBright = lipgloss.Color("#292e42")
)

// Tokyo Night - Borders
var (
	NightBorder   = lipgloss.Color("#3b4261")
	NightBorderHi = lipgloss.Color("#565f89")
)

// Tokyo Night - Neon Accents
var (
	NeonCyan    = lipgloss.Color("#7dcfff")
	NeonMagenta = lipgloss.Color("#bb9af7")
	NeonPink    = lipgloss.Color("#f7768e")
	NeonGreen   = lipgloss.Color("#9ece6a")
	NeonOrange  = lipgloss.Color("#ff9e64")
	NeonYellow  = lipgloss.Color("#e0af68")
	NeonBlue    = lipgloss.Color("#7aa2f7")
)

// Tokyo Night - Text
var (
	TextBright = lipgloss.Color("#c0caf5")
	TextNormal = lipgloss.Color("#a9b1d6")
	TextMuted  = lipgloss.Color("#565f89")
	TextDark   = lipgloss.Color("#414868")
)

// Status icon styles
var (
	StatusOpenStyle       = lipgloss.NewStyle().Foreground(NeonGreen)
	StatusInProgressStyle = lipgloss.NewStyle().Foreground(NeonYellow)
	StatusBlockedStyle    = lipgloss.NewStyle().Foreground(NeonPink)
	StatusClosedStyle     = lipgloss.NewStyle().Foreground(TextMuted)
	StatusDeferredStyle   = lipgloss.NewStyle().Foreground(NeonCyan)
	StatusPinnedStyle     = lipgloss.NewStyle().Foreground(NeonMagenta)
	StatusHookedStyle     = lipgloss.NewStyle().Foreground(NeonCyan)
)

// Priority styles
var (
	PriorityP0Style = lipgloss.NewStyle().Foreground(NeonPink).Bold(true)
	PriorityP1Style = lipgloss.NewStyle().Foreground(NeonOrange)
	PriorityP2Style = lipgloss.NewStyle().Foreground(TextNormal)
	PriorityP3Style = lipgloss.NewStyle().Foreground(TextMuted)
	PriorityP4Style = lipgloss.NewStyle().Foreground(TextDark)
)

// Type badge styles
var (
	TypeEpicStyle    = lipgloss.NewStyle().Foreground(NeonMagenta)
	TypeBugStyle     = lipgloss.NewStyle().Foreground(NeonPink)
	TypeFeatureStyle = lipgloss.NewStyle().Foreground(NeonGreen)
	TypeTaskStyle    = lipgloss.NewStyle().Foreground(TextNormal)
	TypeChoreStyle   = lipgloss.NewStyle().Foreground(TextMuted)
)

// Panel styles
var (
	Panel = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(NightBorder).
		Padding(1, 2)

	PanelHeader = lipgloss.NewStyle().
			Foreground(TextBright).
			Bold(true)
)

// Selection styles
var (
	Selected = lipgloss.NewStyle().
			Background(NightBgHighlight).
			Foreground(NeonCyan)

	Cursor = lipgloss.NewStyle().
		Foreground(NeonMagenta).
		Bold(true)
)

// Help bar styles
var (
	HelpKeyStyle  = lipgloss.NewStyle().Foreground(NeonMagenta)
	HelpDescStyle = lipgloss.NewStyle().Foreground(TextMuted)
	HelpSepStyle  = lipgloss.NewStyle().Foreground(NightBorder)
)

// Filter input styles
var (
	FilterPromptStyle = lipgloss.NewStyle().Foreground(NeonCyan)
	FilterTextStyle   = lipgloss.NewStyle().Foreground(TextBright)
	FilterCursorStyle = lipgloss.NewStyle().Foreground(NeonMagenta)
)

// Text styles
var (
	NormalTextStyle = lipgloss.NewStyle().Foreground(TextNormal)
	BrightTextStyle = lipgloss.NewStyle().Foreground(TextBright)
	MutedTextStyle  = lipgloss.NewStyle().Foreground(TextMuted)
	DarkTextStyle   = lipgloss.NewStyle().Foreground(TextDark)
)

// Status icons - consistent with internal/ui/styles.go
const (
	StatusIconOpen       = "○"
	StatusIconInProgress = "◐"
	StatusIconBlocked    = "●"
	StatusIconClosed     = "✓"
	StatusIconDeferred   = "❄"
	StatusIconPinned     = "📌"
)

// RenderStatusIcon returns the status icon with Tokyo Night styling.
func RenderStatusIcon(status string) string {
	switch status {
	case "open":
		return StatusOpenStyle.Render(StatusIconOpen)
	case "in_progress":
		return StatusInProgressStyle.Render(StatusIconInProgress)
	case "blocked":
		return StatusBlockedStyle.Render(StatusIconBlocked)
	case "closed":
		return StatusClosedStyle.Render(StatusIconClosed)
	case "deferred":
		return StatusDeferredStyle.Render(StatusIconDeferred)
	case "pinned":
		return StatusPinnedStyle.Render(StatusIconPinned)
	case "hooked":
		return StatusHookedStyle.Render(StatusIconPinned)
	default:
		return "?"
	}
}

// RenderPriority returns priority label with Tokyo Night styling.
func RenderPriority(priority int) string {
	label := fmt.Sprintf("P%d", priority)
	switch priority {
	case 0:
		return PriorityP0Style.Render(label)
	case 1:
		return PriorityP1Style.Render(label)
	case 2:
		return PriorityP2Style.Render(label)
	case 3:
		return PriorityP3Style.Render(label)
	case 4:
		return PriorityP4Style.Render(label)
	default:
		return label
	}
}

// RenderType returns issue type with Tokyo Night styling.
func RenderType(issueType string) string {
	switch issueType {
	case "epic":
		return TypeEpicStyle.Render(issueType)
	case "bug":
		return TypeBugStyle.Render(issueType)
	case "feature":
		return TypeFeatureStyle.Render(issueType)
	case "task":
		return TypeTaskStyle.Render(issueType)
	case "chore":
		return TypeChoreStyle.Render(issueType)
	default:
		return issueType
	}
}

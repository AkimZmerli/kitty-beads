// Package issuelist provides a reusable issue list component with fuzzy filtering.
package issuelist

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/types"
)

// Item wraps a types.Issue for the list.Model.
type Item struct {
	Issue  *types.Issue
	Labels []string
}

// Title returns the rendered title line for the list delegate.
func (i Item) Title() string {
	icon := common.RenderStatusIcon(string(i.Issue.Status))
	prio := common.RenderPriority(i.Issue.Priority)
	id := common.MutedTextStyle.Render(i.Issue.ID)
	title := i.Issue.Title
	if i.Issue.Status == types.StatusClosed {
		title = common.MutedTextStyle.Render(title)
	}
	return fmt.Sprintf("%s %s  %s  %s", icon, prio, id, title)
}

// Description returns the type badge and optional labels.
func (i Item) Description() string {
	typeBadge := common.RenderType(string(i.Issue.IssueType))
	if len(i.Labels) > 0 {
		labelStr := ""
		for idx, l := range i.Labels {
			if idx > 0 {
				labelStr += " "
			}
			labelStr += common.DarkTextStyle.Render("#" + l)
		}
		return fmt.Sprintf("  %s  %s", typeBadge, labelStr)
	}
	return fmt.Sprintf("  %s", typeBadge)
}

// FilterValue returns the string used for fuzzy filtering.
func (i Item) FilterValue() string {
	return i.Issue.ID + " " + i.Issue.Title + " " + string(i.Issue.IssueType)
}

// Model wraps bubbles/list.Model with issue-specific helpers.
type Model struct {
	list list.Model
}

// New creates a new issue list component.
func New(width, height int) Model {
	delegate := list.NewDefaultDelegate()

	// Apply Tokyo Night selection styling
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		Foreground(common.NeonCyan).
		BorderLeftForeground(common.NeonMagenta)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		Foreground(common.TextMuted).
		BorderLeftForeground(common.NeonMagenta)
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.
		Foreground(common.TextNormal)
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.
		Foreground(common.TextDark)
	delegate.Styles.DimmedTitle = delegate.Styles.DimmedTitle.
		Foreground(common.TextMuted)
	delegate.Styles.DimmedDesc = delegate.Styles.DimmedDesc.
		Foreground(common.TextDark)

	l := list.New(nil, delegate, width, height)
	l.Title = "Issues"
	l.SetShowHelp(false) // We use our own status bar
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	// Style the list title and filter
	l.Styles.Title = lipgloss.NewStyle().
		Foreground(common.TextBright).
		Bold(true).
		Padding(0, 1)
	l.Styles.FilterPrompt = common.FilterPromptStyle
	l.Styles.FilterCursor = common.FilterCursorStyle

	return Model{list: l}
}

// SetItems populates the list from issues and their labels.
func (m *Model) SetItems(issues []*types.Issue, labels map[string][]string) {
	items := make([]list.Item, len(issues))
	for i, iss := range issues {
		var issLabels []string
		if labels != nil {
			issLabels = labels[iss.ID]
		}
		items[i] = Item{Issue: iss, Labels: issLabels}
	}
	m.list.SetItems(items)
}

// SelectedItem returns the currently selected Item, or nil.
func (m *Model) SelectedItem() *Item {
	sel := m.list.SelectedItem()
	if sel == nil {
		return nil
	}
	item, ok := sel.(Item)
	if !ok {
		return nil
	}
	return &item
}

// SetSize updates the list dimensions.
func (m *Model) SetSize(w, h int) {
	m.list.SetSize(w, h)
}

// FilterState returns the list's current filter state.
func (m *Model) FilterState() list.FilterState {
	return m.list.FilterState()
}

// Update delegates update to the inner list.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders the list.
func (m Model) View() string {
	return m.list.View()
}

// Package issuedetail provides a scrollable issue detail view.
package issuedetail

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/types"
)

// Model is the issue detail view component.
type Model struct {
	viewport viewport.Model
	issue    *types.Issue
	ready    bool
	width    int
	height   int
}

// New creates a new issue detail component.
func New(width, height int) Model {
	vp := viewport.New(width, height)
	vp.Style = lipgloss.NewStyle().Foreground(common.TextNormal)
	return Model{
		viewport: vp,
		width:    width,
		height:   height,
	}
}

// SetContent builds and sets the rendered detail content.
func (m *Model) SetContent(issue *types.Issue, labels []string, deps, dependents []*types.IssueWithDependencyMetadata) {
	m.issue = issue
	m.ready = true

	var b strings.Builder

	// Header
	icon := common.RenderStatusIcon(string(issue.Status))
	prio := common.RenderPriority(issue.Priority)
	typeBadge := common.RenderType(string(issue.IssueType))
	id := common.MutedTextStyle.Render(issue.ID)
	title := common.BrightTextStyle.Render(issue.Title)

	b.WriteString(fmt.Sprintf("%s %s  %s  %s\n", icon, prio, typeBadge, id))
	b.WriteString(title)
	b.WriteString("\n")
	b.WriteString(common.MutedTextStyle.Render(strings.Repeat("─", min(m.width-4, 60))))
	b.WriteString("\n\n")

	// Metadata
	metaStyle := common.NormalTextStyle
	labelStyle := common.HelpKeyStyle

	b.WriteString(metaStyle.Render("Status: "))
	b.WriteString(common.RenderStatusIcon(string(issue.Status)))
	b.WriteString(" " + metaStyle.Render(string(issue.Status)))
	b.WriteString("\n")

	if issue.Assignee != "" {
		b.WriteString(metaStyle.Render("Assignee: "))
		b.WriteString(metaStyle.Render(issue.Assignee))
		b.WriteString("\n")
	}

	if len(labels) > 0 {
		b.WriteString(metaStyle.Render("Labels: "))
		for i, l := range labels {
			if i > 0 {
				b.WriteString(common.MutedTextStyle.Render(", "))
			}
			b.WriteString(labelStyle.Render(l))
		}
		b.WriteString("\n")
	}

	if issue.DueAt != nil {
		b.WriteString(metaStyle.Render("Due: "))
		b.WriteString(metaStyle.Render(issue.DueAt.Format("2006-01-02")))
		b.WriteString("\n")
	}

	b.WriteString("\n")

	// Description via glamour markdown rendering
	if issue.Description != "" {
		b.WriteString(common.PanelHeader.Render("Description"))
		b.WriteString("\n")
		rendered, err := renderMarkdown(issue.Description, m.width-4)
		if err != nil {
			b.WriteString(issue.Description)
		} else {
			b.WriteString(rendered)
		}
		b.WriteString("\n")
	}

	// Dependencies
	if len(deps) > 0 {
		b.WriteString(common.PanelHeader.Render("Dependencies"))
		b.WriteString("\n")
		for _, d := range deps {
			depIcon := common.RenderStatusIcon(string(d.Issue.Status))
			depID := common.MutedTextStyle.Render(d.Issue.ID)
			depType := common.DarkTextStyle.Render(string(d.DependencyType))
			b.WriteString(fmt.Sprintf("  %s %s  %s  %s\n", depIcon, depID, d.Issue.Title, depType))
		}
		b.WriteString("\n")
	}

	// Dependents
	if len(dependents) > 0 {
		b.WriteString(common.PanelHeader.Render("Dependents"))
		b.WriteString("\n")
		for _, d := range dependents {
			depIcon := common.RenderStatusIcon(string(d.Issue.Status))
			depID := common.MutedTextStyle.Render(d.Issue.ID)
			depType := common.DarkTextStyle.Render(string(d.DependencyType))
			b.WriteString(fmt.Sprintf("  %s %s  %s  %s\n", depIcon, depID, d.Issue.Title, depType))
		}
		b.WriteString("\n")
	}

	m.viewport.SetContent(b.String())
	m.viewport.GotoTop()
}

// SetSize updates the viewport dimensions.
func (m *Model) SetSize(w, h int) {
	m.width = w
	m.height = h
	m.viewport.Width = w
	m.viewport.Height = h
}

// Update delegates update to the viewport.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

// View renders the detail view with scroll percentage footer.
func (m Model) View() string {
	if !m.ready {
		return common.MutedTextStyle.Render("Loading...")
	}
	footer := common.MutedTextStyle.Render(
		fmt.Sprintf(" %3.f%%", m.viewport.ScrollPercent()*100),
	)
	return m.viewport.View() + "\n" + footer
}

func renderMarkdown(md string, width int) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(width),
	)
	if err != nil {
		return "", err
	}
	return r.Render(md)
}

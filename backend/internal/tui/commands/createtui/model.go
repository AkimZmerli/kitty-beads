// Package createtui provides the Bubble Tea TUI for the enhanced create-form command.
// It adds template selection, a dependency picker with live search, and a preview step
// on top of the basic field editing experience.
package createtui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/types"
)

// step represents which phase of the form the user is on.
type step int

const (
	stepTemplate step = iota // pick a template
	stepBasic                // title, type, priority
	stepDetails              // description, labels, assignee
	stepDeps                 // dependency picker
	stepPreview              // preview + confirm
)

// Template holds predefined starter values for common issue kinds.
type Template struct {
	Name, TitlePrefix, Description, IssueType string
	Priority                                   int
}

var templates = []Template{
	{"Bug Report", "[Bug] ", "## Summary\n\n## Steps to Reproduce\n\n## Expected Behavior\n\n## Actual Behavior\n", "bug", 1},
	{"Feature Request", "[Feature] ", "## Summary\n\n## Use Case\n\n## Proposed Solution\n", "feature", 2},
	{"Task", "", "", "task", 2},
	{"Epic", "", "## Goal\n\n## Success Criteria\n", "epic", 1},
	{"Chore", "", "", "chore", 3},
	{"Custom", "", "", "task", 2},
}

var issueTypes = []string{"task", "bug", "feature", "epic", "chore"}

var priorityLabels = []string{
	"P0 - Critical",
	"P1 - High",
	"P2 - Medium",
	"P3 - Low",
	"P4 - Backlog",
}

// FormValues holds the final values collected by the TUI, ready for issue creation.
type FormValues struct {
	Title        string
	Description  string
	IssueType    string
	Priority     int
	Assignee     string
	Labels       []string
	Dependencies []string // "type:id" format
}

// Config holds dependencies injected from the CLI layer.
type Config struct {
	Ctx          context.Context
	Store        storage.Storage
	DaemonClient *rpc.Client
	Actor        string
	DbPath       string
	LockTimeout  time.Duration
}

// issuesLoadedMsg is sent when background issue loading completes.
type issuesLoadedMsg struct {
	issues []*types.Issue
	err    error
}

// basicField identifies which field in the basic-info step is active.
type basicField int

const (
	fieldTitle    basicField = iota
	fieldType                // cycle through issueTypes
	fieldPriority            // 0-4, adjusted with +/-
	basicFieldCount
)

// detailField identifies which field in the details step is active.
type detailField int

const (
	fieldDesc     detailField = iota
	fieldLabels
	fieldAssignee
	detailFieldCount
)

// Model is the top-level Bubble Tea model for the create-form TUI.
type Model struct {
	cfg    Config
	step   step
	width  int
	height int

	// stepTemplate
	templateCursor int

	// stepBasic
	activeBasic   basicField
	titleInput    textinput.Model
	typeCursor    int  // index into issueTypes
	priority      int  // 0-4

	// stepDetails
	activeDetail  detailField
	descInput     textinput.Model
	labelsInput   textinput.Model
	assigneeInput textinput.Model

	// stepDeps
	allIssues      []*types.Issue
	searchInput    textinput.Model
	filtered       []*types.Issue
	depsCursor     int
	selectedDeps   map[string]string // issueID -> depType
	depsLoaded     bool

	// stepPreview / result
	confirmed bool
	quitting  bool
}

// New creates and initialises the create-form TUI model.
func New(cfg Config) Model {
	ti := textinput.New()
	ti.Placeholder = "Issue title (required)"
	ti.Focus()
	ti.CharLimit = 500

	di := textinput.New()
	di.Placeholder = "Description (markdown, optional)"
	di.CharLimit = 5000

	li := textinput.New()
	li.Placeholder = "Labels, comma-separated (optional)"
	li.CharLimit = 500

	ai := textinput.New()
	ai.Placeholder = "Assignee (optional)"
	ai.CharLimit = 200

	si := textinput.New()
	si.Placeholder = "Search issues..."
	si.CharLimit = 200

	return Model{
		cfg:           cfg,
		step:          stepTemplate,
		titleInput:    ti,
		descInput:     di,
		labelsInput:   li,
		assigneeInput: ai,
		searchInput:   si,
		selectedDeps:  make(map[string]string),
		typeCursor:    0, // "task"
		priority:      2, // P2 medium
	}
}

// Confirmed returns true if the user confirmed issue creation.
func (m Model) Confirmed() bool { return m.confirmed }

// ToFormValues converts the model state into a FormValues ready for issue creation.
func (m Model) ToFormValues() *FormValues {
	var labels []string
	for _, l := range strings.Split(m.labelsInput.Value(), ",") {
		l = strings.TrimSpace(l)
		if l != "" {
			labels = append(labels, l)
		}
	}

	var deps []string
	for id, dt := range m.selectedDeps {
		deps = append(deps, fmt.Sprintf("%s:%s", dt, id))
	}

	return &FormValues{
		Title:        m.titleInput.Value(),
		Description:  m.descInput.Value(),
		IssueType:    issueTypes[m.typeCursor],
		Priority:     m.priority,
		Assignee:     m.assigneeInput.Value(),
		Labels:       labels,
		Dependencies: deps,
	}
}

// Init returns the initial command (none needed; issues are loaded lazily).
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages and user input.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case issuesLoadedMsg:
		if msg.err == nil {
			m.allIssues = msg.issues
			m.filtered = msg.issues
		}
		m.depsLoaded = true
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Delegate to active text input
	return m.delegateInputUpdate(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "q":
		if m.step == stepPreview || m.step == stepTemplate {
			m.quitting = true
			return m, tea.Quit
		}
	}

	switch m.step {
	case stepTemplate:
		return m.handleTemplateKey(msg)
	case stepBasic:
		return m.handleBasicKey(msg)
	case stepDetails:
		return m.handleDetailsKey(msg)
	case stepDeps:
		return m.handleDepsKey(msg)
	case stepPreview:
		return m.handlePreviewKey(msg)
	}
	return m, nil
}

// --- stepTemplate ---

func (m Model) handleTemplateKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.templateCursor > 0 {
			m.templateCursor--
		}
	case "down", "j":
		if m.templateCursor < len(templates)-1 {
			m.templateCursor++
		}
	case "enter":
		t := templates[m.templateCursor]
		// Pre-fill title prefix and description from template
		if t.TitlePrefix != "" {
			m.titleInput.SetValue(t.TitlePrefix)
			// Position cursor at end
			m.titleInput.CursorEnd()
		}
		if t.Description != "" {
			m.descInput.SetValue(t.Description)
		}
		// Set type and priority from template
		for i, it := range issueTypes {
			if it == t.IssueType {
				m.typeCursor = i
				break
			}
		}
		m.priority = t.Priority
		m.step = stepBasic
		m.activeBasic = fieldTitle
		m.titleInput.Focus()
		m.descInput.Blur()
		m.labelsInput.Blur()
		m.assigneeInput.Blur()
		m.searchInput.Blur()
		return m, textinput.Blink
	case "esc":
		m.quitting = true
		return m, tea.Quit
	}
	return m, nil
}

// --- stepBasic ---

func (m Model) handleBasicKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.step = stepTemplate
		m.titleInput.Blur()
		return m, nil

	case "tab", "down", "j":
		m.activeBasic = (m.activeBasic + 1) % basicFieldCount
		return m.syncBasicFocus()

	case "shift+tab", "up", "k":
		m.activeBasic = (m.activeBasic + basicFieldCount - 1) % basicFieldCount
		return m.syncBasicFocus()

	case "enter":
		if m.activeBasic == fieldPriority {
			// Move to next step
			if strings.TrimSpace(m.titleInput.Value()) == "" {
				// Stay on title if empty
				m.activeBasic = fieldTitle
				return m.syncBasicFocus()
			}
			m.step = stepDetails
			m.activeDetail = fieldDesc
			m.titleInput.Blur()
			m.descInput.Focus()
			return m, textinput.Blink
		}
		// Advance field
		m.activeBasic = (m.activeBasic + 1) % basicFieldCount
		return m.syncBasicFocus()

	case "left", "-":
		if m.activeBasic == fieldPriority {
			if m.priority < 4 {
				m.priority++
			}
		}
		return m, nil

	case "right", "+":
		if m.activeBasic == fieldPriority {
			if m.priority > 0 {
				m.priority--
			}
		}
		return m, nil

	case "0", "1", "2", "3", "4":
		if m.activeBasic == fieldPriority {
			m.priority = int(msg.String()[0] - '0')
		}
		return m, nil
	}

	// Cycle issue type with left/right when on type field
	if m.activeBasic == fieldType {
		switch msg.String() {
		case "left":
			if m.typeCursor > 0 {
				m.typeCursor--
			}
			return m, nil
		case "right":
			if m.typeCursor < len(issueTypes)-1 {
				m.typeCursor++
			}
			return m, nil
		}
		// Don't pass other keys to text input when on type/priority
		return m, nil
	}
	if m.activeBasic == fieldPriority {
		return m, nil
	}

	// Delegate to title text input
	var cmd tea.Cmd
	m.titleInput, cmd = m.titleInput.Update(msg)
	return m, cmd
}

func (m Model) syncBasicFocus() (Model, tea.Cmd) {
	switch m.activeBasic {
	case fieldTitle:
		m.titleInput.Focus()
	default:
		m.titleInput.Blur()
	}
	return m, textinput.Blink
}

// --- stepDetails ---

func (m Model) handleDetailsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.step = stepBasic
		m.descInput.Blur()
		m.labelsInput.Blur()
		m.assigneeInput.Blur()
		m.activeBasic = fieldPriority
		m.titleInput.Blur()
		return m, nil

	case "tab", "ctrl+n":
		m.activeDetail = (m.activeDetail + 1) % detailFieldCount
		return m.syncDetailFocus()

	case "shift+tab", "ctrl+p":
		m.activeDetail = (m.activeDetail + detailFieldCount - 1) % detailFieldCount
		return m.syncDetailFocus()

	case "enter":
		if m.activeDetail == fieldAssignee {
			// Move to deps step
			m.step = stepDeps
			m.descInput.Blur()
			m.labelsInput.Blur()
			m.assigneeInput.Blur()
			m.searchInput.Focus()
			// Start loading issues if not done yet
			if !m.depsLoaded {
				return m, m.loadIssuesCmd()
			}
			return m, textinput.Blink
		}
		m.activeDetail = (m.activeDetail + 1) % detailFieldCount
		return m.syncDetailFocus()
	}

	// Delegate to active input
	var cmd tea.Cmd
	switch m.activeDetail {
	case fieldDesc:
		m.descInput, cmd = m.descInput.Update(msg)
	case fieldLabels:
		m.labelsInput, cmd = m.labelsInput.Update(msg)
	case fieldAssignee:
		m.assigneeInput, cmd = m.assigneeInput.Update(msg)
	}
	return m, cmd
}

func (m Model) syncDetailFocus() (Model, tea.Cmd) {
	m.descInput.Blur()
	m.labelsInput.Blur()
	m.assigneeInput.Blur()
	switch m.activeDetail {
	case fieldDesc:
		m.descInput.Focus()
	case fieldLabels:
		m.labelsInput.Focus()
	case fieldAssignee:
		m.assigneeInput.Focus()
	}
	return m, textinput.Blink
}

// --- stepDeps ---

func (m Model) handleDepsKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.step = stepDetails
		m.searchInput.Blur()
		m.activeDetail = fieldAssignee
		m.assigneeInput.Focus()
		return m, textinput.Blink

	case "enter":
		// Done with deps, go to preview
		m.step = stepPreview
		m.searchInput.Blur()
		return m, nil

	case "up", "k":
		if m.depsCursor > 0 {
			m.depsCursor--
		}
		return m, nil

	case "down", "j":
		if m.depsCursor < len(m.filtered)-1 {
			m.depsCursor++
		}
		return m, nil

	case " ":
		// Toggle selected issue as "blocks" dep
		if len(m.filtered) > 0 && m.depsCursor < len(m.filtered) {
			id := m.filtered[m.depsCursor].ID
			if _, ok := m.selectedDeps[id]; ok {
				delete(m.selectedDeps, id)
			} else {
				m.selectedDeps[id] = "blocks"
			}
		}
		return m, nil

	case "b":
		// Add/toggle as "blocks"
		if len(m.filtered) > 0 && m.depsCursor < len(m.filtered) {
			id := m.filtered[m.depsCursor].ID
			m.selectedDeps[id] = "blocks"
		}
		return m, nil

	case "r":
		// Add/toggle as "related"
		if len(m.filtered) > 0 && m.depsCursor < len(m.filtered) {
			id := m.filtered[m.depsCursor].ID
			m.selectedDeps[id] = "related"
		}
		return m, nil

	case "p":
		// Add/toggle as "parent-child"
		if len(m.filtered) > 0 && m.depsCursor < len(m.filtered) {
			id := m.filtered[m.depsCursor].ID
			m.selectedDeps[id] = "parent-child"
		}
		return m, nil

	case "x":
		// Remove dep
		if len(m.filtered) > 0 && m.depsCursor < len(m.filtered) {
			id := m.filtered[m.depsCursor].ID
			delete(m.selectedDeps, id)
		}
		return m, nil
	}

	// Delegate to search input and update filter
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.filterIssues()
	// Reset cursor when filter changes
	if m.depsCursor >= len(m.filtered) {
		m.depsCursor = 0
	}
	return m, cmd
}

func (m *Model) filterIssues() {
	q := strings.ToLower(strings.TrimSpace(m.searchInput.Value()))
	if q == "" {
		m.filtered = m.allIssues
		return
	}
	var out []*types.Issue
	for _, iss := range m.allIssues {
		if strings.Contains(strings.ToLower(iss.ID), q) ||
			strings.Contains(strings.ToLower(iss.Title), q) {
			out = append(out, iss)
		}
	}
	m.filtered = out
}

// --- stepPreview ---

func (m Model) handlePreviewKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		m.confirmed = true
		return m, tea.Quit
	case "n", "esc":
		// Go back to deps
		m.step = stepDeps
		m.searchInput.Focus()
		return m, textinput.Blink
	}
	return m, nil
}

// delegateInputUpdate forwards non-key messages to whichever input is focused.
func (m Model) delegateInputUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.step {
	case stepBasic:
		if m.activeBasic == fieldTitle {
			m.titleInput, cmd = m.titleInput.Update(msg)
		}
	case stepDetails:
		switch m.activeDetail {
		case fieldDesc:
			m.descInput, cmd = m.descInput.Update(msg)
		case fieldLabels:
			m.labelsInput, cmd = m.labelsInput.Update(msg)
		case fieldAssignee:
			m.assigneeInput, cmd = m.assigneeInput.Update(msg)
		}
	case stepDeps:
		m.searchInput, cmd = m.searchInput.Update(msg)
	}
	return m, cmd
}

// loadIssuesCmd loads open/in-progress issues for the dependency picker.
func (m Model) loadIssuesCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx
		if m.cfg.Store == nil {
			return issuesLoadedMsg{issues: nil, err: nil}
		}
		open := types.StatusOpen
		filter := types.IssueFilter{Status: &open}
		issues, err := m.cfg.Store.SearchIssues(ctx, "", filter)
		if err != nil {
			return issuesLoadedMsg{err: err}
		}
		return issuesLoadedMsg{issues: issues}
	}
}

// ---------- View ----------

// View renders the current step.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}
	if m.quitting && !m.confirmed {
		return common.MutedTextStyle.Render("Canceled.\n")
	}

	var body string
	switch m.step {
	case stepTemplate:
		body = m.viewTemplate()
	case stepBasic:
		body = m.viewBasic()
	case stepDetails:
		body = m.viewDetails()
	case stepDeps:
		body = m.viewDeps()
	case stepPreview:
		body = m.viewPreview()
	}

	return body
}

// stepHeader renders a consistent page header showing current step name and progress.
func (m Model) stepHeader(title string) string {
	steps := []string{"Template", "Basic", "Details", "Deps", "Preview"}
	var sb strings.Builder
	for i, s := range steps {
		if step(i) == m.step {
			sb.WriteString(common.Selected.Render(" " + s + " "))
		} else if step(i) < m.step {
			sb.WriteString(common.MutedTextStyle.Render(" " + s + " "))
		} else {
			sb.WriteString(common.DarkTextStyle.Render(" " + s + " "))
		}
	}
	header := lipgloss.JoinVertical(lipgloss.Left,
		common.BrightTextStyle.Bold(true).Render("  New Issue  ")+
			common.MutedTextStyle.Render(" | ")+
			sb.String(),
		common.MutedTextStyle.Render(strings.Repeat("─", min(m.width, 80))),
		common.BrightTextStyle.Bold(true).Render("  "+title),
		"",
	)
	return header
}

func (m Model) viewTemplate() string {
	var sb strings.Builder
	sb.WriteString(m.stepHeader("Choose a template"))

	for i, t := range templates {
		cursor := "  "
		style := common.NormalTextStyle
		if i == m.templateCursor {
			cursor = common.Cursor.Render("> ")
			style = lipgloss.NewStyle().Foreground(common.NeonCyan).Bold(true)
		}
		sb.WriteString(cursor + style.Render(t.Name) + "\n")
		if i == m.templateCursor && t.Description != "" {
			sb.WriteString(common.MutedTextStyle.Render("    "+firstLine(t.Description)) + "\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(common.MutedTextStyle.Render("  ↑/k ↓/j  navigate    enter  select    esc/q  cancel"))
	return sb.String()
}

func (m Model) viewBasic() string {
	var sb strings.Builder
	sb.WriteString(m.stepHeader("Basic Info"))

	// Title field
	titleLabel := m.fieldLabel("Title", m.activeBasic == fieldTitle)
	sb.WriteString(titleLabel + "\n")
	sb.WriteString("  " + m.titleInput.View() + "\n\n")

	// Type field
	typeLabel := m.fieldLabel("Type", m.activeBasic == fieldType)
	sb.WriteString(typeLabel + "\n  ")
	for i, it := range issueTypes {
		if i == m.typeCursor {
			sb.WriteString(common.Selected.Render(" "+it+" "))
		} else {
			sb.WriteString(common.MutedTextStyle.Render(" "+it+" "))
		}
	}
	sb.WriteString("\n\n")

	// Priority field
	prioLabel := m.fieldLabel("Priority", m.activeBasic == fieldPriority)
	sb.WriteString(prioLabel + "\n  ")
	sb.WriteString(common.RenderPriority(m.priority) + "  " + common.MutedTextStyle.Render(priorityLabels[m.priority]))
	sb.WriteString("\n\n")

	sb.WriteString(common.MutedTextStyle.Render("  tab/j/k  next field    +/-/0-4  priority    enter  next step    esc  back"))
	return sb.String()
}

func (m Model) viewDetails() string {
	var sb strings.Builder
	sb.WriteString(m.stepHeader("Details"))

	sb.WriteString(m.fieldLabel("Description", m.activeDetail == fieldDesc) + "\n")
	sb.WriteString("  " + m.descInput.View() + "\n\n")

	sb.WriteString(m.fieldLabel("Labels", m.activeDetail == fieldLabels) + "\n")
	sb.WriteString("  " + m.labelsInput.View() + "\n\n")

	sb.WriteString(m.fieldLabel("Assignee", m.activeDetail == fieldAssignee) + "\n")
	sb.WriteString("  " + m.assigneeInput.View() + "\n\n")

	sb.WriteString(common.MutedTextStyle.Render("  tab  next field    enter on Assignee  next step    esc  back"))
	return sb.String()
}

func (m Model) viewDeps() string {
	var sb strings.Builder
	sb.WriteString(m.stepHeader("Dependencies"))

	sb.WriteString("  " + m.searchInput.View() + "\n\n")

	if !m.depsLoaded {
		sb.WriteString(common.MutedTextStyle.Render("  Loading issues...") + "\n")
	} else if len(m.filtered) == 0 {
		if m.searchInput.Value() != "" {
			sb.WriteString(common.MutedTextStyle.Render("  No issues match your search.") + "\n")
		} else {
			sb.WriteString(common.MutedTextStyle.Render("  No open issues found.") + "\n")
		}
	} else {
		maxRows := m.height - 14
		if maxRows < 3 {
			maxRows = 3
		}
		start := 0
		if m.depsCursor >= maxRows {
			start = m.depsCursor - maxRows + 1
		}
		end := start + maxRows
		if end > len(m.filtered) {
			end = len(m.filtered)
		}
		for i := start; i < end; i++ {
			iss := m.filtered[i]
			cursor := "  "
			titleStyle := common.NormalTextStyle
			if i == m.depsCursor {
				cursor = common.Cursor.Render("> ")
				titleStyle = lipgloss.NewStyle().Foreground(common.NeonCyan)
			}
			depMark := "  "
			if dt, ok := m.selectedDeps[iss.ID]; ok {
				depMark = lipgloss.NewStyle().Foreground(common.NeonGreen).Render("[*]")
				switch dt {
				case "blocks":
					depMark = lipgloss.NewStyle().Foreground(common.NeonPink).Render("[B]")
				case "related":
					depMark = lipgloss.NewStyle().Foreground(common.NeonYellow).Render("[R]")
				case "parent-child":
					depMark = lipgloss.NewStyle().Foreground(common.NeonMagenta).Render("[P]")
				}
			} else {
				depMark = "   "
			}
			line := fmt.Sprintf("%s%s %s  %s",
				cursor,
				depMark,
				common.MutedTextStyle.Render(iss.ID),
				titleStyle.Render(truncate(iss.Title, m.width-20)),
			)
			sb.WriteString(line + "\n")
		}
	}

	// Show selected deps summary
	if len(m.selectedDeps) > 0 {
		sb.WriteString("\n")
		sb.WriteString(common.BrightTextStyle.Render("  Selected:") + "\n")
		for id, dt := range m.selectedDeps {
			sb.WriteString(common.MutedTextStyle.Render(fmt.Sprintf("    %s:%s", dt, id)) + "\n")
		}
	}

	sb.WriteString("\n")
	sb.WriteString(common.MutedTextStyle.Render("  b=blocks  r=related  p=parent-child  x=remove  space=toggle    enter  continue    esc  back"))
	return sb.String()
}

func (m Model) viewPreview() string {
	var sb strings.Builder
	sb.WriteString(m.stepHeader("Preview"))

	row := func(label, val string) string {
		if val == "" {
			return ""
		}
		return fmt.Sprintf("  %s  %s\n",
			lipgloss.NewStyle().Foreground(common.NeonCyan).Width(14).Render(label+":"),
			common.BrightTextStyle.Render(val),
		)
	}

	sb.WriteString(row("Title", m.titleInput.Value()))
	sb.WriteString(row("Type", issueTypes[m.typeCursor]))
	sb.WriteString(row("Priority", common.RenderPriority(m.priority)+" "+priorityLabels[m.priority]))
	if m.descInput.Value() != "" {
		sb.WriteString(row("Description", truncate(m.descInput.Value(), m.width-20)))
	}
	sb.WriteString(row("Labels", m.labelsInput.Value()))
	sb.WriteString(row("Assignee", m.assigneeInput.Value()))

	if len(m.selectedDeps) > 0 {
		sb.WriteString("  " + lipgloss.NewStyle().Foreground(common.NeonCyan).Render("Deps:") + "\n")
		for id, dt := range m.selectedDeps {
			sb.WriteString(common.MutedTextStyle.Render(fmt.Sprintf("    %s:%s", dt, id)) + "\n")
		}
	}

	if strings.TrimSpace(m.titleInput.Value()) == "" {
		sb.WriteString("\n")
		sb.WriteString(lipgloss.NewStyle().Foreground(common.NeonPink).Render("  Title is required. Press esc to go back.") + "\n")
		sb.WriteString("\n")
		sb.WriteString(common.MutedTextStyle.Render("  esc  back"))
	} else {
		sb.WriteString("\n")
		sb.WriteString(common.BrightTextStyle.Bold(true).Render("  Create this issue? "))
		sb.WriteString(lipgloss.NewStyle().Foreground(common.NeonGreen).Bold(true).Render("y/enter") +
			common.MutedTextStyle.Render(" = yes   ") +
			lipgloss.NewStyle().Foreground(common.NeonPink).Render("n/esc") +
			common.MutedTextStyle.Render(" = back"))
	}

	return sb.String()
}

// --- helpers ---

func (m Model) fieldLabel(name string, active bool) string {
	if active {
		return "  " + lipgloss.NewStyle().Foreground(common.NeonCyan).Bold(true).Render(name)
	}
	return "  " + common.MutedTextStyle.Render(name)
}

func firstLine(s string) string {
	if idx := strings.Index(s, "\n"); idx >= 0 {
		return s[:idx]
	}
	return s
}

func truncate(s string, max int) string {
	if max <= 0 {
		return s
	}
	// Strip newlines for inline display
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

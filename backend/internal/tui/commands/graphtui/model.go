// Package graphtui provides the Bubble Tea model for `bd graph --tui`.
package graphtui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/tui/common"
	"github.com/steveyegge/beads/internal/tui/components/issuedetail"
	"github.com/steveyegge/beads/internal/tui/components/statusbar"
	"github.com/steveyegge/beads/internal/types"
)

// TreeNode represents a single node in the flat display list.
type TreeNode struct {
	Issue       *types.Issue
	Depth       int
	Expanded    bool
	HasChildren bool
	Children    []string // child issue IDs (blocked-by direction: these IDs block this issue)
}

type viewState int

const (
	stateLoading viewState = iota
	stateTree
	stateDetail
)

// Config holds the external dependencies injected from the CLI layer.
type Config struct {
	Ctx          context.Context
	Store        storage.Storage
	DaemonClient *rpc.Client
	RootID       string // empty = show all open issues
	Actor        string
	DbPath       string
	LockTimeout  time.Duration
}

// graphTreeKeys defines keybindings for the graph tree view.
type graphTreeKeys struct {
	Up       key.Binding
	Down     key.Binding
	Expand   key.Binding
	Collapse key.Binding
	Toggle   key.Binding
	Focus    key.Binding
	Back     key.Binding
	Enter    key.Binding
	Quit     key.Binding
	Help     key.Binding
}

func newGraphTreeKeys() graphTreeKeys {
	return graphTreeKeys{
		Up:       key.NewBinding(key.WithKeys("up", "k"), key.WithHelp("↑/k", "up")),
		Down:     key.NewBinding(key.WithKeys("down", "j"), key.WithHelp("↓/j", "down")),
		Expand:   key.NewBinding(key.WithKeys("right", "l"), key.WithHelp("→/l", "expand/detail")),
		Collapse: key.NewBinding(key.WithKeys("left", "h"), key.WithHelp("←/h", "collapse/back")),
		Toggle:   key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle")),
		Focus:    key.NewBinding(key.WithKeys("f"), key.WithHelp("f", "focus subtree")),
		Back:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "unfocus/back")),
		Enter:    key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "detail")),
		Quit:     key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q", "quit")),
		Help:     key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "help")),
	}
}

func (k graphTreeKeys) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Expand, k.Collapse, k.Toggle, k.Focus, k.Quit}
}

func (k graphTreeKeys) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},
		{k.Expand, k.Collapse, k.Toggle},
		{k.Focus, k.Back, k.Enter},
		{k.Help, k.Quit},
	}
}

// graphLoadedMsg is sent when issues and deps have been loaded.
type graphLoadedMsg struct {
	allIssues map[string]*types.Issue
	childMap  map[string][]string // parentID -> childIDs (blocker direction)
	rootIDs   []string            // top-level issue IDs (no parents in set)
}

// Model is the top-level Bubble Tea model for the graph TUI.
type Model struct {
	cfg        Config
	state      viewState
	nodes      []TreeNode // flat list of visible nodes (after expand/collapse)
	allIssues  map[string]*types.Issue
	childMap   map[string][]string // parentID -> childIDs (dep direction: blocker = parent)
	rootIDs    []string            // top-level issue IDs with no parents in the set
	cursor     int                 // current selected index in nodes
	detail     issuedetail.Model
	statusbar  statusbar.Model
	treeKeys   graphTreeKeys
	detailKeys common.DetailKeys
	width      int
	height     int
	focusID    string // if non-empty, only show subtree of this issue
}

// New creates a new graph TUI model.
func New(cfg Config) Model {
	treeKeys := newGraphTreeKeys()
	detailKeys := common.NewDetailKeys()

	return Model{
		cfg:        cfg,
		state:      stateLoading,
		detail:     issuedetail.New(0, 0),
		statusbar:  statusbar.New(treeKeys),
		treeKeys:   treeKeys,
		detailKeys: detailKeys,
	}
}

// Init returns the initial command to load the graph.
func (m Model) Init() tea.Cmd {
	return m.loadGraphCmd()
}

// Update handles all messages and routes to the appropriate component.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		contentH := m.height - 2 // leave room for status bar
		m.detail.SetSize(m.width, contentH)
		m.statusbar.SetWidth(m.width)
		return m, nil

	case graphLoadedMsg:
		m.allIssues = msg.allIssues
		m.childMap = msg.childMap
		m.rootIDs = msg.rootIDs
		m.state = stateTree
		m.cursor = 0
		m.nodes = m.buildFlatList()
		return m, nil

	case common.IssueDetailLoadedMsg:
		m.detail.SetContent(msg.Issue, msg.Labels, msg.Deps, msg.Dependents)
		m.state = stateDetail
		m.statusbar.SetKeyMap(m.detailKeys)
		return m, nil

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

	return m.delegateUpdate(msg)
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global quit
	if key.Matches(msg, m.treeKeys.Quit) {
		return m, tea.Quit
	}

	switch m.state {
	case stateTree:
		switch {
		case key.Matches(msg, m.treeKeys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil

		case key.Matches(msg, m.treeKeys.Down):
			if m.cursor < len(m.nodes)-1 {
				m.cursor++
			}
			return m, nil

		case key.Matches(msg, m.treeKeys.Toggle):
			m.toggleExpand(m.cursor)
			return m, nil

		case key.Matches(msg, m.treeKeys.Expand):
			if m.cursor < len(m.nodes) {
				node := &m.nodes[m.cursor]
				if node.HasChildren && !node.Expanded {
					m.toggleExpand(m.cursor)
					return m, nil
				}
				// Open detail view
				return m, m.loadDetailCmd(node.Issue.ID)
			}
			return m, nil

		case key.Matches(msg, m.treeKeys.Enter):
			if m.cursor < len(m.nodes) {
				return m, m.loadDetailCmd(m.nodes[m.cursor].Issue.ID)
			}
			return m, nil

		case key.Matches(msg, m.treeKeys.Collapse):
			if m.cursor < len(m.nodes) {
				node := &m.nodes[m.cursor]
				if node.HasChildren && node.Expanded {
					m.toggleExpand(m.cursor)
					return m, nil
				}
			}
			return m, nil

		case key.Matches(msg, m.treeKeys.Focus):
			if m.cursor < len(m.nodes) {
				m.focusID = m.nodes[m.cursor].Issue.ID
				m.cursor = 0
				m.nodes = m.buildFlatList()
			}
			return m, nil

		case key.Matches(msg, m.treeKeys.Back):
			if m.focusID != "" {
				m.focusID = ""
				m.cursor = 0
				m.nodes = m.buildFlatList()
			}
			return m, nil
		}

	case stateDetail:
		switch {
		case key.Matches(msg, m.detailKeys.Back):
			m.state = stateTree
			m.statusbar.SetKeyMap(m.treeKeys)
			return m, nil
		case key.Matches(msg, m.treeKeys.Collapse):
			m.state = stateTree
			m.statusbar.SetKeyMap(m.treeKeys)
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
	case stateDetail:
		m.detail, cmd = m.detail.Update(msg)
	}
	return m, cmd
}

// View renders the current view.
func (m Model) View() string {
	if m.width == 0 {
		return ""
	}

	var content string
	switch m.state {
	case stateLoading:
		content = lipgloss.Place(
			m.width, m.height-2,
			lipgloss.Center, lipgloss.Center,
			common.MutedTextStyle.Render("Loading dependency graph..."),
		)
	case stateTree:
		content = m.renderTree()
	case stateDetail:
		content = m.detail.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		content,
		m.statusbar.View(),
	)
}

// renderTree renders the flat visible node list.
func (m Model) renderTree() string {
	if len(m.nodes) == 0 {
		return lipgloss.Place(
			m.width, m.height-2,
			lipgloss.Center, lipgloss.Center,
			common.MutedTextStyle.Render("No issues found"),
		)
	}

	contentH := m.height - 2
	if contentH <= 0 {
		contentH = 1
	}

	// Determine visible window
	start := 0
	end := len(m.nodes)
	if end > contentH {
		// Scroll so cursor stays visible
		if m.cursor >= contentH {
			start = m.cursor - contentH + 1
		}
		end = start + contentH
		if end > len(m.nodes) {
			end = len(m.nodes)
		}
	}

	var sb strings.Builder
	for i := start; i < end; i++ {
		node := m.nodes[i]
		line := m.renderNode(node, i == m.cursor)
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	// Pad remaining lines so statusbar stays at bottom
	rendered := sb.String()
	lineCount := end - start
	for i := lineCount; i < contentH; i++ {
		rendered += "\n"
	}
	return rendered
}

// renderNode renders a single tree node line.
func (m Model) renderNode(node TreeNode, selected bool) string {
	// Indent: 2 spaces per depth level
	indent := strings.Repeat("  ", node.Depth)

	// Expand/collapse indicator
	var indicator string
	switch {
	case node.HasChildren && node.Expanded:
		indicator = lipgloss.NewStyle().Foreground(common.NeonCyan).Render("▼")
	case node.HasChildren && !node.Expanded:
		indicator = lipgloss.NewStyle().Foreground(common.NeonMagenta).Render("▶")
	default:
		indicator = "  "
	}

	statusIcon := common.RenderStatusIcon(string(node.Issue.Status))
	priority := common.RenderPriority(node.Issue.Priority)
	id := common.MutedTextStyle.Render(node.Issue.ID)
	title := node.Issue.Title
	if len([]rune(title)) > 50 {
		title = string([]rune(title)[:49]) + "…"
	}

	line := fmt.Sprintf("%s%s %s %s %s  %s",
		indent,
		indicator,
		statusIcon,
		id,
		priority,
		title,
	)

	if selected {
		// Highlight entire line
		availableWidth := m.width
		if availableWidth <= 0 {
			availableWidth = 80
		}
		line = lipgloss.NewStyle().
			Background(common.NightBgHighlight).
			Foreground(common.TextBright).
			Width(availableWidth).
			Render(line)
	} else {
		line = common.NormalTextStyle.Render(line)
	}

	return line
}

// toggleExpand toggles the expanded state of the node at the given index
// and rebuilds the flat list.
func (m *Model) toggleExpand(idx int) {
	if idx < 0 || idx >= len(m.nodes) {
		return
	}
	issueID := m.nodes[idx].Issue.ID
	// Toggle expanded state in the canonical map
	// We use focusID + expandedSet; instead we track via the nodes slice itself.
	// Find the node for this issue and flip its Expanded flag.
	m.nodes[idx].Expanded = !m.nodes[idx].Expanded
	// Rebuild from scratch preserving expanded states
	expandedSet := make(map[string]bool)
	for _, n := range m.nodes {
		if n.Expanded {
			expandedSet[n.Issue.ID] = true
		}
	}
	_ = issueID
	m.nodes = m.buildFlatListWithExpanded(expandedSet)
}

// buildFlatList builds the visible flat list from root nodes.
// All nodes start collapsed.
func (m Model) buildFlatList() []TreeNode {
	return m.buildFlatListWithExpanded(nil)
}

// buildFlatListWithExpanded builds the flat list, restoring expanded state from expandedSet.
func (m Model) buildFlatListWithExpanded(expandedSet map[string]bool) []TreeNode {
	var result []TreeNode

	roots := m.rootIDs
	if m.focusID != "" {
		roots = []string{m.focusID}
	}

	visited := make(map[string]bool)
	m.appendNodes(&result, roots, 0, expandedSet, visited)
	return result
}

// appendNodes recursively appends nodes to result in DFS order.
func (m Model) appendNodes(result *[]TreeNode, ids []string, depth int, expandedSet map[string]bool, visited map[string]bool) {
	for _, id := range ids {
		if visited[id] {
			continue
		}
		visited[id] = true

		issue, ok := m.allIssues[id]
		if !ok {
			continue
		}

		children := m.childMap[id]
		hasChildren := len(children) > 0
		expanded := expandedSet != nil && expandedSet[id]

		node := TreeNode{
			Issue:       issue,
			Depth:       depth,
			Expanded:    expanded,
			HasChildren: hasChildren,
			Children:    children,
		}
		*result = append(*result, node)

		if expanded && hasChildren {
			m.appendNodes(result, children, depth+1, expandedSet, visited)
		}
	}
}

// loadGraphCmd loads all issues and builds the dep tree.
func (m Model) loadGraphCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := m.cfg.Ctx
		s := m.cfg.Store
		if s == nil {
			return common.ErrMsg{Err: fmt.Errorf("no storage available")}
		}

		// Load all active issues (open, in_progress, blocked)
		allIssues := make(map[string]*types.Issue)
		for _, status := range []types.Status{types.StatusOpen, types.StatusInProgress, types.StatusBlocked} {
			statusCopy := status
			issues, err := s.SearchIssues(ctx, "", types.IssueFilter{Status: &statusCopy})
			if err != nil {
				return common.ErrMsg{Err: fmt.Errorf("loading issues: %w", err)}
			}
			for _, iss := range issues {
				allIssues[iss.ID] = iss
			}
		}

		// If a rootID is specified, do a BFS to find the connected subgraph
		if m.cfg.RootID != "" {
			subgraph, err := buildSubgraph(ctx, s, m.cfg.RootID, allIssues)
			if err != nil {
				return common.ErrMsg{Err: err}
			}
			allIssues = subgraph
		}

		if len(allIssues) == 0 {
			return graphLoadedMsg{
				allIssues: allIssues,
				childMap:  make(map[string][]string),
				rootIDs:   nil,
			}
		}

		// Load all dependency records for issues in the set
		// childMap: issueID -> []IDs that BLOCK issueID (i.e., issueID depends on them)
		// We want tree direction: parent = the issue that blocks others
		// DepBlocks: dep.IssueID is blocked by dep.DependsOnID
		// So dep.DependsOnID is the "parent" (blocker), dep.IssueID is the "child" (blocked)
		blockerToBlocked := make(map[string][]string) // blocker -> list of issues it blocks
		hasParent := make(map[string]bool)

		for id := range allIssues {
			deps, err := s.GetDependencyRecords(ctx, id)
			if err != nil {
				continue
			}
			for _, dep := range deps {
				if dep.Type != types.DepBlocks {
					continue
				}
				// dep.IssueID depends on dep.DependsOnID
				// dep.DependsOnID blocks dep.IssueID
				blocker := dep.DependsOnID
				blocked := dep.IssueID
				// Only include if both ends are in our issue set
				if _, ok := allIssues[blocker]; !ok {
					continue
				}
				if _, ok := allIssues[blocked]; !ok {
					continue
				}
				blockerToBlocked[blocker] = appendUnique(blockerToBlocked[blocker], blocked)
				hasParent[blocked] = true
			}
		}

		// Root nodes: issues with no blocker parents in the set
		var rootIDs []string
		for id := range allIssues {
			if !hasParent[id] {
				rootIDs = append(rootIDs, id)
			}
		}

		// Sort rootIDs by priority then ID for stable ordering
		sortIssueIDs(rootIDs, allIssues)

		// Sort each blocker's children list
		for id := range blockerToBlocked {
			sortIssueIDs(blockerToBlocked[id], allIssues)
		}

		return graphLoadedMsg{
			allIssues: allIssues,
			childMap:  blockerToBlocked,
			rootIDs:   rootIDs,
		}
	}
}

// buildSubgraph does a BFS from rootID to find connected issues.
func buildSubgraph(ctx context.Context, s storage.Storage, rootID string, allIssues map[string]*types.Issue) (map[string]*types.Issue, error) {
	// Make sure the root issue is loaded (may not be active)
	root, ok := allIssues[rootID]
	if !ok {
		var err error
		root, err = s.GetIssue(ctx, rootID)
		if err != nil {
			return nil, fmt.Errorf("loading root issue %s: %w", rootID, err)
		}
		if root == nil {
			return nil, fmt.Errorf("issue %s not found", rootID)
		}
	}

	subgraph := map[string]*types.Issue{root.ID: root}
	queue := []string{root.ID}
	visited := map[string]bool{root.ID: true}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		// Traverse downstream: issues that this one blocks
		dependents, err := s.GetDependents(ctx, cur)
		if err == nil {
			for _, dep := range dependents {
				if !visited[dep.ID] {
					visited[dep.ID] = true
					subgraph[dep.ID] = dep
					queue = append(queue, dep.ID)
				}
			}
		}
	}

	return subgraph, nil
}

// loadDetailCmd loads issue detail for the detail view.
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
			labels = nil
		}

		deps, err := s.GetDependenciesWithMetadata(ctx, id)
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("loading deps for %s: %w", id, err)}
		}

		dependents, err := s.GetDependentsWithMetadata(ctx, id)
		if err != nil {
			return common.ErrMsg{Err: fmt.Errorf("loading dependents for %s: %w", id, err)}
		}

		var issueLabels []string
		if labels != nil {
			issueLabels = labels[id]
		}

		return common.IssueDetailLoadedMsg{
			Issue:      issue,
			Labels:     issueLabels,
			Deps:       deps,
			Dependents: dependents,
		}
	}
}

// appendUnique appends s to slice only if not already present.
func appendUnique(slice []string, s string) []string {
	for _, v := range slice {
		if v == s {
			return slice
		}
	}
	return append(slice, s)
}

// sortIssueIDs sorts a slice of IDs by priority (ascending) then ID (ascending).
func sortIssueIDs(ids []string, issues map[string]*types.Issue) {
	for i := 1; i < len(ids); i++ {
		for j := i; j > 0; j-- {
			a := issues[ids[j-1]]
			b := issues[ids[j]]
			if a == nil || b == nil {
				break
			}
			if a.Priority > b.Priority || (a.Priority == b.Priority && a.ID > b.ID) {
				ids[j-1], ids[j] = ids[j], ids[j-1]
			} else {
				break
			}
		}
	}
}

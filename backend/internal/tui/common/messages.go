package common

import "github.com/steveyegge/beads/internal/types"

// IssuesLoadedMsg is sent when the issue list has been loaded.
type IssuesLoadedMsg struct {
	Issues []*types.Issue
	Labels map[string][]string // issueID -> labels
}

// IssueDetailLoadedMsg is sent when a single issue's detail is loaded.
type IssueDetailLoadedMsg struct {
	Issue      *types.Issue
	Labels     []string
	Deps       []*types.IssueWithDependencyMetadata
	Dependents []*types.IssueWithDependencyMetadata
}

// IssueClosedMsg is sent after an issue is successfully closed.
type IssueClosedMsg struct {
	ID              string
	NewlyUnblocked  []*types.Issue
}

// ErrMsg wraps an error for the TUI message loop.
type ErrMsg struct {
	Err error
}

// StatusMsg sets a transient status message in the status bar.
type StatusMsg struct {
	Text string
}

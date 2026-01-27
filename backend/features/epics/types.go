// Package epics provides the epic management vertical slice.
package epics

import "github.com/steveyegge/beads/internal/types"

// StatusArgs represents arguments for getting epic status.
type StatusArgs struct {
	EligibleOnly bool `json:"eligible_only,omitempty"`
}

// EpicSummary represents a summary of an epic.
type EpicSummary struct {
	ID              string         `json:"id"`
	Title           string         `json:"title"`
	Status          string         `json:"status"`
	TotalChildren   int            `json:"total_children"`
	ClosedChildren  int            `json:"closed_children"`
	OpenChildren    int            `json:"open_children"`
	BlockedChildren int            `json:"blocked_children"`
	Progress        float64        `json:"progress"` // 0.0 to 1.0
	EligibleForClose bool          `json:"eligible_for_close"`
	Children        []*types.Issue `json:"children,omitempty"`
}

// FeatureSummary represents a feature/epic in the feature list.
type FeatureSummary struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Path        string            `json:"path"`
	KanbanStats map[string]int    `json:"kanban_stats"`
	Artifacts   map[string]bool   `json:"artifacts"`
	Meta        map[string]string `json:"meta,omitempty"`
	IsLegacy    bool              `json:"is_legacy"`
}

// FeaturesResponse is the response for listing features.
type FeaturesResponse struct {
	Features       []FeatureSummary  `json:"features"`
	ProjectPath    string            `json:"project_path"`
	WorktreesRoot  *string           `json:"worktrees_root"`
	ActiveWorktree *string           `json:"active_worktree"`
	ActiveMission  map[string]string `json:"active_mission"`
}

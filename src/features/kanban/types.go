// Package kanban provides the kanban/ready work vertical slice.
package kanban

import (
	"time"

	"github.com/steveyegge/beads/internal/types"
)

// ReadyFilter represents filter options for ready work.
type ReadyFilter struct {
	Assignee        string   `json:"assignee,omitempty"`
	Unassigned      bool     `json:"unassigned,omitempty"`
	Priority        *int     `json:"priority,omitempty"`
	Type            string   `json:"type,omitempty"`
	Limit           int      `json:"limit,omitempty"`
	SortPolicy      string   `json:"sort_policy,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	LabelsAny       []string `json:"labels_any,omitempty"`
	ParentID        string   `json:"parent_id,omitempty"`
	MolType         string   `json:"mol_type,omitempty"`
	IncludeDeferred bool     `json:"include_deferred,omitempty"`
}

// BlockedFilter represents filter options for blocked issues.
type BlockedFilter struct {
	ParentID string `json:"parent_id,omitempty"`
}

// StaleFilter represents filter options for stale issues.
type StaleFilter struct {
	Days   int    `json:"days,omitempty"`
	Status string `json:"status,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// EpicStatusFilter represents filter options for epic status.
type EpicStatusFilter struct {
	EligibleOnly bool `json:"eligible_only,omitempty"`
}

// IssueCard is a simplified issue for the kanban board.
type IssueCard struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Priority    int       `json:"priority"`
	Status      string    `json:"status"`
	Assignee    string    `json:"assignee,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Labels      []string  `json:"labels,omitempty"`
	IsBlocked   bool      `json:"is_blocked"`
	Blockers    []string  `json:"blockers,omitempty"`
}

// KanbanResponse is the response for kanban board data.
type KanbanResponse struct {
	Lanes         map[string][]IssueCard `json:"lanes"`
	IsLegacy      bool                   `json:"is_legacy"`
	UpgradeNeeded bool                   `json:"upgrade_needed"`
}

// BlockedIssueResponse wraps blocked issue info.
type BlockedIssueResponse struct {
	Issue    *types.Issue `json:"issue"`
	Blockers []string     `json:"blockers"`
}

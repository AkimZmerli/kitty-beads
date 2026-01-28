// Package kanban provides the kanban/ready work vertical slice.
package kanban

import (
	"context"

	"github.com/steveyegge/beads/internal/types"
)

// Repository defines the data access interface for ready work and blocking.
type Repository interface {
	// GetReadyWork retrieves issues that are ready to work on.
	GetReadyWork(ctx context.Context, filter types.WorkFilter) ([]*types.Issue, error)

	// GetBlockedIssues retrieves issues that are blocked.
	GetBlockedIssues(ctx context.Context, filter types.WorkFilter) ([]*types.BlockedIssue, error)

	// IsBlocked checks if an issue has open blockers.
	IsBlocked(ctx context.Context, issueID string) (bool, []string, error)

	// GetEpicsEligibleForClosure retrieves epics that can be closed.
	GetEpicsEligibleForClosure(ctx context.Context) ([]*types.EpicStatus, error)

	// GetStaleIssues retrieves issues not updated in a given period.
	GetStaleIssues(ctx context.Context, filter types.StaleFilter) ([]*types.Issue, error)

	// GetNewlyUnblockedByClose retrieves issues unblocked by closing an issue.
	GetNewlyUnblockedByClose(ctx context.Context, closedIssueID string) ([]*types.Issue, error)
}

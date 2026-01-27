// Package issues provides the issue management vertical slice.
package issues

import (
	"context"

	"github.com/steveyegge/beads/internal/types"
)

// Repository defines the data access interface for issues.
type Repository interface {
	// Create creates a new issue.
	Create(ctx context.Context, issue *types.Issue, actor string) error

	// CreateBatch creates multiple issues atomically.
	CreateBatch(ctx context.Context, issues []*types.Issue, actor string) error

	// Get retrieves an issue by ID.
	Get(ctx context.Context, id string) (*types.Issue, error)

	// GetByExternalRef retrieves an issue by external reference (e.g., "gh-9").
	GetByExternalRef(ctx context.Context, externalRef string) (*types.Issue, error)

	// Update updates an issue's fields.
	Update(ctx context.Context, id string, updates map[string]interface{}, actor string) error

	// Close marks an issue as closed with a reason.
	Close(ctx context.Context, id string, reason string, actor string, session string) error

	// Delete removes an issue.
	Delete(ctx context.Context, id string) error

	// Search finds issues matching a query and filter.
	Search(ctx context.Context, query string, filter types.IssueFilter) ([]*types.Issue, error)

	// GetNextChildID generates the next child ID for a parent.
	GetNextChildID(ctx context.Context, parentID string) (string, error)

	// UpdateID updates an issue's ID (for prefix rename operations).
	UpdateID(ctx context.Context, oldID, newID string, issue *types.Issue, actor string) error
}

// Package comments provides the comment management vertical slice.
package comments

import (
	"context"
	"time"

	"github.com/steveyegge/beads/internal/types"
)

// Repository defines the data access interface for comments.
type Repository interface {
	// Add adds a comment to an issue.
	Add(ctx context.Context, issueID, author, text string) (*types.Comment, error)

	// Import adds a comment with a preserved timestamp (for JSONL import).
	Import(ctx context.Context, issueID, author, text string, createdAt time.Time) (*types.Comment, error)

	// Get retrieves all comments for an issue.
	Get(ctx context.Context, issueID string) ([]*types.Comment, error)

	// GetForIssues retrieves comments for multiple issues.
	GetForIssues(ctx context.Context, issueIDs []string) (map[string][]*types.Comment, error)
}

// EventRepository defines the data access interface for events (audit trail).
type EventRepository interface {
	// AddComment adds a comment event to the audit trail.
	AddComment(ctx context.Context, issueID, actor, comment string) error

	// Get retrieves events for an issue.
	Get(ctx context.Context, issueID string, limit int) ([]*types.Event, error)
}

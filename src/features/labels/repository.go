// Package labels provides the label management vertical slice.
package labels

import (
	"context"

	"github.com/steveyegge/beads/internal/types"
)

// Repository defines the data access interface for labels.
type Repository interface {
	// Add adds a label to an issue.
	Add(ctx context.Context, issueID, label, actor string) error

	// Remove removes a label from an issue.
	Remove(ctx context.Context, issueID, label, actor string) error

	// Get retrieves all labels for an issue.
	Get(ctx context.Context, issueID string) ([]string, error)

	// GetForIssues retrieves labels for multiple issues.
	GetForIssues(ctx context.Context, issueIDs []string) (map[string][]string, error)

	// GetIssuesByLabel retrieves all issues with a specific label.
	GetIssuesByLabel(ctx context.Context, label string) ([]*types.Issue, error)
}

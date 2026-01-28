// Package dependencies provides the dependency management vertical slice.
package dependencies

import (
	"context"

	"github.com/steveyegge/beads/internal/types"
)

// Repository defines the data access interface for dependencies.
type Repository interface {
	// Add adds a dependency between issues.
	Add(ctx context.Context, dep *types.Dependency, actor string) error

	// Remove removes a dependency between issues.
	Remove(ctx context.Context, issueID, dependsOnID string, actor string) error

	// GetDependencies retrieves issues that an issue depends on.
	GetDependencies(ctx context.Context, issueID string) ([]*types.Issue, error)

	// GetDependents retrieves issues that depend on an issue.
	GetDependents(ctx context.Context, issueID string) ([]*types.Issue, error)

	// GetDependenciesWithMetadata retrieves dependencies with relationship type.
	GetDependenciesWithMetadata(ctx context.Context, issueID string) ([]*types.IssueWithDependencyMetadata, error)

	// GetDependentsWithMetadata retrieves dependents with relationship type.
	GetDependentsWithMetadata(ctx context.Context, issueID string) ([]*types.IssueWithDependencyMetadata, error)

	// GetRecords retrieves raw dependency records for an issue.
	GetRecords(ctx context.Context, issueID string) ([]*types.Dependency, error)

	// GetAllRecords retrieves all dependency records grouped by issue.
	GetAllRecords(ctx context.Context) (map[string][]*types.Dependency, error)

	// GetCounts retrieves dependency counts for multiple issues.
	GetCounts(ctx context.Context, issueIDs []string) (map[string]*types.DependencyCounts, error)

	// GetTree retrieves the dependency tree for an issue.
	GetTree(ctx context.Context, issueID string, maxDepth int, showAllPaths bool, reverse bool) ([]*types.TreeNode, error)

	// DetectCycles finds circular dependencies.
	DetectCycles(ctx context.Context) ([][]*types.Issue, error)

	// RenamePrefix updates dependency prefixes (for prefix rename operations).
	RenamePrefix(ctx context.Context, oldPrefix, newPrefix string) error
}

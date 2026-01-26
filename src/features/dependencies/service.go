// Package dependencies provides the dependency management vertical slice.
package dependencies

import (
	"context"
	"fmt"
	"strings"

	"github.com/steveyegge/beads/internal/types"
)

// Service provides business logic for dependency operations.
type Service struct {
	repo Repository
}

// NewService creates a new dependency service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Add adds a dependency between issues.
func (s *Service) Add(ctx context.Context, args AddArgs, actor string) (*AddResult, error) {
	if args.FromID == "" {
		return nil, fmt.Errorf("from_id is required")
	}
	if args.ToID == "" {
		return nil, fmt.Errorf("to_id is required")
	}
	if args.DepType == "" {
		args.DepType = string(types.DepBlocks)
	}

	// Check for child->parent dependency anti-pattern
	if isChildOf(args.FromID, args.ToID) {
		return nil, fmt.Errorf("cannot add dependency: %s is already a child of %s (children inherit dependency via hierarchy)", args.FromID, args.ToID)
	}

	depType := types.DependencyType(args.DepType)
	if !depType.IsValid() {
		return nil, fmt.Errorf("invalid dependency type '%s' (valid: blocks, related, parent-child, discovered-from)", args.DepType)
	}

	dep := &types.Dependency{
		IssueID:     args.FromID,
		DependsOnID: args.ToID,
		Type:        depType,
	}

	if err := s.repo.Add(ctx, dep, actor); err != nil {
		return nil, fmt.Errorf("failed to add dependency: %w", err)
	}

	return &AddResult{
		Status:      "added",
		IssueID:     args.FromID,
		DependsOnID: args.ToID,
		Type:        args.DepType,
	}, nil
}

// Remove removes a dependency between issues.
func (s *Service) Remove(ctx context.Context, args RemoveArgs, actor string) (*RemoveResult, error) {
	if args.FromID == "" {
		return nil, fmt.Errorf("from_id is required")
	}
	if args.ToID == "" {
		return nil, fmt.Errorf("to_id is required")
	}

	if err := s.repo.Remove(ctx, args.FromID, args.ToID, actor); err != nil {
		return nil, fmt.Errorf("failed to remove dependency: %w", err)
	}

	return &RemoveResult{
		Status:      "removed",
		IssueID:     args.FromID,
		DependsOnID: args.ToID,
	}, nil
}

// GetDependencies retrieves issues that an issue depends on.
func (s *Service) GetDependencies(ctx context.Context, issueID string) ([]*types.Issue, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	deps, err := s.repo.GetDependencies(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	return deps, nil
}

// GetDependents retrieves issues that depend on an issue.
func (s *Service) GetDependents(ctx context.Context, issueID string) ([]*types.Issue, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	deps, err := s.repo.GetDependents(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents: %w", err)
	}

	return deps, nil
}

// GetDependenciesWithMetadata retrieves dependencies with relationship type.
func (s *Service) GetDependenciesWithMetadata(ctx context.Context, issueID string) ([]*types.IssueWithDependencyMetadata, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	deps, err := s.repo.GetDependenciesWithMetadata(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependencies: %w", err)
	}

	return deps, nil
}

// GetDependentsWithMetadata retrieves dependents with relationship type.
func (s *Service) GetDependentsWithMetadata(ctx context.Context, issueID string) ([]*types.IssueWithDependencyMetadata, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	deps, err := s.repo.GetDependentsWithMetadata(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependents: %w", err)
	}

	return deps, nil
}

// GetTree retrieves the dependency tree for an issue.
func (s *Service) GetTree(ctx context.Context, args TreeArgs) (*TreeResult, error) {
	if args.ID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	maxDepth := args.MaxDepth
	if maxDepth <= 0 {
		maxDepth = 10 // Default max depth
	}

	nodes, err := s.repo.GetTree(ctx, args.ID, maxDepth, args.ShowAllPaths, args.Reverse)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency tree: %w", err)
	}

	return &TreeResult{Nodes: nodes}, nil
}

// GetRecords retrieves raw dependency records for an issue.
func (s *Service) GetRecords(ctx context.Context, issueID string) ([]*types.Dependency, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	records, err := s.repo.GetRecords(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency records: %w", err)
	}

	return records, nil
}

// GetCounts retrieves dependency counts for multiple issues.
func (s *Service) GetCounts(ctx context.Context, issueIDs []string) (map[string]*types.DependencyCounts, error) {
	if len(issueIDs) == 0 {
		return map[string]*types.DependencyCounts{}, nil
	}

	counts, err := s.repo.GetCounts(ctx, issueIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get dependency counts: %w", err)
	}

	return counts, nil
}

// DetectCycles finds circular dependencies.
func (s *Service) DetectCycles(ctx context.Context) (*CycleResult, error) {
	cycles, err := s.repo.DetectCycles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to detect cycles: %w", err)
	}

	return &CycleResult{
		Cycles: cycles,
		Count:  len(cycles),
	}, nil
}

// RenamePrefix updates dependency prefixes.
func (s *Service) RenamePrefix(ctx context.Context, oldPrefix, newPrefix string) error {
	if oldPrefix == "" {
		return fmt.Errorf("old prefix is required")
	}
	if newPrefix == "" {
		return fmt.Errorf("new prefix is required")
	}

	if err := s.repo.RenamePrefix(ctx, oldPrefix, newPrefix); err != nil {
		return fmt.Errorf("failed to rename prefix: %w", err)
	}

	return nil
}

// isChildOf returns true if childID is a hierarchical child of parentID.
func isChildOf(childID, parentID string) bool {
	_, actualParentID, depth := types.ParseHierarchicalID(childID)
	if depth == 0 {
		return false
	}
	if actualParentID == parentID {
		return true
	}
	return strings.HasPrefix(childID, parentID+".")
}

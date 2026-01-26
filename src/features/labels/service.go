// Package labels provides the label management vertical slice.
package labels

import (
	"context"
	"fmt"

	"github.com/steveyegge/beads/internal/types"
)

// Service provides business logic for label operations.
type Service struct {
	repo Repository
}

// NewService creates a new label service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Add adds a label to an issue.
func (s *Service) Add(ctx context.Context, args AddArgs, actor string) (*AddResult, error) {
	if args.IssueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}
	if args.Label == "" {
		return nil, fmt.Errorf("label is required")
	}

	if err := s.repo.Add(ctx, args.IssueID, args.Label, actor); err != nil {
		return nil, fmt.Errorf("failed to add label: %w", err)
	}

	return &AddResult{
		Status:  "added",
		IssueID: args.IssueID,
		Label:   args.Label,
	}, nil
}

// Remove removes a label from an issue.
func (s *Service) Remove(ctx context.Context, args RemoveArgs, actor string) (*RemoveResult, error) {
	if args.IssueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}
	if args.Label == "" {
		return nil, fmt.Errorf("label is required")
	}

	if err := s.repo.Remove(ctx, args.IssueID, args.Label, actor); err != nil {
		return nil, fmt.Errorf("failed to remove label: %w", err)
	}

	return &RemoveResult{
		Status:  "removed",
		IssueID: args.IssueID,
		Label:   args.Label,
	}, nil
}

// Get retrieves all labels for an issue.
func (s *Service) Get(ctx context.Context, issueID string) ([]string, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	labels, err := s.repo.Get(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels: %w", err)
	}

	return labels, nil
}

// GetForIssues retrieves labels for multiple issues.
func (s *Service) GetForIssues(ctx context.Context, issueIDs []string) (map[string][]string, error) {
	if len(issueIDs) == 0 {
		return map[string][]string{}, nil
	}

	labels, err := s.repo.GetForIssues(ctx, issueIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get labels: %w", err)
	}

	return labels, nil
}

// GetIssuesByLabel retrieves all issues with a specific label.
func (s *Service) GetIssuesByLabel(ctx context.Context, label string) ([]*types.Issue, error) {
	if label == "" {
		return nil, fmt.Errorf("label is required")
	}

	issues, err := s.repo.GetIssuesByLabel(ctx, label)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues by label: %w", err)
	}

	return issues, nil
}

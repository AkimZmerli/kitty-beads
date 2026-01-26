// Package comments provides the comment management vertical slice.
package comments

import (
	"context"
	"fmt"
	"time"

	"github.com/steveyegge/beads/internal/types"
)

// Service provides business logic for comment operations.
type Service struct {
	repo      Repository
	eventRepo EventRepository
}

// NewService creates a new comment service.
func NewService(repo Repository, eventRepo EventRepository) *Service {
	return &Service{
		repo:      repo,
		eventRepo: eventRepo,
	}
}

// Add adds a comment to an issue.
func (s *Service) Add(ctx context.Context, args AddArgs) (*types.Comment, error) {
	if args.IssueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}
	if args.Text == "" {
		return nil, fmt.Errorf("comment text is required")
	}

	author := args.Author
	if author == "" {
		author = "anonymous"
	}

	comment, err := s.repo.Add(ctx, args.IssueID, author, args.Text)
	if err != nil {
		return nil, fmt.Errorf("failed to add comment: %w", err)
	}

	return comment, nil
}

// Import adds a comment with a preserved timestamp (for JSONL import).
func (s *Service) Import(ctx context.Context, args ImportArgs) (*types.Comment, error) {
	if args.IssueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}
	if args.Text == "" {
		return nil, fmt.Errorf("comment text is required")
	}

	author := args.Author
	if author == "" {
		author = "anonymous"
	}

	createdAt := args.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}

	comment, err := s.repo.Import(ctx, args.IssueID, author, args.Text, createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to import comment: %w", err)
	}

	return comment, nil
}

// List retrieves all comments for an issue.
func (s *Service) List(ctx context.Context, issueID string) ([]*types.Comment, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	comments, err := s.repo.Get(ctx, issueID)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, nil
}

// ListForIssues retrieves comments for multiple issues.
func (s *Service) ListForIssues(ctx context.Context, issueIDs []string) (map[string][]*types.Comment, error) {
	if len(issueIDs) == 0 {
		return map[string][]*types.Comment{}, nil
	}

	comments, err := s.repo.GetForIssues(ctx, issueIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, nil
}

// AddEvent adds a comment event to the audit trail.
func (s *Service) AddEvent(ctx context.Context, issueID, actor, comment string) error {
	if issueID == "" {
		return fmt.Errorf("issue ID is required")
	}

	if err := s.eventRepo.AddComment(ctx, issueID, actor, comment); err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}

	return nil
}

// ListEvents retrieves events for an issue.
func (s *Service) ListEvents(ctx context.Context, issueID string, limit int) ([]*types.Event, error) {
	if issueID == "" {
		return nil, fmt.Errorf("issue ID is required")
	}

	if limit <= 0 {
		limit = 100 // Default limit
	}

	events, err := s.eventRepo.Get(ctx, issueID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	return events, nil
}

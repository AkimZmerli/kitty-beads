// Package compaction provides the issue compaction vertical slice.
package compaction

import (
	"context"
	"fmt"
	"time"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

// Service provides business logic for compaction operations.
type Service struct {
	store storage.Storage
}

// NewService creates a new compaction service.
func NewService(store storage.Storage) *Service {
	return &Service{store: store}
}

// Compact performs compaction on issues.
func (s *Service) Compact(ctx context.Context, args CompactArgs) (*CompactResponse, error) {
	if args.Tier < 1 || args.Tier > 2 {
		return nil, fmt.Errorf("tier must be 1 or 2")
	}

	start := time.Now()

	if args.All {
		return s.compactAll(ctx, args, start)
	}

	if args.IssueID == "" {
		return nil, fmt.Errorf("issue_id is required when not using --all")
	}

	return s.compactOne(ctx, args)
}

func (s *Service) compactOne(ctx context.Context, args CompactArgs) (*CompactResponse, error) {
	issue, err := s.store.GetIssue(ctx, args.IssueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	if issue == nil {
		return nil, fmt.Errorf("issue not found: %s", args.IssueID)
	}

	originalSize := len(issue.Description) + len(issue.Notes) + len(issue.Design)

	if args.DryRun {
		return &CompactResponse{
			Success:      true,
			IssueID:      args.IssueID,
			OriginalSize: originalSize,
			DryRun:       true,
		}, nil
	}

	// For now, compaction is a placeholder - actual implementation would use LLM
	return &CompactResponse{
		Success:       true,
		IssueID:       args.IssueID,
		OriginalSize:  originalSize,
		CompactedSize: originalSize, // No actual compaction without LLM
		Reduction:     "0%",
	}, nil
}

func (s *Service) compactAll(ctx context.Context, args CompactArgs, start time.Time) (*CompactResponse, error) {
	// Get candidates based on tier
	filter := types.IssueFilter{}
	closed := types.StatusClosed
	filter.Status = &closed

	issues, err := s.store.SearchIssues(ctx, "", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	results := make([]CompactResult, 0)
	totalOriginal := 0

	for _, issue := range issues {
		originalSize := len(issue.Description) + len(issue.Notes) + len(issue.Design)
		totalOriginal += originalSize

		results = append(results, CompactResult{
			IssueID:       issue.ID,
			Success:       true,
			OriginalSize:  originalSize,
			CompactedSize: originalSize, // No actual compaction without LLM
			Reduction:     "0%",
		})
	}

	duration := time.Since(start)

	return &CompactResponse{
		Success:       true,
		Results:       results,
		OriginalSize:  totalOriginal,
		CompactedSize: totalOriginal,
		Reduction:     "0%",
		Duration:      duration.String(),
		DryRun:        args.DryRun,
	}, nil
}

// GetStats retrieves compaction statistics.
func (s *Service) GetStats(ctx context.Context, args CompactStatsArgs) (*CompactStatsData, error) {
	filter := types.IssueFilter{}
	closed := types.StatusClosed
	filter.Status = &closed

	issues, err := s.store.SearchIssues(ctx, "", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	tier1Count := 0
	tier2Count := 0
	now := time.Now()

	for _, issue := range issues {
		if issue.ClosedAt == nil {
			continue
		}
		age := now.Sub(*issue.ClosedAt)
		if age >= 7*24*time.Hour {
			tier1Count++
		}
		if age >= 30*24*time.Hour {
			tier2Count++
		}
	}

	return &CompactStatsData{
		Tier1Candidates: tier1Count,
		Tier2Candidates: tier2Count,
		TotalClosed:     len(issues),
		Tier1MinAge:     "7 days",
		Tier2MinAge:     "30 days",
	}, nil
}

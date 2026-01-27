// Package gates provides the async coordination (gates) vertical slice.
package gates

import (
	"context"
	"fmt"
	"strings"

	"github.com/steveyegge/beads/internal/types"
)

// IssueRepository defines issue operations needed by gates.
type IssueRepository interface {
	Create(ctx context.Context, issue *types.Issue, actor string) error
	Get(ctx context.Context, id string) (*types.Issue, error)
	Search(ctx context.Context, query string, filter types.IssueFilter) ([]*types.Issue, error)
	Update(ctx context.Context, id string, updates map[string]interface{}, actor string) error
	Close(ctx context.Context, id string, reason string, actor string, session string) error
}

// Service provides business logic for gate operations.
type Service struct {
	issueRepo IssueRepository
}

// NewService creates a new gate service.
func NewService(issueRepo IssueRepository) *Service {
	return &Service{issueRepo: issueRepo}
}

// isGate checks if an issue is a gate (has AwaitType set).
func isGate(issue *types.Issue) bool {
	return issue != nil && issue.AwaitType != ""
}

// Create creates a new gate.
func (s *Service) Create(ctx context.Context, args CreateArgs, actor string) (*CreateResult, error) {
	if args.Title == "" {
		return nil, fmt.Errorf("title is required")
	}
	if args.AwaitType == "" {
		return nil, fmt.Errorf("await_type is required")
	}

	// Validate await type
	validTypes := []string{"gh:run", "gh:pr", "timer", "human", "mail"}
	valid := false
	for _, t := range validTypes {
		if args.AwaitType == t {
			valid = true
			break
		}
	}
	if !valid {
		return nil, fmt.Errorf("invalid await_type: %s (valid: %v)", args.AwaitType, validTypes)
	}

	issue := &types.Issue{
		Title:     args.Title,
		IssueType: types.TypeTask, // Gates are tasks with AwaitType set
		Status:    types.StatusOpen,
		AwaitType: args.AwaitType,
		AwaitID:   args.AwaitID,
	}

	if err := s.issueRepo.Create(ctx, issue, actor); err != nil {
		return nil, fmt.Errorf("failed to create gate: %w", err)
	}

	// Add waiters if specified
	if len(args.Waiters) > 0 {
		updates := map[string]interface{}{
			"waiters": args.Waiters,
		}
		if err := s.issueRepo.Update(ctx, issue.ID, updates, actor); err != nil {
			return nil, fmt.Errorf("failed to add waiters: %w", err)
		}
	}

	return &CreateResult{ID: issue.ID}, nil
}

// List lists gates.
func (s *Service) List(ctx context.Context, args ListArgs) ([]*GateInfo, error) {
	filter := types.IssueFilter{}

	if !args.All {
		status := types.StatusOpen
		filter.Status = &status
	}

	issues, err := s.issueRepo.Search(ctx, "", filter)
	if err != nil {
		return nil, fmt.Errorf("failed to list gates: %w", err)
	}

	// Filter to only gates (issues with AwaitType set)
	gates := make([]*GateInfo, 0)
	for _, issue := range issues {
		if isGate(issue) {
			gates = append(gates, issueToGateInfo(issue))
		}
	}

	return gates, nil
}

// Show shows a gate by ID.
func (s *Service) Show(ctx context.Context, args ShowArgs) (*GateInfo, error) {
	if args.ID == "" {
		return nil, fmt.Errorf("gate ID is required")
	}

	issue, err := s.issueRepo.Get(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gate: %w", err)
	}
	if issue == nil {
		return nil, fmt.Errorf("gate not found: %s", args.ID)
	}
	if !isGate(issue) {
		return nil, fmt.Errorf("issue %s is not a gate", args.ID)
	}

	return issueToGateInfo(issue), nil
}

// Close closes a gate.
func (s *Service) Close(ctx context.Context, args CloseArgs, actor string) error {
	if args.ID == "" {
		return fmt.Errorf("gate ID is required")
	}

	issue, err := s.issueRepo.Get(ctx, args.ID)
	if err != nil {
		return fmt.Errorf("failed to get gate: %w", err)
	}
	if issue == nil {
		return fmt.Errorf("gate not found: %s", args.ID)
	}
	if !isGate(issue) {
		return fmt.Errorf("issue %s is not a gate", args.ID)
	}

	if err := s.issueRepo.Close(ctx, args.ID, args.Reason, actor, ""); err != nil {
		return fmt.Errorf("failed to close gate: %w", err)
	}

	return nil
}

// Wait adds waiters to a gate.
func (s *Service) Wait(ctx context.Context, args WaitArgs, actor string) (*WaitResult, error) {
	if args.ID == "" {
		return nil, fmt.Errorf("gate ID is required")
	}
	if len(args.Waiters) == 0 {
		return nil, fmt.Errorf("at least one waiter is required")
	}

	issue, err := s.issueRepo.Get(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get gate: %w", err)
	}
	if issue == nil {
		return nil, fmt.Errorf("gate not found: %s", args.ID)
	}
	if !isGate(issue) {
		return nil, fmt.Errorf("issue %s is not a gate", args.ID)
	}

	// Merge with existing waiters
	existingWaiters := make(map[string]bool)
	for _, w := range issue.Waiters {
		existingWaiters[w] = true
	}

	addedCount := 0
	for _, w := range args.Waiters {
		if !existingWaiters[w] {
			issue.Waiters = append(issue.Waiters, w)
			addedCount++
		}
	}

	if addedCount > 0 {
		updates := map[string]interface{}{
			"waiters": issue.Waiters,
		}
		if err := s.issueRepo.Update(ctx, args.ID, updates, actor); err != nil {
			return nil, fmt.Errorf("failed to update waiters: %w", err)
		}
	}

	return &WaitResult{AddedCount: addedCount}, nil
}

// ResolveGateID resolves a partial gate ID to a full ID.
func (s *Service) ResolveGateID(ctx context.Context, partialID string) (string, error) {
	filter := types.IssueFilter{}

	issues, err := s.issueRepo.Search(ctx, "", filter)
	if err != nil {
		return "", fmt.Errorf("failed to search gates: %w", err)
	}

	var matches []*types.Issue
	for _, issue := range issues {
		if !isGate(issue) {
			continue
		}
		if strings.HasPrefix(issue.ID, partialID) || strings.Contains(issue.ID, partialID) {
			matches = append(matches, issue)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no gate found matching: %s", partialID)
	}
	if len(matches) > 1 {
		ids := make([]string, len(matches))
		for i, m := range matches {
			ids[i] = m.ID
		}
		return "", fmt.Errorf("ambiguous gate ID: %s matches %v", partialID, ids)
	}

	return matches[0].ID, nil
}

func issueToGateInfo(issue *types.Issue) *GateInfo {
	info := &GateInfo{
		Issue:     issue,
		AwaitType: issue.AwaitType,
		AwaitID:   issue.AwaitID,
		Waiters:   issue.Waiters,
		Status:    string(issue.Status),
		CreatedAt: issue.CreatedAt,
	}
	if issue.Status == types.StatusClosed && issue.ClosedAt != nil {
		info.ClosedAt = issue.ClosedAt
	}
	return info
}

// Package kanban provides the kanban/ready work vertical slice.
package kanban

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/util"
)

// IssueRepository defines issue operations needed by kanban.
type IssueRepository interface {
	Search(ctx context.Context, query string, filter types.IssueFilter) ([]*types.Issue, error)
}

// Service provides business logic for kanban operations.
type Service struct {
	repo      Repository
	issueRepo IssueRepository
}

// NewService creates a new kanban service.
func NewService(repo Repository, issueRepo IssueRepository) *Service {
	return &Service{
		repo:      repo,
		issueRepo: issueRepo,
	}
}

// GetReady retrieves issues that are ready to work on.
func (s *Service) GetReady(ctx context.Context, filter ReadyFilter) ([]*types.Issue, error) {
	wf := types.WorkFilter{
		Type:            filter.Type,
		Priority:        filter.Priority,
		Unassigned:      filter.Unassigned,
		Limit:           filter.Limit,
		SortPolicy:      types.SortPolicy(filter.SortPolicy),
		Labels:          util.NormalizeLabels(filter.Labels),
		LabelsAny:       util.NormalizeLabels(filter.LabelsAny),
		IncludeDeferred: filter.IncludeDeferred,
	}
	if filter.Assignee != "" && !filter.Unassigned {
		wf.Assignee = &filter.Assignee
	}
	if filter.ParentID != "" {
		wf.ParentID = &filter.ParentID
	}
	if filter.MolType != "" {
		molType := types.MolType(filter.MolType)
		wf.MolType = &molType
	}

	issues, err := s.repo.GetReadyWork(ctx, wf)
	if err != nil {
		return nil, fmt.Errorf("failed to get ready work: %w", err)
	}

	return issues, nil
}

// GetBlocked retrieves issues that are blocked.
func (s *Service) GetBlocked(ctx context.Context, filter BlockedFilter) ([]*types.BlockedIssue, error) {
	var wf types.WorkFilter
	if filter.ParentID != "" {
		wf.ParentID = &filter.ParentID
	}

	issues, err := s.repo.GetBlockedIssues(ctx, wf)
	if err != nil {
		return nil, fmt.Errorf("failed to get blocked issues: %w", err)
	}

	return issues, nil
}

// IsBlocked checks if an issue has open blockers.
func (s *Service) IsBlocked(ctx context.Context, issueID string) (bool, []string, error) {
	if issueID == "" {
		return false, nil, fmt.Errorf("issue ID is required")
	}

	blocked, blockers, err := s.repo.IsBlocked(ctx, issueID)
	if err != nil {
		return false, nil, fmt.Errorf("failed to check blockers: %w", err)
	}

	return blocked, blockers, nil
}

// GetStale retrieves issues not updated in a given period.
func (s *Service) GetStale(ctx context.Context, filter StaleFilter) ([]*types.Issue, error) {
	sf := types.StaleFilter{
		Days:   filter.Days,
		Status: filter.Status,
		Limit:  filter.Limit,
	}

	issues, err := s.repo.GetStaleIssues(ctx, sf)
	if err != nil {
		return nil, fmt.Errorf("failed to get stale issues: %w", err)
	}

	return issues, nil
}

// GetEpicsEligibleForClosure retrieves epics that can be closed.
func (s *Service) GetEpicsEligibleForClosure(ctx context.Context, filter EpicStatusFilter) ([]*types.EpicStatus, error) {
	epics, err := s.repo.GetEpicsEligibleForClosure(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get epic status: %w", err)
	}

	if filter.EligibleOnly {
		filtered := []*types.EpicStatus{}
		for _, epic := range epics {
			if epic.EligibleForClose {
				filtered = append(filtered, epic)
			}
		}
		epics = filtered
	}

	return epics, nil
}

// GetNewlyUnblockedByClose retrieves issues unblocked by closing an issue.
func (s *Service) GetNewlyUnblockedByClose(ctx context.Context, closedIssueID string) ([]*types.Issue, error) {
	if closedIssueID == "" {
		return nil, fmt.Errorf("closed issue ID is required")
	}

	issues, err := s.repo.GetNewlyUnblockedByClose(ctx, closedIssueID)
	if err != nil {
		return nil, fmt.Errorf("failed to get newly unblocked issues: %w", err)
	}

	return issues, nil
}

// GetKanbanBoard retrieves kanban board data for a feature/epic.
func (s *Service) GetKanbanBoard(ctx context.Context, featureID string) (*KanbanResponse, error) {
	if featureID == "" {
		return nil, fmt.Errorf("feature ID is required")
	}

	// Get all issues under this feature
	issues, err := s.issueRepo.Search(ctx, "", types.IssueFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	// Filter to issues under this feature
	var featureIssues []*types.Issue
	for _, issue := range issues {
		if strings.HasPrefix(issue.ID, featureID+".") || issue.ID == featureID {
			featureIssues = append(featureIssues, issue)
		}
	}

	// Group by status into lanes
	lanes := map[string][]IssueCard{
		"planned":    {},
		"doing":      {},
		"for_review": {},
		"done":       {},
	}

	for _, issue := range featureIssues {
		// Skip the epic itself
		if issue.ID == featureID {
			continue
		}

		// Check if blocked
		isBlocked, blockers, _ := s.repo.IsBlocked(ctx, issue.ID)

		card := IssueCard{
			ID:          issue.ID,
			Title:       issue.Title,
			Priority:    issue.Priority,
			Status:      string(issue.Status),
			Assignee:    issue.Assignee,
			Description: truncate(issue.Description, 200),
			CreatedAt:   issue.CreatedAt,
			UpdatedAt:   issue.UpdatedAt,
			Labels:      issue.Labels,
			IsBlocked:   isBlocked,
			Blockers:    blockers,
		}

		lane := statusToLane(issue.Status)
		lanes[lane] = append(lanes[lane], card)
	}

	// Sort each lane by priority
	for _, cards := range lanes {
		sort.Slice(cards, func(i, j int) bool {
			return cards[i].Priority < cards[j].Priority
		})
	}

	return &KanbanResponse{
		Lanes:         lanes,
		IsLegacy:      false,
		UpgradeNeeded: false,
	}, nil
}

func statusToLane(status types.Status) string {
	switch status {
	case types.StatusOpen:
		return "planned"
	case types.StatusInProgress:
		return "doing"
	case types.StatusBlocked:
		return "doing"
	case types.StatusClosed:
		return "done"
	default:
		return "planned"
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

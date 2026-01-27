// Package epics provides the epic management vertical slice.
package epics

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/steveyegge/beads/internal/types"
)

// IssueRepository defines issue operations needed by epics.
type IssueRepository interface {
	Get(ctx context.Context, id string) (*types.Issue, error)
	Search(ctx context.Context, query string, filter types.IssueFilter) ([]*types.Issue, error)
}

// KanbanRepository defines kanban operations needed by epics.
type KanbanRepository interface {
	GetEpicsEligibleForClosure(ctx context.Context) ([]*types.EpicStatus, error)
	IsBlocked(ctx context.Context, issueID string) (bool, []string, error)
}

// Service provides business logic for epic operations.
type Service struct {
	issueRepo  IssueRepository
	kanbanRepo KanbanRepository
	rootDir    string
}

// NewService creates a new epic service.
func NewService(issueRepo IssueRepository, kanbanRepo KanbanRepository, rootDir string) *Service {
	return &Service{
		issueRepo:  issueRepo,
		kanbanRepo: kanbanRepo,
		rootDir:    rootDir,
	}
}

// GetStatus retrieves epic statuses.
func (s *Service) GetStatus(ctx context.Context, args StatusArgs) ([]*types.EpicStatus, error) {
	epics, err := s.kanbanRepo.GetEpicsEligibleForClosure(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get epic status: %w", err)
	}

	if args.EligibleOnly {
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

// GetEpicSummary retrieves a summary for a specific epic.
func (s *Service) GetEpicSummary(ctx context.Context, epicID string) (*EpicSummary, error) {
	if epicID == "" {
		return nil, fmt.Errorf("epic ID is required")
	}

	epic, err := s.issueRepo.Get(ctx, epicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get epic: %w", err)
	}
	if epic == nil {
		return nil, fmt.Errorf("epic not found: %s", epicID)
	}

	// Get all children
	issues, err := s.issueRepo.Search(ctx, "", types.IssueFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	var children []*types.Issue
	closedCount := 0
	blockedCount := 0

	for _, issue := range issues {
		if strings.HasPrefix(issue.ID, epicID+".") {
			children = append(children, issue)
			if issue.Status == types.StatusClosed {
				closedCount++
			}
			if blocked, _, _ := s.kanbanRepo.IsBlocked(ctx, issue.ID); blocked {
				blockedCount++
			}
		}
	}

	totalCount := len(children)
	openCount := totalCount - closedCount
	progress := 0.0
	if totalCount > 0 {
		progress = float64(closedCount) / float64(totalCount)
	}

	return &EpicSummary{
		ID:               epicID,
		Title:            epic.Title,
		Status:           string(epic.Status),
		TotalChildren:    totalCount,
		ClosedChildren:   closedCount,
		OpenChildren:     openCount,
		BlockedChildren:  blockedCount,
		Progress:         progress,
		EligibleForClose: openCount == 0 && totalCount > 0,
		Children:         children,
	}, nil
}

// ListFeatures retrieves all epics/features.
func (s *Service) ListFeatures(ctx context.Context) (*FeaturesResponse, error) {
	issues, err := s.issueRepo.Search(ctx, "", types.IssueFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	// Find epics (top-level issues or issues with type "epic")
	epics := make(map[string]*types.Issue)
	childrenByParent := make(map[string][]*types.Issue)

	for _, issue := range issues {
		if issue.IssueType == types.TypeEpic || !strings.Contains(issue.ID, ".") {
			epics[issue.ID] = issue
		} else {
			parts := strings.Split(issue.ID, ".")
			if len(parts) > 1 {
				parentID := strings.Join(parts[:len(parts)-1], ".")
				childrenByParent[parentID] = append(childrenByParent[parentID], issue)
			}
		}
	}

	// Build feature summaries
	var features []FeatureSummary
	for id, epic := range epics {
		children := childrenByParent[id]
		stats := computeKanbanStats(children)

		features = append(features, FeatureSummary{
			ID:          id,
			Name:        epic.Title,
			Path:        id,
			KanbanStats: stats,
			Artifacts: map[string]bool{
				"spec":       epic.Description != "",
				"plan":       epic.Design != "",
				"tasks":      len(children) > 0,
				"research":   false,
				"contracts":  false,
				"data_model": false,
				"checklists": epic.AcceptanceCriteria != "",
			},
			Meta: map[string]string{
				"mission": "software-dev",
			},
			IsLegacy: false,
		})
	}

	// Sort by ID descending (most recent first)
	sort.Slice(features, func(i, j int) bool {
		return features[i].ID > features[j].ID
	})

	return &FeaturesResponse{
		Features:    features,
		ProjectPath: s.rootDir,
		ActiveMission: map[string]string{
			"name":        "Kitty-Beads",
			"domain":      "software-dev",
			"version":     "1.0.0",
			"slug":        "kitty-beads",
			"description": "Beads issue tracker with Spec Kitty UI",
		},
	}, nil
}

func computeKanbanStats(issues []*types.Issue) map[string]int {
	stats := map[string]int{
		"planned":    0,
		"doing":      0,
		"for_review": 0,
		"done":       0,
	}

	for _, issue := range issues {
		lane := statusToLane(issue.Status)
		stats[lane]++
	}

	return stats
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

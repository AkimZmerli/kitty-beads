// Package issues provides the issue management vertical slice.
package issues

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/storage/sqlite"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/util"
)

// Service provides business logic for issue operations.
type Service struct {
	repo          Repository
	store         storage.Storage
	labelRepo     LabelRepository
	depRepo       DependencyRepository
	onMutation    func(MutationEvent)
}

// LabelRepository defines label operations needed by the service.
type LabelRepository interface {
	Add(ctx context.Context, issueID, label, actor string) error
	Remove(ctx context.Context, issueID, label, actor string) error
	Get(ctx context.Context, issueID string) ([]string, error)
	GetForIssues(ctx context.Context, issueIDs []string) (map[string][]string, error)
}

// DependencyRepository defines dependency operations needed by the service.
type DependencyRepository interface {
	Add(ctx context.Context, dep *types.Dependency, actor string) error
	Remove(ctx context.Context, issueID, dependsOnID string, actor string) error
	GetRecords(ctx context.Context, issueID string) ([]*types.Dependency, error)
	GetCounts(ctx context.Context, issueIDs []string) (map[string]*types.DependencyCounts, error)
}

// NewService creates a new issue service.
func NewService(repo Repository, store storage.Storage, labelRepo LabelRepository, depRepo DependencyRepository) *Service {
	return &Service{
		repo:      repo,
		store:     store,
		labelRepo: labelRepo,
		depRepo:   depRepo,
	}
}

// SetMutationHandler sets a callback for mutation events.
func (s *Service) SetMutationHandler(handler func(MutationEvent)) {
	s.onMutation = handler
}

func (s *Service) emitMutation(event MutationEvent) {
	event.Timestamp = time.Now()
	if s.onMutation != nil {
		s.onMutation(event)
	}
}

// Create creates a new issue with all related data.
func (s *Service) Create(ctx context.Context, args CreateArgs, actor string) (*types.Issue, error) {
	// Check for conflicting flags
	if args.ID != "" && args.Parent != "" {
		return nil, fmt.Errorf("cannot specify both ID and Parent")
	}

	// Warn if creating an issue without a description (unless it's a test issue)
	if args.Description == "" && !strings.Contains(strings.ToLower(args.Title), "test") {
		fmt.Fprintf(os.Stderr, "[WARNING] Creating issue '%s' without description. Issues without descriptions lack context for future work.\n", args.Title)
	}

	// If parent is specified, generate child ID
	issueID := args.ID
	if args.Parent != "" {
		childID, err := s.repo.GetNextChildID(ctx, args.Parent)
		if err != nil {
			return nil, fmt.Errorf("failed to generate child ID: %w", err)
		}
		issueID = childID
	}

	// Parse DueAt if provided
	var dueAt *time.Time
	if args.DueAt != "" {
		t, err := parseTime(args.DueAt)
		if err != nil {
			return nil, fmt.Errorf("invalid due_at format %q. Examples: 2025-01-15, 2025-01-15T10:00:00Z", args.DueAt)
		}
		dueAt = &t
	}

	// Parse DeferUntil if provided
	var deferUntil *time.Time
	if args.DeferUntil != "" {
		t, err := parseTime(args.DeferUntil)
		if err != nil {
			return nil, fmt.Errorf("invalid defer_until format %q. Examples: 2025-01-15, 2025-01-15T10:00:00Z", args.DeferUntil)
		}
		deferUntil = &t
	}

	issue := &types.Issue{
		ID:                 issueID,
		Title:              args.Title,
		Description:        args.Description,
		IssueType:          types.IssueType(args.IssueType),
		Priority:           args.Priority,
		Design:             args.Design,
		AcceptanceCriteria: args.AcceptanceCriteria,
		Notes:              args.Notes,
		Assignee:           args.Assignee,
		EstimatedMinutes:   args.EstimatedMinutes,
		Status:             types.StatusOpen,
		Sender:             args.Sender,
		Ephemeral:          args.Ephemeral,
		IDPrefix:           args.IDPrefix,
		CreatedBy:          args.CreatedBy,
		Owner:              args.Owner,
		MolType:            types.MolType(args.MolType),
		RoleType:           args.RoleType,
		Rig:                args.Rig,
		EventKind:          args.EventCategory,
		Actor:              args.EventActor,
		Target:             args.EventTarget,
		Payload:            args.EventPayload,
		DueAt:              dueAt,
		DeferUntil:         deferUntil,
	}

	if args.ExternalRef != "" {
		issue.ExternalRef = &args.ExternalRef
	}

	// Check if any dependencies are discovered-from type
	// If so, inherit source_repo from the parent issue
	var discoveredFromParentID string
	for _, depSpec := range args.Dependencies {
		depSpec = strings.TrimSpace(depSpec)
		if depSpec == "" {
			continue
		}
		if strings.Contains(depSpec, ":") {
			parts := strings.SplitN(depSpec, ":", 2)
			if len(parts) == 2 {
				depType := types.DependencyType(strings.TrimSpace(parts[0]))
				dependsOnID := strings.TrimSpace(parts[1])
				if depType == types.DepDiscoveredFrom {
					discoveredFromParentID = dependsOnID
					break
				}
			}
		}
	}

	// If we found a discovered-from dependency, inherit source_repo from parent
	if discoveredFromParentID != "" {
		parentIssue, err := s.repo.Get(ctx, discoveredFromParentID)
		if err == nil && parentIssue.SourceRepo != "" {
			issue.SourceRepo = parentIssue.SourceRepo
		}
	}

	if err := s.repo.Create(ctx, issue, actor); err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	// If parent was specified, add parent-child dependency
	if args.Parent != "" {
		dep := &types.Dependency{
			IssueID:     issue.ID,
			DependsOnID: args.Parent,
			Type:        types.DepParentChild,
		}
		if err := s.depRepo.Add(ctx, dep, actor); err != nil {
			return nil, fmt.Errorf("failed to add parent-child dependency %s -> %s: %w", issue.ID, args.Parent, err)
		}
	}

	// If RepliesTo was specified, add replies-to dependency
	if args.RepliesTo != "" {
		dep := &types.Dependency{
			IssueID:     issue.ID,
			DependsOnID: args.RepliesTo,
			Type:        types.DepRepliesTo,
			ThreadID:    args.RepliesTo,
		}
		if err := s.depRepo.Add(ctx, dep, actor); err != nil {
			return nil, fmt.Errorf("failed to add replies-to dependency %s -> %s: %w", issue.ID, args.RepliesTo, err)
		}
	}

	// Add labels if specified
	for _, label := range args.Labels {
		if err := s.labelRepo.Add(ctx, issue.ID, label, actor); err != nil {
			return nil, fmt.Errorf("failed to add label %s: %w", label, err)
		}
	}

	// Auto-add role_type/rig labels for agent beads
	if containsLabel(args.Labels, "gt:agent") {
		if issue.RoleType != "" {
			label := "role_type:" + issue.RoleType
			if err := s.labelRepo.Add(ctx, issue.ID, label, actor); err != nil {
				return nil, fmt.Errorf("failed to add role_type label: %w", err)
			}
		}
		if issue.Rig != "" {
			label := "rig:" + issue.Rig
			if err := s.labelRepo.Add(ctx, issue.ID, label, actor); err != nil {
				return nil, fmt.Errorf("failed to add rig label: %w", err)
			}
		}
	}

	// Add dependencies if specified
	for _, depSpec := range args.Dependencies {
		depSpec = strings.TrimSpace(depSpec)
		if depSpec == "" {
			continue
		}

		var depType types.DependencyType
		var dependsOnID string

		if strings.Contains(depSpec, ":") {
			parts := strings.SplitN(depSpec, ":", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid dependency format '%s', expected 'type:id' or 'id'", depSpec)
			}
			depType = types.DependencyType(strings.TrimSpace(parts[0]))
			dependsOnID = strings.TrimSpace(parts[1])
		} else {
			depType = types.DepBlocks
			dependsOnID = depSpec
		}

		if !depType.IsValid() {
			return nil, fmt.Errorf("invalid dependency type '%s' (valid: blocks, related, parent-child, discovered-from)", depType)
		}

		dep := &types.Dependency{
			IssueID:     issue.ID,
			DependsOnID: dependsOnID,
			Type:        depType,
		}
		if err := s.depRepo.Add(ctx, dep, actor); err != nil {
			return nil, fmt.Errorf("failed to add dependency %s -> %s: %w", issue.ID, dependsOnID, err)
		}
	}

	// Add waits-for dependency if specified
	if args.WaitsFor != "" {
		gate := args.WaitsForGate
		if gate == "" {
			gate = types.WaitsForAllChildren
		}
		if gate != types.WaitsForAllChildren && gate != types.WaitsForAnyChildren {
			return nil, fmt.Errorf("invalid waits_for_gate value '%s' (valid: all-children, any-children)", gate)
		}

		meta := types.WaitsForMeta{Gate: gate}
		metaJSON, err := json.Marshal(meta)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize waits-for metadata: %w", err)
		}

		dep := &types.Dependency{
			IssueID:     issue.ID,
			DependsOnID: args.WaitsFor,
			Type:        types.DepWaitsFor,
			Metadata:    string(metaJSON),
		}
		if err := s.depRepo.Add(ctx, dep, actor); err != nil {
			return nil, fmt.Errorf("failed to add waits-for dependency %s -> %s: %w", issue.ID, args.WaitsFor, err)
		}
	}

	s.emitMutation(MutationEvent{
		Type:     MutationCreate,
		IssueID:  issue.ID,
		Title:    issue.Title,
		Assignee: issue.Assignee,
	})

	return issue, nil
}

// Update updates an existing issue.
func (s *Service) Update(ctx context.Context, args UpdateArgs, actor string) (*types.Issue, error) {
	// Check if issue exists and is not a template
	issue, err := s.repo.Get(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	if issue == nil {
		return nil, fmt.Errorf("issue %s not found", args.ID)
	}
	if issue.IsTemplate {
		return nil, fmt.Errorf("cannot update template %s: templates are read-only; use 'bd molecule instantiate' to create a work item", args.ID)
	}

	// Handle claim operation atomically
	if args.Claim {
		if issue.Assignee != "" {
			return nil, fmt.Errorf("already claimed by %s", issue.Assignee)
		}
		claimUpdates := map[string]interface{}{
			"assignee": actor,
			"status":   "in_progress",
		}
		if err := s.repo.Update(ctx, args.ID, claimUpdates, actor); err != nil {
			return nil, fmt.Errorf("failed to claim issue: %w", err)
		}
	}

	updates, err := updatesFromArgs(args)
	if err != nil {
		return nil, err
	}

	// Apply regular field updates if any
	if len(updates) > 0 {
		if err := s.repo.Update(ctx, args.ID, updates, actor); err != nil {
			return nil, fmt.Errorf("failed to update issue: %w", err)
		}
	}

	// Handle label operations
	if len(args.SetLabels) > 0 {
		currentLabels, err := s.labelRepo.Get(ctx, args.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get current labels: %w", err)
		}
		for _, label := range currentLabels {
			if err := s.labelRepo.Remove(ctx, args.ID, label, actor); err != nil {
				return nil, fmt.Errorf("failed to remove label %s: %w", label, err)
			}
		}
		for _, label := range args.SetLabels {
			if err := s.labelRepo.Add(ctx, args.ID, label, actor); err != nil {
				return nil, fmt.Errorf("failed to set label %s: %w", label, err)
			}
		}
	}

	for _, label := range args.AddLabels {
		if err := s.labelRepo.Add(ctx, args.ID, label, actor); err != nil {
			return nil, fmt.Errorf("failed to add label %s: %w", label, err)
		}
	}

	for _, label := range args.RemoveLabels {
		if err := s.labelRepo.Remove(ctx, args.ID, label, actor); err != nil {
			return nil, fmt.Errorf("failed to remove label %s: %w", label, err)
		}
	}

	// Auto-add role_type/rig labels for agent beads
	issueLabels, _ := s.labelRepo.Get(ctx, args.ID)
	if containsLabel(issueLabels, "gt:agent") {
		if args.RoleType != nil && *args.RoleType != "" {
			for _, l := range issueLabels {
				if strings.HasPrefix(l, "role_type:") {
					_ = s.labelRepo.Remove(ctx, args.ID, l, actor)
				}
			}
			label := "role_type:" + *args.RoleType
			if err := s.labelRepo.Add(ctx, args.ID, label, actor); err != nil {
				return nil, fmt.Errorf("failed to add role_type label: %w", err)
			}
		}
		if args.Rig != nil && *args.Rig != "" {
			for _, l := range issueLabels {
				if strings.HasPrefix(l, "rig:") {
					_ = s.labelRepo.Remove(ctx, args.ID, l, actor)
				}
			}
			label := "rig:" + *args.Rig
			if err := s.labelRepo.Add(ctx, args.ID, label, actor); err != nil {
				return nil, fmt.Errorf("failed to add rig label: %w", err)
			}
		}
	}

	// Handle reparenting
	if args.Parent != nil {
		newParentID := *args.Parent

		if newParentID != "" {
			newParent, err := s.repo.Get(ctx, newParentID)
			if err != nil {
				return nil, fmt.Errorf("failed to get new parent: %w", err)
			}
			if newParent == nil {
				return nil, fmt.Errorf("parent issue %s not found", newParentID)
			}
		}

		deps, err := s.depRepo.GetRecords(ctx, args.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependencies: %w", err)
		}
		for _, dep := range deps {
			if dep.Type == types.DepParentChild {
				if err := s.depRepo.Remove(ctx, args.ID, dep.DependsOnID, actor); err != nil {
					return nil, fmt.Errorf("failed to remove old parent dependency: %w", err)
				}
				break
			}
		}

		if newParentID != "" {
			newDep := &types.Dependency{
				IssueID:     args.ID,
				DependsOnID: newParentID,
				Type:        types.DepParentChild,
			}
			if err := s.depRepo.Add(ctx, newDep, actor); err != nil {
				return nil, fmt.Errorf("failed to add parent dependency: %w", err)
			}
		}
	}

	// Emit mutation event
	if len(updates) > 0 || len(args.SetLabels) > 0 || len(args.AddLabels) > 0 || len(args.RemoveLabels) > 0 || args.Parent != nil {
		effectiveAssignee := issue.Assignee
		if args.Assignee != nil && *args.Assignee != "" {
			effectiveAssignee = *args.Assignee
		}

		if args.Status != nil && *args.Status != string(issue.Status) {
			s.emitMutation(MutationEvent{
				Type:      MutationStatus,
				IssueID:   args.ID,
				Title:     issue.Title,
				Assignee:  effectiveAssignee,
				Actor:     actor,
				OldStatus: string(issue.Status),
				NewStatus: *args.Status,
			})
		} else {
			s.emitMutation(MutationEvent{
				Type:     MutationUpdate,
				IssueID:  args.ID,
				Title:    issue.Title,
				Assignee: effectiveAssignee,
				Actor:    actor,
			})
		}
	}

	updatedIssue, err := s.repo.Get(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated issue: %w", err)
	}

	return updatedIssue, nil
}

// Close closes an issue.
func (s *Service) Close(ctx context.Context, args CloseArgs, actor string) (*CloseResult, error) {
	issue, err := s.repo.Get(ctx, args.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	if issue != nil && issue.IsTemplate {
		return nil, fmt.Errorf("cannot close template %s: templates are read-only", args.ID)
	}

	// Check if issue has open blockers
	if !args.Force {
		blocked, blockers, err := s.store.IsBlocked(ctx, args.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check blockers: %w", err)
		}
		if blocked && len(blockers) > 0 {
			return nil, fmt.Errorf("cannot close %s: blocked by open issues %v (use --force to override)", args.ID, blockers)
		}
	}

	oldStatus := ""
	if issue != nil {
		oldStatus = string(issue.Status)
	}

	if err := s.repo.Close(ctx, args.ID, args.Reason, actor, args.Session); err != nil {
		return nil, fmt.Errorf("failed to close issue: %w", err)
	}

	s.emitMutation(MutationEvent{
		Type:      MutationStatus,
		IssueID:   args.ID,
		Title:     issue.Title,
		Assignee:  issue.Assignee,
		OldStatus: oldStatus,
		NewStatus: "closed",
	})

	closedIssue, _ := s.repo.Get(ctx, args.ID)

	if args.SuggestNext {
		unblocked, err := s.store.GetNewlyUnblockedByClose(ctx, args.ID)
		if err != nil {
			unblocked = nil
		}
		return &CloseResult{
			Closed:    closedIssue,
			Unblocked: unblocked,
		}, nil
	}

	return &CloseResult{Closed: closedIssue}, nil
}

// Delete deletes one or more issues.
func (s *Service) Delete(ctx context.Context, args DeleteArgs, actor string) (*DeleteResult, error) {
	if len(args.IDs) == 0 {
		return nil, fmt.Errorf("no issue IDs provided for deletion")
	}

	// Use batch delete for cascade/multi-issue operations on SQLite storage
	if sqlStore, ok := s.store.(*sqlite.SQLiteStorage); ok {
		useBatchDelete := args.Cascade || args.Force || len(args.IDs) > 1 || args.DryRun
		if useBatchDelete {
			result, err := sqlStore.DeleteIssues(ctx, args.IDs, args.Cascade, args.Force, args.DryRun)
			if err != nil {
				return nil, fmt.Errorf("delete failed: %w", err)
			}

			if !args.DryRun {
				for _, issueID := range args.IDs {
					s.emitMutation(MutationEvent{
						Type:    MutationDelete,
						IssueID: issueID,
					})
				}
			}

			return &DeleteResult{
				DeletedCount:        result.DeletedCount,
				TotalCount:          len(args.IDs),
				DryRun:              args.DryRun,
				IssueCount:          result.DeletedCount,
				DependenciesRemoved: result.DependenciesCount,
				LabelsRemoved:       result.LabelsCount,
				EventsRemoved:       result.EventsCount,
				OrphanedIssues:      result.OrphanedIssues,
			}, nil
		}
	}

	// Simple single-issue delete path
	if args.DryRun {
		return &DeleteResult{
			DryRun:     true,
			IssueCount: len(args.IDs),
			TotalCount: len(args.IDs),
		}, nil
	}

	deletedCount := 0
	errors := make([]string, 0)

	for _, issueID := range args.IDs {
		issue, err := s.repo.Get(ctx, issueID)
		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", issueID, err))
			continue
		}
		if issue == nil {
			errors = append(errors, fmt.Sprintf("%s: not found", issueID))
			continue
		}
		if issue.IsTemplate {
			errors = append(errors, fmt.Sprintf("%s: cannot delete template (templates are read-only)", issueID))
			continue
		}

		// Create tombstone instead of hard delete
		type tombstoner interface {
			CreateTombstone(ctx context.Context, id string, actor string, reason string) error
		}
		if t, ok := s.store.(tombstoner); ok {
			reason := args.Reason
			if reason == "" {
				reason = "deleted via daemon"
			}
			if err := t.CreateTombstone(ctx, issueID, "daemon", reason); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", issueID, err))
				continue
			}
		} else {
			if err := s.repo.Delete(ctx, issueID); err != nil {
				errors = append(errors, fmt.Sprintf("%s: %v", issueID, err))
				continue
			}
		}

		s.emitMutation(MutationEvent{
			Type:     MutationDelete,
			IssueID:  issueID,
			Title:    issue.Title,
			Assignee: issue.Assignee,
		})
		deletedCount++
	}

	result := &DeleteResult{
		DeletedCount: deletedCount,
		TotalCount:   len(args.IDs),
	}

	if len(errors) > 0 {
		result.Errors = errors
		if deletedCount == 0 {
			return nil, fmt.Errorf("failed to delete all issues: %v", errors)
		}
		result.PartialSuccess = true
	}

	return result, nil
}

// List lists issues with optional filtering.
func (s *Service) List(ctx context.Context, filter ListFilter) ([]*types.IssueWithCounts, error) {
	issueFilter := buildIssueFilter(filter)

	// Guard against excessive ID lists
	const maxIDs = 1000
	if len(issueFilter.IDs) > maxIDs {
		return nil, fmt.Errorf("--id flag supports at most %d issue IDs, got %d", maxIDs, len(issueFilter.IDs))
	}

	issues, err := s.repo.Search(ctx, filter.Query, issueFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list issues: %w", err)
	}

	// Populate labels for each issue
	for _, issue := range issues {
		labels, _ := s.labelRepo.Get(ctx, issue.ID)
		issue.Labels = labels
	}

	// Get dependency counts in bulk
	issueIDs := make([]string, len(issues))
	for i, issue := range issues {
		issueIDs[i] = issue.ID
	}
	depCounts, _ := s.depRepo.GetCounts(ctx, issueIDs)

	// Build response with counts
	issuesWithCounts := make([]*types.IssueWithCounts, len(issues))
	for i, issue := range issues {
		counts := depCounts[issue.ID]
		if counts == nil {
			counts = &types.DependencyCounts{DependencyCount: 0, DependentCount: 0}
		}
		issuesWithCounts[i] = &types.IssueWithCounts{
			Issue:           issue,
			DependencyCount: counts.DependencyCount,
			DependentCount:  counts.DependentCount,
		}
	}

	return issuesWithCounts, nil
}

// Count counts issues matching the filter.
func (s *Service) Count(ctx context.Context, filter CountFilter) (interface{}, error) {
	issueFilter := types.IssueFilter{}

	if filter.Status != "" && filter.Status != "all" {
		status := types.Status(filter.Status)
		issueFilter.Status = &status
	}
	if filter.IssueType != "" {
		issueType := types.IssueType(filter.IssueType)
		issueFilter.IssueType = &issueType
	}
	if filter.Assignee != "" {
		issueFilter.Assignee = &filter.Assignee
	}
	if filter.Priority != nil {
		issueFilter.Priority = filter.Priority
	}

	labels := util.NormalizeLabels(filter.Labels)
	labelsAny := util.NormalizeLabels(filter.LabelsAny)
	if len(labels) > 0 {
		issueFilter.Labels = labels
	}
	if len(labelsAny) > 0 {
		issueFilter.LabelsAny = labelsAny
	}
	if len(filter.IDs) > 0 {
		ids := util.NormalizeLabels(filter.IDs)
		if len(ids) > 0 {
			issueFilter.IDs = ids
		}
	}

	issueFilter.TitleContains = filter.TitleContains
	issueFilter.DescriptionContains = filter.DescriptionContains
	issueFilter.NotesContains = filter.NotesContains
	issueFilter.EmptyDescription = filter.EmptyDescription
	issueFilter.NoAssignee = filter.NoAssignee
	issueFilter.NoLabels = filter.NoLabels
	issueFilter.PriorityMin = filter.PriorityMin
	issueFilter.PriorityMax = filter.PriorityMax

	// Parse date filters
	if filter.CreatedAfter != "" {
		t, err := parseTime(filter.CreatedAfter)
		if err != nil {
			return nil, fmt.Errorf("invalid --created-after date: %w", err)
		}
		issueFilter.CreatedAfter = &t
	}
	if filter.CreatedBefore != "" {
		t, err := parseTime(filter.CreatedBefore)
		if err != nil {
			return nil, fmt.Errorf("invalid --created-before date: %w", err)
		}
		issueFilter.CreatedBefore = &t
	}
	if filter.UpdatedAfter != "" {
		t, err := parseTime(filter.UpdatedAfter)
		if err != nil {
			return nil, fmt.Errorf("invalid --updated-after date: %w", err)
		}
		issueFilter.UpdatedAfter = &t
	}
	if filter.UpdatedBefore != "" {
		t, err := parseTime(filter.UpdatedBefore)
		if err != nil {
			return nil, fmt.Errorf("invalid --updated-before date: %w", err)
		}
		issueFilter.UpdatedBefore = &t
	}
	if filter.ClosedAfter != "" {
		t, err := parseTime(filter.ClosedAfter)
		if err != nil {
			return nil, fmt.Errorf("invalid --closed-after date: %w", err)
		}
		issueFilter.ClosedAfter = &t
	}
	if filter.ClosedBefore != "" {
		t, err := parseTime(filter.ClosedBefore)
		if err != nil {
			return nil, fmt.Errorf("invalid --closed-before date: %w", err)
		}
		issueFilter.ClosedBefore = &t
	}

	issues, err := s.repo.Search(ctx, filter.Query, issueFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to count issues: %w", err)
	}

	// If no grouping, just return the count
	if filter.GroupBy == "" {
		return &CountResult{Count: len(issues)}, nil
	}

	counts := make(map[string]int)

	// For label grouping, fetch all labels in one query
	var labelsMap map[string][]string
	if filter.GroupBy == "label" {
		issueIDs := make([]string, len(issues))
		for i, issue := range issues {
			issueIDs[i] = issue.ID
		}
		labelsMap, err = s.labelRepo.GetForIssues(ctx, issueIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get labels: %w", err)
		}
	}

	for _, issue := range issues {
		var groupKey string
		switch filter.GroupBy {
		case "status":
			groupKey = string(issue.Status)
		case "priority":
			groupKey = fmt.Sprintf("P%d", issue.Priority)
		case "type":
			groupKey = string(issue.IssueType)
		case "assignee":
			if issue.Assignee == "" {
				groupKey = "(unassigned)"
			} else {
				groupKey = issue.Assignee
			}
		case "label":
			issueLabels := labelsMap[issue.ID]
			if len(issueLabels) > 0 {
				for _, label := range issueLabels {
					counts[label]++
				}
				continue
			} else {
				groupKey = "(no labels)"
			}
		default:
			return nil, fmt.Errorf("invalid group_by value: %s (must be one of: status, priority, type, assignee, label)", filter.GroupBy)
		}
		counts[groupKey]++
	}

	groups := make([]GroupCount, 0, len(counts))
	for group, count := range counts {
		groups = append(groups, GroupCount{Group: group, Count: count})
	}

	return &GroupedCountResult{
		Total:  len(issues),
		Groups: groups,
	}, nil
}

// Show retrieves detailed information about an issue.
func (s *Service) Show(ctx context.Context, id string) (*types.IssueDetails, error) {
	issue, err := s.repo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get issue: %w", err)
	}
	if issue == nil {
		return nil, fmt.Errorf("issue not found: %s", id)
	}

	labels, _ := s.labelRepo.Get(ctx, issue.ID)

	// Get dependencies and dependents with metadata
	var deps []*types.IssueWithDependencyMetadata
	var dependents []*types.IssueWithDependencyMetadata

	if sqliteStore, ok := s.store.(*sqlite.SQLiteStorage); ok {
		deps, _ = sqliteStore.GetDependenciesWithMetadata(ctx, issue.ID)
		dependents, _ = sqliteStore.GetDependentsWithMetadata(ctx, issue.ID)
	} else {
		regularDeps, _ := s.store.GetDependencies(ctx, issue.ID)
		for _, d := range regularDeps {
			deps = append(deps, &types.IssueWithDependencyMetadata{
				Issue:          *d,
				DependencyType: types.DepBlocks,
			})
		}
		regularDependents, _ := s.store.GetDependents(ctx, issue.ID)
		for _, d := range regularDependents {
			dependents = append(dependents, &types.IssueWithDependencyMetadata{
				Issue:          *d,
				DependencyType: types.DepBlocks,
			})
		}
	}

	comments, _ := s.store.GetIssueComments(ctx, issue.ID)

	return &types.IssueDetails{
		Issue:        *issue,
		Labels:       labels,
		Dependencies: deps,
		Dependents:   dependents,
		Comments:     comments,
	}, nil
}

// GetReady retrieves ready work items.
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

	return s.store.GetReadyWork(ctx, wf)
}

// GetBlocked retrieves blocked issues.
func (s *Service) GetBlocked(ctx context.Context, filter BlockedFilter) ([]*types.BlockedIssue, error) {
	var wf types.WorkFilter
	if filter.ParentID != "" {
		wf.ParentID = &filter.ParentID
	}
	return s.store.GetBlockedIssues(ctx, wf)
}

// GetStale retrieves stale issues.
func (s *Service) GetStale(ctx context.Context, filter StaleFilter) ([]*types.Issue, error) {
	return s.store.GetStaleIssues(ctx, types.StaleFilter{
		Days:   filter.Days,
		Status: filter.Status,
		Limit:  filter.Limit,
	})
}

// GetStats retrieves issue statistics.
func (s *Service) GetStats(ctx context.Context) (*types.Statistics, error) {
	return s.store.GetStatistics(ctx)
}

// GetEpicStatus retrieves epic statuses.
func (s *Service) GetEpicStatus(ctx context.Context, filter EpicStatusFilter) ([]*types.EpicStatus, error) {
	epics, err := s.store.GetEpicsEligibleForClosure(ctx)
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

// Helper functions

func containsLabel(labels []string, label string) bool {
	for _, l := range labels {
		if l == label {
			return true
		}
	}
	return false
}

func parseTime(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("unsupported date format: %q (use YYYY-MM-DD or RFC3339)", s)
}

func updatesFromArgs(a UpdateArgs) (map[string]interface{}, error) {
	u := map[string]interface{}{}
	if a.Title != nil {
		u["title"] = *a.Title
	}
	if a.Description != nil {
		u["description"] = *a.Description
	}
	if a.Status != nil {
		u["status"] = *a.Status
	}
	if a.Priority != nil {
		u["priority"] = *a.Priority
	}
	if a.Design != nil {
		u["design"] = *a.Design
	}
	if a.AcceptanceCriteria != nil {
		u["acceptance_criteria"] = *a.AcceptanceCriteria
	}
	if a.Notes != nil {
		u["notes"] = *a.Notes
	}
	if a.Assignee != nil {
		u["assignee"] = *a.Assignee
	}
	if a.ExternalRef != nil {
		u["external_ref"] = *a.ExternalRef
	}
	if a.EstimatedMinutes != nil {
		u["estimated_minutes"] = *a.EstimatedMinutes
	}
	if a.IssueType != nil {
		u["issue_type"] = *a.IssueType
	}
	if a.Sender != nil {
		u["sender"] = *a.Sender
	}
	if a.Ephemeral != nil {
		u["wisp"] = *a.Ephemeral
	}
	if a.RepliesTo != nil {
		u["replies_to"] = *a.RepliesTo
	}
	if a.RelatesTo != nil {
		u["relates_to"] = *a.RelatesTo
	}
	if a.DuplicateOf != nil {
		u["duplicate_of"] = *a.DuplicateOf
	}
	if a.SupersededBy != nil {
		u["superseded_by"] = *a.SupersededBy
	}
	if a.Pinned != nil {
		u["pinned"] = *a.Pinned
	}
	if a.HookBead != nil {
		u["hook_bead"] = *a.HookBead
	}
	if a.RoleBead != nil {
		u["role_bead"] = *a.RoleBead
	}
	if a.AgentState != nil {
		u["agent_state"] = *a.AgentState
	}
	if a.LastActivity != nil && *a.LastActivity {
		u["last_activity"] = time.Now()
	}
	if a.RoleType != nil {
		u["role_type"] = *a.RoleType
	}
	if a.Rig != nil {
		u["rig"] = *a.Rig
	}
	if a.EventCategory != nil {
		u["event_category"] = *a.EventCategory
	}
	if a.EventActor != nil {
		u["event_actor"] = *a.EventActor
	}
	if a.EventTarget != nil {
		u["event_target"] = *a.EventTarget
	}
	if a.EventPayload != nil {
		u["event_payload"] = *a.EventPayload
	}
	if a.AwaitID != nil {
		u["await_id"] = *a.AwaitID
	}
	if len(a.Waiters) > 0 {
		u["waiters"] = a.Waiters
	}
	if a.Holder != nil {
		u["holder"] = *a.Holder
	}
	if a.DueAt != nil {
		if *a.DueAt == "" {
			u["due_at"] = nil
		} else {
			if t, err := time.ParseInLocation("2006-01-02", *a.DueAt, time.Local); err == nil {
				u["due_at"] = t
			} else if t, err := time.Parse(time.RFC3339, *a.DueAt); err == nil {
				u["due_at"] = t
			} else {
				return nil, fmt.Errorf("invalid due_at format %q: use YYYY-MM-DD or RFC3339", *a.DueAt)
			}
		}
	}
	if a.DeferUntil != nil {
		if *a.DeferUntil == "" {
			u["defer_until"] = nil
		} else {
			if t, err := time.ParseInLocation("2006-01-02", *a.DeferUntil, time.Local); err == nil {
				u["defer_until"] = t
			} else if t, err := time.Parse(time.RFC3339, *a.DeferUntil); err == nil {
				u["defer_until"] = t
			} else {
				return nil, fmt.Errorf("invalid defer_until format %q: use YYYY-MM-DD or RFC3339", *a.DeferUntil)
			}
		}
	}
	return u, nil
}

func buildIssueFilter(filter ListFilter) types.IssueFilter {
	issueFilter := types.IssueFilter{
		Limit: filter.Limit,
	}

	if filter.Status != "" && filter.Status != "all" {
		status := types.Status(filter.Status)
		issueFilter.Status = &status
	}
	if filter.IssueType != "" {
		issueType := types.IssueType(filter.IssueType)
		issueFilter.IssueType = &issueType
	}
	if filter.Assignee != "" {
		issueFilter.Assignee = &filter.Assignee
	}
	if filter.Priority != nil {
		issueFilter.Priority = filter.Priority
	}

	labels := util.NormalizeLabels(filter.Labels)
	labelsAny := util.NormalizeLabels(filter.LabelsAny)
	if len(labels) > 0 {
		issueFilter.Labels = labels
	} else if filter.Label != "" {
		issueFilter.Labels = []string{strings.TrimSpace(filter.Label)}
	}
	if len(labelsAny) > 0 {
		issueFilter.LabelsAny = labelsAny
	}
	if len(filter.IDs) > 0 {
		ids := util.NormalizeLabels(filter.IDs)
		if len(ids) > 0 {
			issueFilter.IDs = ids
		}
	}

	issueFilter.TitleContains = filter.TitleContains
	issueFilter.DescriptionContains = filter.DescriptionContains
	issueFilter.NotesContains = filter.NotesContains
	issueFilter.EmptyDescription = filter.EmptyDescription
	issueFilter.NoAssignee = filter.NoAssignee
	issueFilter.NoLabels = filter.NoLabels
	issueFilter.PriorityMin = filter.PriorityMin
	issueFilter.PriorityMax = filter.PriorityMax
	issueFilter.Pinned = filter.Pinned
	issueFilter.Ephemeral = filter.Ephemeral
	issueFilter.Deferred = filter.Deferred
	issueFilter.Overdue = filter.Overdue

	if !filter.IncludeTemplates {
		isTemplate := false
		issueFilter.IsTemplate = &isTemplate
	}

	if filter.ParentID != "" {
		issueFilter.ParentID = &filter.ParentID
	}

	if filter.MolType != "" {
		molType := types.MolType(filter.MolType)
		issueFilter.MolType = &molType
	}

	if len(filter.ExcludeStatus) > 0 {
		for _, s := range filter.ExcludeStatus {
			issueFilter.ExcludeStatus = append(issueFilter.ExcludeStatus, types.Status(s))
		}
	}

	if len(filter.ExcludeTypes) > 0 {
		for _, t := range filter.ExcludeTypes {
			issueFilter.ExcludeTypes = append(issueFilter.ExcludeTypes, types.IssueType(t))
		}
	}

	// Parse date filters
	if filter.CreatedAfter != "" {
		if t, err := parseTime(filter.CreatedAfter); err == nil {
			issueFilter.CreatedAfter = &t
		}
	}
	if filter.CreatedBefore != "" {
		if t, err := parseTime(filter.CreatedBefore); err == nil {
			issueFilter.CreatedBefore = &t
		}
	}
	if filter.UpdatedAfter != "" {
		if t, err := parseTime(filter.UpdatedAfter); err == nil {
			issueFilter.UpdatedAfter = &t
		}
	}
	if filter.UpdatedBefore != "" {
		if t, err := parseTime(filter.UpdatedBefore); err == nil {
			issueFilter.UpdatedBefore = &t
		}
	}
	if filter.ClosedAfter != "" {
		if t, err := parseTime(filter.ClosedAfter); err == nil {
			issueFilter.ClosedAfter = &t
		}
	}
	if filter.ClosedBefore != "" {
		if t, err := parseTime(filter.ClosedBefore); err == nil {
			issueFilter.ClosedBefore = &t
		}
	}
	if filter.DeferAfter != "" {
		if t, err := parseTime(filter.DeferAfter); err == nil {
			issueFilter.DeferAfter = &t
		}
	}
	if filter.DeferBefore != "" {
		if t, err := parseTime(filter.DeferBefore); err == nil {
			issueFilter.DeferBefore = &t
		}
	}
	if filter.DueAfter != "" {
		if t, err := parseTime(filter.DueAfter); err == nil {
			issueFilter.DueAfter = &t
		}
	}
	if filter.DueBefore != "" {
		if t, err := parseTime(filter.DueBefore); err == nil {
			issueFilter.DueBefore = &t
		}
	}

	return issueFilter
}

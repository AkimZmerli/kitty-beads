package template

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/utils"
)

// LoadSubgraph loads a template epic and all its descendants
func LoadSubgraph(ctx context.Context, s storage.Storage, templateID string) (*Subgraph, error) {
	if s == nil {
		return nil, fmt.Errorf("no database connection")
	}

	// Get the root issue
	root, err := s.GetIssue(ctx, templateID)
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}
	if root == nil {
		return nil, fmt.Errorf("template %s not found", templateID)
	}

	subgraph := &Subgraph{
		Root:     root,
		Issues:   []*types.Issue{root},
		IssueMap: map[string]*types.Issue{root.ID: root},
	}

	// Recursively load all children
	if err := loadDescendants(ctx, s, subgraph, root.ID); err != nil {
		return nil, err
	}

	// Load all dependencies within the subgraph
	for _, issue := range subgraph.Issues {
		deps, err := s.GetDependencyRecords(ctx, issue.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get dependencies for %s: %w", issue.ID, err)
		}
		for _, dep := range deps {
			// Only include dependencies where both ends are in the subgraph
			if _, ok := subgraph.IssueMap[dep.DependsOnID]; ok {
				subgraph.Dependencies = append(subgraph.Dependencies, dep)
			}
		}
	}

	return subgraph, nil
}

// loadDescendants recursively loads all child issues
// It uses two strategies to find children:
// 1. Check dependency records for parent-child relationships
// 2. Check for hierarchical IDs (parent.N) to catch children with missing/wrong deps
func loadDescendants(ctx context.Context, s storage.Storage, subgraph *Subgraph, parentID string) error {
	// Track children we've already added to avoid duplicates
	addedChildren := make(map[string]bool)

	// Strategy 1: GetDependents returns issues that depend on parentID
	dependents, err := s.GetDependents(ctx, parentID)
	if err != nil {
		return fmt.Errorf("failed to get dependents of %s: %w", parentID, err)
	}

	// Check each dependent to see if it's a child (has parent-child relationship)
	for _, dependent := range dependents {
		if _, exists := subgraph.IssueMap[dependent.ID]; exists {
			continue // Already in subgraph
		}

		// Check if this dependent has a parent-child relationship with parentID
		depRecs, err := s.GetDependencyRecords(ctx, dependent.ID)
		if err != nil {
			continue
		}

		isChild := false
		for _, depRec := range depRecs {
			if depRec.DependsOnID == parentID && depRec.Type == types.DepParentChild {
				isChild = true
				break
			}
		}

		if !isChild {
			continue
		}

		// Add to subgraph
		subgraph.Issues = append(subgraph.Issues, dependent)
		subgraph.IssueMap[dependent.ID] = dependent
		addedChildren[dependent.ID] = true

		// Recurse to get children of this child
		if err := loadDescendants(ctx, s, subgraph, dependent.ID); err != nil {
			return err
		}
	}

	// Strategy 2: Find hierarchical children by ID pattern
	// This catches children that have missing or incorrect dependency types.
	// Hierarchical IDs follow the pattern: parentID.N (e.g., "gt-abc.1", "gt-abc.2")
	hierarchicalChildren, err := findHierarchicalChildren(ctx, s, parentID)
	if err != nil {
		// Non-fatal: continue with what we have
		return nil
	}

	for _, child := range hierarchicalChildren {
		if addedChildren[child.ID] {
			continue // Already added via dependency
		}
		if _, exists := subgraph.IssueMap[child.ID]; exists {
			continue // Already in subgraph
		}

		// Add to subgraph
		subgraph.Issues = append(subgraph.Issues, child)
		subgraph.IssueMap[child.ID] = child
		addedChildren[child.ID] = true

		// Recurse to get children of this child
		if err := loadDescendants(ctx, s, subgraph, child.ID); err != nil {
			return err
		}
	}

	return nil
}

// findHierarchicalChildren finds issues with IDs that match the pattern parentID.N
// This catches hierarchical children that may be missing parent-child dependencies.
func findHierarchicalChildren(ctx context.Context, s storage.Storage, parentID string) ([]*types.Issue, error) {
	// Look for issues with IDs starting with "parentID."
	// We need to query by ID pattern, which requires listing issues
	pattern := parentID + "."

	// Use the storage's search capability with a filter
	allIssues, err := s.SearchIssues(ctx, "", types.IssueFilter{})
	if err != nil {
		return nil, err
	}

	var children []*types.Issue
	for _, issue := range allIssues {
		// Check if ID starts with pattern and is a direct child (no further dots after the pattern)
		if len(issue.ID) > len(pattern) && issue.ID[:len(pattern)] == pattern {
			// Check it's a direct child, not a grandchild
			// e.g., "parent.1" is a child, "parent.1.2" is a grandchild
			remaining := issue.ID[len(pattern):]
			if !strings.Contains(remaining, ".") {
				children = append(children, issue)
			}
		}
	}

	return children, nil
}

// LoadSubgraphViaDaemon loads a template subgraph using daemon RPC calls
func LoadSubgraphViaDaemon(client *rpc.Client, templateID string) (*Subgraph, error) {
	// Get root issue with dependencies/dependents
	resp, err := client.Show(&rpc.ShowArgs{ID: templateID})
	if err != nil {
		return nil, fmt.Errorf("failed to get template: %w", err)
	}

	var rootDetails IssueDetailsFromShow
	if err := json.Unmarshal(resp.Data, &rootDetails); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	root := &rootDetails.Issue
	subgraph := &Subgraph{
		Root:     root,
		Issues:   []*types.Issue{root},
		IssueMap: map[string]*types.Issue{root.ID: root},
	}

	// Find children from dependents (those with parent-child relationship)
	// and recursively load them
	if err := loadDescendantsViaDaemon(client, subgraph, rootDetails.Dependents); err != nil {
		return nil, err
	}

	// Now build dependencies list by examining each issue's dependencies
	// We need to get the dependency records, which Show provides
	for _, issue := range subgraph.Issues {
		resp, err := client.Show(&rpc.ShowArgs{ID: issue.ID})
		if err != nil {
			continue
		}

		var details IssueDetailsFromShow
		if err := json.Unmarshal(resp.Data, &details); err != nil {
			continue
		}

		// Dependencies are issues that THIS issue depends on
		for _, dep := range details.Dependencies {
			// Only include if the dependency target is also in the subgraph
			if _, ok := subgraph.IssueMap[dep.Issue.ID]; ok {
				subgraph.Dependencies = append(subgraph.Dependencies, &types.Dependency{
					IssueID:     issue.ID,
					DependsOnID: dep.Issue.ID,
					Type:        dep.DependencyType,
				})
			}
		}
	}

	return subgraph, nil
}

// loadDescendantsViaDaemon recursively loads child issues via daemon RPC
func loadDescendantsViaDaemon(client *rpc.Client, subgraph *Subgraph, dependents []*types.IssueWithDependencyMetadata) error {
	for _, dep := range dependents {
		// Check if this is a child (parent-child relationship)
		if dep.DependencyType != types.DepParentChild {
			continue
		}

		if _, exists := subgraph.IssueMap[dep.Issue.ID]; exists {
			continue // Already in subgraph
		}

		// Add to subgraph
		issue := &dep.Issue
		subgraph.Issues = append(subgraph.Issues, issue)
		subgraph.IssueMap[issue.ID] = issue

		// Get this issue's dependents for recursion
		resp, err := client.Show(&rpc.ShowArgs{ID: issue.ID})
		if err != nil {
			continue
		}

		var details IssueDetailsFromShow
		if err := json.Unmarshal(resp.Data, &details); err != nil {
			continue
		}

		// Recurse on children
		if err := loadDescendantsViaDaemon(client, subgraph, details.Dependents); err != nil {
			return err
		}
	}

	return nil
}

// CloneSubgraph creates new issues from the template with variable substitution.
// Uses CloneOptions to control all spawn/bond behavior including dynamic bonding.
func CloneSubgraph(ctx context.Context, s storage.Storage, subgraph *Subgraph, opts CloneOptions) (*InstantiateResult, error) {
	if s == nil {
		return nil, fmt.Errorf("no database connection")
	}

	// Generate new IDs and create mapping
	idMapping := make(map[string]string)

	// Use transaction for atomicity
	err := s.RunInTransaction(ctx, func(tx storage.Transaction) error {
		// First pass: create all issues with new IDs
		for _, oldIssue := range subgraph.Issues {
			// Determine assignee: use override for root epic, otherwise keep template's
			issueAssignee := oldIssue.Assignee
			if oldIssue.ID == subgraph.Root.ID && opts.Assignee != "" {
				issueAssignee = opts.Assignee
			}

			newIssue := &types.Issue{
				// ID will be set below based on bonding options
				Title:              SubstituteVariables(oldIssue.Title, opts.Vars),
				Description:        SubstituteVariables(oldIssue.Description, opts.Vars),
				Design:             SubstituteVariables(oldIssue.Design, opts.Vars),
				AcceptanceCriteria: SubstituteVariables(oldIssue.AcceptanceCriteria, opts.Vars),
				Notes:              SubstituteVariables(oldIssue.Notes, opts.Vars),
				Status:             types.StatusOpen, // Always start fresh
				Priority:           oldIssue.Priority,
				IssueType:          oldIssue.IssueType,
				Assignee:           issueAssignee,
				EstimatedMinutes:   oldIssue.EstimatedMinutes,
				Ephemeral:          opts.Ephemeral, // mark for cleanup when closed
				IDPrefix:           opts.Prefix,   // distinct prefixes for mols/wisps
				// Gate fields (for async coordination)
				AwaitType: oldIssue.AwaitType,
				AwaitID:   SubstituteVariables(oldIssue.AwaitID, opts.Vars),
				Timeout:   oldIssue.Timeout,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			// Generate custom ID for dynamic bonding if ParentID is set
			if opts.ParentID != "" {
				bondedID, err := GenerateBondedID(oldIssue.ID, subgraph.Root.ID, opts)
				if err != nil {
					return fmt.Errorf("failed to generate bonded ID for %s: %w", oldIssue.ID, err)
				}
				newIssue.ID = bondedID
			}

			if err := tx.CreateIssue(ctx, newIssue, opts.Actor); err != nil {
				return fmt.Errorf("failed to create issue from %s: %w", oldIssue.ID, err)
			}

			idMapping[oldIssue.ID] = newIssue.ID
		}

		// Second pass: recreate dependencies with new IDs
		for _, dep := range subgraph.Dependencies {
			newFromID, ok1 := idMapping[dep.IssueID]
			newToID, ok2 := idMapping[dep.DependsOnID]
			if !ok1 || !ok2 {
				continue // Skip if either end is outside the subgraph
			}

			newDep := &types.Dependency{
				IssueID:     newFromID,
				DependsOnID: newToID,
				Type:        dep.Type,
			}
			if err := tx.AddDependency(ctx, newDep, opts.Actor); err != nil {
				return fmt.Errorf("failed to create dependency: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &InstantiateResult{
		NewEpicID: idMapping[subgraph.Root.ID],
		IDMapping: idMapping,
		Created:   len(subgraph.Issues),
	}, nil
}

// CloneSubgraphViaDaemon creates new issues from the template using daemon RPC calls
func CloneSubgraphViaDaemon(client *rpc.Client, subgraph *Subgraph, opts CloneOptions) (*InstantiateResult, error) {
	// Generate new IDs and create mapping
	idMapping := make(map[string]string)

	// First pass: create all issues with new IDs
	for _, oldIssue := range subgraph.Issues {
		// Determine assignee: use override for root epic, otherwise keep template's
		issueAssignee := oldIssue.Assignee
		if oldIssue.ID == subgraph.Root.ID && opts.Assignee != "" {
			issueAssignee = opts.Assignee
		}

		// Build create args
		createArgs := &rpc.CreateArgs{
			Title:              SubstituteVariables(oldIssue.Title, opts.Vars),
			Description:        SubstituteVariables(oldIssue.Description, opts.Vars),
			IssueType:          string(oldIssue.IssueType),
			Priority:           oldIssue.Priority,
			Design:             SubstituteVariables(oldIssue.Design, opts.Vars),
			AcceptanceCriteria: SubstituteVariables(oldIssue.AcceptanceCriteria, opts.Vars),
			Assignee:           issueAssignee,
			EstimatedMinutes:   oldIssue.EstimatedMinutes,
			Ephemeral:          opts.Ephemeral,
			IDPrefix:           opts.Prefix, // distinct prefixes for mols/wisps
		}

		// Generate custom ID for dynamic bonding if ParentID is set
		if opts.ParentID != "" {
			bondedID, err := GenerateBondedID(oldIssue.ID, subgraph.Root.ID, opts)
			if err != nil {
				return nil, fmt.Errorf("failed to generate bonded ID for %s: %w", oldIssue.ID, err)
			}
			createArgs.ID = bondedID
		}

		resp, err := client.Create(createArgs)
		if err != nil {
			return nil, fmt.Errorf("failed to create issue from %s: %w", oldIssue.ID, err)
		}

		// Parse response to get the new issue ID
		var newIssue types.Issue
		if err := json.Unmarshal(resp.Data, &newIssue); err != nil {
			return nil, fmt.Errorf("failed to parse created issue: %w", err)
		}

		idMapping[oldIssue.ID] = newIssue.ID
	}

	// Second pass: recreate dependencies with new IDs
	for _, dep := range subgraph.Dependencies {
		newFromID, ok1 := idMapping[dep.IssueID]
		newToID, ok2 := idMapping[dep.DependsOnID]
		if !ok1 || !ok2 {
			continue // Skip if either end is outside the subgraph
		}

		_, err := client.AddDependency(&rpc.DepAddArgs{
			FromID:  newFromID,
			ToID:    newToID,
			DepType: string(dep.Type),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create dependency: %w", err)
		}
	}

	return &InstantiateResult{
		NewEpicID: idMapping[subgraph.Root.ID],
		IDMapping: idMapping,
		Created:   len(subgraph.Issues),
	}, nil
}

// ResolveProtoIDOrTitle resolves a proto by ID or title.
// It first tries to resolve as an ID (via ResolvePartialID).
// If that fails, it searches for protos with matching titles.
// Returns the proto ID if found, or an error if not found or ambiguous.
func ResolveProtoIDOrTitle(ctx context.Context, s storage.Storage, input string) (string, error) {
	// Strategy 1: Try to resolve as an ID
	protoID, err := utils.ResolvePartialID(ctx, s, input)
	if err == nil {
		// Verify it's a proto (has template label)
		issue, getErr := s.GetIssue(ctx, protoID)
		if getErr == nil && issue != nil {
			labels, _ := s.GetLabels(ctx, protoID)
			for _, label := range labels {
				if label == BeadsTemplateLabel {
					return protoID, nil // Found a valid proto by ID
				}
			}
		}
		// ID resolved but not a proto - continue to title search
	}

	// Strategy 2: Search for protos by title
	protos, err := s.GetIssuesByLabel(ctx, BeadsTemplateLabel)
	if err != nil {
		return "", fmt.Errorf("failed to search protos: %w", err)
	}

	var matches []*types.Issue
	var exactMatch *types.Issue

	for _, proto := range protos {
		// Check for exact title match (case-insensitive)
		if strings.EqualFold(proto.Title, input) {
			exactMatch = proto
			break
		}
		// Check for partial title match (case-insensitive)
		if strings.Contains(strings.ToLower(proto.Title), strings.ToLower(input)) {
			matches = append(matches, proto)
		}
	}

	if exactMatch != nil {
		return exactMatch.ID, nil
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("no proto found matching %q (by ID or title)", input)
	}

	if len(matches) == 1 {
		return matches[0].ID, nil
	}

	// Multiple matches - show them all for disambiguation
	var matchNames []string
	for _, m := range matches {
		matchNames = append(matchNames, fmt.Sprintf("%s: %s", m.ID, m.Title))
	}
	return "", fmt.Errorf("ambiguous: %q matches %d protos:\n  %s\nUse the ID or a more specific title", input, len(matches), strings.Join(matchNames, "\n  "))
}

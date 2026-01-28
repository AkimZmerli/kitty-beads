// Package issues provides the issue management vertical slice.
package issues

import (
	"time"

	"github.com/steveyegge/beads/internal/types"
)

// CreateArgs represents arguments for creating a new issue.
type CreateArgs struct {
	ID                 string   `json:"id,omitempty"`
	Parent             string   `json:"parent,omitempty"` // Parent ID for hierarchical issues
	Title              string   `json:"title"`
	Description        string   `json:"description,omitempty"`
	IssueType          string   `json:"issue_type"`
	Priority           int      `json:"priority"`
	Design             string   `json:"design,omitempty"`
	AcceptanceCriteria string   `json:"acceptance_criteria,omitempty"`
	Notes              string   `json:"notes,omitempty"`
	Assignee           string   `json:"assignee,omitempty"`
	ExternalRef        string   `json:"external_ref,omitempty"`
	EstimatedMinutes   *int     `json:"estimated_minutes,omitempty"`
	Labels             []string `json:"labels,omitempty"`
	Dependencies       []string `json:"dependencies,omitempty"`
	// Waits-for dependencies
	WaitsFor     string `json:"waits_for,omitempty"`
	WaitsForGate string `json:"waits_for_gate,omitempty"`
	// Messaging fields
	Sender    string `json:"sender,omitempty"`
	Ephemeral bool   `json:"ephemeral,omitempty"`
	RepliesTo string `json:"replies_to,omitempty"`
	// ID generation
	IDPrefix  string `json:"id_prefix,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
	Owner     string `json:"owner,omitempty"`
	// Molecule type
	MolType string `json:"mol_type,omitempty"`
	// Agent identity fields
	RoleType string `json:"role_type,omitempty"`
	Rig      string `json:"rig,omitempty"`
	// Event fields
	EventCategory string `json:"event_category,omitempty"`
	EventActor    string `json:"event_actor,omitempty"`
	EventTarget   string `json:"event_target,omitempty"`
	EventPayload  string `json:"event_payload,omitempty"`
	// Time-based scheduling
	DueAt      string `json:"due_at,omitempty"`
	DeferUntil string `json:"defer_until,omitempty"`
}

// UpdateArgs represents arguments for updating an issue.
type UpdateArgs struct {
	ID                 string   `json:"id"`
	Title              *string  `json:"title,omitempty"`
	Description        *string  `json:"description,omitempty"`
	Status             *string  `json:"status,omitempty"`
	Priority           *int     `json:"priority,omitempty"`
	Design             *string  `json:"design,omitempty"`
	AcceptanceCriteria *string  `json:"acceptance_criteria,omitempty"`
	Notes              *string  `json:"notes,omitempty"`
	Assignee           *string  `json:"assignee,omitempty"`
	ExternalRef        *string  `json:"external_ref,omitempty"`
	EstimatedMinutes   *int     `json:"estimated_minutes,omitempty"`
	IssueType          *string  `json:"issue_type,omitempty"`
	AddLabels          []string `json:"add_labels,omitempty"`
	RemoveLabels       []string `json:"remove_labels,omitempty"`
	SetLabels          []string `json:"set_labels,omitempty"`
	// Messaging fields
	Sender    *string `json:"sender,omitempty"`
	Ephemeral *bool   `json:"ephemeral,omitempty"`
	RepliesTo *string `json:"replies_to,omitempty"`
	// Graph link fields
	RelatesTo    *string `json:"relates_to,omitempty"`
	DuplicateOf  *string `json:"duplicate_of,omitempty"`
	SupersededBy *string `json:"superseded_by,omitempty"`
	// Pinned field
	Pinned *bool `json:"pinned,omitempty"`
	// Reparenting field
	Parent *string `json:"parent,omitempty"`
	// Agent slot fields
	HookBead *string `json:"hook_bead,omitempty"`
	RoleBead *string `json:"role_bead,omitempty"`
	// Agent state fields
	AgentState   *string `json:"agent_state,omitempty"`
	LastActivity *bool   `json:"last_activity,omitempty"`
	// Agent identity fields
	RoleType *string `json:"role_type,omitempty"`
	Rig      *string `json:"rig,omitempty"`
	// Event fields
	EventCategory *string `json:"event_category,omitempty"`
	EventActor    *string `json:"event_actor,omitempty"`
	EventTarget   *string `json:"event_target,omitempty"`
	EventPayload  *string `json:"event_payload,omitempty"`
	// Work queue claim operation
	Claim bool `json:"claim,omitempty"`
	// Time-based scheduling
	DueAt      *string `json:"due_at,omitempty"`
	DeferUntil *string `json:"defer_until,omitempty"`
	// Gate fields
	AwaitID *string  `json:"await_id,omitempty"`
	Waiters []string `json:"waiters,omitempty"`
	// Slot fields
	Holder *string `json:"holder,omitempty"`
}

// CloseArgs represents arguments for closing an issue.
type CloseArgs struct {
	ID          string `json:"id"`
	Reason      string `json:"reason,omitempty"`
	Session     string `json:"session,omitempty"`
	SuggestNext bool   `json:"suggest_next,omitempty"`
	Force       bool   `json:"force,omitempty"`
}

// CloseResult is returned when closing an issue with SuggestNext=true.
type CloseResult struct {
	Closed    *types.Issue   `json:"closed"`
	Unblocked []*types.Issue `json:"unblocked,omitempty"`
}

// DeleteArgs represents arguments for deleting issues.
type DeleteArgs struct {
	IDs     []string `json:"ids"`
	Force   bool     `json:"force,omitempty"`
	DryRun  bool     `json:"dry_run,omitempty"`
	Cascade bool     `json:"cascade,omitempty"`
	Reason  string   `json:"reason,omitempty"`
}

// DeleteResult represents the result of a delete operation.
type DeleteResult struct {
	DeletedCount        int      `json:"deleted_count"`
	TotalCount          int      `json:"total_count"`
	DryRun              bool     `json:"dry_run,omitempty"`
	IssueCount          int      `json:"issue_count,omitempty"`
	DependenciesRemoved int      `json:"dependencies_removed,omitempty"`
	LabelsRemoved       int      `json:"labels_removed,omitempty"`
	EventsRemoved       int      `json:"events_removed,omitempty"`
	OrphanedIssues      []string `json:"orphaned_issues,omitempty"`
	Errors              []string `json:"errors,omitempty"`
	PartialSuccess      bool     `json:"partial_success,omitempty"`
}

// ListFilter represents filter options for listing issues.
type ListFilter struct {
	Query     string   `json:"query,omitempty"`
	Status    string   `json:"status,omitempty"`
	Priority  *int     `json:"priority,omitempty"`
	IssueType string   `json:"issue_type,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	Label     string   `json:"label,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	LabelsAny []string `json:"labels_any,omitempty"`
	IDs       []string `json:"ids,omitempty"`
	Limit     int      `json:"limit,omitempty"`
	// Pattern matching
	TitleContains       string `json:"title_contains,omitempty"`
	DescriptionContains string `json:"description_contains,omitempty"`
	NotesContains       string `json:"notes_contains,omitempty"`
	// Date ranges
	CreatedAfter  string `json:"created_after,omitempty"`
	CreatedBefore string `json:"created_before,omitempty"`
	UpdatedAfter  string `json:"updated_after,omitempty"`
	UpdatedBefore string `json:"updated_before,omitempty"`
	ClosedAfter   string `json:"closed_after,omitempty"`
	ClosedBefore  string `json:"closed_before,omitempty"`
	// Empty/null checks
	EmptyDescription bool `json:"empty_description,omitempty"`
	NoAssignee       bool `json:"no_assignee,omitempty"`
	NoLabels         bool `json:"no_labels,omitempty"`
	// Priority range
	PriorityMin *int `json:"priority_min,omitempty"`
	PriorityMax *int `json:"priority_max,omitempty"`
	// Pinned filtering
	Pinned *bool `json:"pinned,omitempty"`
	// Template filtering
	IncludeTemplates bool `json:"include_templates,omitempty"`
	// Parent filtering
	ParentID string `json:"parent_id,omitempty"`
	// Ephemeral filtering
	Ephemeral *bool `json:"ephemeral,omitempty"`
	// Molecule type filtering
	MolType string `json:"mol_type,omitempty"`
	// Status exclusion
	ExcludeStatus []string `json:"exclude_status,omitempty"`
	// Type exclusion
	ExcludeTypes []string `json:"exclude_types,omitempty"`
	// Time-based scheduling filters
	Deferred    bool   `json:"deferred,omitempty"`
	DeferAfter  string `json:"defer_after,omitempty"`
	DeferBefore string `json:"defer_before,omitempty"`
	DueAfter    string `json:"due_after,omitempty"`
	DueBefore   string `json:"due_before,omitempty"`
	Overdue     bool   `json:"overdue,omitempty"`
	// Staleness control
	AllowStale bool `json:"allow_stale,omitempty"`
}

// CountFilter represents filter options for counting issues.
type CountFilter struct {
	Query     string   `json:"query,omitempty"`
	Status    string   `json:"status,omitempty"`
	Priority  *int     `json:"priority,omitempty"`
	IssueType string   `json:"issue_type,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	LabelsAny []string `json:"labels_any,omitempty"`
	IDs       []string `json:"ids,omitempty"`
	// Pattern matching
	TitleContains       string `json:"title_contains,omitempty"`
	DescriptionContains string `json:"description_contains,omitempty"`
	NotesContains       string `json:"notes_contains,omitempty"`
	// Date ranges
	CreatedAfter  string `json:"created_after,omitempty"`
	CreatedBefore string `json:"created_before,omitempty"`
	UpdatedAfter  string `json:"updated_after,omitempty"`
	UpdatedBefore string `json:"updated_before,omitempty"`
	ClosedAfter   string `json:"closed_after,omitempty"`
	ClosedBefore  string `json:"closed_before,omitempty"`
	// Empty/null checks
	EmptyDescription bool `json:"empty_description,omitempty"`
	NoAssignee       bool `json:"no_assignee,omitempty"`
	NoLabels         bool `json:"no_labels,omitempty"`
	// Priority range
	PriorityMin *int `json:"priority_min,omitempty"`
	PriorityMax *int `json:"priority_max,omitempty"`
	// Grouping
	GroupBy string `json:"group_by,omitempty"`
}

// CountResult represents the result of counting issues.
type CountResult struct {
	Count int `json:"count"`
}

// GroupCount represents a count grouped by a field.
type GroupCount struct {
	Group string `json:"group"`
	Count int    `json:"count"`
}

// GroupedCountResult represents grouped count results.
type GroupedCountResult struct {
	Total  int          `json:"total"`
	Groups []GroupCount `json:"groups"`
}

// ShowArgs represents arguments for showing an issue.
type ShowArgs struct {
	ID string `json:"id"`
}

// ReadyFilter represents filter options for ready work.
type ReadyFilter struct {
	Assignee        string   `json:"assignee,omitempty"`
	Unassigned      bool     `json:"unassigned,omitempty"`
	Priority        *int     `json:"priority,omitempty"`
	Type            string   `json:"type,omitempty"`
	Limit           int      `json:"limit,omitempty"`
	SortPolicy      string   `json:"sort_policy,omitempty"`
	Labels          []string `json:"labels,omitempty"`
	LabelsAny       []string `json:"labels_any,omitempty"`
	ParentID        string   `json:"parent_id,omitempty"`
	MolType         string   `json:"mol_type,omitempty"`
	IncludeDeferred bool     `json:"include_deferred,omitempty"`
}

// BlockedFilter represents filter options for blocked issues.
type BlockedFilter struct {
	ParentID string `json:"parent_id,omitempty"`
}

// StaleFilter represents filter options for stale issues.
type StaleFilter struct {
	Days   int    `json:"days,omitempty"`
	Status string `json:"status,omitempty"`
	Limit  int    `json:"limit,omitempty"`
}

// EpicStatusFilter represents filter options for epic status.
type EpicStatusFilter struct {
	EligibleOnly bool `json:"eligible_only,omitempty"`
}

// MutationEvent represents a mutation event emitted by the service.
type MutationEvent struct {
	Type      MutationType `json:"type"`
	IssueID   string       `json:"issue_id"`
	Title     string       `json:"title,omitempty"`
	Assignee  string       `json:"assignee,omitempty"`
	Actor     string       `json:"actor,omitempty"`
	OldStatus string       `json:"old_status,omitempty"`
	NewStatus string       `json:"new_status,omitempty"`
	Timestamp time.Time    `json:"timestamp"`
}

// MutationType represents the type of mutation.
type MutationType string

const (
	MutationCreate MutationType = "create"
	MutationUpdate MutationType = "update"
	MutationDelete MutationType = "delete"
	MutationStatus MutationType = "status"
)

// ResolveIDArgs represents arguments for resolving a partial ID.
type ResolveIDArgs struct {
	ID string `json:"id"`
}

// Package storage provides the storage adapter for vertical slice architecture.
//
// The StorageAdapter wraps the existing storage.Storage interface and provides
// access to feature-specific repository interfaces. This enables gradual migration:
//   - Existing code continues using storage.Storage directly
//   - New vertical slice code uses focused repository interfaces
//   - Both access the same underlying database
package storage

import (
	"context"
	"time"

	"github.com/steveyegge/beads/features/comments"
	"github.com/steveyegge/beads/features/dependencies"
	"github.com/steveyegge/beads/features/export"
	"github.com/steveyegge/beads/features/issues"
	"github.com/steveyegge/beads/features/kanban"
	"github.com/steveyegge/beads/features/labels"
	"github.com/steveyegge/beads/features/statistics"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
)

// Adapter provides access to feature-specific repositories.
// It wraps the underlying storage.Storage and delegates to it.
type Adapter struct {
	store storage.Storage
}

// NewAdapter creates a new storage adapter.
func NewAdapter(store storage.Storage) *Adapter {
	return &Adapter{store: store}
}

// Storage returns the underlying storage.Storage for backward compatibility.
func (a *Adapter) Storage() storage.Storage {
	return a.store
}

// Issues returns the issues repository.
func (a *Adapter) Issues() issues.Repository {
	return &issueRepo{a.store}
}

// Labels returns the labels repository.
func (a *Adapter) Labels() labels.Repository {
	return &labelRepo{a.store}
}

// Comments returns the comments repository.
func (a *Adapter) Comments() comments.Repository {
	return &commentRepo{a.store}
}

// Events returns the events repository.
func (a *Adapter) Events() comments.EventRepository {
	return &eventRepo{a.store}
}

// Dependencies returns the dependencies repository.
func (a *Adapter) Dependencies() dependencies.Repository {
	return &dependencyRepo{a.store}
}

// Statistics returns the statistics repository.
func (a *Adapter) Statistics() statistics.Repository {
	return &statisticsRepo{a.store}
}

// Kanban returns the kanban/ready work repository.
func (a *Adapter) Kanban() kanban.Repository {
	return &kanbanRepo{a.store}
}

// Export returns the export repository.
func (a *Adapter) Export() export.Repository {
	return &exportRepo{a.store}
}

// ===== Issue Repository Adapter =====

type issueRepo struct {
	store storage.Storage
}

func (r *issueRepo) Create(ctx context.Context, issue *types.Issue, actor string) error {
	return r.store.CreateIssue(ctx, issue, actor)
}

func (r *issueRepo) CreateBatch(ctx context.Context, issues []*types.Issue, actor string) error {
	return r.store.CreateIssues(ctx, issues, actor)
}

func (r *issueRepo) Get(ctx context.Context, id string) (*types.Issue, error) {
	return r.store.GetIssue(ctx, id)
}

func (r *issueRepo) GetByExternalRef(ctx context.Context, externalRef string) (*types.Issue, error) {
	return r.store.GetIssueByExternalRef(ctx, externalRef)
}

func (r *issueRepo) Update(ctx context.Context, id string, updates map[string]interface{}, actor string) error {
	return r.store.UpdateIssue(ctx, id, updates, actor)
}

func (r *issueRepo) Close(ctx context.Context, id string, reason string, actor string, session string) error {
	return r.store.CloseIssue(ctx, id, reason, actor, session)
}

func (r *issueRepo) Delete(ctx context.Context, id string) error {
	return r.store.DeleteIssue(ctx, id)
}

func (r *issueRepo) Search(ctx context.Context, query string, filter types.IssueFilter) ([]*types.Issue, error) {
	return r.store.SearchIssues(ctx, query, filter)
}

func (r *issueRepo) GetNextChildID(ctx context.Context, parentID string) (string, error) {
	return r.store.GetNextChildID(ctx, parentID)
}

func (r *issueRepo) UpdateID(ctx context.Context, oldID, newID string, issue *types.Issue, actor string) error {
	return r.store.UpdateIssueID(ctx, oldID, newID, issue, actor)
}

// ===== Label Repository Adapter =====

type labelRepo struct {
	store storage.Storage
}

func (r *labelRepo) Add(ctx context.Context, issueID, label, actor string) error {
	return r.store.AddLabel(ctx, issueID, label, actor)
}

func (r *labelRepo) Remove(ctx context.Context, issueID, label, actor string) error {
	return r.store.RemoveLabel(ctx, issueID, label, actor)
}

func (r *labelRepo) Get(ctx context.Context, issueID string) ([]string, error) {
	return r.store.GetLabels(ctx, issueID)
}

func (r *labelRepo) GetForIssues(ctx context.Context, issueIDs []string) (map[string][]string, error) {
	return r.store.GetLabelsForIssues(ctx, issueIDs)
}

func (r *labelRepo) GetIssuesByLabel(ctx context.Context, label string) ([]*types.Issue, error) {
	return r.store.GetIssuesByLabel(ctx, label)
}

// ===== Comment Repository Adapter =====

type commentRepo struct {
	store storage.Storage
}

func (r *commentRepo) Add(ctx context.Context, issueID, author, text string) (*types.Comment, error) {
	return r.store.AddIssueComment(ctx, issueID, author, text)
}

func (r *commentRepo) Import(ctx context.Context, issueID, author, text string, createdAt time.Time) (*types.Comment, error) {
	return r.store.ImportIssueComment(ctx, issueID, author, text, createdAt)
}

func (r *commentRepo) Get(ctx context.Context, issueID string) ([]*types.Comment, error) {
	return r.store.GetIssueComments(ctx, issueID)
}

func (r *commentRepo) GetForIssues(ctx context.Context, issueIDs []string) (map[string][]*types.Comment, error) {
	return r.store.GetCommentsForIssues(ctx, issueIDs)
}

// ===== Event Repository Adapter =====

type eventRepo struct {
	store storage.Storage
}

func (r *eventRepo) AddComment(ctx context.Context, issueID, actor, comment string) error {
	return r.store.AddComment(ctx, issueID, actor, comment)
}

func (r *eventRepo) Get(ctx context.Context, issueID string, limit int) ([]*types.Event, error) {
	return r.store.GetEvents(ctx, issueID, limit)
}

// ===== Dependency Repository Adapter =====

type dependencyRepo struct {
	store storage.Storage
}

func (r *dependencyRepo) Add(ctx context.Context, dep *types.Dependency, actor string) error {
	return r.store.AddDependency(ctx, dep, actor)
}

func (r *dependencyRepo) Remove(ctx context.Context, issueID, dependsOnID string, actor string) error {
	return r.store.RemoveDependency(ctx, issueID, dependsOnID, actor)
}

func (r *dependencyRepo) GetDependencies(ctx context.Context, issueID string) ([]*types.Issue, error) {
	return r.store.GetDependencies(ctx, issueID)
}

func (r *dependencyRepo) GetDependents(ctx context.Context, issueID string) ([]*types.Issue, error) {
	return r.store.GetDependents(ctx, issueID)
}

func (r *dependencyRepo) GetDependenciesWithMetadata(ctx context.Context, issueID string) ([]*types.IssueWithDependencyMetadata, error) {
	return r.store.GetDependenciesWithMetadata(ctx, issueID)
}

func (r *dependencyRepo) GetDependentsWithMetadata(ctx context.Context, issueID string) ([]*types.IssueWithDependencyMetadata, error) {
	return r.store.GetDependentsWithMetadata(ctx, issueID)
}

func (r *dependencyRepo) GetRecords(ctx context.Context, issueID string) ([]*types.Dependency, error) {
	return r.store.GetDependencyRecords(ctx, issueID)
}

func (r *dependencyRepo) GetAllRecords(ctx context.Context) (map[string][]*types.Dependency, error) {
	return r.store.GetAllDependencyRecords(ctx)
}

func (r *dependencyRepo) GetCounts(ctx context.Context, issueIDs []string) (map[string]*types.DependencyCounts, error) {
	return r.store.GetDependencyCounts(ctx, issueIDs)
}

func (r *dependencyRepo) GetTree(ctx context.Context, issueID string, maxDepth int, showAllPaths bool, reverse bool) ([]*types.TreeNode, error) {
	return r.store.GetDependencyTree(ctx, issueID, maxDepth, showAllPaths, reverse)
}

func (r *dependencyRepo) DetectCycles(ctx context.Context) ([][]*types.Issue, error) {
	return r.store.DetectCycles(ctx)
}

func (r *dependencyRepo) RenamePrefix(ctx context.Context, oldPrefix, newPrefix string) error {
	return r.store.RenameDependencyPrefix(ctx, oldPrefix, newPrefix)
}

// ===== Statistics Repository Adapter =====

type statisticsRepo struct {
	store storage.Storage
}

func (r *statisticsRepo) Get(ctx context.Context) (*types.Statistics, error) {
	return r.store.GetStatistics(ctx)
}

func (r *statisticsRepo) GetMoleculeProgress(ctx context.Context, moleculeID string) (*types.MoleculeProgressStats, error) {
	return r.store.GetMoleculeProgress(ctx, moleculeID)
}

// ===== Kanban Repository Adapter =====

type kanbanRepo struct {
	store storage.Storage
}

func (r *kanbanRepo) GetReadyWork(ctx context.Context, filter types.WorkFilter) ([]*types.Issue, error) {
	return r.store.GetReadyWork(ctx, filter)
}

func (r *kanbanRepo) GetBlockedIssues(ctx context.Context, filter types.WorkFilter) ([]*types.BlockedIssue, error) {
	return r.store.GetBlockedIssues(ctx, filter)
}

func (r *kanbanRepo) IsBlocked(ctx context.Context, issueID string) (bool, []string, error) {
	return r.store.IsBlocked(ctx, issueID)
}

func (r *kanbanRepo) GetEpicsEligibleForClosure(ctx context.Context) ([]*types.EpicStatus, error) {
	return r.store.GetEpicsEligibleForClosure(ctx)
}

func (r *kanbanRepo) GetStaleIssues(ctx context.Context, filter types.StaleFilter) ([]*types.Issue, error) {
	return r.store.GetStaleIssues(ctx, filter)
}

func (r *kanbanRepo) GetNewlyUnblockedByClose(ctx context.Context, closedIssueID string) ([]*types.Issue, error) {
	return r.store.GetNewlyUnblockedByClose(ctx, closedIssueID)
}

// ===== Export Repository Adapter =====

type exportRepo struct {
	store storage.Storage
}

func (r *exportRepo) GetDirtyIssues(ctx context.Context) ([]string, error) {
	return r.store.GetDirtyIssues(ctx)
}

func (r *exportRepo) GetDirtyIssueHash(ctx context.Context, issueID string) (string, error) {
	return r.store.GetDirtyIssueHash(ctx, issueID)
}

func (r *exportRepo) ClearDirtyIssuesByID(ctx context.Context, issueIDs []string) error {
	return r.store.ClearDirtyIssuesByID(ctx, issueIDs)
}

func (r *exportRepo) GetExportHash(ctx context.Context, issueID string) (string, error) {
	return r.store.GetExportHash(ctx, issueID)
}

func (r *exportRepo) SetExportHash(ctx context.Context, issueID, contentHash string) error {
	return r.store.SetExportHash(ctx, issueID, contentHash)
}

func (r *exportRepo) ClearAllExportHashes(ctx context.Context) error {
	return r.store.ClearAllExportHashes(ctx)
}

func (r *exportRepo) GetJSONLFileHash(ctx context.Context) (string, error) {
	return r.store.GetJSONLFileHash(ctx)
}

func (r *exportRepo) SetJSONLFileHash(ctx context.Context, fileHash string) error {
	return r.store.SetJSONLFileHash(ctx, fileHash)
}

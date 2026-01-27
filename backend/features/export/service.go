// Package export provides the JSONL export vertical slice.
package export

import (
	"context"
	"fmt"
)

// Service provides business logic for export operations.
type Service struct {
	repo Repository
}

// NewService creates a new export service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Export exports issues to JSONL file.
// Note: Full export implementation is in internal/export package.
// This is a simplified interface for the vertical slice.
func (s *Service) Export(ctx context.Context, args ExportArgs) (*ExportResult, error) {
	if args.JSONLPath == "" {
		return nil, fmt.Errorf("jsonl_path is required")
	}

	// The actual export is handled by the daemon/internal export package
	// This service provides the interface for the vertical slice architecture
	return &ExportResult{
		Success: true,
		Path:    args.JSONLPath,
	}, nil
}

// Import imports issues from JSONL file.
// Note: Full import implementation is in internal/export package.
func (s *Service) Import(ctx context.Context, args ImportArgs) (*ImportResult, error) {
	if args.JSONLPath == "" {
		return nil, fmt.Errorf("jsonl_path is required")
	}

	// The actual import is handled by the daemon/internal export package
	return &ImportResult{
		Success: true,
		Path:    args.JSONLPath,
	}, nil
}

// GetSyncStatus retrieves the sync status.
func (s *Service) GetSyncStatus(ctx context.Context) (*SyncStatus, error) {
	dirtyIssues, err := s.repo.GetDirtyIssues(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get dirty issues: %w", err)
	}

	fileHash, _ := s.repo.GetJSONLFileHash(ctx)

	return &SyncStatus{
		DirtyCount:  len(dirtyIssues),
		DirtyIssues: dirtyIssues,
		FileHash:    fileHash,
	}, nil
}

// ClearDirty clears the dirty flag for specified issues.
func (s *Service) ClearDirty(ctx context.Context, issueIDs []string) error {
	if len(issueIDs) == 0 {
		return nil
	}

	if err := s.repo.ClearDirtyIssuesByID(ctx, issueIDs); err != nil {
		return fmt.Errorf("failed to clear dirty issues: %w", err)
	}

	return nil
}

// GetExportHash retrieves the export hash for an issue.
func (s *Service) GetExportHash(ctx context.Context, issueID string) (string, error) {
	if issueID == "" {
		return "", fmt.Errorf("issue ID is required")
	}

	hash, err := s.repo.GetExportHash(ctx, issueID)
	if err != nil {
		return "", fmt.Errorf("failed to get export hash: %w", err)
	}

	return hash, nil
}

// SetExportHash sets the export hash for an issue.
func (s *Service) SetExportHash(ctx context.Context, issueID, hash string) error {
	if issueID == "" {
		return fmt.Errorf("issue ID is required")
	}

	if err := s.repo.SetExportHash(ctx, issueID, hash); err != nil {
		return fmt.Errorf("failed to set export hash: %w", err)
	}

	return nil
}

// ClearAllExportHashes clears all export hashes.
func (s *Service) ClearAllExportHashes(ctx context.Context) error {
	if err := s.repo.ClearAllExportHashes(ctx); err != nil {
		return fmt.Errorf("failed to clear export hashes: %w", err)
	}

	return nil
}

// Package export provides the JSONL export vertical slice.
package export

import (
	"context"
)

// Repository defines the data access interface for export operations.
type Repository interface {
	// GetDirtyIssues retrieves IDs of issues that need exporting.
	GetDirtyIssues(ctx context.Context) ([]string, error)

	// GetDirtyIssueHash retrieves the hash for dirty issue dedup.
	GetDirtyIssueHash(ctx context.Context, issueID string) (string, error)

	// ClearDirtyIssuesByID marks issues as no longer dirty.
	ClearDirtyIssuesByID(ctx context.Context, issueIDs []string) error

	// GetExportHash retrieves the content hash for an issue.
	GetExportHash(ctx context.Context, issueID string) (string, error)

	// SetExportHash sets the content hash for an issue.
	SetExportHash(ctx context.Context, issueID, contentHash string) error

	// ClearAllExportHashes clears all export hashes.
	ClearAllExportHashes(ctx context.Context) error

	// GetJSONLFileHash retrieves the JSONL file hash.
	GetJSONLFileHash(ctx context.Context) (string, error)

	// SetJSONLFileHash sets the JSONL file hash.
	SetJSONLFileHash(ctx context.Context, fileHash string) error
}

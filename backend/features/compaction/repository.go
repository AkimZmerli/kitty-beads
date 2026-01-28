// Package compaction provides the issue compaction vertical slice.
package compaction

// Repository defines the data access interface for compaction.
// Compaction reduces issue size by archiving old history.
type Repository interface {
	// Compaction-specific operations will be added here as they're extracted
	// from the RPC handlers in Phase 3-4.
}

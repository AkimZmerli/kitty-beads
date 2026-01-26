// Package epics provides the epic management vertical slice.
package epics

// Repository defines the data access interface for epics.
// Epic-specific operations (beyond issue CRUD) will be extracted here.
type Repository interface {
	// Epic-specific operations will be added here as they're extracted
	// from the RPC handlers in Phase 3-4.
}

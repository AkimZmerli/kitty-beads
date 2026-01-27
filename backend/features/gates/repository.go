// Package gates provides the async coordination (gates) vertical slice.
package gates

// Repository defines the data access interface for gates.
// Gates are async coordination primitives (await conditions).
// The actual gate logic is implemented via issue fields (AwaitType, AwaitID, Timeout).
// This package will contain gate-specific queries and operations.
type Repository interface {
	// Gate-specific operations will be added here as they're extracted
	// from the RPC handlers in Phase 3-4.
}

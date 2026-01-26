// Package gates provides the async coordination (gates) vertical slice.
package gates

import (
	"time"

	"github.com/steveyegge/beads/internal/types"
)

// CreateArgs represents arguments for creating a gate.
type CreateArgs struct {
	Title     string        `json:"title"`
	AwaitType string        `json:"await_type"` // gh:run, gh:pr, timer, human, mail
	AwaitID   string        `json:"await_id"`   // ID/value for the await type
	Timeout   time.Duration `json:"timeout"`    // Timeout duration
	Waiters   []string      `json:"waiters"`    // Mail addresses to notify when gate clears
}

// CloseArgs represents arguments for closing a gate.
type CloseArgs struct {
	ID     string `json:"id"`
	Reason string `json:"reason,omitempty"`
}

// WaitArgs represents arguments for adding waiters to a gate.
type WaitArgs struct {
	ID      string   `json:"id"`
	Waiters []string `json:"waiters"`
}

// ListArgs represents arguments for listing gates.
type ListArgs struct {
	All bool `json:"all"` // Include closed gates
}

// ShowArgs represents arguments for showing a gate.
type ShowArgs struct {
	ID string `json:"id"` // Gate ID (partial or full)
}

// CreateResult represents the result of creating a gate.
type CreateResult struct {
	ID string `json:"id"`
}

// WaitResult represents the result of adding waiters.
type WaitResult struct {
	AddedCount int `json:"added_count"`
}

// GateInfo represents gate information for display.
type GateInfo struct {
	Issue     *types.Issue `json:"issue"`
	AwaitType string       `json:"await_type"`
	AwaitID   string       `json:"await_id"`
	Waiters   []string     `json:"waiters"`
	Status    string       `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	ClosedAt  *time.Time   `json:"closed_at,omitempty"`
}

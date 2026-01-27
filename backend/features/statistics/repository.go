// Package statistics provides the statistics vertical slice.
package statistics

import (
	"context"

	"github.com/steveyegge/beads/internal/types"
)

// Repository defines the data access interface for statistics.
type Repository interface {
	// Get retrieves aggregate statistics.
	Get(ctx context.Context) (*types.Statistics, error)

	// GetMoleculeProgress retrieves progress stats for a molecule.
	GetMoleculeProgress(ctx context.Context, moleculeID string) (*types.MoleculeProgressStats, error)
}

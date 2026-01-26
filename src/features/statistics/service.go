// Package statistics provides the statistics vertical slice.
package statistics

import (
	"context"
	"fmt"

	"github.com/steveyegge/beads/internal/types"
)

// Service provides business logic for statistics operations.
type Service struct {
	repo Repository
}

// NewService creates a new statistics service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// Get retrieves aggregate statistics.
func (s *Service) Get(ctx context.Context) (*types.Statistics, error) {
	stats, err := s.repo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get statistics: %w", err)
	}

	return stats, nil
}

// GetMoleculeProgress retrieves progress stats for a molecule.
func (s *Service) GetMoleculeProgress(ctx context.Context, moleculeID string) (*types.MoleculeProgressStats, error) {
	if moleculeID == "" {
		return nil, fmt.Errorf("molecule ID is required")
	}

	progress, err := s.repo.GetMoleculeProgress(ctx, moleculeID)
	if err != nil {
		return nil, fmt.Errorf("failed to get molecule progress: %w", err)
	}

	return progress, nil
}

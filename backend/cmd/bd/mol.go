package main

import (
	"context"

	molcmd "github.com/steveyegge/beads/cmd/bd/commands/molecules"
	"github.com/steveyegge/beads/internal/storage"
)

// Molecule commands - work templates for agent workflows
//
// This file provides backward-compatible exports for the mol commands.
// The actual command is now in the molecules package (commands/molecules/).
// Subcommand files (mol_bond.go, pour.go, etc.) still use molCmd.AddCommand().

// MoleculeLabel is the label used to identify molecules (templates)
// Molecules use the same label as templates - they ARE templates with workflow semantics
const MoleculeLabel = molcmd.MoleculeLabel

// MoleculeSubgraph is an alias for TemplateSubgraph
// Molecules and templates share the same subgraph structure
type MoleculeSubgraph = molcmd.MoleculeSubgraph

// =============================================================================
// Molecule Helper Functions
// =============================================================================

// spawnMolecule creates new issues from the proto with variable substitution.
// This instantiates a proto (template) into a molecule (real issues).
// If ephemeral is true, spawned issues are marked for bulk deletion when closed.
// The prefix parameter overrides the default issue prefix.
func spawnMolecule(ctx context.Context, s storage.Storage, subgraph *MoleculeSubgraph, vars map[string]string, assignee string, actorName string, ephemeral bool, prefix string) (*InstantiateResult, error) {
	return molcmd.SpawnMolecule(ctx, s, subgraph, vars, assignee, actorName, ephemeral, prefix)
}

// Note: init() removed - molCmd is now registered by molcmd.Register() in main.go

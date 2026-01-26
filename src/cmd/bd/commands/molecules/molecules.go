// Package molecules implements the bd CLI molecule commands (work templates).
//
// Terminology:
//   - Proto: Uninstantiated template (easter egg: 'protomolecule' alias)
//   - Molecule: A spawned instance of a proto
//   - Spawn: Instantiate a proto, creating real issues from the template
//   - Bond: Polymorphic combine operation (proto+proto, proto+mol, mol+mol)
//   - Distill: Extract ad-hoc epic → reusable proto
//   - Compound: Result of bonding
package molecules

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/cmd/bd/cli"
	"github.com/steveyegge/beads/cmd/bd/commands/shared/template"
	"github.com/steveyegge/beads/internal/storage"
)

// MoleculeLabel is the label used to identify molecules (templates)
// Molecules use the same label as templates - they ARE templates with workflow semantics
const MoleculeLabel = template.BeadsTemplateLabel

// MoleculeSubgraph is an alias for template.Subgraph
// Molecules and templates share the same subgraph structure
type MoleculeSubgraph = template.Subgraph

// molCmd is the parent command for molecule operations
var molCmd = &cobra.Command{
	Use:     "mol",
	Aliases: []string{"protomolecule"}, // Easter egg for The Expanse fans
	Short:   "Molecule commands (work templates)",
	Long: `Manage molecules - work templates for agent workflows.

Protos are template epics with the "template" label. They define a DAG of work
that can be spawned to create real issues (molecules).

The molecule metaphor:
  - A proto is an uninstantiated template (reusable work pattern)
  - Spawning creates a molecule (real issues) from the proto
  - Variables ({{key}}) are substituted during spawning
  - Bonding combines protos or molecules into compounds
  - Distilling extracts a proto from an ad-hoc epic

Commands:
  show       Show proto/molecule structure and variables
  pour       Instantiate proto as persistent mol (liquid phase)
  wisp       Instantiate proto as ephemeral wisp (vapor phase)
  bond       Polymorphic combine: proto+proto, proto+mol, mol+mol
  squash     Condense molecule to digest
  burn       Discard wisp
  distill    Extract proto from ad-hoc epic

Use "bd formula list" to list available formulas.`,
}

// =============================================================================
// Molecule Helper Functions
// =============================================================================

// SpawnMolecule creates new issues from the proto with variable substitution.
// This instantiates a proto (template) into a molecule (real issues).
// If ephemeral is true, spawned issues are marked for bulk deletion when closed.
// The prefix parameter overrides the default issue prefix.
func SpawnMolecule(ctx context.Context, s storage.Storage, subgraph *MoleculeSubgraph, vars map[string]string, assignee string, actorName string, ephemeral bool, prefix string) (*template.InstantiateResult, error) {
	opts := template.CloneOptions{
		Vars:      vars,
		Assignee:  assignee,
		Actor:     actorName,
		Ephemeral: ephemeral,
		Prefix:    prefix,
	}
	return template.CloneSubgraph(ctx, s, subgraph, opts)
}

// SpawnMoleculeWithOptions creates new issues from the proto using CloneOptions.
// This allows full control over dynamic bonding, variable substitution, and wisp phase.
func SpawnMoleculeWithOptions(ctx context.Context, s storage.Storage, subgraph *MoleculeSubgraph, opts template.CloneOptions) (*template.InstantiateResult, error) {
	return template.CloneSubgraph(ctx, s, subgraph, opts)
}

// PrintMoleculeTree prints the molecule structure as a tree
func PrintMoleculeTree(subgraph *MoleculeSubgraph, parentID string, depth int, isRoot bool) {
	cliCtx := cli.Get()
	if cliCtx.IsJSONOutput() {
		return // Don't print tree in JSON mode
	}
	template.PrintTree(subgraph, parentID, depth, isRoot, func(format string, args ...interface{}) {
		// Use fmt.Printf for output (imported by callers)
		// This function is typically called in contexts where fmt is available
	})
}

// IsProto checks if an issue has the template label
func IsProto(labels []string) bool {
	for _, label := range labels {
		if label == MoleculeLabel {
			return true
		}
	}
	return false
}

// Register adds molecule commands to the root command.
// Called from main package during initialization.
func Register(root *cobra.Command) {
	// Register migrated subcommands
	registerShowCmd()
	registerCookCmd()

	// Note: Other subcommands (pour, wisp, bond, etc.) are still in main package
	// and register themselves via their init() functions using GetMolCmd()

	root.AddCommand(molCmd)

	// cook is a top-level command (not under mol)
	root.AddCommand(cookCmd)
}

// GetMolCmd returns the mol command for subcommand registration.
// Used by main package to register subcommands that are still in main.
func GetMolCmd() *cobra.Command {
	return molCmd
}

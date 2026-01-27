// Package admin implements the bd CLI admin command group.
// The actual admin subcommands (cleanup, compact, reset) remain in the main
// package due to their complex dependencies on global state and helper functions.
package admin

import (
	"github.com/spf13/cobra"
)

var adminCmd = &cobra.Command{
	Use:     "admin",
	GroupID: "advanced",
	Short:   "Administrative commands for database maintenance",
	Long: `Administrative commands for beads database maintenance.

These commands are for advanced users and should be used carefully:
  cleanup   Delete closed issues and prune expired tombstones
  compact   Compact old closed issues to save space
  reset     Remove all beads data and configuration

For routine operations, prefer 'bd doctor --fix'.`,
}

// GetCommand returns the admin parent command.
// Main package uses this to add subcommands.
func GetCommand() *cobra.Command {
	return adminCmd
}

// AddSubcommand adds a subcommand to the admin command.
// This allows main package to register cleanup, compact, reset.
func AddSubcommand(cmd *cobra.Command) {
	adminCmd.AddCommand(cmd)
}

// Register adds the admin command to the root command.
// Called from main package during initialization.
func Register(root *cobra.Command) {
	root.AddCommand(adminCmd)
}

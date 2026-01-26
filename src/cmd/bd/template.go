package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/cmd/bd/commands/shared/template"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/storage"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/ui"
	"github.com/steveyegge/beads/internal/utils"
)

// BeadsTemplateLabel is the label used to identify Beads-based templates
// Deprecated: Use template.BeadsTemplateLabel from commands/shared/template instead
const BeadsTemplateLabel = template.BeadsTemplateLabel

// variablePattern matches {{variable}} placeholders
// Deprecated: Use template.VariablePattern from commands/shared/template instead
var variablePattern = template.VariablePattern

// TemplateSubgraph holds a template epic and all its descendants
// Deprecated: Use template.Subgraph from commands/shared/template instead
type TemplateSubgraph = template.Subgraph

// InstantiateResult holds the result of template instantiation
// Deprecated: Use template.InstantiateResult from commands/shared/template instead
type InstantiateResult = template.InstantiateResult

// CloneOptions controls how the subgraph is cloned during spawn/bond
// Deprecated: Use template.CloneOptions from commands/shared/template instead
type CloneOptions = template.CloneOptions

// IssueDetailsFromShow represents the response structure from daemon Show RPC
// Deprecated: Use template.IssueDetailsFromShow from commands/shared/template instead
type IssueDetailsFromShow = template.IssueDetailsFromShow

// bondedIDPattern validates bonded IDs (alphanumeric, dash, underscore, dot)
// Deprecated: Use template.BondedIDPattern from commands/shared/template instead
var bondedIDPattern = template.BondedIDPattern

var templateCmd = &cobra.Command{
	Use:        "template",
	GroupID:    "setup",
	Short:      "Manage issue templates",
	Deprecated: "use 'bd mol' instead (will be removed in v1.0.0)",
	Long: `Manage Beads templates for creating issue hierarchies.

Templates are epics with the "template" label. They can have child issues
with {{variable}} placeholders that get substituted during instantiation.

To create a template:
  1. Create an epic with child issues
  2. Add the 'template' label: bd label add <epic-id> template
  3. Use {{variable}} placeholders in titles/descriptions

To use a template:
  bd template instantiate <id> --var key=value`,
}

var templateListCmd = &cobra.Command{
	Use:        "list",
	Short:      "List available templates",
	Deprecated: "use 'bd formula list' instead (will be removed in v1.0.0)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := rootCtx
		var beadsTemplates []*types.Issue

		if daemonClient != nil {
			resp, err := daemonClient.List(&rpc.ListArgs{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading templates: %v\n", err)
				os.Exit(1)
			}
			var allIssues []*types.Issue
			if err := json.Unmarshal(resp.Data, &allIssues); err == nil {
				for _, issue := range allIssues {
					for _, label := range issue.Labels {
						if label == BeadsTemplateLabel {
							beadsTemplates = append(beadsTemplates, issue)
							break
						}
					}
				}
			}
		} else if store != nil {
			var err error
			beadsTemplates, err = store.GetIssuesByLabel(ctx, BeadsTemplateLabel)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error loading templates: %v\n", err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: no database connection\n")
			os.Exit(1)
		}

		if jsonOutput {
			outputJSON(beadsTemplates)
			return
		}

		// Human-readable output
		if len(beadsTemplates) == 0 {
			fmt.Println("No templates available.")
			fmt.Println("\nTo create a template:")
			fmt.Println("  1. Create an epic with child issues")
			fmt.Println("  2. Add the 'template' label: bd label add <epic-id> template")
			fmt.Println("  3. Use {{variable}} placeholders in titles/descriptions")
			return
		}

		fmt.Printf("%s\n", ui.RenderPass("Templates (for bd template instantiate):"))
		for _, tmpl := range beadsTemplates {
			vars := extractVariables(tmpl.Title + " " + tmpl.Description)
			varStr := ""
			if len(vars) > 0 {
				varStr = fmt.Sprintf(" (vars: %s)", strings.Join(vars, ", "))
			}
			fmt.Printf("  %s: %s%s\n", ui.RenderAccent(tmpl.ID), tmpl.Title, varStr)
		}
		fmt.Println()
	},
}

var templateShowCmd = &cobra.Command{
	Use:        "show <template-id>",
	Short:      "Show template details",
	Deprecated: "use 'bd mol show' instead (will be removed in v1.0.0)",
	Args:       cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := rootCtx
		var templateID string

		if daemonClient != nil {
			resolveArgs := &rpc.ResolveIDArgs{ID: args[0]}
			resp, err := daemonClient.ResolveID(resolveArgs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: template '%s' not found\n", args[0])
				os.Exit(1)
			}
			if err := json.Unmarshal(resp.Data, &templateID); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else if store != nil {
			var err error
			templateID, err = utils.ResolvePartialID(ctx, store, args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: template '%s' not found\n", args[0])
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: no database connection\n")
			os.Exit(1)
		}

		// Load and show Beads template
		var subgraph *TemplateSubgraph
		var err error
		if daemonClient != nil {
			subgraph, err = loadTemplateSubgraphViaDaemon(daemonClient, templateID)
		} else {
			subgraph, err = loadTemplateSubgraph(ctx, store, templateID)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
			os.Exit(1)
		}

		showBeadsTemplate(subgraph)
	},
}

func showBeadsTemplate(subgraph *TemplateSubgraph) {
	if jsonOutput {
		outputJSON(map[string]interface{}{
			"root":         subgraph.Root,
			"issues":       subgraph.Issues,
			"dependencies": subgraph.Dependencies,
			"variables":    extractAllVariables(subgraph),
		})
		return
	}

	fmt.Printf("\n%s Template: %s\n", ui.RenderAccent("📋"), subgraph.Root.Title)
	fmt.Printf("   ID: %s\n", subgraph.Root.ID)
	fmt.Printf("   Issues: %d\n", len(subgraph.Issues))

	// Show variables
	vars := extractAllVariables(subgraph)
	if len(vars) > 0 {
		fmt.Printf("\n%s Variables:\n", ui.RenderWarn("📝"))
		for _, v := range vars {
			fmt.Printf("   {{%s}}\n", v)
		}
	}

	// Show structure
	fmt.Printf("\n%s Structure:\n", ui.RenderPass("🌲"))
	printTemplateTree(subgraph, subgraph.Root.ID, 0, true)
	fmt.Println()
}

var templateInstantiateCmd = &cobra.Command{
	Use:        "instantiate <template-id>",
	Short:      "Create issues from a Beads template",
	Deprecated: "use 'bd mol bond' instead (will be removed in v1.0.0)",
	Long: `Instantiate a Beads template by cloning its subgraph and substituting variables.

Variables are specified with --var key=value flags. The template's {{key}}
placeholders will be replaced with the corresponding values.

Example:
  bd template instantiate bd-abc123 --var version=1.2.0 --var date=2024-01-15`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CheckReadonly("template instantiate")

		ctx := rootCtx
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		varFlags, _ := cmd.Flags().GetStringArray("var")
		assignee, _ := cmd.Flags().GetString("assignee")

		// Parse variables
		vars := make(map[string]string)
		for _, v := range varFlags {
			parts := strings.SplitN(v, "=", 2)
			if len(parts) != 2 {
				fmt.Fprintf(os.Stderr, "Error: invalid variable format '%s', expected 'key=value'\n", v)
				os.Exit(1)
			}
			vars[parts[0]] = parts[1]
		}

		// Resolve template ID
		var templateID string
		if daemonClient != nil {
			resolveArgs := &rpc.ResolveIDArgs{ID: args[0]}
			resp, err := daemonClient.ResolveID(resolveArgs)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error resolving template ID %s: %v\n", args[0], err)
				os.Exit(1)
			}
			if err := json.Unmarshal(resp.Data, &templateID); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %v\n", err)
				os.Exit(1)
			}
		} else if store != nil {
			var err error
			templateID, err = utils.ResolvePartialID(ctx, store, args[0])
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error resolving template ID %s: %v\n", args[0], err)
				os.Exit(1)
			}
		} else {
			fmt.Fprintf(os.Stderr, "Error: no database connection\n")
			os.Exit(1)
		}

		// Load the template subgraph
		var subgraph *TemplateSubgraph
		var err error
		if daemonClient != nil {
			subgraph, err = loadTemplateSubgraphViaDaemon(daemonClient, templateID)
		} else {
			subgraph, err = loadTemplateSubgraph(ctx, store, templateID)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading template: %v\n", err)
			os.Exit(1)
		}

		// Check for missing variables
		requiredVars := extractAllVariables(subgraph)
		var missingVars []string
		for _, v := range requiredVars {
			if _, ok := vars[v]; !ok {
				missingVars = append(missingVars, v)
			}
		}
		if len(missingVars) > 0 {
			fmt.Fprintf(os.Stderr, "Error: missing required variables: %s\n", strings.Join(missingVars, ", "))
			fmt.Fprintf(os.Stderr, "Provide them with: --var %s=<value>\n", missingVars[0])
			os.Exit(1)
		}

		if dryRun {
			// Preview what would be created
			fmt.Printf("\nDry run: would create %d issues from template %s\n\n", len(subgraph.Issues), templateID)
			for _, issue := range subgraph.Issues {
				newTitle := substituteVariables(issue.Title, vars)
				suffix := ""
				if issue.ID == subgraph.Root.ID && assignee != "" {
					suffix = fmt.Sprintf(" (assignee: %s)", assignee)
				}
				fmt.Printf("  - %s (from %s)%s\n", newTitle, issue.ID, suffix)
			}
			if len(vars) > 0 {
				fmt.Printf("\nVariables:\n")
				for k, v := range vars {
					fmt.Printf("  {{%s}} = %s\n", k, v)
				}
			}
			return
		}

		// Clone the subgraph (deprecated command, non-wisp for backwards compatibility)
		opts := CloneOptions{
			Vars:     vars,
			Assignee: assignee,
			Actor:    actor,
			Ephemeral:     false,
		}
		var result *InstantiateResult
		if daemonClient != nil {
			result, err = cloneSubgraphViaDaemon(daemonClient, subgraph, opts)
		} else {
			result, err = cloneSubgraph(ctx, store, subgraph, opts)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error instantiating template: %v\n", err)
			os.Exit(1)
		}

		// Schedule auto-flush
		markDirtyAndScheduleFlush()

		if jsonOutput {
			outputJSON(result)
			return
		}

		fmt.Printf("%s Created %d issues from template\n", ui.RenderPass("✓"), result.Created)
		fmt.Printf("  New epic: %s\n", result.NewEpicID)
	},
}

func init() {
	templateInstantiateCmd.Flags().StringArray("var", []string{}, "Variable substitution (key=value)")
	templateInstantiateCmd.Flags().Bool("dry-run", false, "Preview what would be created")
	templateInstantiateCmd.Flags().String("assignee", "", "Assign the root epic to this agent/user")

	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateShowCmd)
	templateCmd.AddCommand(templateInstantiateCmd)
	rootCmd.AddCommand(templateCmd)
}

// =============================================================================
// Beads Template Functions
// =============================================================================

// loadTemplateSubgraph loads a template epic and all its descendants
// Deprecated: Use template.LoadSubgraph from commands/shared/template instead
func loadTemplateSubgraph(ctx context.Context, s storage.Storage, templateID string) (*TemplateSubgraph, error) {
	return template.LoadSubgraph(ctx, s, templateID)
}

// =============================================================================
// Proto Lookup Functions
// =============================================================================

// resolveProtoIDOrTitle resolves a proto by ID or title.
// Deprecated: Use template.ResolveProtoIDOrTitle from commands/shared/template instead
func resolveProtoIDOrTitle(ctx context.Context, s storage.Storage, input string) (string, error) {
	return template.ResolveProtoIDOrTitle(ctx, s, input)
}

// =============================================================================
// Daemon-compatible Template Functions
// =============================================================================

// loadTemplateSubgraphViaDaemon loads a template subgraph using daemon RPC calls
// Deprecated: Use template.LoadSubgraphViaDaemon from commands/shared/template instead
func loadTemplateSubgraphViaDaemon(client *rpc.Client, templateID string) (*TemplateSubgraph, error) {
	return template.LoadSubgraphViaDaemon(client, templateID)
}

// cloneSubgraphViaDaemon creates new issues from the template using daemon RPC calls
// Deprecated: Use template.CloneSubgraphViaDaemon from commands/shared/template instead
func cloneSubgraphViaDaemon(client *rpc.Client, subgraph *TemplateSubgraph, opts CloneOptions) (*InstantiateResult, error) {
	return template.CloneSubgraphViaDaemon(client, subgraph, opts)
}

// extractVariables finds all {{variable}} patterns in text
// Deprecated: Use template.ExtractVariables from commands/shared/template instead
func extractVariables(text string) []string {
	return template.ExtractVariables(text)
}

// extractAllVariables finds all variables across the entire subgraph
// Deprecated: Use template.ExtractAllVariables from commands/shared/template instead
func extractAllVariables(subgraph *TemplateSubgraph) []string {
	return template.ExtractAllVariables(subgraph)
}

// extractRequiredVariables returns only variables that don't have defaults.
// Deprecated: Use template.ExtractRequiredVariables from commands/shared/template instead
func extractRequiredVariables(subgraph *TemplateSubgraph) []string {
	return template.ExtractRequiredVariables(subgraph)
}

// applyVariableDefaults merges formula default values with provided variables.
// Deprecated: Use template.ApplyVariableDefaults from commands/shared/template instead
func applyVariableDefaults(vars map[string]string, subgraph *TemplateSubgraph) map[string]string {
	return template.ApplyVariableDefaults(vars, subgraph)
}

// substituteVariables replaces {{variable}} with values
// Deprecated: Use template.SubstituteVariables from commands/shared/template instead
func substituteVariables(text string, vars map[string]string) string {
	return template.SubstituteVariables(text, vars)
}

// generateBondedID creates a custom ID for dynamically bonded molecules.
// Deprecated: Use template.GenerateBondedID from commands/shared/template instead
func generateBondedID(oldID string, rootID string, opts CloneOptions) (string, error) {
	return template.GenerateBondedID(oldID, rootID, opts)
}

// extractIDSuffix extracts a suffix from an ID for use when IDs aren't hierarchical.
// Deprecated: Use template.ExtractIDSuffix from commands/shared/template instead
func extractIDSuffix(id string) string {
	return template.ExtractIDSuffix(id)
}

// getRelativeID extracts the relative portion of a child ID from its parent.
// Deprecated: Use template.GetRelativeID from commands/shared/template instead
func getRelativeID(oldID, rootID string) string {
	return template.GetRelativeID(oldID, rootID)
}

// cloneSubgraph creates new issues from the template with variable substitution.
// Deprecated: Use template.CloneSubgraph from commands/shared/template instead
func cloneSubgraph(ctx context.Context, s storage.Storage, subgraph *TemplateSubgraph, opts CloneOptions) (*InstantiateResult, error) {
	return template.CloneSubgraph(ctx, s, subgraph, opts)
}

// printTemplateTree prints the template structure as a tree
// This function stays in main package as it uses fmt.Printf for output
func printTemplateTree(subgraph *TemplateSubgraph, parentID string, depth int, isRoot bool) {
	template.PrintTree(subgraph, parentID, depth, isRoot, func(format string, args ...interface{}) {
		fmt.Printf(format, args...)
	})
}

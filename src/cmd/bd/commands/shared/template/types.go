// Package template provides shared template processing types and functions
// for the molecular chemistry CLI commands (pour, wisp, cook, mol, etc.).
//
// This package enables vertical slice architecture by providing reusable
// template functionality without circular dependencies.
package template

import (
	"regexp"

	"github.com/steveyegge/beads/internal/formula"
	"github.com/steveyegge/beads/internal/types"
)

// BeadsTemplateLabel is the label used to identify Beads-based templates
const BeadsTemplateLabel = "template"

// VariablePattern matches {{variable}} placeholders
var VariablePattern = regexp.MustCompile(`\{\{([a-zA-Z_][a-zA-Z0-9_]*)\}\}`)

// BondedIDPattern validates bonded IDs (alphanumeric, dash, underscore, dot)
var BondedIDPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+$`)

// Subgraph holds a template epic and all its descendants
type Subgraph struct {
	Root         *types.Issue              // The template epic
	Issues       []*types.Issue            // All issues in the subgraph (including root)
	Dependencies []*types.Dependency       // All dependencies within the subgraph
	IssueMap     map[string]*types.Issue   // ID -> Issue for quick lookup
	VarDefs      map[string]formula.VarDef // Variable definitions from formula (for defaults)
	Phase        string                    // Recommended phase: "liquid" (pour) or "vapor" (wisp)
}

// InstantiateResult holds the result of template instantiation
type InstantiateResult struct {
	NewEpicID string            `json:"new_epic_id"`
	IDMapping map[string]string `json:"id_mapping"` // old ID -> new ID
	Created   int               `json:"created"`    // number of issues created
}

// CloneOptions controls how the subgraph is cloned during spawn/bond
type CloneOptions struct {
	Vars      map[string]string // Variable substitutions for {{key}} placeholders
	Assignee  string            // Assign the root epic to this agent/user
	Actor     string            // Actor performing the operation
	Ephemeral bool              // If true, spawned issues are marked for bulk deletion
	Prefix    string            // Override prefix for ID generation (bd-hobo: distinct prefixes)

	// Dynamic bonding fields (for Christmas Ornament pattern)
	ParentID string // Parent molecule ID to bond under (e.g., "patrol-x7k")
	ChildRef string // Child reference with variables (e.g., "arm-{{polecat_name}}")
}

// IssueDetailsFromShow represents the response structure from daemon Show RPC
type IssueDetailsFromShow struct {
	types.Issue
	Labels       []string                              `json:"labels,omitempty"`
	Dependencies []*types.IssueWithDependencyMetadata `json:"dependencies,omitempty"`
	Dependents   []*types.IssueWithDependencyMetadata `json:"dependents,omitempty"`
}

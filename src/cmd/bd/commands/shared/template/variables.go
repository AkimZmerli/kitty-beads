package template

import (
	"fmt"
	"strings"

	"github.com/steveyegge/beads/internal/types"
)

// ExtractVariables finds all {{variable}} patterns in text
func ExtractVariables(text string) []string {
	matches := VariablePattern.FindAllStringSubmatch(text, -1)
	seen := make(map[string]bool)
	var vars []string
	for _, match := range matches {
		if len(match) >= 2 && !seen[match[1]] {
			vars = append(vars, match[1])
			seen[match[1]] = true
		}
	}
	return vars
}

// ExtractAllVariables finds all variables across the entire subgraph
func ExtractAllVariables(subgraph *Subgraph) []string {
	allText := ""
	for _, issue := range subgraph.Issues {
		allText += issue.Title + " " + issue.Description + " "
		allText += issue.Design + " " + issue.AcceptanceCriteria + " " + issue.Notes + " "
	}
	return ExtractVariables(allText)
}

// ExtractRequiredVariables returns only variables that don't have defaults.
// If VarDefs is available (from a cooked formula), uses it to filter out defaulted vars.
// Otherwise, falls back to returning all variables.
func ExtractRequiredVariables(subgraph *Subgraph) []string {
	allVars := ExtractAllVariables(subgraph)

	// If no VarDefs, assume all variables are required
	if subgraph.VarDefs == nil || len(subgraph.VarDefs) == 0 {
		return allVars
	}

	// Filter to only required variables (no default and marked as required, or not defined in VarDefs)
	var required []string
	for _, v := range allVars {
		def, exists := subgraph.VarDefs[v]
		// A variable is required if:
		// 1. It's not defined in VarDefs at all, OR
		// 2. It's defined with Required=true and no Default, OR
		// 3. It's defined with no Default (even if Required is false)
		if !exists {
			required = append(required, v)
		} else if def.Default == "" {
			required = append(required, v)
		}
		// If exists and has default, it's not required
	}
	return required
}

// ApplyVariableDefaults merges formula default values with provided variables.
// Returns a new map with defaults applied for any missing variables.
func ApplyVariableDefaults(vars map[string]string, subgraph *Subgraph) map[string]string {
	if subgraph.VarDefs == nil {
		return vars
	}

	result := make(map[string]string)
	for k, v := range vars {
		result[k] = v
	}

	// Apply defaults for missing variables
	for name, def := range subgraph.VarDefs {
		if _, exists := result[name]; !exists && def.Default != "" {
			result[name] = def.Default
		}
	}

	return result
}

// SubstituteVariables replaces {{variable}} with values
func SubstituteVariables(text string, vars map[string]string) string {
	return VariablePattern.ReplaceAllStringFunc(text, func(match string) string {
		// Extract variable name from {{name}}
		name := match[2 : len(match)-2]
		if val, ok := vars[name]; ok {
			return val
		}
		return match // Leave unchanged if not found
	})
}

// GenerateBondedID creates a custom ID for dynamically bonded molecules.
// When bonding a proto to a parent molecule, this generates IDs like:
//   - Root: parent.childref (e.g., "patrol-x7k.arm-ace")
//   - Children: parent.childref.step (e.g., "patrol-x7k.arm-ace.capture")
//
// The childRef is variable-substituted before use.
// Returns empty string if not a bonded operation (opts.ParentID empty).
func GenerateBondedID(oldID string, rootID string, opts CloneOptions) (string, error) {
	if opts.ParentID == "" {
		return "", nil // Not a bonded operation
	}

	// Substitute variables in childRef
	childRef := SubstituteVariables(opts.ChildRef, opts.Vars)

	// Validate childRef after substitution
	if childRef == "" {
		return "", fmt.Errorf("childRef is empty after variable substitution")
	}
	if !BondedIDPattern.MatchString(childRef) {
		return "", fmt.Errorf("invalid childRef '%s': must be alphanumeric, dash, underscore, or dot only", childRef)
	}

	if oldID == rootID {
		// Root issue: parent.childref
		newID := fmt.Sprintf("%s.%s", opts.ParentID, childRef)
		return newID, nil
	}

	// Child issue: parent.childref.relative
	// Extract the relative portion of the old ID (part after root)
	relativeID := GetRelativeID(oldID, rootID)
	if relativeID == "" {
		// No hierarchical relationship - use a suffix from the old ID to ensure uniqueness.
		// Extract the last part of the old ID (after any prefix or dash)
		suffix := ExtractIDSuffix(oldID)
		newID := fmt.Sprintf("%s.%s.%s", opts.ParentID, childRef, suffix)
		return newID, nil
	}

	newID := fmt.Sprintf("%s.%s.%s", opts.ParentID, childRef, relativeID)
	return newID, nil
}

// ExtractIDSuffix extracts a suffix from an ID for use when IDs aren't hierarchical.
// For "patrol-abc123", returns "abc123".
// For "bd-xyz.1", returns "1".
// This ensures child IDs remain unique when bonding.
func ExtractIDSuffix(id string) string {
	// First try to get the part after the last dot (for hierarchical IDs)
	if lastDot := strings.LastIndex(id, "."); lastDot >= 0 {
		return id[lastDot+1:]
	}
	// Otherwise, get the part after the last dash (for prefix-hash IDs)
	if lastDash := strings.LastIndex(id, "-"); lastDash >= 0 {
		return id[lastDash+1:]
	}
	// Fallback: use the whole ID
	return id
}

// GetRelativeID extracts the relative portion of a child ID from its parent.
// For example: GetRelativeID("bd-abc.step1.sub", "bd-abc") returns "step1.sub"
// Returns empty string if oldID equals rootID or doesn't start with rootID.
func GetRelativeID(oldID, rootID string) string {
	if oldID == rootID {
		return ""
	}
	// Check if oldID starts with rootID followed by a dot
	prefix := rootID + "."
	if strings.HasPrefix(oldID, prefix) {
		return oldID[len(prefix):]
	}
	return ""
}

// PrintTree prints the subgraph structure as a tree
func PrintTree(subgraph *Subgraph, parentID string, depth int, isRoot bool, printFunc func(format string, args ...interface{})) {
	indent := strings.Repeat("  ", depth)

	// Print root
	if isRoot {
		printFunc("%s   %s (root)\n", indent, subgraph.Root.Title)
	}

	// Find children of this parent
	var children []*types.Issue
	for _, dep := range subgraph.Dependencies {
		if dep.DependsOnID == parentID && dep.Type == types.DepParentChild {
			if child, ok := subgraph.IssueMap[dep.IssueID]; ok {
				children = append(children, child)
			}
		}
	}

	// Print children
	for i, child := range children {
		connector := "├──"
		if i == len(children)-1 {
			connector = "└──"
		}
		vars := ExtractVariables(child.Title)
		varStr := ""
		if len(vars) > 0 {
			varStr = fmt.Sprintf(" [%s]", strings.Join(vars, ", "))
		}
		printFunc("%s   %s %s%s\n", indent, connector, child.Title, varStr)
		PrintTree(subgraph, child.ID, depth+1, false, printFunc)
	}
}

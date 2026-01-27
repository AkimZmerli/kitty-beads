// Package dependencies provides the dependency management vertical slice.
package dependencies

import "github.com/steveyegge/beads/internal/types"

// AddArgs represents arguments for adding a dependency.
type AddArgs struct {
	FromID  string `json:"from_id"`
	ToID    string `json:"to_id"`
	DepType string `json:"dep_type"`
}

// RemoveArgs represents arguments for removing a dependency.
type RemoveArgs struct {
	FromID  string `json:"from_id"`
	ToID    string `json:"to_id"`
	DepType string `json:"dep_type,omitempty"`
}

// TreeArgs represents arguments for getting a dependency tree.
type TreeArgs struct {
	ID           string `json:"id"`
	MaxDepth     int    `json:"max_depth,omitempty"`
	ShowAllPaths bool   `json:"show_all_paths,omitempty"`
	Reverse      bool   `json:"reverse,omitempty"`
}

// AddResult represents the result of adding a dependency.
type AddResult struct {
	Status      string `json:"status"`
	IssueID     string `json:"issue_id"`
	DependsOnID string `json:"depends_on_id"`
	Type        string `json:"type"`
}

// RemoveResult represents the result of removing a dependency.
type RemoveResult struct {
	Status      string `json:"status"`
	IssueID     string `json:"issue_id"`
	DependsOnID string `json:"depends_on_id"`
}

// TreeResult represents a dependency tree.
type TreeResult struct {
	Nodes []*types.TreeNode `json:"nodes"`
}

// CycleResult represents detected cycles.
type CycleResult struct {
	Cycles [][]*types.Issue `json:"cycles"`
	Count  int              `json:"count"`
}

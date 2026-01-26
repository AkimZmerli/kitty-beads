// Package labels provides the label management vertical slice.
package labels

// AddArgs represents arguments for adding a label.
type AddArgs struct {
	IssueID string `json:"id"`
	Label   string `json:"label"`
}

// RemoveArgs represents arguments for removing a label.
type RemoveArgs struct {
	IssueID string `json:"id"`
	Label   string `json:"label"`
}

// ListArgs represents arguments for listing labels.
type ListArgs struct {
	IssueID string `json:"id"`
}

// AddResult represents the result of adding a label.
type AddResult struct {
	Status  string `json:"status"`
	IssueID string `json:"issue_id"`
	Label   string `json:"label"`
}

// RemoveResult represents the result of removing a label.
type RemoveResult struct {
	Status  string `json:"status"`
	IssueID string `json:"issue_id"`
	Label   string `json:"label"`
}

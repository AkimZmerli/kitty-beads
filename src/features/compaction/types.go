// Package compaction provides the issue compaction vertical slice.
package compaction

// CompactArgs represents arguments for the compact operation.
type CompactArgs struct {
	IssueID   string `json:"issue_id,omitempty"`
	Tier      int    `json:"tier"`
	DryRun    bool   `json:"dry_run"`
	Force     bool   `json:"force"`
	All       bool   `json:"all"`
	APIKey    string `json:"api_key,omitempty"`
	Workers   int    `json:"workers,omitempty"`
	BatchSize int    `json:"batch_size,omitempty"`
}

// CompactStatsArgs represents arguments for compact stats operation.
type CompactStatsArgs struct {
	Tier int `json:"tier,omitempty"`
}

// CompactResult represents the result of compacting a single issue.
type CompactResult struct {
	IssueID       string `json:"issue_id"`
	Success       bool   `json:"success"`
	Error         string `json:"error,omitempty"`
	OriginalSize  int    `json:"original_size,omitempty"`
	CompactedSize int    `json:"compacted_size,omitempty"`
	Reduction     string `json:"reduction,omitempty"`
}

// CompactResponse represents the response from a compact operation.
type CompactResponse struct {
	Success       bool               `json:"success"`
	IssueID       string             `json:"issue_id,omitempty"`
	Results       []CompactResult    `json:"results,omitempty"`
	Stats         *CompactStatsData  `json:"stats,omitempty"`
	OriginalSize  int                `json:"original_size,omitempty"`
	CompactedSize int                `json:"compacted_size,omitempty"`
	Reduction     string             `json:"reduction,omitempty"`
	Duration      string             `json:"duration,omitempty"`
	DryRun        bool               `json:"dry_run,omitempty"`
}

// CompactStatsData represents compaction statistics.
type CompactStatsData struct {
	Tier1Candidates  int    `json:"tier1_candidates"`
	Tier2Candidates  int    `json:"tier2_candidates"`
	TotalClosed      int    `json:"total_closed"`
	Tier1MinAge      string `json:"tier1_min_age"`
	Tier2MinAge      string `json:"tier2_min_age"`
	EstimatedSavings string `json:"estimated_savings,omitempty"`
}

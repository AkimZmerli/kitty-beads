// Package export provides the JSONL export vertical slice.
package export

// ExportArgs represents arguments for the export operation.
type ExportArgs struct {
	JSONLPath string `json:"jsonl_path"`
	Force     bool   `json:"force,omitempty"`
}

// ImportArgs represents arguments for the import operation.
type ImportArgs struct {
	JSONLPath string `json:"jsonl_path"`
	Force     bool   `json:"force,omitempty"`
}

// ExportResult represents the result of an export operation.
type ExportResult struct {
	Success      bool   `json:"success"`
	IssuesExported int  `json:"issues_exported"`
	Path         string `json:"path"`
	Duration     string `json:"duration,omitempty"`
}

// ImportResult represents the result of an import operation.
type ImportResult struct {
	Success      bool   `json:"success"`
	IssuesImported int  `json:"issues_imported"`
	Path         string `json:"path"`
	Duration     string `json:"duration,omitempty"`
}

// SyncStatus represents the sync status.
type SyncStatus struct {
	DirtyCount    int      `json:"dirty_count"`
	DirtyIssues   []string `json:"dirty_issues,omitempty"`
	LastExportAt  string   `json:"last_export_at,omitempty"`
	FileHash      string   `json:"file_hash,omitempty"`
}

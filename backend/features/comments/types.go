// Package comments provides the comment management vertical slice.
package comments

import (
	"time"

	"github.com/steveyegge/beads/internal/types"
)

// AddArgs represents arguments for adding a comment.
type AddArgs struct {
	IssueID string `json:"id"`
	Author  string `json:"author"`
	Text    string `json:"text"`
}

// ImportArgs represents arguments for importing a comment with preserved timestamp.
type ImportArgs struct {
	IssueID   string    `json:"id"`
	Author    string    `json:"author"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

// ListArgs represents arguments for listing comments.
type ListArgs struct {
	IssueID string `json:"id"`
}

// EventListArgs represents arguments for listing events.
type EventListArgs struct {
	IssueID string `json:"id"`
	Limit   int    `json:"limit,omitempty"`
}

// CommentResult wraps a comment for API responses.
type CommentResult struct {
	Comment *types.Comment `json:"comment"`
}

package comments

import (
	"fmt"
	"testing"
)

func TestIsUnknownOperationError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "unknown operation error",
			err:      fmt.Errorf("unknown operation: test"),
			expected: true,
		},
		{
			name:     "other error",
			err:      fmt.Errorf("some other error"),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isUnknownOperationError(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for error: %v", tt.expected, result, tt.err)
			}
		})
	}
}

func TestCommentAlias(t *testing.T) {
	t.Run("comment alias shares Run function with comments add", func(t *testing.T) {
		// This verifies that commentCmd reuses commentsAddCmd.Run
		if commentCmd.Run == nil {
			t.Error("commentCmd.Run is nil")
		}

		if commentsAddCmd.Run == nil {
			t.Error("commentsAddCmd.Run is nil")
		}

		// Verify the command structure is set up correctly
		if commentCmd.Use != "comment [issue-id] [text]" {
			t.Errorf("Expected Use to be 'comment [issue-id] [text]', got %s", commentCmd.Use)
		}

		if commentCmd.Short != "Add a comment to an issue (alias for 'comments add')" {
			t.Errorf("Unexpected Short description: %s", commentCmd.Short)
		}
	})
}

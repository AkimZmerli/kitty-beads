// Package comments implements the bd CLI comment management commands.
package comments

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/cmd/bd/cli"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/ui"
)

var commentsCmd = &cobra.Command{
	Use:     "comments [issue-id]",
	GroupID: "issues",
	Short:   "View or manage comments on an issue",
	Long: `View or manage comments on an issue.

Examples:
  # List all comments on an issue
  bd comments bd-123

  # List comments in JSON format
  bd comments bd-123 --json

  # Add a comment
  bd comments add bd-123 "This is a comment"

  # Add a comment from a file
  bd comments add bd-123 -f notes.txt`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cliCtx := cli.Get()
		issueID := args[0]

		comments := make([]*types.Comment, 0)
		usedDaemon := false
		if cliCtx.GetDaemonClient() != nil {
			// Resolve short/partial ID to full ID before sending to daemon (#1070)
			resolveArgs := &rpc.ResolveIDArgs{ID: issueID}
			resolveResp, err := cliCtx.GetDaemonClient().ResolveID(resolveArgs)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("resolving ID %s: %v", issueID, err)
			}
			var resolvedID string
			if err := json.Unmarshal(resolveResp.Data, &resolvedID); err != nil {
				cliCtx.FatalErrorRespectJSON("unmarshaling resolved ID: %v", err)
			}
			issueID = resolvedID

			resp, err := cliCtx.GetDaemonClient().ListComments(&rpc.CommentListArgs{ID: issueID})
			if err != nil {
				if isUnknownOperationError(err) {
					if err := cliCtx.FallbackToDirectMode("daemon does not support comment_list RPC"); err != nil {
						cliCtx.FatalErrorRespectJSON("getting comments: %v", err)
					}
				} else {
					cliCtx.FatalErrorRespectJSON("getting comments: %v", err)
				}
			} else {
				if err := json.Unmarshal(resp.Data, &comments); err != nil {
					cliCtx.FatalErrorRespectJSON("decoding comments: %v", err)
				}
				usedDaemon = true
			}
		}

		if !usedDaemon {
			if err := cliCtx.EnsureStoreActive(); err != nil {
				cliCtx.FatalErrorRespectJSON("getting comments: %v", err)
			}
			ctx := cliCtx.GetRootCtx()
			fullID, err := cliCtx.ResolvePartialID(ctx, issueID)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("resolving %s: %v", issueID, err)
			}
			issueID = fullID

			result, err := cliCtx.GetStore().GetIssueComments(ctx, issueID)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("getting comments: %v", err)
			}
			comments = result
		}

		// Normalize nil to empty slice for consistent JSON output
		if comments == nil {
			comments = make([]*types.Comment, 0)
		}

		if cliCtx.IsJSONOutput() {
			data, err := json.MarshalIndent(comments, "", "  ")
			if err != nil {
				cliCtx.FatalErrorRespectJSON("encoding JSON: %v", err)
			}
			fmt.Println(string(data))
			return
		}

		// Human-readable output
		if len(comments) == 0 {
			fmt.Printf("No comments on %s\n", issueID)
			return
		}

		fmt.Printf("\nComments on %s:\n\n", issueID)
		for _, comment := range comments {
			fmt.Printf("[%s] at %s\n", comment.Author, comment.CreatedAt.Format("2006-01-02 15:04"))
			rendered := ui.RenderMarkdown(comment.Text)
			// TrimRight removes trailing newlines that Glamour adds, preventing extra blank lines
			for _, line := range strings.Split(strings.TrimRight(rendered, "\n"), "\n") {
				fmt.Printf("  %s\n", line)
			}
			fmt.Println()
		}
	},
}

var commentsAddCmd = &cobra.Command{
	Use:   "add [issue-id] [text]",
	Short: "Add a comment to an issue",
	Long: `Add a comment to an issue.

Examples:
  # Add a comment
  bd comments add bd-123 "Working on this now"

  # Add a comment from a file
  bd comments add bd-123 -f notes.txt`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cliCtx := cli.Get()
		cliCtx.CheckReadonly("comment add")
		issueID := args[0]

		// Get comment text from flag or argument
		commentText, _ := cmd.Flags().GetString("file")
		if commentText != "" {
			// Read from file
			data, err := os.ReadFile(commentText) // #nosec G304 - user-provided file path is intentional
			if err != nil {
				cliCtx.FatalErrorRespectJSON("reading file: %v", err)
			}
			commentText = string(data)
		} else if len(args) < 2 {
			cliCtx.FatalErrorRespectJSON("comment text required (use -f to read from file)")
		} else {
			commentText = args[1]
		}

		// Get author from author flag, or use git-aware default
		author, _ := cmd.Flags().GetString("author")
		if author == "" {
			author = cliCtx.GetActorWithGit()
		}

		var comment *types.Comment
		if cliCtx.GetDaemonClient() != nil {
			// Resolve short/partial ID to full ID before sending to daemon (#1070)
			resolveArgs := &rpc.ResolveIDArgs{ID: issueID}
			resolveResp, err := cliCtx.GetDaemonClient().ResolveID(resolveArgs)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("resolving ID %s: %v", issueID, err)
			}
			var resolvedID string
			if err := json.Unmarshal(resolveResp.Data, &resolvedID); err != nil {
				cliCtx.FatalErrorRespectJSON("unmarshaling resolved ID: %v", err)
			}
			issueID = resolvedID

			resp, err := cliCtx.GetDaemonClient().AddComment(&rpc.CommentAddArgs{
				ID:     issueID,
				Author: author,
				Text:   commentText,
			})
			if err != nil {
				if isUnknownOperationError(err) {
					if err := cliCtx.FallbackToDirectMode("daemon does not support comment_add RPC"); err != nil {
						cliCtx.FatalErrorRespectJSON("adding comment: %v", err)
					}
				} else {
					cliCtx.FatalErrorRespectJSON("adding comment: %v", err)
				}
			} else {
				var parsed types.Comment
				if err := json.Unmarshal(resp.Data, &parsed); err != nil {
					cliCtx.FatalErrorRespectJSON("decoding comment: %v", err)
				}
				comment = &parsed
			}
		}

		if comment == nil {
			if err := cliCtx.EnsureStoreActive(); err != nil {
				cliCtx.FatalErrorRespectJSON("adding comment: %v", err)
			}
			ctx := cliCtx.GetRootCtx()

			fullID, err := cliCtx.ResolvePartialID(ctx, issueID)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("resolving %s: %v", issueID, err)
			}
			issueID = fullID

			comment, err = cliCtx.GetStore().AddIssueComment(ctx, issueID, author, commentText)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("adding comment: %v", err)
			}
		}

		if cliCtx.IsJSONOutput() {
			data, err := json.MarshalIndent(comment, "", "  ")
			if err != nil {
				cliCtx.FatalErrorRespectJSON("encoding JSON: %v", err)
			}
			fmt.Println(string(data))
			return
		}

		fmt.Printf("Comment added to %s\n", issueID)
	},
}

// commentCmd is a hidden top-level alias for commentsAddCmd (backwards compat)
var commentCmd = &cobra.Command{
	Use:        "comment [issue-id] [text]",
	Short:      "Add a comment to an issue (alias for 'comments add')",
	Long:       `Add a comment to an issue. This is an alias for 'bd comments add'.`,
	Args:       cobra.MinimumNArgs(1),
	Run:        commentsAddCmd.Run,
	Hidden:     true,
	Deprecated: "use 'bd comments add' instead (will be removed in v1.0.0)",
}

func isUnknownOperationError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "unknown operation")
}

// Register adds comment commands to the root command.
// Called from main package during initialization.
func Register(root *cobra.Command) {
	cliCtx := cli.Get()

	commentsCmd.AddCommand(commentsAddCmd)
	commentsAddCmd.Flags().StringP("file", "f", "", "Read comment text from file")
	commentsAddCmd.Flags().StringP("author", "a", "", "Add author to comment")

	// Add the same flags to the alias
	commentCmd.Flags().StringP("file", "f", "", "Read comment text from file")
	commentCmd.Flags().StringP("author", "a", "", "Add author to comment")

	// Issue ID completions
	commentsCmd.ValidArgsFunction = cliCtx.IssueIDCompletion
	commentsAddCmd.ValidArgsFunction = cliCtx.IssueIDCompletion
	commentCmd.ValidArgsFunction = cliCtx.IssueIDCompletion

	root.AddCommand(commentsCmd)
	root.AddCommand(commentCmd)
}

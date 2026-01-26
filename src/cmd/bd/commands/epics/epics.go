// Package epics implements the bd CLI epic management commands.
package epics

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/cmd/bd/cli"
	"github.com/steveyegge/beads/internal/rpc"
	"github.com/steveyegge/beads/internal/types"
	"github.com/steveyegge/beads/internal/ui"
)

var epicCmd = &cobra.Command{
	Use:     "epic",
	GroupID: "deps",
	Short:   "Epic management commands",
}

var epicStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show epic completion status",
	Run: func(cmd *cobra.Command, args []string) {
		cliCtx := cli.Get()
		eligibleOnly, _ := cmd.Flags().GetBool("eligible-only")

		var epics []*types.EpicStatus
		var err error

		if cliCtx.GetDaemonClient() != nil {
			resp, err := cliCtx.GetDaemonClient().EpicStatus(&rpc.EpicStatusArgs{
				EligibleOnly: eligibleOnly,
			})
			if err != nil {
				cliCtx.FatalErrorRespectJSON("communicating with daemon: %v", err)
			}
			if !resp.Success {
				cliCtx.FatalErrorRespectJSON("getting epic status: %s", resp.Error)
			}
			if err := json.Unmarshal(resp.Data, &epics); err != nil {
				cliCtx.FatalErrorRespectJSON("parsing response: %v", err)
			}
		} else {
			ctx := cliCtx.GetRootCtx()
			epics, err = cliCtx.GetStore().GetEpicsEligibleForClosure(ctx)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("getting epic status: %v", err)
			}
			if eligibleOnly {
				filtered := []*types.EpicStatus{}
				for _, epic := range epics {
					if epic.EligibleForClose {
						filtered = append(filtered, epic)
					}
				}
				epics = filtered
			}
		}

		if cliCtx.IsJSONOutput() {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(epics); err != nil {
				cliCtx.FatalErrorRespectJSON("encoding JSON: %v", err)
			}
			return
		}

		// Human-readable output
		if len(epics) == 0 {
			fmt.Println("No open epics found")
			return
		}

		for _, epicStatus := range epics {
			epic := epicStatus.Epic
			percentage := 0
			if epicStatus.TotalChildren > 0 {
				percentage = (epicStatus.ClosedChildren * 100) / epicStatus.TotalChildren
			}
			statusIcon := ""
			if epicStatus.EligibleForClose {
				statusIcon = ui.RenderPass("✓")
			} else if percentage > 0 {
				statusIcon = ui.RenderWarn("○")
			} else {
				statusIcon = "○"
			}
			fmt.Printf("%s %s %s\n", statusIcon, ui.RenderAccent(epic.ID), ui.RenderBold(epic.Title))
			fmt.Printf("   Progress: %d/%d children closed (%d%%)\n",
				epicStatus.ClosedChildren, epicStatus.TotalChildren, percentage)
			if epicStatus.EligibleForClose {
				fmt.Printf("   %s\n", ui.RenderPass("Eligible for closure"))
			}
			fmt.Println()
		}
	},
}

var closeEligibleEpicsCmd = &cobra.Command{
	Use:   "close-eligible",
	Short: "Close epics where all children are complete",
	Run: func(cmd *cobra.Command, args []string) {
		cliCtx := cli.Get()
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		// Block writes in readonly mode (closing modifies data)
		if !dryRun {
			cliCtx.CheckReadonly("epic close-eligible")
		}

		var eligibleEpics []*types.EpicStatus
		if cliCtx.GetDaemonClient() != nil {
			resp, err := cliCtx.GetDaemonClient().EpicStatus(&rpc.EpicStatusArgs{
				EligibleOnly: true,
			})
			if err != nil {
				cliCtx.FatalErrorRespectJSON("communicating with daemon: %v", err)
			}
			if !resp.Success {
				cliCtx.FatalErrorRespectJSON("getting eligible epics: %s", resp.Error)
			}
			if err := json.Unmarshal(resp.Data, &eligibleEpics); err != nil {
				cliCtx.FatalErrorRespectJSON("parsing response: %v", err)
			}
		} else {
			ctx := cliCtx.GetRootCtx()
			epics, err := cliCtx.GetStore().GetEpicsEligibleForClosure(ctx)
			if err != nil {
				cliCtx.FatalErrorRespectJSON("getting eligible epics: %v", err)
			}
			for _, epic := range epics {
				if epic.EligibleForClose {
					eligibleEpics = append(eligibleEpics, epic)
				}
			}
		}

		if len(eligibleEpics) == 0 {
			if !cliCtx.IsJSONOutput() {
				fmt.Println("No epics eligible for closure")
			} else {
				fmt.Println("[]")
			}
			return
		}

		if dryRun {
			if cliCtx.IsJSONOutput() {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				if err := enc.Encode(eligibleEpics); err != nil {
					cliCtx.FatalErrorRespectJSON("encoding JSON: %v", err)
				}
			} else {
				fmt.Printf("Would close %d epic(s):\n", len(eligibleEpics))
				for _, epicStatus := range eligibleEpics {
					fmt.Printf("  - %s: %s\n", epicStatus.Epic.ID, epicStatus.Epic.Title)
				}
			}
			return
		}

		// Actually close the epics
		closedIDs := []string{}
		for _, epicStatus := range eligibleEpics {
			if cliCtx.GetDaemonClient() != nil {
				resp, err := cliCtx.GetDaemonClient().CloseIssue(&rpc.CloseArgs{
					ID:     epicStatus.Epic.ID,
					Reason: "All children completed",
				})
				if err != nil || !resp.Success {
					errMsg := ""
					if err != nil {
						errMsg = err.Error()
					} else if !resp.Success {
						errMsg = resp.Error
					}
					fmt.Fprintf(os.Stderr, "Error closing %s: %s\n", epicStatus.Epic.ID, errMsg)
					continue
				}
			} else {
				ctx := cliCtx.GetRootCtx()
				err := cliCtx.GetStore().CloseIssue(ctx, epicStatus.Epic.ID, "All children completed", "system", "")
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error closing %s: %v\n", epicStatus.Epic.ID, err)
					continue
				}
			}
			closedIDs = append(closedIDs, epicStatus.Epic.ID)
		}

		if cliCtx.IsJSONOutput() {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			if err := enc.Encode(map[string]interface{}{
				"closed": closedIDs,
				"count":  len(closedIDs),
			}); err != nil {
				cliCtx.FatalErrorRespectJSON("encoding JSON: %v", err)
			}
		} else {
			fmt.Printf("✓ Closed %d epic(s)\n", len(closedIDs))
			for _, id := range closedIDs {
				fmt.Printf("  - %s\n", id)
			}
		}
	},
}

// Register adds epic commands to the root command.
// Called from main package during initialization.
func Register(root *cobra.Command) {
	epicCmd.AddCommand(epicStatusCmd)
	epicCmd.AddCommand(closeEligibleEpicsCmd)

	epicStatusCmd.Flags().Bool("eligible-only", false, "Show only epics eligible for closure")
	closeEligibleEpicsCmd.Flags().Bool("dry-run", false, "Preview what would be closed without making changes")

	root.AddCommand(epicCmd)
}

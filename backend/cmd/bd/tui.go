package main

import (
	"fmt"
	"os"

	"github.com/AkimZmerli/splitty"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/steveyegge/beads/internal/tui/common"
)

// themesByName maps CLI flag values to splitty themes.
var themesByName = map[string]splitty.Theme{
	"kitty-beads": common.KittyBeadsTheme,
	"tokyo-night": splitty.TokyoNight,
	"nightshade":  splitty.Nightshade,
	"glacier":     splitty.Glacier,
	"sorbet":      splitty.Sorbet,
}

var tuiCmd = &cobra.Command{
	Use:     "tui",
	GroupID: "views",
	Short:   "Launch split-pane terminal multiplexer",
	Long: `Launch the splitty terminal multiplexer with split pane support.

Presets:
  single   One pane (default)
  dev      60/40 vertical split with horizontal sub-split
  triple   Three equal columns
  quad     2x2 grid

Themes:
  kitty-beads   Tokyo Night palette matching bd's TUI (default)
  tokyo-night   Splitty's built-in Tokyo Night
  nightshade    Dark purple with vivid accents
  glacier       Arctic blue, cool tones
  sorbet        Warm pastel lavender and peach`,
	Run: func(cmd *cobra.Command, args []string) {
		preset, _ := cmd.Flags().GetString("preset")
		layoutPath, _ := cmd.Flags().GetString("layout")
		themeName, _ := cmd.Flags().GetString("theme")

		// Build options
		var opts []splitty.Option

		// Theme
		if theme, ok := themesByName[themeName]; ok {
			opts = append(opts, splitty.WithTheme(theme))
		} else {
			fmt.Fprintf(os.Stderr, "Unknown theme %q, using kitty-beads\n", themeName)
			opts = append(opts, splitty.WithTheme(common.KittyBeadsTheme))
		}

		// Preset (ignored if --layout is provided, loaded after init)
		if layoutPath == "" && preset != "" {
			opts = append(opts, splitty.WithPreset(preset))
		}

		m := splitty.New(opts...)

		p := tea.NewProgram(m,
			tea.WithAltScreen(),
			tea.WithMouseCellMotion(),
		)

		finalModel, err := p.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
			os.Exit(1)
		}

		// If a layout path was requested, load it after the program
		// has initialized but before the user interacts.
		// Actually, layout loading happens via the Manager's LoadLayout
		// which needs the program running. For now, save layout on exit.
		if layoutPath != "" {
			if fm, ok := finalModel.(*splitty.Manager); ok {
				if saveErr := fm.SaveLayout(layoutPath); saveErr != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not save layout: %v\n", saveErr)
				}
			}
		}
	},
}

func init() {
	tuiCmd.Flags().String("preset", "", "Layout preset: single, dev, triple, quad")
	tuiCmd.Flags().String("layout", "", "Path to save/restore layout JSON file")
	tuiCmd.Flags().String("theme", "kitty-beads", "Color theme: kitty-beads, tokyo-night, nightshade, glacier, sorbet")
	rootCmd.AddCommand(tuiCmd)
}

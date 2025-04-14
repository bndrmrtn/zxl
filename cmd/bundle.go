package cmd

import (
	"github.com/bndrmrtn/flare/internal/bundler"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// bundleCmd represents the bundle command
var bundleCmd = &cobra.Command{
	Use:   "bundle <folder>",
	Short: "Bundles all Flare Packages into a single archive",
	Run:   execBundle,
}

func init() {
	// Add the download command to the root command
	rootCmd.AddCommand(bundleCmd)
	bundleCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execBundle executes the bundle command
func execBundle(cmd *cobra.Command, args []string) {
	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	if len(args) != 1 {
		color.Red("Please provide a folder path")
		return
	}

	bndlr := bundler.New(args[0])
	if err := bndlr.Bundle("bundle"); err != nil {
		cmd.PrintErrln(err)
	}
}

package cmd

import (
	"github.com/bndrmrtn/zxl/pkg/pkgman"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes a new Zx project",
	Run:   execInit,
}

func init() {
	// Add the init command to the root command
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execInit executes the init command
func execInit(cmd *cobra.Command, args []string) {
	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	pm, err := pkgman.NewInitializer(".")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if err := pm.Save(); err != nil {
		cmd.PrintErrln(err)
		return
	}
}

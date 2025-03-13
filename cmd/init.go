package cmd

import (
	"github.com/bndrmrtn/zxl/pkg/pkgman"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <packageName>",
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
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErrln("No package name specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErrln("Only one package name can be specified")
		return
	}

	pm, err := pkgman.New(args[0], ".")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if err := pm.Save(); err != nil {
		cmd.PrintErrln(err)
		return
	}
}

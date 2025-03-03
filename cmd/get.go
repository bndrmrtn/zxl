package cmd

import (
	"github.com/bndrmrtn/zxl/pkg/pkgman"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:     "get packageUrl",
	Aliases: []string{"install", "i"},
	Short:   "Installs a Zx Package",
	Run:     execGet,
}

func init() {
	// Add the format command to the root command
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execGet executes the get command
func execGet(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErr("No url specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErr("Only one url can be specified")
		return
	}

	pm, err := pkgman.New(".")
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	if err := pm.Add(args[0]); err != nil {
		cmd.PrintErr(err)
		return
	}
}

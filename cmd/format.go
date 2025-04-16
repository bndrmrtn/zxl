package cmd

import (
	"os"

	"github.com/flarelang/flare/pkg/formatter"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// formatCmd represents the run command
var formatCmd = &cobra.Command{
	Use:     "format <folder>",
	Aliases: []string{"fmt"},
	Short:   "Format Flare (.fl) files in a directory",
	Run:     execFormat,
}

func init() {
	// Add the format command to the root command
	rootCmd.AddCommand(formatCmd)
	formatCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execRun executes the run command
func execFormat(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErrln("No directory specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErrln("Only one directory can be specified")
		return
	}

	info, err := os.Stat(args[0])
	if err != nil {
		cmd.PrintErrln("Error: " + err.Error())
		return
	}
	if !info.IsDir() {
		cmd.PrintErrln("Specified path is not a directory")
		return
	}

	f := formatter.New(args[0])
	if err := f.Format(); err != nil {
		cmd.PrintErrln("Error: " + err.Error())
		return
	}
}

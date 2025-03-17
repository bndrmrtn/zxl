package cmd

import (
	"github.com/bndrmrtn/zxl/pkg/pkgman"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <packageUrl>",
	Aliases: []string{"uninstall", "u"},
	Short:   "Uninstalls a Zx Package",
	Run:     execRemove,
}

func init() {
	// Add the format command to the root command
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execRemove executes the remove command
func execRemove(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErrln("No url specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErrln("Only one url can be specified")
		return
	}

	pm, err := pkgman.New(".")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if pm.PackageName == "" {
		cmd.PrintErrln("Package name cannot be empty\nUse `zxl init <packageName>` to initialize your package")
		return
	}

	if err := pm.Remove(args[0]); err != nil {
		cmd.PrintErrln(err)
		return
	}
}

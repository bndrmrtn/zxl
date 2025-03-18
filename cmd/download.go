package cmd

import (
	"github.com/bndrmrtn/zxl/pkg/pkgman"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download <packageUrl>",
	Short: "Downloads all Zx Packages from zxpack.yaml",
	Run:   execDownload,
}

func init() {
	// Add the download command to the root command
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execDownload executes the download command
func execDownload(cmd *cobra.Command, args []string) {
	colors := cmd.Flag("nocolor").Value.String() == "false"

	if !colors {
		color.NoColor = true
	}

	pm, err := pkgman.New(".")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	if err := pm.Download(); err != nil {
		cmd.PrintErrln(err)
		return
	}
}

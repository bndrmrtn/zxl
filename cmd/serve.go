package cmd

import (
	"net/http"
	"os"

	"github.com/bndrmrtn/zexlang/pkg/language"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var serveCmd = &cobra.Command{
	Use:     "serve folder",
	Aliases: []string{"server"},
	Short:   "Start an HTTP server to serve Zex (.zx) files from a directory",
	Run:     execServe,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
	serveCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
	serveCmd.Flags().StringP("listenAddr", "l", ":3000", "Server listen address")
}

// execRun executes the run command
func execServe(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"
	debug := cmd.Flag("debug").Value.String() == "true"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErr("No directory specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErr("Only one directory can be specified")
		return
	}

	info, err := os.Stat(args[0])
	if err != nil {
		cmd.PrintErr("Error: " + err.Error())
		return
	}
	if !info.IsDir() {
		cmd.PrintErr("Specified path is not a directory")
		return
	}

	var mode language.InterpreterMode
	if debug {
		mode = language.ModeDebug
	} else {
		mode = language.ModeProduction
	}

	interpreter := language.NewInterpreter(mode, true)
	if err := interpreter.Serve(args[0], cmd.Flag("listenAddr").Value.String(), colors); err != nil && err != http.ErrServerClosed {
		cmd.PrintErr("Error: " + err.Error())
		return
	}
}

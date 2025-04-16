package cmd

import (
	"net/http"
	"os"

	"github.com/flarelang/flare/pkg/language"
	"github.com/flarelang/flare/pkg/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:     "serve <folder or file>",
	Aliases: []string{"server"},
	Short:   "Start an HTTP server to serve Flare (.fl) files from a directory",
	Run:     execServe,
}

func init() {
	// Add the serve command to the root command
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolP("debug", "", false, "Run the program in debug mode")
	serveCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
	serveCmd.Flags().StringP("listenAddr", "l", ":3000", "Server listen address")
	serveCmd.Flags().BoolP("cache", "c", false, "Cache the parsed files for faster execution")
	serveCmd.Flags().BoolP("dev", "d", false, "Run the program in development mode")
}

// execRun executes the run command
func execServe(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"
	debug := cmd.Flag("debug").Value.String() == "true"
	cache := cmd.Flag("cache").Value.String() == "true"
	listenAddr := cmd.Flag("listenAddr").Value.String()
	dev := cmd.Flag("dev").Value.String() == "true"

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

	var mode language.InterpreterMode
	if debug {
		mode = language.ModeDebug
	} else {
		mode = language.ModeProduction
	}

	interpreter := language.NewInterpreter(mode, true)
	httpServer := server.New(interpreter, args[0], info.IsDir(), cache, colors, dev)

	if err := httpServer.Serve(listenAddr); err != nil && err != http.ErrServerClosed {
		cmd.PrintErrln("Error: " + err.Error())
		return
	}
}

package cmd

import (
	"net/http"
	"os"

	"github.com/bndrmrtn/zxl/pkg/language"
	"github.com/bndrmrtn/zxl/pkg/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var serveCmd = &cobra.Command{
	Use:     "serve folder",
	Aliases: []string{"server"},
	Short:   "Start an HTTP server to serve Zx (.zx) files from a directory",
	Run:     execServe,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
	serveCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
	serveCmd.Flags().StringP("listenAddr", "l", ":3000", "Server listen address")
	serveCmd.Flags().BoolP("cache", "c", false, "Cache the parsed files for faster execution")
}

// execRun executes the run command
func execServe(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"
	debug := cmd.Flag("debug").Value.String() == "true"
	cache := cmd.Flag("cache").Value.String() == "true"
	listenAddr := cmd.Flag("listenAddr").Value.String()

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
	httpServer := server.New(interpreter, args[0], cache, colors)

	if err := httpServer.Serve(listenAddr); err != nil && err != http.ErrServerClosed {
		cmd.PrintErr("Error: " + err.Error())
		return
	}
}

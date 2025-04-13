package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/bndrmrtn/flare/pkg/language"
	"github.com/bndrmrtn/flare/pkg/pkgman"
	"github.com/bndrmrtn/flare/pkg/server"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// startCmd represents the init command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a Flare project with package configuration",
	Run:   execStart,
}

func init() {
	// Add the start command to the root command
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
	startCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
}

// execStart executes the init command
func execStart(cmd *cobra.Command, args []string) {
	colors := cmd.Flag("nocolor").Value.String() == "false"
	debug := cmd.Flag("debug").Value.String() == "true"

	if !colors {
		color.NoColor = true
	}

	pm, err := pkgman.New(".")
	if err != nil {
		cmd.PrintErrln(err)
		return
	}

	var mode language.InterpreterMode
	if debug {
		mode = language.ModeDebug
	} else {
		mode = language.ModeProduction
	}

	entryAny, ok := pm.PackageConfig["entry"]
	if !ok {
		cmd.PrintErrln("Error: entry not found in package configuration")
		return
	}
	entry := fmt.Sprint(entryAny)

	switch pm.PackageType {
	case pkgman.TypeCLI:
		file, err := os.Open(entry)
		if err != nil {
			cmd.PrintErrln(err)
			return
		}
		defer file.Close()

		interpreter := language.NewInterpreter(mode, true)

		if _, err = interpreter.Interpret(entry, file); err != nil {
			cmd.PrintErrln(err)
			return
		}
	case pkgman.TypeWeb:
		info, err := os.Stat(entry)
		if err != nil {
			cmd.PrintErrln("Error: " + err.Error())
			return
		}

		var listenAddr = ":3000"
		if webConfigAny, ok := pm.PackageConfig["web"]; ok {
			if webConfig, ok := webConfigAny.(map[string]any); ok {
				host, hostOk := webConfig["host"]
				port, portOk := webConfig["port"]

				if hostOk && portOk {
					listenAddr = fmt.Sprintf("%s:%v", host, port)
				} else if hostOk {
					listenAddr = fmt.Sprintf("%s:3000", host)
				} else if portOk {
					listenAddr = fmt.Sprintf(":%v", port)
				}
			}
		}

		interpreter := language.NewInterpreter(mode, true)
		httpServer := server.New(interpreter, entry, info.IsDir(), true, colors)

		if err := httpServer.Serve(listenAddr); err != nil && err != http.ErrServerClosed {
			cmd.PrintErrln("Error: " + err.Error())
			return
		}
	case pkgman.TypeModule:
		cmd.PrintErrf("Cannot start type: %s\n", pm.PackageType)
	}
}

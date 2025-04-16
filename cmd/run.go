package cmd

import (
	"bytes"
	"errors"
	"io"
	"os"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/pkg/language"
	"github.com/flarelang/flare/pkg/prettycode"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run <file.fl>",
	Aliases: []string{"r", "exec"},
	Short:   "Interpret and execute Flare (.fl) files",
	Run:     execRun,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
	runCmd.Flags().BoolP("cache", "c", false, "Allow or disallow caching")
	runCmd.Flags().BoolP("nocolor", "n", false, "Enable or disable colorized output")
}

// execRun executes the run command
func execRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		_ = cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"
	debug := cmd.Flag("debug").Value.String() == "true"
	caching := cmd.Flag("cache").Value.String() == "true"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErrln("No file specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErrln("Only one file can be run at a time")
		return
	}

	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		cmd.PrintErrln("File does not exist")
		return
	}

	var mode language.InterpreterMode
	if debug {
		mode = language.ModeDebug
	} else {
		mode = language.ModeProduction
	}

	file, err := os.Open(args[0])
	if err != nil {
		cmd.PrintErrln(err)
		return
	}
	defer file.Close()

	interpreter := language.NewInterpreter(mode, caching)

	if _, err = interpreter.Interpret(args[0], file); err != nil {
		var de errs.DebugError

		if errors.As(err, &de) {
			s := de.PrettyError(func(r io.Reader) string {
				b, _ := io.ReadAll(r)

				pc, err := prettycode.New(bytes.NewReader(b))
				if err != nil {
					return string(b)
				}

				return pc.HighlightConsole()
			})

			cmd.PrintErrln(errors.New(s))
			return
		}

		cmd.PrintErrln(err)
		return
	}
}

package cmd

import (
	"os"

	"github.com/bndrmrtn/zxl/pkg/language"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:     "run <file.zx>",
	Aliases: []string{"r", "exec"},
	Short:   "Interpret and execute Zx (.zx) files",
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
		cmd.Help()
		return
	}

	colors := cmd.Flag("nocolor").Value.String() == "false"
	debug := cmd.Flag("debug").Value.String() == "true"
	caching := cmd.Flag("cache").Value.String() == "true"

	if !colors {
		color.NoColor = true
	}

	if len(args) == 0 {
		cmd.PrintErr("No file specified")
		return
	}

	if len(args) > 1 {
		cmd.PrintErr("Only one file can be run at a time")
		return
	}

	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		cmd.PrintErr("File does not exist")
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
		cmd.PrintErr(err)
		return
	}
	defer file.Close()

	interpreter := language.NewInterpreter(mode, caching)

	if _, err = interpreter.Interpret(args[0], file); err != nil {
		cmd.PrintErr(err)
		return
	}
}

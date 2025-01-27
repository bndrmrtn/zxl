package cmd

import (
	"os"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/ast"
	"github.com/bndrmrtn/zexlang/internal/lexer"
	"github.com/bndrmrtn/zexlang/internal/runtime"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var runCmd = &cobra.Command{
	Use:     "run filename.zx",
	Aliases: []string{"r", "exec"},
	Short:   "Interpret and execute Zex (.zx) files",
	Run:     execRun,
}

func init() {
	// Add the run command to the root command
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolP("debug", "d", false, "Run the program in debug mode")
	runCmd.Flags().BoolP("color", "c", true, "Enable or disable colorized output")
}

func execRun(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		cmd.Help()
		return
	}

	debug := cmd.Flag("debug").Value.String() == "true"

	for _, arg := range args {
		if !strings.HasSuffix(arg, ".zx") {
			cmd.Println("Zex can only run files that has .zx extesion.")
			return
		}
	}

	file, err := os.Open(args[0])
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	// Tokenize the source code with lexer
	lx := lexer.New(args[0])
	tokens, err := lx.Parse(file)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	if debug {
		writeDebug("lexer.yaml", tokens)
	}

	// Generate abstract syntax tree from tokens
	builder := ast.NewBuilder()
	nodes, err := builder.Build(tokens)
	if err != nil {
		cmd.PrintErr(err)
		return
	}

	if debug {
		writeDebug("ast.yaml", nodes)
	}

	run := runtime.New()
	if _, err = run.Execute(nodes); err != nil {
		cmd.PrintErr(err)
		return
	}
}

func writeDebug(file string, v any) {
	f, err := os.Create("debug/" + file)
	if err != nil {
		return
	}
	defer f.Close()

	_ = yaml.NewEncoder(f).Encode(v)
}

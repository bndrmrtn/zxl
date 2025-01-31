package cmd

import (
	"log"

	"github.com/bndrmrtn/zexlang/internal/version"
	"github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
	Use:     "zexlang",
	Short:   "Zex ✨ A small programming language for template rendering.",
	Version: version.Version,
}

// Execute executes the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

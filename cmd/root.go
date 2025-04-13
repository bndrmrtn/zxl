package cmd

import (
	"log"

	"github.com/bndrmrtn/flare/internal/version"
	"github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
	Use:     "flare",
	Short:   "Flare ✨ A small programming language for template rendering.",
	Version: version.Version,
}

// Execute executes the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

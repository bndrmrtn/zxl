package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// rootCmd is the root command for the CLI
var rootCmd = &cobra.Command{
	Use:   "zex",
	Short: "Zex ✨ A small programming language for template rendering.",
}

// Execute executes the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

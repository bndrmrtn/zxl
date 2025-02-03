// ZexLang is a simple templating programming language
// that is designed to be easy to use and understand.
//
// Zex is built in Go and is designed to be a simple.
// Created by Martin Binder
package main

import (
	"fmt"
	"os"

	"github.com/bndrmrtn/zexlang/cmd"
	"github.com/fatih/color"
)

func main() {
	// Execute the command
	if os.Getenv("DEBUG") != "true" {
		defer fatal()
	}
	cmd.Execute()
}

// fatal is a helper function to recover from a panic
func fatal() {
	if r := recover(); r != nil {
		f := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s: %v\n", f("Fatal error"), r)
	}
}

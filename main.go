// ZexLang is a simple templating programming language
// that is designed to be easy to use and understand.
//
// Zex is built in Go and is designed to be a simple.
// Created by Martin Binder
// Website: https://mrtn.vip
package main

import (
	"fmt"
	"os"

	"github.com/bndrmrtn/zexlang/cmd"
	"github.com/fatih/color"
)

// main is the entry point of the ZexLang program
func main() {
	// Recover from a panic if the program is not in debug mode
	if os.Getenv("DEBUG") != "true" {
		defer fatal()
	}
	// Execute the command
	cmd.Execute()
}

// fatal is a helper function to recover from a panic
func fatal() {
	if r := recover(); r != nil {
		f := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s: %v\n", f("Fatal error"), r)
	}
}

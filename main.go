// FlareLang is a simple templating programming language
// that is designed to be easy to use and understand.
//
// Flare is built in Go and is designed to be a simple.
// Created by Martin Binder
// Website: https://mrtn.vip
package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/flarelang/flare/cmd"
	"github.com/fatih/color"
	"go.uber.org/zap"
)

var logger *zap.Logger

// initLogger initializes the logger
func initLogger() {
	if os.Getenv("DEBUG") == "true" {
		logger, _ = zap.NewDevelopment()
	} else {
		logger = zap.NewNop()
	}
	zap.ReplaceGlobals(logger)
}

// main is the entry point of the FlareLang program
func main() {
	// Recover from a panic if the program is not in debug mode
	if os.Getenv("DEBUG") != "true" {
		defer fatal()
	}

	initLogger()
	defer logger.Sync()

	// Execute the command
	cmd.Execute()
}

// fatal is a helper function to recover from a panic
func fatal() {
	if r := recover(); r != nil {
		f := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Printf("%s: %v\n", f("Fatal error"), r)

		// Show Go's stack trace if the SHOW_STACK environment variable is set to true
		if os.Getenv("SHOW_STACK") == "true" {
			fmt.Printf("Stack trace: %s\n", debug.Stack())
		}
	}
}

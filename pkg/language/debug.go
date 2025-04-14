package language

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// InterpreterMode is the mode of the interpreter
type InterpreterMode int

const (
	// ModeDebug writes debug information to file
	ModeDebug InterpreterMode = iota
	// ModeProduction is the default mode
	ModeProduction
	// ModeTest is the test mode
	ModeTest
)

// writeDebug writes debug information to file
func (ir *Interpreter) writeDebug(file, suffix string, v any) {
	// Create debug directory if it does not exist
	_ = os.MkdirAll(".flare/", os.ModePerm)

	file = strings.ReplaceAll(file, "/", ".")
	file = strings.ReplaceAll(file, "\\", ".")
	file = strings.Trim(file, ".")
	file = "debug-" + file + suffix + ".yaml"

	// Write debug information to file
	f, err := os.Create(".flare/" + file)
	if err != nil {
		return
	}
	defer f.Close()

	// Write debug information to file
	_ = yaml.NewEncoder(f).Encode(v)
}

package language

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/ast"
	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/lexer"
	"github.com/bndrmrtn/zexlang/internal/runtime"
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

// Interpreter is a language interpreter
type Interpreter struct {
	mode InterpreterMode
}

// NewInterpreter creates a new interpreter
func NewInterpreter(mode InterpreterMode) *Interpreter {
	return &Interpreter{
		mode: mode,
	}
}

// Interpret interprets the given data
func (ir *Interpreter) Interpret(fileName string, data io.Reader) ([]*builtin.FuncReturn, error) {
	if !strings.HasSuffix(fileName, ".zx") {
		return nil, fmt.Errorf("Zex can only run files that has .zx extesion.")
	}

	// Tokenize the source code with lexer
	lx := lexer.New(fileName)
	ts, err := lx.Parse(data)
	if err != nil {
		return nil, err
	}

	if ir.mode == ModeDebug {
		ir.writeDebug(fileName, "lexer", ts)
	}

	// Build the abstract syntax tree from tokens
	builder := ast.NewBuilder()
	nodes, err := builder.Build(ts)
	if err != nil {
		return nil, err
	}

	if ir.mode == ModeDebug {
		ir.writeDebug(fileName, "ast", nodes)
	}

	// Write debug information to file
	return runtime.New(runtime.EntryPoint).Execute(nodes)
}

// writeDebug writes debug information to file
func (ir *Interpreter) writeDebug(file, suffix string, v any) {
	// Create debug directory if it does not exist
	_ = os.MkdirAll("debug/", os.ModePerm)

	file = strings.ReplaceAll(file, "/", ".")
	file = strings.Trim(file, ".")
	file = file + "." + suffix + ".yaml"

	// Write debug information to file
	f, err := os.Create("debug/" + file)
	if err != nil {
		return
	}
	defer f.Close()

	// Write debug information to file
	_ = yaml.NewEncoder(f).Encode(v)
}

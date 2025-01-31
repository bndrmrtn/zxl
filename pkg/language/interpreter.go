package language

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/ast"
	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/cache"
	"github.com/bndrmrtn/zexlang/internal/lexer"
	"github.com/bndrmrtn/zexlang/internal/models"
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
	// mode is the mode of the interpreter
	mode InterpreterMode
	// cache is the cache flag
	cache bool
}

// NewInterpreter creates a new interpreter
func NewInterpreter(mode InterpreterMode, cache bool) *Interpreter {
	return &Interpreter{
		mode:  mode,
		cache: cache,
	}
}

// Interpret interprets the given data
func (ir *Interpreter) Interpret(fileName string, data io.Reader) ([]*builtin.FuncReturn, error) {
	if !strings.HasSuffix(fileName, ".zx") {
		return nil, fmt.Errorf("Zex can only run files that has .zx extesion.")
	}

	// Get the nodes from the given data
	nodes, err := ir.getNodes(fileName, data)
	if err != nil {
		return nil, err
	}

	// Write debug information to file
	return runtime.New(runtime.EntryPoint).Execute(nodes)
}

// GetNodes gets the nodes from the given data
func (ir *Interpreter) GetNodes(fileName string, data io.Reader) ([]*models.Node, error) {
	if !strings.HasSuffix(fileName, ".zx") {
		return nil, fmt.Errorf("Zex can only run files that has .zx extesion.")
	}
	return ir.getNodes(fileName, data)
}

func (ir *Interpreter) Serve(dir string, addr string, colors bool) error {
	server := NewServer(ir, dir, colors)
	return server.Serve(addr)
}

// getNodes gets the nodes from the given data
func (ir *Interpreter) getNodes(fileName string, data io.Reader) ([]*models.Node, error) {
	b, err := io.ReadAll(data)
	if err != nil {
		return nil, err
	}

	// Get the nodes from cache if it exists
	if ir.cache {
		if nodes, ok := cache.Get(fileName, b); ok {
			return nodes, nil
		}
	}

	// Tokenize the source code with lexer
	lx := lexer.New(fileName)
	ts, err := lx.Parse(bytes.NewReader(b))
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

	// Store cache information
	if ir.cache {
		cache.Store(fileName, b, nodes)
	}

	return nodes, nil
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

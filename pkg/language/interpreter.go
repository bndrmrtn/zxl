package language

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/bndrmrtn/zxl/internal/ast"
	"github.com/bndrmrtn/zxl/internal/cache"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/runtimev2"
	"github.com/bndrmrtn/zxl/lang"
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
func (ir *Interpreter) Interpret(fileName string, data io.Reader) (lang.Object, error) {
	if !strings.HasSuffix(fileName, ".zx") {
		return nil, fmt.Errorf("Zx can only run files that has .zx extesion.")
	}

	// Get the nodes from the given data
	nodes, err := ir.getNodes(fileName, data)
	if err != nil {
		return nil, err
	}

	run, err := runtimev2.New()
	if err != nil {
		return nil, err
	}

	return run.Execute(nodes)
}

// GetNodes gets the nodes from the given data
func (ir *Interpreter) GetNodes(fileName string, data io.Reader) ([]*models.Node, error) {
	if !strings.HasSuffix(fileName, ".zx") {
		return nil, fmt.Errorf("Zx can only run files that has .zx extesion.")
	}
	return ir.getNodes(fileName, data)
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

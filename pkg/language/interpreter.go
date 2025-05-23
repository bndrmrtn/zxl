package language

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/flarelang/flare/internal/ast"
	"github.com/flarelang/flare/internal/cache"
	"github.com/flarelang/flare/internal/lexer"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/runtimev2"
	"github.com/flarelang/flare/internal/state"
	"github.com/flarelang/flare/lang"
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
	if !strings.HasSuffix(fileName, ".fl") && !strings.HasSuffix(fileName, ".flb") && !strings.HasSuffix(fileName, ".flare") {
		return nil, fmt.Errorf("Flare can only run files that has .fl, .flb or .flare extesion.")
	}

	// Get the nodes from the given data
	nodes, err := ir.getNodes(fileName, data)
	if err != nil {
		return nil, err
	}

	run, err := runtimev2.New(state.Default())
	if err != nil {
		return nil, err
	}

	return run.Execute(nodes)
}

// GetNodes gets the nodes from the given data
func (ir *Interpreter) GetNodes(fileName string, data io.Reader) ([]*models.Node, error) {
	if !strings.HasSuffix(fileName, ".fl") {
		return nil, fmt.Errorf("Flare can only run files that has .fl extesion.")
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

package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// RuntimeMode is the runtime mode
type Runtime struct {
	// mode is the runtime mode
	mode RuntimeMode

	// funcs is a map of function names to functions
	funcs map[string]builtin.Function
	// variables is a map of variable names to variables
	variables map[string]*builtin.Variable

	// pkgs is a map of package names to packages
	pkgs map[string]builtin.Package

	// executers is a map of namespace names to exec
	executers map[string]*Executer
}

// New creates a new runtime
func New(mode RuntimeMode) *Runtime {
	return &Runtime{
		mode:      mode,
		funcs:     builtin.GetBuiltins(),
		variables: make(map[string]*builtin.Variable),
		pkgs:      make(map[string]builtin.Package),
		executers: make(map[string]*Executer),
	}
}

// Execute executes the given nodes
func (r *Runtime) Execute(nodes []*models.Node) (*builtin.FuncReturn, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	if nodes[0].Type == tokens.Namespace {

	}

	return r.Exec(ExecuterScopeGlobal, nil, "", nodes)
}

// Exec executes the given nodes in the given namespace
func (r *Runtime) Exec(scope ExecuterScope, parent *Executer, namespace string, nodes []*models.Node) (*builtin.FuncReturn, error) {
	ex, ok := r.executers[namespace]
	if !ok {
		ex = NewExecuter(scope, r, parent)
		r.executers[namespace] = ex
	}
	return ex.Execute(nodes)
}

// GetNamespaceExecuter gets the executer for the given namespace
func (r *Runtime) GetNamespaceExecuter(namespace string) (*Executer, error) {
	ex, ok := r.executers[namespace]
	if !ok {
		return nil, fmt.Errorf("namespace %v not found", namespace)
	}
	return ex, nil
}

// BindModule binds the given package to the given name
func (r *Runtime) BindModule(name string, pkg builtin.Package) {
	r.pkgs[name] = pkg
}

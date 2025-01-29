package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

type Runtime struct {
	mode RuntimeMode

	// funcs is a map of function names to functions
	funcs map[string]builtin.Function
	// variables is a map of variable names to variables
	variables map[string]*builtin.Variable

	executers map[string]*Executer
}

func New(mode RuntimeMode) *Runtime {
	return &Runtime{
		mode:      mode,
		funcs:     builtin.GetBuiltins(),
		variables: make(map[string]*builtin.Variable),
		executers: make(map[string]*Executer),
	}
}

func (r *Runtime) Execute(nodes []*models.Node) ([]*builtin.FuncReturn, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	if nodes[0].Type == tokens.Namespace {

	}

	return r.Exec(ExecuterScopeGlobal, nil, "", nodes)
}

func (r *Runtime) Exec(scope ExecuterScope, parent *Executer, namespace string, nodes []*models.Node) ([]*builtin.FuncReturn, error) {
	ex, ok := r.executers[namespace]
	if !ok {
		ex = NewExecuter(scope, r, parent)
		r.executers[namespace] = ex
	}
	return ex.Execute(nodes)
}

func (r *Runtime) GetNamespaceExecuter(namespace string) (*Executer, error) {
	ex, ok := r.executers[namespace]
	if !ok {
		return nil, fmt.Errorf("namespace %v not found", namespace)
	}
	return ex, nil
}

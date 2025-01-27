package runtime

import (
	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/models"
)

type Runtime struct {
	// funcs is a map of function names to functions
	funcs map[string]builtin.Function
	// variables is a map of variable names to variables
	variables map[string]*builtin.Variable
}

func New() *Runtime {
	return &Runtime{
		funcs:     builtin.GetBuiltins(),
		variables: make(map[string]*builtin.Variable),
	}
}

func (r *Runtime) Execute(tokens []*models.Node) ([]*builtin.FuncReturn, error) {
	ex := NewExecuter(r, nil)
	_, err := ex.Execute(tokens)
	if err != nil {
		return nil, err
	}
	return ex.ExecuteFn("main", nil)
}

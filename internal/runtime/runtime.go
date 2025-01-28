package runtime

import (
	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/models"
)

type Runtime struct {
	mode RuntimeMode

	// funcs is a map of function names to functions
	funcs map[string]builtin.Function
	// variables is a map of variable names to variables
	variables map[string]*builtin.Variable
}

func New(mode RuntimeMode) *Runtime {
	return &Runtime{
		mode:      mode,
		funcs:     builtin.GetBuiltins(),
		variables: make(map[string]*builtin.Variable),
	}
}

func (r *Runtime) Execute(tokens []*models.Node) ([]*builtin.FuncReturn, error) {
	ex := NewExecuter(r, nil)
	return ex.Execute(tokens)
}

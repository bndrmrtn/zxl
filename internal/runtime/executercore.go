package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (e *Executer) executeFn(token *models.Node) ([]*builtin.FuncReturn, error) {
	name := token.Content
	args := token.Args

	fn, ok := e.fns[name]
	if ok {
		ex := NewExecuter(e.runtime, e)
		for _, arg := range args {
			if arg.VariableType == tokens.ReferenceVariable {
				arg, ok = e.vars[arg.Content]
				if !ok {
					return nil, fmt.Errorf("variable %v not found", arg.Content)
				}
			}
			ex.Bind(arg)
		}
		return ex.Execute(fn.Children)
	}

	builtinFn, ok := e.runtime.funcs[name]
	if ok {
		var convArgs []*builtin.Variable
		for _, arg := range args {
			if arg.VariableType == tokens.ReferenceVariable {
				arg, ok = e.vars[arg.Content]
				if !ok {
					return nil, fmt.Errorf("variable %v not found", arg.Content)
				}
			}

			convArgs = append(convArgs, &builtin.Variable{
				Type:  arg.VariableType,
				Value: arg.Value,
			})
		}

		return builtinFn(convArgs)
	}

	return nil, fmt.Errorf("function %v not found", name)
}

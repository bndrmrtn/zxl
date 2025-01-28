package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (e *Executer) GetVariableValue(name string) (*models.Node, error) {
	v, ok := e.vars[name]
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("variable %v cannot be referenced", name), nil)
	}

	if v.Reference {
		return e.GetVariableValue(v.Content)
	}

	if v.VariableType == tokens.ExpressionVariable {
		return e.evaluateExpression(v)
	}

	return v, nil
}

func (e *Executer) executeFn(token *models.Node) ([]*builtin.FuncReturn, error) {
	name := token.Content
	args := token.Args

	fn, ok := e.fns[name]
	if ok {
		if len(args) != len(fn.Args) {
			return nil, errs.WithDebug(fmt.Errorf("function %v expects %v arguments, got %v", name, len(fn.Args), len(args)), fn.Debug)
		}

		ex := NewExecuter(e.runtime, e)
		for i, arg := range args {
			if arg.VariableType == tokens.ReferenceVariable {
				arg, err := e.GetVariableValue(arg.Content)
				if err != nil {
					return nil, err
				}
				arg.Content = fn.Args[i].Content
				ex.Bind(arg)
				continue
			}

			arg.Content = fn.Args[i].Content
			ex.Bind(arg)
		}
		ret, err := ex.Execute(fn.Children)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		return ret, nil
	}

	builtinFn, ok := e.runtime.funcs[name]
	if ok {
		var convArgs []*builtin.Variable

		for _, arg := range args {
			if arg.VariableType == tokens.ReferenceVariable {
				var err error

				arg, err = e.GetVariableValue(arg.Content)
				if err != nil {
					return nil, errs.WithDebug(err, token.Debug)
				}
			}

			if arg.VariableType == tokens.InlineValue {
				arg.VariableType = arg.Type.ToVariableType()
			}

			convArgs = append(convArgs, &builtin.Variable{
				Type:  arg.VariableType,
				Value: arg.Value,
			})
		}

		ret, err := builtinFn(convArgs)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		return ret, nil
	}

	return nil, errs.WithDebug(fmt.Errorf("function %v not found", name), token.Debug)
}

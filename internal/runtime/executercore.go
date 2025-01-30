package runtime

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// GetVariableValue gets the value of a variable
func (e *Executer) GetVariableValue(name string) (*models.Node, error) {
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		last := parts[len(parts)-1]
		parts = parts[:len(parts)-1]
		ex, _, err := e.accessUnderlyingVariable(parts)
		if err != nil {
			return nil, err
		}

		return ex.GetVariableValue(last)
	}

	v, ok := e.vars[name]
	if !ok {
		if e.scope == ExecuterScopeBlock {
			return e.parent.GetVariableValue(name)
		}
		return nil, fmt.Errorf("variable '%v' cannot be referenced", name)
	}

	if v.Reference {
		return e.GetVariableValue(v.Content)
	}

	if v.VariableType == tokens.ReferenceVariable {
		return e.GetVariableValue(v.Value.(string))
	}

	if v.VariableType == tokens.ExpressionVariable {
		return e.evaluateExpression(v)
	}

	return v, nil
}

// executeFn executes a function
func (e *Executer) executeFn(token *models.Node) ([]*builtin.FuncReturn, error) {
	name := token.Content
	args := token.Args

	switch name {
	case "eval":
		convArgs, err := e.convertArgument(args)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		return e.runFuncEval(token.Debug, convArgs)
	case "import":
		convArgs, err := e.convertArgument(args)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		return e.runFuncImport(token.Debug, convArgs)
	}

	// if
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		last := parts[len(parts)-1]
		parts = parts[:len(parts)-1]
		return e.executeComplexFuncCall(token, parts, last, args)
	}

	fn, ok := e.fns[name]
	if ok {
		if len(args) != len(fn.Args) {
			return nil, errs.WithDebug(fmt.Errorf("function %v expects %v arguments, got %v", name, len(fn.Args), len(args)), fn.Debug)
		}

		ex := NewExecuter(ExecuterScopeFunction, e.runtime, e)
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
		convArgs, err := e.convertArgument(args)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		ret, err := builtinFn(convArgs)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		return ret, nil
	}

	block, ok := e.blocks[name]
	if ok {
		return e.newBlock(block, args)
	}

	return nil, errs.WithDebug(fmt.Errorf("function %v not found", name), token.Debug)
}

// executeComplexFuncCall executes a function with a complex name
func (e *Executer) executeComplexFuncCall(token *models.Node, args []string, fn string, variables []*models.Node) ([]*builtin.FuncReturn, error) {
	if fn == "construct" {
		return nil, errs.WithDebug(fmt.Errorf("construct is a reserved method"), token.Debug)
	}

	ex, _, err := e.accessUnderlyingVariable(args)
	if err != nil {
		return nil, errs.WithDebug(err, token.Debug)
	}

	convertArgs, err := e.convertArgument(variables)
	if err != nil {
		return nil, errs.WithDebug(err, token.Debug)
	}

	return ex.ExecuteFn(fn, convertArgs)
}

// convertArgument converts arguments to builtin variables
func (e *Executer) convertArgument(args []*models.Node) ([]*builtin.Variable, error) {
	var convArgs []*builtin.Variable

	for _, arg := range args {
		if arg.VariableType == tokens.ReferenceVariable {
			var err error

			arg, err = e.GetVariableValue(arg.Content)
			if err != nil {
				return nil, err
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

	return convArgs, nil
}

// newBlock creates a new block
func (e *Executer) newBlock(block *models.Node, args []*models.Node) ([]*builtin.FuncReturn, error) {
	ex := NewExecuter(ExecuterScopeDefinition, e.runtime, e)
	_, err := ex.Execute(block.Children)
	if err != nil {
		return nil, errs.WithDebug(err, block.Debug)
	}

	if _, ok := ex.fns["construct"]; ok {
		convArgs, err := e.convertArgument(args)
		if err != nil {
			return nil, errs.WithDebug(err, block.Debug)
		}
		_, err = ex.ExecuteFn("construct", convArgs)
		if err != nil {
			return nil, errs.WithDebug(err, block.Debug)
		}
	}

	return []*builtin.FuncReturn{
		{
			Type:  tokens.DefinitionBlock,
			Value: ex,
		},
	}, nil
}

package runtime

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// executeFn executes a function
func (e *Executer) executeFn(token *models.Node) (*builtin.FuncReturn, error) {
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
	case "ref":
		return e.runFuncRef(token)
	}

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

		ex := NewExecuter(ExecuterScopeFunction, e.runtime, e).WithName(e.name + ".{" + name + "}")
		for i, arg := range args {
			if arg.VariableType == tokens.ReferenceVariable {
				debug := arg.Debug
				arg, err := e.GetVariableValue(arg.Content)
				if err != nil {
					return nil, errs.WithDebug(err, debug)
				}
				arg.Content = fn.Args[i].Content
				ex.Bind(arg)
				continue
			}

			if arg.VariableType == tokens.ExpressionVariable {
				debug := arg.Debug
				arg, err := e.evaluateExpression(arg)
				if err != nil {
					return nil, errs.WithDebug(err, debug)
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

	if e.scope == ExecuterScopeFunction || e.scope == ExecuterScopeBlock {
		convArgs, err := e.convertArgument(args)
		if err != nil {
			return nil, err
		}
		if e.parent != nil {
			return e.parent.ExecuteFn(name, convArgs)
		}
	}

	return nil, errs.WithDebug(fmt.Errorf("%w: function '%v' not found", errs.RuntimeError, name), token.Debug)
}

// executeComplexFuncCall executes a function with a complex name
func (e *Executer) executeComplexFuncCall(token *models.Node, args []string, fn string, variables []*models.Node) (*builtin.FuncReturn, error) {
	if fn == "construct" {
		return nil, errs.WithDebug(fmt.Errorf("construct is a reserved method"), token.Debug)
	}

	if len(args) == 0 {
		return e.executeFn(token)
	}

	if pkg, err := e.GetPackage(args[0]); err == nil {
		convArgs, err := e.convertArgument(variables)
		if err != nil {
			return nil, errs.WithDebug(err, token.Debug)
		}
		return pkg.Execute(fn, convArgs)
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

		if arg.VariableType == tokens.ExpressionVariable {
			arg, err := e.evaluateExpression(arg)
			if err != nil {
				return nil, err
			}
			convArgs = append(convArgs, &builtin.Variable{
				Type:  arg.VariableType,
				Value: arg.Value,
			})
			continue
		}

		convArgs = append(convArgs, &builtin.Variable{
			Type:  arg.VariableType,
			Value: arg.Value,
		})
	}

	return convArgs, nil
}

// newBlock creates a new block
func (e *Executer) newBlock(block *models.Node, args []*models.Node) (*builtin.FuncReturn, error) {
	ex := NewExecuter(ExecuterScopeDefinition, e.runtime, e).WithName(e.name + ".[" + block.Content + "]")
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

	block.VariableType = tokens.DefinitionReference

	return &builtin.FuncReturn{
		Type:  tokens.DefinitionReference,
		Value: ex,
	}, nil
}

func (e *Executer) GetPackage(name string) (builtin.Package, error) {
	pkgName, ok := e.packages[name]

	if ok {
		pkg, ok := e.runtime.pkgs[pkgName]
		if !ok {
			return nil, fmt.Errorf("package '%v' not found", name)
		}
		return pkg, nil
	}

	if e.scope != ExecuterScopeGlobal && e.parent != nil {
		return e.parent.GetPackage(name)
	}

	return nil, fmt.Errorf("package '%v' not found", name)
}

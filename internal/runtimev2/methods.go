package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/lang"
	"go.uber.org/zap"
)

// callFunctionFromNode calls a function from a node
func (e *Executer) callFunctionFromNode(n *models.Node) (lang.Object, error) {
	name := n.Content

	method, err := e.GetMethod(name)
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}

	var (
		variadicArgName string
		isVariadic      bool
	)

	variadicMethod, ok := method.(lang.VariadicMethod)
	if ok {
		isVariadic = variadicMethod.HasVariadicArg()
		variadicArgName = variadicMethod.GetVariadicArg()
	}

	expectedArgs := len(method.Args())
	givenArgs := len(n.Args)

	zap.L().Debug("calling function", zap.String("name", name), zap.Int("expectedArgs", expectedArgs), zap.Int("givenArgs", givenArgs), zap.Bool("isVariadic", isVariadic))

	if !isVariadic {
		// argument number should match
		if expectedArgs != givenArgs {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s' expects %d arguments, got %d", errs.InvalidArguments, name, expectedArgs, givenArgs), n.Debug)
		}
	} else {
		// if the arguments are less than expected return error
		if expectedArgs > givenArgs {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s' expects %d arguments, got %d", errs.InvalidArguments, name, expectedArgs, givenArgs), n.Debug)
		}
	}

	args, err := e.getFunctionArguments(n.Args, method.Args())
	if err != nil {
		return nil, err
	}

	if givenArgs > expectedArgs && isVariadic {
		var variadicArgList = make([]lang.Object, givenArgs-expectedArgs)

		for i := expectedArgs; i < givenArgs; i++ {
			obj, err := e.getFunctionArgFromNode(n.Args[i])
			if err != nil {
				return nil, err
			}

			variadicArgList[i-expectedArgs] = obj
		}

		variadicArg := lang.NewList(variadicArgName, variadicArgList, n.Debug)
		args = append(args, variadicArg)
	} else if isVariadic {
		args = append(args, lang.NewList(variadicArgName, make([]lang.Object, 0), n.Debug))
	}

	r, err := method.Execute(args)

	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}

	if len(n.Children) > 0 {
		return e.getObjectValueByNodes(r, n.Children)
	}

	return r, nil
}

func (e *Executer) getFunctionArguments(nodeArgs []*models.Node, argNames []string) ([]lang.Object, error) {
	args := make([]lang.Object, 0, len(argNames))

	for i, argName := range argNames {
		obj, err := e.getFunctionArgFromNode(nodeArgs[i])
		if err != nil {
			return nil, err
		}

		obj.Rename(argName)
		args = append(args, obj)
	}

	zap.L().Debug("function arguments", zap.Any("args", args))

	return args, nil
}

func (e *Executer) getFunctionArgFromNode(child *models.Node) (lang.Object, error) {
	if child.Reference {
		obj, err := e.GetVariable(child.Content)
		if err != nil {
			return nil, errs.WithDebug(err, child.Debug)
		}

		return obj.Copy(), nil
	}

	_, obj, err := e.createObjectFromNode(child)

	if err != nil {
		return nil, errs.WithDebug(err, child.Debug)
	}

	return obj.Copy(), nil
}

// createMethodFromNode creates a method from a node
func (e *Executer) createMethodFromNode(n *models.Node) (string, lang.Method, error) {
	name := n.Content

	argsLen := len(n.Args)
	var args []string

	if argsLen > 0 {
		args = make([]string, argsLen)
	}

	for i, arg := range n.Args {
		args[i] = arg.Content
	}

	method := lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
		ex := NewExecuter(ExecuterScopeFunction, e.runtime, e).WithName(e.name + ".{" + name + "}")

		for _, arg := range args {
			ex.BindObject(arg.Name(), arg)
		}

		r, err := ex.Execute(n.Children)
		return r, err
	}).WithArgs(args).WithDebug(n.Debug)

	zap.L().Debug("creating method from node", zap.String("name", name), zap.Any("args", args))

	return name, method, nil
}

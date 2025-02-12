package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/lang"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// createMethodFromNode creates a method from a node
func (e *Executer) createMethodFromNode(n *models.Node) (string, lang.Method, error) {
	name := n.Content

	argsLen := len(n.Children)
	args := make([]string, argsLen, argsLen)

	for i, arg := range n.Args {
		args[i] = arg.Content
	}

	method := lang.NewFunction(args, func(args []lang.Object) (lang.Object, error) {
		ex := NewExecuter(ExecuterScopeFunction, e.runtime, e).WithName(e.name + ".{" + name + "}")

		for _, arg := range args {
			ex.BindObject(arg.Name(), arg)
		}

		return ex.Execute(n.Children)
	}, n.Debug)

	return name, method, nil
}

// createObjectFromNode creates an object from a node
func (e *Executer) createObjectFromNode(n *models.Node) (string, lang.Object, error) {
	name := n.Content
	var obj lang.Object

	switch n.VariableType {
	case tokens.StringVariable:
		s, ok := n.Value.(string)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not string", errs.ValueError), n.Debug)
		}
		obj = lang.NewString(name, s, n.Debug)
	case tokens.IntVariable:
		i, ok := n.Value.(int)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not a number", errs.ValueError), n.Debug)
		}
		obj = lang.NewInteger(name, i, n.Debug)
	case tokens.ListVariable:
		li, err := e.createListFromNode(n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = li
	}

	// if the object is a constant, make it immutable
	if n.Type == tokens.Const {
		obj.Immute()
	}

	return name, obj, nil
}

// callFunctionFromNode calls a function from a node
func (e *Executer) callFunctionFromNode(n *models.Node) (lang.Object, error) {
	name := n.Content

	method, err := e.GetMethod(name)
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}

	expectedArgs := len(method.Args())
	givenArgs := len(n.Args)

	if expectedArgs != givenArgs {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s' expects %d arguments, got %d", errs.InvalidArguments, name, expectedArgs, givenArgs), n.Debug)
	}

	expectedArgNames := method.Args()
	args := make([]lang.Object, 0, givenArgs)

	for i, child := range n.Args {
		if child.Reference {
			obj, err := e.GetVariable(child.Content)
			if err != nil {
				return nil, errs.WithDebug(err, child.Debug)
			}
			obj.Rename(expectedArgNames[i])
			args = append(args, obj)
			continue
		}

		_, obj, err := e.createObjectFromNode(child)
		if err != nil {
			return nil, errs.WithDebug(err, child.Debug)
		}
		obj.Rename(expectedArgNames[i])
		args = append(args, obj)
	}

	r, err := method.Execute(args)
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}
	return r, nil
}

func (e *Executer) assignObjectFromNode(n *models.Node) error {
	name := n.Content

	_, obj, err := e.createObjectFromNode(n)
	if err != nil {
		return errs.WithDebug(err, n.Debug)
	}

	return e.AssignVariable(name, obj)
}

// createListFromNode creates a list from a node
func (e *Executer) createListFromNode(n *models.Node) (lang.Object, error) {
	name := n.Content

	childrens := len(n.Children)
	li := make([]lang.Object, childrens, childrens)

	for i, child := range n.Children {
		_, obj, err := e.createObjectFromNode(child)
		if err != nil {
			return nil, err
		}
		li[i] = obj
	}

	return lang.NewList(name, li, n.Debug), nil
}

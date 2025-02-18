package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

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

	method := lang.NewFunction(args, func(args []lang.Object) (lang.Object, error) {
		ex := NewExecuter(ExecuterScopeFunction, e.runtime, e).WithName(e.name + ".{" + name + "}")

		for _, arg := range args {
			ex.BindObject(arg.Name(), arg)
		}

		r, err := ex.Execute(n.Children)
		return r, err
	}, n.Debug)

	return name, method, nil
}

// createObjectFromNode creates an object from a node
func (e *Executer) createObjectFromNode(n *models.Node) (string, lang.Object, error) {
	name := n.Content
	var obj lang.Object

	switch n.VariableType {
	default:
		return "", nil, fmt.Errorf("unknown variable type: '%s'", n.VariableType)
	case tokens.StringVariable:
		s, ok := n.Value.(string)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not string", errs.ValueError), n.Debug)
		}
		obj = lang.NewString(name, s, n.Debug)
	case tokens.IntVariable:
		i, ok := n.Value.(int)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not an integer", errs.ValueError), n.Debug)
		}
		obj = lang.NewInteger(name, i, n.Debug)
	case tokens.FloatVariable:
		f, ok := n.Value.(float64)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not a float", errs.ValueError), n.Debug)
		}
		obj = lang.NewFloat(name, f, n.Debug)
	case tokens.ListVariable:
		li, err := e.createListFromNode(n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = li
	case tokens.InlineValue:
		typ := e.getVariableTypeFromType(n)
		n.VariableType = typ
		return e.createObjectFromNode(n)
	case tokens.ExpressionVariable:
		expr, err := e.evaluateExpression(n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = expr
	case tokens.ReferenceVariable:
		// ERR: Maybe this should be a reference to a variable, not a reference to a value
		refName := n.Content
		if refVal, ok := n.Value.(string); ok {
			refName = refVal
		}

		ref, err := e.GetVariable(refName)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = ref
	case tokens.BoolVariable:
		b, ok := n.Value.(bool)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not boolean", errs.ValueError), n.Debug)
		}
		obj = lang.NewBool(name, b, n.Debug)
	case tokens.NilVariable:
		obj = lang.NewNil(name, n.Debug)
	}

	if n.ObjectAccessors != nil {
		accessed, err := e.accessObject(obj, n.ObjectAccessors)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = accessed
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
			if obj.Type() != lang.TNil {
				obj = obj.Copy()
				obj.Rename(expectedArgNames[i])
			}
			args = append(args, obj)
			continue
		}

		_, obj, err := e.createObjectFromNode(child)
		if err != nil {
			return nil, errs.WithDebug(err, child.Debug)
		}
		obj = obj.Copy()
		obj.Rename(expectedArgNames[i])
		args = append(args, obj)
	}

	r, err := method.Execute(args)
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}
	return r, nil
}

// assignObjectFromNode assigns an object from a node
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
	li := make([]lang.Object, childrens)

	for i, child := range n.Children {
		_, obj, err := e.createObjectFromNode(child)
		if err != nil {
			return nil, err
		}
		li[i] = obj
	}

	return lang.NewList(name, li, n.Debug), nil
}

// accessObject accesses an object
func (e *Executer) accessObject(obj lang.Object, accessors []*models.Node) (lang.Object, error) {
	if obj.Type() != lang.TList && obj.Type() != lang.TDefinition {
		return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
	}

	var access []any

	for _, accessor := range accessors {
		if accessor.VariableType == tokens.ReferenceVariable {
			obj, err := e.GetVariable(accessor.Content)
			if err != nil {
				return nil, errs.WithDebug(err, accessor.Debug)
			}
			access = append(access, obj.Value())
		} else {
			access = append(access, accessor.Value)
		}
	}

	if obj.Type() == lang.TList {
		li, ok := obj.Value().([]lang.Object)
		if !ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
		}

		var value any = li
		for _, a := range access {
			i, ok := a.(int)
			if !ok {
				s, ok := a.(string)
				if !ok {
					return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
				}

				ob, err := e.GetVariable(s)
				if err != nil || ob.Type() != lang.TInt {
					return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
				}
				i = ob.Value().(int)
			}

			v, ok := value.([]lang.Object)
			if !ok {
				return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
			}

			if i < 0 || i >= len(v) {
				return nil, errs.WithDebug(fmt.Errorf("%w: %v, length: %d", errs.IndexOutOfRange, i, len(v)), accessors[0].Debug)
			}

			value = v[i]
		}

		if v, ok := value.(lang.Object); ok {
			value = v.Value()
		}

		_, ob, err := e.createObjectFromNode(&models.Node{
			VariableType: tokens.InlineValue,
			Type:         e.getTypeFromValue(value),
			Content:      obj.Name(),
			Value:        value,
			Debug:        obj.Debug(),
		})
		if err != nil {
			return nil, errs.WithDebug(err, obj.Debug())
		}

		return ob, nil
	}

	return obj, nil
}

// createDefinitionFromNode creates a definition from a node
func (e *Executer) createObjectFromDefinitionNode(n *models.Node) (string, lang.Object, error) {
	name := n.Content

	ex := NewExecuter(ExecuterScopeDefinition, e.runtime, e).WithName(e.name + ".[" + name + "]")
	_, err := ex.Execute(n.Children)
	if err != nil {
		return "", nil, errs.WithDebug(err, n.Debug)
	}

	return name, lang.NewDefinition(e.name+"."+name, name, n.Debug, ex), nil
}

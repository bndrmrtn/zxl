package runtimev2

import (
	"fmt"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tmpl"
	"github.com/flarelang/flare/internal/tokens"
	"github.com/flarelang/flare/lang"
	"go.uber.org/zap"
)

// createObjectFromNode creates an object from a node
func (e *Executer) createObjectFromNode(n *models.Node) (string, lang.Object, error) {
	name := n.Content
	var obj lang.Object

	zap.L().Debug("creating object from node", zap.String("name", name))

	switch n.VariableType {
	default:
		return "", nil, Error(ErrInvalidValue, n.Debug, fmt.Sprintf("unknown variable type: %s", n.VariableType))
	case tokens.StringVariable:
		s, ok := n.Value.(string)
		if !ok {
			return "", nil, Error(ErrInvalidValue, n.Debug, valueIsNotErr("string"))
		}
		obj = lang.NewString(name, s, n.Debug)
	case tokens.TemplateVariable:
		str, ok := n.Value.(string)
		if !ok {
			return "", nil, Error(ErrInvalidValue, n.Debug, valueIsNotErr("template"))
		}

		template, err := tmpl.NewTemplate(str)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}

		str, err = e.parseTemplate(template)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}

		obj = lang.NewString(name, str, n.Debug)
	case tokens.IntVariable:
		i, ok := n.Value.(int)
		if !ok {
			if f, ok := n.Value.(float64); ok {
				obj = lang.NewFloat(name, f, n.Debug)
				break
			} else {
				return "", nil, Error(ErrInvalidValue, n.Debug, valueIsNotErr("int"))
			}
		}
		obj = lang.NewInteger(name, i, n.Debug)
	case tokens.FloatVariable:
		f, ok := n.Value.(float64)
		if !ok {
			return "", nil, Error(ErrInvalidValue, n.Debug, valueIsNotErr("float"))
		}
		obj = lang.NewFloat(name, f, n.Debug)
	case tokens.ListVariable:
		li, err := e.createListFromNode(n)
		if err != nil {
			return "", nil, Error(err, n.Debug)
		}
		obj = li
	case tokens.InlineValue:
		typ := e.getVariableTypeFromType(n)
		n.VariableType = typ
		return e.createObjectFromNode(n)
	case tokens.ExpressionVariable:
		expr, err := e.evaluateExpression(n)
		if err != nil {
			return "", nil, Error(err, n.Debug)
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
			return "", nil, Error(err, n.Debug)
		}
		obj = ref
	case tokens.BoolVariable:
		b, ok := n.Value.(bool)
		if !ok {
			return "", nil, Error(ErrInvalidValue, n.Debug, valueIsNotErr("bool"))
		}
		obj = lang.NewBool(name, b, n.Debug)
	case tokens.NilVariable:
		obj = lang.NewNil(name, n.Debug)
	case tokens.ArrayVariable:
		arr, err := e.createObjectFromArrayNode(name, n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = arr
	case tokens.FunctionVariable:
		name, fn, err := e.createMethodFromNode(n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}

		if name != "fn" {
			return "", nil, Error(ErrNamedInlineFunction, n.Debug)
		}

		obj = lang.NewFn("<inlineFn>", n.Debug, fn)
	case tokens.FunctionCallVariable:
		var err error
		obj, err = e.callFunctionFromNode(n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
	}

	if obj == nil {
		return "", nil, Error(ErrInvalidObject, n.Debug)
	}

	zap.L().Debug("object created", zap.String("name", name), zap.String("type", obj.Type().String()))

	if len(n.ObjectAccessors) > 0 {
		accessed, err := e.accessObject(obj, n.ObjectAccessors)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = accessed
	}

	// if the object is a constant, make it immutable
	if obj.Type() != lang.TNil && n.Type == tokens.Const {
		zap.L().Debug("making object immutable", zap.String("name", name))
		obj.Immute()
	}

	return name, obj, nil
}

// assignObjectFromNode assigns an object from a node
func (e *Executer) assignObjectFromNode(n *models.Node) error {
	name := n.Content

	_, obj, err := e.createObjectFromNode(n)
	if err != nil {
		return errs.WithDebug(err, n.Debug)
	}

	zap.L().Debug("assigning object from node", zap.String("name", name), zap.Any("object", obj))

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

	zap.L().Debug("creating list from node", zap.String("name", name), zap.Any("list", li))

	liObj := lang.NewList(name, li, n.Debug)

	if len(n.ObjectAccessors) > 0 {
		return e.accessObject(liObj, n.ObjectAccessors)
	}

	return liObj, nil
}

// createDefinitionFromNode creates a definition from a node
func (e *Executer) createObjectFromDefinitionNode(n *models.Node) (string, lang.Object, error) {
	name := n.Content
	ex := NewExecuter(ExecuterScopeDefinition, e.runtime, e).WithName(e.name + ".[" + name + "]")

	zap.L().Debug("creating object from definition node", zap.String("name", name))

	return name, lang.NewDefinition(e.name+"."+name, name, n.Debug, n.Children, ex), nil
}

func (e *Executer) getObjectValueByNodes(obj lang.Object, nodes []*models.Node) (lang.Object, error) {
	zap.L().Debug("getting object value by nodes", zap.Any("nodes", nodes))

	for _, node := range nodes {
		if node.Type != tokens.FuncCall && node.Type != tokens.Identifier {
			return nil, Error(ErrInvalidObjectAccess, node.Debug, obj.Type())
		}

		if node.Type == tokens.FuncCall {
			m := obj.Method(node.Content)
			if m == nil {
				return nil, Error(ErrInvalidObjectAccess, node.Debug, fnErr(node.Content))
			}

			if len(node.Args) != len(m.Args()) {
				return nil, Error(ErrInvalidArguments, node.Debug, expectedErr(len(m.Args()), node.Args))
			}

			args := make([]lang.Object, len(node.Args))
			for i, arg := range node.Args {
				_, o, err := e.createObjectFromNode(arg)
				if err != nil {
					return nil, errs.WithDebug(err, arg.Debug)
				}
				args[i] = o
			}

			r, err := m.Execute(args)
			if err != nil {
				return nil, errs.WithDebug(err, node.Debug)
			}

			obj = r
		}

		if node.Type == tokens.Identifier {
			m := obj.Variable(node.Content)
			if m == nil {
				return nil, Error(ErrInvalidObjectAccess, node.Debug, node.Content)
			}

			obj = m
		}
	}

	return obj, nil
}

func (e *Executer) parseTemplate(template []tmpl.Part) (string, error) {
	var result string
	for _, part := range template {
		if part.Static {
			result += part.Content
		} else {
			_, ob, err := e.createObjectFromNode(part.Node)
			if err != nil {
				return "", err
			}
			result += ob.String()
		}
	}

	zap.L().Debug("parsed template", zap.String("result", result))

	return result, nil
}

// createObjectFromArrayNode creates an array object from a node.
func (e *Executer) createObjectFromArrayNode(name string, n *models.Node) (lang.Object, error) {
	if len(n.Children) == 0 {
		return lang.NewArray(name, n.Debug, nil, nil), nil
	}

	var keys []lang.Object
	var values []lang.Object

	for _, child := range n.Children {
		key := child.Args[0]

		if key.Type == tokens.Identifier {
			keys = append(keys, lang.NewString(key.Content, key.Content, key.Debug))
		} else {
			_, val, err := e.createObjectFromNode(key)
			if err != nil {
				return nil, errs.WithDebug(err, key.Debug)
			}
			keys = append(keys, val)
		}

		value := child.Children[0]
		_, val, err := e.createObjectFromNode(value)
		if err != nil {
			return nil, errs.WithDebug(err, value.Debug)
		}
		values = append(values, val)
	}

	zap.L().Debug("creating array from node", zap.String("name", name), zap.Any("keys", keys), zap.Any("values", values))

	return lang.NewArray(name, n.Debug, keys, values), nil
}

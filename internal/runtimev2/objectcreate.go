package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tmpl"
	"github.com/bndrmrtn/zxl/internal/tokens"
	"github.com/bndrmrtn/zxl/lang"
)

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
	case tokens.TemplateVariable:
		str, ok := n.Value.(string)
		if !ok {
			return "", nil, errs.WithDebug(fmt.Errorf("%w: value is not a template", errs.ValueError), n.Debug)
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
	case tokens.ArrayVariable:
		arr, err := e.createObjectFromArrayNode(name, n)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = arr
	}

	if len(n.ObjectAccessors) > 0 {
		accessed, err := e.accessObject(obj, n.ObjectAccessors)
		if err != nil {
			return "", nil, errs.WithDebug(err, n.Debug)
		}
		obj = accessed
	}

	// if the object is a constant, make it immutable
	if obj.Type() != lang.TNil && n.Type == tokens.Const {
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

// createDefinitionFromNode creates a definition from a node
func (e *Executer) createObjectFromDefinitionNode(n *models.Node) (string, lang.Object, error) {
	name := n.Content
	ex := NewExecuter(ExecuterScopeDefinition, e.runtime, e).WithName(e.name + ".[" + name + "]")

	return name, lang.NewDefinition(e.name+"."+name, name, n.Debug, n.Children, ex), nil
}

func (e *Executer) getObjectValueByNodes(obj lang.Object, nodes []*models.Node) (lang.Object, error) {
	for _, node := range nodes {
		if node.Type != tokens.FuncCall && node.Type != tokens.Identifier {
			return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), node.Debug)
		}

		if node.Type == tokens.FuncCall {
			m := obj.Method(node.Content)
			if m == nil {
				return nil, errs.WithDebug(fmt.Errorf("%w: method '%s()' not found", errs.RuntimeError, node.Content), node.Debug)
			}

			if len(node.Args) != len(m.Args()) {
				return nil, errs.WithDebug(fmt.Errorf("%w: '%s' expects %d arguments, got %d", errs.InvalidArguments, node.Content, len(m.Args()), len(node.Args)), node.Debug)
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
				return nil, errs.WithDebug(fmt.Errorf("%w: variable '%s' not found", errs.RuntimeError, node.Content), node.Debug)
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

	return lang.NewArray(name, n.Debug, keys, values), nil
}

package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
	"github.com/bndrmrtn/zxl/lang"
)

// executeNode executes a node
func (e *Executer) executeNode(node *models.Node) (lang.Object, error) {
	switch node.Type {
	default:
		return nil, errs.WithDebug(fmt.Errorf("unhandled node type: %s", node.Type), node.Debug)
	case tokens.Use:
		using := node.Content
		as := node.Value.(string)
		if _, ok := e.usedNamespaces[as]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s' as '%s'", errs.CannotReUseNamespace, using, as), node.Debug)
		}
		e.usedNamespaces[as] = using
	case tokens.Function:
		name, method, err := e.createMethodFromNode(node)
		if err != nil {
			return nil, err
		}
		if _, ok := e.functions[name]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s(...)'", errs.CannotRedecareFunction, name), node.Debug)
		}

		e.mu.Lock()
		e.functions[name] = method
		e.mu.Unlock()
	case tokens.Let, tokens.Const:
		name, object, err := e.createObjectFromNode(node)
		if err != nil {
			return nil, err
		}
		if _, ok := e.objects[name]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotRedeclareVariable, name), node.Debug)
		}
		e.mu.Lock()
		e.objects[name] = object
		e.mu.Unlock()
	case tokens.FuncCall:
		_, err := e.callFunctionFromNode(node)
		if err != nil {
			return nil, errs.WithDebug(err, node.Debug)
		}
	case tokens.Assign:
		err := e.assignObjectFromNode(node)
		if err != nil {
			return nil, errs.WithDebug(err, node.Debug)
		}
	case tokens.Increment, tokens.Decrement:
		v, err := e.GetVariable(node.Content)
		if err != nil {
			return nil, errs.WithDebug(err, node.Debug)
		}

		if v.Type() != lang.TInt {
			return nil, errs.WithDebug(fmt.Errorf("%w: cannot increment type: '%s'", errs.RuntimeError, v.Type()), node.Debug)
		}

		add := 1
		if node.Type == tokens.Decrement {
			add = -1
		}

		e.AssignVariable(node.Content, lang.NewInteger(node.Content, v.Value().(int)+add, node.Debug))
	case tokens.Define:
		name, object, err := e.createObjectFromDefinitionNode(node)
		if err != nil {
			return nil, err
		}
		if _, ok := e.objects[name]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotRedeclareVariable, name), node.Debug)
		}
		e.mu.Lock()
		e.objects[name] = object
		e.mu.Unlock()
	case tokens.Return, tokens.EmptyReturn:
		return e.handleReturn(node)
	case tokens.If:
		return e.handleIf(node)
	case tokens.While:
		return e.handleWhile(node)
	case tokens.For:
		return e.handleFor(node)
	case tokens.Spin:
		return e.handleSpin(node)
	}

	return nil, nil
}

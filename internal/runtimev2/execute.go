package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/lang"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// executeNode executes a node
func (e *Executer) executeNode(node *models.Node) (lang.Object, error) {
	switch node.Type {
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
		return nil, errs.WithDebug(err, node.Debug)
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
	case tokens.EmptyReturn:
		return lang.NilObject, nil
	case tokens.Return:
		return e.handleReturn(node)
	}
	return nil, nil
}

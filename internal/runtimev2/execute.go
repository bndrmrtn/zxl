package runtimev2

import (
	"errors"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
	"github.com/flarelang/flare/lang"
)

// executeNode executes a node
func (e *Executer) executeNode(node *models.Node) (lang.Object, error) {
	switch node.Type {
	default:
		return nil, Error(ErrUnhandledNodeType, node.Debug, node.Type)
	case tokens.Use:
		using := node.Content
		as := node.Value.(string)
		if _, ok := e.usedNamespaces[as]; ok {
			return nil, Error(ErrNamespaceInUse, node.Debug, nsErr(using, as))
		}
		e.usedNamespaces[as] = using
	case tokens.Function:
		name, method, err := e.createMethodFromNode(node)
		if err != nil {
			return nil, err
		}
		if _, ok := e.functions[name]; ok {
			return nil, Error(ErrFunctionRedeclared, node.Debug, fnErr(name))
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
			return nil, Error(ErrVariableNotDeclared, node.Debug, name)
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
			return nil, Error(ErrInvalidIncrementTarget, node.Debug, v.Type())
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
			return nil, Error(ErrVariableRedeclared, node.Debug, name)
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
	case tokens.Error:
		if len(node.Children) == 0 {
			return nil, Error(ErrEmptyErrorBlock, node.Debug)
		}

		_, err := e.Execute(node.Children)
		if err != nil {
			var de errs.DebugError
			if errors.As(err, &de) {
				err = de.GetParentError()
			}
			e.BindObject(node.Content, lang.NewString(node.Content, err.Error(), node.Debug))
		} else {
			e.BindObject(node.Content, lang.NilObject)
		}
	}

	return nil, nil
}

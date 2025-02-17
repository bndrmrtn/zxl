package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

// handleReturn handles return tokens
func (e *Executer) handleReturn(node *models.Node) (lang.Object, error) {
	if node.Type == tokens.EmptyReturn {
		if e.scope == ExecuterScopeBlock && e.parent != nil {
			return e.parent.handleReturn(node)
		}

		return lang.NilObject, nil
	}

	// Evaluate return value
	value, err := e.evaluateExpression(node)
	if err != nil {
		return nil, err
	}

	if e.scope == ExecuterScopeBlock && e.parent != nil {
		return e.parent.handleReturn(node)
	}

	if value == nil {
		return nil, nil
	}

	return value, nil
}

// handleIf handles if tokens
func (e *Executer) handleIf(node *models.Node) (lang.Object, error) {
	// Evaluate condition
	condition, err := e.evaluateExpression(&models.Node{
		Type:         tokens.If,
		VariableType: tokens.ExpressionVariable,
		Children:     node.Args,
		Debug:        node.Debug,
	})
	if err != nil {
		return nil, errs.WithDebug(err, node.Debug)
	}

	if condition.Type() != lang.TBool {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected boolean", errs.ValueError), node.Debug)
	}

	ok := condition.Value().(bool)

	if ok {
		if len(node.Children) == 0 {
			return nil, nil
		}

		child := node.Children[0]
		if child.Type == tokens.Then {
			ex := NewExecuter(ExecuterScopeBlock, e.runtime, e).WithName(e.name)
			return ex.Execute(child.Children)
		}
	} else {
		if len(node.Children) < 2 {
			return nil, nil
		}

		child := node.Children[1]
		if child.Type == tokens.Else {
			ex := NewExecuter(ExecuterScopeBlock, e.runtime, e).WithName(e.name)
			return ex.Execute(child.Children)
		}
	}

	return nil, nil
}

// handleWhile handles while tokens
func (e *Executer) handleWhile(node *models.Node) (lang.Object, error) {
	ex := NewExecuter(ExecuterScopeBlock, e.runtime, e).WithName(e.name)

	for {

		// Evaluate condition
		condition, err := e.evaluateExpression(&models.Node{
			Type:         tokens.While,
			VariableType: tokens.ExpressionVariable,
			Children:     node.Args,
			Debug:        node.Debug,
		})
		if err != nil {
			return nil, err
		}

		if condition.Type() != lang.TBool {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected boolean", errs.ValueError), node.Debug)
		}

		ok := condition.Value().(bool)
		if !ok {
			break
		}

		ret, err := ex.Execute(node.Children)
		if ret != nil || err != nil {
			return ret, err
		}
	}

	return nil, nil
}

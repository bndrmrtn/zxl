package runtimev2

import (
	"fmt"
	"strings"

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

	if value != nil {
		return value, nil
	}

	if e.scope == ExecuterScopeBlock && e.parent != nil {
		return e.parent.handleReturn(node)
	}

	return nil, nil
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

func (e *Executer) handleFor(node *models.Node) (lang.Object, error) {
	ex := NewExecuter(ExecuterScopeBlock, e.runtime, e).WithName(e.name)

	if len(node.Args) != 2 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected 1 identifier and 1 iterable expression", errs.ValueError), node.Debug)
	}

	// stage 1: setting up the iterator and iterable

	// make sure the variable is a let, not a reference
	node.Args[0].Type = tokens.Let
	node.Args[0].VariableType = tokens.NilVariable
	node.Args[0].Reference = false
	name, _, err := e.createObjectFromNode(node.Args[0])
	if err != nil {
		return nil, err
	}

	_, iterable, err := e.createObjectFromNode(&models.Node{
		VariableType: tokens.ExpressionVariable,
		Children: []*models.Node{
			node.Args[1],
		},
	})
	if err != nil {
		return nil, err
	}

	switch iterable.Type() {
	default:
		return nil, errs.WithDebug(fmt.Errorf("%w: expected iterable value or expression, got '%s'", errs.ValueError, iterable.Type()), node.Debug)
	case lang.TList:
		for _, item := range iterable.Value().([]lang.Object) {
			exec := NewExecuter(ExecuterScopeBlock, ex.runtime, ex)
			exec.mu.Lock()
			exec.objects[name] = item
			exec.mu.Unlock()

			ret, err := exec.Execute(node.Children)
			if ret != nil || err != nil {
				return ret, err
			}
		}
	case lang.TString:
		str := strings.Split(iterable.Value().(string), "")
		for _, item := range str {
			exec := NewExecuter(ExecuterScopeBlock, ex.runtime, ex)
			exec.mu.Lock()
			exec.objects[name] = lang.NewString(name, string(item), node.Debug)
			exec.mu.Unlock()

			ret, err := exec.Execute(node.Children)
			if ret != nil || err != nil {
				return ret, err
			}
		}
	case lang.TInt:
		for item := 0; item < iterable.Value().(int); item++ {
			exec := NewExecuter(ExecuterScopeBlock, ex.runtime, ex)
			exec.mu.Lock()
			exec.objects[name] = lang.NewInteger(name, item, node.Debug)
			exec.mu.Unlock()

			ret, err := exec.Execute(node.Children)
			if ret != nil || err != nil {
				return ret, err
			}
		}
	}

	return nil, nil
}

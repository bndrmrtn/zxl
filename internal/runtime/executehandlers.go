package runtime

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (e *Executer) handleLetConst(token *models.Node) error {
	// Check if variable is already declared
	if _, ok := e.vars[token.Content]; ok {
		return errs.WithDebug(fmt.Errorf("%w: %v", errs.CannotRedeclareVariable, token.Content), token.Debug)
	}

	variableName := token.Content

	// Check if variable is an expression
	if token.VariableType == tokens.ExpressionVariable {
		v, err := e.evaluateExpression(token)
		if err != nil {
			return err
		}

		// Required for the variable to be accessible in the future without re-evaluating
		token.Value = v.Value
		token.Type = v.Type
		token.VariableType = v.VariableType
		e.vars[token.Content] = v
		return nil
	}

	e.vars[variableName] = token
	return nil
}

func (e *Executer) handleAssignment(token *models.Node) error {
	// Get variable name
	variableName := token.Content
	ex := e

	if strings.Contains(variableName, ".") {
		parts := strings.Split(variableName, ".")
		last := parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		exec, _, err := e.accessUnderlyingVariable(parts)
		if err != nil {
			return errs.WithDebug(err, token.Debug)
		}
		ex = exec
		variableName = last

		if token.VariableType == tokens.ReferenceVariable {
			ref, err := e.GetVariableValue(token.Value.(string))
			if err != nil {
				return errs.WithDebug(err, token.Debug)
			}
			token = ref
		}
	}

	// Check if variable is declared
	v, ok := ex.vars[variableName]
	if !ok {
		if ex.scope == ExecuterScopeBlock {
			return ex.parent.handleAssignment(token)
		}
		return errs.WithDebug(fmt.Errorf("%w: %v", errs.VariableNotDeclared, token.Content), token.Debug)
	}
	// Check if variable is a constant
	if v.Type == tokens.Const {
		return errs.WithDebug(fmt.Errorf("%w: %v", errs.CannotReassignConstant, token.Content), token.Debug)
	}

	if token.VariableType == tokens.ExpressionVariable {
		value, err := e.evaluateExpression(token)
		if err != nil {
			return err
		}
		token.Value = value.Value
		token.Type = value.Type
		token.VariableType = value.VariableType
	}

	if token.VariableType == tokens.ReferenceVariable && e.scope == ExecuterScopeBlock {
		ref, err := e.GetVariableValue(token.Value.(string))
		if err != nil {
			return err
		}
		token = ref
	}

	ex.vars[variableName] = token
	return nil
}

func (e *Executer) handleReturn(token *models.Node) ([]*builtin.FuncReturn, error) {
	if token.Type == tokens.EmptyReturn {
		if e.scope == ExecuterScopeBlock && e.parent != nil {
			return e.parent.handleReturn(token)
		}

		return []*builtin.FuncReturn{
			{
				Type:  tokens.EmptyReturnValue,
				Value: nil,
			},
		}, nil
	}

	// Evaluate return value
	value, err := e.evaluateExpression(token)
	if err != nil {
		return nil, err
	}

	if e.scope == ExecuterScopeBlock && e.parent != nil {
		return e.parent.handleReturn(token)
	}

	if value == nil {
		return nil, nil
	}

	return []*builtin.FuncReturn{
		{
			Type:  value.VariableType,
			Value: value.Value,
		},
	}, nil
}

func (e *Executer) handleIf(token *models.Node) ([]*builtin.FuncReturn, error) {
	// Evaluate condition
	condition, err := e.evaluateExpression(&models.Node{
		Type:         tokens.If,
		VariableType: tokens.ExpressionVariable,
		Children:     token.Args,
	})
	if err != nil {
		return nil, err
	}

	// Check if condition is a boolean
	if condition.VariableType != tokens.BoolVariable {
		return nil, errs.WithDebug(fmt.Errorf("expected bool value, but got '%v'", condition.Type), token.Debug)
	}

	ok := condition.Value.(bool)

	// Evaluate block
	if ok {
		if len(token.Children) == 0 {
			return nil, nil
		}

		child := token.Children[0]
		if child.Type == tokens.Then {
			ex := NewExecuter(ExecuterScopeBlock, e.runtime, e)
			return ex.Execute(child.Children)
		}
	} else {
		if len(token.Children) < 2 {
			return nil, nil
		}

		child := token.Children[1]
		if child.Type == tokens.Else {
			ex := NewExecuter(ExecuterScopeBlock, e.runtime, e)
			return ex.Execute(child.Children)
		}
	}

	return nil, nil
}

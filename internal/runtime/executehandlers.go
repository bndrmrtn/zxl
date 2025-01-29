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
	}

	e.vars[token.Content] = token
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
	}

	// Check if variable is declared
	v, ok := ex.vars[variableName]
	if !ok {
		return errs.WithDebug(fmt.Errorf("%w: %v", errs.VariableNotDeclared, token.Content), token.Debug)
	}
	// Check if variable is a constant
	if v.Type == tokens.Const {
		return errs.WithDebug(fmt.Errorf("%w: %v", errs.CannotReassignConstant, token.Content), token.Debug)
	}

	ex.vars[token.Content] = token
	return nil
}

func (e *Executer) handleReturn(token *models.Node) ([]*builtin.FuncReturn, error) {
	// Evaluate return value
	value, err := e.evaluateExpression(token)
	if err != nil {
		return nil, err
	}

	return []*builtin.FuncReturn{
		{
			Type:  value.VariableType,
			Value: value.Value,
		},
	}, nil
}

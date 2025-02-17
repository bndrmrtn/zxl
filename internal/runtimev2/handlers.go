package runtimev2

import (
	"github.com/bndrmrtn/zexlang/internal/lang"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
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

package ast

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (b *Builder) getType(t *models.Token) (tokens.VariableType, error) {
	switch t.Type {
	case tokens.Number:
		if t.Map["isFloat"] == true {
			return tokens.FloatVariable, nil
		}
		return tokens.IntVariable, nil
	case tokens.String:
		return tokens.StringVariable, nil
	case tokens.Bool:
		return tokens.BoolVariable, nil
	case tokens.Identifier:
		if t.Value == "nil" {
			return tokens.NilVariable, nil
		}
		return tokens.ReferenceVariable, nil
	default:
		return 0, errs.WithDebug(fmt.Errorf("%w: unknown type", errs.SyntaxError), t.Debug)
	}
}

package ast

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

var SemiColonToken = &models.Token{
	Type:  tokens.Semicolon,
	Value: ";",
}

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
	case tokens.Identifier, tokens.Nil:
		if t.Value == "nil" {
			return tokens.NilVariable, nil
		}
		return tokens.ReferenceVariable, nil
	default:
		return 0, errs.WithDebug(fmt.Errorf("%w: unknown type", errs.SyntaxError), t.Debug)
	}
}

func (b *Builder) getValue(t *models.Token) any {
	switch t.Type {
	case tokens.String:
		return b.handleEscapedString(t.Value)
	case tokens.Number:
		isFloat := t.Map["isFloat"]

		if isFloat.(bool) {
			val, _ := strconv.ParseFloat(t.Value, 64)
			return val
		}

		val, _ := strconv.Atoi(t.Value)
		return val
	case tokens.Bool:
		val, _ := strconv.ParseBool(t.Value)
		return val
	case tokens.Nil:
		return nil
	case tokens.Identifier:
		return t.Value
	}
	return t.Value
}

// handleEscapedString replaces escaped characters in a string
func (b *Builder) handleEscapedString(s string) string {
	if len(s) < 2 {
		return s
	}

	quote := s[0]
	s = s[1 : len(s)-1]

	// Replace escaped characters
	s = strings.ReplaceAll(s, `\\`, `\`)
	s = strings.ReplaceAll(s, `\n`, "\n")
	s = strings.ReplaceAll(s, `\t`, "\t")

	// Replace escaped quotes
	if quote == '"' {
		s = strings.ReplaceAll(s, `\"`, `"`)
	} else if quote == '\'' {
		s = strings.ReplaceAll(s, `\'`, `'`)
	}

	return s
}

func (b *Builder) isExpression(n *models.Token) bool {
	return n.Type == tokens.Addition ||
		n.Type == tokens.Subtraction ||
		n.Type == tokens.Multiplication ||
		n.Type == tokens.Division ||
		n.Type == tokens.Less ||
		n.Type == tokens.LessOrEqual ||
		n.Type == tokens.Greater ||
		n.Type == tokens.GreaterOrEqual ||
		n.Type == tokens.Equation ||
		n.Type == tokens.NotEquation ||
		n.Type == tokens.And ||
		n.Type == tokens.Or ||
		n.Type == tokens.Not ||
		n.Type == tokens.Power
}

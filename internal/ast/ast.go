package ast

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// Builder is the AST builder
type Builder struct{}

// NewBuilder creates a new AST builder
func NewBuilder() *Builder {
	return &Builder{}
}

// Build builds the AST from the tokens
func (b *Builder) Build(tokens []*models.Token) ([]*models.Node, error) {
	var (
		inx   int
		nodes []*models.Node
	)

	for inx < len(tokens) {
		node, err := b.buildNode(tokens, &inx)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	return nodes, nil
}

func (b *Builder) buildNode(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]

	switch token.Type {
	case tokens.String, tokens.Number, tokens.Bool, tokens.Nil:
		return b.parseInlineValue(ts, inx)
	case tokens.Let, tokens.Const:
		return b.parseLetConst(ts, inx)
	case tokens.Define:
		return b.parseDefine(ts, inx)
	case tokens.Function:
		return b.parseFunction(ts, inx)
	case tokens.Identifier, tokens.This:
		return b.parseIdentifier(ts, inx)
	default:
		return nil, errs.WithDebug(fmt.Errorf("%w: invalid token '%v'", errs.SyntaxError, token.Value), token.Debug)
	}
}

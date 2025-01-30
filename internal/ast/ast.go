package ast

import (
	"context"
	"fmt"
	"time"

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
func (b *Builder) Build(ts []*models.Token) ([]*models.Node, error) {
	var (
		inx   int
		nodes []*models.Node
	)

	if len(ts) == 0 {
		return nil, fmt.Errorf("%w: no tokens to build AST", errs.SyntaxError)
	}

	if ts[0].Type == tokens.Namespace {
		node, err := b.parseNamespace(ts, &inx)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, node)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for inx < len(ts) {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("%w: timeout while building AST", errs.RuntimeError)
		default:
			node, err := b.buildNode(ts, &inx)
			if err != nil {
				return nil, err
			}
			if node == nil {
				continue
			}
			nodes = append(nodes, node)
		}
	}

	return nodes, nil
}

func (b *Builder) buildNode(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]

	switch token.Type {
	case tokens.Namespace:
		return nil, errs.WithDebug(fmt.Errorf("%w: namespace can only be at the beginning of the file", errs.SyntaxError), token.Debug)
	case tokens.Addition, tokens.Subtraction, tokens.Multiplication, tokens.Division, tokens.Equation, tokens.NotEquation, tokens.Greater, tokens.GreaterOrEqual, tokens.Less, tokens.LessOrEqual, tokens.And, tokens.Or, tokens.Not:
		*inx++
		return &models.Node{
			Type:    token.Type,
			Content: token.Value,
			Debug:   token.Debug,
		}, nil
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
	case tokens.Return:
		return b.parseReturn(ts, inx)
	case tokens.If:
		return b.parseIf(ts, inx)
	case tokens.Semicolon:
		*inx++
		return nil, nil
	default:
		return nil, errs.WithDebug(fmt.Errorf("%w: invalid token '%v'", errs.SyntaxError, token.Value), token.Debug)
	}
}

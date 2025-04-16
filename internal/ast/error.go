package ast

import (
	"fmt"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
)

func (b *Builder) parseError(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++
	node := &models.Node{
		Type: tokens.Error,
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got EOF", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected : or {, but got EOF", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Colon && ts[*inx].Type != tokens.LeftBrace {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected : or {, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	var (
		isBrace    = ts[*inx-1].Type == tokens.LeftBrace
		braceCount = 1
	)

	var child []*models.Token
	if isBrace {
		for {
			if *inx >= len(ts) {
				break
			}

			if ts[*inx].Type == tokens.LeftBrace {
				braceCount++
			}

			if ts[*inx].Type == tokens.RightBrace {
				braceCount--
			}

			if braceCount == 0 {
				*inx++
				break
			}

			child = append(child, ts[*inx])
			*inx++
		}
	} else {
		for {
			if *inx >= len(ts) {
				break
			}

			if ts[*inx].Type == tokens.Semicolon {
				break
			}

			child = append(child, ts[*inx])
			*inx++
		}
	}

	if *inx < len(ts) && ts[*inx].Type == tokens.Semicolon {
		*inx++
	}

	children, err := b.Build(append(child, SemiColonToken))
	if err != nil {
		return nil, err
	}

	node.Children = children
	return node, nil
}

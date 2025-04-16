package ast

import (
	"fmt"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
)

func (b *Builder) parseObjectAccess(ts []*models.Token, inx *int, node *models.Node) (*models.Node, error) {
	token := ts[*inx]

	if token.Type != tokens.LeftBracket {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '[', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	var (
		children     []*models.Token
		bracketCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ']', but got 'EOF'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.RightBracket {
			bracketCount--
			if bracketCount == 0 {
				*inx++
				break
			}
		}

		if ts[*inx].Type == tokens.LeftBracket {
			bracketCount++
		}

		children = append(children, ts[*inx])
		*inx++
	}

	if bracketCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ']', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	children = append(children, SemiColonToken)

	child, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	if len(child) != 1 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected one child, but got %d", errs.SyntaxError, len(child)), token.Debug)
	}

	node.ObjectAccessors = append(node.ObjectAccessors, child...)

	if *inx+1 < len(ts) {
		if ts[*inx].Type == tokens.LeftBracket {
			return b.parseObjectAccess(ts, inx, node)
		}

		if ts[*inx].Type == tokens.Semicolon {
			*inx++
		}
	}

	return node, nil
}

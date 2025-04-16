package ast

import (
	"fmt"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
)

func (b *Builder) parseNamespace(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:  token.Type,
		Debug: token.Debug,
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	*inx++

	if *inx >= len(ts) {
		return node, nil
	}

	if ts[*inx].Type == tokens.Semicolon {
		*inx++
	}

	return node, nil
}

func (b *Builder) parseUse(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:  token.Type,
		Debug: token.Debug,
	}

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	node.Value = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return node, nil
	}

	if ts[*inx].Type == tokens.Colon {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type != tokens.Identifier {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
		}

		node.Content += ":" + ts[*inx].Value
		node.Value = ts[*inx].Value
		*inx++
	}

	if ts[*inx].Type != tokens.Semicolon && ts[*inx].Type != tokens.As {
		return node, nil
	}

	if ts[*inx].Type == tokens.Semicolon {
		*inx++
		return node, nil
	}

	if ts[*inx].Type != tokens.As {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected 'as', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Value = ts[*inx].Value

	*inx++

	if *inx >= len(ts) {
		return node, nil
	}

	if ts[*inx].Type == tokens.Semicolon {
		*inx++
		return node, nil
	}

	return node, nil
}

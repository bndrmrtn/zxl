package ast

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (b *Builder) parseLetConst(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type: token.Type,
	}
	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), token.Debug)
	}
	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type == tokens.Semicolon {
		node.VariableType = tokens.NilVariable
		node.Value = "nil"
		*inx++
		return node, nil
	}

	if ts[*inx].Type != tokens.Assign {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator", errs.SyntaxError), ts[*inx].Debug)
	}

	var values []*models.Token
	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.Semicolon {
			*inx++
			break
		}

		values = append(values, ts[*inx])
	}

	if len(values) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
	}

	if len(values) == 1 {
		node.Value = values[0].Value
		typ, err := b.getType(values[0])
		if err != nil {
			return nil, err
		}
		node.VariableType = typ
	}

	return node, nil
}

func (b *Builder) parseFunction(ts []*models.Token, inx *int) (m *models.Node, e error) {
	token := ts[*inx]
	node := &models.Node{
		Type: token.Type,
	}
	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), token.Debug)
	}
	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '('", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftParenthesis {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '('", errs.SyntaxError), ts[*inx].Debug)
	}

	*inx++

	var (
		args       []*models.Node
		parenCount = 1
	)
	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.RightParenthesis {
			parenCount--
			if parenCount == 0 {
				*inx++
				break
			}
		}

		if ts[*inx].Type == tokens.LeftParenthesis {
			parenCount++
		}

		arg, err := b.buildNode(ts, inx)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	node.Args = args

	if parenCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ')'", errs.SyntaxError), token.Debug)
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftBrace {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{'", errs.SyntaxError), ts[*inx].Debug)
	}

	var (
		children   []*models.Token
		braceCount = 1
	)

	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.RightBrace {
			braceCount--
			if braceCount == 0 {
				*inx++
				break
			}
		}

		if ts[*inx].Type == tokens.LeftBrace {
			braceCount += 1
		}

		children = append(children, ts[*inx])
	}

	if braceCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}'", errs.SyntaxError), token.Debug)
	}

	child, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	node.Children = child
	return node, nil
}

func (b *Builder) parseIdentifier(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	node := &models.Node{
		Type:    token.Type,
		Content: token.Value,
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: unexpected end of file", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type == tokens.Semicolon || (ts[*inx].Type != tokens.Assign && ts[*inx].Type != tokens.Dot && ts[*inx].Type != tokens.LeftParenthesis) {
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.ReferenceVariable,
			Content:      token.Value,
		}, nil
	}

	if ts[*inx].Type == tokens.LeftParenthesis {
		return b.parseFunctionCall(ts, inx, node)
	}

	if ts[*inx].Type == tokens.Dot {
		node.Content += "."
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type != tokens.Identifier {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), ts[*inx].Debug)
		}

		node.Content += ts[*inx].Value
		node.Reference = true

		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type != tokens.Assign {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator", errs.SyntaxError), ts[*inx].Debug)
		}

		*inx--
	}

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
	}

	var values []*models.Token
	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.Semicolon {
			*inx++
			break
		}

		values = append(values, ts[*inx])
	}

	if len(values) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
	}

	if len(values) == 1 {
		node.Value = values[0].Value
		typ, err := b.getType(values[0])
		if err != nil {
			return nil, err
		}
		node.VariableType = typ
	}

	return node, nil
}

func (b *Builder) parseDefine(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type: token.Type,
	}
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier", errs.SyntaxError), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftBrace {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{'", errs.SyntaxError), ts[*inx].Debug)
	}

	var (
		children   []*models.Token
		braceCount = 1
	)

	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.RightBrace {
			braceCount--
			if braceCount == 0 {
				*inx++
				break
			}
		}

		if ts[*inx].Type == tokens.LeftBrace {
			braceCount++
		}

		if braceCount > 0 {
			children = append(children, ts[*inx])
		}
	}

	if braceCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}'", errs.SyntaxError), token.Debug)
	}

	child, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	node.Children = child
	return node, nil
}

func (b *Builder) parseFunctionCall(ts []*models.Token, inx *int, node *models.Node) (*models.Node, error) {
	*inx++

	var (
		args       []*models.Node
		parenCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')'", errs.SyntaxError), node.Debug)
		}

		if ts[*inx].Type == tokens.RightParenthesis {
			parenCount--
			if parenCount == 0 {
				*inx++
				break
			}
		}

		if ts[*inx].Type == tokens.LeftParenthesis {
			parenCount++
		}

		arg, err := b.buildNode(ts, inx)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	if parenCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ')'", errs.SyntaxError), node.Debug)
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';'", errs.SyntaxError), node.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';'", errs.SyntaxError), ts[*inx].Debug)
	}

	*inx++
	node.Args = args
	node.Type = tokens.FuncCall

	return node, nil
}

func (b *Builder) parseInlineValue(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value", errs.SyntaxError), ts[*inx].Debug)
	}

	if ts[*inx].Type == tokens.Semicolon || ts[*inx].Type == tokens.Comma || ts[*inx].Type == tokens.RightParenthesis {
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.InlineValue,
			Content:      token.Value,
		}, nil
	}

	node := &models.Node{
		Type:         tokens.Unkown,
		VariableType: tokens.ExpressionVariable,
	}

	var values []*models.Token
	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.Semicolon || ts[*inx].Type == tokens.Comma || ts[*inx].Type == tokens.RightParenthesis {
			*inx++
			break
		}

		values = append(values, ts[*inx])
	}

	if len(values) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression", errs.SyntaxError), token.Debug)
	}

	if len(values) == 1 {
		node.Value = values[0].Value
		typ, err := b.getType(values[0])
		if err != nil {
			return nil, err
		}
		node.VariableType = typ
	} else {
		children, err := b.Build(values)
		if err != nil {
			return nil, err
		}
		node.Children = children
	}

	return node, nil
}

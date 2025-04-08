package ast

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

func (b *Builder) parseLetConst(ts []*models.Token, inx *int) (*models.Node, error) {
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s' 3", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if b.isExpression(ts[*inx]) {
		return &models.Node{
			Type:         ts[*inx].Type,
			VariableType: tokens.ExpressionVariable,
			Content:      ts[*inx].Value,
			Debug:        ts[*inx].Debug,
		}, nil
	}

	if ts[*inx].Type == tokens.Semicolon {
		node.VariableType = tokens.NilVariable
		node.Value = nil
		*inx++
		return node, nil
	}

	if ts[*inx].Type != tokens.Assign {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	var (
		values     []*models.Token
		braceCount int
	)
	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.LeftBrace {
			braceCount++
		}

		if ts[*inx].Type == tokens.RightBrace {
			braceCount--
		}

		if ts[*inx].Type == tokens.Semicolon && braceCount == 0 {
			*inx++
			break
		}

		values = append(values, ts[*inx])
	}

	if len(values) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if len(values) == 1 {
		node.Value = b.getValue(values[0])
		typ, err := b.getType(values[0])
		if err != nil {
			return nil, err
		}
		node.VariableType = typ
	} else {
		values = append(values, SemiColonToken)
		children, err := b.Build(values)
		if err != nil {
			return nil, err
		}

		if len(children) == 1 && children[0].Type == tokens.List {
			node.VariableType = tokens.ListVariable
			node.Children = children[0].Children
		} else {
			node.Children = children
			node.VariableType = tokens.ExpressionVariable
		}
	}

	return node, nil
}

func (b *Builder) parseDefine(ts []*models.Token, inx *int) (*models.Node, error) {
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{', but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftBrace {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	var (
		children   []*models.Token
		braceCount = 1
	)

	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got 'EOF'", errs.SyntaxError), token.Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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
		Debug:   token.Debug,
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: unexpected end of file, got 'EOF'", errs.SyntaxError), token.Debug)
	}

	// Handle multi-part identifiers with dots (e.g., "app.core.func")
	for *inx < len(ts) && ts[*inx].Type == tokens.Dot {
		node.Content += "."
		*inx++

		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier after dot, but got EOF", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type != tokens.Identifier {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier after dot, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
		}

		node.Content += ts[*inx].Value
		*inx++
	}

	if *inx >= len(ts) {
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.ReferenceVariable,
			Reference:    true,
			Content:      node.Content,
			Debug:        node.Debug,
		}, nil
	}

	if ts[*inx].Type != tokens.Assign && ts[*inx].Type != tokens.LeftParenthesis && ts[*inx].Type != tokens.LeftBracket && ts[*inx].Type != tokens.Increment && ts[*inx].Type != tokens.Decrement {
		if ts[*inx].Type == tokens.Semicolon {
			*inx++
		}
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.ReferenceVariable,
			Reference:    true,
			Content:      node.Content,
			Debug:        node.Debug,
		}, nil
	}

	if b.isExpression(ts[*inx]) {
		return &models.Node{
			Type:         ts[*inx].Type,
			VariableType: tokens.ExpressionVariable,
			Content:      ts[*inx].Value,
			Debug:        ts[*inx].Debug,
		}, nil
	}

	if ts[*inx].Type == tokens.LeftParenthesis {
		return b.parseFunctionCall(ts, inx, node)
	}

	if ts[*inx].Type == tokens.LeftBracket {
		node.VariableType = tokens.ReferenceVariable
		return b.parseObjectAccess(ts, inx, node)
	}

	if ts[*inx].Type == tokens.Increment || ts[*inx].Type == tokens.Decrement {
		node.Type = ts[*inx].Type
		*inx++

		if *inx > len(ts) {
			return node, nil
		}

		if ts[*inx].Type != tokens.Semicolon {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected semicolon after increment/decrement operator", errs.SyntaxError), ts[*inx].Debug)
		}

		*inx++
		return node, nil
	}

	if ts[*inx].Type != tokens.Assign {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator or expression, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Type = tokens.Assign

	// Process assignment and value expressions
	*inx++
	var values []*models.Token
	for *inx < len(ts) && ts[*inx].Type != tokens.Semicolon {
		values = append(values, ts[*inx])
		*inx++
	}

	if len(values) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got nothing", errs.SyntaxError), token.Debug)
	}

	if len(values) == 1 {
		node.Value = b.getValue(values[0])
		typ, err := b.getType(values[0])
		if err != nil {
			return nil, err
		}
		node.VariableType = typ
	} else if len(values) > 1 {
		if values[0].Type == tokens.LeftBracket {
			values = append(values, SemiColonToken)
			children, err := b.Build(values)
			if err != nil {
				return nil, err
			}
			node.VariableType = tokens.ListVariable
			node.Type = tokens.List
			node.Children = children[0].Children
			return node, nil
		}
		values = append(values, SemiColonToken)
		children, err := b.Build(values)
		if err != nil {
			return nil, err
		}
		node.VariableType = tokens.ExpressionVariable
		node.Children = children
	}

	return node, nil
}

func (b *Builder) parseInlineValue(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type == tokens.Semicolon || ts[*inx].Type == tokens.Comma || ts[*inx].Type == tokens.RightParenthesis || b.isExpression(ts[*inx]) {
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.InlineValue,
			Content:      "inlineValue",
			Value:        b.getValue(token),
			Debug:        token.Debug,
		}, nil
	}

	node := &models.Node{
		Type:         tokens.Unkown,
		VariableType: tokens.ExpressionVariable,
		Content:      "inlineValue",
		Debug:        token.Debug,
	}

	var values []*models.Token
	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.Semicolon || ts[*inx].Type == tokens.Comma || ts[*inx].Type == tokens.RightParenthesis {
			*inx++
			break
		}

		values = append(values, ts[*inx])
	}

	if len(values) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if len(values) == 1 {
		node.Value = b.getValue(values[0])
		typ, err := b.getType(values[0])
		if err != nil {
			return nil, err
		}
		node.VariableType = typ
	} else {
		values = append(values, SemiColonToken)
		children, err := b.Build(values)
		if err != nil {
			return nil, err
		}

		if len(children) == 1 && children[0].Type == tokens.List {
			node.VariableType = tokens.ListVariable
			node.Children = children[0].Children
		} else {
			node.VariableType = tokens.ExpressionVariable
			node.Children = children
		}
	}

	return node, nil
}

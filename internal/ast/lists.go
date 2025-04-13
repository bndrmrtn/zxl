package ast

import (
	"fmt"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/internal/tokens"
)

func (b *Builder) parseParenthesis(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	var (
		children   []*models.Token
		parenCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got 'EOF'", errs.SyntaxError), ts[*inx-1].Debug)
		}

		if ts[*inx].Type == tokens.RightParenthesis {
			parenCount--
			if parenCount == 0 {
				break
			}
		}

		if ts[*inx].Type == tokens.LeftParenthesis {
			parenCount++
		}

		if ts[*inx].Type == tokens.Comma && parenCount == 1 {
			*inx++
			break
		}

		children = append(children, ts[*inx])
		*inx++
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.RightParenthesis {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	children = append(children, SemiColonToken)

	nodes, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	if len(nodes) == 1 && nodes[0].Type != tokens.FuncCall {
		return nodes[0], nil
	}

	return &models.Node{
		Type:         tokens.LeftParenthesis,
		VariableType: tokens.ExpressionVariable,
		Content:      "expression",
		Children:     nodes,
		Debug:        token.Debug,
	}, nil
}

func (b *Builder) parseList(ts []*models.Token, inx *int) (*models.Node, error) {
	node := &models.Node{
		Type:         tokens.List,
		VariableType: tokens.ListVariable,
		Content:      "list",
		Debug:        ts[*inx].Debug,
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value, but got 'EOF'", errs.SyntaxError), node.Debug)
	}

	if ts[*inx].Type == tokens.RightBracket {
		*inx++
		return node, nil
	}

	var (
		children     []*models.Node
		bracketCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ']', but got 'EOF'", errs.SyntaxError), node.Debug)
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

		value, last, err := b.parseListValue(ts, inx)
		if err != nil {
			return nil, err
		}
		children = append(children, value)
		if last {
			bracketCount = 0
			break
		}
	}

	if bracketCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ']', but got '%s'", errs.SyntaxError, node.Type), node.Debug)
	}

	if ts[*inx].Type == tokens.LeftBracket {
		node.Children = children
		return b.parseObjectAccess(ts, inx, node)
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), node.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon && ts[*inx].Type != tokens.LeftBracket {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++
	node.Children = children
	return node, nil
}

func (b *Builder) parseListValue(ts []*models.Token, inx *int) (*models.Node, bool, error) {
	var (
		children     []*models.Token
		bracketCount = 1
		paranCount   = 0
		braceCount   = 0
	)

	for {
		if *inx >= len(ts) {
			return nil, false, errs.WithDebug(fmt.Errorf("%w: expected ']', but got 'EOF'", errs.SyntaxError), ts[*inx-1].Debug)
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

		if ts[*inx].Type == tokens.LeftParenthesis {
			paranCount++
		}

		if ts[*inx].Type == tokens.RightParenthesis {
			paranCount--
		}

		if ts[*inx].Type == tokens.LeftBrace {
			braceCount++
		}

		if ts[*inx].Type == tokens.RightBrace {
			braceCount--
		}

		if ts[*inx].Type == tokens.Comma && bracketCount == 1 && paranCount == 0 && braceCount == 0 {
			*inx++
			break
		}

		children = append(children, ts[*inx])
		*inx++
	}

	children = append(children, SemiColonToken)

	nodes, err := b.Build(children)
	if err != nil {
		return nil, false, err
	}

	if len(nodes) == 1 && nodes[0].Type != tokens.FuncCall {
		return nodes[0], bracketCount == 0, nil
	}

	return &models.Node{
		Type:         tokens.ListValue,
		Content:      "listValue",
		VariableType: tokens.ExpressionVariable,
		Children:     nodes,
		Debug:        ts[*inx-1].Debug,
	}, bracketCount == 0, nil
}

func (b *Builder) parseArray(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	if token.Type == tokens.Array {
		*inx++
	}

	if token.Type != tokens.Array && token.Type != tokens.LeftBrace {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected array or '{', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if token.Type == tokens.Array && *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{', but got EOF", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftBrace {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	var (
		values     []*models.Node
		braceCount = 1
	)

	for {
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

		keyValue, err := b.parseArrayKeyValues(ts, inx)
		if err != nil {
			return nil, err
		}

		values = append(values, keyValue)
	}

	if braceCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
	}

	if *inx < len(ts) && ts[*inx].Type == tokens.Semicolon {
		*inx++
	}

	return &models.Node{
		Type:         tokens.Array,
		VariableType: tokens.ArrayVariable,
		Content:      "Array{}",
		Children:     values,
	}, nil
}

func (b *Builder) parseArrayKeyValues(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	if token.Type != tokens.Identifier &&
		token.Type != tokens.Number &&
		token.Type != tokens.String {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, number or string, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	keyNode, err := b.Build(append([]*models.Token{}, token, SemiColonToken))
	if err != nil {
		return nil, err
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ':', but got EOF", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Colon {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ':', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	var (
		children     []*models.Token
		braceCount   int
		bracketCount int
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got 'EOF'", errs.SyntaxError), token.Debug)
		}

		if braceCount == 0 && bracketCount == 0 && ts[*inx].Type == tokens.RightBrace {
			break
		}

		if ts[*inx].Type == tokens.LeftBrace {
			braceCount++
		}

		if ts[*inx].Type == tokens.RightBrace {
			braceCount--
		}

		if ts[*inx].Type == tokens.LeftBracket {
			bracketCount++
		}

		if ts[*inx].Type == tokens.RightBracket {
			bracketCount--
		}

		if braceCount == 0 && bracketCount == 0 {
			if ts[*inx].Type == tokens.Comma {
				*inx++
				break
			}
		}

		children = append(children, ts[*inx])
		*inx++
	}

	childNodes, err := b.Build(append(children, SemiColonToken))
	if err != nil {
		return nil, err
	}

	return &models.Node{
		Type:     tokens.ArrayKeyValuePair,
		Content:  "ArrayKey-ValuePair",
		Args:     keyNode,
		Children: childNodes,
	}, nil
}

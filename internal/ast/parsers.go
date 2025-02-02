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
		Type:  token.Type,
		Debug: token.Debug,
	}
	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}
	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected assignment operator, but got '%s'", errs.SyntaxError, ts[*inx].Type), token.Debug)
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

	var values []*models.Token
	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
		}

		if ts[*inx].Type == tokens.Semicolon {
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
		node.VariableType = tokens.ExpressionVariable
		node.Children = children
	}

	return node, nil
}

func (b *Builder) parseFunction(ts []*models.Token, inx *int) (m *models.Node, e error) {
	token := ts[*inx]
	node := &models.Node{
		Type:    token.Type,
		Debug:   token.Debug,
		Content: "func",
	}
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '(', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftParenthesis {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '(', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	var (
		args       []*models.Node
		parenCount = 1
	)

	expectArg := true
	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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

		if ts[*inx].Type == tokens.Comma {
			expectArg = true
			*inx++
			continue
		}

		if !expectArg {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ',' between arguments", errs.SyntaxError), ts[*inx].Debug)
		}

		arg, err := b.buildNode(ts, inx)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		expectArg = false
	}

	if len(args) > 0 && expectArg {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected token: ','", errs.SyntaxError), ts[*inx].Debug)
	}

	node.Args = args

	if parenCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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

		children = append(children, ts[*inx])
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
		return nil, errs.WithDebug(fmt.Errorf("%w: unexpected end of file, got '%s'", errs.SyntaxError, token.Type), token.Debug)
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
			Content:      node.Content,
		}, nil
	}

	if ts[*inx].Type == tokens.Semicolon || (ts[*inx].Type != tokens.Assign && ts[*inx].Type != tokens.LeftParenthesis) {
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.ReferenceVariable,
			Content:      node.Content,
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
	} else {
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

func (b *Builder) parseDefine(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:  token.Type,
		Debug: token.Debug,
	}
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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

func (b *Builder) parseFunctionCall(ts []*models.Token, inx *int, node *models.Node) (*models.Node, error) {
	*inx++

	var (
		args       []*models.Node
		parenCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got '%s'", errs.SyntaxError, node.Type), node.Debug)
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

		arg, err := b.parseFuncCallArg(ts, inx)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}

	if parenCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got '%s'", errs.SyntaxError, node.Type), node.Debug)
	}

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, node.Type), node.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++
	node.Args = args
	node.Type = tokens.FuncCall
	node.VariableType = tokens.FunctionCallVariable

	return node, nil
}

func (b *Builder) parseInlineValue(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value", errs.SyntaxError), ts[*inx].Debug)
	}

	if ts[*inx].Type == tokens.Semicolon || ts[*inx].Type == tokens.Comma || ts[*inx].Type == tokens.RightParenthesis || b.isExpression(ts[*inx]) {
		return &models.Node{
			Type:         token.Type,
			VariableType: tokens.InlineValue,
			Content:      token.Value,
			Value:        b.getValue(token),
			Debug:        ts[*inx].Debug,
		}, nil
	}

	node := &models.Node{
		Type:         tokens.Unkown,
		VariableType: tokens.ExpressionVariable,
		Debug:        token.Debug,
	}

	var values []*models.Token
	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
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
		node.VariableType = tokens.ExpressionVariable
		node.Children = children
	}

	return node, nil
}

func (b *Builder) parseNamespace(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:  token.Type,
		Debug: token.Debug,
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++
	return node, nil
}

func (b *Builder) parseReturn(ts []*models.Token, inx *int) (*models.Node, error) {
	node := &models.Node{
		Type:         ts[*inx].Type,
		VariableType: tokens.ExpressionVariable,
		Content:      "return",
	}

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got EOF", errs.SyntaxError), node.Debug)
	}

	if ts[*inx].Type == tokens.Semicolon {
		node.Value = nil
		node.Type = tokens.EmptyReturn
		node.VariableType = tokens.EmptyReturnValue
		*inx++
		return node, nil
	}

	var children = []*models.Token{}
	for {
		if ts[*inx].Type == tokens.Semicolon {
			*inx++
			break
		}

		children = append(children, ts[*inx])
		*inx++
	}

	children = append(children, SemiColonToken)
	child, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	node.Children = child
	return node, nil
}

func (b *Builder) parseIf(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]

	node := &models.Node{
		Type:    token.Type,
		Debug:   token.Debug,
		Content: "if",
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got EOF", errs.SyntaxError), token.Debug)
	}

	var args []*models.Token
	for {
		if ts[*inx].Type == tokens.LeftBrace {
			break
		}

		args = append(args, ts[*inx])
		*inx++
	}

	if len(args) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got EOF", errs.SyntaxError), token.Debug)
	}

	args = append(args, SemiColonToken)
	bArgs, err := b.Build(args)
	if err != nil {
		return nil, err
	}
	node.Args = bArgs

	*inx++

	var (
		children   []*models.Token
		braceCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
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

		children = append(children, ts[*inx])
		*inx++
	}

	if braceCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
	}

	var statement []*models.Node

	child, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	statement = append(statement, &models.Node{
		Type:     tokens.Then,
		Children: child,
		Content:  "then",
	})

	if *inx >= len(ts) || (ts[*inx].Type != tokens.Else && ts[*inx].Type != tokens.ElseIf) {
		node.Children = statement
		return node, nil
	}

	if ts[*inx].Type == tokens.Else {
		elseNode, err := b.handleElse(ts, inx)
		if err != nil {
			return nil, err
		}
		statement = append(statement, elseNode)
		node.Children = statement
	}

	return node, nil
}

func (b *Builder) handleElse(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:    tokens.Else,
		Debug:   token.Debug,
		Content: "else",
	}
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{' or 'if', but got EOF", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftBrace && ts[*inx].Type != tokens.If {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{' or 'if', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	if ts[*inx].Type == tokens.If {
		ifNode, err := b.parseIf(ts, inx)
		if err != nil {
			return nil, err
		}
		node.Children = []*models.Node{ifNode}
		return node, nil
	}

	*inx++

	var (
		elseChildren   []*models.Token
		elseBraceCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type == tokens.RightBrace {
			elseBraceCount--
			if elseBraceCount == 0 {
				*inx++
				break
			}
		}

		if ts[*inx].Type == tokens.LeftBrace {
			elseBraceCount++
		}

		elseChildren = append(elseChildren, ts[*inx])
		*inx++
	}

	if elseBraceCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
	}

	elseChild, err := b.Build(elseChildren)
	if err != nil {
		return nil, err
	}

	node.Children = elseChild
	return node, nil
}

func (b *Builder) parseFuncCallArg(ts []*models.Token, inx *int) (*models.Node, error) {
	var (
		children   []*models.Token
		parenCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got EOF", errs.SyntaxError), ts[*inx-1].Debug)
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

	children = append(children, SemiColonToken)

	nodes, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	return &models.Node{
		Type:         tokens.FuncArg,
		Content:      "argument",
		VariableType: tokens.ExpressionVariable,
		Children:     nodes,
	}, nil
}

func (b *Builder) parseUse(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:  token.Type,
		Debug: token.Debug,
	}

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	node.Value = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon && ts[*inx].Type != tokens.As {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Value = ts[*inx].Value

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, token.Type), token.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

	return node, nil
}

func (b *Builder) parseWhile(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:    token.Type,
		Content: "while",
		Debug:   token.Debug,
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got EOF", errs.SyntaxError), token.Debug)
	}

	var args []*models.Token
	for {
		if ts[*inx].Type == tokens.LeftBrace {
			break
		}

		args = append(args, ts[*inx])
		*inx++
	}

	if len(args) == 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got EOF", errs.SyntaxError), token.Debug)
	}

	args = append(args, SemiColonToken)
	bArgs, err := b.Build(args)
	if err != nil {
		return nil, err
	}
	node.Args = bArgs

	*inx++

	var (
		children   []*models.Token
		braceCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
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

		children = append(children, ts[*inx])
		*inx++
	}

	if braceCount != 0 {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got EOF", errs.SyntaxError), token.Debug)
	}

	child, err := b.Build(children)
	if err != nil {
		return nil, err
	}

	node.Children = child

	return node, nil
}

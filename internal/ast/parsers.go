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

	var values []*models.Token
	for {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), token.Debug)
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

func (b *Builder) parseFunction(ts []*models.Token, inx *int) (m *models.Node, e error) {
	token := ts[*inx]
	node := &models.Node{
		Type:    token.Type,
		Debug:   token.Debug,
		Content: "func",
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '(', but got 'EOF'", errs.SyntaxError), token.Debug)
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
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got 'EOF'", errs.SyntaxError), token.Debug)
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

	if ts[*inx].Type == tokens.Semicolon || (ts[*inx].Type != tokens.Assign && ts[*inx].Type != tokens.LeftParenthesis && ts[*inx].Type != tokens.LeftBracket) {
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

func (b *Builder) parseFunctionCall(ts []*models.Token, inx *int, node *models.Node) (*models.Node, error) {
	*inx++

	var (
		args       []*models.Node
		parenCount = 1
	)

	for {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected ')', but got 'EOF'", errs.SyntaxError), node.Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), node.Debug)
	}

	var funcCallTokens []*models.Token

	if ts[*inx].Type == tokens.Dot {
		*inx++
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), node.Debug)
		}

		for {
			if ts[*inx].Type == tokens.Semicolon {
				break
			}

			funcCallTokens = append(funcCallTokens, ts[*inx])
			*inx++
		}
	}

	funcCallTokens = append(funcCallTokens, SemiColonToken)
	funcCallChild, err := b.Build(funcCallTokens)
	if err != nil {
		return nil, err
	}

	if ts[*inx].Type != tokens.Semicolon {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++
	node.Args = args
	node.Type = tokens.FuncCall
	node.VariableType = tokens.FunctionCallVariable
	node.Children = funcCallChild

	return node, nil
}

func (b *Builder) parseInlineValue(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value, but got 'EOF'", errs.SyntaxError), ts[*inx].Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), token.Debug)
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
		Debug:        ts[*inx].Debug,
	}

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), node.Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), token.Debug)
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

	if *inx >= len(ts) || ts[*inx].Type != tokens.Else {
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{' or 'if', but got 'EOF'", errs.SyntaxError), token.Debug)
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
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got 'EOF'", errs.SyntaxError), token.Debug)
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

		if ts[*inx].Type == tokens.LeftBracket {
			bracketCount := 0
			for {
				if *inx >= len(ts) {
					return nil, errs.WithDebug(fmt.Errorf("%w: expected ']', but got 'EOF'", errs.SyntaxError), ts[*inx-1].Debug)
				}

				if ts[*inx].Type == tokens.LeftBracket {
					bracketCount++
				}

				if ts[*inx].Type == tokens.RightBracket {
					bracketCount--
					if bracketCount == 0 {
						break
					}
				}

				children = append(children, ts[*inx])
				*inx++
			}
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

	if len(nodes) == 1 && nodes[0].Type != tokens.FuncCall {
		return nodes[0], nil
	}

	return &models.Node{
		Type:         tokens.FuncArg,
		Content:      "argument",
		VariableType: tokens.ExpressionVariable,
		Children:     nodes,
		Debug:        ts[*inx-1].Debug,
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Content = ts[*inx].Value
	node.Value = ts[*inx].Value

	*inx++
	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), token.Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	node.Value = ts[*inx].Value

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), token.Debug)
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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected value or expression, but got 'EOF'", errs.SyntaxError), token.Debug)
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

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), node.Debug)
	}

	if ts[*inx].Type != tokens.Semicolon && ts[*inx].Type != tokens.LeftBracket {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	if ts[*inx].Type == tokens.LeftBracket {
		node.Children = children
		return b.parseObjectAccess(ts, inx, node)
	}

	*inx++
	node.Children = children
	return node, nil
}

func (b *Builder) parseListValue(ts []*models.Token, inx *int) (*models.Node, bool, error) {
	var (
		children     []*models.Token
		bracketCount = 1
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

		if ts[*inx].Type == tokens.Comma && bracketCount == 1 {
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

func (b *Builder) parseFor(ts []*models.Token, inx *int) (*models.Node, error) {
	token := ts[*inx]
	node := &models.Node{
		Type:    token.Type,
		Content: "for",
		Debug:   token.Debug,
	}

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected iterable expression, but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier, but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	iteratorVariable := ts[*inx]

	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected 'in', but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.In {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected 'in', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	*inx++

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

	iterator, err := b.Build([]*models.Token{iteratorVariable, SemiColonToken})
	if err != nil {
		return nil, err
	}

	args = append(args, SemiColonToken)
	bArgs, err := b.Build(args)
	if err != nil {
		return nil, err
	}

	node.Args = append(node.Args, iterator...)
	node.Args = append(node.Args, bArgs...)

	*inx++

	var (
		children   []*models.Token
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

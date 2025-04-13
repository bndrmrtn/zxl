package ast

import (
	"fmt"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/internal/tokens"
)

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

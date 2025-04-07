package ast

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

func (b *Builder) parseFunction(ts []*models.Token, inx *int) (m *models.Node, e error) {
	token := ts[*inx]
	node := &models.Node{
		Type:    tokens.Function,
		Debug:   token.Debug,
		Content: "fn",
		Map:     map[string]interface{}{},
	}
	*inx++

	if *inx >= len(ts) {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier or '(', but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.Identifier && ts[*inx].Type != tokens.LeftParenthesis {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected identifier or '(', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	if ts[*inx].Type == tokens.LeftParenthesis {
		node.Type = tokens.InlineFunction
		node.VariableType = tokens.FunctionVariable
	} else {
		node.Content = ts[*inx].Value
	}

	*inx++

	if node.Type == tokens.Function {
		if *inx >= len(ts) {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '(', but got 'EOF'", errs.SyntaxError), token.Debug)
		}

		if ts[*inx].Type != tokens.LeftParenthesis {
			return nil, errs.WithDebug(fmt.Errorf("%w: expected '(', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
		}

		*inx++
	}

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
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{' or '=>', but got 'EOF'", errs.SyntaxError), token.Debug)
	}

	if ts[*inx].Type != tokens.LeftBrace && ts[*inx].Type != tokens.Arrow {
		return nil, errs.WithDebug(fmt.Errorf("%w: expected '{' or '=>', but got '%s'", errs.SyntaxError, ts[*inx].Type), ts[*inx].Debug)
	}

	var children []*models.Token

	if ts[*inx].Type == tokens.Arrow {
		*inx++
		// Bind return token for auto-returning arrow function values
		children = append(children, &models.Token{
			Type: tokens.Return,
		})
		for {
			if *inx >= len(ts) {
				return nil, errs.WithDebug(fmt.Errorf("%w: expected ';', but got 'EOF'", errs.SyntaxError), token.Debug)
			}

			if ts[*inx].Type == tokens.Semicolon {
				break
			}

			children = append(children, ts[*inx])
			*inx++
		}
	} else {
		var braceCount = 1

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
	}

	child, err := b.Build(append(children, SemiColonToken))
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

	node.Args = args
	node.Type = tokens.FuncCall
	node.VariableType = tokens.FunctionCallVariable

	if *inx >= len(ts) {
		return node, nil
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

	if ts[*inx].Type == tokens.Semicolon {
		*inx++
	}

	node.Children = funcCallChild

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
		if *inx >= len(ts) {
			break
		}

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

		if ts[*inx].Type == tokens.LeftBrace {
			braceCount := 0
			for {
				if *inx >= len(ts) {
					return nil, errs.WithDebug(fmt.Errorf("%w: expected '}', but got 'EOF'", errs.SyntaxError), ts[*inx-1].Debug)
				}

				if ts[*inx].Type == tokens.LeftBrace {
					braceCount++
				}

				if ts[*inx].Type == tokens.RightBrace {
					braceCount--
					if braceCount == 0 {
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

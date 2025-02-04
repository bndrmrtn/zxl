package prettycode

import (
	"html"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

type PrettyCode struct {
	tokens []*models.Token
}

func New(ts []*models.Token) *PrettyCode {
	return &PrettyCode{
		tokens: ts,
	}
}

func (p *PrettyCode) HighlightHtml() string {
	var sb strings.Builder
	defer sb.Reset()

	for i, token := range p.tokens {
		next := p.nextToken(i)

		sb.WriteString(p.highlightToken(HtmlMode, token, next))
	}

	return sb.String()
}

func (p *PrettyCode) highlightToken(mode Mode, token *models.Token, next *models.Token) string {
	token.Value = html.EscapeString(token.Value)
	switch token.Type {
	case tokens.Let, tokens.Const,
		tokens.Function, tokens.Define, tokens.Return,
		tokens.Namespace, tokens.Use, tokens.As, tokens.From,
		tokens.While, tokens.For,
		tokens.If, tokens.Else:
		return p.highlightKeyword(mode, token.Value)
	case tokens.Identifier:
		if next != nil && next.Type == tokens.LeftParenthesis {
			return p.highlightFunction(mode, token.Value)
		}
		if next != nil && next.Type == tokens.Dot {
			return p.highlightIdentifierWithDot(mode, token.Value)
		}
		return p.highlightIdentifier(mode, token.Value)
	case tokens.String:
		return p.highlightString(mode, token.Value)
	case tokens.Number:
		return p.highlightNumber(mode, token.Value)
	case tokens.LeftParenthesis, tokens.RightParenthesis,
		tokens.LeftBrace, tokens.RightBrace,
		tokens.LeftBracket, tokens.RightBracket,
		tokens.Comma, tokens.Dot, tokens.Colon:
		return p.highlightBracket(mode, token.Value)
	case tokens.Addition, tokens.Subtraction, tokens.Multiplication, tokens.Division,
		tokens.Equation, tokens.NotEquation, tokens.Greater, tokens.GreaterOrEqual, tokens.Less, tokens.LessOrEqual,
		tokens.And, tokens.Or, tokens.Not,
		tokens.Assign:
		return p.highlightOperator(mode, token.Value)
	case tokens.NewLine, tokens.WhiteSpace:
		return token.Value
	default:
		return p.highlightUnknown(mode, token.Value)
	}
}

func (p *PrettyCode) nextToken(i int) *models.Token {
	var next *models.Token

	if i+1 < len(p.tokens) {
		next = p.tokens[i+1]
		if next.Type == tokens.WhiteSpace || next.Type == tokens.NewLine {
			next = p.nextToken(i + 1)
		}
	}

	return next
}

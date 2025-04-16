package prettycode

import (
	"html"
	"io"
	"strings"

	"github.com/flarelang/flare/internal/lexer"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
)

// PrettyCode is a struct that holds the tokens of the code
type PrettyCode struct {
	tokens []*models.Token
}

func New(r io.Reader) (*PrettyCode, error) {
	lx := lexer.New("snippet")

	ts, err := lx.Parse(r)
	if err != nil {
		return nil, err
	}

	return &PrettyCode{
		tokens: ts,
	}, nil
}

// NewToken creates a new PrettyCode struct from tokens
func NewToken(ts []*models.Token) *PrettyCode {
	return &PrettyCode{
		tokens: ts,
	}
}

// HighlightHtml highlights the code with html tags
func (p *PrettyCode) HighlightHtml() string {
	var sb strings.Builder
	defer sb.Reset()

	for i, token := range p.tokens {
		next := p.nextToken(i)

		sb.WriteString(p.highlightToken(HtmlMode, token, next))
	}

	return sb.String()
}

// HighlightConsole highlights the code with console
func (p *PrettyCode) HighlightConsole() string {
	var sb strings.Builder
	defer sb.Reset()

	for i, token := range p.tokens {
		next := p.nextToken(i)

		sb.WriteString(p.highlightToken(ConsoleMode, token, next))
	}

	return sb.String()
}

// highlightToken highlights the token with the given mode
func (p *PrettyCode) highlightToken(mode Mode, token *models.Token, next *models.Token) string {
	if mode == HtmlMode {
		token.Value = html.EscapeString(token.Value)
	}

	switch token.Type {
	case tokens.Let, tokens.Const,
		tokens.Function, tokens.Define, tokens.Return,
		tokens.Namespace, tokens.Use, tokens.As, tokens.From,
		tokens.While, tokens.For, tokens.Spin,
		tokens.If, tokens.Else, tokens.In, tokens.Array, tokens.Error:
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
	case tokens.Number, tokens.Bool, tokens.Nil:
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

// nextToken returns the next token in the list
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

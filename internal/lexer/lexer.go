package lexer

import (
	"fmt"
	"io"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// Lexer is a lexical analyzer
type Lexer struct {
	filename string
}

// New creates a new lexer
func New(filename string) *Lexer {
	return &Lexer{
		filename: filename,
	}
}

// Parse reads the content of the reader and returns the tokens
func (lx *Lexer) Parse(r io.Reader) ([]*models.Token, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", errs.ErrFailedToReadContent, err)
	}

	return lx.parse(string(b))
}

// parse reads the content of the string and returns the tokens
func (lx *Lexer) parse(s string) ([]*models.Token, error) {
	var (
		pos     int
		line    int = 1
		col     int
		fileLen = len(s)
		parsed  []*models.Token
	)

	for pos < len(s) {
		switch s[pos] {
		// Handle new lines
		case '\n':
			line++
			col = 0
		// Handle single line comments
		case '/':
			if pos+1 < len(s) && s[pos+1] == '/' {
				// Skip the entire comment line
				pos += 2
				for pos < len(s) && s[pos] != '\n' {
					pos++
				}
				// New line reached
				line++
				col = 0
			} else if pos+1 < len(s) && s[pos+1] == '*' {
				pos += 2
				for pos < len(s) {
					if s[pos] == '*' && pos+1 < len(s) && s[pos+1] == '/' {
						pos += 2
						break
					} else if s[pos] == '\n' {
						line++
						col = 0
					} else {
						pos++
					}
				}
			} else {
				col++
			}
		// Handle strings
		case '"', '\'':
			// Skip the opening quote
			quote := s[pos] // Store the opening quote character
			pos++

			var value string
			start := pos // Start of the string content

			for pos < len(s) {
				if s[pos] == '\\' {
					// Escape sequence detected
					if pos+1 < len(s) {
						pos += 2
					} else {
						// Syntax error: escape character at the end
						return nil, errs.WithDebug(fmt.Errorf("%w: incomplete escape sequence", errs.SyntaxError), &models.Debug{
							Line:   line,
							Column: col,
							File:   lx.filename,
							Near:   lx.near(s, pos, fileLen),
						})
					}
				} else if s[pos] == quote {
					// Closing quote found
					value = s[start:pos] // Extract the string content
					pos++                // Skip the closing quote
					break
				} else {
					pos++
				}
			}

			// If we exited the loop without finding a closing quote
			if value == "" && pos >= len(s) {
				return nil, errs.WithDebug(fmt.Errorf("%w: missing closing quote for string starting", errs.SyntaxError), &models.Debug{
					Line:   line,
					Column: col,
					File:   lx.filename,
					Near:   lx.near(s, pos, fileLen),
				})
			}

			// Add the string token to the parsed tokens
			parsed = append(parsed, &models.Token{
				Type:  tokens.String,
				Value: string(quote) + value + string(quote),
				Debug: &models.Debug{
					Line:   line,
					Column: col,
					File:   lx.filename,
					Near:   lx.near(s, pos, fileLen),
				},
				Map: map[string]any{
					"quote": quote,
				},
			})
			pos--
		case ' ', '\t':
			col++
		case ';', ':', ',', '.', '(', ')', '{', '}', '[', ']':
			ch := s[pos]
			parsed = append(parsed, &models.Token{
				Type:  lx.getCharIdent(ch),
				Value: string(ch),
				Debug: &models.Debug{
					Line:   line,
					Column: col,
					File:   lx.filename,
					Near:   lx.near(s, pos, fileLen),
				},
			})
		case '=':
			if pos+1 < len(s) && s[pos+1] == '=' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Equation,
					Value: "==",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				pos++
			} else {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Assign,
					Value: "=",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
		case '+':
			if pos+1 < len(s) && s[pos+1] == '+' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Increment,
					Value: "++",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				pos++
			} else {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Addition,
					Value: "+",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
			col++

		default:
			start := pos
			if isLetter(s[pos]) || s[pos] == '_' {
				// Identifier parsing
				for pos < len(s) && (isLetter(s[pos]) || isDigit(s[pos]) || s[pos] == '_') {
					pos++
				}
				value := s[start:pos]
				// Appending the identifier
				parsed = append(parsed, &models.Token{
					Type:  lx.getIdentType(value),
					Value: value,
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				pos--
			} else if isDigit(s[pos]) || (s[pos] == '.' && pos+1 < len(s) && isDigit(s[pos+1])) {
				// Number parsing (integer or float)
				isFloat := false
				for pos < len(s) && (isDigit(s[pos]) || s[pos] == '.') {
					if s[pos] == '.' {
						if isFloat {
							// Second dot found, invalid number
							return nil, errs.WithDebug(fmt.Errorf("%w invalid number format", errs.SyntaxError), &models.Debug{
								Line:   line,
								Column: col,
								File:   lx.filename,
								Near:   lx.near(s, pos, fileLen),
							})
						}
						isFloat = true
					}
					pos++
				}
				value := s[start:pos]
				// Appending the number
				parsed = append(parsed, &models.Token{
					Type:  tokens.Number,
					Value: value,
					Map: map[string]any{
						"isFloat": isFloat,
					},
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				pos--
			} else {
				// Invalid token
				return nil, errs.WithDebug(fmt.Errorf("%w: unexpected character '%c'", errs.SyntaxError, s[pos]), &models.Debug{
					Line:   line,
					Column: col,
					File:   lx.filename,
					Near:   lx.near(s, pos, fileLen),
				})
			}
		}
		pos++
	}

	return parsed, nil
}

func (lx *Lexer) getCharIdent(ch byte) tokens.TokenType {
	switch ch {
	case '*':
		return tokens.Multiplication
	case '/':
		return tokens.Division
	case ';':
		return tokens.Semicolon
	case ':':
		return tokens.Colon
	case ',':
		return tokens.Comma
	case '.':
		return tokens.Dot
	case '(':
		return tokens.LeftParenthesis
	case ')':
		return tokens.RightParenthesis
	case '{':
		return tokens.LeftBrace
	case '}':
		return tokens.RightBrace
	case '[':
		return tokens.LeftBracket
	case ']':
		return tokens.RightBracket
	default:
		return tokens.Unkown
	}
}

func (lx *Lexer) getIdentType(s string) tokens.TokenType {
	switch s {
	case "use":
		return tokens.Use
	case "as":
		return tokens.As
	case "from":
		return tokens.From
	case "let":
		return tokens.Let
	case "const":
		return tokens.Const
	case "define":
		return tokens.Define
	case "return":
		return tokens.Return
	case "fn":
		return tokens.Function
	case "new":
		return tokens.New
	case "true", "false":
		return tokens.Bool
	case "nil":
		return tokens.Nil
	case "if":
		return tokens.If
	case "else":
		return tokens.Else
	case "elseif":
		return tokens.ElseIf
	default:
		return tokens.Identifier
	}
}

func (lx *Lexer) near(s string, pos int, fileLen int) string {
	if pos < 0 || pos >= fileLen || len(s) == 0 {
		return ""
	}

	start := max(0, pos-20)
	end := min(pos+20, fileLen)
	substr := s[start:end]

	return strings.TrimSpace(substr)
}

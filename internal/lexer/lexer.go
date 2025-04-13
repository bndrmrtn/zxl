package lexer

import (
	"fmt"
	"io"
	"strings"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
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
	// Fix Carriage Return error on Windows PCs
	s = strings.ReplaceAll(s, "\r", "")

	var (
		runes   = []rune(s)
		pos     int
		line    int = 1
		col     int = 1
		fileLen     = len(s)
		parsed  []*models.Token
	)

	for pos < len(runes) {
		switch runes[pos] {
		// Handle new lines
		case '\n':
			// Add the new line token for debugging purposes
			parsed = append(parsed, &models.Token{
				Type:  tokens.NewLine,
				Value: "\n",
				Debug: &models.Debug{
					Line:   line,
					Column: col,
					File:   lx.filename,
					Near:   lx.near(s, pos, fileLen),
				},
			})
			line++
			col = 1
		// Handle single line comments
		case '/':
			if pos+1 < len(s) && runes[pos+1] == '/' {
				// Skip the entire comment line
				pos += 2

				var sb strings.Builder
				sb.WriteString("//")

				for pos < len(s) && runes[pos] != '\n' {
					sb.WriteByte(s[pos])
					pos++
				}

				// Add the single line comment token for debugging purposes
				sb.WriteByte('\n')
				parsed = append(parsed, &models.Token{
					Type:  tokens.SingleLineComment,
					Value: sb.String(),
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				sb.Reset()
				// New line reached
				line++
				col = 1
			} else if pos+1 < len(s) && runes[pos+1] == '*' {
				pos += 2

				var sb strings.Builder
				sb.WriteString("/*")

				for pos < len(s) {
					if runes[pos] == '*' && pos+1 < len(s) && runes[pos+1] == '/' {
						pos += 2

						// Add the multi line comment token for debugging purposes
						sb.WriteString("*/")
						parsed = append(parsed, &models.Token{
							Type:  tokens.MultiLineComment,
							Value: sb.String(),
							Debug: &models.Debug{
								Line:   line,
								Column: col,
								File:   lx.filename,
								Near:   lx.near(s, pos, fileLen),
							},
						})
						sb.Reset()
						break
					} else if runes[pos] == '\n' {
						sb.WriteByte('\n')
						line++
						col = 1
					} else {
						sb.WriteByte(s[pos])
						pos++
					}
				}
			} else {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Division,
					Value: "/",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col++
			}
		// Handle strings
		case '"', '\'', '`':
			// Skip the opening quote
			quote := runes[pos] // Store the opening quote character
			pos++

			var value string
			start := pos // Start of the string content

			for pos < len(s) {
				if runes[pos] == '\\' {
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
				} else if runes[pos] == quote {
					// Closing quote found
					value = string(runes[start:pos]) // Extract the string content
					pos++                            // Skip the closing quote
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

			typ := tokens.String
			if quote == '`' {
				typ = tokens.TemplateLiteral
			}

			// Add the string token to the parsed tokens
			parsed = append(parsed, &models.Token{
				Type:  typ,
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
			// Whitespace token for debugging purposes
			parsed = append(parsed, &models.Token{
				Type:  tokens.WhiteSpace,
				Value: string(s[pos]),
				Debug: &models.Debug{
					Line:   line,
					Column: col,
					File:   lx.filename,
					Near:   lx.near(s, pos, fileLen),
				},
			})
			col++
		case ';', ':', ',', '.', '(', ')', '{', '}', '[', ']':
			ch := runes[pos]
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
			col++
		case '=':
			if pos+1 < len(s) && runes[pos+1] == '=' {
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
			} else if pos+1 < len(s) && runes[pos+1] == '>' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Arrow,
					Value: "=>",
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
		case '!':
			if pos+1 < len(s) && runes[pos+1] == '=' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.NotEquation,
					Value: "!=",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col += 2
				pos++
			} else {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Not,
					Value: "!",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
		case '&':
			if pos+1 < len(s) && runes[pos+1] == '&' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.And,
					Value: "&&",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col += 2
				pos++
			} else {
				parsed = append(parsed, &models.Token{
					Type:  tokens.And,
					Value: "&",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col++
			}
		case '|':
			if pos+1 < len(s) && runes[pos+1] == '|' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Or,
					Value: "||",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col += 2
				pos++
			} else {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Or,
					Value: "|",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col++
			}
		case '+':
			if pos+1 < len(s) && runes[pos+1] == '+' {
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
		case '-':
			if pos+1 < len(s) && runes[pos+1] == '-' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Decrement,
					Value: "--",
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
					Type:  tokens.Subtraction,
					Value: "-",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
			col++
		case '*':
			if pos+1 < len(s) && runes[pos+1] == '*' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.Power,
					Value: "**",
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
					Type:  tokens.Multiplication,
					Value: "*",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
		case '<':
			if pos+1 < len(s) && runes[pos+1] == '>' {
				var str strings.Builder
				pos += 2
				str.WriteString("<>")
				for pos < len(s)-2 && !(runes[pos] == '<' && (runes[pos+1] == '/' && runes[pos+2] == '>')) {
					str.WriteByte(s[pos])
					pos++
					col++

					if runes[pos] == '\n' {
						line++
						col = 0
					}
				}
				str.WriteString("</>")
				line--
				pos += 2
				parsed = append(parsed, &models.Token{
					Type:  tokens.TemplateLiteral,
					Value: str.String(),
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			} else if pos+1 < len(s) && runes[pos+1] == '=' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.LessOrEqual,
					Value: "<=",
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
					Type:  tokens.Less,
					Value: "<",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
		case '>':
			if pos+1 < len(s) && runes[pos+1] == '=' {
				parsed = append(parsed, &models.Token{
					Type:  tokens.GreaterOrEqual,
					Value: ">=",
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
					Type:  tokens.Greater,
					Value: ">",
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
			}
		default:
			start := pos
			if isLetter(runes[pos]) {
				// Identifier parsing
				for pos < len(runes) && (isLetter(runes[pos]) || isDigit(runes[pos])) {
					pos++
				}
				value := string(runes[start:pos])
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
			} else if isDigit(runes[pos]) || (s[pos] == '.' && pos+1 < len(s) && isDigit(runes[pos+1])) {
				// Number parsing (integer or float)
				isFloat := false
				for pos < len(s) && (isDigit(runes[pos]) || runes[pos] == '.') {
					if runes[pos] == '.' {
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
				value := string(runes[start:pos])
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
				parsed = append(parsed, &models.Token{
					Type:  tokens.Unkown,
					Value: string(s[pos]),
					Debug: &models.Debug{
						Line:   line,
						Column: col,
						File:   lx.filename,
						Near:   lx.near(s, pos, fileLen),
					},
				})
				col++
			}
		}
		pos++
	}

	return parsed, nil
}

func (lx *Lexer) getCharIdent(ch rune) tokens.TokenType {
	switch ch {
	case '*':
		return tokens.Multiplication
	case '/':
		return tokens.Division
	case '!':
		return tokens.Not
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
	case "namespace":
		return tokens.Namespace
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
	case "throw":
		return tokens.Throw
	case "true", "false":
		return tokens.Bool
	case "nil":
		return tokens.Nil
	case "if":
		return tokens.If
	case "else":
		return tokens.Else
	case "for":
		return tokens.For
	case "while":
		return tokens.While
	case "in":
		return tokens.In
	case "array":
		return tokens.Array
	case "spin":
		return tokens.Spin
	case "error":
		return tokens.Error
	default:
		return tokens.Identifier
	}
}

func (lx *Lexer) near(s string, pos int, fileLen int) string {
	if pos < 0 || pos >= fileLen || len(s) == 0 {
		return ""
	}

	start := max(0, pos-30)
	end := min(pos+30, fileLen)
	substr := s[start:end]

	return strings.TrimSpace(substr)
}

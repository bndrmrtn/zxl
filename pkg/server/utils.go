package server

import (
	"bytes"
	"strings"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/lexer"
	"github.com/bndrmrtn/flare/pkg/prettycode"
)

func (s *Server) makePrettyCode(htmlErr *errs.HtmlError) string {
	lx := lexer.New("source")
	ts, err := lx.Parse(bytes.NewReader(htmlErr.Code))
	if err != nil {
		return htmlErr.Error()
	}

	pretty := prettycode.NewToken(ts)
	code := pretty.HighlightHtml()
	hlCodeParts := strings.Split(code, "\n")
	codeParts := strings.Split(string(htmlErr.Code), "\n")
	errLine := htmlErr.Debug().GetLine()

	if len(hlCodeParts) < errLine {
		return htmlErr.Error()
	}

	hlCodeParts[errLine-1] = "<span style=\"color:#fb2c36;text-decoration:line-through\">" + codeParts[errLine-1] + "</span>"
	htmlErr.Code = []byte(strings.Join(hlCodeParts, "\n"))

	return htmlErr.Error()
}

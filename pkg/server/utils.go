package server

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/pkg/prettycode"
)

// getExecutablePath gets the executable path
func (s *Server) getExecutablePath(path string) (string, error) {
	if filepath.Ext(path) != ".zx" {
		stat, err := os.Stat(path + ".zx")
		if err == nil && !stat.IsDir() {
			return path + ".zx", nil
		}
	}

	stat, err := os.Stat(path)
	if err == nil && !stat.IsDir() {
		return path, nil
	}

	path = filepath.Join(path, "index.zx")
	stat, err = os.Stat(path)
	if err == nil && !stat.IsDir() {
		return path, nil
	}

	return "", fmt.Errorf("%w: file not found", errNotFound)
}

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

package server

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
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

	return "", fmt.Errorf("file not found")
}

// handleError handles the error
func (s *Server) handleError(err error, w http.ResponseWriter) {
	var de errs.DebugError
	if errors.As(err, &de) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		htmlErr := de.HttpError()
		if htmlErr == nil {
			w.Write([]byte(de.Error()))
			return
		}
		w.Write([]byte(s.makePrettyCode(htmlErr)))
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func (s *Server) makePrettyCode(htmlErr *errs.HtmlError) string {
	lx := lexer.New("")
	ts, err := lx.Parse(bytes.NewReader(htmlErr.Code))
	if err != nil {
		return htmlErr.Error()
	}

	pretty := prettycode.New(ts)
	code := pretty.HighlightHtml()
	hlCodeParts := strings.Split(code, "\n")
	codeParts := strings.Split(string(htmlErr.Code), "\n")
	errLine := htmlErr.Debug().GetLine()

	if len(hlCodeParts) < errLine {
		return htmlErr.Error()
	}

	hlCodeParts[errLine-1] = "<span style=\"color:#fb2c36;text-decoration:line-through\">" + codeParts[errLine-1] + "</span>"
	htmlErr.Code = []byte(strings.Join(hlCodeParts, "\n"))

	htmlErr.Debug()
	return htmlErr.Error()
}

package errs

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"os"
	"strconv"

	"github.com/bndrmrtn/zxl/internal/version"
)

//go:embed assets/*
var htmlFiles embed.FS

// HtmlError is an error with HTML debug information
type HtmlError struct {
	err *DebugError

	Code []byte
}

// NewHtmlError creates a new HtmlError
func NewHtmlError(de *DebugError) *HtmlError {
	fileContent, err := os.ReadFile(de.debug.File)
	if err != nil {
		return &HtmlError{err: de}
	}

	he := &HtmlError{err: de}
	he.Code = fileContent

	return he
}

func (he *HtmlError) Debug() *DebugError {
	return he.err
}

// Error returns the error message with HTML debug information
func (he *HtmlError) Error() string {
	tmpl, err := template.ParseFS(htmlFiles, "assets/error.html")
	if err != nil {
		return he.err.Error()
	}

	props, err := he.getTemplateProps()
	if err != nil {
		return he.err.Error()
	}

	var wr bytes.Buffer
	err = tmpl.Execute(&wr, props)

	if err != nil {
		return err.Error()
	}

	return wr.String()
}

func (he *HtmlError) getTemplateProps() (map[string]interface{}, error) {
	cssFile, err := htmlFiles.Open("assets/main.css")
	if err != nil {
		return nil, he.err
	}
	defer cssFile.Close()

	css, err := io.ReadAll(cssFile)
	if err != nil {
		return nil, he.err
	}

	startLine, endLine := he.getLineNumbers(he.err.debug.Line, bytes.Count(he.Code, []byte("\n")))
	code := he.getFileLines(he.Code, startLine, endLine)

	var lineNumStr string

	for i := startLine; i <= endLine; i++ {
		lineNumStr += strconv.Itoa(i) + "\n"
	}

	return map[string]interface{}{
		"Version": version.Version,
		"Style":   template.HTML("<style>" + string(css) + "</style>"),
		"Message": he.err.err.Error(),
		"File":    he.err.debug.File + ":" + strconv.Itoa(he.err.debug.Line) + ":" + strconv.Itoa(he.err.debug.Column),
		"Line":    he.err.debug.Line,
		"Lines":   lineNumStr,
		"Code":    template.HTML(code),
	}, nil
}

func (he *HtmlError) getFileLines(content []byte, start, end int) string {
	lines := bytes.Split(content, []byte("\n"))
	var result bytes.Buffer

	if start < 1 {
		start = 1
	}

	if end > len(lines) {
		end = len(lines)
	}

	for i := start - 1; i < end; i++ {
		result.WriteString(string(lines[i]) + "\n")
	}

	return result.String()
}

func (he *HtmlError) getLineNumbers(line int, maxLines int) (int, int) {
	var diff = 5

	start := line - diff
	if start < 1 {
		start = 1
	}

	end := line + diff

	if end > maxLines {
		end = maxLines
	}

	return start, end
}

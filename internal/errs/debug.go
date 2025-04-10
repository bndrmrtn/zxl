package errs

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/version"
	"github.com/fatih/color"
)

// WithDebug adds debug information to an error
func WithDebug(err error, debug *models.Debug) error {
	// Check if the error is already a DebugError
	var de DebugError
	if errors.As(err, &de) {
		if de.debug != nil {
			return de
		}
	}

	// If not, create a new DebugError
	return DebugError{err: err, debug: debug}
}

// DebugError is an error with debug information
type DebugError struct {
	err   error
	debug *models.Debug
}

// Error returns the error message with debug information
func (de DebugError) Error() string {
	return de.PrettyError(func(r io.Reader) string {
		b, _ := io.ReadAll(r)
		return string(b)
	})
}

func (de DebugError) PrettyError(pf func(r io.Reader) string) string {
	if de.debug == nil {
		return de.err.Error()
	}

	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	near := de.getNear(pf)

	return fmt.Sprintf("%s\n%s\nat %s:%d:%d\n%s", color.New(color.FgBlue, color.Bold).Sprint("Zx - ", version.Version), redBold(de.err.Error()), de.debug.File, de.debug.Line, de.debug.Column, near)
}

func (de DebugError) GetParentError() error {
	return de.err
}

func (de DebugError) getNear(pf func(r io.Reader) string) string {
	if de.debug == nil {
		return ""
	}

	var near string

	parts := strings.Split(pf(strings.NewReader(de.debug.Near)), "\n")
	maxLineNumLen := len(strconv.Itoa(de.debug.Line + len(parts) - 1))

	for i, part := range parts {
		lineNum := de.debug.Line + i
		lineNumStr := strconv.Itoa(lineNum)
		lineNumStr = strings.Repeat(" ", maxLineNumLen-len(lineNumStr)) + lineNumStr
		near += fmt.Sprintf("%s %s", color.New(color.FgHiBlack).Sprint(lineNumStr+" |"), part)
		if i < len(parts)-1 {
			near += "\n"
		}
	}

	return fmt.Sprintf("near:\n%s", near)
}

// HttpError returns the error message with debug information for HTTP
func (de DebugError) HttpError() *HtmlError {
	if de.debug == nil {
		return nil
	}

	htmlErr := NewHtmlError(&de)
	return htmlErr
}

func (de DebugError) GetLine() int {
	if de.debug == nil {
		return 1
	}

	return de.debug.Line
}

func (de DebugError) GetColumn() int {
	if de.debug == nil {
		return 1
	}

	return de.debug.Column
}

func (de DebugError) GetFile() string {
	if de.debug == nil {
		return ""
	}

	return de.debug.File
}

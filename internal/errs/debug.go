package errs

import (
	"errors"
	"fmt"
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
	if de.debug == nil {
		return de.err.Error() + "\n"
	}

	redBold := color.New(color.FgRed, color.Bold).SprintFunc()
	near := de.getNear()

	return fmt.Sprintf("%s\n%s\nat %s:%d:%d\n%s", color.New(color.FgBlue, color.Bold).Sprint("Zx - ", version.Version), redBold(de.err.Error()), de.debug.File, de.debug.Line, de.debug.Column, near)
}

func (de DebugError) getNear() string {
	if de.debug == nil {
		return ""
	}

	var near string

	parts := strings.Split(de.debug.Near, "\n")
	maxLineNumLen := len(strconv.Itoa(de.debug.Line + len(parts) - 1))

	for i, part := range parts {
		lineNum := de.debug.Line + i
		lineNumStr := strconv.Itoa(lineNum)
		lineNumStr = strings.Repeat(" ", maxLineNumLen-len(lineNumStr)) + lineNumStr
		near += fmt.Sprintf("%s | %s\n", color.New(color.FgHiBlack).Sprint(lineNumStr), part)
	}

	near = fmt.Sprintf("near:\n%s\n", color.New(color.FgHiBlack).Sprint(near))

	return near
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

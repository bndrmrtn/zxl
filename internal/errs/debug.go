package errs

import (
	"errors"
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/version"
	"github.com/fatih/color"
)

// WithDebug adds debug information to an error
func WithDebug(err error, debug *models.Debug) error {
	// Check if the error is already a DebugError
	var de DebugError
	if errors.As(err, &de) {
		return de
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

	near := ""
	if de.debug.Near != "" {
		near = fmt.Sprintf("near:\n%s\n", color.New(color.FgHiBlack).Sprint(de.debug.Near))
	}

	return fmt.Sprintf("%s\n%s\nat %s:%d:%d\n%s", color.New(color.FgBlue, color.Bold).Sprint("Zex - ", version.Version), redBold(de.err.Error()), de.debug.File, de.debug.Line, de.debug.Column, near)
}

// HttpError returns the error message with debug information for HTTP
func (de DebugError) HttpError() string {
	if de.debug == nil {
		return fmt.Sprintf("<p>%s</p>", de.err.Error())
	}

	redBold := "<span style=\"color: #fb2c36; font-weight: bold;\">" + de.err.Error() + "</span>"

	near := ""
	if de.debug.Near != "" {
		near = fmt.Sprintf("<pre style=\"background:#1e2939;padding:5px;color: #f2f2f2;\"><span style=\"font-weight:bold\">near:</span><br>%s</pre>", de.debug.Near)
	}

	return fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; margin: 20px">
			<h1 style="color: #615fff; font-weight: bold;">Zex - %s</h1>
			<p>%s</p>
			<p>at <strong>%s:%d:%d</strong></p>
			%s
		</div>
	`, version.Version, redBold, de.debug.File, de.debug.Line, de.debug.Column, near)
}

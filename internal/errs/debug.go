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

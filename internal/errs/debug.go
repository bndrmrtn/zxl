package errs

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/models"
)

// WithDebug adds debug information to an error
func WithDebug(err error, debug *models.Debug) error {
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
		return de.err.Error()
	}
	return fmt.Sprintf("%s at line %d position %d", de.err.Error(), de.debug.Line, de.debug.Column)
}

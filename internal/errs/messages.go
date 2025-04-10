package errs

import "fmt"

var (
	// Lexer errors
	ErrFailedToReadContent = fmt.Errorf("could not read content")
	SyntaxError            = fmt.Errorf("syntax error")

	// RuntimeError is the base error for all runtime issues.
	RuntimeError = fmt.Errorf("runtime error")
)

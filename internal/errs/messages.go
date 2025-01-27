package errs

import "fmt"

var (
	// Lexer errors
	ErrFailedToReadContent = fmt.Errorf("failed to read content")
	SyntaxError            = fmt.Errorf("syntax error")

	// Ast errors

	// Runtime errors
	VariableAlreadyExists = fmt.Errorf("variable already exists")
)

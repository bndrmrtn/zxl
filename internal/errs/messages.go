package errs

import "fmt"

var (
	// Lexer errors
	ErrFailedToReadContent = fmt.Errorf("failed to read content")
	SyntaxError            = fmt.Errorf("syntax error")

	// Ast errors

	// Runtime errors
	RuntimeError            = fmt.Errorf("runtime error")
	CannotRedeclareVariable = fmt.Errorf("%w: cannot redeclare variable", RuntimeError)
	VariableNotDeclared     = fmt.Errorf("%w: variable not declared", RuntimeError)
	CannotReassignConstant  = fmt.Errorf("%w: cannot reassign constant", RuntimeError)
)

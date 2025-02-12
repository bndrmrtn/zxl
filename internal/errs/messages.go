package errs

import "fmt"

var (
	// Lexer errors
	ErrFailedToReadContent = fmt.Errorf("failed to read content")
	SyntaxError            = fmt.Errorf("syntax error")

	// Ast errors

	// Runtime errors
	RuntimeError            = fmt.Errorf("runtime error")
	ValueError              = fmt.Errorf("%w: value error", RuntimeError)
	CannotRedeclareVariable = fmt.Errorf("%w: cannot redeclare variable", RuntimeError)
	VariableNotDeclared     = fmt.Errorf("%w: variable not declared", RuntimeError)
	CannotReassignConstant  = fmt.Errorf("%w: cannot reassign constant", RuntimeError)
	CannotAccessVariable    = fmt.Errorf("%w: cannot access variable", RuntimeError)
	CannotRedecareFunction  = fmt.Errorf("%w: cannot redeclare function", RuntimeError)
	CannotReUseNamespace    = fmt.Errorf("%w: cannot reuse namespace", RuntimeError)
	CannotAccessFunction    = fmt.Errorf("%w: cannot access function", RuntimeError)
	InvalidArguments        = fmt.Errorf("%w: invalid arguments", RuntimeError)
)

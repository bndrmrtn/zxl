package errs

import "fmt"

var (
	// Lexer errors
	ErrFailedToReadContent = fmt.Errorf("failed to read content")
	SyntaxError            = fmt.Errorf("syntax error")

	// Ast errors

	// Runtime errors
	RuntimeError             = fmt.Errorf("runtime error")
	ValueError               = fmt.Errorf("%w: value error", RuntimeError)
	CannotRedeclareVariable  = fmt.Errorf("%w: cannot redeclare variable", RuntimeError)
	VariableNotDeclared      = fmt.Errorf("%w: variable not declared", RuntimeError)
	CannotReassignConstant   = fmt.Errorf("%w: cannot reassign constant", RuntimeError)
	CannotAccessVariable     = fmt.Errorf("%w: cannot access variable", RuntimeError)
	CannotRedecareFunction   = fmt.Errorf("%w: cannot redeclare function", RuntimeError)
	CannotReUseNamespace     = fmt.Errorf("%w: cannot reuse namespace", RuntimeError)
	CannotAccessNamespace    = fmt.Errorf("%w: cannot access namespace", RuntimeError)
	CannotAccessFunction     = fmt.Errorf("%w: cannot access function", RuntimeError)
	InvalidArguments         = fmt.Errorf("%w: invalid arguments", RuntimeError)
	IndexOutOfRange          = fmt.Errorf("%w: index out of range", RuntimeError)
	CannotReassignDefinition = fmt.Errorf("%w: cannot reassign definition", RuntimeError)
	ThisNotInMethod          = fmt.Errorf("%w: 'this' cannot be used outside a definition", RuntimeError)
)

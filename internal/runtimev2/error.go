package runtimev2

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/models"
)

// Error creates a new error with debug information.
func Error(err error, debug *models.Debug, customText ...any) error {
	if len(customText) > 0 {
		custom := make([]string, len(customText))
		for i, s := range customText {
			custom[i] = fmt.Sprint(s)
		}

		return fmt.Errorf("%w - %w: %v", errs.RuntimeError, err, strings.Join(custom, " "))
	}

	return errs.WithDebug(fmt.Errorf("%w - %w", errs.RuntimeError, err), debug)
}

var (
	ErrInvalidValue           = fmt.Errorf("invalid value")
	ErrVariableRedeclared     = fmt.Errorf("variable already declared")
	ErrVariableNotDeclared    = fmt.Errorf("variable not declared")
	ErrConstantReassignment   = fmt.Errorf("cannot reassign constant")
	ErrVariableNotFound       = fmt.Errorf("variable not found")
	ErrFunctionRedeclared     = fmt.Errorf("function already declared")
	ErrNamespaceInUse         = fmt.Errorf("namespace already in use")
	ErrNamespaceNotFound      = fmt.Errorf("namespace not found")
	ErrFunctionNotFound       = fmt.Errorf("function not found")
	ErrInvalidArguments       = fmt.Errorf("invalid arguments")
	ErrIndexOutOfBounds       = fmt.Errorf("index out of bounds")
	ErrDefinitionReassignment = fmt.Errorf("cannot reassign this definition")
	ErrThisOutsideMethod      = fmt.Errorf("'this' cannot be used outside of a method")
	ErrUnhandledNodeType      = fmt.Errorf("unhandled node type")
	ErrInvalidIncrementTarget = fmt.Errorf("cannot increment this value")
	ErrEmptyErrorBlock        = fmt.Errorf("error block has no child nodes")
	ErrNamedInlineFunction    = fmt.Errorf("inline function cannot have a name")
	ErrExpectedExpression     = fmt.Errorf("expected expression")
	ErrExpectedBoolean        = fmt.Errorf("%w: expected boolean", ErrInvalidValue)
	ErrExpectedIterable       = fmt.Errorf("%w: expected an iterable value (e.g. array, list, or string)", ErrInvalidValue)
	ErrInvalidIndexAccess     = fmt.Errorf("%w: cannot access value: invalid index or not indexable", ErrInvalidValue)
	ErrKeyNotFound            = fmt.Errorf("key not found in array")
	ErrInvalidObject          = fmt.Errorf("invalid object")
	ErrInvalidObjectAccess    = fmt.Errorf("invalid object member access")
)

func fnErr(name string) string {
	return fmt.Sprintf("%s(...)", name)
}

func nsErr(ns, as string) string {
	if ns == as {
		return ns
	}
	return fmt.Sprintf("%s as %s", ns, as)
}

func gotErr(got any) string {
	return fmt.Sprintf("got %v", got)
}

func expectedErr(expected, got any) string {
	return fmt.Sprintf("expected %v, got %v", expected, got)
}

func valueIsNotErr(typ string) string {
	return fmt.Sprintf("value is not of type %s", typ)
}

func joinCustom(messages ...string) string {
	return strings.Join(messages, ", ")
}

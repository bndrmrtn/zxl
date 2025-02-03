package builtin

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// Variable is a runtime variable that can be used in the language
type Variable struct {
	Type  tokens.VariableType
	Value any
}

func (v *Variable) String() string {
	if v == nil {
		return "nil: this variable does not exist"
	}
	return fmt.Sprintf("%s: %v", v.Type, v.Value)
}

// FuncReturn is a return value from a function
type FuncReturn struct {
	Type  tokens.VariableType
	Value any
}

// Function is a function that can be used in the language
type Function func(args []*Variable) (*FuncReturn, error)

// Package is a package that can be used in the language
type Package interface {
	// Execute runs a function in the package
	Execute(fn string, args []*Variable) (*FuncReturn, error)
	// Access gets a constant variable from the package
	Access(variable string) (*Variable, error)
}

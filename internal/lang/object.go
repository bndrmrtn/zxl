package lang

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

// ObjType is the type of the object
type ObjType string

const (
	TString     ObjType = "<String>"
	TInt        ObjType = "<Integer>"
	TFloat      ObjType = "<Float>"
	TBool       ObjType = "<Bool>"
	TList       ObjType = "<List>"
	TNothing    ObjType = "<Nothing>"
	TNil        ObjType = "<Nil>"
	TDefinition ObjType = "<Definition>"
	TInstance   ObjType = "<Instance>"
	TFunction   ObjType = "<function>"
	TIOStream   ObjType = "<IOStream>"
	TArray      ObjType = "<Array>"
	TFnRef      ObjType = "<Function Reference>"
)

var errNotImplemented = fmt.Errorf("not implemented")

func (o ObjType) TokenType() tokens.TokenType {
	switch o {
	default:
		return tokens.Unkown
	case TString:
		return tokens.String
	case TInt, TFloat:
		return tokens.Number
	case TBool:
		return tokens.Bool
	case TList:
		return tokens.List
	case TDefinition:
		return tokens.Define
	case TFunction:
		return tokens.Function
	case TNil:
		return tokens.Nil
	case TNothing:
		return tokens.EmptyReturn
	}
}

// Object represents an object in the language
type Object interface {
	// Type returns the type of the object
	Type() ObjType
	// Name is the name of the object, or empty string if the object is inlined
	Name() string
	// Rename renames the object
	Rename(s string)
	// Value returns the value of the object
	Value() any

	// IsMutable returns if the object is mutable
	IsMutable() bool
	// Immute makes the object immutable
	Immute()

	// Method returns a method by name
	Method(name string) Method
	// Methods returns the method names of the object
	Methods() []string

	// Variable returns a variable by name
	Variable(name string) Object
	// Variables returns the variable names of the object
	Variables() []string
	// SetVariable sets a variable by name
	SetVariable(name string, object Object) error

	// Debug returns the debug information of the object
	Debug() *models.Debug

	// String returns the string representation of the object
	String() string

	// Copy returns a copy of the object
	Copy() Object
}

// Method represents a method in the language
type Method interface {
	// Args returns the arguments of the method
	Args() []string
	// Execute executes the method with the given arguments
	Execute(args []Object) (Object, error)
}

// VariadicMethod represents a variadic method in the language
type VariadicMethod interface {
	Method
	// HasVariadicArg returns if the method has a variadic argument
	HasVariadicArg() bool
	// GetVariadicArg returns the variadic argument of the method
	GetVariadicArg() string
}

// Module represents a module in the language
type Module interface {
	// Namespace returns the namespace of the module
	Namespace() string
	// Objects returns the objects of the module
	Objects() map[string]Object
	// Methods returns the methods of the module
	Methods() map[string]Method
}

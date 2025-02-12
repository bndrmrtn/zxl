package lang

import "github.com/bndrmrtn/zexlang/internal/models"

// ObjType is the type of the object
type ObjType string

const (
	TString     ObjType = "string"
	TInt        ObjType = "int"
	TFloat      ObjType = "float"
	TBool       ObjType = "bool"
	TList       ObjType = "list"
	TNothing    ObjType = "nothing"
	TNil        ObjType = "nil"
	TDefinition ObjType = "definition"
	TFunction   ObjType = "function"
)

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
	// Variable(name string) Object
	// Variables returns the variable names of the object
	// Variables() []string

	// Get returns an underlying object by key
	// Get(key string) (Object, bool)
	// Set sets an underlying object by key
	// Set(key string, value Object) bool
	// GetUnderlying returns the underlying objects
	// GetUnderlying() []string

	// Debug returns the debug information of the object
	Debug() *models.Debug
}

// Method represents a method in the language
type Method interface {
	// Args returns the arguments of the method
	Args() []string
	// Execute executes the method with the given arguments
	Execute(args []Object) (Object, error)
}

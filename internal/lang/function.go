package lang

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/models"
)

// Function represents a function method
type Function struct {
	// Method is a method of the function
	Method

	// args is a list of arguments for the function
	args []string
	// typesafeArgs is a map of argument names to their types
	typesafeArgs []TypeSafeArg
	// variadicArg is the variadic argument of the function
	variadicArg string
	// exec is the function to execute when the function is called
	exec func(args []Object) (Object, error)

	debug *models.Debug
}

type TypeSafeArg struct {
	Name string
	Type ObjType
}

// NewFunction creates a new function method
func NewFunction(exec func(args []Object) (Object, error)) *Function {
	return &Function{
		exec: exec,
	}
}

func (f *Function) WithArgs(args []string) *Function {
	f.args = args
	return f
}

func (f *Function) WithArg(arg string) *Function {
	f.args = append(f.args, arg)
	return f
}

func (f *Function) WithVariadicArg(arg string) *Function {
	f.variadicArg = arg
	return f
}

func (f *Function) WithTypeSafeArgs(typesafeArgs ...TypeSafeArg) *Function {
	var args = make([]string, len(typesafeArgs))

	for i, arg := range typesafeArgs {
		args[i] = arg.Name
	}

	f.args = args
	f.typesafeArgs = typesafeArgs
	return f
}

func (f *Function) WithDebug(debug *models.Debug) *Function {
	f.debug = debug
	return f
}

func (f *Function) Args() []string {
	return f.args
}

func (f *Function) Execute(args []Object) (Object, error) {
	if len(f.typesafeArgs) == len(f.args) {
		for i, name := range f.args {
			if args[i].Type() != f.typesafeArgs[i].Type {
				return nil, fmt.Errorf("argument %s is not of type %s", name, f.typesafeArgs[i].Type)
			}
		}
	}

	return f.exec(args)
}

func (f *Function) Debug() *models.Debug {
	return f.debug
}

func (f *Function) HasVariadicArg() bool {
	return f.variadicArg != ""
}

func (f *Function) GetVariadicArg() string {
	return f.variadicArg
}

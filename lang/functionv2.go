package lang

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/models"
)

type ExecFuncV2 func(args *ArgCtx) (any, error)

type FunctionV2 struct {
	v1Args []string
	args   []string
	exec   ExecFuncV2
	debug  *models.Debug

	// v1
	typesafeArgs []TypeSafeArg
	variadicArg  string
}

func FnV2(exec ExecFuncV2) *FunctionV2 {
	return &FunctionV2{
		exec: exec,
	}
}

func (f *FunctionV2) WithDebug(debug *models.Debug) *FunctionV2 {
	f.debug = debug
	return f
}

func (f *FunctionV2) WithTypeSafeArgs(typesafeArgs ...TypeSafeArg) *FunctionV2 {
	var args = make([]string, len(typesafeArgs))

	for i, arg := range typesafeArgs {
		args[i] = arg.Name
	}

	f.v1Args = args
	f.args = args
	f.typesafeArgs = typesafeArgs
	return f
}

func (f *FunctionV2) WithArg(arg string) *FunctionV2 {
	f.args = append(f.v1Args, arg)
	f.args = append(f.args, arg)
	return f
}

func (f *FunctionV2) WithVariadicArg(arg string) *FunctionV2 {
	f.variadicArg = arg
	f.args = append(f.args, arg)
	return f
}

func (f *FunctionV2) Adapter() *Function {
	return &Function{
		args:         f.v1Args,
		typesafeArgs: f.typesafeArgs,
		variadicArg:  f.variadicArg,
		debug:        f.debug,
		exec: func(args []Object) (Object, error) {
			ctx := NewArgCtx(f, args)
			value, err := f.exec(ctx)
			if err != nil {
				return nil, err
			}
			return FromValue(value)
		},
	}
}

type ArgCtx struct {
	fn *FunctionV2

	args []Object
}

func NewArgCtx(fn *FunctionV2, args []Object) *ArgCtx {
	return &ArgCtx{
		fn:   fn,
		args: args,
	}
}

func (ctx *ArgCtx) Arg(name string) (Object, error) {
	for i, arg := range ctx.fn.args {
		if arg == name {
			return ctx.args[i], nil
		}
	}
	return nil, fmt.Errorf("argument not found")
}

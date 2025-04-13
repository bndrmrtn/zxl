package builtin

import (
	"fmt"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/lang"
)

func setIsMethods(m map[string]lang.Method) map[string]lang.Method {
	m["isInt"] = lang.NewFunction(is(toInt)).WithArg("object")
	m["isFloat"] = lang.NewFunction(is(toFloat)).WithArg("object")
	m["isBool"] = lang.NewFunction(is(toBool)).WithArg("object")
	m["isInstanceOf"] = lang.NewFunction(isInstaceOf).
		WithArg("type").WithArg("value")

	return m
}

func is(exec lang.ExecFunc) lang.ExecFunc {
	return func(args []lang.Object) (lang.Object, error) {
		_, ok := exec(args)
		return lang.NewBool("ok", ok == nil, args[0].Debug()), nil
	}
}

func isInstaceOf(args []lang.Object) (lang.Object, error) {
	typ, ok := args[0].(*lang.Definition)
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("type must be a definition"), args[0].Debug())
	}
	inst, ok := args[1].(*lang.Instance)
	if !ok {
		return lang.NewBool("ok", false, args[1].Debug()), nil
	}

	return lang.NewBool("ok", inst.Definition().Type() == typ.Type(), args[1].Debug()), nil
}

package builtin

import (
	"github.com/bndrmrtn/zxl/lang"
)

func setIsMethods(m map[string]lang.Method) map[string]lang.Method {
	m["isInt"] = lang.NewFunction(is(toInt)).WithArg("object")
	m["isFloat"] = lang.NewFunction(is(toFloat)).WithArg("object")
	m["isBool"] = lang.NewFunction(is(toBool)).WithArg("object")
	m["isInstanceOf"] = lang.NewFunction(isInstaceOf).
		WithTypeSafeArgs(lang.TypeSafeArg{Name: "type", Type: lang.TDefinition}, lang.TypeSafeArg{Name: "object", Type: lang.TInstance})

	return m
}

func is(exec lang.ExecFunc) lang.ExecFunc {
	return func(args []lang.Object) (lang.Object, error) {
		_, ok := exec(args)
		return lang.NewBool("ok", ok == nil, args[0].Debug()), nil
	}
}

func isInstaceOf(args []lang.Object) (lang.Object, error) {
	typ := args[0].(*lang.Definition)
	inst := args[1].(*lang.Instance).Definition()

	return lang.NewBool("ok", inst.Type() == typ.Type(), args[1].Debug()), nil
}

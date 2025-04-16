package modules

import (
	"fmt"
	"strconv"

	"github.com/flarelang/flare/lang"
)

type Convert struct {
	dotEnvLoaded bool
}

func NewConvert() *Convert {
	return &Convert{}
}

func (*Convert) Namespace() string {
	return "conv"
}

func (*Convert) Objects() map[string]lang.Object {
	return nil
}

func (c *Convert) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"object": lang.NewFunction(c.fnConvertObject).WithArgs([]string{"object", "toType"}),
	}
}

func (*Convert) fnConvertObject(args []lang.Object) (lang.Object, error) {
	obj := args[0]
	to, ok := args[1].Value().(string)
	if args[1].Type() != lang.TString || !ok {
		return nil, fmt.Errorf("conversion second parameter should be a type")
	}

	switch to {
	case lang.TString.String():
		return lang.NewString("string", obj.String(), nil), nil
	case lang.TInt.String():
		if lang.ObjType(obj.Type()) == lang.TString {
			i, err := strconv.Atoi(obj.Value().(string))
			if err != nil {
				return nil, err
			}
			return lang.NewInteger("integer", i, nil), nil
		}
	}

	return nil, fmt.Errorf("invalid or unsupported conversion")
}

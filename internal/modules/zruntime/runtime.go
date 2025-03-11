package zruntime

import (
	"fmt"
	"os"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type Runtime struct{}

func New() *Runtime {
	return &Runtime{}
}

func (*Runtime) Namespace() string {
	return "runtime"
}

func (h *Runtime) Objects() map[string]lang.Object {
	return nil
}

func (h *Runtime) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"copy": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			return args[0].Copy(), nil
		}).WithArg("object"),
		"exit": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			os.Exit(args[0].Value().(int))
			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "code", Type: lang.TInt}),
		"$deepObjDump": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			fmt.Printf("%#v\n", args[0])
			return nil, nil
		}).WithArg("object"),
	}
}

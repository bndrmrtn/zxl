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
		"copy": lang.NewFunction([]string{"object"}, func(args []lang.Object) (lang.Object, error) {
			return args[0].Copy(), nil
		}, nil),
		"exit": lang.NewFunction([]string{"code"}, func(args []lang.Object) (lang.Object, error) {
			if args[0].Type() != lang.TInt {
				return nil, fmt.Errorf("expected int, got %s", args[0].Type())
			}
			os.Exit(int(args[0].Value().(int64)))
			return nil, nil
		}, nil),
	}
}

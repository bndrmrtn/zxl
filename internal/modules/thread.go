package modules

import (
	"fmt"
	"time"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type Thread struct{}

func NewThreadModule() *Thread {
	return &Thread{}
}

func (*Thread) Namespace() string {
	return "thread"
}

func (*Thread) Objects() map[string]lang.Object {
	return nil
}

func (*Thread) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"sleep": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			duration := args[0].Value().(int)
			time.Sleep(time.Duration(duration) * time.Millisecond)
			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{
			Name: "duration",
			Type: lang.TInt,
		}),
		"spawn": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			fn := args[0].Value().(lang.Method)
			if len(fn.Args()) != 0 {
				return nil, fmt.Errorf("spawn expected no arguments, got %d", len(fn.Args()))
			}

			go func() {
				fn.Execute(nil)
			}()

			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "method", Type: lang.TFnRef}),
	}
}

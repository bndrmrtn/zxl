package thread

import (
	"fmt"
	"time"

	"github.com/flarelang/flare/lang"
)

type Thread struct {
	id uint
}

func New() *Thread {
	return &Thread{}
}

func (*Thread) Namespace() string {
	return "thread"
}

func (*Thread) Objects() map[string]lang.Object {
	return nil
}

func (t *Thread) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"sleep": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			duration := args[0].Value().(int)
			time.Sleep(time.Duration(duration) * time.Millisecond)
			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{
			Name: "duration",
			Type: lang.TInt,
		}),
		"portal": lang.NewFunction(func(variadicArgs []lang.Object) (lang.Object, error) {
			portalBufferSize := 10
			args := variadicArgs[0].Value().([]lang.Object)

			if len(args) == 1 && args[0].Type() == lang.TInt {
				portalBufferSize = args[0].Value().(int)
			} else if len(args) > 0 {
				return nil, fmt.Errorf("portal expected one argument of type int, got %d", len(args))
			}

			portal := NewPortal(t.id, portalBufferSize)
			t.id++
			return portal, nil
		}).WithVariadicArg("buffer"),
		"spawn": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			fn := args[0].(*lang.Fn).Fn
			if len(fn.Args()) != 0 {
				return nil, fmt.Errorf("spawn handler must have 0 arguments, got %d", len(fn.Args()))
			}

			go func() {
				_, _ = fn.Execute(nil)
			}()

			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "method", Type: lang.TFnRef}),
		"spawner": lang.NewFunction(func(variadicArgs []lang.Object) (lang.Object, error) {
			spawnerBufferSize := 10
			args := variadicArgs[0].Value().([]lang.Object)

			if len(args) == 1 && args[0].Type() == lang.TInt {
				spawnerBufferSize = args[0].Value().(int)
			} else if len(args) > 0 {
				return nil, fmt.Errorf("spawner expected one argument of type int, got %d", len(args))
			}

			spawner := NewSpawner(spawnerBufferSize)
			t.id++
			return spawner, nil
		}).WithVariadicArg("buffer"),
	}
}

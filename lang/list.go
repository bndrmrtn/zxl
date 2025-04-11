package lang

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/bndrmrtn/zxl/internal/models"
)

// List represents a list object
type List struct {
	Base

	value  []Object
	length int
}

func GetListValue(l Object) []Object {
	if l.Type() != TList {
		panic("Object is not a List")
	}
	return l.(*List).value
}

// NewList creates a new list object
func NewList(name string, li []Object, debug *models.Debug) Object {
	return &List{
		Base: Base{
			name:    name,
			debug:   debug,
			mutable: true,
		},
		value:  li,
		length: len(li),
	}
}

func (l *List) Type() ObjType {
	return TList
}

func (l *List) Name() string {
	return l.name
}

func (l *List) Value() any {
	return l.value
}

func (l *List) Method(name string) Method {
	switch name {
	default:
		return nil
	case "append":
		return NewFunction(func(args []Object) (Object, error) {
			l.value = append(l.value, args[0].Copy())
			l.length++
			return nil, nil
		}).WithArgs([]string{"item"}).WithDebug(l.debug)
	case "contains":
		return NewFunction(func(args []Object) (Object, error) {
			if l.length < 1 {
				return NewBool("contains", false, l.debug), nil
			}

			search := args[0].Value()

			for _, v := range l.value {
				if reflect.DeepEqual(v.Value(), search) {
					return NewBool("contains", true, l.debug), nil
				}
			}

			return NewBool("contains", false, l.debug), nil
		}).WithDebug(l.debug).WithArgs([]string{"item"})
	case "filter":
		return NewFunction(func(args []Object) (Object, error) {
			fn := args[0].(*Fn).Fn
			fnArgs := fn.Args()

			if len(fn.Args()) != 1 {
				return nil, fmt.Errorf("filter function must have one argument")
			}

			var filtered []Object

			for _, v := range l.value {
				arg := v.Copy()
				arg.Rename(fnArgs[0])

				obj, err := fn.Execute([]Object{arg})
				if err != nil {
					return nil, err
				}

				if obj.Type() != TBool {
					return nil, fmt.Errorf("filter function must return a boolean")
				}

				if obj.Value().(bool) {
					filtered = append(filtered, v)
				}
			}

			return NewList(l.name, filtered, l.debug), nil
		}).WithTypeSafeArgs(TypeSafeArg{"filterFunc", TFnRef})
	case "insert":
		return NewFunction(func(args []Object) (Object, error) {
			if args[0].Type() != TInt {
				return nil, fmt.Errorf("index must be an integer")
			}

			index := args[0].Value().(int)
			item := args[1].Copy()

			if index < 0 || index > l.length {
				return nil, fmt.Errorf("index out of range")
			}

			l.value = append(l.value[:index], append([]Object{item}, l.value[index:]...)...)
			l.length++

			return nil, nil
		}).WithArgs([]string{"index", "item"}).WithDebug(l.debug)
	}
}

func (l *List) Methods() []string {
	return []string{"append", "contains", "filter", "insert"}
}

func (l *List) Variable(variable string) Object {
	switch variable {
	default:
		return nil
	case "length":
		return NewInteger("length", l.length, l.debug)
	case "$addr":
		return addr(l)
	}
}

func (l *List) Variables() []string {
	return []string{"length", "$addr"}
}

func (l *List) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (l *List) String() string {
	sb := strings.Builder{}
	defer sb.Reset()

	for i, v := range l.value {
		sb.WriteString(v.String())
		if i < len(l.value)-1 {
			sb.WriteString(", ")
		}
	}

	return fmt.Sprintf("[%s]", sb.String())
}

func (l *List) Copy() Object {
	var value []Object

	for _, v := range l.value {
		value = append(value, v.Copy())
	}

	return NewList(l.name, value, l.debug)
}

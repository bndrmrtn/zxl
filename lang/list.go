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
	}
}

func (l *List) Methods() []string {
	return []string{"append", "contains"}
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

package lang

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/models"
)

// List represents a list object
type List struct {
	Base

	value []Object
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
			name:  name,
			debug: debug,
		},
		value: li,
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
	return nil
}

func (l *List) Methods() []string {
	return []string{"append"}
}

func (l *List) Variable(_ string) Object {
	return nil
}

func (l *List) Variables() []string {
	return []string{"length"}
}

func (l *List) SetVariable(_ string, _ Object) error {
	return notImplemented
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

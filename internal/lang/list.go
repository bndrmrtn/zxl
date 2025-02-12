package lang

import "github.com/bndrmrtn/zexlang/internal/models"

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
	return []string{"append", "len"}
}

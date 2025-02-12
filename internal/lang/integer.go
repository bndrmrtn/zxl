package lang

import (
	"github.com/bndrmrtn/zexlang/internal/models"
)

type Integer struct {
	Base
	value int
}

func NewInteger(name string, i int, debug *models.Debug) Object {
	return &Integer{
		Base:  NewBase(name, debug),
		value: i,
	}
}

func (i *Integer) Type() ObjType {
	return TInt
}

func (i *Integer) Name() string {
	return i.name
}

func (i *Integer) Value() any {
	return i.value
}

func (i *Integer) Method(name string) Method {
	return nil
}

func (i *Integer) Methods() []string {
	return []string{"toString"}
}

package lang

import (
	"github.com/bndrmrtn/zxl/internal/models"
)

type Bool struct {
	Base
	value bool
}

func NewBool(name string, b bool, debug *models.Debug) Object {
	return &Bool{
		Base:  NewBase(name, debug),
		value: b,
	}
}

func (*Bool) Type() ObjType {
	return TBool
}

func (b *Bool) Name() string {
	return b.name
}

func (b *Bool) Value() any {
	return b.value
}

func (*Bool) Method(name string) Method {
	return nil
}

func (*Bool) Methods() []string {
	return nil
}

func (b *Bool) Variable(name string) Object {
	switch name {
	default:
		return nil
	case "$addr":
		return addr(b)
	}
}

func (*Bool) Variables() []string {
	return []string{"$addr"}
}

func (*Bool) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (b *Bool) String() string {
	if b.value {
		return "true"
	}
	return "false"
}

func (b *Bool) Copy() Object {
	return NewBool(b.name, b.value, b.debug)
}

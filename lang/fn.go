package lang

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/models"
)

type Fn struct {
	Base
	Fn Method
}

func NewFn(name string, debug *models.Debug, fn Method) Object {
	return &Fn{
		Base: NewBase(name, debug),
		Fn:   fn,
	}
}

func (f *Fn) Type() ObjType {
	return TFnRef
}

func (f *Fn) Name() string {
	return f.name
}

func (f *Fn) Value() any {
	return f
}

func (f *Fn) Method(name string) Method {
	return nil
}

func (f *Fn) Methods() []string {
	return nil
}

func (f *Fn) Variable(name string) Object {
	switch name {
	default:
		return nil
	case "$addr":
		return addr(f.Fn)
	}
}

func (f *Fn) Variables() []string {
	return []string{"$addr"}
}

func (f *Fn) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (f *Fn) String() string {
	return fmt.Sprintf("Fn(reference:%s)", f.name)
}

func (f *Fn) Copy() Object {
	return f
}

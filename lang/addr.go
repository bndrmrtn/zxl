package lang

import (
	"fmt"
)

func Addr(s Object) Object {
	return addr(s)
}

func addr(s Object) Object {
	return &ObjAddr{
		Base: Base{
			name: "$addr",
		},
		value: fmt.Sprintf("%p", s),
		obj:   &s,
	}
}

type ObjAddr struct {
	Base
	value string
	obj   *Object
}

func (*ObjAddr) Type() ObjType {
	return TAddr
}

func (o *ObjAddr) Value() any {
	return o
}

func (o *ObjAddr) Method(name string) Method {
	switch name {
	case "value":
		return NewFunction(func(args []Object) (Object, error) {
			return *o.obj, nil
		})
	}
	return nil
}

func (*ObjAddr) Methods() []string {
	return []string{"value"}
}

func (*ObjAddr) Variable(_ string) Object {
	return nil
}

func (*ObjAddr) Variables() []string {
	return nil
}

func (*ObjAddr) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (o *ObjAddr) String() string {
	return o.value
}

func (o *ObjAddr) Copy() Object {
	return o
}

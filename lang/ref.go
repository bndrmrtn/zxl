package lang

import (
	"github.com/bndrmrtn/zxl/internal/models"
)

type Ref struct {
	Base
	value Object
}

func NewRef(name string, debug *models.Debug, r Object) Object {
	return &Ref{
		Base:  NewBase(name, debug),
		value: r,
	}
}

func (r *Ref) Type() ObjType {
	return r.value.Type()
}

func (r *Ref) Name() string {
	return r.value.Name()
}

func (r *Ref) Value() any {
	return r.value.Value()
}

func (r *Ref) Method(name string) Method {
	if name == "*assign" {
		return NewFunction(func(args []Object) (Object, error) {
			r.value = args[0]
			return nil, nil
		}).WithArg("value")
	}

	return r.value.Method(name)
}

func (r *Ref) Methods() []string {
	return append(r.value.Methods(), "*assign")
}

func (r *Ref) Variable(name string) Object {
	if name == "$value" {
		return r.value
	}

	if name == "$addr" {
		return addr(r)
	}

	return r.value.Variable(name)
}

func (r *Ref) Variables() []string {
	return append(r.value.Variables(), "$value")
}

func (r *Ref) SetVariable(name string, value Object) error {
	return r.value.SetVariable(name, value)
}

func (r *Ref) String() string {
	return r.value.String()
}

func (r *Ref) Copy() Object {
	// prevent ref from being copied to become a reference to the same object
	return r
}

package lang

import (
	"strconv"

	"github.com/bndrmrtn/zxl/internal/models"
)

// Float represents a float object
type Float struct {
	Base
	value float64
}

// NewFloat creates a new float object
func NewFloat(name string, f float64, debug *models.Debug) Object {
	return &Float{
		Base:  NewBase(name, debug),
		value: f,
	}
}

func (f *Float) Type() ObjType {
	return TFloat
}

func (f *Float) Name() string {
	return f.name
}

func (f *Float) Value() any {
	return f.value
}

func (f *Float) Method(name string) Method {
	return nil
}

func (f *Float) Methods() []string {
	return []string{"toString"}
}

func (f *Float) Variable(name string) Object {
	switch name {
	default:
		return nil
	case "$addr":
		return addr(f)
	}
}

func (f *Float) Variables() []string {
	return []string{"$addr"}
}

func (f *Float) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (f *Float) String() string {
	return strconv.FormatFloat(f.value, 'f', -1, 64)
}

func (f *Float) Copy() Object {
	return NewFloat(f.name, f.value, f.debug)
}

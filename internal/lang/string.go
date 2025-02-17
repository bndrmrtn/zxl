package lang

import (
	"github.com/bndrmrtn/zxl/internal/models"
)

// String represents a string object
type String struct {
	Base

	value  string
	length int
}

// NewString creates a new string object
func NewString(name, s string, debug *models.Debug) Object {
	return &String{
		Base:   NewBase(name, debug),
		value:  s,
		length: -1,
	}
}

func (s *String) Type() ObjType {
	return TString
}

func (s *String) Value() any {
	return s.value
}

func (s *String) Method(name string) Method {
	return nil
}

func (s *String) Methods() []string {
	return []string{"split"}
}

func (s *String) Variable(variable string) Object {
	switch variable {
	default:
		return nil
	case "length":
		if s.length == -1 {
			s.length = len(s.value)
		}
		return NewInteger("length", s.length, s.debug)
	case "$addr":
		return addr(s)
	}
}

func (s *String) Variables() []string {
	return []string{"length", "$addr"}
}

func (s *String) SetVariable(_ string, _ Object) error {
	return errNotImplemented
}

func (s *String) String() string {
	return s.value
}

func (s *String) Copy() Object {
	return NewString(s.name, s.value, s.debug)
}

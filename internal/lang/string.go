package lang

import (
	"github.com/bndrmrtn/zexlang/internal/models"
)

// String represents a string object
type String struct {
	Base

	value string
}

// NewString creates a new string object
func NewString(name, s string, debug *models.Debug) Object {
	return &String{
		Base:  NewBase(name, debug),
		value: s,
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

func (s *String) Variable(_ string) Object {
	return nil
}

func (s *String) Variables() []string {
	return []string{"length"}
}

func (s *String) SetVariable(_ string, _ Object) error {
	return notImplemented
}

func (s *String) String() string {
	return s.value
}

func (s *String) Copy() Object {
	return NewString(s.name, s.value, s.debug)
}

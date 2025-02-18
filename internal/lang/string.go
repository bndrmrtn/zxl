package lang

import (
	"fmt"
	"strings"

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
	switch name {
	case "split":
		return NewFunction(nil, func(args []Object) (Object, error) {
			if len(args) != 1 {
				return nil, fmt.Errorf("wrong number of arguments for split: expected 1, got %d", len(args))
			}

			separator := args[0]

			if separator.Type() != TString {
				return nil, fmt.Errorf("separator must be a string")
			}

			value := strings.Split(s.value, separator.Value().(string))
			var parts []Object
			for _, part := range value {
				parts = append(parts, NewString("part", part, s.debug))
			}

			return NewList("split", parts, s.debug), nil
		}, nil)
	case "trim":
		return NewFunction(nil, func(args []Object) (Object, error) {
			return NewString("trim", strings.TrimSpace(s.value), s.debug), nil
		}, nil)
	default:
		return nil
	}
}

func (s *String) Methods() []string {
	return []string{"split", "trim"}
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

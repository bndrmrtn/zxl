package lang

import (
	"strings"

	"github.com/bndrmrtn/flare/internal/models"
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
		return NewFunction(func(args []Object) (Object, error) {
			separator := args[0]
			value := strings.Split(s.value, separator.Value().(string))
			var parts []Object
			for _, part := range value {
				parts = append(parts, NewString("part", part, s.debug))
			}
			return NewList("split", parts, s.debug), nil
		}).WithTypeSafeArgs(TypeSafeArg{"separator", TString})
	case "trim":
		return NewFunction(func(args []Object) (Object, error) {
			return NewString("trim", strings.TrimSpace(s.value), s.debug), nil
		})
	case "lower":
		return NewFunction(func(args []Object) (Object, error) {
			return NewString("lower", strings.ToLower(s.value), s.debug), nil
		})
	case "upper":
		return NewFunction(func(args []Object) (Object, error) {
			return NewString("upper", strings.ToUpper(s.value), s.debug), nil
		})
	case "replace":
		return NewFunction(func(args []Object) (Object, error) {
			oldStr, newStr := args[0], args[1]
			return NewString("replace", strings.ReplaceAll(s.value, oldStr.Value().(string), newStr.Value().(string)), s.debug), nil
		}).WithTypeSafeArgs(TypeSafeArg{"old", TString}, TypeSafeArg{"new", TString})
	case "contains":
		return NewFunction(func(args []Object) (Object, error) {
			substr := args[0]
			return NewBool("contains", strings.Contains(s.value, substr.Value().(string)), s.debug), nil
		}).WithTypeSafeArgs(TypeSafeArg{"substring", TString})
	case "startsWith":
		return NewFunction(func(args []Object) (Object, error) {
			prefix := args[0]
			return NewBool("startswith", strings.HasPrefix(s.value, prefix.Value().(string)), s.debug), nil
		}).WithTypeSafeArgs(TypeSafeArg{"prefix", TString})
	case "endsWith":
		return NewFunction(func(args []Object) (Object, error) {
			suffix := args[0]
			return NewBool("endswith", strings.HasSuffix(s.value, suffix.Value().(string)), s.debug), nil
		}).WithTypeSafeArgs(TypeSafeArg{"suffix", TString})
	default:
		return nil
	}
}

func (s *String) Methods() []string {
	return []string{"split", "trim", "lower", "upper", "replace", "contains", "startsWith", "endsWith"}
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

package lang

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zxl/internal/models"
)

type Array struct {
	Base

	// Keys are the keys of the array.
	Keys []Object
	// Map is the map of the array.
	Map map[any]Object
}

func NewArray(name string, debug *models.Debug, keys []Object, values []Object) Object {
	if len(keys) != len(values) {
		panic("keys and values must have the same length")
	}
	array := &Array{
		Base: NewBase(name, debug),
		Keys: keys,
		Map:  make(map[any]Object, len(keys)),
	}
	for i, key := range keys {
		array.Map[key] = values[i]
	}
	return array
}

func (a *Array) Type() ObjType {
	return TArray
}

func (a *Array) Value() any {
	return a
}

func (a *Array) Method(name string) Method {
	switch name {
	case "values":
		return NewFunction(func(args []Object) (Object, error) {
			var values = make([]Object, len(a.Keys))

			for i, key := range a.Keys {
				values[i] = a.Map[key]
			}

			return NewList("values", values, a.debug), nil
		}).WithDebug(a.debug)
	case "$bind":
		return NewFunction(func(args []Object) (Object, error) {
			key := args[0]
			value := args[1]

			if _, ok := a.Map[key]; !ok {
				a.Keys = append(a.Keys, key)
			}

			if a.Map == nil {
				a.Map = make(map[any]Object)
			}

			a.Map[key] = value
			return nil, nil
		}).WithDebug(a.debug).WithArgs([]string{"key", "value"})
	default:
		return nil
	}
}

func (a *Array) Methods() []string {
	return []string{"values", "$bind"}
}

func (a *Array) Variable(variable string) Object {
	switch variable {
	default:
		acc, _ := a.Access(variable)
		return acc
	case "$addr":
		return addr(a)
	case "keys":
		return NewList("keys", a.Keys, a.debug)
	}
}

func (a *Array) Variables() []string {
	return []string{"$addr"}
}

func (a *Array) SetVariable(name string, value Object) error {
	key, ok := a.realKey(name)
	if !ok {
		return fmt.Errorf("key %v not found", name)
	}

	a.Map[key] = value
	return nil
}

func (a *Array) String() string {
	sb := strings.Builder{}

	for i, key := range a.Keys {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(key.String())
		sb.WriteString(": ")
		sb.WriteString(a.Map[key].String())
	}

	return "array{" + sb.String() + "}"
}

func (a *Array) Copy() Object {
	values := make([]Object, len(a.Keys))
	for i, key := range a.Keys {
		values[i] = a.Map[key].Copy()
	}
	return NewArray(a.name, a.debug, a.Keys, values)
}

func (a *Array) Access(access any) (Object, bool) {
	if access == nil {
		return nil, false
	}

	for _, key := range a.Keys {
		if key.Value() == access {
			return a.Map[key], true
		}
	}

	return nil, false
}
func (a *Array) realKey(access any) (Object, bool) {
	if access == nil {
		return nil, false
	}

	for _, key := range a.Keys {
		if key.Value() == access {
			return key, true
		}
	}

	return nil, false
}

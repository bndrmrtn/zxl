package builtin

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/state"
	"github.com/bndrmrtn/zxl/lang"
)

type State struct {
	lang.Base
	name string

	state state.State
}

func NewState(name string, debug *models.Debug, state state.State) *State {
	return &State{
		Base:  lang.NewBase(name, debug),
		name:  name,
		state: state,
	}
}

func (s *State) Type() lang.ObjType {
	return lang.TInstance
}

func (s *State) TypeString() string {
	return "zx.stateClient"
}

func (s *State) Value() any {
	return s
}

func (s *State) Method(name string) lang.Method {
	switch name {
	default:
		return nil
	case "get":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].Value().(string)

			value, ok := s.state.Get(key)
			if !ok {
				return lang.NilObject, nil
			}

			return value, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString})
	case "set":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0]
			if key.Type() != lang.TString {
				return nil, errs.WithDebug(fmt.Errorf("name must be a string, got %T", key.Type()), key.Debug())
			}
			value := args[1]
			ok := s.state.Set(key.Value().(string), value)
			return lang.NewBool("ok", ok, value.Debug()), nil
		}).WithArg("name").WithArg("value")
	}
}

func (s *State) Methods() []string {
	return []string{"get", "set"}
}

func (s *State) Variable(variable string) lang.Object {
	return nil
}

func (s *State) Variables() []string {
	return nil
}

func (s *State) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (s *State) String() string {
	return fmt.Sprintf("<zx.stateClient %s>", s.name)
}

func (s *State) Copy() lang.Object {
	return s
}

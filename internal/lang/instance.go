package lang

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/models"
)

// Executer represents a node executer in the runtime
type Executer interface {
	GetVariable(name string) (Object, error)
	Variables() []string
	AssignVariable(name string, object Object) error
	GetMethod(name string) (Method, error)
	Execute(nodes []*models.Node) (Object, error)
	GetNew() Executer
}

type Instance struct {
	Base

	def *Definition
	ex  Executer
}

func NewDefinitionInstance(def *Definition, ex Executer) *Instance {
	return &Instance{
		Base: def.Base,
		def:  def,
		ex:   ex,
	}
}

func (*Instance) Type() ObjType {
	return TInstance
}

func (i *Instance) TypeString() string {
	return i.def.defName
}

func (i *Instance) Value() any {
	return i
}

func (i *Instance) Method(name string) Method {
	if name == "$init" {
		construct, err := i.ex.GetMethod("construct")
		if err != nil {
			return NewFunction(nil, func(args []Object) (Object, error) {
				return i, nil
			}, i.debug)
		}

		return NewFunction(construct.Args(), func(args []Object) (Object, error) {
			construct, err := i.ex.GetMethod("construct")
			if err != nil {
				return i, nil
			}

			_, err = construct.Execute(args)
			return i, err
		}, i.debug)
	}

	m, _ := i.ex.GetMethod(name)
	return m
}

func (i *Instance) Methods() []string {
	return nil
}

func (i *Instance) Variable(variable string) Object {
	if variable == "$addr" {
		return addr(i)
	}

	obj, _ := i.ex.GetVariable(variable)
	return obj
}

func (i *Instance) Variables() []string {
	return i.ex.Variables()
}

func (i *Instance) SetVariable(name string, value Object) error {
	return i.ex.AssignVariable(name, value)
}

func (i *Instance) String() string {
	str := addr(i).String()

	method, err := i.ex.GetMethod("string")
	if err == nil && len(method.Args()) == 0 {
		val, err := method.Execute(nil)
		if err == nil {
			str = val.String()
		}
	}

	return fmt.Sprintf("<%s %s>", i.def.defName, str)
}

func (i *Instance) Copy() Object {
	return NewDefinitionInstance(i.def, i.ex)
}

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
			return NewFunction(func(args []Object) (Object, error) {
				return i, nil
			}).WithDebug(i.debug)
		}

		return NewFunction(func(args []Object) (Object, error) {
			construct, err := i.ex.GetMethod("construct")
			if err != nil {
				return i, nil
			}

			_, err = construct.Execute(args)
			return i, err
		}).WithArgs(construct.Args()).WithDebug(i.debug)
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
	method, err := i.ex.GetMethod("string")
	if err == nil && len(method.Args()) == 0 {
		val, err := method.Execute(nil)
		if err == nil {
			return val.String()
		}
	}

	return fmt.Sprintf("<%s %s>", i.def.defName, addr(i).String())
}

func (i *Instance) Copy() Object {
	return NewDefinitionInstance(i.def, i.ex)
}

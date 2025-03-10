package lang

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zxl/internal/models"
)

// Executer represents a node executer in the runtime
type Executer interface {
	GetVariable(name string) (Object, error)
	AssignVariable(name string, object Object) error

	GetMethod(name string) (Method, error)
	Copy() Executer
}

type Definition struct {
	Base

	defName string
	Ex      Executer
}

func NewDefinition(defName, name string, debug *models.Debug, ex Executer) Object {
	return &Definition{
		defName: strings.TrimLeft(defName, "."),
		Base:    NewBase(name, debug),
		Ex:      ex,
	}
}

func (d *Definition) Type() ObjType {
	return TDefinition
}

func (d *Definition) TypeString() string {
	return d.defName
}

func (d *Definition) Value() any {
	return d
}

func (d *Definition) Method(name string) Method {
	if name == "$init" {
		construct, err := d.Ex.GetMethod("construct")
		if err != nil {
			return NewFunction(nil, func(args []Object) (Object, error) {
				return d.Copy(), nil
			}, d.debug)
		}

		return NewFunction(construct.Args(), func(args []Object) (Object, error) {
			obj := d.Copy().(*Definition)

			construct, err := d.Ex.GetMethod("construct")
			if err != nil {
				return obj, nil
			}

			_, err = construct.Execute(args)
			return obj, err
		}, d.debug)
	}

	m, _ := d.Ex.GetMethod(name)
	return m
}

func (d *Definition) Methods() []string {
	return []string{"split"}
}

func (d *Definition) Variable(variable string) Object {
	if variable == "$addr" {
		return addr(d)
	}

	obj, _ := d.Ex.GetVariable(variable)
	return obj
}

func (d *Definition) Variables() []string {
	return []string{"length"}
}

func (d *Definition) SetVariable(name string, value Object) error {
	return d.Ex.AssignVariable(name, value)
}

func (d *Definition) String() string {
	method, err := d.Ex.GetMethod("string")
	if err == nil && len(method.Args()) == 0 {
		val, err := method.Execute(nil)
		if err == nil {
			return val.String()
		}
	}

	return fmt.Sprintf("<%s %s>", d.defName, addr(d))
}

func (d *Definition) Copy() Object {
	return NewDefinition(d.defName, d.name, d.debug, d.Ex.Copy())
}

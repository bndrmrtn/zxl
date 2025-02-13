package lang

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/models"
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
		defName: defName,
		Base:    NewBase(name, debug),
		Ex:      ex,
	}
}

func (d *Definition) Type() ObjType {
	return TDefinition
}

func (d *Definition) Value() any {
	return d
}

func (d *Definition) Method(name string) Method {
	if name == "$init" {
		return NewFunction([]string{}, func(args []Object) (Object, error) {
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

func (d *Definition) Variable(_ string) Object {
	return nil
}

func (d *Definition) Variables() []string {
	return []string{"length"}
}

func (d *Definition) SetVariable(_ string, _ Object) error {
	return notImplemented
}

func (d *Definition) String() string {
	return fmt.Sprintf("<%s %v>", d.defName, &d)
}

func (d *Definition) Copy() Object {
	return NewDefinition(d.defName, d.name, d.debug, d.Ex.Copy())
}

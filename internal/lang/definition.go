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
	Execute(nodes []*models.Node) (Object, error)
	GetNew() Executer
}

type Definition struct {
	Base

	defName string

	ex    Executer
	nodes []*models.Node
}

func NewDefinition(defName, name string, debug *models.Debug, nodes []*models.Node, ex Executer) *Definition {
	return &Definition{
		Base:    NewBase(name, debug),
		defName: strings.TrimLeft(defName, "."),
		nodes:   nodes,
		ex:      ex,
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
	return nil
}

func (d *Definition) Methods() []string {
	return nil
}

func (d *Definition) Variable(variable string) Object {
	return nil
}

func (d *Definition) Variables() []string {
	return nil
}

func (d *Definition) SetVariable(name string, value Object) error {
	return nil
}

func (d *Definition) String() string {
	return fmt.Sprintf("<%s>", d.defName)
}

func (d *Definition) Copy() Object {
	return d
}

func (d *Definition) NewInstance() Object {
	newDef := *&d
	exec := newDef.ex.GetNew()

	exec.Execute(newDef.nodes)
	return NewDefinitionInstance(newDef, exec)
}

package lang

import "github.com/bndrmrtn/zxl/internal/models"

type Nil struct {
	Object

	name  string
	debug *models.Debug
}

var NilObject = Nil{}

func NewNil(name string, debug *models.Debug) Object {
	return Nil{
		name:  name,
		debug: debug,
	}
}

func (n Nil) Type() ObjType {
	return TNil
}

func (n Nil) Name() string {
	return ""
}

func (n Nil) Value() any {
	return nil
}

func (n Nil) String() string {
	return "<Nil>"
}

func (n Nil) Debug() *models.Debug {
	return n.debug
}

func (n Nil) IsMutable() bool {
	return true
}

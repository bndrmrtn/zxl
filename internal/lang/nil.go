package lang

type Nil struct {
	Object
}

var NilObject = Nil{}

func (n Nil) Type() ObjType {
	return TNil
}

func (n Nil) Name() string {
	return ""
}

func (n Nil) Value() any {
	return nil
}

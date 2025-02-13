package lang

type EmptyReturn struct {
	Object
}

func NewEmptyReturn() Object {
	return &EmptyReturn{}
}

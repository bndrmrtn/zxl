package tokens

// VariableType represents the type of the variable
type VariableType int

const (
	NilVariable VariableType = iota
	IntVariable
	FloatVariable
	StringVariable
	BoolVariable
	FunctionCallVariable
	ExpressionVariable
	ReferenceVariable
	InlineValue
	DefinitionBlock
)

func (v VariableType) String() string {
	switch v {
	case NilVariable:
		return "nil"
	case IntVariable:
		return "int"
	case FloatVariable:
		return "float"
	case StringVariable:
		return "string"
	case BoolVariable:
		return "bool"
	case FunctionCallVariable:
		return "call:fn()"
	case ExpressionVariable:
		return "expression"
	case ReferenceVariable:
		return ":ref:"
	default:
		return "unknown"
	}
}

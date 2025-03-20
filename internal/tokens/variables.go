package tokens

// VariableType represents the type of the variable
type VariableType int

const (
	NilVariable VariableType = iota
	IntVariable
	FloatVariable
	StringVariable
	TemplateVariable
	BoolVariable
	FunctionCallVariable
	FunctionVariable
	ExpressionVariable
	ReferenceVariable
	InlineValue
	DefinitionBlock
	DefinitionReference
	EmptyReturnValue
	ListVariable
	ArrayVariable
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
		return "functionCall"
	case ExpressionVariable:
		return "expression"
	case ReferenceVariable:
		return "reference(variable)"
	case DefinitionReference:
		return "reference(definition)"
	case DefinitionBlock:
		return "definition"
	case EmptyReturnValue:
		return "return(empty)"
	case ListVariable:
		return "list"
	default:
		return "unknown"
	}
}

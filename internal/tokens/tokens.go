package tokens

// TokenType represents the type of the token
type TokenType int

const (
	Unkown TokenType = iota

	Use TokenType = iota + 10
	As
	From

	Let TokenType = iota + 100
	Const
	Define
	Return
	New
	This // Reference for blocks
	Identifier
	FuncCall

	Number TokenType = iota + 1000
	String
	Bool
	Nil
	Function
	Definiton

	Addition TokenType = iota + 10000
	Subtraction
	Multiplication
	Division
	Modulo
	Increment
	Decrement
	Assign
	Equation
	NotEquation
	Greater
	Less
	GreaterOrEqual
	LessOrEqual
	And
	Or
	Not

	LeftParenthesis TokenType = iota + 100000
	RightParenthesis
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket

	Comma TokenType = iota + 1000000
	Dot
	Semicolon
	Colon
	At

	If TokenType = iota + 10000000
	ElseIf
	Else
)

func (t TokenType) IsOperator() bool {
	return t >= Addition && t <= Not
}

func (t TokenType) ToVariableType() VariableType {
	switch t {
	case Number:
		return IntVariable
	case String:
		return StringVariable
	case Bool:
		return BoolVariable
	default:
		return NilVariable
	}
}

func (t TokenType) String() string {
	switch t {
	case Unkown:
		return "unkown"
	case Use:
		return "use"
	case As:
		return "as"
	case From:
		return "from"
	case Let:
		return "let"
	case Const:
		return "const"
	case Define:
		return "define"
	case Return:
		return "return"
	case New:
		return "new"
	case This:
		return "this"
	case Identifier:
		return "identifier"
	case FuncCall:
		return "funcCall"
	case Number:
		return "number"
	case String:
		return "string"
	case Bool:
		return "bool"
	case Nil:
		return "nil"
	case Function:
		return "function"
	case Definiton:
		return "definiton"
	case Addition:
		return "+"
	case Subtraction:
		return "-"
	case Multiplication:
		return "*"
	case Division:
		return "/"
	case Modulo:
		return "%"
	case Increment:
		return "++"
	case Decrement:
		return "--"
	case Assign:
		return "="
	case Equation:
		return "=="
	case NotEquation:
		return "!="
	case Greater:
		return ">"
	case Less:
		return "<"
	case GreaterOrEqual:
		return ">="
	case LessOrEqual:
		return "<="
	case And:
		return "&&"
	case Or:
		return "||"
	case Not:
		return "!"
	case LeftParenthesis:
		return "("
	case RightParenthesis:
		return ")"
	case LeftBrace:
		return "{"
	case RightBrace:
		return "}"
	case LeftBracket:
		return "["
	case RightBracket:
		return "]"
	case Comma:
		return ","
	case Dot:
		return "."
	case Semicolon:
		return ";"
	case Colon:
		return ":"
	case At:
		return "@"
	case If:
		return "if"
	case ElseIf:
		return "else if"
	case Else:
		return "else"
	default:
		return "unkown"
	}
}

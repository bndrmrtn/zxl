package tokens

// TokenType represents the type of the token
type TokenType int

const (
	Unkown TokenType = iota

	Namespace TokenType = iota + 10
	Use
	As
	From

	Let TokenType = iota + 100
	Const
	Define
	Return
	This // Reference for blocks
	Identifier
	FuncCall
	InlineFunction
	EmptyReturn
	For
	In
	While
	Throw
	Spin

	Number TokenType = iota + 1000
	String
	Bool
	Nil
	Function
	FuncArg
	List
	ListValue
	Definiton
	TemplateLiteral
	Array
	ArrayKeyValuePair

	Addition TokenType = iota + 10000
	Subtraction
	Multiplication
	Power
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
	Arrow

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
	Then
	Else

	NewLine TokenType = iota + 100000000
	SingleLineComment
	MultiLineComment
	WhiteSpace
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
	case InlineFunction:
		return "inlineFunction"
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
	case Else:
		return "else"
	case In:
		return "in"
	case NewLine:
		return "new line"
	case SingleLineComment:
		return "single line comment"
	case MultiLineComment:
		return "multi line comment"
	case WhiteSpace:
		return "white space"
	case For:
		return "for"
	case Throw:
		return "throw"
	case List:
		return "list"
	case ListValue:
		return "list value"
	case TemplateLiteral:
		return "<></>"
	case Array:
		return "array"
	case Spin:
		return "spin"
	default:
		return "unkown"
	}
}

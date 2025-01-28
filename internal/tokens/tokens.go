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

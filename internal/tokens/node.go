package tokens

// NodeType represents the type of the node
type NodeType int

const (
	UnkownNode NodeType = iota

	UseNode NodeType = iota + 10

	VariableDeclarationNode NodeType = iota + 100
	VariableAssignmentNode
	FunctionDeclarationNode
	FunctionCallNode
	FunctionReturnNode

	DefinitionNode NodeType = iota + 1000

	IfStatementNode NodeType = iota + 10000
	ElseIfStatementNode
	ElseStatementNode
	WhileStatementNode
	ForStatementNode
)

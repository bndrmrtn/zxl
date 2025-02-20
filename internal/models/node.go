package models

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/tokens"
)

// Node represents a node in the AST
type Node struct {
	// Type represents the type of the node
	Type tokens.TokenType `yaml:"type"`
	// VariableType represents the type of the variable
	VariableType tokens.VariableType `yaml:"variableType"`

	// Content represents the content of the node
	Content string `yaml:"content,omitempty"`
	// Value represents the value of the node
	Value any `yaml:"value,omitempty"`

	// Reference is a bool to determine if the node is a reference
	Reference bool `yaml:"reference,omitempty"`
	// Children represents the children of the node
	Children []*Node `yaml:"children,omitempty"`
	// Args represents the arguments of the node
	Args []*Node `yaml:"args,omitempty"`
	// ObjectAccessors represents the object accessors of the node like x[0]
	ObjectAccessors []*Node `yaml:"objectAccessors,omitempty"`

	// Map is a key value store for the node
	Map map[string]any `yaml:"map,omitempty"`

	// Flags are custom operations for a node
	Flags []string `yaml:"flags,omitempty"`

	// Debug represents the debug information of the node
	Debug *Debug `yaml:"debug,omitempty"`
}

func (n *Node) String() string {
	if n == nil {
		return "nil: this node does not exist"
	}

	return fmt.Sprintf("type: %d, variable: %d, content: %s, value: %v, ref: %v, children: %v, args: %v",
		n.Type,
		n.VariableType,
		n.Content,
		n.Value,
		n.Reference,
		n.Children,
		n.Args,
	)
}

func (n *Node) ValueString() string {
	switch n.VariableType {
	case tokens.StringVariable:
		return fmt.Sprintf("%q", n.Value)
	case tokens.BoolVariable:
		return fmt.Sprintf("%t", n.Value)
	case tokens.IntVariable:
		return fmt.Sprintf("%d", n.Value)
	case tokens.FloatVariable:
		return fmt.Sprintf("%f", n.Value)
	}

	return fmt.Sprintf("%v", n.Value)
}

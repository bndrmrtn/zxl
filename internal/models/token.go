package models

import "github.com/flarelang/flare/internal/tokens"

// Token represents a token in the lexer
type Token struct {
	// TokenType is the type of the token
	Type tokens.TokenType `yaml:"type"`
	// Value is the value of the token
	Value string `yaml:"value"`

	// Map is a map of key value pairs
	Map map[string]any `yaml:"map,omitempty"`
	// Debug is a debug object
	Debug *Debug `yaml:"-"`
}

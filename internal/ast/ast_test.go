package ast

import (
	"fmt"
	"testing"

	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func Test_AstVarDelc(t *testing.T) {
	ts := []*models.Token{
		{
			Type:  tokens.Let,
			Value: "let",
		},
		{
			Type:  tokens.Identifier,
			Value: "a",
		},
		{
			Type:  tokens.Assign,
			Value: "=",
		},
		{
			Type:  tokens.Number,
			Value: "25",
		},
		{
			Type:  tokens.Semicolon,
			Value: ";",
		},
	}

	b := NewBuilder()
	nodes, err := b.Build(ts)
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range nodes {
		fmt.Println(node.String())
	}
}

func Test_AstFuncDecl(t *testing.T) {
	ts := []*models.Token{
		{
			Type:  tokens.Function,
			Value: "fn",
		},
		{
			Type:  tokens.Identifier,
			Value: "hi",
		},
		{
			Type:  tokens.LeftParenthesis,
			Value: "(",
		},
		{
			Type:  tokens.Identifier,
			Value: "username",
		},
		{
			Type:  tokens.RightParenthesis,
			Value: ")",
		},
		{
			Type:  tokens.LeftBrace,
			Value: "{",
		},
		{
			Type:  tokens.Let,
			Value: "let",
		},
		{
			Type:  tokens.Identifier,
			Value: "a",
		},
		{
			Type:  tokens.Assign,
			Value: "=",
		},
		{
			Type:  tokens.Number,
			Value: "25",
		},
		{
			Type:  tokens.Semicolon,
			Value: ";",
		},
		{
			Type:  tokens.RightBrace,
			Value: "}",
		},
	}

	b := NewBuilder()
	nodes, err := b.Build(ts)
	if err != nil {
		t.Fatal(err)
	}

	for _, node := range nodes {
		fmt.Println(node.String())
	}
}

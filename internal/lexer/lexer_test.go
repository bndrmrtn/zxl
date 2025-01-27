package lexer

import (
	"fmt"
	"strings"
	"testing"
)

func TestLexer_Strings_Comments(t *testing.T) {
	testCode := `// this is a comment
	"this is a string"

	"this is an invalid string'`

	lx := New("test")
	_, err := lx.Parse(strings.NewReader(testCode))
	if err == nil {
		t.Errorf("error should not be nil, but got %v", err)
	}
}

func TestLexer_VariableDeclaration(t *testing.T) {
	testCode := `let x = "hello";`

	lx := New("test")
	ts, err := lx.Parse(strings.NewReader(testCode))
	if err != nil {
		t.Fatal(err)
	}

	if len(ts) != 5 {
		t.Errorf("expected 5 tokens, but got %d", len(ts))
	}

	for _, token := range ts {
		fmt.Println(token)
	}
}

func TestLexer_UsePackage(t *testing.T) {
	testCode := `use http;`

	lx := New("test")
	ts, err := lx.Parse(strings.NewReader(testCode))
	if err != nil {
		t.Fatal(err)
	}

	if len(ts) != 3 {
		t.Errorf("expected 3 tokens, but got %d", len(ts))
	}

	for _, token := range ts {
		fmt.Println(token)
	}
}

func TestLexer_FuncDecl(t *testing.T) {
	testCode := `fn main(user) {
		let x = 4;
	}`

	lx := New("test")
	ts, err := lx.Parse(strings.NewReader(testCode))
	if err != nil {
		t.Fatal(err)
	}

	if len(ts) != 16 {
		t.Errorf("expected 16 tokens, but got %d", len(ts))
	}

	for _, token := range ts {
		fmt.Println(token)
	}
}

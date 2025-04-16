package ast

import (
	"strings"
	"testing"

	"github.com/flarelang/flare/internal/lexer"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
	"github.com/stretchr/testify/assert"
)

func build(t *testing.T, s string) []*models.Node {
	ts, err := lexer.New("<test>").Parse(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}

	b := NewBuilder()
	nodes, err := b.Build(ts)
	if err != nil {
		t.Fatal(err)
	}

	return nodes
}

func Test_AstVarDelc(t *testing.T) {
	nodes := build(t, "let x = 5;")

	assert.Equal(t, 1, len(nodes), "must create on variable node")
	variable := nodes[0]

	assert.Equal(t, "x", variable.Content, "variable name must be x")
	assert.Equal(t, 5, variable.Value, "variable value must be int 5")
}

func Test_FnDecl(t *testing.T) {
	nodes := build(t, `
		fn main() {
			return 0;
		}
	`)

	assert.Equal(t, 1, len(nodes), "must create on function node")
	fn := nodes[0]

	assert.Equal(t, "main", fn.Content, "fn name must be main")
	assert.Equal(t, 0, fn.Args, "args must be zero")
	assert.Equal(t, 1, len(fn.Children), "must create return children")
}

func Test_Define(t *testing.T) {
	nodes := build(t, `
		define MyApp {
			let id = 0;

			fn construct(id) {
				this.id = id;
			}
		}
	`)

	assert.Equal(t, 1, len(nodes), "must create one definition node")
	definition := nodes[0]

	assert.Equal(t, "MyApp", definition.Content, "definition name must be MyApp")
	assert.Equal(t, 2, len(definition.Children), "definition must have 2 children")
	variable := definition.Children[0]
	fn := definition.Children[1]

	assert.Equal(t, "id", variable.Content, "variable name should be id")
	assert.Equal(t, 0, variable.Value, "variable value must be 0")

	assert.Equal(t, "construct", fn.Content, "fn name should be construct")
	assert.Equal(t, 1, len(fn.Args), "fn args must be 1")
	assert.Equal(t, 1, len(fn.Children), "fn body children must be 1")
	arg := fn.Args[0]
	body := fn.Children[0]

	assert.Equal(t, "id", arg.Content, "fn arg must be id")

	assert.Equal(t, tokens.Assign, body.Type, "fn content must be an assignment")
	assert.Equal(t, "this.id", body.Content, "fn must assign to this.id")
	assert.Equal(t, "id", body.Value, "fn must assign id")
}

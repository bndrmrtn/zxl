package source

import (
	"bytes"
	"embed"

	"github.com/bndrmrtn/flare/internal/ast"
	"github.com/bndrmrtn/flare/internal/lexer"
	"github.com/bndrmrtn/flare/internal/models"
)

//go:embed data/*.fl
var sourceFiles embed.FS

func Get() (map[string][]*models.Node, error) {
	files, err := sourceFiles.ReadDir("data")
	if err != nil {
		return nil, err
	}

	var m = make(map[string][]*models.Node)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := sourceFiles.ReadFile("data/" + file.Name())
		if err != nil {
			return nil, err
		}

		// Tokenize the source code with lexer
		fileName := "@zx/" + file.Name()
		lx := lexer.New(fileName)
		ts, err := lx.Parse(bytes.NewReader(content))
		if err != nil {
			return nil, err
		}

		// Build the abstract syntax tree from tokens
		builder := ast.NewBuilder()
		nodes, err := builder.Build(ts)
		if err != nil {
			return nil, err
		}

		m[fileName] = nodes
	}

	return m, nil
}

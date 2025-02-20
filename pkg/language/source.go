package language

import (
	"bytes"
	"embed"

	"github.com/bndrmrtn/zxl/internal/ast"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/internal/runtimev2"
)

//go:embed source/*.zx
var sourceFiles embed.FS

// executeSourceFiles executes all source files in the "source" directory.
func (ir *Interpreter) ExecuteSourceFiles(run *runtimev2.Runtime) error {
	files, err := sourceFiles.ReadDir("source")
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		content, err := sourceFiles.ReadFile("source/" + file.Name())
		if err != nil {
			return err
		}

		// Tokenize the source code with lexer
		lx := lexer.New("@zx/" + file.Name())
		ts, err := lx.Parse(bytes.NewReader(content))
		if err != nil {
			return err
		}

		// Build the abstract syntax tree from tokens
		builder := ast.NewBuilder()
		nodes, err := builder.Build(ts)
		if err != nil {
			return err
		}

		if _, err = run.Execute(nodes); err != nil {
			return err
		}
	}

	return nil
}

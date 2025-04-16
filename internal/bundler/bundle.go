package bundler

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flarelang/flare/internal/ast"
	"github.com/flarelang/flare/internal/lexer"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/version"
)

type File struct {
	File    string
	Content []byte
	Nodes   []*models.Node
	Sum     string
}

type Bundle struct {
	Version string
	Files   []*File
}

func NewBundle(root string) (*Bundle, error) {
	files := make([]*File, 0)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		filePath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(file)
		file.Close()
		if err != nil {
			return err
		}

		lx := lexer.New(filepath.Base(path))
		tokens, err := lx.Parse(bytes.NewReader(data))
		if err != nil {
			return nil
		}

		ast := ast.NewBuilder()
		nodes, err := ast.Build(tokens)
		if err != nil {
			return err
		}

		files = append(files, &File{
			File:    filePath,
			Content: data,
			Nodes:   nodes,
			Sum:     fmt.Sprintf("%x", sha256.Sum256(data)),
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &Bundle{
		Files:   files,
		Version: version.Version,
	}, nil
}

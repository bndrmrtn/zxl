package runtimev2

import (
	"bytes"
	"io"
	"os"
	"path/filepath"

	"github.com/bndrmrtn/zxl/internal/ast"
	"github.com/bndrmrtn/zxl/internal/cache"
	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/internal/models"
)

func (r *Runtime) importer(filename string, dg *models.Debug) (lang.Object, error) {
	var root = ""
	if dg != nil {
		root = dg.File
	}

	path := filepath.Join(filepath.Dir(root), filename)
	path = filepath.Clean(path)

	file, err := os.Open(path)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}

	builder := ast.NewBuilder()
	nodes, ok := cache.Get(path, b)
	if !ok {
		lx := lexer.New(path)
		ts, err := lx.Parse(bytes.NewReader(b))
		if err != nil {
			return nil, errs.WithDebug(err, dg)
		}

		nodes, err = builder.Build(ts)
		if err != nil {
			return nil, errs.WithDebug(err, dg)
		}

		if len(nodes) == 0 {
			return nil, nil
		}

	}

	cache.Store(path, b, nodes)

	ret, err := r.Execute(nodes)
	if err != nil {
		return nil, errs.WithDebug(err, dg)
	}

	return ret, nil
}

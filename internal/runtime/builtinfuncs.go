package runtime

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/ast"
	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/cache"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/lexer"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// runFuncEval runs the eval function
func (e *Executer) runFuncEval(debug *models.Debug, args []*builtin.Variable) (*builtin.FuncReturn, error) {
	if len(args) != 1 {
		return nil, errs.WithDebug(fmt.Errorf("eval function takes only one argument"), debug)
	}
	if args[0].Type != tokens.StringVariable {
		return nil, errs.WithDebug(fmt.Errorf("eval function takes only string argument"), debug)
	}

	lx := lexer.New(debug.File + ":eval")
	ts, err := lx.Parse(strings.NewReader(args[0].Value.(string)))
	if err != nil {
		return nil, errs.WithDebug(err, debug)
	}

	builder := ast.NewBuilder()
	nodes, err := builder.Build(ts)
	if err != nil {
		return nil, errs.WithDebug(err, debug)
	}

	ex := NewExecuter(ExecuterScopeGlobal, e.runtime, e).WithName(e.name + ".{eval}")
	return ex.Execute(nodes)
}

// runFuncImport runs the import function
func (e *Executer) runFuncImport(debug *models.Debug, args []*builtin.Variable) (*builtin.FuncReturn, error) {
	if len(args) < 1 {
		return nil, errs.WithDebug(fmt.Errorf("import function takes minimum one argument"), debug)
	}

	for _, arg := range args {
		if arg.Type != tokens.StringVariable {
			return nil, errs.WithDebug(fmt.Errorf("import function takes only string arguments"), debug)
		}
	}

	builder := ast.NewBuilder()

	for _, arg := range args {
		path := filepath.Join(filepath.Dir(debug.File), arg.Value.(string))
		path = filepath.Clean(path)

		file, err := os.Open(path)
		if err != nil {
			return nil, errs.WithDebug(err, debug)
		}
		defer file.Close()

		b, err := io.ReadAll(file)
		if err != nil {
			return nil, errs.WithDebug(err, debug)
		}

		nodes, ok := cache.Get(path, b)
		if !ok {
			lx := lexer.New(path)
			ts, err := lx.Parse(bytes.NewReader(b))
			if err != nil {
				return nil, errs.WithDebug(err, debug)
			}

			nodes, err = builder.Build(ts)
			if err != nil {
				return nil, errs.WithDebug(err, debug)
			}

			if len(nodes) == 0 {
				continue
			}
		}

		cache.Store(path, b, nodes)

		var ns string
		if nodes[0].Type == tokens.Namespace {
			ns = nodes[0].Content
			nodes = nodes[1:]
		}

		ret, err := e.runtime.Exec(ExecuterScopeGlobal, e, ns, nodes)
		if err != nil {
			return nil, errs.WithDebug(err, debug)
		}

		if ret != nil {
			return ret, nil
		}
	}

	return nil, nil
}

// runFuncRef runs the ref function
func (e *Executer) runFuncRef(token *models.Node) (*builtin.FuncReturn, error) {
	if len(token.Args) != 1 {
		return nil, errs.WithDebug(fmt.Errorf("ref function takes only one argument"), token.Debug)
	}

	n, ok := e.vars[token.Args[0].Content]
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("ref: variable not found"), token.Debug)
	}

	if n.VariableType == tokens.ReferenceVariable {
		if n.Value == nil {
			return nil, errs.WithDebug(fmt.Errorf("ref: referencing a nil value"), token.Debug)
		}

		return &builtin.FuncReturn{
			Type:  tokens.ReferenceVariable,
			Value: n.Value.(string),
		}, nil
	}

	if n.VariableType == tokens.DefinitionReference {
		ex, ok := n.Value.(*Executer)
		if !ok {
			return nil, errs.WithDebug(fmt.Errorf("ref: invalid reference"), token.Debug)
		}
		return &builtin.FuncReturn{
			Type:  tokens.DefinitionReference,
			Value: ex.name,
		}, nil
	}

	return &builtin.FuncReturn{
		Type:  tokens.BoolVariable,
		Value: false,
	}, nil
}

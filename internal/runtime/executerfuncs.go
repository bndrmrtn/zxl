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

func (e *Executer) runFuncEval(debug *models.Debug, args []*builtin.Variable) ([]*builtin.FuncReturn, error) {
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

	ex := NewExecuter(ExecuterScopeGlobal, e.runtime, e)
	return ex.Execute(nodes)
}

func (e *Executer) runFuncImport(debug *models.Debug, args []*builtin.Variable) ([]*builtin.FuncReturn, error) {
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

		nodes, ok := cache.Get(b)
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

		cache.Store(b, nodes)

		var ns string
		if nodes[0].Type == tokens.Namespace {
			ns = nodes[0].Content
			nodes = nodes[1:]
		}

		ret, err := e.runtime.Exec(ExecuterScopeGlobal, e, ns, nodes)
		if err != nil {
			return nil, errs.WithDebug(err, debug)
		}

		if len(ret) != 0 {
			return ret, nil
		}
	}

	return nil, nil
}

func (e *Executer) runFuncRef(token *models.Node) ([]*builtin.FuncReturn, error) {
	if len(token.Args) != 1 {
		return nil, errs.WithDebug(fmt.Errorf("ref function takes only one argument"), token.Debug)
	}

	n, err := e.GetVariableValue(token.Args[0].Content)
	if err != nil {
		return nil, err
	}

	if n.VariableType == tokens.ReferenceVariable {
		if token.Value == nil {
			return nil, errs.WithDebug(fmt.Errorf("ref: referencing a nil value"), token.Debug)
		}

		node, err := e.GetVariableValue(token.Value.(string))
		if err != nil {
			return nil, err
		}
		return e.runFuncRef(node)
	}

	return []*builtin.FuncReturn{
		{
			Type:  tokens.StringVariable,
			Value: fmt.Sprintf("{Type: %s, Value: %v}", n.VariableType, n.Value),
		},
	}, nil
}

package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// Executer is the runtime executer
type Executer struct {
	scope ExecuterScope

	runtime *Runtime
	parent  *Executer

	fns    map[string]*models.Node
	vars   map[string]*models.Node
	blocks map[string]*models.Node
}

// NewExecuter creates a new runtime executer
func NewExecuter(scope ExecuterScope, r *Runtime, parent *Executer) *Executer {
	return &Executer{
		scope:   scope,
		runtime: r,
		parent:  parent,
		fns:     make(map[string]*models.Node),
		vars:    make(map[string]*models.Node),
		blocks:  make(map[string]*models.Node),
	}
}

// Bind binds a variable to the executer
func (e *Executer) Bind(variable *models.Node) {
	e.vars[variable.Content] = variable
}

// Execute executes the given nodes
func (e *Executer) Execute(ts []*models.Node) ([]*builtin.FuncReturn, error) {
	for _, token := range ts {
		switch token.Type {
		// Handle function declarations
		case tokens.Function:
			e.fns[token.Content] = token
		// Handle block definitions
		case tokens.Define:
			e.blocks[token.Content] = token
		// Handle variable declarations
		case tokens.Let, tokens.Const:
			err := e.handleLetConst(token)
			if err != nil {
				return nil, err
			}
			break
		// Handle function calls
		case tokens.FuncCall:
			_, err := e.executeFn(token)
			if err != nil {
				return nil, err
			}
			break
		case tokens.Assign:
			err := e.handleAssignment(token)
			if err != nil {
				return nil, err
			}
			break
		case tokens.If:
			ret, err := e.handleIf(token)
			if err != nil {
				return nil, err
			}
			if ret != nil {
				return ret, nil
			}
			break
		case tokens.Return:
			return e.handleReturn(token)
		}
	}
	return nil, nil
}

func (e *Executer) ExecuteFn(name string, args []*builtin.Variable) ([]*builtin.FuncReturn, error) {
	fn, ok := e.fns[name]
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("function %v not found", name), nil)
	}

	ex := NewExecuter(ExecuterScopeFunction, e.runtime, e)

	for i, arg := range args {
		ex.Bind(&models.Node{
			Type:         tokens.Let,
			VariableType: arg.Type,
			Content:      fn.Args[i].Content,
			Value:        arg.Value,
		})
	}

	return ex.Execute(fn.Children)
}

package runtime

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/builtin"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// Executer is the runtime executer
type Executer struct {
	// name of the executer
	name string

	// scope of the executer
	scope ExecuterScope

	// runtime reference
	runtime *Runtime
	// parent executer
	parent *Executer

	// fns is the map of functions
	fns map[string]*models.Node
	// vars is the map of variables
	vars map[string]*models.Node
	// blocks is the map of blocks
	blocks map[string]*models.Node

	packages map[string]string
}

// NewExecuter creates a new runtime executer
func NewExecuter(scope ExecuterScope, r *Runtime, parent *Executer) *Executer {
	return &Executer{
		scope:    scope,
		runtime:  r,
		parent:   parent,
		fns:      make(map[string]*models.Node),
		vars:     make(map[string]*models.Node),
		blocks:   make(map[string]*models.Node),
		packages: make(map[string]string),
	}
}

func (e *Executer) WithName(name string) *Executer {
	e.name = strings.TrimPrefix(e.name+"."+name, ".")
	return e
}

// Bind binds a variable to the executer
func (e *Executer) Bind(variable *models.Node) {
	name := variable.Content
	v, ok := e.vars[name]
	if !ok {
		if e.scope == ExecuterScopeBlock {
			if e.parent != nil {
				e.parent.Bind(variable)
				return
			}
		}
		e.vars[variable.Content] = variable
		return
	}
	v.Value = variable.Value
	v.Type = variable.Type
	v.VariableType = tokens.ReferenceVariable
	v.Children = variable.Children
	v.Args = variable.Args
	v.Reference = variable.Reference
}

// Execute executes the given nodes
func (e *Executer) Execute(ts []*models.Node) (*builtin.FuncReturn, error) {
	for _, token := range ts {
		switch token.Type {
		// Handle package imports
		case tokens.Use:
			using := token.Content
			as := token.Value.(string)
			if _, ok := e.packages[as]; ok {
				return nil, errs.WithDebug(fmt.Errorf("%w: package '%v' already imported", errs.RuntimeError, as), token.Debug)
			}
			e.packages[as] = using
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
		// Handle function calls
		case tokens.FuncCall:
			_, err := e.executeFn(token)
			if err != nil {
				return nil, err
			}
		case tokens.Assign:
			err := e.handleAssignment(token)
			if err != nil {
				return nil, err
			}
		case tokens.If:
			ret, err := e.handleIf(token)
			if err != nil {
				return nil, err
			}
			if ret != nil {
				return ret, nil
			}
		case tokens.Return, tokens.EmptyReturn:
			return e.handleReturn(token)
		case tokens.While:
			ret, err := e.handleWhile(token)
			if err != nil {
				return nil, err
			}
			if ret != nil {
				return ret, nil
			}
		}
	}
	return nil, nil
}

func (e *Executer) ExecuteFn(name string, args []*builtin.Variable) (*builtin.FuncReturn, error) {
	fn, ok := e.fns[name]
	if !ok {
		if (e.scope == ExecuterScopeBlock || e.scope == ExecuterScopeFunction) && e.parent != nil {
			return e.parent.ExecuteFn(name, args)
		}
		return nil, errs.WithDebug(fmt.Errorf("%w: function '%v' not found", errs.RuntimeError, name), nil)
	}

	ex := NewExecuter(ExecuterScopeFunction, e.runtime, e).WithName(e.name + ".{" + name + "}")

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

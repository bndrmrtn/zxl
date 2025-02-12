package runtimev2

import (
	"fmt"
	"strings"
	"sync"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/lang"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// Executer represents a node executer in the runtime
type Executer struct {
	// name of the executer
	name string

	// scope of the executer
	scope ExecuterScope

	// runtime reference
	runtime *Runtime
	// parent executer
	parent *Executer

	// functions is the map of functions
	functions map[string]lang.Method
	// objects is the map of objects
	objects map[string]lang.Object
	// usednamespaces is the list of used namespaces
	usedNamespaces map[string]string

	// mu is the mutex
	mu sync.RWMutex
}

// NewExecuter creates a new runtime executer
func NewExecuter(scope ExecuterScope, r *Runtime, parent *Executer) *Executer {
	return &Executer{
		scope:          scope,
		runtime:        r,
		parent:         parent,
		functions:      make(map[string]lang.Method),
		objects:        make(map[string]lang.Object),
		usedNamespaces: make(map[string]string),
	}
}

// WithName sets the name of the executer
func (e *Executer) WithName(name string) *Executer {
	e.name = strings.TrimLeft(e.name+"."+name, ".")
	return e
}

// BindObject binds an object to the executer
func (e *Executer) BindObject(name string, object lang.Object) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.objects[name] = object
}

// GetMethod gets a method by name
func (e *Executer) GetMethod(name string) (lang.Method, error) {
	ex := e

	ex.mu.RLock()
	defer ex.mu.RUnlock()

	if method, ok := e.functions[name]; ok {
		return method, nil
	}

	if ex.parent != nil && (ex.scope == ExecuterScopeBlock || ex.scope == ExecuterScopeFunction || ex.scope == ExecuterScopeDefinition) {
		return ex.parent.GetMethod(name)
	}

	if fn, ok := ex.runtime.functions[name]; ok {
		return fn, nil
	}

	return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
}

// GetVariable gets a variable by name
func (e *Executer) GetVariable(name string) (lang.Object, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if obj, ok := e.objects[name]; ok {
		return obj, nil
	}

	if e.parent != nil && e.scope == ExecuterScopeBlock {
		return e.parent.GetVariable(name)
	}

	return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
}

func (e *Executer) AssignVariable(name string, object lang.Object) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if _, ok := e.objects[name]; !ok {
		if e.parent != nil && e.scope == ExecuterScopeBlock {
			return e.parent.AssignVariable(name, object)
		}
		return errs.WithDebug(fmt.Errorf("%w: '%s'", errs.VariableNotDeclared, name), nil)
	}

	if !e.objects[name].IsMutable() {
		return errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotReassignConstant, name), nil)
	}

	e.objects[name] = object
	return nil
}

// Execute executes the given nodes
func (e *Executer) Execute(nodes []*models.Node) (lang.Object, error) {
	for _, node := range nodes {
		ret, err := e.executeNode(node)
		if err != nil {
			return nil, err
		}
		if ret != nil {
			return ret, nil
		}
	}
	return nil, nil
}

// executeNode executes a node
func (e *Executer) executeNode(node *models.Node) (lang.Object, error) {
	switch node.Type {
	case tokens.Use:
		using := node.Content
		as := node.Value.(string)
		if _, ok := e.usedNamespaces[as]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s' as '%s'", errs.CannotReUseNamespace, using, as), node.Debug)
		}
		e.usedNamespaces[as] = using
	case tokens.Function:
		name, method, err := e.createMethodFromNode(node)
		if err != nil {
			return nil, err
		}
		if _, ok := e.functions[name]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s(...)'", errs.CannotRedecareFunction, name), node.Debug)
		}

		e.mu.Lock()
		e.functions[name] = method
		e.mu.Unlock()
	case tokens.Let, tokens.Const:
		name, object, err := e.createObjectFromNode(node)
		if err != nil {
			return nil, err
		}
		if _, ok := e.objects[name]; ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotRedeclareVariable, name), node.Debug)
		}
		e.mu.Lock()
		e.objects[name] = object
		e.mu.Unlock()
	case tokens.FuncCall:
		_, err := e.callFunctionFromNode(node)
		if err != nil {
			return nil, errs.WithDebug(err, node.Debug)
		}
	case tokens.Assign:
		err := e.assignObjectFromNode(node)
		return nil, errs.WithDebug(err, node.Debug)
	}
	return nil, nil
}

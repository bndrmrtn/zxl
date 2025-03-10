package runtimev2

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bndrmrtn/zxl/internal/ast"
	"github.com/bndrmrtn/zxl/internal/cache"
	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/internal/models"
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

// BindMethod binds a method to the executer
func (e *Executer) BindMethod(name string, fn lang.Method) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.functions[name] = fn
}

func (e *Executer) AssignVariable(name string, object lang.Object) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		first := names[0]
		middle := names[1 : len(names)-1]
		last := names[len(names)-1]

		if first == "this" {
			if e.parent != nil && e.parent.scope == ExecuterScopeDefinition {
				return e.parent.AssignVariable(strings.Join(append(middle, last), "."), object)
			}
			return errs.WithDebug(fmt.Errorf("%w: '%s'", errs.ThisNotInMethod, name), nil)
		}

		variable := strings.Join(append([]string{first}, middle...), ".")

		e.mu.Unlock()
		v, err := e.GetVariable(variable)
		e.mu.Lock()

		if err == nil {
			return v.SetVariable(last, object)
		}
	}

	obj, ok := e.objects[name]
	if !ok {
		if e.parent != nil && e.scope == ExecuterScopeBlock {
			return e.parent.AssignVariable(name, object)
		}
		return errs.WithDebug(fmt.Errorf("%w: '%s'", errs.VariableNotDeclared, name), nil)
	}

	if !obj.IsMutable() {
		return errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotReassignConstant, name), nil)
	}

	if obj.Type() == lang.TDefinition {
		return errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotReassignDefinition, name), nil)
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

// Copy creates a copy of the executer
func (e *Executer) Copy() lang.Executer {
	newExec := &Executer{
		name:           e.name,
		scope:          e.scope,
		runtime:        e.runtime,
		functions:      e.functions,
		objects:        make(map[string]lang.Object),
		mu:             sync.RWMutex{},
		usedNamespaces: e.usedNamespaces,
		parent:         e.parent,
	}

	for key, obj := range e.objects {
		if obj == nil || obj.Type() == lang.TNil {
			newExec.objects[key] = obj
		} else {
			newExec.objects[key] = obj.Copy()
		}
	}

	return newExec
}

func (e *Executer) GetNamespaceExecuter(namespace string) (*Executer, error) {
	nses := e.getUsedNamespaces()
	if nses == nil {
		return nil, errs.WithDebug(errs.CannotAccessNamespace, nil)
	}

	ns, ok := nses[namespace]
	if !ok {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessNamespace, namespace), nil)
	}

	exec, err := e.runtime.GetNamespaceExecuter(ns)
	if err != nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: %w", errs.CannotAccessNamespace, err), nil)
	}

	return exec, nil
}

func (e *Executer) getUsedNamespaces() map[string]string {
	ns := e.usedNamespaces
	if ns != nil {
		return ns
	}

	return e.parent.getUsedNamespaces()
}

func (e *Executer) LoadFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	builder := ast.NewBuilder()
	nodes, ok := cache.Get(path, b)
	if !ok {
		lx := lexer.New(path)
		ts, err := lx.Parse(bytes.NewReader(b))
		if err != nil {
			return err
		}

		nodes, err = builder.Build(ts)
		if err != nil {
			return err
		}

		if len(nodes) == 0 {
			return nil
		}

	}

	cache.Store(path, b, nodes)

	_, err = e.runtime.Execute(nodes)
	return err
}

package runtimev2

import (
	"bytes"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bndrmrtn/flare/internal/ast"
	"github.com/bndrmrtn/flare/internal/cache"
	"github.com/bndrmrtn/flare/internal/lexer"
	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/lang"
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

func (e *Executer) GetNew() lang.Executer {
	return &Executer{
		scope:          e.scope,
		runtime:        e.runtime,
		parent:         e.parent,
		functions:      make(map[string]lang.Method),
		objects:        make(map[string]lang.Object),
		usedNamespaces: e.usedNamespaces,
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
	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		first := names[0]
		middle := names[1 : len(names)-1]
		last := names[len(names)-1]

		if first == "this" {
			def := e.isInsideDefinitionMethod(e)
			if def != nil {
				return def.AssignVariable(strings.Join(append(middle, last), "."), object)
			}
			return Error(ErrThisOutsideMethod, nil, name)
		}

		variable := strings.Join(append([]string{first}, middle...), ".")

		v, err := e.GetVariable(variable)

		if err == nil {
			return v.SetVariable(last, object)
		}
	}

	e.mu.RLock()
	obj, ok := e.objects[name]
	e.mu.RUnlock()
	if !ok {
		if e.parent != nil && e.scope == ExecuterScopeBlock {
			return e.parent.AssignVariable(name, object)
		}
		return Error(ErrVariableNotDeclared, nil, name)
	}

	if !obj.IsMutable() {
		return Error(ErrConstantReassignment, nil, name)
	}

	if obj.Type() == lang.TDefinition {
		return Error(ErrDefinitionReassignment, nil, name)
	}

	e.mu.Lock()
	e.objects[name] = object
	e.mu.Unlock()
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
	ex := &Executer{
		name:           e.name,
		scope:          e.scope,
		runtime:        e.runtime,
		parent:         e.parent,
		functions:      e.functions,
		mu:             sync.RWMutex{},
		usedNamespaces: e.usedNamespaces,
	}

	ex.objects = make(map[string]lang.Object)
	for k, v := range e.objects {
		ex.objects[k] = v
	}

	return ex
}

func (e *Executer) GetNamespaceExecuter(namespace string) (*Executer, error) {
	nses := e.getUsedNamespaces()
	if nses == nil {
		return nil, Error(ErrNamespaceNotFound, nil, namespace)
	}

	ns, ok := nses[namespace]
	if !ok {
		return nil, Error(ErrNamespaceNotFound, nil, namespace)
	}

	exec, err := e.runtime.GetNamespaceExecuter(ns)
	if err != nil {
		return nil, Error(ErrNamespaceNotFound, nil, namespace)
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

func (e *Executer) Variables() []string {
	var vars = make([]string, 0, len(e.objects))

	for i := range e.objects {
		vars = append(vars, e.objects[i].Name())
	}

	return vars
}

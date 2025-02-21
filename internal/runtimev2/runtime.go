package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/builtin"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/modules"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

// RuntimeMode is the runtime mode
type Runtime struct {
	functions map[string]lang.Method

	// executers is a map of namespace names to exec
	executers map[string]*Executer
}

// New creates a new runtime
func New() *Runtime {
	r := &Runtime{
		executers: make(map[string]*Executer),
	}
	r.functions = builtin.GetMethods(r.importer)
	r.BindModule(modules.NewRandModule())
	r.BindModule(modules.NewIOModule())
	r.BindModule(modules.NewHttpModule())
	r.BindModule(modules.NewJSONModule())
	return r
}

// Execute executes the given nodes
func (r *Runtime) Execute(nodes []*models.Node) (lang.Object, error) {
	if len(nodes) == 0 {
		return nil, nil
	}

	if nodes[0].Type == tokens.Namespace {
		namespace := nodes[0].Content
		nodes = nodes[1:]
		return r.Exec(ExecuterScopeGlobal, nil, namespace, nodes)
	}

	return r.Exec(ExecuterScopeGlobal, nil, "", nodes)
}

// Exec executes the given nodes in the given namespace
func (r *Runtime) Exec(scope ExecuterScope, parent *Executer, namespace string, nodes []*models.Node) (lang.Object, error) {
	ex, ok := r.executers[namespace]
	if !ok {
		ex = NewExecuter(scope, r, parent).WithName(namespace)
		r.executers[namespace] = ex
	}
	return ex.Execute(nodes)
}

// GetNamespaceExecuter gets the executer for the given namespace
func (r *Runtime) GetNamespaceExecuter(namespace string) (*Executer, error) {
	ex, ok := r.executers[namespace]
	if !ok {
		return nil, fmt.Errorf("namespace %v not found", namespace)
	}
	return ex, nil
}

// BindModule binds the given package to the given name
func (r *Runtime) BindModule(module lang.Module) {
	namespace := module.Namespace()

	ex, ok := r.executers[namespace]
	if !ok {
		ex = NewExecuter(ExecuterScopeGlobal, r, nil).WithName(namespace)
		r.executers[namespace] = ex
	}

	for name, object := range module.Objects() {
		ex.BindObject(name, object)
	}

	for name, method := range module.Methods() {
		ex.BindMethod(name, method)
	}
}

package runtimev2

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bndrmrtn/flare/internal/builtin"
	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/internal/modules"
	"github.com/bndrmrtn/flare/internal/state"
	"github.com/bndrmrtn/flare/internal/tokens"
	"github.com/bndrmrtn/flare/lang"
	"go.uber.org/zap"
)

const PackageDirectory = ".flmod"

// RuntimeMode is the runtime mode
type Runtime struct {
	// functions are a map of function names to function objects
	functions map[string]lang.Method

	// executers is a map of namespace names to exec
	executers map[string]*Executer
	// builtinModules is a map of builtin module names to module objects
	builtinModules map[string]lang.Module
	// sourceNamespaces is a list of source namespaces
	sourceNamespaces map[string]*Namespace

	// packages is a map of package names to runtime
	packages map[string]*Runtime

	stateProvider *state.Provider

	mu sync.RWMutex
}

// New creates a new runtime
func New(provider *state.Provider) (*Runtime, error) {
	zap.L().Info("creating new runtime")
	modules := modules.Get()

	r := &Runtime{
		executers:      make(map[string]*Executer),
		packages:       make(map[string]*Runtime),
		builtinModules: make(map[string]lang.Module, len(modules)),
		stateProvider:  provider,
	}
	r.functions = builtin.GetMethods(r.importer, r.evaler, provider)

	for _, module := range modules {
		r.BindModule(module)
	}

	if err := r.LoadSourceNamespaces(); err != nil {
		return nil, err
	}

	return r, nil
}

// Execute executes the given nodes
func (r *Runtime) Execute(nodes []*models.Node) (lang.Object, error) {
	namespace, nodes, err := r.GetNamespace(nodes)
	if err != nil {
		return nil, err
	}
	zap.L().Info("executing nodes", zap.String("namespace", namespace))
	return r.Exec(ExecuterScopeGlobal, nil, namespace, nodes)
}

func (r *Runtime) GetNamespace(nodes []*models.Node) (string, []*models.Node, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(nodes) == 0 {
		zap.L().Debug("no nodes to execute")
		return "", nil, nil
	}

	if nodes[0].Type == tokens.Namespace {
		namespace := nodes[0].Content
		nodes = nodes[1:]

		if _, ok := r.builtinModules[namespace]; ok {
			return "", nil, errs.WithDebug(fmt.Errorf("cannot use builtin '%s' as namespace", namespace), nodes[0].Debug)
		}

		zap.L().Debug("using namespace", zap.String("namespace", namespace))
		return namespace, nodes, nil
	}

	return "", nodes, nil
}

// Exec executes the given nodes in the given namespace
func (r *Runtime) Exec(scope ExecuterScope, parent *Executer, namespace string, nodes []*models.Node) (lang.Object, error) {
	r.mu.RLock()
	ex, ok := r.executers[namespace]
	r.mu.RUnlock()
	if !ok {
		ex = NewExecuter(scope, r, parent).WithName(namespace)
		r.mu.Lock()
		r.executers[namespace] = ex
		r.mu.Unlock()
	}
	return ex.Execute(nodes)
}

// GetNamespaceExecuter gets the executer for the given namespace
func (r *Runtime) GetNamespaceExecuter(namespace string) (*Executer, error) {
	zap.L().Debug("getting namespace executer", zap.String("namespace", namespace))
	r.mu.RLock()
	ns, ok := r.sourceNamespaces[namespace]
	r.mu.RUnlock()
	if ok {
		if !ns.Loaded {
			if _, err := r.Execute(ns.Nodes); err != nil {
				return nil, err
			}
			ns.Loaded = true
			return r.GetNamespaceExecuter(ns.Name)
		}
	}

	r.mu.RLock()
	ex, ok := r.executers[namespace]
	r.mu.RUnlock()
	if ok {
		return ex, nil
	}

	r.mu.RLock()
	mod, ok := r.builtinModules[namespace]
	r.mu.RUnlock()
	if ok {
		return r.UseModule(mod), nil
	}

	if !strings.Contains(namespace, ":") {
		return nil, fmt.Errorf("namespace %v not found", namespace)
	}

	parts := strings.Split(namespace, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid namespace %v", namespace)
	}

	r.mu.RLock()
	run, ok := r.packages[namespace]
	r.mu.RUnlock()
	if ok {
		ex, err := run.GetNamespaceExecuter(parts[1])
		if err == nil {
			return ex, nil
		}
	}

	ex, err := r.loadPackage(parts[0], parts[1])
	if err != nil {
		return nil, err
	}

	return ex, nil
}

// UseModule returns a module by its namespace
func (r *Runtime) UseModule(module lang.Module) *Executer {
	namespace := module.Namespace()
	zap.L().Debug("using module", zap.String("namespace", namespace))

	r.mu.RLock()
	ex, ok := r.executers[namespace]
	r.mu.RUnlock()
	if ok {
		return ex
	}

	ex = NewExecuter(ExecuterScopeGlobal, r, nil).WithName(namespace)

	r.mu.Lock()
	r.executers[namespace] = ex
	r.mu.Unlock()

	for name, object := range module.Objects() {
		ex.BindObject(name, object)
	}

	for name, method := range module.Methods() {
		ex.BindMethod(name, method)
	}

	return ex
}

// BindModule binds the given module to the runtime
func (r *Runtime) BindModule(module lang.Module) {
	r.mu.Lock()
	defer r.mu.Unlock()

	namespace := module.Namespace()
	r.builtinModules[namespace] = module
}

func (r *Runtime) loadPackage(author string, pkg string) (*Executer, error) {
	zap.L().Debug("loading package", zap.String("author", author), zap.String("package", pkg))

	run, err := New(r.stateProvider)
	if err != nil {
		return nil, err
	}

	r.mu.Lock()
	r.packages[author+":"+pkg] = run
	r.mu.Unlock()
	if _, err := run.Exec(ExecuterScopeGlobal, nil, pkg, nil); err != nil {
		return nil, err
	}

	ex, err := run.GetNamespaceExecuter(pkg)
	if err != nil {
		return nil, err
	}

	stat, err := os.Stat(filepath.Join(PackageDirectory, author, pkg))
	if err != nil {
		return nil, err
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("package %v not found", pkg)
	}

	var files []string
	err = filepath.WalkDir(filepath.Join(PackageDirectory, author, pkg), func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(d.Name()) == ".fl" {
			files = append(files, s)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if err := ex.LoadFile(file); err != nil {
			return nil, err
		}
	}

	return ex, nil
}

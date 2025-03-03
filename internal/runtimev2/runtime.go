package runtimev2

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bndrmrtn/zxl/internal/builtin"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/modules"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

const PackageDirectory = ".zxpack"

// RuntimeMode is the runtime mode
type Runtime struct {
	functions map[string]lang.Method

	// executers is a map of namespace names to exec
	executers map[string]*Executer

	// packages is a map of package names to runtime
	packages map[string]*Runtime
}

// New creates a new runtime
func New() *Runtime {
	r := &Runtime{
		executers: make(map[string]*Executer),
		packages:  make(map[string]*Runtime),
	}
	r.functions = builtin.GetMethods(r.importer)

	for _, module := range modules.Get() {
		r.BindModule(module)
	}

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
		if !strings.Contains(namespace, ":") {
			return nil, fmt.Errorf("namespace %v not found", namespace)
		}

		parts := strings.Split(namespace, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid namespace %v", namespace)
		}

		if run, ok := r.packages[namespace]; ok {
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

func (r *Runtime) loadPackage(author string, pkg string) (*Executer, error) {
	run := New()
	r.packages[author+":"+pkg] = run
	run.Exec(ExecuterScopeGlobal, nil, pkg, nil)

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
	err = filepath.WalkDir(filepath.Join(PackageDirectory, author, pkg), func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".zx" {
			files = append(files, s)
		}
		return nil
	})

	for _, file := range files {
		if err := ex.LoadFile(file); err != nil {
			return nil, err
		}
	}

	return ex, nil
}

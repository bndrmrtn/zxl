package runtimev2

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/lang"
)

// GetMethod gets a method by name
func (e *Executer) GetMethod(name string) (lang.Method, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		first := names[0]
		middle := names[1 : len(names)-1]
		last := names[len(names)-1]

		exec, err := e.runtime.GetNamespaceExecuter(first)
		if err == nil {
			return e.accessNamespace(exec, name, middle, last)
		}

		if obj, err := e.GetVariable(first); err == nil {
			if obj.Type() == lang.TDefinition {
				return e.accessDefinition(obj, name, middle, last)
			}
		}
	}

	if method, ok := e.functions[name]; ok {
		return method, nil
	}

	if obj, ok := e.objects[name]; ok {
		if obj.Type() == lang.TDefinition {
			return obj.Method("$init"), nil
		}
	}

	if e.parent != nil && (e.scope == ExecuterScopeBlock || e.scope == ExecuterScopeFunction || e.scope == ExecuterScopeDefinition) {
		return e.parent.GetMethod(name)
	}

	if fn, ok := e.runtime.functions[name]; ok {
		return fn, nil
	}

	return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
}

// GetVariable gets a variable by name
func (e *Executer) GetVariable(name string) (lang.Object, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		first := names[0]
		middle := names[1 : len(names)-1]
		last := names[len(names)-1]

		exec, err := e.runtime.GetNamespaceExecuter(first)
		if err == nil {
			return e.accessObjNamespace(exec, name, middle, last)
		}

		if obj, err := e.GetVariable(first); err == nil {
			if obj.Type() == lang.TDefinition {
				return e.accessObjDefinition(obj, name, middle, last)
			}
		}
	}

	if obj, ok := e.objects[name]; ok {
		return obj, nil
	}

	if e.parent != nil && e.scope == ExecuterScopeBlock {
		return e.parent.GetVariable(name)
	}

	return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
}

// accessNamespace accesses a method in a namespace
func (e *Executer) accessNamespace(exec *Executer, name string, middle []string, last string) (lang.Method, error) {
	if len(middle) == 0 {
		return exec.GetMethod(last)
	}

	first := middle[0]
	middle = middle[1:]

	obj, err := exec.GetVariable(first)
	if err != nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
	}

	// Middle elemek feldolgozása
	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
		}
	}

	method := obj.Method(last)
	if method == nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
	}

	return method, nil
}

// accessDefinition accesses a method in a definition
func (e *Executer) accessDefinition(def lang.Object, name string, middle []string, last string) (lang.Method, error) {
	var obj lang.Object = def

	// Middle elemek feldolgozása
	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
		}
	}

	method := obj.Method(last)
	if method == nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s()'", errs.CannotAccessFunction, name), nil)
	}

	return method, nil
}

// accessObjNamespace accesses an object in a namespace
func (e *Executer) accessObjNamespace(exec *Executer, name string, middle []string, last string) (lang.Object, error) {
	if len(middle) == 0 {
		return exec.GetVariable(last)
	}

	first := middle[0]
	middle = middle[1:]

	obj, err := exec.GetVariable(first)
	if err != nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
	}

	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
		}
	}

	obj = obj.Variable(last)
	if obj == nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
	}

	return obj, nil
}

// accessObjDefinition accesses an object in a definition
func (e *Executer) accessObjDefinition(def lang.Object, name string, middle []string, last string) (lang.Object, error) {
	obj := def

	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
		}
	}

	obj = obj.Variable(last)
	if obj == nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: '%s'", errs.CannotAccessVariable, name), nil)
	}

	return obj, nil
}

package runtimev2

import (
	"strings"

	"github.com/bndrmrtn/zxl/lang"
)

// GetMethod gets a method by name
func (e *Executer) GetMethod(name string) (lang.Method, error) {
	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		first := names[0]
		middle := names[1 : len(names)-1]
		last := names[len(names)-1]

		if first == "this" {
			def := e.isInsideDefinition(e)
			if def != nil {
				return e.GetMethod(strings.Join(append(middle, last), "."))
			}

			return nil, Error(ErrThisOutsideMethod, nil, fnErr(name))
		}

		exec, err := e.GetNamespaceExecuter(first)
		if err == nil {
			ob, err := e.accessNamespace(exec, name, middle, last)

			if ob != nil || err != nil {
				return ob, err
			}
		}

		if obj, err := e.GetVariable(first); err == nil {
			ob, err := e.accessDefinition(obj, name, middle, last)
			if ob != nil || err != nil {
				return ob, err
			}
		}
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	if method, ok := e.functions[name]; ok {
		return method, nil
	}

	if obj, ok := e.objects[name]; ok {
		if obj.Type() == lang.TDefinition {
			def := obj.(*lang.Definition)
			instance, err := def.NewInstance()
			if err != nil {
				return nil, err
			}
			return instance.Method("$init"), nil
		}
		if obj.Type() == lang.TFnRef {
			return obj.(*lang.Fn).Fn, nil
		}
	}

	if e.parent != nil && (e.scope == ExecuterScopeBlock || e.scope == ExecuterScopeFunction || e.scope == ExecuterScopeDefinition) {
		return e.parent.GetMethod(name)
	}

	if fn, ok := e.runtime.functions[name]; ok {
		return fn, nil
	}

	return nil, Error(ErrFunctionNotFound, nil, fnErr(name))
}

// GetVariable gets a variable by name
func (e *Executer) GetVariable(name string) (lang.Object, error) {
	if strings.Contains(name, ".") {
		names := strings.Split(name, ".")
		first := names[0]
		middle := names[1 : len(names)-1]
		last := names[len(names)-1]

		if first == "this" {
			ex := e
			for ex.parent != nil {
				ex = ex.parent
				if ex.scope == ExecuterScopeDefinition {
					break
				}
			}
			ob, err := ex.GetVariable(strings.Join(append(middle, last), "."))
			if ob != nil || err != nil {
				return ob, err
			}
		}

		exec, err := e.runtime.GetNamespaceExecuter(first)
		if err == nil {
			ob, err := e.accessObjNamespace(exec, name, middle, last)
			if ob != nil || err != nil {
				return ob, err
			}
		}

		if obj, err := e.GetVariable(first); err == nil {
			ob, err := e.accessObjDefinition(obj, name, middle, last)
			if ob != nil || err != nil {
				return ob, err
			}
		}
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	if obj, ok := e.objects[name]; ok {
		return obj, nil
	}

	if fn, ok := e.functions[name]; ok {
		return lang.NewFn(name, nil, fn), nil
	}

	if fn, ok := e.GetMethod(name); ok == nil {
		return lang.NewFn(name, nil, fn), nil
	}

	if e.parent != nil && (e.scope == ExecuterScopeBlock || e.scope == ExecuterScopeFunction) {
		return e.parent.GetVariable(name)
	}

	return nil, Error(ErrVariableNotFound, nil, name)
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
		return nil, Error(ErrFunctionNotFound, nil, fnErr(name))
	}

	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, Error(ErrFunctionNotFound, nil, fnErr(name))
		}
	}

	method := obj.Method(last)
	if method == nil {
		return nil, Error(ErrFunctionNotFound, nil, fnErr(name))
	}

	return method, nil
}

// accessDefinition accesses a method in a definition
func (e *Executer) accessDefinition(def lang.Object, name string, middle []string, last string) (lang.Method, error) {
	var obj = def

	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, Error(ErrFunctionNotFound, nil, fnErr(name))
		}
	}

	method := obj.Method(last)
	if method == nil {
		return nil, Error(ErrFunctionNotFound, nil, fnErr(name))
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
		return nil, Error(ErrVariableNotFound, nil, name)
	}

	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, Error(ErrVariableNotFound, nil, name)
		}
	}

	obj = obj.Variable(last)
	if obj == nil {
		return nil, Error(ErrVariableNotFound, nil, name)
	}

	return obj, nil
}

// accessObjDefinition accesses an object in a definition
func (e *Executer) accessObjDefinition(def lang.Object, name string, middle []string, last string) (lang.Object, error) {
	obj := def

	for _, part := range middle {
		obj = obj.Variable(part)
		if obj == nil {
			return nil, Error(ErrVariableNotFound, nil, name)
		}
	}

	obj = obj.Variable(last)
	if obj == nil {
		return nil, Error(ErrVariableNotFound, nil, name)
	}

	return obj, nil
}

func (e *Executer) isInsideDefinition(ex *Executer) *Executer {
	for ex.parent != nil {
		ex = ex.parent
		if ex.scope == ExecuterScopeDefinition {
			return ex
		}
	}
	return nil
}

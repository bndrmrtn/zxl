package runtime

import (
	"fmt"
	"strings"

	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// GetVariableValue gets the value of a variable
func (e *Executer) GetVariableValue(name string) (*models.Node, error) {
	if strings.Contains(name, ".") {
		parts := strings.Split(name, ".")
		last := parts[len(parts)-1]
		parts = parts[:len(parts)-1]

		// Check if the variable is a package
		if pkg, err := e.GetPackage(parts[0]); err == nil {
			variable, err := pkg.Access(last)
			if err == nil {
				return &models.Node{
					VariableType: variable.Type,
					Value:        variable.Value,
				}, nil
			}
		}

		ex, _, err := e.accessUnderlyingVariable(parts)
		if err != nil {
			return nil, err
		}

		return ex.GetVariableValue(last)
	}

	v, ok := e.vars[name]
	if !ok {
		if e.scope == ExecuterScopeBlock {
			return e.parent.GetVariableValue(name)
		}
		return nil, fmt.Errorf("%w: variable '%v' cannot be referenced", errs.RuntimeError, name)
	}

	if v.Reference {
		return e.GetVariableValue(v.Content)
	}

	if v.VariableType == tokens.ReferenceVariable {
		return e.GetVariableValue(v.Value.(string))
	}

	if v.VariableType == tokens.ExpressionVariable {
		return e.evaluateExpression(v)
	}

	return v, nil
}

// accessUnderlyingVariable accesses the underlying variable
func (e *Executer) accessUnderlyingVariable(args []string) (*Executer, *models.Node, error) {
	var (
		executer = e
		node     *models.Node
	)

	if len(args) == 0 {
		return nil, nil, fmt.Errorf("no variable to access")
	}

	if pkgName, ok := e.packages[args[0]]; ok {
		if nsEx, err := e.runtime.GetNamespaceExecuter(pkgName); err == nil {
			executer = nsEx
			args = args[1:]
			if len(args) == 0 {
				return executer, nil, nil
			}
		}
	}

	if len(args) == 0 {
		return nil, nil, fmt.Errorf("no variable to access")
	}

	if args[0] == "this" {
		if executer.parent == nil || executer.parent.scope == ExecuterScopeGlobal {
			return nil, nil, fmt.Errorf("'this' can only be accessed in a definition block function")
		}
		args = args[1:]
		executer = executer.parent
	}

	for _, part := range args {
		variable, err := executer.GetVariableValue(part)
		if err != nil {
			return nil, nil, err
		}

		switch variable.VariableType {
		case tokens.DefinitionReference:
			exec, ok := variable.Value.(*Executer)
			if !ok {
				return nil, nil, fmt.Errorf("variable %v is not a block", part)
			}
			executer = exec
			node = variable
		case tokens.FunctionCallVariable:
			return nil, nil, fmt.Errorf("cannot access function %v as variable", part)
		}
	}

	return executer, node, nil
}

package runtime

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (e *Executer) accessUnderlyingVariable(args []string) (*Executer, *models.Node, error) {
	var (
		executer = e
		node     *models.Node
	)

	if len(args) == 0 {
		return nil, nil, fmt.Errorf("no variable to access")
	}

	if nsEx, err := e.runtime.GetNamespaceExecuter(args[0]); err == nil {
		executer = nsEx
		args = args[1:]
		if len(args) == 0 {
			return executer, nil, nil
		}
	}

	if len(args) == 0 {
		return nil, nil, fmt.Errorf("no variable to access")
	}

	if args[0] == "this" {
		if executer.parent.scope != ExecuterScopeDefinition {
			return nil, nil, fmt.Errorf("this can only be accessed in a definition block function")
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

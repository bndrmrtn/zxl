package runtime

import (
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func getVariableRealValue(node *models.Node) any {
	switch node.VariableType {
	case tokens.NilVariable:
		return nil
	}
	return nil
}

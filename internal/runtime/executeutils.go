package runtime

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bndrmrtn/zexlang/internal/errs"
	"github.com/bndrmrtn/zexlang/internal/models"
	"github.com/bndrmrtn/zexlang/internal/tokens"
)

func (ex *Executer) evaluateExpression(n *models.Node) (*models.Node, error) {
	if n.VariableType != tokens.ExpressionVariable {
		return nil, errs.WithDebug(fmt.Errorf("cannot evaluate non expression variable"), n.Debug)
	}

	var expressionList []string
	var args map[string]any = make(map[string]any)

	for _, node := range n.Children {
		typ := node.Type
		if typ.IsOperator() {
			expressionList = append(expressionList, node.Content)
			continue
		}

		if node.VariableType == tokens.ReferenceVariable {
			var err error
			variable := node.Content
			node, err = ex.GetVariableValue(node.Content)
			if err != nil {
				return nil, err
			}
			node.Content = variable
			args[variable] = node.Value
		}

		if node.VariableType == tokens.FunctionCallVariable {
			ret, err := ex.executeFn(node)
			if err != nil {
				return nil, err
			}
			if len(ret) != 1 {
				return nil, errs.WithDebug(fmt.Errorf("function call in expression must return a single value"), node.Debug)
			}

			sum := fmt.Sprintf("%x", md5.Sum([]byte(node.Content)))
			node.Content = sum
			args[sum] = ret[0].Value
		}

		expressionList = append(expressionList, node.Content)
	}

	expression, err := govaluate.NewEvaluableExpression(strings.Join(expressionList, " "))
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}

	result, err := expression.Evaluate(args)
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}

	return &models.Node{
		VariableType: ex.getVarType(result),
		Type:         n.Type,
		Content:      fmt.Sprintf("%v", result),
		Value:        result,
		Debug:        n.Debug,
	}, nil
}

func (ex *Executer) getVarType(v any) tokens.VariableType {
	switch v.(type) {
	case int:
		return tokens.IntVariable
	case float64:
		return tokens.FloatVariable
	case string:
		return tokens.StringVariable
	case bool:
		return tokens.BoolVariable
	default:
		return tokens.NilVariable
	}
}

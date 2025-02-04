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

// evaluateExpression evaluates an expression
func (ex *Executer) evaluateExpression(n *models.Node) (*models.Node, error) {
	if n.VariableType != tokens.ExpressionVariable {
		return nil, errs.WithDebug(fmt.Errorf("cannot evaluate non expression variable"), n.Debug)
	}

	var (
		variableName   = n.Content
		expressionList []string
		args           map[string]any = make(map[string]any)
	)

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
				return nil, errs.WithDebug(err, n.Debug)
			}

			// Reset the content after the evaluation
			oldContent := node.Content
			defer func() {
				node.Content = oldContent
			}()

			sum := fmt.Sprintf("var_%x", md5.Sum([]byte(variable)))
			sum = sum[:10]

			node.Content = sum
			args[sum] = node.Value
		}

		if node.VariableType == tokens.FunctionCallVariable {
			ret, err := ex.executeFn(node)
			if err != nil {
				return nil, errs.WithDebug(err, n.Debug)
			}
			if ret == nil {
				return nil, errs.WithDebug(fmt.Errorf("function call in expression must return a single value"), node.Debug)
			}

			if ret.Type == tokens.DefinitionReference {
				return &models.Node{
					VariableType: tokens.DefinitionReference,
					Value:        ret.Value,
				}, nil
			}

			sum := fmt.Sprintf("var_%x", md5.Sum([]byte(node.Content)))
			sum = sum[:10]

			// Reset the content after the evaluation
			oldContent := node.Content
			defer func() {
				node.Content = oldContent
			}()

			node.Content = sum
			args[sum] = ret.Value
		}

		if node.VariableType == tokens.InlineValue {
			sum := fmt.Sprintf("var_%x", md5.Sum([]byte(node.Content)))
			sum = sum[:10]

			oldContent := node.Content
			defer func() {
				node.Content = oldContent
			}()

			node.Content = sum
			args[sum] = node.Value
		}

		expressionList = append(expressionList, node.Content)
	}

	if len(expressionList) == 0 {
		return nil, nil
	}

	value := strings.Join(expressionList, " ")
	expression, err := govaluate.NewEvaluableExpression(value)
	if err != nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: %w: %s", errs.RuntimeError, err, value), n.Debug)
	}

	result, err := expression.Evaluate(args)
	if err != nil {
		return nil, errs.WithDebug(fmt.Errorf("%w: %w: %s", errs.RuntimeError, err, value), n.Debug)
	}

	if fmt.Sprintf("%T", result) == "float64" && result.(float64) == float64(int64(result.(float64))) {
		result = int(result.(float64))
	}

	return &models.Node{
		VariableType: ex.getVarType(result),
		Type:         n.Type,
		Content:      variableName,
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

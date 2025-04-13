package runtimev2

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/Knetic/govaluate"
	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/internal/tmpl"
	"github.com/bndrmrtn/flare/internal/tokens"
	"github.com/bndrmrtn/flare/lang"
	"github.com/google/uuid"
)

// evaluateExpression evaluates an expression node
func (e *Executer) evaluateExpression(n *models.Node) (lang.Object, error) {
	if n.VariableType != tokens.ExpressionVariable {
		return nil, Error(ErrExpectedExpression, n.Debug)
	}

	var (
		variableName   string         = n.Content
		expressionList []string       = make([]string, 0, len(n.Children))
		args           map[string]any = make(map[string]any, len(n.Children))
		nameVal        map[string]any = make(map[string]any, len(n.Children))
	)

	for _, node := range n.Children {
		var (
			nodeType     = node.Type
			variableType = node.VariableType
		)

		// If the node is an operator, add it to the expression list
		if nodeType.IsOperator() {
			expressionList = append(expressionList, node.Content)
			continue
		}

		// If the node is a variable reference, get the variable value and add it to the expression list
		if variableType == tokens.ReferenceVariable {
			obj, err := e.GetVariable(node.Content)
			if err != nil {
				return nil, errs.WithDebug(err, n.Debug)
			}

			if len(node.ObjectAccessors) > 0 {
				acc, err := e.accessObject(obj, node.ObjectAccessors)
				if err != nil {
					return nil, errs.WithDebug(err, n.Debug)
				}
				obj = acc
			}

			sum := newVariableName()

			expressionList = append(expressionList, sum)

			if obj.Type() == lang.TList {
				args[sum] = obj
				nameVal[sum] = obj
				continue
			}

			args[sum] = obj.Value()
			nameVal[sum] = obj.Value()
			continue
		}

		// If the node is a function call, call the function and add the result to the expression list
		if variableType == tokens.FunctionCallVariable {
			obj, err := e.callFunctionFromNode(node)
			if err != nil {
				return nil, errs.WithDebug(err, n.Debug)
			}

			if obj == nil {
				// Allow empty returns as nil
				obj = lang.NewNil("nil", node.Debug)
			}

			if len(node.ObjectAccessors) > 0 {
				acc, err := e.accessObject(obj, node.ObjectAccessors)
				if err != nil {
					return nil, errs.WithDebug(err, n.Debug)
				}
				obj = acc
			}

			sum := newVariableName()

			expressionList = append(expressionList, sum)

			if obj.Type() == lang.TList {
				args[sum] = obj
				nameVal[sum] = obj
				continue
			}

			args[sum] = obj.Value()
			nameVal[sum] = obj.Value()
			continue
		}

		if variableType == tokens.FunctionVariable {
			name, method, err := e.createMethodFromNode(node)
			if err != nil {
				return nil, errs.WithDebug(err, n.Debug)
			}

			if name != "fn" {
				return nil, Error(ErrNamedInlineFunction, n.Debug)
			}

			obj := lang.NewFn(name, node.Debug, method)

			sum := newVariableName()

			expressionList = append(expressionList, sum)

			args[sum] = obj
			nameVal[sum] = obj
			continue
		}

		if variableType == tokens.InlineValue {
			sum := newVariableName()

			if node.Type == tokens.TemplateLiteral {
				rawTmpl, err := tmpl.NewTemplate(node.Value.(string))
				if err != nil {
					return nil, errs.WithDebug(err, n.Debug)
				}
				s, err := e.parseTemplate(rawTmpl)
				if err != nil {
					return nil, errs.WithDebug(err, n.Debug)
				}

				expressionList = append(expressionList, sum)
				args[sum] = s
				nameVal[sum] = fmt.Sprintf("%q", s)
				continue
			}

			expressionList = append(expressionList, sum)
			args[sum] = node.Value
			nameVal[sum] = node.Value
			continue
		}

		if variableType == tokens.ExpressionVariable {
			obj, err := e.evaluateExpression(node)
			if err != nil {
				return nil, errs.WithDebug(err, n.Debug)
			}

			sum := newVariableName()

			expressionList = append(expressionList, sum)

			if obj.Type() == lang.TList {
				args[sum] = obj
				nameVal[sum] = obj
				continue
			}

			args[sum] = obj.Value()
			nameVal[sum] = obj.Value()
			continue
		}

		if variableType == tokens.ArrayVariable {
			_, obj, err := e.createObjectFromNode(node)
			if err != nil {
				return nil, errs.WithDebug(err, n.Debug)
			}

			sum := newVariableName()

			expressionList = append(expressionList, sum)

			if obj.Type() == lang.TList {
				args[sum] = obj
				nameVal[sum] = obj
				continue
			}

			args[sum] = obj.Value()
			nameVal[sum] = obj.Value()
			continue
		}
	}

	if len(expressionList) == 0 {
		return nil, nil
	}

	// Handle the case where the expression is a single object
	if len(expressionList) == 1 && len(args) == 1 {
		expression := expressionList[0]
		if obj, ok := args[expression]; ok {
			if v, ok := obj.(lang.Object); ok {
				v = v.Copy()
				v.Rename(variableName)
				return v, nil
			}
		}
	}

	value := strings.Join(expressionList, " ")
	expression, err := govaluate.NewEvaluableExpression(value)
	if err != nil {
		for name, val := range nameVal {
			value = strings.ReplaceAll(value, name, fmt.Sprintf("%v", val))
		}

		return nil, Error(err, n.Debug, value)
	}

	result, err := expression.Evaluate(args)
	if err != nil {
		for name, val := range nameVal {
			value = strings.ReplaceAll(value, name, fmt.Sprintf("%v", val))
		}

		return nil, Error(err, n.Debug, value)
	}

	if fmt.Sprintf("%T", result) == "float64" && result.(float64) == float64(int64(result.(float64))) {
		result = int(result.(float64))
	}

	_, obj, err := e.createObjectFromNode(&models.Node{
		VariableType: e.getVarType(result),
		Type:         n.Type,
		Content:      variableName,
		Value:        result,
		Debug:        n.Debug,
	})
	if err != nil {
		return nil, errs.WithDebug(err, n.Debug)
	}
	return obj, nil
}

// getVarType returns the variable type of the given value
func (e *Executer) getVarType(v any) tokens.VariableType {
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

// getVariableTypeFromType returns the variable type of the given node type
func (e *Executer) getVariableTypeFromType(n *models.Node) tokens.VariableType {
	switch n.Type {
	case tokens.String:
		return tokens.StringVariable
	case tokens.Number:
		if val, ok := n.Map["isFloat"].(bool); val && ok {
			return tokens.FloatVariable
		}
		return tokens.IntVariable
	case tokens.Bool:
		return tokens.BoolVariable
	case tokens.TemplateLiteral:
		return tokens.TemplateVariable
	case tokens.Array:
		return tokens.ArrayVariable
	default:
		return tokens.NilVariable
	}
}

func (e *Executer) getTypeFromValue(value any) tokens.TokenType {
	if v, ok := value.(lang.Object); ok {
		return v.Type().TokenType()
	}

	switch value.(type) {
	case int, float64:
		return tokens.Number
	case string:
		return tokens.String
	case bool:
		return tokens.Bool
	default:
		return tokens.Nil
	}
}

func newVariableName() string {
	sum := fmt.Sprintf("var_%x", md5.Sum([]byte(uuid.NewString())))
	return sum[:10]
}

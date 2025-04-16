package formatter

import (
	"fmt"
	"os"
	"strings"

	"github.com/flarelang/flare/internal/ast"
	"github.com/flarelang/flare/internal/lexer"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
)

type FileFmt struct {
	fileName string

	hasNamespace  bool
	defaultIndent int
}

func NewFileFmt(file string) *FileFmt {
	return &FileFmt{
		fileName: file,
	}
}

func (f *FileFmt) WithIndent(indent int) *FileFmt {
	f.defaultIndent = indent
	return f
}

func (f *FileFmt) Format() error {
	file, err := os.Open(f.fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	lx := lexer.New(f.fileName)
	ts, err := lx.Parse(file)
	if err != nil {
		return err
	}

	builder := ast.NewBuilder()
	nodes, err := builder.Build(ts)
	if err != nil {
		return err
	}

	sb := strings.Builder{}
	nodesLen := len(nodes)
	for i, node := range nodes {
		var nextNode *models.Node
		if i < nodesLen-1 {
			nextNode = nodes[i+1]
		}

		err := f.formatNode(ScopeGlobal, 0, &sb, node, nextNode)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(f.fileName, []byte(sb.String()), os.ModePerm)
}

func (f *FileFmt) formatNode(scope Scope, indent int, sb *strings.Builder, node *models.Node, next *models.Node) error {
	tab := strings.Repeat("\t", indent)
	sb.WriteString(tab)

	switch node.Type {
	case tokens.Namespace:
		sb.WriteString("namespace " + node.Content + ";\n\n")
		f.hasNamespace = true
		break
	case tokens.Use:
		sb.WriteString("use " + node.Content)
		if as, ok := node.Value.(string); ok {
			sb.WriteString(" as " + as)
		}
		sb.WriteString(";\n")
		if next != nil && next.Type != tokens.Use {
			sb.WriteString("\n")
		}
		break
	case tokens.Function:
		sb.WriteString("fn " + node.Content + "(")
		for i, param := range node.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(param.Content)
		}
		sb.WriteString(") {")
		if len(node.Children) != 0 {
			sb.WriteString("\n")

			var (
				childLen int = len(node.Children)
				nextNode *models.Node
			)

			for i, stmt := range node.Children {
				if i < childLen-1 {
					nextNode = node.Children[i+1]
				}

				err := f.formatNode(ScopeFunction, indent+1, sb, stmt, nextNode)
				if err != nil {
					return err
				}
			}
			sb.WriteString(tab)
		}
		sb.WriteString("}\n")
		if next != nil && scope != ScopeFunction && scope != ScopeDefinition {
			sb.WriteString("\n")
		}
		break
	case tokens.Let:
		if node.VariableType == tokens.NilVariable {
			sb.WriteString("let " + node.Content + ";\n")
			break
		}

		sb.WriteString("let " + node.Content + " = ")
		sb.WriteString(f.formatValue(node, indent))
		sb.WriteString(";\n")

		if next != nil {
			if next.Type != tokens.Let && next.Type != tokens.Const && next.Type != tokens.Assign {
				sb.WriteString("\n")
			}
		}
		break
	case tokens.Const:
		if node.VariableType == tokens.NilVariable {
			sb.WriteString("const " + node.Content + ";\n")
			break
		}

		sb.WriteString("const " + node.Content + " = ")
		sb.WriteString(f.formatValue(node))
		sb.WriteString(";\n")

		if next != nil {
			if next.Type != tokens.Let && next.Type != tokens.Const && next.Type != tokens.Assign {
				sb.WriteString("\n")
			}
		}
		break
	case tokens.Assign:
		sb.WriteString(node.Content + " = ")
		sb.WriteString(f.formatValue(node))
		sb.WriteString(";\n")

		if next != nil {
			if next.Type != tokens.Let && next.Type != tokens.Const && next.Type != tokens.Assign {
				sb.WriteString("\n")
			}
		}
		break
	case tokens.Define:
		sb.WriteString("define " + node.Content + " {")
		if len(node.Children) != 0 {
			sb.WriteString("\n")
		}

		var (
			childLen int = len(node.Children)
			nextNode *models.Node
		)

		for i, stmt := range node.Children {
			if i < childLen-1 {
				nextNode = node.Children[i+1]
			}

			err := f.formatNode(ScopeDefinition, indent+1, sb, stmt, nextNode)
			if err != nil {
				return err
			}
		}
		sb.WriteString(tab + "}\n")

		if next != nil {
			sb.WriteRune('\n')
		}
		break
	case tokens.Addition, tokens.Subtraction, tokens.Multiplication, tokens.Division, tokens.Equation, tokens.NotEquation, tokens.Greater, tokens.GreaterOrEqual, tokens.Less, tokens.LessOrEqual, tokens.And, tokens.Or, tokens.Not, tokens.Power:
		sb.WriteString(node.Content)
		break
	case tokens.String, tokens.Number, tokens.Bool:
		sb.WriteString(node.ValueString())
		break
	case tokens.Identifier:
		if node.VariableType == tokens.ReferenceVariable {
			sb.WriteString(node.Content)
			break
		}
		sb.WriteString(f.formatValue(node))
		break
	case tokens.FuncCall:
		sb.WriteString(node.Content + "(")
		for i, arg := range node.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.formatValue(arg))
		}
		sb.WriteString(");\n")
		break
	case tokens.LeftParenthesis:
		newSb := strings.Builder{}
		sb.WriteString("(")

		for i, child := range node.Children {
			err := f.formatNode(ScopeGlobal, 0, &newSb, child, nil)
			if err != nil {
				return err
			}
			if i < len(node.Children)-1 {
				newSb.WriteString(" ")
			}
		}
		sb.WriteString(newSb.String())
		sb.WriteString(")")
		break
	case tokens.List:
		sb.WriteString("[")
		for i, child := range node.Children {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.formatValue(child))
		}
		sb.WriteString("]")
		break
	case tokens.Array:
		sb.WriteString("array {")
		if len(node.Children) == 0 {
			sb.WriteString("}")
			break
		}

		sb.WriteRune('\n')
		for _, child := range node.Children {
			sb.WriteString(tab + "\t")
			var keyB strings.Builder
			if err := f.formatNode(ScopeGlobal, 0, &keyB, child.Args[0], nil); err != nil {
				return err
			}
			key := keyB.String()
			if lexer.IsIdentifier(key) {
				sb.WriteString(key)
			} else {
				sb.WriteString(fmt.Sprintf("%q", child.Args[0].Value))
			}

			sb.WriteString(": ")

			var vB strings.Builder
			nInd := 0
			if child.Children[0].Type == tokens.Array {
				nInd = indent + 1
			}
			if err := f.formatNode(ScopeGlobal, nInd, &vB, child.Children[0], nil); err != nil {
				return err
			}

			sb.WriteString(strings.TrimLeft(vB.String(), " \t"))

			sb.WriteString(",\n")
		}

		sb.WriteString(tab + "}")
		break
	}
	return nil
}

func (f *FileFmt) formatValue(node *models.Node, indentArgs ...int) string {
	indent := 0
	if len(indentArgs) > 0 {
		indent = indentArgs[0]
	}

	if node.VariableType == tokens.ExpressionVariable {
		return f.formatExpression(false, node.Children, indent)
	}

	if node.VariableType == tokens.ListVariable {
		var sb strings.Builder
		node.Type = tokens.List
		_ = f.formatNode(ScopeGlobal, 0, &sb, node, nil)
		return sb.String()
	}

	val := fmt.Sprintf("%v", node.ValueString())
	if val == "<nil>" {
		return "nil"
	}

	return val
}

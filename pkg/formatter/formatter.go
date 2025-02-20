package formatter

import (
	"fmt"
	"os"
	"strings"

	"github.com/bndrmrtn/zxl/internal/ast"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

type Formatter struct {
	folder           string
	useTokenStart    bool
	skipNextLineFunc bool
}

func New(folder string) *Formatter {
	return &Formatter{folder: folder}
}

func (f *Formatter) Format() error {
	files, err := getFiles(f.folder, ".zx")
	if err != nil {
		return err
	}
	for _, file := range files {
		err := f.formatFile(file)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *Formatter) formatFile(name string) error {
	file, err := os.Open(name)
	if err != nil {
		return err
	}
	defer file.Close()

	lx := lexer.New(name)
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
	for _, node := range nodes {
		err := f.formatNode(0, &sb, node)
		if err != nil {
			return err
		}
	}

	return os.WriteFile(name, []byte(sb.String()), os.ModePerm)
}

func (f *Formatter) formatNode(indent int, sb *strings.Builder, node *models.Node) error {
	if node.Type != tokens.Use && f.useTokenStart {
		sb.WriteString("\n")
		f.useTokenStart = false
	}

	if node.Type == tokens.Function && f.skipNextLineFunc {
		sb.WriteString("\n")
		f.skipNextLineFunc = false
	}

	tab := strings.Repeat("\t", indent)
	sb.WriteString(tab)

	switch node.Type {
	case tokens.Namespace:
		sb.WriteString("namespace " + node.Content + ";\n")
	case tokens.Use:
		if !f.useTokenStart {
			sb.WriteString("\n")
		}
		sb.WriteString("use " + node.Content)
		if as, ok := node.Value.(string); ok {
			sb.WriteString(" as " + as)
		}
		sb.WriteString(";\n")
		f.useTokenStart = true
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
		}
		for _, stmt := range node.Children {
			err := f.formatNode(indent+1, sb, stmt)
			if err != nil {
				return err
			}
		}
		sb.WriteString(tab + "}\n")
		f.skipNextLineFunc = true
	case tokens.Let:
		if node.VariableType == tokens.NilVariable {
			sb.WriteString("let " + node.Content + ";\n")
			break
		}

		sb.WriteString("let " + node.Content + " = ")
		sb.WriteString(f.formatValue(node))
		sb.WriteString(";\n")
		break
	case tokens.Define:
		sb.WriteString("\ndefine " + node.Content + " {")
		if len(node.Children) != 0 {
			sb.WriteString("\n")
		}
		for _, stmt := range node.Children {
			err := f.formatNode(indent+1, sb, stmt)
			if err != nil {
				return err
			}
		}
		sb.WriteString(tab + "}\n")
	case tokens.Addition, tokens.Subtraction, tokens.Multiplication, tokens.Division, tokens.Equation, tokens.NotEquation, tokens.Greater, tokens.GreaterOrEqual, tokens.Less, tokens.LessOrEqual, tokens.And, tokens.Or, tokens.Not, tokens.Power:
		sb.WriteString(node.Content)
	case tokens.String, tokens.Number, tokens.Bool:
		sb.WriteString(node.ValueString())
	case tokens.Identifier:
		sb.WriteString(node.Content)
	case tokens.FuncCall:
		sb.WriteString(node.Content + "(")
		for i, arg := range node.Args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(f.formatValue(arg))
		}
		sb.WriteString(");\n")
	}
	return nil
}

func (f *Formatter) formatValue(node *models.Node) string {
	if node.VariableType == tokens.ExpressionVariable {
		return f.formatExpression(false, node.Children)
	}

	return fmt.Sprintf("%v", node.ValueString())
}

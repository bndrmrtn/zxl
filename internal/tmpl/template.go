package tmpl

import (
	"regexp"
	"strings"

	"github.com/bndrmrtn/zxl/internal/ast"
	"github.com/bndrmrtn/zxl/internal/lexer"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

// Part represents a part of a template.
type Part struct {
	// Static indicates whether the content is static or dynamic.
	Static bool
	// Content represents the content of the template part.
	Content string
	// Node represents the node associated with the template part.
	Node *models.Node
}

// NewTemplate creates a new template from the given string.
func NewTemplate(s string) ([]Part, error) {
	var result []Part
	re := regexp.MustCompile(`\{\{(.*?)}}`)
	matches := re.FindAllStringSubmatchIndex(s, -1)

	lastIndex := 0
	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		if match[0] > lastIndex {
			result = append(result, Part{
				Static:  true,
				Content: s[lastIndex:match[0]],
			})
		}

		content := strings.TrimSpace(s[match[2]:match[3]])
		lx := lexer.New("")
		ts, err := lx.Parse(strings.NewReader(content + ";"))
		if err != nil {
			return nil, err
		}

		nodes, err := ast.NewBuilder().Build(ts)
		if err != nil {
			return nil, err
		}

		result = append(result, Part{
			Static:  false,
			Content: content,
			Node: &models.Node{
				VariableType: tokens.ExpressionVariable,
				Children:     nodes,
			},
		})

		lastIndex = match[1]
	}

	if lastIndex < len(s) {
		result = append(result, Part{
			Static:  true,
			Content: s[lastIndex:],
		})
	}

	return result, nil
}

package formatter

import (
	"strings"

	"github.com/bndrmrtn/zxl/internal/models"
)

func (f *Formatter) formatExpression(wrap bool, nodes []*models.Node) string {
	result := strings.Builder{}

	if wrap {
		result.WriteString("(")
	}

	for i, node := range nodes {
		if i > 0 {
			result.WriteString(" ")
		}
		f.formatNode(0, &result, node)
	}

	if wrap {
		result.WriteString(")")
	}

	return result.String()
}

package formatter

import (
	"strings"

	"github.com/bndrmrtn/flare/internal/models"
)

func (f *FileFmt) formatExpression(wrap bool, nodes []*models.Node, indent int) string {
	result := strings.Builder{}

	if wrap {
		result.WriteString("(")
	}

	for i, node := range nodes {
		if i > 0 {
			result.WriteString(" ")
		}
		var vB strings.Builder
		f.formatNode(ScopeGlobal, indent, &vB, node, nil)
		result.WriteString(strings.TrimLeft(vB.String(), " \t"))
	}

	if wrap {
		result.WriteString(")")
	}

	return result.String()
}

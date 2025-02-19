package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
)

// accessObject accesses an object
func (e *Executer) accessObject(obj lang.Object, accessors []*models.Node) (lang.Object, error) {
	if obj.Type() != lang.TList && obj.Type() != lang.TDefinition {
		return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
	}

	access, err := e.getObjAccessors(accessors)
	if err != nil {
		return nil, err
	}

	if obj.Type() == lang.TList {
		li, ok := obj.Value().([]lang.Object)
		if !ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
		}

		var value any = li
		for inx, a := range access {
			i, ok := a.(int)
			if !ok {
				return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
			}

			v, ok := value.([]lang.Object)
			if !ok {
				if v, ok := value.(lang.Object); ok {
					return e.accessObject(v, accessors[inx:])
				}

				return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
			}

			if i < 0 || i >= len(v) {
				return nil, errs.WithDebug(fmt.Errorf("%w: %v, length: %d", errs.IndexOutOfRange, i, len(v)), accessors[0].Debug)
			}

			value = v[i]
		}

		if v, ok := value.(lang.Object); ok {
			value = v.Value()
		}

		_, ob, err := e.createObjectFromNode(&models.Node{
			VariableType: tokens.InlineValue,
			Type:         e.getTypeFromValue(value),
			Content:      obj.Name(),
			Value:        value,
			Debug:        obj.Debug(),
		})
		if err != nil {
			return nil, errs.WithDebug(err, obj.Debug())
		}

		return ob, nil
	}

	return obj, nil
}

func (e *Executer) getObjAccessors(accessors []*models.Node) ([]any, error) {
	var access = make([]any, 0, len(accessors))

	for _, accessor := range accessors {
		if accessor.VariableType == tokens.ReferenceVariable {
			obj, err := e.GetVariable(accessor.Content)
			if err != nil {
				return nil, errs.WithDebug(err, accessor.Debug)
			}
			access = append(access, obj.Value())
		} else {
			access = append(access, accessor.Value)
		}
	}

	return access, nil
}

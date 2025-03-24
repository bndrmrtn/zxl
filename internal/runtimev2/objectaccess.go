package runtimev2

import (
	"fmt"

	"github.com/bndrmrtn/zxl/internal/errs"
	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/internal/tokens"
	"github.com/bndrmrtn/zxl/lang"
	"go.uber.org/zap"
)

// accessObject accesses an object
func (e *Executer) accessObject(obj lang.Object, accessors []*models.Node) (lang.Object, error) {
	if len(accessors) == 0 {
		return obj, nil
	}

	if obj.Type() != lang.TList && obj.Type() != lang.TDefinition && obj.Type() != lang.TArray {
		return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
	}

	access, err := e.getObjAccessors(accessors[0])
	if err != nil {
		return nil, err
	}

	zap.L().Debug("accessing object", zap.Any("object", obj), zap.Any("accessor", access))

	var value any
	if obj.Type() == lang.TList {
		li, ok := obj.Value().([]lang.Object)
		if !ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
		}

		i, ok := access.(int)
		if !ok || i < 0 || i >= len(li) {
			return nil, errs.WithDebug(fmt.Errorf("%w: %v, length: %d", errs.IndexOutOfRange, i, len(li)), accessors[0].Debug)
		}
		value = li[i]
	} else if obj.Type() == lang.TArray {
		arr, ok := obj.(*lang.Array)
		if !ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: cannot access object with type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
		}

		o, ok := arr.Access(access)
		if !ok {
			return nil, errs.WithDebug(fmt.Errorf("%w: key not found: '%v'", errs.ValueError, access), accessors[0].Debug)
		}
		value = o
	} else {
		return nil, errs.WithDebug(fmt.Errorf("%w: unsupported object type '%s'", errs.ValueError, obj.Type()), accessors[0].Debug)
	}

	if len(accessors) > 1 {
		if nextObj, ok := value.(lang.Object); ok {
			return e.accessObject(nextObj, accessors[1:])
		}
		return nil, errs.WithDebug(fmt.Errorf("%w: unexpected value type", errs.ValueError), accessors[0].Debug)
	}

	if v, ok := value.(lang.Object); ok {
		return v.Copy(), nil
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

func (e *Executer) getObjAccessors(accessor *models.Node) (any, error) {
	if accessor.VariableType == tokens.ReferenceVariable {
		obj, err := e.GetVariable(accessor.Content)
		if err != nil {
			return nil, errs.WithDebug(err, accessor.Debug)
		}

		zap.L().Debug("accessing object", zap.Any("object", obj))

		return obj.Value(), nil
	}

	zap.L().Debug("accessing object", zap.Any("accessor", accessor))

	return accessor.Value, nil
}

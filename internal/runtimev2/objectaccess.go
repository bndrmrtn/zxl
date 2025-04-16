package runtimev2

import (
	"fmt"

	"github.com/flarelang/flare/internal/errs"
	"github.com/flarelang/flare/internal/models"
	"github.com/flarelang/flare/internal/tokens"
	"github.com/flarelang/flare/lang"
	"go.uber.org/zap"
)

// accessObject accesses an object
func (e *Executer) accessObject(obj lang.Object, accessors []*models.Node) (lang.Object, error) {
	if len(accessors) == 0 {
		return obj, nil
	}

	currentAccessor := accessors[0]

	if obj.Type() != lang.TList && obj.Type() != lang.TDefinition && obj.Type() != lang.TArray {
		return nil, Error(ErrInvalidIndexAccess, accessors[0].Debug, obj.Type())
	}

	access, err := e.getObjAccessors(currentAccessor)
	if err != nil {
		return nil, err
	}

	zap.L().Debug("accessing object", zap.Any("object", obj), zap.Any("accessor", access))

	var value any
	if obj.Type() == lang.TList {
		li, ok := obj.Value().([]lang.Object)
		if !ok {
			return nil, Error(ErrInvalidIndexAccess, currentAccessor.Debug, obj.Type())
		}

		i, ok := access.(int)
		if !ok || i < 0 || i >= len(li) {
			return nil, Error(ErrIndexOutOfBounds, currentAccessor.Debug, fmt.Sprintf("%d length: %d", i, len(li)))
		}
		value = li[i]
	} else if obj.Type() == lang.TArray {
		arr, ok := obj.(*lang.Array)
		if !ok {
			return nil, Error(ErrInvalidIndexAccess, currentAccessor.Debug, obj.Type())
		}

		o, ok := arr.Access(access)
		if !ok {
			return nil, Error(ErrKeyNotFound, currentAccessor.Debug, access)
		}
		value = o
	} else {
		return nil, Error(ErrInvalidValue, currentAccessor.Debug, fmt.Sprintf("unsupported object type: %s", obj.Type()))
	}

	if len(accessors) > 1 {
		if nextObj, ok := value.(lang.Object); ok {
			return e.accessObject(nextObj, accessors[1:])
		}
		return nil, Error(ErrInvalidValue, currentAccessor.Debug, "unexpected value type")
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

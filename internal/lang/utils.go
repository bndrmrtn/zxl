package lang

import "fmt"

func Addr(s any) Object {
	return addr(s)
}

func addr(s any) Object {
	return NewString("$addr", fmt.Sprintf("%p", s), nil)
}

func FromValue(data any) (Object, error) {
	switch value := data.(type) {
	case map[string]interface{}:
		var obj = NewArray("object", nil, nil, nil)
		for key, val := range value {
			value, err := FromValue(val)
			if err != nil {
				return nil, err
			}

			_, err = obj.Method("$bind").Execute([]Object{
				NewString("key", key, nil),
				value,
			})
			if err != nil {
				return nil, err
			}
		}
		return obj, nil
	case []interface{}:
		var li = make([]Object, len(value))
		for i, val := range value {
			value, err := FromValue(val)
			if err != nil {
				return nil, err
			}

			li[i] = value
		}
		return NewList("list", li, nil), nil
	case int:
		return NewInteger("number", int(value), nil), nil
	case int64:
		return NewInteger("number", int(value), nil), nil
	case float64:
		return NewFloat("number", value, nil), nil
	case float32:
		return NewFloat("number", float64(value), nil), nil
	case string:
		return NewString("string", value, nil), nil
	case bool:
		return NewBool("bool", value, nil), nil
	case nil:
		return NewNil("nil", nil), nil
	case interface{}:
		return FromValue(value)
	}
	return nil, fmt.Errorf("unsupported type %T", data)
}

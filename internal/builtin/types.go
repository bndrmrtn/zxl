package builtin

import (
	"fmt"
	"strconv"

	"github.com/bndrmrtn/flare/lang"
)

func setTypeMethods(m map[string]lang.Method) map[string]lang.Method {
	m["string"] = lang.NewFunction(toString).WithArg("object")
	m["int"] = lang.NewFunction(toInt).WithArg("object")
	m["float"] = lang.NewFunction(toFloat).WithArg("object")
	m["bool"] = lang.NewFunction(toBool).WithArg("object")

	return m
}

func toString(args []lang.Object) (lang.Object, error) {
	return lang.NewString("string", args[0].String(), args[0].Debug()), nil
}

func toInt(args []lang.Object) (lang.Object, error) {
	var value int

	switch v := args[0].Value().(type) {
	case int:
		value = v
	case int8:
		value = int(v)
	case int16:
		value = int(v)
	case int32:
		value = int(v)
	case int64:
		value = int(v)
	case float32:
		value = int(v)
	case float64:
		value = int(v)
	case string:
		val, err := strconv.Atoi(v)
		if err != nil {
			return nil, err
		}
		value = val
		break
	case bool:
		if v {
			value = 1
		}
		value = 0
		break
	default:
		return nil, fmt.Errorf("unsupported type: %T", args[0].Type())
	}

	return lang.NewInteger("convert", value, args[0].Debug()), nil
}

func toFloat(args []lang.Object) (lang.Object, error) {
	var value float64

	switch v := args[0].Value().(type) {
	case int:
		value = float64(v)
	case int8:
		value = float64(v)
	case int16:
		value = float64(v)
	case int32:
		value = float64(v)
	case int64:
		value = float64(v)
	case float32:
		value = float64(v)
	case float64:
		value = v
	case string:
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, err
		}
		value = val
		break
	case bool:
		if v {
			value = 1.0
		} else {
			value = 0.0
		}
		break
	default:
		return nil, fmt.Errorf("unsupported type: %T", args[0].Type())
	}

	return lang.NewFloat("convert", value, args[0].Debug()), nil
}

func toBool(args []lang.Object) (lang.Object, error) {
	var value bool

	switch v := args[0].Value().(type) {
	case int:
		value = v != 0
	case int8:
		value = v != 0
	case int16:
		value = v != 0
	case int32:
		value = v != 0
	case int64:
		value = v != 0
	case float32:
		value = v != 0
	case float64:
		value = v != 0
	case string:
		value = v != ""
	case bool:
		value = v
	default:
		return nil, fmt.Errorf("unsupported type: %T", args[0].Type())
	}

	return lang.NewBool("convert", value, args[0].Debug()), nil
}

package modules

import (
	"encoding/json"
	"fmt"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type JSON struct{}

func NewJSONModule() *JSON {
	return &JSON{}
}

func (*JSON) Namespace() string {
	return "json"
}

func (j *JSON) Objects() map[string]lang.Object {
	return nil
}

func (j *JSON) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"parse":    lang.NewFunction([]string{"string"}, j.parse, nil),
		"toString": lang.NewFunction([]string{"object"}, j.toString, nil),
	}
}

func (j *JSON) parse(args []lang.Object) (lang.Object, error) {
	if args[0].Type() != lang.TString {
		return nil, fmt.Errorf("expected string, got %s", args[0].Type())
	}

	value := args[0].Value().(string)

	var data any
	err := json.Unmarshal([]byte(value), &data)
	if err != nil {
		return nil, err
	}

	obj, err := j.traverseJSON(data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (j *JSON) traverseJSON(data interface{}) (lang.Object, error) {
	return lang.FromValue(data)
}

func (j *JSON) toString(args []lang.Object) (lang.Object, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("expected at least one argument")
	}

	jsonData, err := j.convertToJSON(args[0])
	if err != nil {
		return nil, err
	}

	return lang.NewString("string", string(jsonData), nil), nil
}

func (j *JSON) convertToJSON(obj lang.Object) ([]byte, error) {
	switch obj.Type() {
	case lang.TArray:
		arr := obj.(*lang.Array)
		keys := arr.Keys
		data := make(map[string]interface{}, len(keys))

		for _, keyObj := range keys {
			key := fmt.Sprintf("%v", keyObj.Value()) // Ensure key is a string
			valObj := arr.Map[keyObj]
			val, err := j.convertToJSON(valObj)
			if err != nil {
				return nil, err
			}
			var value interface{}
			json.Unmarshal(val, &value)
			data[key] = value
		}

		return json.Marshal(data)
	case lang.TList:
		var (
			items = obj.Value().([]lang.Object)
			arr   []interface{}
		)

		for _, item := range items {
			val, err := j.convertToJSON(item)
			if err != nil {
				return nil, err
			}
			var value interface{}
			json.Unmarshal(val, &value)
			arr = append(arr, value)
		}

		return json.Marshal(arr)
	case lang.TNil, lang.TBool, lang.TInt, lang.TFloat, lang.TString:
		return json.Marshal(obj.Value())
	default:
		return nil, fmt.Errorf("unsupported type: %s", obj.Type())
	}
}

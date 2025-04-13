package modules

import (
	"encoding/json"
	"fmt"

	"github.com/bndrmrtn/flare/lang"
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
		"parse":    lang.NewFunction(j.parse).WithTypeSafeArgs(lang.TypeSafeArg{Name: "string", Type: lang.TString}),
		"toString": lang.NewFunction(j.toString).WithArg("object"),
	}
}

func (j *JSON) parse(args []lang.Object) (lang.Object, error) {
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
			key := fmt.Sprintf("%v", keyObj.Value())
			valObj := arr.Map[keyObj]
			val, err := j.convertToJSON(valObj)
			if err != nil {
				return nil, err
			}
			var value interface{}
			if err := json.Unmarshal(val, &value); err != nil {
				return nil, err
			}
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
			if err := json.Unmarshal(val, &value); err != nil {
				return nil, err
			}
			arr = append(arr, value)
		}

		return json.Marshal(arr)
	case lang.TNil, lang.TBool, lang.TInt, lang.TFloat, lang.TString:
		return json.Marshal(obj.Value())
	case lang.TInstance:
		method := obj.Method("value")
		if method != nil && len(method.Args()) == 0 {
			data, err := method.Execute(nil)
			if err != nil {
				return nil, err
			}
			return j.convertToJSON(data)
		}

		variableNames := obj.Variables()
		keys := make([]lang.Object, len(variableNames))
		vars := make([]lang.Object, len(variableNames))
		for i, name := range variableNames {
			keys[i] = lang.NewString("key", name, nil)
			vars[i] = obj.Variable(name)
		}

		return j.convertToJSON(lang.NewArray("variables", nil, keys, vars))
	}
	return nil, fmt.Errorf("unsupported type: %s", obj.Type())
}

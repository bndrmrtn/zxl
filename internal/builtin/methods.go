package builtin

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/internal/state"
	"github.com/bndrmrtn/flare/lang"
)

type ImportFunc func(file string, d *models.Debug) (lang.Object, error)
type EvalFunc func(code string) (lang.Object, error)

func GetMethods(importer ImportFunc, evaler EvalFunc, provider *state.Provider) map[string]lang.Method {
	m := make(map[string]lang.Method)

	m["print"] = &Print{}
	m["println"] = &Print{true}
	m["import"] = &Import{importer}
	m["type"] = lang.NewFunction(fnType).WithArg("object")
	m["range"] = lang.NewFunction(fnRange).WithVariadicArg("range")
	m["read"] = lang.NewFunction(fnRead).WithArg("text")
	m["string"] = lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
		return lang.NewString("string", args[0].String(), args[0].Debug()), nil
	}).WithArg("object")
	m["fail"] = lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
		return nil, errors.New(args[0].String())
	}).WithArg("message")
	m["eval"] = lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
		return evaler(args[0].String())
	}).WithTypeSafeArgs(lang.TypeSafeArg{Type: lang.TString, Name: "code"})
	m["ref"] = lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
		return lang.NewRef(args[0].Name(), args[0].Debug(), args[0]), nil
	}).WithArg("object")

	m["map"] = lang.NewFunction(fnMap).WithArgs([]string{"fn", "object"})
	m["fetch"] = lang.NewFunction(fnFetch).
		WithArg("url").WithVariadicArg("config")
	m["state"] = lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
		if provider == nil {
			panic("Fatal error: No provider found for state.")
		}

		name := args[0].Value().(string)
		state := provider.State(name)
		return NewState(name, args[0].Debug(), state), nil
	}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "name", Type: lang.TString})

	m = setTypeMethods(m)
	m = setIsMethods(m)
	return m
}

type Print struct {
	newLine bool
}

func (p *Print) Args() []string {
	return nil
}

func (p *Print) Execute(args []lang.Object) (lang.Object, error) {
	var sb strings.Builder
	argList, ok := args[0].Value().([]lang.Object)
	if !ok {
		return nil, fmt.Errorf("expected []lang.Object, got %T", args[0].Value())
	}

	for i, arg := range argList {
		sb.WriteString(arg.String())
		if i < len(argList)-1 {
			sb.WriteString(" ")
		}
	}

	if p.newLine {
		fmt.Println(sb.String())
	} else {
		fmt.Print(sb.String())
	}
	return nil, nil
}

func (p *Print) HasVariadicArg() bool {
	return true
}

func (p *Print) GetVariadicArg() string {
	return "message"
}

type Import struct {
	importer ImportFunc
}

func (*Import) Args() []string {
	return []string{"file"}
}

func (i *Import) Execute(args []lang.Object) (lang.Object, error) {
	obj := args[0]
	if obj.Type() != lang.TString {
		return nil, fmt.Errorf("expected string, got %v", obj.Type())
	}

	file := obj.Value().(string)
	return i.importer(file, obj.Debug())
}

// Functions

// fnType returns the type of the given object.
func fnType(args []lang.Object) (lang.Object, error) {
	obj := args[0]

	if ts, ok := obj.(interface{ TypeString() string }); ok {
		return lang.NewString("type", ts.TypeString(), obj.Debug()), nil
	}

	return lang.NewString("type", string(obj.Type()), obj.Debug()), nil
}

// fnRange returns a range object.
func fnRange(arguments []lang.Object) (lang.Object, error) {
	var start, stop, step int

	args, err := rangeArgs(arguments)
	if err != nil {
		return nil, err
	}

	switch len(args) {
	case 1:
		stop, _ = args[0].Value().(int)
		start, step = 0, 1
	case 2:
		start, _ = args[0].Value().(int)
		stop, _ = args[1].Value().(int)
		step = 1
	case 3:
		start, _ = args[0].Value().(int)
		stop, _ = args[1].Value().(int)
		step, _ = args[2].Value().(int)
	default:
		return nil, errors.New("invalid number of arguments")
	}

	if step == 0 {
		return nil, errors.New("step cannot be zero")
	}

	var result []lang.Object
	if step > 0 {
		for i := start; i < stop; i += step {
			result = append(result, lang.NewInteger("range", i, nil))
		}
	} else {
		for i := start; i > stop; i += step {
			result = append(result, lang.NewInteger("range", i, nil))
		}
	}

	return lang.NewList("range", result, nil), nil
}

// rangeArgs parses the arguments for the range function.
func rangeArgs(args []lang.Object) ([]lang.Object, error) {

	if args[0].Type() == lang.TInt {
		return args, nil
	}

	if args[0].Type() == lang.TList {
		li := args[0].Value().([]lang.Object)
		if len(li) > 3 || len(li) < 1 {
			return nil, fmt.Errorf("expected to have 1, 2 or 3 elements, got %v", len(li))
		}

		for _, arg := range li {
			if arg.Type() != lang.TInt {
				return nil, fmt.Errorf("expected values to be <Int>, got %v", arg.Type())
			}
		}

		return li, nil
	}

	return nil, fmt.Errorf("expected <Int> or <List> [start, stop?, step?], got %v", args[0].Type())
}

func fnRead(args []lang.Object) (lang.Object, error) {
	if len(args) != 1 {
		return nil, errors.New("expected 1 argument")
	}

	if args[0].Type() != lang.TString {
		return nil, fmt.Errorf("expected <String>, got %v", args[0].Type())
	}

	input := args[0].Value().(string)

	fmt.Print(input)
	var value string
	fmt.Scanln(&value)
	return lang.NewString("read", value, nil), nil
}

func fnMap(args []lang.Object) (lang.Object, error) {
	if args[0].Type() != lang.TFnRef {
		return nil, fmt.Errorf("expected %s, got %v", lang.TFnRef.String(), args[0].Type())
	}

	fn := args[0].Value().(*lang.Fn).Fn
	fnArgs := fn.Args()

	if len(fnArgs) != 1 {
		return nil, fmt.Errorf("expected %s, got %v", lang.TFnRef.String(), args[0].Type())
	}

	value := args[1]
	if value.Type() != lang.TList {
		return fn.Execute([]lang.Object{value})
	}

	list := value.Value().([]lang.Object)
	var (
		result []lang.Object = make([]lang.Object, len(list))
		err    error
	)

	for i, item := range list {
		cp := item.Copy()
		cp.Rename(fnArgs[0])

		result[i], err = fn.Execute([]lang.Object{cp})
		if err != nil {
			return nil, err
		}
	}

	return lang.NewList("map", result, args[1].Debug()), nil
}

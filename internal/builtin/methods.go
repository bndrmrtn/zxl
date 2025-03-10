package builtin

import (
	"errors"
	"fmt"

	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
)

type ImportFunc func(file string, d *models.Debug) (lang.Object, error)

func GetMethods(importer ImportFunc) map[string]lang.Method {
	return map[string]lang.Method{
		"print":   &Print{},
		"println": &Print{true},
		"import":  &Import{importer},
		"type":    lang.NewFunction([]string{"value"}, fnType, nil),
		"range":   lang.NewFunction([]string{"range"}, fnRange, nil),
		"read":    lang.NewFunction([]string{"text"}, fnRead, nil),
		"string": lang.NewFunction([]string{"object"}, func(args []lang.Object) (lang.Object, error) {
			return lang.NewString("string", args[0].String(), args[0].Debug()), nil
		}, nil),
		"fail": lang.NewFunction([]string{"message"}, func(args []lang.Object) (lang.Object, error) {
			return nil, errors.New(args[0].String())
		}, nil),
	}
}

type Print struct {
	newLine bool
}

func (p *Print) Args() []string {
	return []string{"value"}
}

func (p *Print) Execute(args []lang.Object) (lang.Object, error) {
	if p.newLine {
		fmt.Println(args[0].String())
	} else {
		fmt.Print(args[0].String())
	}
	return nil, nil
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
	if def, ok := obj.(*lang.Definition); ok {
		return lang.NewString("type", string(def.TypeString()), def.Debug()), nil
	}
	if inst, ok := obj.(*lang.Instance); ok {
		return lang.NewString("type", string(inst.TypeString()), inst.Debug()), nil
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
			return nil, fmt.Errorf("expected list to have 1, 2 or 3 elements, got %v", len(li))
		}

		for _, arg := range li {
			if arg.Type() != lang.TInt {
				return nil, fmt.Errorf("expected list values to be <Int>, got %v", arg.Type())
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

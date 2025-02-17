package builtin

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/lang"
	"github.com/bndrmrtn/zexlang/internal/models"
)

type ImportFunc func(file string, d *models.Debug) (lang.Object, error)

func GetMethods(importer ImportFunc) map[string]lang.Method {
	return map[string]lang.Method{
		"print":         &Print{},
		"println":       &Print{true},
		"type":          Type{},
		"import":        &Import{importer},
		"zx_langbridge": LangBridge{},
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

type LangBridge struct{}

func (l LangBridge) Args() []string {
	return []string{"value", "args"}
}

func (l LangBridge) Execute(args []lang.Object) (lang.Object, error) {
	if args[0].Type() != lang.TString {
		return nil, fmt.Errorf("expected string, got %v", args[0].Type())
	}

	if args[1].Type() != lang.TList {
		return nil, fmt.Errorf("expected list, got %v", args[1].Type())
	}

	fmt.Println("calling", args[0].Value())
	return nil, nil
}

type Type struct{}

func (t Type) Args() []string {
	return []string{"value"}
}

func (t Type) Execute(args []lang.Object) (lang.Object, error) {
	obj := args[0]
	if def, ok := obj.(*lang.Definition); ok {
		return lang.NewString("type", string(def.TypeString()), def.Debug()), nil
	}

	return lang.NewString("type", string(obj.Type()), obj.Debug()), nil
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

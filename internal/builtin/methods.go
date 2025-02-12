package builtin

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/lang"
)

var Methods = map[string]lang.Method{
	"print":         &Print{},
	"println":       &Print{true},
	"type":          Type{},
	"zx_langbridge": LangBridge{},
}

type Print struct {
	newLine bool
}

func (p *Print) Args() []string {
	return []string{"value"}
}

func (p *Print) Execute(args []lang.Object) (lang.Object, error) {
	if p.newLine {
		fmt.Println(args[0].Value())
	} else {
		fmt.Print(args[0].Value())
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
	return lang.NewString("type", string(args[0].Type()), nil), nil
}

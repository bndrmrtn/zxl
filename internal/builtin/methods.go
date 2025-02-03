package builtin

import (
	"fmt"

	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// GetBuiltins returns the built-in functions
func GetBuiltins() map[string]Function {
	return map[string]Function{
		"print":   print,
		"println": println,
		"printf":  printf,
		"type":    getType,
		"read":    read,
	}
}

func print(args []*Variable) (*FuncReturn, error) {
	var values []any

	for _, arg := range args {
		values = append(values, arg.Value)
	}

	n, err := fmt.Print(values...)
	return &FuncReturn{
		Type:  tokens.IntVariable,
		Value: n,
	}, err
}

func println(args []*Variable) (*FuncReturn, error) {
	var values []any

	for _, arg := range args {
		values = append(values, arg.Value)
	}

	n, err := fmt.Println(values...)
	return &FuncReturn{
		Type:  tokens.IntVariable,
		Value: n,
	}, err
}

func printf(args []*Variable) (*FuncReturn, error) {
	var (
		format string
		values []any
	)

	for i, arg := range args {
		if i == 0 {
			if arg.Type != tokens.StringVariable {
				return nil, fmt.Errorf("expected string, got %v", arg.Type)
			}
			format = arg.Value.(string)
		} else {
			values = append(values, arg.Value)
		}
	}

	n, err := fmt.Printf(format, values...)
	return &FuncReturn{
		Type:  tokens.IntVariable,
		Value: n,
	}, err
}

func getType(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 argument, got %d", len(args))
	}

	return &FuncReturn{
		Type:  tokens.StringVariable,
		Value: args[0].Type.String(),
	}, nil
}

func read(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("expected 1 arguments, got %d", len(args))
	}

	if args[0].Type != tokens.StringVariable {
		return nil, fmt.Errorf("expected string argument, got %v", args[0].Type)
	}

	var value string
	fmt.Print(args[0].Value)
	_, err := fmt.Scan(&value)
	return &FuncReturn{
		Type:  tokens.StringVariable,
		Value: value,
	}, err
}

package ztime

import (
	"fmt"
	"time"

	"github.com/bndrmrtn/zxl/internal/lang"
)

// TimeNamespace implements a namespace for time operations
type TimeNamespace struct{}

// New creates a new time namespace
func New() *TimeNamespace {
	return &TimeNamespace{}
}

// Namespace returns the name of the namespace
func (*TimeNamespace) Namespace() string {
	return "time"
}

// Objects returns the objects in the namespace
func (*TimeNamespace) Objects() map[string]lang.Object {
	return map[string]lang.Object{}
}

// Methods returns the methods in the namespace
func (*TimeNamespace) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"now": lang.NewFunction(fnNow),
		"from": lang.NewFunction(fnFrom).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "format", Type: lang.TString}, lang.TypeSafeArg{Name: "value", Type: lang.TString}),
		"since": lang.NewFunction(fnSince).WithArg("time"),
		"until": lang.NewFunction(fnUntil).WithArg("time"),
	}
}

func fnNow(args []lang.Object) (lang.Object, error) {
	fmt.Println("NEW", NewTime(time.Now()))
	return NewTime(time.Now()), nil
}

func fnFrom(args []lang.Object) (lang.Object, error) {
	format := args[0]
	value := args[1]

	goFormat, err := ZxTimeFormatToGo(format.Value().(string))
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(goFormat, value.Value().(string))
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}

	return NewTime(t), nil
}

func fnSince(args []lang.Object) (lang.Object, error) {
	t, ok := args[0].(*Time)
	if !ok {
		return nil, fmt.Errorf("argument must be a time.time object")
	}

	duration := time.Since(t.time)
	return lang.NewFloat("duration", float64(duration)/float64(time.Second), nil), nil
}

func fnUntil(args []lang.Object) (lang.Object, error) {
	t, ok := args[0].(*Time)
	if !ok {
		return nil, fmt.Errorf("argument must be a time.time object")
	}

	duration := time.Until(t.time)
	return lang.NewFloat("duration", float64(duration)/float64(time.Second), nil), nil
}

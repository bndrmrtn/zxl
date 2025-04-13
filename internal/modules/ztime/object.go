package ztime

import (
	"fmt"
	"time"

	"github.com/bndrmrtn/flare/lang"
)

// Time wraps the Go time.Time type
type Time struct {
	lang.Base
	time time.Time
}

// NewTime creates a new Time object
func NewTime(t time.Time) *Time {
	return &Time{
		Base: lang.NewBase("time", nil),
		time: t,
	}
}

// Type returns the object type
func (t *Time) Type() lang.ObjType {
	return lang.TInstance
}

func (t *Time) TypeString() string {
	return "time.time"
}

// Value returns the underlying value
func (t *Time) Value() any {
	return t
}

// Method returns the named method
func (t *Time) Method(name string) lang.Method {
	switch name {
	case "format":
		return lang.NewFunction(t.fnFormat).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "format", Type: lang.TString})
	case "unix":
		return lang.NewFunction(t.fnUnix)
	case "add":
		return lang.NewFunction(t.fnAdd).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "seconds", Type: lang.TInt})
	case "sub":
		return lang.NewFunction(t.fnSub).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "seconds", Type: lang.TInt})
	case "addDays":
		return lang.NewFunction(t.fnAddDays).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "days", Type: lang.TInt})
	case "addMonths":
		return lang.NewFunction(t.fnAddMonths).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "months", Type: lang.TInt})
	case "addYears":
		return lang.NewFunction(t.fnAddYears).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "years", Type: lang.TInt})
	case "year":
		return lang.NewFunction(t.fnYear)
	case "month":
		return lang.NewFunction(t.fnMonth)
	case "day":
		return lang.NewFunction(t.fnDay)
	case "hour":
		return lang.NewFunction(t.fnHour)
	case "minute":
		return lang.NewFunction(t.fnMinute)
	case "second":
		return lang.NewFunction(t.fnSecond)
	case "weekday":
		return lang.NewFunction(t.fnWeekday)
	case "isZero":
		return lang.NewFunction(t.fnIsZero)
	case "before":
		return lang.NewFunction(t.fnBefore).WithArg("other")
	case "after":
		return lang.NewFunction(t.fnAfter).WithArg("other")
	case "equal":
		return lang.NewFunction(t.fnEqual).WithArg("other")
	default:
		return nil
	}
}

// Methods returns the available method names
func (t *Time) Methods() []string {
	return []string{
		"format", "unix", "add", "sub", "addDays", "addMonths", "addYears",
		"year", "month", "day", "hour", "minute", "second", "weekday",
		"isZero", "before", "after", "equal",
	}
}

// Variable returns a variable by name
func (t *Time) Variable(variable string) lang.Object {
	switch variable {
	case "$addr":
		return lang.Addr(t)
	default:
		return nil
	}
}

// Variables returns the available variable names
func (t *Time) Variables() []string {
	return []string{"$addr"}
}

// SetVariable sets a variable value
func (t *Time) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

// String returns a string representation
func (t *Time) String() string {
	return fmt.Sprintf("<Time %s %s>", t.time.Format(time.RFC3339), lang.Addr(t))
}

// Copy returns a copy of the object
func (t *Time) Copy() lang.Object {
	return NewTime(t.time)
}

// Time object methods
func (t *Time) fnFormat(args []lang.Object) (lang.Object, error) {
	format, ok := args[0].(*lang.String)
	if !ok {
		return nil, fmt.Errorf("format argument must be a string")
	}

	goFormat, err := FlareTimeFormatToGo(format.Value().(string))
	if err != nil {
		return nil, err
	}

	formatted := t.time.Format(goFormat)
	return lang.NewString("formatted", formatted, nil), nil
}

func (t *Time) fnUnix(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("unix", int(t.time.Unix()), nil), nil
}

func (t *Time) fnAdd(args []lang.Object) (lang.Object, error) {
	seconds := args[0].Value().(int)
	newTime := t.time.Add(time.Duration(seconds) * time.Second)

	return NewTime(newTime), nil
}

func (t *Time) fnSub(args []lang.Object) (lang.Object, error) {
	seconds := args[0].Value().(int)
	newTime := t.time.Add(-time.Duration(seconds) * time.Second)

	return NewTime(newTime), nil
}

func (t *Time) fnAddDays(args []lang.Object) (lang.Object, error) {
	days := args[0].Value().(int)
	newTime := t.time.AddDate(0, 0, days)

	return NewTime(newTime), nil
}

func (t *Time) fnAddMonths(args []lang.Object) (lang.Object, error) {
	months := args[0].Value().(int)
	newTime := t.time.AddDate(0, months, 0)

	return NewTime(newTime), nil
}

func (t *Time) fnAddYears(args []lang.Object) (lang.Object, error) {
	years := args[0].Value().(int)
	newTime := t.time.AddDate(years, 0, 0)

	return NewTime(newTime), nil
}

func (t *Time) fnYear(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("year", t.time.Year(), nil), nil
}

func (t *Time) fnMonth(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("month", int(t.time.Month()), nil), nil
}

func (t *Time) fnDay(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("day", t.time.Day(), nil), nil
}

func (t *Time) fnHour(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("hour", t.time.Hour(), nil), nil
}

func (t *Time) fnMinute(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("minute", t.time.Minute(), nil), nil
}

func (t *Time) fnSecond(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("second", t.time.Second(), nil), nil
}

func (t *Time) fnWeekday(args []lang.Object) (lang.Object, error) {
	return lang.NewInteger("weekday", int(t.time.Weekday()), nil), nil
}

func (t *Time) fnIsZero(args []lang.Object) (lang.Object, error) {
	return lang.NewBool("isZero", t.time.IsZero(), nil), nil
}

func (t *Time) fnBefore(args []lang.Object) (lang.Object, error) {
	other, ok := args[0].(*Time)
	if !ok {
		return nil, fmt.Errorf("other argument must be a time.time object")
	}

	return lang.NewBool("before", t.time.Before(other.time), nil), nil
}

func (t *Time) fnAfter(args []lang.Object) (lang.Object, error) {
	other, ok := args[0].(*Time)
	if !ok {
		return nil, fmt.Errorf("other argument must be a time.time object")
	}

	return lang.NewBool("after", t.time.After(other.time), nil), nil
}

func (t *Time) fnEqual(args []lang.Object) (lang.Object, error) {
	other, ok := args[0].(*Time)
	if !ok {
		return nil, fmt.Errorf("other argument must be a time.time object")
	}

	return lang.NewBool("equal", t.time.Equal(other.time), nil), nil
}

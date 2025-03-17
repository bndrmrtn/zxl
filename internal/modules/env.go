package modules

import (
	"os"

	"github.com/bndrmrtn/zxl/lang"
	"github.com/joho/godotenv"
)

type Env struct {
	dotEnvLoaded bool
}

func NewEnv() *Env {
	err := godotenv.Load()
	env := &Env{
		dotEnvLoaded: err == nil,
	}
	return env
}

func (*Env) Namespace() string {
	return "env"
}

func (h *Env) Objects() map[string]lang.Object {
	return map[string]lang.Object{
		"dotenvLoaded": lang.Immute(lang.NewBool("dotenvLoaded", h.dotEnvLoaded, nil)),
	}
}

func (h *Env) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"get": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].Value().(string)
			return lang.NewString(key, key, nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString}),
		"set": lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].Value().(string)
			value := args[1].Value().(string)
			err := os.Setenv(key, value)
			return lang.NewBool("set", err == nil, nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString}, lang.TypeSafeArg{Name: "value", Type: lang.TString}),
	}
}

package modules

import (
	"math/rand"
	"time"

	"github.com/bndrmrtn/flare/lang"
)

type Rand struct {
	random *rand.Rand
}

func NewRandModule() *Rand {
	return &Rand{
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (*Rand) Namespace() string {
	return "rand"
}

func (r *Rand) Objects() map[string]lang.Object {
	return nil
}

func (r *Rand) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"int": lang.NewFunction(r.fnInt).
			WithTypeSafeArgs(lang.TypeSafeArg{Name: "min", Type: lang.TInt}, lang.TypeSafeArg{Name: "max", Type: lang.TInt}),
		"bool": lang.NewFunction(r.fnBool),
	}
}

func (r *Rand) fnInt(args []lang.Object) (lang.Object, error) {
	min := args[0].Value().(int)
	max := args[1].Value().(int)

	return lang.NewInteger("number", min+r.random.Intn(max-min), nil), nil
}

func (r *Rand) fnBool(args []lang.Object) (lang.Object, error) {
	var b bool
	if r.random.Intn(2) == 1 {
		b = true
	}
	return lang.NewBool("bool", b, nil), nil
}

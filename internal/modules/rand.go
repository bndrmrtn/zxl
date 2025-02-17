package modules

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bndrmrtn/zxl/internal/lang"
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
		"int":  lang.NewFunction([]string{"min", "max"}, r.fnInt, nil),
		"bool": lang.NewFunction(nil, r.fnBool, nil),
	}
}

func (r *Rand) fnInt(args []lang.Object) (lang.Object, error) {
	if args[0].Type() != lang.TInt || args[1].Type() != lang.TInt {
		return nil, fmt.Errorf("min and max must be integers")
	}

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

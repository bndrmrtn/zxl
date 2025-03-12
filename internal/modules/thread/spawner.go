package thread

import (
	"fmt"
	"sync"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type Spawner struct {
	lang.Base

	max int

	sem chan struct{}

	wg sync.WaitGroup
}

func NewSpawner(max int) *Spawner {
	return &Spawner{
		Base: lang.NewBase("spawner", nil),
		max:  max,
		sem:  make(chan struct{}, max),
	}
}

func (s *Spawner) Type() lang.ObjType {
	return lang.TDefinition
}

func (s *Spawner) Value() any {
	return s
}

func (s *Spawner) Method(name string) lang.Method {
	switch name {
	case "spawn":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			fn := args[0].Value().(lang.Method)
			if len(fn.Args()) != 0 {
				return nil, fmt.Errorf("spawn handler must have 0 arguments, got %d", len(fn.Args()))
			}

			s.wg.Add(1)
			go func(wg *sync.WaitGroup, sem chan struct{}) {
				defer wg.Done()
				defer func() { <-sem }()
				sem <- struct{}{}

				_, _ = fn.Execute(nil)
			}(&s.wg, s.sem)

			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "method", Type: lang.TFnRef})
	default:
		return nil
	}
}

func (s *Spawner) Methods() []string {
	return []string{"spawn"}
}

func (s *Spawner) Variable(variable string) lang.Object {
	return nil
}

func (s *Spawner) Variables() []string {
	return nil
}

func (s *Spawner) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (s *Spawner) String() string {
	return fmt.Sprintf("<Spawner %d>", lang.Addr(s))
}

func (s *Spawner) Copy() lang.Object {
	return s
}

package thread

import (
	"context"
	"fmt"
	"time"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type Portal struct {
	lang.Base

	id uint

	portal chan lang.Object
}

func NewPortal(id uint, portalBufferSize int) *Portal {
	return &Portal{
		Base:   lang.NewBase("portal", nil),
		id:     id,
		portal: make(chan lang.Object, portalBufferSize),
	}
}

func (p *Portal) Type() lang.ObjType {
	return lang.TInstance
}

func (p *Portal) TypeString() string {
	return "thread.portal"
}

func (p *Portal) Value() any {
	return p
}

func (p *Portal) Method(name string) lang.Method {
	switch name {
	case "send":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			p.portal <- args[0]
			return lang.NewBool("ok", true, args[0].Debug()), nil
		}).WithArg("message")
	case "receive":
		return lang.NewFunction(func(_ []lang.Object) (lang.Object, error) {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				select {
				case <-ctx.Done():
					return lang.NewBool("ok", false, nil), nil
				case msg := <-p.portal:
					return msg, nil
				}
			}
		})
	default:
		return nil
	}
}

func (p *Portal) Methods() []string {
	return []string{"send", "receive"}
}

func (p *Portal) Variable(variable string) lang.Object {
	return nil
}

func (p *Portal) Variables() []string {
	return nil
}

func (p *Portal) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (p *Portal) String() string {
	return fmt.Sprintf("<Portal %d>", p.id)
}

func (p *Portal) Copy() lang.Object {
	return p
}

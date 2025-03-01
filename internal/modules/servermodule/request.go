package servermodule

import (
	"fmt"
	"net/http"

	"github.com/bndrmrtn/zxl/internal/lang"
)

type Request struct {
	lang.Base

	r *http.Request
}

func NewRequest(r *http.Request) *Request {
	return &Request{
		Base: lang.NewBase("request", nil),
		r:    r,
	}
}

func (r *Request) Type() lang.ObjType {
	return lang.TDefinition
}

func (r *Request) Value() any {
	return r
}

func (r *Request) Method(name string) lang.Method {
	switch name {
	default:
		return nil
	}
}

func (r *Request) Methods() []string {
	return nil
}

func (r *Request) Variable(variable string) lang.Object {
	switch variable {
	default:
		return nil
	case "$addr":
		return lang.Addr(r)
	case "method":
		return lang.NewString("method", r.r.Method, nil)
	case "url":
		return lang.NewString("url", r.r.URL.String(), nil)
	case "remoteAddr":
		return lang.NewString("remoteAddr", r.r.RemoteAddr, nil)
	}
}

func (r *Request) Variables() []string {
	return []string{"$addr"}
}

func (r *Request) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("not implemented")
}

func (r *Request) String() string {
	return fmt.Sprintf("<Request %s %s>", r.r.Method, lang.Addr(r))
}

func (r *Request) Copy() lang.Object {
	return r
}

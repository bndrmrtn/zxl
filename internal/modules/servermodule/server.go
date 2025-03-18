package servermodule

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/bndrmrtn/zxl/lang"
)

type HttpServer struct {
	w       http.ResponseWriter
	r       *http.Request
	Body    bytes.Buffer
	Code    int
	Written bool
}

func New(w http.ResponseWriter, r *http.Request) *HttpServer {
	return &HttpServer{
		w:       w,
		r:       r,
		Code:    http.StatusOK,
		Written: false,
	}
}

func (*HttpServer) Namespace() string {
	return "server"
}

func (h *HttpServer) Objects() map[string]lang.Object {
	return map[string]lang.Object{
		"request": lang.Immute(NewRequest(h.r)),
		"header":  lang.Immute(lang.NewDefinitionInstance(lang.NewDefinition("server.header", "header", nil, nil, nil), newHeader(h.r.Header, h.w.Header()))),
	}
}

func (h *HttpServer) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"write":    lang.NewFunction(h.fnWrite).WithArg("data"),
		"status":   lang.NewFunction(h.fnStatus).WithTypeSafeArgs(lang.TypeSafeArg{Name: "code", Type: lang.TInt}),
		"json":     lang.NewFunction(h.fnContentType("json")),
		"html":     lang.NewFunction(h.fnContentType("html")),
		"text":     lang.NewFunction(h.fnContentType("text")),
		"redirect": lang.NewFunction(h.fnRedirect).WithTypeSafeArgs(lang.TypeSafeArg{Name: "url", Type: lang.TString}, lang.TypeSafeArg{Name: "code", Type: lang.TInt}),
		"sendFile": lang.NewFunction(h.fnSendFile).WithTypeSafeArgs(lang.TypeSafeArg{Name: "path", Type: lang.TString}),
	}
}

func (h *HttpServer) fnWrite(args []lang.Object) (lang.Object, error) {
	data := fmt.Sprint(args[0].Value())
	h.Body.WriteString(data)
	return nil, nil
}

func (h *HttpServer) fnStatus(args []lang.Object) (lang.Object, error) {
	code := args[0].Value().(int)
	if code < 100 || code > 599 {
		return nil, fmt.Errorf("status code must be between 100 and 599")
	}
	h.Code = code
	return nil, nil
}

func (h *HttpServer) fnContentType(typ string) func([]lang.Object) (lang.Object, error) {
	var contentType = "text/plain"
	switch typ {
	case "json":
		contentType = "application/json"
	case "html":
		contentType = "text/html"
	case "text":
		contentType = "text/plain"
	}
	return func(_ []lang.Object) (lang.Object, error) {
		h.w.Header().Set("Content-Type", contentType)
		return nil, nil
	}
}

func (h *HttpServer) fnRedirect(args []lang.Object) (lang.Object, error) {
	url := args[0].String()
	code := args[1].Value().(int)
	if code < 300 || code > 399 {
		return nil, fmt.Errorf("redirect status code must be between 300 and 399")
	}
	http.Redirect(h.w, h.r, url, code)
	h.Written = true
	return nil, nil
}

func (h *HttpServer) fnSendFile(args []lang.Object) (lang.Object, error) {
	path := args[0].String()
	http.ServeFile(h.w, h.r, path)
	h.Written = true
	return nil, nil
}

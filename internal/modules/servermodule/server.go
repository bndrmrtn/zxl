package servermodule

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/bndrmrtn/zxl/internal/lang"
	"github.com/bndrmrtn/zxl/internal/models"
)

type HttpServer struct {
	w http.ResponseWriter
	r *http.Request

	Body bytes.Buffer
	Code int
}

func New(w http.ResponseWriter, r *http.Request) *HttpServer {
	return &HttpServer{
		w:    w,
		r:    r,
		Code: http.StatusOK,
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
		"write":  lang.NewFunction(h.fnWrite).WithArg("data"),
		"status": lang.NewFunction(h.fnStatus).WithTypeSafeArgs(lang.TypeSafeArg{Name: "code", Type: lang.TInt}),
		"body":   lang.NewFunction(h.fnBody),
		"json":   lang.NewFunction(h.fnContentType("json")),
		"html":   lang.NewFunction(h.fnContentType("html")),
		"text":   lang.NewFunction(h.fnContentType("text")),
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

func (h *HttpServer) fnBody(_ []lang.Object) (lang.Object, error) {
	body, err := io.ReadAll(h.r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body content: '%v'", err)
	}

	return lang.NewString("body", string(body), nil), nil
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

// Header

type header struct {
	request  http.Header
	response http.Header
}

func newHeader(r http.Header, w http.Header) *header {
	return &header{
		request:  r,
		response: w,
	}
}

func (h *header) GetVariable(variable string) (lang.Object, error) {
	return nil, fmt.Errorf("variable '%s' not found on http.Header", variable)
}

func (h *header) AssignVariable(variable string, value lang.Object) error {
	return fmt.Errorf("cannot set variable '%s' on http.Header", variable)
}

func (h *header) GetMethod(name string) (lang.Method, error) {
	switch name {
	default:
		return nil, fmt.Errorf("method '%s' not found on http.Header", name)
	case "set":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0]
			h.response.Set(key.String(), args[1].String())
			return nil, nil
		}).WithTypeSafeArgs(
			lang.TypeSafeArg{Name: "key", Type: lang.TString},
			lang.TypeSafeArg{Name: "value", Type: lang.TString},
		), nil
	case "get":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].String()
			return lang.NewString(key, h.request.Get(key), nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString}), nil
	}
}

func (h *header) Execute(_ []*models.Node) (lang.Object, error) {
	return nil, nil
}

func (h *header) GetNew() lang.Executer {
	return newHeader(h.request, h.response)
}

func (h *header) Get(_ []*models.Node) (lang.Object, error) {
	return nil, nil
}

func (h *header) Variables() []string {
	return []string{}
}

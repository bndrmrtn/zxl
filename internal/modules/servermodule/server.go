package servermodule

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/lang"
)

type HttpServer struct {
	w       http.ResponseWriter
	r       *http.Request
	Body    bytes.Buffer
	Code    int
	Written bool

	Params lang.Object
}

func New(w http.ResponseWriter, r *http.Request) *HttpServer {
	params := r.Context().Value("__params__").(map[string]string)

	hs := &HttpServer{
		w:       w,
		r:       r,
		Code:    http.StatusOK,
		Written: false,
	}

	var (
		keys   []lang.Object
		values []lang.Object
	)
	for key, value := range params {
		keys = append(keys, lang.NewString("key", key, nil))
		values = append(values, lang.NewString("value", value, nil))
	}

	hs.Params = lang.NewArray("params", nil, keys, values)

	return hs
}

func (*HttpServer) Namespace() string {
	return "server"
}

func (h *HttpServer) Objects() map[string]lang.Object {
	return map[string]lang.Object{
		"request": lang.Immute(NewRequest(h.r)),
		"header":  lang.Immute(lang.NewDefinitionInstance(lang.NewDefinition("server.header", "header", nil, nil, nil), newHeader(h.r.Header, h.w.Header()))),
		"params":  lang.Immute(h.Params),
	}
}

func (h *HttpServer) Methods() map[string]lang.Method {
	return map[string]lang.Method{
		"write":     lang.NewFunction(h.fnWrite).WithArg("data"),
		"status":    lang.NewFunction(h.fnStatus).WithTypeSafeArgs(lang.TypeSafeArg{Name: "code", Type: lang.TInt}),
		"json":      lang.NewFunction(h.fnContentType("json")),
		"html":      lang.NewFunction(h.fnContentType("html")),
		"text":      lang.NewFunction(h.fnContentType("text")),
		"redirect":  lang.NewFunction(h.fnRedirect).WithTypeSafeArgs(lang.TypeSafeArg{Name: "url", Type: lang.TString}, lang.TypeSafeArg{Name: "code", Type: lang.TInt}),
		"sendFile":  lang.NewFunction(h.fnSendFile).WithTypeSafeArgs(lang.TypeSafeArg{Name: "path", Type: lang.TString}),
		"setCookie": lang.NewFunction(h.fnSetCookie).WithTypeSafeArgs(lang.TypeSafeArg{Name: "options", Type: lang.TArray}),
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

func (h *HttpServer) fnSetCookie(args []lang.Object) (lang.Object, error) {
	cookieConfig := args[0].(*lang.Array)

	var cookie http.Cookie

	if name, ok := cookieConfig.Access("name"); ok {
		if name.Type() != lang.TString {
			return nil, errs.WithDebug(fmt.Errorf("invalid type for cookie.name, expected string"), name.Debug())
		}
		cookie.Name = name.String()
	} else {
		return nil, errs.WithDebug(fmt.Errorf("cookie.name is required"), name.Debug())
	}

	if value, ok := cookieConfig.Access("value"); ok {
		if value.Type() != lang.TString {
			return nil, errs.WithDebug(fmt.Errorf("invalid type for cookie.value, expected string"), value.Debug())
		}
		cookie.Value = value.String()
	} else {
		return nil, errs.WithDebug(fmt.Errorf("cookie.value is required"), value.Debug())
	}

	if path, ok := cookieConfig.Access("path"); ok && path.Type() == lang.TString {
		cookie.Path = path.String()
	}

	if domain, ok := cookieConfig.Access("domain"); ok && domain.Type() == lang.TString {
		cookie.Domain = domain.String()
	}

	if secure, ok := cookieConfig.Access("secure"); ok && secure.Type() == lang.TBool {
		cookie.Secure = secure.Value().(bool)
	}

	if httpOnly, ok := cookieConfig.Access("httpOnly"); ok && httpOnly.Type() == lang.TBool {
		cookie.HttpOnly = httpOnly.Value().(bool)
	}

	if expires, ok := cookieConfig.Access("expires"); ok && expires.Type() == lang.TString {
		t, err := time.Parse(time.RFC1123, expires.String())
		if err != nil {
			return nil, errs.WithDebug(fmt.Errorf("invalid expires format, expected RFC1123: %w", err), expires.Debug())
		}
		cookie.Expires = t
	}

	http.SetCookie(h.w, &cookie)

	return lang.NewBool("setCookie", true, args[0].Debug()), nil
}

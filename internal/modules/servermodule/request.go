package servermodule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/bndrmrtn/flare/lang"
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
	return lang.TInstance
}

func (*Request) TypeString() string {
	return "server.request"
}

func (r *Request) Value() any {
	return r
}

func (r *Request) Method(name string) lang.Method {
	switch name {
	default:
		return nil
	case "param":
		return lang.NewFunction(r.fnParam).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString})
	case "cookie":
		return lang.NewFunction(r.fnCookie).WithTypeSafeArgs(lang.TypeSafeArg{Name: "name", Type: lang.TString})
	case "query":
		return lang.NewFunction(r.fnQuery).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString})
	case "isSecure":
		return lang.NewFunction(r.fnIsSecure)
	case "isXHR":
		return lang.NewFunction(r.fnIsXHR)
	case "accept":
		return lang.NewFunction(r.fnAccept).WithTypeSafeArgs(lang.TypeSafeArg{Name: "contentType", Type: lang.TString})
	case "parseForm":
		return lang.NewFunction(r.fnParseForm)
	case "parseMultipartForm":
		return lang.NewFunction(r.fnParseMultipartForm).WithTypeSafeArgs(lang.TypeSafeArg{Name: "maxMemory", Type: lang.TInt})
	case "body":
		return lang.NewFunction(r.fnBody)
	case "bodyJson":
		return lang.NewFunction(r.fnBodyJson)
	}
}

func (r *Request) fnParam(args []lang.Object) (lang.Object, error) {
	key := args[0].String()
	return lang.NewString(key, r.r.PathValue(key), nil), nil
}

func (r *Request) fnHeader(args []lang.Object) (lang.Object, error) {
	key := args[0].String()
	return lang.NewString(key, r.r.Header.Get(key), nil), nil
}

func (r *Request) fnCookie(args []lang.Object) (lang.Object, error) {
	name := args[0].String()
	cookie, err := r.r.Cookie(name)
	if err != nil {
		if err == http.ErrNoCookie {
			return lang.NewString(name, "", nil), nil
		}
		return nil, fmt.Errorf("could not get cookie '%s': %v", name, err)
	}
	return lang.NewString(name, cookie.Value, nil), nil
}

func (r *Request) fnQuery(args []lang.Object) (lang.Object, error) {
	key := args[0].String()
	values := r.r.URL.Query()[key]
	if len(values) == 0 {
		return lang.NewString(key, "", nil), nil
	}
	if len(values) == 1 {
		return lang.NewString(key, values[0], nil), nil
	}
	queryValues := make([]lang.Object, len(values))
	for i, value := range values {
		queryValues[i] = lang.NewString(fmt.Sprintf("%s[%d]", key, i), value, nil)
	}
	return lang.NewList("queryValues", queryValues, nil), nil
}

func (r *Request) fnIsSecure(_ []lang.Object) (lang.Object, error) {
	return lang.NewBool("isSecure", r.r.TLS != nil, nil), nil
}

func (r *Request) fnIsXHR(_ []lang.Object) (lang.Object, error) {
	requestedWith := r.r.Header.Get("X-Requested-With")
	return lang.NewBool("isXHR", strings.ToLower(requestedWith) == "xmlhttprequest", nil), nil
}

func (r *Request) fnAccept(args []lang.Object) (lang.Object, error) {
	contentType := args[0].String()
	accept := r.r.Header.Get("Accept")
	return lang.NewBool("accept", strings.Contains(accept, contentType), nil), nil
}

func (r *Request) fnParseForm(_ []lang.Object) (lang.Object, error) {
	if err := r.r.ParseForm(); err != nil {
		return nil, fmt.Errorf("could not parse form: '%v'", err)
	}

	form := make(map[string]lang.Object)
	for key, values := range r.r.Form {
		if len(values) == 1 {
			form[key] = lang.NewString(key, values[0], nil)
		} else {
			formValues := make([]lang.Object, len(values))
			for i, value := range values {
				formValues[i] = lang.NewString(fmt.Sprintf("%s[%d]", key, i), value, nil)
			}
			form[key] = lang.NewList("formValues", formValues, nil)
		}
	}

	return lang.NewArrayMap("form", nil, form), nil
}

func (r *Request) fnParseMultipartForm(args []lang.Object) (lang.Object, error) {
	maxMemory := args[0].Value().(int)
	if err := r.r.ParseMultipartForm(int64(maxMemory)); err != nil {
		return nil, fmt.Errorf("could not parse multipart form: '%v'", err)
	}

	form := make(map[string]lang.Object)
	for key, values := range r.r.MultipartForm.Value {
		if len(values) == 1 {
			form[key] = lang.NewString(key, values[0], nil)
		} else {
			formValues := make([]lang.Object, len(values))
			for i, value := range values {
				formValues[i] = lang.NewString(fmt.Sprintf("%s[%d]", key, i), value, nil)
			}
			form[key] = lang.NewList("value", formValues, nil)
		}
	}

	// Add files
	for key, fileHeaders := range r.r.MultipartForm.File {
		fileArray := make([]lang.Object, len(fileHeaders))
		for i, fileHeader := range fileHeaders {
			fileInfo := map[string]lang.Object{
				"filename": lang.NewString("filename", fileHeader.Filename, nil),
				"size":     lang.NewInteger("size", int(fileHeader.Size), nil),
				"header":   lang.NewString("header", fmt.Sprintf("%v", fileHeader.Header), nil),
			}
			fileArray[i] = lang.NewArrayMap("fileInfo", nil, fileInfo)
		}
		form[key] = lang.NewList("file", fileArray, nil)
	}

	return lang.NewArrayMap("form", nil, form), nil
}

func (r *Request) fnBody(_ []lang.Object) (lang.Object, error) {
	body, err := io.ReadAll(r.r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body content: '%v'", err)
	}
	return lang.NewString("body", string(body), nil), nil
}

func (r *Request) fnBodyJson(_ []lang.Object) (lang.Object, error) {
	body, err := io.ReadAll(r.r.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body content: '%v'", err)
	}

	var data interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("could not parse JSON body: '%v'", err)
	}

	return lang.FromValue(data)
}

func (r *Request) Methods() []string {
	return []string{
		"param",
		"header",
		"cookie",
		"query",
		"isSecure",
		"isXHR",
		"accept",
		"parseForm",
		"parseMultipartForm",
		"body",
		"bodyJson",
	}
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
	case "path":
		return lang.NewString("path", r.r.URL.Path, nil)
	case "host":
		return lang.NewString("host", r.r.Host, nil)
	case "scheme":
		if r.r.TLS != nil {
			return lang.NewString("scheme", "https", nil)
		}
		return lang.NewString("scheme", "http", nil)
	case "protocol":
		return lang.NewString("protocol", r.r.Proto, nil)
	case "contentType":
		return lang.NewString("contentType", r.r.Header.Get("Content-Type"), nil)
	case "accept":
		return lang.NewString("accept", r.r.Header.Get("Accept"), nil)
	case "contentLength":
		return lang.NewInteger("contentLength", int(r.r.ContentLength), nil)
	case "params":
		params := make(map[string]lang.Object)
		for name, values := range r.r.URL.Query() {
			if len(values) == 1 {
				params[name] = lang.NewString(name, values[0], nil)
			} else {
				valuesArray := make([]lang.Object, len(values))
				for i, v := range values {
					valuesArray[i] = lang.NewString(fmt.Sprintf("%s[%d]", name, i), v, nil)
				}
				params[name] = lang.NewList("values", valuesArray, nil)
			}
		}
		return lang.NewArrayMap("params", nil, params)
	case "headers":
		headers := make(map[string]lang.Object)
		for name, values := range r.r.Header {
			if len(values) == 1 {
				headers[name] = lang.NewString(name, values[0], nil)
			} else {
				valuesArray := make([]lang.Object, len(values))
				for i, v := range values {
					valuesArray[i] = lang.NewString(fmt.Sprintf("%s[%d]", name, i), v, nil)
				}
				headers[name] = lang.NewList("header", valuesArray, nil)
			}
		}
		return lang.NewArrayMap("headers", nil, headers)
	}
}

func (r *Request) Variables() []string {
	return []string{
		"$addr",
		"method",
		"url",
		"remoteAddr",
		"path",
		"host",
		"scheme",
		"protocol",
		"contentType",
		"accept",
		"contentLength",
		"params",
		"headers",
	}
}

func (r *Request) SetVariable(_ string, _ lang.Object) error {
	return fmt.Errorf("request variables are read-only")
}

func (r *Request) String() string {
	return fmt.Sprintf("<Request %s %s>", r.r.Method, lang.Addr(r))
}

func (r *Request) Copy() lang.Object {
	return r
}

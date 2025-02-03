package builtin

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/bndrmrtn/zexlang/internal/tokens"
)

// HttpModule is a module for handling HTTP requests
type HttpModule struct {
	w http.ResponseWriter
	r *http.Request

	// StatusCode is the HTTP status code
	StatusCode int
	// Body is the body of the request
	Body bytes.Buffer
}

// NewHttpModule creates a new instance of HttpModule
func NewHttpModule(w http.ResponseWriter, r *http.Request) *HttpModule {
	return &HttpModule{
		w:          w,
		r:          r,
		StatusCode: http.StatusOK,
	}
}

// Access checks for a variable in the HttpModule (not used in your code, so it returns an error)
func (hm *HttpModule) Access(variable string) (*Variable, error) {
	switch variable {
	case "method":
		return &Variable{
			Type:  tokens.StringVariable,
			Value: hm.r.Method,
		}, nil
	case "body":
		body, err := io.ReadAll(hm.r.Body)
		if err != nil {
			return nil, fmt.Errorf("could not read body content: %v", err)
		}
		return &Variable{
			Type:  tokens.StringVariable,
			Value: string(body),
		}, nil
	case "url":
		return &Variable{
			Type:  tokens.StringVariable,
			Value: hm.r.URL.String(),
		}, nil
	default:
		return nil, fmt.Errorf("variable %s not exists", variable)
	}
}

// Execute runs the function passed as `fn` on the HttpModule
func (hm *HttpModule) Execute(fn string, args []*Variable) (*FuncReturn, error) {
	switch fn {
	case "write":
		return hm.write(args, false)
	case "writeln":
		return hm.write(args, true)
	case "status":
		return hm.status(args)
	case "setHeader":
		return hm.setHeader(args)
	case "getHeader":
		return hm.getHeader(args)
	case "query":
		return hm.query(args)
	case "html", "json", "text":
		hm.w.Header().Set("Content-Type", hm.getContentType(fn))
		return nil, nil
	default:
		return nil, fmt.Errorf("function %s not exists in http module", fn)
	}
}

// write sends the specified argument to the ResponseWriter
func (hm *HttpModule) write(args []*Variable, nl bool) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("write function expects 1 argument, got %d", len(args))
	}

	value := fmt.Sprintf("%v", args[0].Value)
	if nl {
		value += "\n"
	}

	_, err := hm.Body.WriteString(value)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// status sets the HTTP status code for the response
func (hm *HttpModule) status(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("status function expects 1 argument, got %d", len(args))
	}

	statusCode, ok := args[0].Value.(int)
	if !ok {
		return nil, fmt.Errorf("status code must be an integer, got %T", args[0].Value)
	}

	hm.StatusCode = statusCode
	return nil, nil
}

// setHeader sets an HTTP header for the response
func (hm *HttpModule) setHeader(args []*Variable) (*FuncReturn, error) {
	if len(args) != 2 {
		return nil, fmt.Errorf("setHeader function expects 2 arguments, got %d", len(args))
	}

	headerName, ok := args[0].Value.(string)
	if !ok {
		return nil, fmt.Errorf("header name must be a string, got %T", args[0].Value)
	}

	headerValue, ok := args[1].Value.(string)
	if !ok {
		return nil, fmt.Errorf("header value must be a string, got %T", args[1].Value)
	}

	hm.w.Header().Set(headerName, headerValue)
	return nil, nil
}

// getHeader returns the value of an HTTP header
func (hm *HttpModule) getHeader(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("getHeader function expects 1 arguments, got %d", len(args))
	}

	headerName, ok := args[0].Value.(string)
	if !ok {
		return nil, fmt.Errorf("header name must be a string, got %T", args[0].Value)
	}

	return &FuncReturn{
		Type:  tokens.StringVariable,
		Value: hm.w.Header().Get(headerName),
	}, nil
}

func (hm *HttpModule) getContentType(fn string) string {
	switch fn {
	case "html":
		return "text/html"
	case "json":
		return "application/json"
	case "text":
		return "text/plain"
	default:
		return "text/plain"
	}
}

func (hm *HttpModule) query(args []*Variable) (*FuncReturn, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("query function expects 1 argument, got %d", len(args))
	}

	queryKey, ok := args[0].Value.(string)
	if !ok {
		return nil, fmt.Errorf("query key must be a string, got %T", args[0].Value)
	}

	return &FuncReturn{
		Type:  tokens.StringVariable,
		Value: hm.r.URL.Query().Get(queryKey),
	}, nil
}

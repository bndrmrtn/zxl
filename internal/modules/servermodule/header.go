package servermodule

import (
	"fmt"
	"net/http"

	"github.com/bndrmrtn/zxl/internal/models"
	"github.com/bndrmrtn/zxl/lang"
)

// Header
type Header struct {
	request  http.Header
	response http.Header
}

func newHeader(r http.Header, w http.Header) *Header {
	return &Header{
		request:  r,
		response: w,
	}
}

func (h *Header) GetVariable(variable string) (lang.Object, error) {
	return nil, fmt.Errorf("variable '%s' not found on server.header", variable)
}

func (h *Header) AssignVariable(variable string, value lang.Object) error {
	return fmt.Errorf("cannot set variable '%s' on server.header", variable)
}

func (h *Header) GetMethod(name string) (lang.Method, error) {
	switch name {
	default:
		return nil, fmt.Errorf("method '%s' not found on server.header", name)
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
	case "add":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].String()
			h.response.Add(key, args[1].String())
			return nil, nil
		}).WithTypeSafeArgs(
			lang.TypeSafeArg{Name: "key", Type: lang.TString},
			lang.TypeSafeArg{Name: "value", Type: lang.TString},
		), nil
	case "del":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].String()
			h.response.Del(key)
			return nil, nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString}), nil
	case "values":
		return lang.NewFunction(func(args []lang.Object) (lang.Object, error) {
			key := args[0].String()
			values := h.request.Values(key)
			valuesArray := make([]lang.Object, len(values))
			for i, v := range values {
				valuesArray[i] = lang.NewString(fmt.Sprintf("%s[%d]", key, i), v, nil)
			}
			return lang.NewList("values", valuesArray, nil), nil
		}).WithTypeSafeArgs(lang.TypeSafeArg{Name: "key", Type: lang.TString}), nil
	}
}

func (h *Header) Execute(_ []*models.Node) (lang.Object, error) {
	return nil, nil
}

func (h *Header) GetNew() lang.Executer {
	return newHeader(h.request, h.response)
}

func (h *Header) Get(_ []*models.Node) (lang.Object, error) {
	return nil, nil
}

func (h *Header) Variables() []string {
	return []string{}
}

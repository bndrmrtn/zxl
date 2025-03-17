package modules

import (
	"net/http"

	"github.com/bndrmrtn/zxl/lang"
)

type Http struct{}

func NewHttpModule() *Http {
	return &Http{}
}

func (*Http) Namespace() string {
	return "http"
}

func (h *Http) Objects() map[string]lang.Object {
	m := h.getStatusMap()
	return m
}

func (h *Http) Methods() map[string]lang.Method {
	return nil
}

func (h *Http) getStatusMap() map[string]lang.Object {
	return map[string]lang.Object{
		"statusOK":                            lang.NewInteger("statusOk", http.StatusOK, nil),
		"statusCreated":                       lang.NewInteger("statusCreated", http.StatusCreated, nil),
		"statusAccepted":                      lang.NewInteger("statusAccepted", http.StatusAccepted, nil),
		"statusNoContent":                     lang.NewInteger("statusNoContent", http.StatusNoContent, nil),
		"statusResetContent":                  lang.NewInteger("statusResetContent", http.StatusResetContent, nil),
		"statusPartialContent":                lang.NewInteger("statusPartialContent", http.StatusPartialContent, nil),
		"statusMultiStatus":                   lang.NewInteger("statusMultiStatus", http.StatusMultiStatus, nil),
		"statusAlreadyReported":               lang.NewInteger("statusAlreadyReported", http.StatusAlreadyReported, nil),
		"statusIMUsed":                        lang.NewInteger("statusIMUsed", http.StatusIMUsed, nil),
		"statusMultipleChoices":               lang.NewInteger("statusMultipleChoices", http.StatusMultipleChoices, nil),
		"statusMovedPermanently":              lang.NewInteger("statusMovedPermanently", http.StatusMovedPermanently, nil),
		"statusFound":                         lang.NewInteger("statusFound", http.StatusFound, nil),
		"statusSeeOther":                      lang.NewInteger("statusSeeOther", http.StatusSeeOther, nil),
		"statusNotModified":                   lang.NewInteger("statusNotModified", http.StatusNotModified, nil),
		"statusUseProxy":                      lang.NewInteger("statusUseProxy", http.StatusUseProxy, nil),
		"statusTemporaryRedirect":             lang.NewInteger("statusTemporaryRedirect", http.StatusTemporaryRedirect, nil),
		"statusPermanentRedirect":             lang.NewInteger("statusPermanentRedirect", http.StatusPermanentRedirect, nil),
		"statusBadRequest":                    lang.NewInteger("statusBadRequest", http.StatusBadRequest, nil),
		"statusUnauthorized":                  lang.NewInteger("statusUnauthorized", http.StatusUnauthorized, nil),
		"statusPaymentRequired":               lang.NewInteger("statusPaymentRequired", http.StatusPaymentRequired, nil),
		"statusForbidden":                     lang.NewInteger("statusForbidden", http.StatusForbidden, nil),
		"statusNotFound":                      lang.NewInteger("statusNotFound", http.StatusNotFound, nil),
		"statusMethodNotAllowed":              lang.NewInteger("statusMethodNotAllowed", http.StatusMethodNotAllowed, nil),
		"statusNotAcceptable":                 lang.NewInteger("statusNotAcceptable", http.StatusNotAcceptable, nil),
		"statusRequestTimeout":                lang.NewInteger("statusRequestTimeout", http.StatusRequestTimeout, nil),
		"statusConflict":                      lang.NewInteger("statusConflict", http.StatusConflict, nil),
		"statusGone":                          lang.NewInteger("statusGone", http.StatusGone, nil),
		"statusLengthRequired":                lang.NewInteger("statusLengthRequired", http.StatusLengthRequired, nil),
		"statusPreconditionFailed":            lang.NewInteger("statusPreconditionFailed", http.StatusPreconditionFailed, nil),
		"statusUnsupportedMediaType":          lang.NewInteger("statusUnsupportedMediaType", http.StatusUnsupportedMediaType, nil),
		"statusExpectationFailed":             lang.NewInteger("statusExpectationFailed", http.StatusExpectationFailed, nil),
		"statusMisdirectedRequest":            lang.NewInteger("statusMisdirectedRequest", http.StatusMisdirectedRequest, nil),
		"statusUnprocessableEntity":           lang.NewInteger("statusUnprocessableEntity", http.StatusUnprocessableEntity, nil),
		"statusLocked":                        lang.NewInteger("statusLocked", http.StatusLocked, nil),
		"statusFailedDependency":              lang.NewInteger("statusFailedDependency", http.StatusFailedDependency, nil),
		"statusUpgradeRequired":               lang.NewInteger("statusUpgradeRequired", http.StatusUpgradeRequired, nil),
		"statusPreconditionRequired":          lang.NewInteger("statusPreconditionRequired", http.StatusPreconditionRequired, nil),
		"statusTooManyRequests":               lang.NewInteger("statusTooManyRequests", http.StatusTooManyRequests, nil),
		"statusRequestHeaderFieldsTooLarge":   lang.NewInteger("statusRequestHeaderFieldsTooLarge", http.StatusRequestHeaderFieldsTooLarge, nil),
		"statusUnavailableForLegalReasons":    lang.NewInteger("statusUnavailableForLegalReasons", http.StatusUnavailableForLegalReasons, nil),
		"statusInternalServerError":           lang.NewInteger("statusInternalServerError", http.StatusInternalServerError, nil),
		"statusNotImplemented":                lang.NewInteger("statusNotImplemented", http.StatusNotImplemented, nil),
		"statusBadGateway":                    lang.NewInteger("statusBadGateway", http.StatusBadGateway, nil),
		"statusServiceUnavailable":            lang.NewInteger("statusServiceUnavailable", http.StatusServiceUnavailable, nil),
		"statusGatewayTimeout":                lang.NewInteger("statusGatewayTimeout", http.StatusGatewayTimeout, nil),
		"statusHTTPVersionNotSupported":       lang.NewInteger("statusHTTPVersionNotSupported", http.StatusHTTPVersionNotSupported, nil),
		"statusVariantAlsoNegotiates":         lang.NewInteger("statusVariantAlsoNegotiates", http.StatusVariantAlsoNegotiates, nil),
		"statusInsufficientStorage":           lang.NewInteger("statusInsufficientStorage", http.StatusInsufficientStorage, nil),
		"statusLoopDetected":                  lang.NewInteger("statusLoopDetected", http.StatusLoopDetected, nil),
		"statusNotExtended":                   lang.NewInteger("statusNotExtended", http.StatusNotExtended, nil),
		"statusNetworkAuthenticationRequired": lang.NewInteger("statusNetworkAuthenticationRequired", http.StatusNetworkAuthenticationRequired, nil),
	}
}

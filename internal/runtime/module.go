package runtime

import (
	"net/http"

	"github.com/bndrmrtn/zexlang/internal/builtin"
)

// NewHttpModule creates a new instance of HttpModule
func NewHttpModule(w http.ResponseWriter, r *http.Request) *builtin.HttpModule {
	return builtin.NewHttpModule(w, r)
}

// NewHtmlModule creates a new instance of HttpModule
func NewHtmlModule(hm *builtin.HttpModule) *builtin.HtmlModule {
	return builtin.NewHtmlModule(hm)
}

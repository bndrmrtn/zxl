package server

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bndrmrtn/flare/internal/errs"
	"github.com/bndrmrtn/flare/internal/models"
	"github.com/bndrmrtn/flare/internal/modules/servermodule"
	"github.com/bndrmrtn/flare/internal/runtimev2"
	"github.com/bndrmrtn/flare/lang"
)

var errNotFound = errors.New("not found")

// handleError handles the error
func (s *Server) handleError(err error, w http.ResponseWriter, r *http.Request) {
	var statusCode = http.StatusInternalServerError

	if errors.Is(err, errNotFound) {
		statusCode = http.StatusNotFound
	}

	var de errs.DebugError

	if errors.As(err, &de) {
		if s.handleCustomErrorHandler(de.GetParentError(), statusCode, w, r) {
			return
		}

		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(statusCode)
		htmlErr := de.HttpError()
		if htmlErr == nil {
			_, _ = w.Write([]byte(de.Error()))
			return
		}
		_, _ = w.Write([]byte(s.makePrettyCode(htmlErr)))
		return
	}

	if s.handleCustomErrorHandler(err, statusCode, w, r) {
		return
	}

	http.Error(w, err.Error(), statusCode)
}

func (s *Server) handleCustomErrorHandler(serverErr error, code int, w http.ResponseWriter, r *http.Request) bool {
	path := filepath.Join(s.root, "err.fl")

	// Execute the cached nodes if they exist
	if s.useCaching {
		if nodes, ok := s.getCache(path); ok {
			return s.executeErrorHandler(nodes, serverErr, code, w, r)
		}
	}

	// Open the file
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Get the nodes
	nodes, err := s.ir.GetNodes(path, file)
	if err != nil {
		return false
	}

	// Cache the nodes for faster execution
	if s.useCaching {
		s.setCache(path, nodes)
	}

	// Execute the error handler
	return s.executeErrorHandler(nodes, serverErr, code, w, r)
}

// executeErrorHandler executes the error handler nodes.
func (s *Server) executeErrorHandler(nodes []*models.Node, serverErr error, code int, w http.ResponseWriter, r *http.Request) bool {
	// Execute the nodes
	run, err := runtimev2.New(s.serverStateProvider)
	if err != nil {
		return false
	}

	httpModule := servermodule.New(w, r)
	httpModule.Code = code
	run.BindModule(httpModule)

	// Execute the nodes
	if _, err := run.Execute(nodes); err != nil {
		return false
	}

	ex, err := run.GetNamespaceExecuter("err")
	if err != nil {
		return false
	}

	handler, err := ex.GetMethod("handle")
	if err != nil {
		return false
	}

	if len(handler.Args()) != 2 {
		return false
	}

	if _, err := handler.Execute([]lang.Object{
		lang.NewInteger("code", code, nil),
		lang.NewString("message", serverErr.Error(), nil),
	}); err != nil {
		return false
	}

	// Write the response
	w.WriteHeader(httpModule.Code)
	_, _ = w.Write(httpModule.Body.Bytes())
	return true
}

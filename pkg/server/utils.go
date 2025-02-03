package server

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/bndrmrtn/zexlang/internal/errs"
)

// getExecutablePath gets the executable path
func (s *Server) getExecutablePath(path string) (string, error) {
	if filepath.Ext(path) != ".zx" {
		stat, err := os.Stat(path + ".zx")
		if err == nil && !stat.IsDir() {
			return path + ".zx", nil
		}
	}

	stat, err := os.Stat(path)
	if err == nil && !stat.IsDir() {
		return path, nil
	}

	path = filepath.Join(path, "index.zx")
	stat, err = os.Stat(path)
	if err == nil && !stat.IsDir() {
		return path, nil
	}

	return "", fmt.Errorf("file not found")
}

// handleError handles the error
func (s *Server) handleError(err error, w http.ResponseWriter) {
	var de errs.DebugError
	if errors.As(err, &de) {
		w.Header().Add("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(de.HttpError()))
		return
	}

	http.Error(w, err.Error(), http.StatusInternalServerError)
}
